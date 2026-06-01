# WeDo 车队推荐系统技术实现文档


---

## 需求背景

### 业务场景

WeDo 是一款面向微信/QQ 生态的游戏社交平台。**车队（Convoy）** 是核心社交单元——玩家可以创建或加入车队，与队友协同组队参与 PvP 对局。

推荐系统的目标是在玩家进入"发现车队"频道时，从大量可用车队中找出最适合该用户加入的队伍，提升玩家的社交活跃度与匹配质量。

### 核心挑战

| 挑战 | 描述 |
|------|------|
| **社交图谱复杂** | 好友关系分为游戏好友（type=1）、平台好友（type=2）、双端好友（type=3），亲密度、组队时间需实时更新 |
| **数据时效性** | 车队信息（成员进出、解散、改名）变化频繁，需近实时同步 |
| **多维度排序** | 推荐结果需综合好友数量、亲密度、最近组队时间等多维信号 |
| **请求低延迟** | 推荐接口 P99 需控制在 2000ms 以内 |
| **多场景支持** | v1 支持通用分页浏览，v2 支持以"某好友"为锚点的定向推荐 |

### 技术选型背景

- **Flink**：处理游戏服务器上报的 Kafka 流水，构建实时好友关系链和车队索引
- **Elasticsearch**：存储结构化车队信息，支持快速按 EnvId/PublishID 多条件查询
- **Redis (TendisPlus)**：缓存好友关系哈希表、用户推荐结果；支持高并发读写
- **tRPC + RecDAG**：腾讯内部微服务框架，DAG 方式编排推荐算法，无需改代码即可通过修改 TOML 配置调整算法流程
- **Go**：推荐服务用 Go 实现，充分利用 goroutine 并发能力

---

## 系统整体架构

```
┌─────────────────────────────────────────────────────────────────────────┐
│                          游戏客户端 / 业务服务                             │
└──────────────────────────────┬──────────────────────────────────────────┘
                               │ HTTP (tRPC)
                               ▼
┌──────────────────────────────────────────────────────────────────────────┐
│                    推荐服务 (Go / tRPC)  :30000                           │
│                                                                          │
│   ┌──────────────┐     ┌──────────────────────────────────┐              │
│   │ RequestDecoder│───▶│         DAG Graph Engine         │              │
│   └──────────────┘     │  (RecDAG / graph.toml 驱动)       │              │
│                        │                                  │              │
│   ┌──────────────┐     │  configstore → Redis → ES →      │              │
│   │ResponseComposer◀───│  merge → setAttr → cache         │              │
│   └──────────────┘     └──────────────────────────────────┘              │
│                                    │ │ │                                 │
└────────────────────────────────────┼─┼─┼─────────────────────────────────┘
                                     │ │ │
          ┌──────────────────────────┘ │ └─────────────────────────┐
          ▼                            ▼                           ▼
┌──────────────────┐       ┌──────────────────┐        ┌─────────────────┐
│  ConfigStore     │       │  Redis (TendisPlus│        │ Elasticsearch  │
│  (配置中心)       │       │  推荐/SNS集群)     │        │ (车队索引)       │
│  reason_config   │       │  - 好友关系链      │        │ wedo_convoy_env │
└──────────────────┘       │  - 用户推荐缓存    │        │ wedo_video_rec  │
                           │  - 车队标签       │         └─────────────────┘
                           │  - 版本元数据      │                   ▲
                           └──────────────────┘                   │
                                     ▲                            │
                                     │                            │
┌────────────────────────────────────┴────────────────────────────┴──────┐
│                    Flink 实时数据处理层 (Java)                           │
│                                                                        │
│  ┌──────────────────────┐         ┌─────────────────┐                  │
│  │  FriendsOperatorAsync│         │  ConvoyInfoToES │                  │
│  │  (好友关系链维护)      │         │  (车队信息同步)   │                  │
│  └──────────────────────┘         └─────────────────┘                  │
│                 ▲                           ▲                          │
└─────────────────┼───────────────────────────┼──────────────────────────┘
                  │                    │
┌─────────────────┴────────────────────┴─────────────────────────────────┐
│                              Kafka 消息队列                             │
│  billow_ex_wedo_playerfriendslistflow (全量好友关系)                     │
│  billow_ex_wedo_friendflow            (好友增量变更)                     │
│  billow_ex_wedo_modifyconvoyflow      (车队变更流水)                     │
│  billow_ex_wedo_pvpendflow            (PvP 对局结束)                    │
└────────────────────────────────────────────────────────────────────────┘
                  ▲
┌─────────────────┴──────────────────────────────────────────────────────┐
│                         游戏服务器 / TLog 上报                           │
│              (wedo_dbfile_tlog.xml / wedo_tdw_tlog.xml 定义字段规范)     │
└────────────────────────────────────────────────────────────────────────┘
```

---

## 数据流转链路

### 数据写入链路（Flink → 存储）

```
游戏事件
  │
  ├─► [好友关系变更]
  │     Kafka: wedo_playerfriendslistflow (全量登录)
  │           + wedo_friendflow (增量变更)
  │     Flink: FriendsOperatorAsync
  │       ├── 全量登录时：HGet 历史 + 合并 + HMSet 写 Redis
  │       └── 增量变更时：HSet/HDel 单条更新 Redis
  │     Redis Key: wedo_car_friendlist04_{env}_{envid}_{roleid}
  │     Value: {friendRuid} → {亲密度}|{组队时间}|{好友类型}|{亲密度类型}
  │
  ├─► [车队信息变更]
  │     Kafka: wedo_modifyconvoyflow
  │     Flink: ConvoyInfoToES
  │       ├── KeyBy(ConvoyId) 聚合 5s 窗口内的多次变更
  │       └── 写 ES: wedo_convoy_{env}
  │     ES 字段: ConvoyId, ConvoyName, MemberList, MemberNum,
  │              HeadRuid, MemberPublishConvoyList, CreateTime,
  │              JoinLimitType, ConvoyTagList, ...
  │
  ├─► [PvP 对局数据]
  │     Kafka: wedo_pvpendflow + wedo_uploadvideoflow
  │     Flink: PvpToES
  │     ES: wedo_video_rec_{env}
```

### 数据读取链路（推荐服务）

```
客户端 HTTP 请求
  │
  ├── RequestDecoder: JSON → Request 结构体
  │
  ├── DAG 并行执行
  │     ├── configstoreDataOp    → 加载 reason_config 推荐理由配置
  │     ├── gjsonPickStrOp       → 提取请求字段(roleid/areaid/envid/...)
  │     ├── stringFormat*Op      → 生成 Redis Key 模板
  │     ├── redisMgetOp          → 批量拉离线推荐列表
  │     ├── redisHgetallOp       → 获取好友关系哈希、车队标签
  │     ├── customEsGetOp        → 查询 ES 车队详情
  │     ├── customUserCacheOp    → 检查用户推荐缓存（v1 专有）
  │     └── customConvPictureOp  → 聚合成员头像
  │
  ├── mergeRecDataOp             → 核心合并：车队信息 + 好友信息 + 推荐理由 + 分页
  │
  ├── redisSetEXOp               → 缓存结果（TTL 3600s）
  │
  └── ResponseComposer           → 组装响应 + 生成 RecommendID (UUID)
```

---

## Flink 实时流计算详解

### 整体 Job 架构

车队推荐涉及两个核心 Flink Job，分工明确：

| Job 名称 | 入口类 | 输出目标 | 作用 |
|---------|--------|---------|------|
| `wedo_car_friend_online` | FriendsOperatorAsync | Redis Hash | 实时维护用户好友关系链，是推荐服务的核心数据源 |
| `wedo_write_convoy_task` | ConvoyInfoToES | Elasticsearch | 实时同步车队变更，供推荐服务 ES 查询 |

两个 Job 均遵循相同的容灾模式：**主从双 Kafka + `.union()` 合并**，任意一路断流不影响数据处理。

代码模块结构：

```
bus/wedo/
├── car/
│   ├── online/
│   │   ├── FriendsOperatorAsync.java   ← 好友关系链（核心）
│   │   └── ConvoyInfoToES.java         ← 车队 ES 同步（核心）
│   └── redis/
│       ├── AsyncRedisHGet.java         ← 异步 HGETALL 封装
│       └── AsyncRedisSMembers.java
└── video/
    └── online/
        └── PvpToES.java                ← PvP 对局数据（与车队推荐无直接关联）
```

---

### FriendsOperatorAsync —— 好友关系链维护

这是推荐系统最核心的 Flink Job。它消费两路 Kafka 数据源，将玩家的好友关系实时写入 Redis Hash，为推荐服务提供"哪些好友在哪支车队"的数据基础。

**执行拓扑**

```
[Source] playerfriendslistflow (master)  ─┐
[Source] playerfriendslistflow (slave)   ─┼─► union ─► flatMap ① ─► keyBy ②
                                          │               (TLog解析)  (envid_ruid)
                                          │                               │
                                          │                               ▼
                                          │                  KeyedProcessFunction ③
                                          │                  (ListState + ValueState)
                                          │                  (2s Processing Time Timer)
                                          │                               │
                                          │                               ▼
                                          │                  AsyncDataStream ④
                                          │                  unorderedWait(AsyncRedisHGet)
                                          │                  timeout=60s, capacity=20
                                          │                               │
                                          │                               ▼
                                          │                  flatMap ⑤ (全量合并逻辑)
                                          │                  → asyncCommands.hmset()
                                          │                  → asyncCommands.hdel()
                                          │                               │
[Source] friendflow (master)  ─┐          │                               │
[Source] friendflow (slave)   ─┼─► union ─┼─► flatMap ⑥                  │
                               │          │   (增量解析)                   │
                               └──────────┘       │                       │
                                                  └──────── union ⑦ ─────┘
                                                                 │
                                                             keyBy ⑧
                                                                 │
                                                                 ▼
                                                  KeyedProcessFunction ⑨
                                                  (MapState TTL=1day)
                                                  (事件类型分发)
                                                  → asyncCommands.hset()
                                                  → asyncCommands.hdel()
```

**算子逐一详解**

**① flatMap（RichFlatMapFunction）—— TLog 流水解析**

TLog 流水以 `|` 分隔，索引固定。`flatMap` 按字段索引提取关键字段，过滤无效记录后向下游 emit：

```java
.flatMap(new RichFlatMapFunction<String, Tuple5<...>>() {
    public void flatMap(String s, Collector<...> collector) {
        String[] info = s.split("\\|");
        if (info.length > 25) {
            String rUid       = info[12];
            String friendType = info[21];  // 1=游戏 2=平台 3=双端
            String friends    = info[24];  // 分号分隔的好友列表
            String envid      = info[25];
            if (envid.equals("") || envid.equals("0")) return; // 过滤无效环境
            collector.collect(Tuple5.of(envid + "_" + rUid, friends, "", friendType, loginType));
        }
    }
})
```

**② keyBy —— 按 envid_roleid 分区**

将同一用户的所有流水路由到同一个 TaskManager 子任务，确保状态一致性：

```java
.keyBy(data -> data.f0)  // 例："34_576465154108032693"
```

**③ KeyedProcessFunction + Timer —— 2 秒微批聚合窗口**

全量登录流水中，同一用户可能上报多条（游戏好友和平台好友各一条），需在 2s 内合并后再统一处理：

```java
// 状态均配置 10s TTL，避免 Job 重启后状态积压导致 OOM
ListState<String>   friendListState;  // 累积同次登录的多条好友数据
ValueState<Boolean> isRegister;       // Timer 防重复注册标记

public void processElement(Tuple5<...> ins, Context ctx, Collector<...> out) {
    if (isRegister.value() == null) {
        // 仅第一条记录触发：注册 2 秒后的 Processing Time Timer
        ctx.timerService().registerProcessingTimeTimer(
            ctx.timerService().currentProcessingTime() + 2000
        );
        isRegister.update(true);
    }
    friendListState.add(ins.f1 + "|" + ins.f3);  // 追加到 ListState
}

public void onTimer(long timestamp, OnTimerContext ctx, Collector<...> out) {
    // 2s 后：汇总 ListState 中所有条目，合并为一条数据发出
    List<String> friendList = new ArrayList<>();
    for (String element : friendListState.get()) friendList.add(element);
    out.collect(new Tuple4<>("", "", ctx.getCurrentKey(), friendList));
    friendListState.clear();
    isRegister.clear();
}
```

> **TTL 设计要点**：状态 TTL 设为 10s（远大于 2s Timer），是为了防止 TaskManager 重启导致 Timer 尚未触发时状态已提前过期、数据丢失。Timer 触发后主动 `clear()`，状态实际占用时间很短。

**④ AsyncDataStream.unorderedWait —— 异步 Redis HGETALL**

全量好友合并前，需先读取该用户在 Redis 中的历史好友哈希。使用 Flink 异步 IO 算子，避免同步阻塞流水线：

```java
AsyncDataStream.unorderedWait(
    login_full_friend,           // 上游流（Tuple4）
    new AsyncRedisHGet(...),     // 自定义 RichAsyncFunction
    60, TimeUnit.SECONDS,        // 超时 60s
    20                           // 最大并发异步请求数（控制背压）
)
```

`AsyncRedisHGet` 内部通过 `CompletableFuture` 非阻塞执行 Redis `HGETALL`：

```java
public void asyncInvoke(Tuple4<...> input, ResultFuture<...> resultFuture) {
    CompletableFuture.supplyAsync(() -> client.hgetAll(prefix + input.f2))
        .thenAccept(res ->
            resultFuture.complete(Collections.singleton(Tuple2.of(res, input)))
        );
}

// 超时降级：不下发数据，跳过本次全量合并（相当于保留旧数据）
public void timeout(...) {
    resultFuture.complete(Collections.emptyList());
}
```

> 选用 `unorderedWait` 而非 `orderedWait`：允许不同 key 的 Redis 响应乱序返回，吞吐量更高。同一 key 的数据顺序由上游 `keyBy` 已保证，不存在乱序问题。

**⑤ flatMap（RichFlatMapFunction）—— 全量好友三路合并 + 写 Redis**

拿到 Redis 历史数据后，与日志好友列表进行合并，分三种情况处理：

| 情况 | 处理逻辑 |
|------|---------|
| 历史有 & 日志也有 | 保留历史的**组队时间**，更新亲密度和好友类型 |
| 历史有 & 日志没有 | 加入删除列表，执行 `hdel` 清除旧好友 |
| 历史没有 & 日志有 | 新增记录，组队时间初始化为空 |

写入使用 Lettuce 异步命令，不阻塞算子主线程：

```java
// 批量写入新/更新的好友
asyncCommands.hmset(friendListPrefix + uidKey, result);
// 删除日志中已不存在的好友
asyncCommands.hdel(friendListPrefix + uidKey, deleteKey.toArray(new String[0]));
```

**⑥ flatMap（RichFlatMapFunction）—— 增量好友流水解析**

增量流 (`friendflow`) 字段含事件类型 `friendFlowType`，按类型 emit 不同语义的 Tuple6：

```java
// 1: 申请；2: 同意添加   → addFriendType
// 4: 删除好友           → deleteFriendType
// 5: 查看好友信息       → updateIntimacyType（亲密度有变化时）
// 7: 亲密度变化         → updateIntimacyType
// 8: 好友组队          → updatePlayTimeType
```

**⑦⑧ union + 二次 keyBy —— 合流后重新分区**

全量处理流与增量流 `union` 合并，再次 `keyBy` 路由到同一有状态算子，统一处理：

```java
processFullList.union(incr_friend)
    .keyBy(data -> data.f0)  // 同样按 envid_roleid 分区
    .process(new KeyedProcessFunction<...>() { ... })
```

**⑨ KeyedProcessFunction —— 事件分发 + Redis 写入**

维护一份天级快照 `MapState<String, String>`（好友 ruid → 关系值字符串），按事件类型分发处理：

```java
MapState<String, String> loginFriendState;  // TTL = 1 天

public void processElement(Tuple6<...> data, Context ctx, Collector<...> out) {
    switch (data.f4) {
        case loginType:          // 全量登录：覆盖当日好友快照
            loginFriendState.clear();
            loginFriendState.putAll(data.f5);
            break;

        case addFriendType:      // 新增好友：双向写入（我记他，他也记我）
            asyncCommands.hset(friendListPrefix + data.f0, data.f1, value);
            asyncCommands.hset(friendListPrefix + idInfo[0] + "_" + data.f1, idInfo[1], value);
            break;

        case deleteFriendType:   // 删除好友：只删游戏好友，平台好友不删
            if (!reportType.equals("2")) {
                asyncCommands.hdel(friendListPrefix + data.f0, data.f1);
            }
            break;

        case updateIntimacyType:
        case updatePlayTimeType: // 更新亲密度或组队时间：重新拼接 value 后 hset
            value = data.f2 + "|" + data.f3 + "|" + rmpFriendType + "|" + intimacyType;
            asyncCommands.hset(friendListPrefix + data.f0, data.f1, value);
            break;
    }
}
```

> **好友类型保护逻辑**：删除事件中，若对方为平台好友（`isPlatFriend=1`），不执行 `hdel`，保证跨平台好友关系不因游戏侧删除操作而误删。

**最终 Redis 写入格式：**

```
Key:   wedo_car_friendlist_online_{envid}_{roleid}
Type:  Hash
Field: {friendRuid}
Value: {亲密度}|{最近组队时间}|{好友类型(1/2/3)}|{亲密度类型}

示例：
HGETALL wedo_car_friendlist_online_34_576465154108032693
→ "123456789" → "1500|1718250000|3|1"
   "987654321" → "800||1|0"
```

---

### ConvoyInfoToES —— 车队信息实时同步

**执行拓扑**

```
[Source] modifyconvoyflow (master) ─┐
[Source] modifyconvoyflow (slave)  ─┼─► union ─► keyBy ① ─► KeyedProcessFunction ② ─► EsSink ③
                                                  (ConvoyId)  (ValueState + MapState        (upsert)
                                                               5s Timer + 乱序过滤)
```

**算子详解**

**① keyBy —— 按 ConvoyId 分区**

同一支车队短时间内可能产生多次变更流水（改名、成员加入、踢人等），需路由到同一子任务聚合：

```java
.keyBy(value -> value.split("\\|")[14])  // ConvoyId 在第 14 个字段
```

**② KeyedProcessFunction —— 5 秒聚合窗口 + 乱序过滤**

与 FriendsOperatorAsync 的 2s Timer 方案类似，但增加了**乱序流水保护**：

```java
ValueState<Boolean>     isRegister;    // Timer 注册标记（无 TTL，持久）
MapState<String,Object> mapConvoyHash; // 车队字段暂存（无 TTL，持久）

public void processElement(String data, Context ctx, Collector<...> out) {
    String[] fields = data.trim().split("\\|");
    if (fields.length <= 33) return;  // 字段不足，丢弃无效流水

    // ── 乱序保护：忽略比已有 UpdateTime 更早的流水 ──
    long currModifyTime = Long.parseLong(fields[33]);
    Long lastUpdateTime = (Long) mapConvoyHash.get("UpdateTime");
    if (lastUpdateTime != null && lastUpdateTime > currModifyTime) return;

    // ── 增量更新 MapState ──
    mapConvoyHash.put("ConvoyId",    fields[14]);
    mapConvoyHash.put("ConvoyName",  fields[15]);
    mapConvoyHash.put("MemberList",  fields[24].split(","));
    mapConvoyHash.put("MemberNum",   Integer.parseInt(fields[25]));
    mapConvoyHash.put("UpdateTime",  currModifyTime);
    // ModifyType=1(创建)：记录 CreateTime
    if (ModifyType.equals(1)) mapConvoyHash.put("CreateTime", System.currentTimeMillis() / 1000);
    // ModifyType=6(解散)：标记已删除
    if (ModifyType.equals(6)) mapConvoyHash.put("DeletedConvoy", "1");

    // ── 首条流水触发 5s Timer，后续流水只更新状态 ──
    if (isRegister.value() == null) {
        ctx.timerService().registerProcessingTimeTimer(
            ctx.timerService().currentProcessingTime() + 5000
        );
        isRegister.update(true);
    }
}

public void onTimer(long timestamp, OnTimerContext ctx, Collector<...> out) {
    // 5s 后：将 MapState 中所有字段组装为 HashMap 发往 ES Sink
    HashMap<String, Object> json = new HashMap<>();
    for (Map.Entry<String, Object> e : mapConvoyHash.iterator()) {
        json.put(e.getKey(), e.getValue());
    }
    out.collect(json);
    isRegister.clear();
    mapConvoyHash.clear();
}
```

> **与 FriendsOperatorAsync 的 Timer 设计对比**：
> - FriendsOperatorAsync 用 `ListState + 2s Timer`：聚合同一用户**同次登录的多条好友流水**（不同 friendType），目的是去重合并
> - ConvoyInfoToES 用 `MapState + 5s Timer`：聚合同一车队**短时间内的多次字段变更**，减少对 ES 的写入次数，降低 IO 压力

**③ EsSink —— 按 ConvoyId 做 Upsert**

ES Sink 以 `ConvoyId` 为文档 ID，天然支持 upsert 语义，重复写入安全：

```java
env.addSink(
    EsSink.getESSinkBuilder(esIp, esPort, esUser, esPassword, esIndex, "ConvoyId")
          .build()
);
```

**ES 文档结构示例：**

```json
{
  "ConvoyId": "convoy_123",
  "ConvoyName": "JaydenAn的车队",
  "PlatID": 1,
  "PublishID": 1,
  "EnvId": "34",
  "MemberList": ["ruid_001", "ruid_002"],
  "MemberNum": 5,
  "HeadRuid": "ruid_001",
  "MemberPublishConvoyList": ["ruid_003"],
  "ConvoyTagList": "tag1,tag2,,",
  "CreateTime": 1718250000,
  "UpdateTime": 1718260000,
  "DeletedConvoy": "0"
}
```

---

### 状态管理汇总

| Job | 状态描述符名称 | 类型 | TTL | 用途 |
|-----|--------------|------|-----|------|
| FriendsOperatorAsync | `wedo_car_is_register_online_nnew03` | `ValueState<Boolean>` | 10s | 全量流：防 Timer 重复注册 |
| FriendsOperatorAsync | `wedo_car_loginonce_friend_online_new03` | `ListState<String>` | 10s | 全量流：聚合同次登录的多条好友数据 |
| FriendsOperatorAsync | `wedo_car_loginoneday_friend_online_new04` | `MapState<String,String>` | 1 天 | 增量流：当日好友关系快照，供亲密度/组队时间更新查询 |
| ConvoyInfoToES | `wedo_convoy_is_register` | `ValueState<Boolean>` | 无 | 防 Timer 重复注册 |
| ConvoyInfoToES | `wedo_convoy_info` | `MapState<String,Object>` | 无 | 车队多字段暂存，等待 5s 后统一写 ES |

**TTL 策略说明：**

- **10s TTL**：与 2s Timer 对齐，Timer 触发并 `clear()` 后状态自然过期，防止因 Job 异常重启导致 ListState 数据长期滞留
- **1 天 TTL**：对应用户日活跃周期，超过 1 天未登录的用户状态自动清理，控制 TaskManager 内存占用
- **无 TTL**：ConvoyInfoToES 的状态理论上随 `onTimer → clear()` 被清空；若 Job 重启，状态会保留并在下次 Timer 触发时再次写出，但 ES upsert 语义保证幂等性

---

### 双 Kafka 容灾模式

两个 Job 均采用相同的容灾模式，以 FriendsOperatorAsync 为例：

```java
FlinkKafkaConsumer<String> consumer1 = KafkaUtil.initKafka(kafkaMaster, topic, groupID);
FlinkKafkaConsumer<String> consumer2 = KafkaUtil.initKafka(kafkaSlave,  topic, groupID);

env.addSource(consumer1)
   .union(env.addSource(consumer2))  // 两路数据流合并
   ...
```

| Kafka 集群 | 地址 | 角色 |
|-----------|------|------|
| Master | `204-kafka-clb.kpmq.woa.com:9092` | 主集群 |
| Slave  | `199-kafka-clb.kpmq.woa.com:9092` | 备集群 |

两路 Consumer 使用**相同的 `groupID`**，Kafka 消费组语义保证同一 partition 的消息只被消费一次。若两路均正常，消息会被重复消费；下游通过以下机制保证幂等性：

| Job | 幂等保障机制 |
|-----|------------|
| FriendsOperatorAsync | `hmset`/`hset` 覆盖写，重复写入值不变 |
| ConvoyInfoToES | ES 按 ConvoyId upsert；MapState 内 UpdateTime 乱序过滤 |

---

## 具体实现（推荐服务）

### 推荐服务 v1（1126）- 通用分页模式
 
**请求 / 响应结构**

```go
// 请求
type Request struct {
    Credid  string `json:"credid"`   // 凭证ID: "wedo_teamrec_26933"
    Flowid  string `json:"flowid"`   // 链路追踪ID
    ReqTime string `json:"req_time"` // 请求时间戳
    Userid  string `json:"userid"`   // 用户ID
    Data struct {
        Roleid    string `json:"roleid"`     // 角色ID（主键）
        Areaid    string `json:"areaid"`     // 区域ID
        EnvID     string `json:"env_id"`     // 环境ID（服务器分区）
        Platid    string `json:"platid"`     // 平台ID（iOS/Android）
        ReqNum    string `json:"req_num"`    // 请求数量
        PageStart string `json:"page_start"` // 分页起始偏移 ← v1 特有
    }
}

// 响应
type RecList struct {
    Convoy      string       `json:"convoy"`
    ConvoyName  string       `json:"convoy_name"`
    MemberNum   string       `json:"member_num"`
    FriendList  string       `json:"friend_list"`
    Ext1        string       `json:"ext1"`
    UserLists   []UserList   `json:"user_list"`
    ReasonLists []ReasonList `json:"reason_list"`
    // 内部排序字段（不序列化）
    ConvoyFriendNum  int
    IntimacyValue    int64
    RecentBrPlayTime int64
}
```

**DAG 执行图（核心路径）**

```
[A] configstoreDataOp          → 加载推荐理由配置（configstore）
[B] gjsonPickStrOp(×6)         → 提取 roleid/areaid/platid/env_id/req_num/page_start
[C] setUserStringAttrOp(×6)    → 注入 DAG Context
[D] stringFormat3Op            → 生成 Redis Key: wedo_car_{env}_{area}_{roleid}
[E] redisMgetOp                → 批量获取离线推荐列表
[F] redisHgetallOp             → 获取好友关系链（FriendsOperatorAsync 写入）
[G] redisHgetallOp             → 获取车队标签
[H] redisLrangeWithCacheOp     → 获取版本列表
[I] customUserCacheOp          → 检查分页缓存（page_start > 0 时读缓存）
[J] customEsGetOp              → ES 查询车队详情
    ├── 输出: vecConvoyInfo     → 车队详情切片
    ├── 输出: recallConvoyList  → 召回列表（合并在线/离线）
    ├── 输出: offlineRecList    → 离线召回子集
    └── 输出: onlineRecList     → 在线召回子集
[K] customConvPictureOp        → 聚合多来源成员头像 map
[L] mergeRecDataOp             → 核心合并（11个输入参数，见下表）
[M] redisSetEXOp               → 缓存结果（TTL=3600s）
[N] metricsCounterReportOp     → 上报监控指标
```

**mergeRecDataOp 输入（v1，11 个参数）：**

| # | 参数名 | 类型 | 说明 |
|---|--------|------|------|
| 0 | reasonConfig | `[]string` | 推荐理由配置 |
| 1 | recallConvoyList | `[]string` | 总召回列表 |
| 2 | userLabel | `map[string]string` | 用户标签 |
| 3 | carLabel | `map[string]string` | 车队标签 |
| 4 | vecConvoyInfo | `[]ConvoyInfo` | 车队详情切片 |
| 5 | mapFriendList | `map[string]string` | 好友关系哈希（Flink 写入） |
| 6 | offlineMemberPicture | `map[string]string` | 成员头像 |
| 7 | offlineRecConvoyList | `[]string` | 离线召回子集 |
| 8 | onlineRecConvoyList | `[]string` | 在线召回子集 |
| 9 | pageSize | `string` | 每页大小 |
| 10 | pageStart | `string` | 分页偏移 |

**用户缓存机制（customUserCacheOp）**

v1 引入用户级缓存，当 `page_start > 0` 时，直接从 Redis 读取已缓存的推荐结果切片，避免重复查询：

- `page_start=0`：全量查询 ES + Redis → 合并 → 缓存到 Redis（TTL=3600s）
- `page_start>0`：直接读 Redis 缓存，按 `[page_start : page_start+page_size]` 切片返回

**图片批量获取（DataMore 服务）**

```go
// 每批最多 20 个 ruid，并发请求 DataMore 服务
// DataMore URL: polaris://64987585:1114112
// 超时: 800ms
func getBatchRuidPicture(ruidList map[string]struct{}) map[string]string { ... }
```

---

## 方案对比分析
 
| 策略维度 | v1 | v2 |
|---------|----|----|
| **召回** | Redis 离线推荐列表（主）+ ES 在线补充 | 直接 ES 全量扫描（依赖 friend_ruid 过滤）|
| **排序信号** | 好友数 + 亲密度 + 最近组队时间 | 基于 mapMemberPublish 成员发布信息 |
| **缓存策略** | 分页级别缓存（1h TTL）| 无缓存 |
| **个性化程度** | 中（基于好友关系）| 高（聚焦到单个好友的车队）|
| **冷启动** | 有离线兜底 | 强依赖实时 ES 数据 |

### 架构演进方向

```
通用推荐（v1）：
  用户进入"发现车队"频道 → 系统为其推荐多支满足条件的车队
  适用场景：首页推荐流、分页浏览

定向推荐（v2）：
  用户查看某个好友 → 系统展示"这个好友所在的车队"
  适用场景：社交关系链驱动的推荐、"你的好友 XXX 在这个车队"
```

### 技术方案横向对比

**方案 A（当前实现）：Flink + DAG 推荐服务**

**优点：**
- 算法可热更新（修改 graph.toml 无需重新编译）
- Flink 保证实时数据的准确性和幂等性（乱序流水过滤）
- Redis 缓存保证低延迟

**缺点：**
- 架构链路长，排查问题复杂
- DAG 配置文件 TOML 体积大，可读性差
- 离线推荐和在线服务数据格式需要对齐
 
**方案 B（对比）：向量化召回 + 实时精排**

引入向量数据库存储车队 embedding，离线训练用户偏好 embedding，在线通过 ANN 检索。

| | 方案 A（当前）| 方案 C（对比）|
|-|--------------|--------------|
| 推荐质量 | 中（规则+协同过滤）| 高（深度学习）|
| 工程复杂度 | 中 | 高 |
| 训练成本 | 低 | 高 |
| 适用规模 | 中等体量 | 大规模 |

> **结论**：当前方案 A 适合现阶段业务规模，随着车队数量和用户数增长，可逐步向方案 C 演进。

---

## 核心算法与数据结构

### 好友关系链数据结构

```
Redis Hash: wedo_car_friendlist04_{env}_{envid}_{roleid}

Field (Key)  = friendRuid
Value格式    = {亲密度}|{最近组队时间}|{好友类型}|{亲密度类型}

好友类型:
  1 = 游戏好友（仅游戏内添加）
  2 = 平台好友（微信/QQ好友导入）
  3 = 双端好友（两种都有）
```

### 车队排序算法（v1 merge 阶段）

基于以下维度综合排序（优先级从高到低）：

1. **车队内好友数** (`ConvoyFriendNum`) — 社交关系最强信号
2. **好友亲密度之和** (`IntimacyValue`) — 衡量关系深度
3. **最近组队时间** (`RecentBrPlayTime`) — 衡量活跃程度
4. **成员数** (`ConvoyMemberNum`) — 辅助参考

```go
sort.Slice(recList, func(i, j int) bool {
    if recList[i].ConvoyFriendNum != recList[j].ConvoyFriendNum {
        return recList[i].ConvoyFriendNum > recList[j].ConvoyFriendNum
    }
    if recList[i].IntimacyValue != recList[j].IntimacyValue {
        return recList[i].IntimacyValue > recList[j].IntimacyValue
    }
    if recList[i].RecentBrPlayTime != recList[j].RecentBrPlayTime {
        return recList[i].RecentBrPlayTime > recList[j].RecentBrPlayTime
    }
    return recList[i].ConvoyMemberNum > recList[j].ConvoyMemberNum
})
```

### 推荐理由（ReasonList）配置结构

通过 ConfigStore 配置中心动态管理，无需重新部署：

```json
{
  "reason_id": "1001",
  "reason_name": "好友在此车队",
  "type": "friend",
  "source": "friend_list",
  "calc_type": "count",
  "redis_format": "wedo_car_friendlist04_{env}_{envid}_{roleid}"
}
```

---

## 存储设计

### Redis 键空间设计

| 键模板 | 类型 | 说明 | TTL |
|--------|------|------|-----|
| `wedo_car_friendlist04_{env}_{envid}_{roleid}` | Hash | 好友关系链 | 无（持久）|
| `wedo_car_user_req_cache_{env}_{area}_{roleid}` | String | 推荐结果缓存 | 3600s |
| `wedo_offline_rec_{env}_{area}_{roleid}` | List | 离线推荐车队列表 | 无 |
| `wedo_car_label_{env}_{convoyId}` | Hash | 车队标签 | 无 |
 
### Elasticsearch 索引设计

```
索引: wedo_convoy_{env}
主键: ConvoyId

关键字段:
  - ConvoyId (keyword)
  - EnvId (keyword)        ← 必过滤
  - PublishID (integer)
  - MemberNum (integer)
  - CreateTime (long)
  - UpdateTime (long)
  - DeletedConvoy (keyword) ← 过滤已解散
  - MemberList (keyword[])
  - MemberPublishConvoyList (keyword[])
```

---

## 性能与可靠性设计

### 延迟优化

| 优化点 | 方式 | 效果 |
|--------|------|------|
| DAG 并行执行 | 无依赖节点并发运行 | 整体耗时取最长依赖路径 |
| Redis 用户缓存 | 第二页起直接读缓存 | 命中时 <50ms |
| DataMore 批量查询 | 20个/批并发请求 | 减少 HTTP 次数 |
| jsoniter 快速 JSON | 替代标准库 `encoding/json` | 解析加速 2-3x |
| gjson 零反射解析 | 不全量反序列化 | 降低 CPU |

### 容灾设计

| 组件 | 容灾方式 |
|------|---------|
| Kafka 数据源 | 主从双集群 `.union()` 合并，任意一个故障不影响 |
| 离线推荐 | Redis 持久存储，ES 查询作为兜底 |
| 请求超时 | 全局 2000ms，子操作 800ms |
| Flink 乱序处理 | UpdateTime 去重，忽略旧版本流水 |
| Redis 连接池 | max_active=120，is_wait=true 队列等待 |

### 幂等性保证

- `ConvoyInfoToES`：ES 按 ConvoyId 做 upsert，重复写入安全
- `FriendsOperatorAsync`：hmset/hset 覆盖写；天级 TTL 状态防止 OOM
- 推荐服务：只读操作 + 缓存写，天然幂等

---

## 部署与运维

### 环境管理

| 环境 | 说明 | 对应代码包 |
|------|------|----------|
| test | 本地/测试 | `bus/wedo/*/test/` |
| pre | 预发布 | `bus/wedo/*/pre/` |
| online | 生产 | `bus/wedo/*/online/` |

### Flink 任务提交

| Job 名称 | 主类 | 说明 |
|---------|------|------|
| `wedo_write_convoy_task` | ConvoyInfoToES | 车队信息同步 |
| `wedo_car_friend_online` | FriendsOperatorAsync | 好友关系链 |
| `wedo_pvp_end_to_es` | PvpToES | 对局数据写 ES |
 
---
 

## 附录：关键配置速查

### Kafka Topic 速查

| Topic | 用途 | 消费 Flink 任务 |
|-------|------|----------------|
| `billow_ex_wedo_playerfriendslistflow` | 全量好友关系（每日登录）| FriendsOperatorAsync |
| `billow_ex_wedo_friendflow` | 增量好友变更 | FriendsOperatorAsync |
| `billow_ex_wedo_modifyconvoyflow` | 车队信息变更 | ConvoyInfoToES |
| `billow_ex_wedo_pvpendflow` | PvP 对局结束 | PvpToES |

### Flink 算子速查

| 算子 | 使用 Job | 用途 |
|------|---------|------|
| `addSource().union()` | 全部 | 双 Kafka 容灾 |
| `flatMap(RichFlatMapFunction)` | FriendsOperatorAsync | TLog 管道符字段解析 |
| `keyBy(KeySelector)` | 全部 | 按 envid_roleid / ConvoyId 分区 |
| `process(KeyedProcessFunction)` + `ValueState` + `ListState` + `Timer` | FriendsOperatorAsync | 2s 微批聚合全量好友流水 |
| `AsyncDataStream.unorderedWait(RichAsyncFunction)` | FriendsOperatorAsync | 异步 Redis HGETALL，不阻塞主流水线 |
| `process(KeyedProcessFunction)` + `MapState` + `Timer` + 乱序过滤 | ConvoyInfoToES | 5s 聚合车队字段变更，写 ES |
| `addSink(EsSink)` | ConvoyInfoToES | ES upsert 写入 |
