# 直播后台面试深挖点——具体到每一行代码、每一个数字

> 面试官最烦的回答："我们用了Redis做缓存"。他想听的是："Redis Cluster 256分片，单分片10万QPS，热点房间打散用{room_id}%256做slot分配，有一次线上大V开播导致单分片热点过载，我是这么解的——"

---

## 一、WS长连接网关：面试官一定会追问的细节

### 1.1 "单机5万连接，这个数字怎么来的？"

```
内存计算：
  每个连接的内存消耗：
  ├── TCP socket buffer：读4KB + 写4KB = 8KB（调小了默认的128KB）
  ├── goroutine栈：8KB（Go默认，实际用不满）
  ├── 应用层buffer（Room订阅列表+消息队列）：~4KB
  ├── 连接元数据（user_id, room_id, 时间戳等）：~1KB
  └── 合计：~21KB/连接

  5万连接 × 21KB = 1.05GB，16核32GB机器内存完全够用

CPU瓶颈才是真正限制：
  每个连接每秒下推10条消息 × 200B = 2KB/s/连接
  5万连接 × 2KB/s = 100MB/s 出口带宽
  单机千兆网卡出口带宽 = 125MB/s → 网卡利用率80%，到极限了

  实际瓶颈不是连接数，是出口带宽。
  所以我们用的是双万兆网卡（25Gbps bonding），
  这时候CPU成为瓶颈——大量消息序列化+系统调用。

  最终极限：16核机器，epoll + goroutine pool，
  实测极限是6万连接@10msg/s，留20%余量取5万。
```

### 1.2 "epoll具体怎么用的？水平触发还是边缘触发？"

```go
// 实际用的是Go的netpoll（底层自动用epoll ET模式）
// 但如果面试官追问，要能说清楚区别：

// 水平触发(LT)：只要fd可读，epoll_wait就返回 → 简单但可能多次唤醒
// 边缘触发(ET)：只在状态变化时返回一次 → 高效但必须一次读完

// Go的net库用的是ET+non-blocking IO，每个goroutine阻塞在Read上
// 底层由runtime的netpoll统一管理所有fd的epoll事件

// 我们的WS网关实际实现（关键代码）：
func (s *Server) handleConn(conn *websocket.Conn, uid, roomID int64) {
    // 注册到房间
    s.roomMgr.Join(roomID, uid, conn)
    defer s.roomMgr.Leave(roomID, uid)

    // 只处理心跳pong，不处理业务上行
    conn.SetReadDeadline(time.Now().Add(30 * time.Second))
    conn.SetPongHandler(func(string) error {
        conn.SetReadDeadline(time.Now().Add(30 * time.Second))
        return nil
    })

    // 读循环只为检测断连（不处理业务消息！）
    for {
        _, _, err := conn.ReadMessage()
        if err != nil {
            break // 连接断开
        }
    }
}
```

### 1.3 "心跳机制具体怎么实现的？30秒超时怎么定的？"

```
为什么是30秒？
├── 太短（如5秒）：移动端锁屏/后台时频繁误判断连
├── 太长（如120秒）：真正断连的用户长时间占用资源
├── 30秒是经验值：覆盖大部分运营商NAT超时（通常60-120秒）
│   每30秒ping一次，保证在NAT超时前刷新映射

具体实现：
├── 服务端每15秒发ping帧（WebSocket control frame）
├── 客户端收到ping后回pong
├── 服务端30秒没收到任何数据（含pong）→ 判定断连
├── 客户端侧同理：15秒没收到任何数据 → 主动重连

为什么不让客户端主动ping？
├── 如果100万客户端同时ping → 服务端瞬间涌入100万消息
├── 服务端主动ping可以错峰发送（每秒分批ping一部分连接）
└── 我们的实现：用时间轮(TimeWheel)管理心跳，
    将5万连接分成15个桶，每秒处理一个桶（约3333个连接）

时间轮具体实现：
slots: [15][]ConnID  // 15个槽位，对应15秒一轮
tick: 每秒前进一格
当前slot的连接：发送ping + 检查上一轮是否有pong响应
```

### 1.4 "房间消息广播具体怎么做的？50万人的房间怎么推？"

```
关键问题：一条弹幕要推给50万人，如果逐个写socket → 单机要写50万次syscall

实际方案——分层扇出：

1. 业务服务发送消息 → MQ(topic: room_broadcast)
2. Job路由层消费MQ，查询"room_id=xxx的连接分布在哪些网关节点"
   └── Redis Set: room:gateways:{room_id} → {gw1, gw2, gw3, ...}
   └── 50万人分布在约10-15台网关上
3. Job层向这10-15台网关发gRPC请求（并行，每台一条）
4. 每台网关收到后，在本机内存中找到该room的所有连接
   └── 本机约3-5万该房间的连接
5. 网关内部广播：

   // 关键优化：消息只序列化一次，发给所有连接
   func (r *Room) Broadcast(msg []byte) {
       // msg已经是序列化好的WebSocket frame
       for _, conn := range r.conns {
           // 非阻塞写入每个连接的发送队列
           select {
           case conn.sendCh <- msg:
           default:
               // 队列满了 → 这个连接太慢，标记待踢出
               conn.markSlow()
           }
       }
   }

6. 每个连接有独立的写goroutine，从sendCh批量取出消息合并写入socket

性能数据：
├── 一条消息从进MQ到所有观众收到：P99 < 200ms
├── 网关单机广播5万连接耗时：~10ms（非阻塞入队）
├── 真正的耗时在网络传输（TCP write + 网卡发包）
└── 瓶颈：网卡带宽，不是CPU
```

### 1.5 "如果某台网关机器挂了，5万连接怎么办？"

```
实际发生过的线上事故：一台网关OOM导致进程crash

处理过程：
1. 5万连接瞬间断开 → 客户端检测到连接断开
2. 客户端自动重连（SDK内置指数退避重连）
   ├── 第一次：立即重连（大部分在1秒内回来）
   ├── 第二次：1秒后重连
   ├── 第三次：2秒后重连
   └── 最多重试5次，之后提示用户"网络异常"
3. 重连时通过负载均衡分配到其他网关（不会再分到挂掉的那台）
4. 重连建立后，新网关注册连接 + 上报路由变更

对用户的影响：
├── 1-3秒无法收到弹幕/系统消息
├── 音视频不受影响（走CDN独立链路！）
└── 自动恢复，不需要用户手动操作

事后优化：
├── 增加了内存水位告警（>80%开始限流新连接）
├── 网关实现了graceful shutdown（先通知客户端迁移，再关进程）
└── 单房间连接打散到至少3台网关（避免单点影响过大）
```

---

## 二、Redis热点问题：面试最爱问的"线上故障"

### 2.1 "大V开播瞬间100万人涌入，Redis扛得住吗？"

```
实际发生过的场景：某明星直播，10秒内涌入80万人

问题暴露：
  进房请求 → SADD room:viewers:{room_id} user_id
  80万次SADD打到同一个Redis key → 这个key所在的slot过载
  单slot承载的分片QPS飙升到50万（正常上限10万）→ 延迟从1ms→50ms→超时

为什么Redis Cluster也扛不住？
├── Redis Cluster是按slot分片的，一个key只在一个分片上
├── 超热key（如大V的room_id） → 流量全打到单分片 → 分片过载
├── 其他分片很闲 → 集群总容量没到上限，但体验已崩

我的解决方案——热点key打散：

方案A（已上线）：大房间不用Set，改用计数器+抽样列表
  // 大房间进房（>1万人时自动切换策略）
  func JoinRoom(roomID, userID int64) {
      if isLargeRoom(roomID) {
          // 计数器分散到16个子key，最后求和
          subKey := fmt.Sprintf("room:count:%d:%d", roomID, userID%16)
          redis.Incr(subKey)
          // 在线列表只保留最近100人（环形buffer）
          redis.LPush(fmt.Sprintf("room:recent:%d", roomID), userID)
          redis.LTrim(fmt.Sprintf("room:recent:%d", roomID), 0, 99)
      } else {
          redis.SAdd(fmt.Sprintf("room:viewers:%d", roomID), userID)
      }
  }

方案B（评估中）：本地缓存+异步聚合
  // 网关本地计数，每5秒批量上报Redis
  localCounter[roomID]++
  // 定时任务每5秒：
  redis.IncrBy("room:count:{roomID}", localDelta)
  localDelta = 0
  // 代价：在线数有5秒延迟，但大房间这个延迟完全可接受
```

### 2.2 "Redis和MySQL数据不一致怎么办？"

```
场景：主播开播，流状态写Redis成功，写MySQL失败（网络抖动）

我们的一致性策略——先Redis后MySQL + 异步对账：

  func StartStream(streamID int64) error {
      // 1. 先写Redis（高频读写走Redis）
      err := redis.HSet("stream:info:"+streamID, "status", "live")
      if err != nil {
          return err  // Redis写失败 → 直接失败，不往下走
      }

      // 2. 异步写MySQL（MQ保证最终一致）
      msg := StreamStartedEvent{StreamID: streamID, Status: "live"}
      mq.Send("topic_stream_sync", msg)  // 异步

      return nil  // 对外返回成功
  }

  // 消费端：MySQL写入 + 重试
  func handleStreamSync(msg StreamStartedEvent) {
      for retry := 0; retry < 3; retry++ {
          err := db.Exec("UPDATE stream SET status=? WHERE id=?", msg.Status, msg.StreamID)
          if err == nil { return }
          time.Sleep(time.Second * time.Duration(retry+1))
      }
      // 3次失败 → 告警 + 人工介入
      alertOnCall("stream_sync_failed", msg)
  }

对账机制（兜底）：
  // 每5分钟跑一次对账Job
  func reconcile() {
      // 从Redis拿所有"live"状态的流
      liveStreams := redis.Keys("stream:info:*") // 实际用scan
      for _, s := range liveStreams {
          redisStatus := redis.HGet(s, "status")
          mysqlStatus := db.Query("SELECT status FROM stream WHERE id=?", id)
          if redisStatus != mysqlStatus {
              // 以Redis为准修复MySQL（Redis是热路径的source of truth）
              db.Exec("UPDATE stream SET status=? WHERE id=?", redisStatus, id)
              metric.Inc("reconcile_fix_count")
          }
      }
  }
```

### 2.3 "TRTC事件回调乱序怎么处理？"

```
真实场景：
  TRTC回调按HTTP投递，不保证顺序。可能收到：
  10:00:01 - 事件103（主播开始推流）
  10:00:03 - 事件104（主播退出房间）  ← 网络抖动导致的误报
  10:00:02 - 事件103（主播开始推流）  ← 重连成功的回调，但晚到了

  如果按收到顺序处理：pushing → live → closing → live ???

解决方案——状态机 + 事件时间戳 + 幂等：

  func handleTRTCCallback(event TRTCEvent) error {
      streamInfo := redis.HGetAll("stream:info:" + event.StreamID)

      // 1. 幂等检查：用事件序列号去重
      lastSeq := streamInfo["last_event_seq"]
      if event.Sequence <= lastSeq {
          return nil  // 旧事件，忽略
      }

      // 2. 时间戳检查：拒绝过时事件
      lastEventTime := streamInfo["last_event_time"]
      if event.Timestamp < lastEventTime {
          log.Warn("out-of-order event", event)
          return nil  // 乱序事件，忽略
      }

      // 3. 状态机合法性检查
      currentStatus := streamInfo["status"]
      if !isValidTransition(currentStatus, event.Type) {
          log.Warn("invalid transition", currentStatus, event.Type)
          return nil  // 非法转换，忽略
      }

      // 4. 执行状态转换
      newStatus := getNextStatus(currentStatus, event.Type)
      redis.HMSet("stream:info:"+event.StreamID, map[string]interface{}{
          "status":          newStatus,
          "last_event_seq":  event.Sequence,
          "last_event_time": event.Timestamp,
      })

      return nil
  }

  // 合法状态转换表
  var validTransitions = map[string]map[string]string{
      "idle":    {"103": "pushing"},        // 进房通知
      "pushing": {"103": "live", "104": "closing"}, // 推流成功 or 进房失败退出
      "live":    {"104": "closing"},        // 退房/断流
      "closing": {"103": "live"},           // 断线重连成功
  }
```

---

## 三、连麦实现：面试官最爱追问的"具体怎么做"

### 3.1 "连麦信令交互的具体流程？每一步的延迟是多少？"

```
完整的连麦时序（从观众点击"申请连麦"到画面出现）：

T+0ms     观众B点击"申请连麦"
            │ POST /api/v1/mic/apply {room_id, seat_index}
            ▼
T+50ms    连麦信令服务处理
            │ 1. 检查席位是否空闲（Redis HGET）  ~1ms
            │ 2. 检查用户是否有连麦资格           ~1ms
            │ 3. 设置席位状态为"邀请中"           ~1ms
            │ 4. 发MQ通知主播                    ~5ms
            ▼
T+100ms   主播端收到连麦申请通知（WS推送）
            │ 弹窗："用户B申请连麦，是否同意？"
            ▼
T+3000ms  主播点击"同意"（假设3秒后操作）
            │ POST /api/v1/mic/accept {room_id, user_id, seat_index}
            ▼
T+3050ms  连麦信令服务处理
            │ 1. 更新席位状态："邀请中"→"已上麦"      ~1ms
            │ 2. 为观众B生成TRTC UserSig          ~5ms
            │ 3. 通知观众B进TRTC房间（WS推送）      ~50ms
            │ 4. 发MQ广播席位变更                  ~5ms
            ▼
T+3100ms  观众B客户端收到"上麦"指令
            │ 1. TRTC SDK切换角色：Audience → Anchor   ~100ms
            │ 2. 开启摄像头+麦克风采集               ~200ms
            │ 3. 开始推流到TRTC房间                  ~100ms
            ▼
T+3500ms  主播端看到观众B的画面（TRTC房间内<300ms延迟）
            │
            ▼
T+3600ms  云端混流更新布局（TRTC MCU自动触发）
            │ 将主播A+连麦者B画面合成一路
            │ 旁路推流到CDN更新
            ▼
T+5000ms  普通观众通过CDN看到连麦画面（CDN延迟1-2秒）

总结：
├── 主播同意到连麦者出画面：~500ms（TRTC房间内）
├── 普通观众看到：再多1-2秒（CDN延迟）
└── 用户全程感知：点击后约3秒看到（主要是等主播操作）
```

### 3.2 "云端混流的布局怎么配置的？具体参数？"

```
TRTC混流API调用（实际代码）：

func updateMixLayout(roomID int64, seats []MicSeat) error {
    // 计算混流布局
    layout := buildLayout(seats)

    // 调用TRTC REST API设置混流
    req := &trtc.StartMCUMixTranscodeRequest{
        SdkAppId: sdkAppId,
        RoomId:   roomID,
        OutputParams: &trtc.OutputParams{
            StreamId: fmt.Sprintf("mix_%d", roomID),  // CDN流ID
            PureAudioStream: 0,
        },
        EncodeParams: &trtc.EncodeParams{
            VideoWidth:  1080,
            VideoHeight: 1920,  // 竖屏直播
            VideoBitrate: 3000, // kbps
            VideoFramerate: 25,
            AudioSampleRate: 48000,
            AudioBitrate: 64,
            AudioChannels: 2,
        },
        LayoutParams: layout,
    }

    _, err := trtcClient.StartMCUMixTranscode(req)
    return err
}

// 布局计算——画中画模式
func buildLayout(seats []MicSeat) *trtc.LayoutParams {
    params := &trtc.LayoutParams{
        Template: 0,  // 自定义布局
    }

    for i, seat := range seats {
        if seat.Status != StatusOnMic { continue }

        var region trtc.MixLayoutInfo
        if i == 0 {
            // 主播：全屏底层
            region = trtc.MixLayoutInfo{
                UserId: seat.TRTCUserId,
                X: 0, Y: 0, Width: 1080, Height: 1920,
                ZOrder: 1,
            }
        } else {
            // 连麦者：右下角小窗（按序排列）
            region = trtc.MixLayoutInfo{
                UserId: seat.TRTCUserId,
                X: 1080 - 270,  // 右侧
                Y: 200 + (i-1)*380,  // 纵向排列
                Width: 240, Height: 360,
                ZOrder: 2,
            }
        }
        params.MixLayoutList = append(params.MixLayoutList, &region)
    }

    return params
}
```

### 3.3 "连麦者网络差导致画面卡顿，怎么处理？"

```
实际遇到过的问题：连麦者用4G网络，上行带宽不足

TRTC SDK侧的自适应策略（SDK内置，我们配置参数）：

// 设置编码参数的降级策略
trtcCloud.setVideoEncoderParam(TRTCVideoEncParam{
    videoResolution: .resolution_640_360,  // 连麦者不需要1080P
    videoFps: 15,        // 帧率也可以低一些
    videoBitrate: 800,   // 800kbps
    minVideoBitrate: 200, // 最低降到200kbps（关键！）
    enableAdjustRes: true, // 允许动态降分辨率
})

// SDK根据网络状况自动调整：
// 网络好：640x360 @ 15fps @ 800kbps
// 网络一般：480x270 @ 15fps @ 500kbps
// 网络差：320x180 @ 10fps @ 200kbps
// 网络极差：纯音频模式（视频关闭）

业务侧的降级策略（我们自己实现的）：

func onNetworkQualityChanged(userID string, quality int) {
    switch {
    case quality >= 4:  // 网络极差
        // 1. 通知该用户关闭视频，只保留音频
        pushToUser(userID, MsgType_CloseVideo, nil)
        // 2. 通知主播："连麦者网络不佳"
        pushToAnchor(roomID, MsgType_GuestNetworkPoor, userID)
        // 3. 如果持续30秒不恢复，自动下麦
        scheduleKickIfNotRecover(userID, 30*time.Second)
    case quality == 3:  // 网络一般
        // 提示但不处理
        pushToUser(userID, MsgType_NetworkWarning, nil)
    }
}
```

---

## 四、首屏秒开：面试官最爱问"具体优化了哪些点"

### 4.1 "你们首屏从多少优化到多少？具体做了什么？"

```
优化前：P95 = 2.8秒
优化后：P95 = 680ms

逐项优化及其贡献：

优化1：HTTPDNS + IP直连（节省150ms）
  before: 系统DNS解析 → 可能走运营商LocalDNS → 劫持/慢
  after:  App启动时预拉取IP列表，拉流直接用IP
  代码：
    // 预热DNS
    func preResolveCDN() {
        ips := httpdns.Resolve("cdn.example.com")
        cache.Set("cdn_ips", ips, 5*time.Minute)
    }
    // 拉流时直接用IP
    url := fmt.Sprintf("http://%s/live/%s.flv", cachedIP, streamID)
    req.Header.Set("Host", "cdn.example.com")  // Host头保持域名

优化2：TCP连接预建（节省200ms）
  before: 点击进入直播间 → 新建TCP连接 → TLS握手 → 请求
  after:  推荐列表可见时就预建连接
  代码：
    // 列表item可见时
    func onLiveCardVisible(streamURL string) {
        // 预建TCP+TLS连接，放入连接池
        go preConnect(streamURL)
    }

优化3：GOP缓存（节省0-2000ms，平均节省800ms）——最大贡献
  CDN边缘节点配置：
    # Nginx-RTMP模块配置
    application live {
        gop_cache on;        # 开启GOP缓存
        gop_cache_count 1;   # 缓存最近1个GOP
    }
  效果：观众进入时立即收到最近的I帧，不用等下一个GOP

优化4：并行请求（节省200ms）
  before: 串行 → 获取房间信息 → 拿到播放地址 → 开始拉流
  after:  并行 → 房间信息 + WS建连 + 拉流同时发起
  客户端代码：
    func enterRoom(roomID) {
        // 三个请求并行
        go fetchRoomSnapshot(roomID)    // HTTP
        go connectWebSocket(roomID)      // WS
        go startPullStream(roomID)       // 拉流（用预缓存的地址）
    }
  关键：拉流地址从"进入时请求"改为"列表页预下发"

优化5：解码器预热（节省50ms）
  before: 拉到数据才初始化解码器
  after:  进入直播间页面时就初始化H.264硬解码器
  iOS代码：
    // viewWillAppear时
    let session = VTDecompressionSession(...)  // 预创建解码会话
```

### 4.2 "首屏监控怎么做的？怎么定义首屏时间？"

```
首屏时间定义：
  T_start = 用户点击直播间封面的时间戳
  T_end   = 播放器渲染第一帧视频的时间戳（onFirstFrameRendered回调）
  首屏时间 = T_end - T_start

上报方式：
  客户端SDK在onFirstFrameRendered回调中计算并上报：
  {
      "event": "first_frame",
      "duration_ms": 680,
      "room_id": "xxx",
      "network_type": "4G",
      "cdn_node": "sz-ct-01",  // 深圳电信节点1
      "codec": "H264",
      "resolution": "720P",
      "breakdown": {
          "dns_ms": 0,       // 用了IP直连
          "tcp_ms": 0,       // 用了预连接
          "first_byte_ms": 120,  // 首包到达
          "gop_cache_ms": 350,   // GOP缓存数据接收
          "decode_ms": 80,       // 解码
          "render_ms": 16        // 渲染
      }
  }

监控大盘：按CDN节点/运营商/网络类型/地域 分维度看P50/P95/P99
  发现异常：某地区P95突然从800ms飙到3s
  → 排查CDN节点 → 发现该节点GOP缓存被关了 → 联系CDN恢复
```

---

## 五、在线人数：一个"看起来简单"的功能

### 5.1 "进退房的计数为什么会不准？怎么解决的？"

```
实际遇到过的坑——计数漂移：

问题：在线人数只增不减，运行一天后比实际多了20%
原因：
  ├── 用户App闪退 → 没有正常退房 → WS心跳超时判断断连
  ├── 但心跳超时的回调偶尔丢失（网关OOM重启）
  ├── 导致DECR没有执行 → 计数偏高
  └── 积累一天后偏差越来越大

解决方案——定时校准：

  // 方案1：小房间（<1万人）用Set，天然精确
  // Set的SCARD就是精确值，不存在漂移问题

  // 方案2：大房间定时校准
  func calibrateOnlineCount(roomID int64) {
      // 每30秒执行一次
      // 统计当前WS网关上该房间的实际连接数
      actualCount := 0
      for _, gw := range getGatewaysForRoom(roomID) {
          count := gw.GetRoomConnCount(roomID)  // gRPC查询
          actualCount += count
      }

      // 用实际值覆盖Redis计数器
      redis.Set(fmt.Sprintf("room:online_count:%d", roomID), actualCount)
  }

  // 这个校准有成本（要查询所有网关），所以只对大房间每30秒做一次
  // 小房间用Set就没这个问题

方案3（最终上线版本）——三重保障：
  1. 正常路径：进房INCR + 退房DECR（实时）
  2. 兜底路径：WS心跳超时触发DECR（30秒延迟）
  3. 校准路径：每30秒比对网关实际连接数修正（最终一致）
```

### 5.2 "假如面试官追问：为什么不直接用网关连接数当在线人数？"

```
原因：网关连接 ≠ 在线观众

反例：
├── 一个用户可能有多个连接（多端登录、重连闪断）
├── 连接建立了但还没完成进房（鉴权阶段）
├── 旧连接还没超时断开，新连接已经建立（短暂双计）
├── WebSocket只是信令通道，不代表用户在看直播
│   （用户可能只是停留在直播间但最小化了App）
└── 某些用户走CDN拉流但不建WS连接（精简模式/省电模式）

正确做法：
  在线人数 = 业务层管理（进房/退房API驱动）
  网关连接数 = 基础设施指标（用于监控和容量规划）
  两者独立统计，定期校准
```

---

## 六、线上故障案例——"你遇到过最棘手的问题"

### 6.1 "一次大型活动的雪崩事件"

```
背景：某明星直播，预告20:00开播，粉丝涌入

时间线：
19:55 - 粉丝开始进入直播间（房间状态还是idle，等待开播）
        进房QPS从5万→20万→50万，还在涨
19:58 - 进房QPS达到80万，房间服务开始出现超时
        原因：每次进房都要查Redis + 发MQ + 更新计数
        MQ积压从0→50万条
19:59 - Redis单分片（该房间所在slot）QPS达到40万
        延迟从1ms→20ms→200ms
        触发上游服务超时→重试→雪崩
20:00 - 主播开播，但进房接口已经超时，新观众进不来
        大量用户反馈"进不去"，微博热搜预警

紧急处理（5分钟内）：
1. [20:00] 开启进房限流：令牌桶限制10万QPS入口
2. [20:01] 关闭非核心逻辑：
   - 关闭欢迎消息广播（大房间本来就不需要每人进来都广播）
   - 关闭实时在线列表更新（只更新计数器）
   - MQ消息降级为批量聚合（每秒聚合一次而非逐条）
3. [20:02] Redis分片扩容：将热点room_id的计数器打散到16个子key
4. [20:03] 系统恢复正常，QPS稳定在15万

事后复盘改进：
├── 预案1：大型活动提前开启"高峰模式"（自动关闭重型逻辑）
├── 预案2：进房改为"乐观进房"——先返回成功，后台异步处理
├── 预案3：热点房间的Redis key提前打散（活动开始前手动触发）
├── 预案4：进房接口增加本地缓存（房间信息5秒内不重复查Redis）
└── 预案5：MQ消费增加批量聚合能力（100条消息合并为1条广播）
```

### 6.2 "音画不同步的线上bug"

```
现象：约0.5%的观众反馈"声音比画面快了2-3秒"

排查过程：
1. 初步怀疑：播放器bug → 替换播放器版本后仍有问题
2. 抓包分析：CDN推下来的FLV流本身音视频时间戳就对不齐
3. 定位源头：云端混流环节
   └── 连麦场景下，主播的音频和连麦者的音频时钟不同步
   └── MCU混音时以主播时钟为基准，但混合后的时间戳计算有偏移

根因：
  TRTC云端混流在处理"纯音频流"和"音视频流"合并时，
  如果连麦者只开了麦克风（没开摄像头），
  MCU输出的混流中音频PTS基准不统一

解决：
  1. 临时方案：强制要求连麦者同时开启摄像头（哪怕用黑屏占位）
  2. 正式方案：升级TRTC SDK版本（腾讯云修复了MCU的PTS对齐逻辑）
  3. 兜底方案：播放器端增加音画同步检测
     if abs(audioPTS - videoPTS) > 200ms {
         // 主动丢弃偏移过大的音频帧，等待重新对齐
         resync()
     }

教训：
  混流场景的时钟同步是一个非常隐蔽的问题，
  只在特定条件下触发（纯音频+音视频混合），
  测试环境很难复现（因为网络条件不同）。
  最终靠线上抓包+客户端PTS打印日志定位。
```

### 6.3 "CDN节点故障导致部分地区黑屏"

```
现象：
  河南联通用户大面积反馈"黑屏"，其他地区正常

排查（3分钟内定位）：
  1. 查监控大盘：拉流成功率按地域+运营商分组
     → 河南联通拉流成功率从99.9%掉到60%
  2. 查CDN节点状态：河南联通边缘节点zz-cu-03报错
     → 节点磁盘满了，无法写入新的缓存 → 所有新请求回源 → 源站过载
  3. 该节点承载了河南联通80%的直播流量

处理：
  T+0min: 发现告警
  T+1min: 确认是单节点故障
  T+2min: 调度系统摘除该节点（新请求不再分配到此节点）
  T+3min: 已连接用户的播放器检测到卡顿 → 自动重连 → 被调度到其他节点
  T+5min: 全面恢复

播放器侧的自愈逻辑：
  func onPlayError(err error) {
      // 1. 换CDN节点重试
      newURL := requestNewPlayURL(roomID)  // 重新请求调度
      player.play(newURL)

      // 2. 如果3次重试都失败
      if retryCount >= 3 {
          // 降级到HLS（另一套CDN链路，更稳定但延迟高）
          player.play(hlsURL)
          showToast("网络不佳，已切换到流畅模式")
      }
  }
```

---

## 七、性能优化实战

### 7.1 "MQ消费积压怎么处理的？"

```
场景：晚高峰MQ消息量从200万/s涨到500万/s，消费跟不上

积压表现：
  topic_room_signal消费延迟从10ms涨到30秒
  → 观众进房后30秒才收到欢迎消息 → 体验崩溃

根因分析：
  Job路由层消费MQ后，要gRPC调用WS网关推送
  每条消息一次gRPC调用 → 500万次/s gRPC调用 → RTT开销太大

优化方案——批量+合并：

  // 优化前：逐条推送
  func consume(msg Message) {
      gw := findGateway(msg.RoomID)
      gw.Push(msg)  // 每条消息一次RPC
  }

  // 优化后：批量聚合推送
  var batchBuffer = make(map[string][]Message)  // key=gateway_addr
  var batchTicker = time.NewTicker(50 * time.Millisecond)

  func consume(msg Message) {
      gwAddr := findGatewayAddr(msg.RoomID)
      batchBuffer[gwAddr] = append(batchBuffer[gwAddr], msg)
  }

  func flushBatch() {
      for gwAddr, msgs := range batchBuffer {
          // 一次gRPC调用推送一批消息（最多500条）
          gw.BatchPush(msgs)
      }
      batchBuffer = make(map[string][]Message)
  }

  // 每50ms或buffer满500条时flush
  // 500万条/s ÷ 500条/批 = 1万次gRPC/s（减少500倍！）

效果：
  消费延迟从30秒恢复到50ms
  CPU利用率从90%降到40%
  代价：消息推送延迟增加50ms（从"几乎实时"变为"最多50ms延迟"，完全可接受）
```

### 7.2 "拉流调度的具体实现？QPS怎么扛500万？"

```
拉流调度做的事：给用户返回最优的CDN播放地址
  输入：user_ip, room_id, protocol(flv/hls), quality(720p/1080p)
  输出：最优CDN节点的播放URL

为什么QPS这么高？
  5000万在线观众 × (10%切换清晰度 + 5%断线重连) / 每分钟 ≈ 500万QPS

设计要点——轻+快+无状态：

  func getPlayURL(req PlayRequest) string {
      // 1. 解析用户IP → 运营商+地域（本地MaxMind GeoIP库，0.01ms）
      geo := geoIP.Lookup(req.UserIP)

      // 2. 查本地缓存的CDN节点健康状态（每5秒从中控同步）
      nodes := localCache.GetHealthyNodes(geo.ISP, geo.Province)

      // 3. 加权随机选择（权重=节点质量分数）
      node := weightedRandom(nodes)

      // 4. 拼接URL + 防盗链签名
      url := fmt.Sprintf("http://%s/live/%s_%s.flv?auth=%s",
          node.IP, req.StreamID, req.Quality,
          genAuthToken(req.StreamID, time.Now().Unix()))

      return url
  }

性能关键：
├── 全内存计算，零IO（GeoIP库+节点列表都在本地内存）
├── 单次请求耗时 < 0.5ms
├── 单机QPS > 10万（纯CPU计算，8核机器）
├── 500万QPS ÷ 10万/机 = 50台，冗余取100台
└── 无状态，可以任意水平扩展

防盗链签名计算：
  // 防止别人盗用我们的CDN带宽
  token = md5(secret + streamID + expireTimestamp + clientIP)
  // CDN边缘节点验签：收到请求后用同样算法验证，不匹配则403
```

---

## 八、面试"杀手级"回答模板

### 当面试官问"你做了什么优化"

**错误回答**：
> "我们用了Redis做缓存，用了MQ做异步，用了CDN做分发。"

**正确回答**：
> "举个具体的例子——首屏优化。我们把P95从2.8秒优化到680ms。最大的贡献来自GOP缓存，它解决了'观众进入时必须等待下一个I帧'的问题，单项贡献约800ms。但光有GOP缓存不够，我们还做了HTTPDNS预解析省150ms、TCP连接预建省200ms、并行请求省200ms。这些组合在一起才达到680ms。其中我踩过一个坑——预连接在4G切WiFi时会导致连接失效但池里还认为有效，所以我们加了连接健康检查机制，每次取连接前验证一下是否alive。"

### 当面试官问"遇到过什么困难"

**错误回答**：
> "我们系统有一次挂了，然后我们修了。"

**正确回答**：
> "印象最深的是一次80万人同时进房导致Redis单分片过载。问题的根因是Redis Cluster按slot分片，同一个room_id的所有操作打到同一分片。我们的解决方案是分层：小房间仍用Set精确管理，大房间切换为分散计数器（room:count:{room_id}:{user_id%16}，打散到16个不同slot）+ 环形buffer近似列表。关键决策点在于'什么时候判定为大房间'——我们用了渐进式切换：当Set大小超过1万时，下次进房自动切换策略，而不是一开始就用近似方案，因为小房间精确列表有业务价值（展示头像墙）。"

### 当面试官问"系统怎么保证高可用"

**错误回答**：
> "我们做了主从、做了集群、做了监控。"

**正确回答**：
> "以推流高可用为例，我们做了三层防线。SDK层：TRTC SDK内置断线重连，30%丢包仍可维持推流。服务层：流状态机有'closing'状态+30秒断流保护期，避免网络抖动误判关播。CDN层：双CDN同时推流，A挂了B无缝接管。有一次真实案例：主播WiFi断了切4G，TRTC SDK 2秒内自动重连成功，观众端看到的是画面卡了1秒然后恢复——全程无人工干预。我们监控大盘上看到的就是一个流状态从live→closing→live的跳变，持续时间2秒。"
