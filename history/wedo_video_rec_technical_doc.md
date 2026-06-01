# WeDo 对局推荐系统技术实现文档（实时流计算部分）

---

## 需求背景

### 业务场景

WeDo 平台的**对局回放（Video）** 功能允许玩家上传、观看、点赞对局录像。推荐系统的目标是将"最值得看"的对局视频推给合适的玩家——即那些高质量、高竞技水平、热门程度适中的对局。

实时流计算层的目标是：**从游戏服务器上报的对局结束事件中提取多维度信号**（对局元信息、阵营表现、角色伤害、玩家胜率、视频点赞数、排名分数），并将这些信号实时写入 Elasticsearch 和 Redis，为推荐服务的**召回**和**排序**阶段提供准确的特征数据。

### 核心挑战

| 挑战 | 描述 |
|------|------|
| **多 Job 协作拼接 ES 文档** | `wedo_video_rec_online` 索引由 4 个 Job 分别写入不同字段，依赖 ES upsert 语义逐步拼装完整文档 |
| **延迟读取依赖** | 部分 Job 需要读取其他 Job 的计算结果。为保证数据就绪，通过 ProcessingTime Timer 引入固定延迟（3s/4s/8s）再触发读取 |
| **排行榜内存与持久化双写** | 排行榜需要高吞吐的内存 PriorityQueue，同时通过 Redis 持久化保证 Job 重启后状态恢复 |
| **玩家胜率滑动窗口** | 每个玩家最近 100 场对局的多粒度胜率（近 10/20/30/40/50/全部场次）需维护在 Redis Hash 中，并清除 7 天前的旧数据 |
| **点赞排行榜多维度分组** | 点赞和排名排行榜需按（env+版本）、（env+版本+任务类型）、（env+版本+任务类型+角色）等 4 个维度分别维护 |

### 技术选型

- **Flink**：处理多路 Kafka 流水，完成对局元信息解析、阵营/角色聚合、胜率计算、排行榜维护等实时计算
- **Elasticsearch**：存储对局结构化信息，索引名 `wedo_video_rec_online`，主键 `key`（`env_roomId`）
- **Redis (TendisPlus)**：缓存对局详情（供下游 Job 读取）、玩家胜率历史记录、排行榜数据、点赞计数

---

## Flink 实时流计算详解

### 整体 Job 架构

对局推荐涉及 7 个 Flink Job，构成一条**分阶段的数据流水线**：

| Job 名称 | 入口类 | 核心输入 | 输出目标 | 作用 |
|---------|--------|---------|---------|------|
| `wedo_pvp_end_to_es` | PvpToES | pvpendflow + uploadvideoflow | ES | 写入对局元信息、标记视频上传状态 |
| `wedo_pvp_player_to_es` | PvpPlayerToES | pvpplayerendflow | ES | 聚合对局玩家列表（uid、昵称、角色、排名分） |
| `wedo_pvp_roleid_damage_to_es` | PvpRoleDamageToES | pvproleendflow | ES | 聚合对局角色伤害/击杀表现 |
| `wedo_pvp_battle_detail_redis` | BattleDetailToRedis | uploadvideoflow | Redis | 延迟 3s 后将完整 ES 文档缓存到 Redis |
| `wedo_pvp_win_percent` | ComputerUserWinPercent | uploadvideoflow | Redis + ES | 延迟 4s 后计算玩家胜率，写入 ES |
| `wedo_pvp_rank_sort` | PvpRankSort | uploadvideoflow | Redis | 延迟 8s 后更新对局排名排行榜 |
| `wedo_pvp_like_num_sort` | VideoLoveToES | gameplayrecordlikeflow | Redis + ES | 实时维护点赞数、点赞排行榜、玩家点赞 Hash |

**数据流依赖关系：**

```
[对局结束事件]
        │
        ├─► PvpToES ──────────────────────────────┐
        │                                         │ (4个Job各写部分字段，ES upsert 拼装)
        ├─► PvpPlayerToES ───────────────────────►│
        │                                         │ ES: wedo_video_rec_online
        ├─► PvpRoleDamageToES ───────────────────►│
        │                                         │
        └─► ComputerUserWinPercent ──────────────►┘
                (延迟 4s, 读 Redis 缓存)

[视频上传事件]
        │
        ├─► PvpToES (标记 is_success)
        │
        ├─► BattleDetailToRedis (延迟 3s → 读 ES → 缓存到 Redis)
        │                     Redis: wedo_video_online_details_{key}
        │
        ├─► ComputerUserWinPercent (延迟 4s → 读 Redis 缓存 → 计算胜率)
        │
        └─► PvpRankSort (延迟 8s → 读 Redis 缓存 → 更新排名排行榜)

[点赞事件]
        └─► VideoLoveToES → 点赞数 ES/Redis + 点赞排行榜 Redis + 玩家点赞 Hash Redis
```

> **延迟读取模式**：`BattleDetailToRedis` 在视频上传后 3s 才将 ES 文档写入 Redis；`ComputerUserWinPercent` 和 `PvpRankSort` 分别延迟 4s 和 8s 触发，确保读取时 Redis 缓存已就绪。

代码模块结构：

```
bus/wedo/
└── video/
    ├── online/
    │   ├── PvpToES.java                ← 对局元信息 + 视频上传标记
    │   ├── PvpPlayerToES.java          ← 对局玩家列表聚合
    │   ├── PvpRoleDamageToES.java      ← 角色伤害/击杀聚合
    │   ├── BattleDetailToRedis.java    ← ES 文档缓存到 Redis（3s 延迟）
    │   ├── VideoLoveToES.java          ← 点赞数维护 + 点赞排行榜
    │   ├── PvpRankSort.java            ← 排名排行榜（8s 延迟）
    │   └── ComputerUserWinPercent.java ← 玩家胜率计算（4s 延迟）
    ├── pojo/
    │   └── PvpBatterInfo.java          ← 玩家对局信息 POJO
    └── redis/
        ├── AsyncRedisPvpDetail.java    ← 异步读取 Redis 对局缓存
        └── AsyncRedisPlayerRecordHGet.java ← 异步读取玩家胜率 Hash
```

所有 Job 均采用**主从双 Kafka + `.union()` 合并**容灾模式。

---

### PvpToES —— 对局元信息同步

消费两路独立的 Kafka 流水（对局结束流、视频上传流），将对局的基础信息和视频上传状态分别写入 ES。

**执行拓扑**

```
[Source] pvpendflow    (master) ─┐
[Source] pvpendflow    (slave)  ─┼─► union ─► ProcessFunction ① ─► ES Sink
                                             (解析对局元信息)

[Source] uploadvideoflow (master) ─┐
[Source] uploadvideoflow (slave)  ─┼─► union ─► ProcessFunction ② ─► ES Sink
                                              (标记视频上传状态)
```

**① ProcessFunction —— pvpendflow 对局结束事件**

提取对局元信息写入 ES，字段 `pvp_end_success="1"` 作为"对局结束数据就绪"的标记：

```java
// 过滤条件：
// data[37] == "0"：非 AI 对局
// data[39] 非空且包含逗号：有有效的游戏标签（game_play_record_tags）
if (!data[37].equals("0") || CommonUtil.isEmpty(data[39]) || !data[39].contains(",")) return;

String env     = data[17];
String roomId  = data[11];
String key     = env + "_" + roomId;

HashMap<String, Object> json = new HashMap<>();
json.put("key",         key);
json.put("roomid",      roomId);
json.put("fubenid",     data[22]);       // 关卡 ID
json.put("fubenhardlv", data[23]);       // 关卡难度
json.put("score",       data[26]);       // 胜方分数
json.put("score_sum",   data[27]);       // 总分数
json.put("wincampid",   data[28]);       // 胜利阵营 ID
json.put("is_finaldeath", data[34]);     // 是否全灭
json.put("mmr",         data[35]);       // MMR 分
json.put("time",        data[36]);       // 对局时长
json.put("start_time",  data[38]);       // 对局开始时间
json.put("game_play_record_tags", data[39].split(","));  // 标签数组
json.put("env",         env);
json.put("version",     data[15]);
json.put("pvp_end_success", "1");       // 标记对局结束数据就绪
out.collect(json);
```

**② ProcessFunction —— uploadvideoflow 视频上传事件**

标记该对局的视频已完成上传，并写入随机分桶字段：

```java
// 过滤条件：isAI == "0"（非 AI 对局才纳入推荐）
String isAI  = data[18].replaceAll("\n|\r", "");
if (!isAI.equals("0")) return;

int rand = new Random().nextInt(20);  // [0,20) 随机分桶

HashMap<String, Object> json = new HashMap<>();
json.put("key",        data[17] + "_" + data[11]);  // env_roomId
json.put("is_success", "1");   // 标记视频上传成功
json.put("is_ai",      isAI);
json.put("rand",       rand);
out.collect(json);
```

---

### PvpPlayerToES —— 对局玩家列表聚合

消费对局玩家结算流水，聚合同一房间内所有玩家的信息（uid 列表、昵称、角色 ID、排名分等），写入 ES。

**执行拓扑**

```
[Source] pvpplayerendflow (master) ─┐
[Source] pvpplayerendflow (slave)  ─┼─► union ─► ProcessFunction ① ─► keyBy ②
                                                  (解析 PvpBatterInfo)  (env_roomId)
                                                                            │
                                                              KeyedProcessFunction ③
                                                              (ListState 聚合)
                                                              (2s Processing Time Timer)
                                                                            │
                                                                       ES Sink
```

**① ProcessFunction —— 流水解析**

每条流水对应一个玩家的结算数据，解析为 `PvpBatterInfo` POJO：

```java
// data.length > 74 才解析（字段足够）
String env        = data[74];
String uid        = data[18];
String campId     = data[20];   // 阵营 ID（0,1,2,…）
String roomId     = data[11];
String roleId     = data[27];   // 主角色
String roleId2    = data[54];   // 副角色2
String roleId3    = data[55];   // 副角色3
String nickname   = data[67];
String isAI       = data[21];
String rankOrder  = data[28];   // 本场排名
String showRankScore = data[68]; // 展示用排名分
```

**③ KeyedProcessFunction —— 2 秒 Timer 聚合 + 数据校验**

同一房间的多个玩家流水在 2s 内积累到 `ListState`，Timer 触发时做完整性校验：

```java
// 校验规则（任意不满足则标记为 isAI=true，视为无效数据）
// 1. 各阵营人数相等（campToUidNum 中所有 value 相同）
// 2. 各阵营角色总数相等（campToRoleIdNum 中所有 value 相同）
// 3. 角色总数 % 玩家总数 == 0
// 4. 无 AI 玩家（battle.getIsAI() != "1"）

// 通过校验后，按阵营顺序（0,1,2,…）排列玩家，写入 ES
json.put("key",              ctx.getCurrentKey());  // env_roomId
json.put("uids",             uidList.toArray());
json.put("nicknames",        nicknameList.toArray());
json.put("roleids",          roleIdList.toArray());
json.put("max_showrankscore", maxShowRankScore);
json.put("min_showrankscore", minShowRankScore);
json.put("rank_orders",      orderList.toArray());
json.put("player_ranks",     playerRankList.toArray());
json.put("camp_num",         campToUidNum.size());
json.put("pvp_player_success", "1");  // 标记玩家数据就绪
```

> **校验逻辑的意义**：若数据不完整（如部分玩家流水未到），胜率计算和视频推荐会使用错误的阵营分布。isAI 校验用于过滤与机器人对局，这类对局没有推荐价值。

---

### PvpRoleDamageToES —— 角色伤害/击杀聚合

消费对局角色结算流水，聚合同一房间内各角色的最高伤害和最高击杀得分，写入 ES。

**执行拓扑**

```
[Source] pvproleendflow (master) ─┐
[Source] pvproleendflow (slave)  ─┼─► union ─► ProcessFunction ① ─► keyBy(env_roomId) ─►
                                                                         KeyedProcessFunction ②
                                                                         (ListState 聚合)
                                                                         (2s Processing Time Timer)
                                                                              │
                                                                         ES Sink
```

**② KeyedProcessFunction —— 2 秒 Timer 聚合角色表现**

同一房间的角色流水积累 2s 后，对每个 roleId 取最大伤害和最大击杀分：

```java
// 状态: ListState<Tuple4<env_roomId, roleId, damage, score>>
// Timer 触发时：按 roleId 聚合，取 max(damage), max(score)

HashMap<String, long[]> rolePerformance = new HashMap<>();
for (Tuple4<...> record : records) {
    long[] cur = rolePerformance.getOrDefault(record.f1, new long[]{0L, 0L});
    cur[0] = Math.max(cur[0], record.f2);  // max damage
    cur[1] = Math.max(cur[1], record.f3);  // max kill score
    rolePerformance.put(record.f1, cur);
}

// 拼接结果字符串：roleId:score,damage+roleId:score,damage+...
StringBuilder sb = new StringBuilder();
for (Map.Entry<String, long[]> entry : rolePerformance.entrySet()) {
    sb.append(entry.getKey())
      .append(":").append(entry.getValue()[1])
      .append(",").append(entry.getValue()[0])
      .append("+");
}

HashMap<String, Object> json = new HashMap<>();
json.put("key",             env_roomId);
json.put("roleperformance", sb.toString());  // e.g. "101:980,45000+102:760,38000+"
```

---

### BattleDetailToRedis —— 对局详情缓存中转

消费视频上传流水，**延迟 3 秒**后将该对局的完整 ES 文档读出并缓存到 Redis。这是一个"中转 Job"，专门为 `PvpRankSort`（需要延迟 8s）和 `ComputerUserWinPercent`（需要延迟 4s）提供快速的 Redis 读取接口。

**执行拓扑**

```
[Source] uploadvideoflow (master) ─┐
[Source] uploadvideoflow (slave)  ─┼─► union ─► ProcessFunction ① ─► keyBy(roomId) ─►
                                                  (提取 env_roomId)
                                                               KeyedProcessFunction ②
                                                               (ValueState + 3s Timer)
                                                                        │
                                                               AsyncDataStream ③
                                                               (异步读 ES 文档)
                                                                        │
                                                               ProcessFunction ④
                                                               (写 Redis)
```

**② KeyedProcessFunction —— 3 秒防抖 Timer**

视频上传后等待 3s，让 `PvpToES`、`PvpPlayerToES`、`PvpRoleDamageToES` 等 Job 完成 ES 写入后再读取：

```java
private ValueState<String> roomIdState;

public void processElement(String s, Context context, Collector<String> out) {
    roomIdState.update(s);  // 记录 env_roomId
    long curTime = context.timerService().currentProcessingTime();
    context.timerService().registerProcessingTimeTimer(curTime + 3000);
}

public void onTimer(long timestamp, OnTimerContext ctx, Collector<String> out) {
    out.collect(roomIdState.value());  // 3s 后下发 env_roomId
    roomIdState.clear();
}
```

**③ AsyncDataStream —— 异步读取 ES 文档**

使用 `AsyncESQuery` 异步读取 `wedo_video_rec_online` 索引中该 key 的完整文档：

```java
AsyncDataStream.unorderedWait(
    timerStream,
    new AsyncESQuery(esIp, esPort, esIndex, esId),
    100, TimeUnit.SECONDS, 10
)
```

**④ ProcessFunction —— 写 Redis**

将 ES 文档（JSON 字符串）写入 Redis，设置 6 个月 TTL：

```java
// TTL = 60 * 60 * 24 * 30 * 6 秒 ≈ 6 个月
asyncCommands.setex(
    redisPrefix + key,         // wedo_video_online_details_{env_roomId}
    60 * 60 * 24 * 30 * 6,
    esDocJson                  // ES 文档的完整 JSON 字符串
);
```

**最终 Redis 写入格式：**

```
Key:   wedo_video_online_details_{env_roomId}
Type:  String
Value: {"key":"online_xxx","pvp_end_success":"1","pvp_player_success":"1","is_success":"1",...}
TTL:   15552000 秒（6 个月）
```

> **为什么用 Redis 而非直接读 ES？** 排行榜 Job（PvpRankSort）每处理一场对局都需要读取完整文档，若直接访问 ES 会产生大量随机读压力。将文档预先写入 Redis，后续读取延迟从毫秒级别降到微秒级别，且不影响 ES 的写入吞吐。

---

### VideoLoveToES —— 点赞数维护与点赞排行榜

这是对局推荐中最复杂的 Flink Job。它消费点赞流水，维护三个维度的数据：对局点赞总数（ES + Redis）、玩家点赞视频 Hash（Redis）、多维度点赞排行榜（Redis PriorityQueue）。

**执行拓扑**

```
[Source] gameplayrecordlikeflow (master) ─┐
[Source] gameplayrecordlikeflow (slave)  ─┼─► union ─► keyBy(env_roomId) ─►
                                                         TimeWindow(5s) ①
                                                         (汇总净点赞数)
                                                                │
                                                    AsyncRedisMGet ②
                                                    (读对局详情 + 现有点赞数)
                                                                │
                                                    ProcessFunction ③
                                                    (校验 + 分发)
                                                    ├── ES Sink: total_like_num
                                                    ├── Redis: wedo_video_online_like_{key}
                                                    ├── SideOutput: videoPlayerLikeInfo ──► 玩家点赞 Hash 分支
                                                    └── SideOutput: videoLikeInfo       ──► 点赞排行榜分支

[videoPlayerLikeInfo] ─► keyBy(uid) ─► TimeWindow(500ms) ─► AsyncRedisHGet ─► ProcessFunction ④
                                                                                (top-50 合并写 Hash)

[videoLikeInfo] ─► keyBy(groupByKey) ─► TimeWindow(5s) ─► ProcessFunction ⑤
                                                           (PriorityQueue top-200 + 写 Redis)
```

**① ProcessWindowFunction（5s ProcessingTime 窗口）—— 净点赞数计算**

5s 内同一对局的点赞/取消点赞事件合并，计算净变化量：

```java
// isLike="1" → +1（点赞）; isLike="0" → -1（取消点赞）
// 同一用户同一窗口内的多次操作取最后一次
int netLike = 0;
for (String element : elements) {
    String isLike = data[18];
    netLike += isLike.equals("1") ? 1 : -1;
}
// 输出: Tuple2(env_roomId, netLike)
```

**② AsyncRedisMGet —— 并行读取对局详情和现有点赞数**

同时读取两个 Redis Key（对局 JSON 缓存 + 已有点赞数），减少网络往返：

```java
// 读取：
// wedo_video_online_details_{key}  → 对局详情 JSON
// wedo_video_online_like_{key}     → 现有总点赞数
asyncCommands.mget(detailKey, likeKey)
```

**③ ProcessFunction —— 校验完整性 + 三路分发**

拿到 Redis 数据后，校验对局数据是否完整，然后同时输出三路数据：

```java
JsonObject detail = parseJson(redisResult.f0);

// 完整性校验：三个 success 标志全部存在
if (!detail.has("pvp_end_success") || !detail.has("pvp_player_success") || !detail.has("is_success")) {
    return;  // 数据未就绪，跳过
}

int totalLikeNum = parseInt(redisResult.f1) + netLike;

// 路径 A：写 ES 总点赞数
json.put("key", key);
json.put("total_like_num", totalLikeNum);
out.collect(json);

// 路径 B：写 Redis 点赞数
asyncCommands.set("wedo_video_online_like_" + key, String.valueOf(totalLikeNum));

// 路径 C：SideOutput 给玩家点赞 Hash 分支（每个参赛玩家一条记录）
String[] uids = gson.fromJson(detail.get("uids"), String[].class);
for (String uid : uids) {
    ctx.output(videoPlayerLikeTag, Tuple3.of(key, env + "_" + version + "_" + uid, totalLikeNum));
}

// 路径 D：SideOutput 给点赞排行榜分支（含角色列表、任务类型等排行 key 所需信息）
ctx.output(videoLikeTag, Tuple5.of(key, totalLikeNum, missionType, roles, metadata));
```

**④ ProcessFunction（玩家点赞 Hash）—— 维护每个玩家观看过的点赞视频 top-50**

500ms 窗口聚合同一 uid 的点赞变更，读取历史 Hash 后合并，只保留点赞数最高的前 50 条：

```java
// Redis Hash: wedo_video_online_like_hash_{uid}
// Field: env_roomId, Value: likeNum

// 读取现有 Hash → 合并新数据 → 按 likeNum 降序排序 → 截取 top-50
Map<String, String> existingHash = asyncRedisHGet.get(hashKey);
existingHash.put(key, String.valueOf(newLikeNum));

List<Map.Entry<String, String>> sorted = existingHash.entrySet().stream()
    .sorted((a, b) -> parseInt(b.getValue()) - parseInt(a.getValue()))
    .collect(Collectors.toList());

Map<String, String> top50 = new LinkedHashMap<>();
sorted.stream().limit(50).forEach(e -> top50.put(e.getKey(), e.getValue()));

asyncCommands.hmset("wedo_video_online_like_hash_" + uid, top50);
```

**⑤ ProcessWindowFunction（点赞排行榜）—— 多维度 top-200 排行榜**

5s 窗口内聚合同一分组维度的点赞变更，更新内存 PriorityQueue，持久化到 Redis：

```java
// 4 个分组维度（groupByKey）：
// env_version__                    （全局，无任务类型，无角色）
// env_version_missionType_         （按任务类型）
// env_version_missionType_role     （按任务类型+角色）
// env_version__role                （按角色）

// 每个 groupByKey 在 TaskManager 内维护一个 PriorityQueue（top-200）
// 首次启动时从 Redis 加载历史数据恢复
// 新数据入队 → 超过 200 条则剔除最小值

PriorityQueue<VideoItem> queue = getOrInitQueue(groupByKey);  // 从 Redis 初始化
queue.add(new VideoItem(key, likeNum));
if (queue.size() > 200) queue.poll();  // 移除最小值

// 写 Redis: wedo_video_online_{groupByKey}_sort
// 格式: "key1:likeNum1;key2:likeNum2;..." 按 likeNum 降序
asyncCommands.set("wedo_video_online_" + groupByKey + "_sort", serialize(queue));
```

**最终 Redis 写入格式：**

```
// 对局点赞总数
Key:   wedo_video_online_like_{env_roomId}
Type:  String
Value: "142"

// 玩家点赞视频 Hash（top-50）
Key:   wedo_video_online_like_hash_{uid}
Type:  Hash
Field: {env_roomId}
Value: {likeNum}

// 多维度点赞排行榜
Key:   wedo_video_online_{env}_{version}__sort
Key:   wedo_video_online_{env}_{version}_{missionType}_sort
Key:   wedo_video_online_{env}_{version}_{missionType}_{role}_sort
Key:   wedo_video_online_{env}_{version}__role_sort
Type:  String
Value: "online_room1:142;online_room2:98;online_room3:55;..."  （降序，分号分隔）
```

---

### PvpRankSort —— 排名排行榜

消费视频上传流水，**延迟 8 秒**后读取 Redis 中的对局缓存，根据对局最低排名分（`max_showrankscore`，代表本场最低竞技门槛）维护多维度 top-500 排行榜。

**延迟 8s 的原因**：需要等待 `BattleDetailToRedis`（延迟 3s 写完 Redis），再预留额外 5s 容忍网络和处理抖动。

**执行拓扑**

```
[Source] uploadvideoflow (master) ─┐
[Source] uploadvideoflow (slave)  ─┼─► union ─► ProcessFunction ① ─► keyBy(roomId) ─►
                                                  (提取 env_roomId)
                                                               KeyedProcessFunction ②
                                                               (ValueState + 8s Timer)
                                                                        │
                                                               AsyncRedisPvpDetail ③
                                                               (读 Redis 对局缓存)
                                                                        │
                                                               KeyedProcessFunction ④
                                                               (PriorityQueue 排行榜更新)
                                                               → Redis Sink
```

**② KeyedProcessFunction —— 8 秒防抖 Timer**

与 BattleDetailToRedis 类似，用 ValueState + Timer 引入延迟：

```java
private ValueState<String> roomIdState;

public void processElement(String s, Context context, Collector<String> out) {
    roomIdState.update(s);
    long curTime = context.timerService().currentProcessingTime();
    context.timerService().registerProcessingTimeTimer(curTime + 8000);
}
```

**③ AsyncRedisPvpDetail —— 读取 Redis 对局缓存**

```java
// param "2" 表示读 wedo_video_online_details_{roomId}（与 ComputerUserWinPercent 共用，
// param "1" 读同前缀但加版本环境 hash，param "2" 读 string 缓存）
AsyncDataStream.unorderedWait(timerStream, new AsyncRedisPvpDetail(host, port, password, env, "2"), ...)
```

**④ KeyedProcessFunction —— PriorityQueue 维护 top-500 排名排行榜**

```java
// 完整性校验：pvp_end_success + pvp_player_success + is_success 三标志均须存在
JsonObject detail = parseJson(pvpDetailJson);
if (!detail.has("pvp_end_success") || !detail.has("pvp_player_success") || !detail.has("is_success")) return;

int maxShowRankScore = detail.get("max_showrankscore").getAsInt();

// 4 个分组维度与 VideoLoveToES 一致
for (String groupByKey : groupByKeys) {
    PriorityQueue<VideoItem> queue = getOrInitFromRedis(groupByKey);  // 首次从 Redis 加载

    queue.add(new VideoItem(key, maxShowRankScore));
    if (queue.size() > 500) queue.poll();  // top-500，移除最小值

    asyncCommands.set("wedo_video_online_grade_" + groupByKey + "_sort", serialize(queue));
}

// 定期清理：每 24h 的 in-process Timer 触发，清除超过 24h 的过期数据
if (isFirstLoad) {
    registerCleanupTimer(currentTime + 24 * 60 * 60 * 1000);
}
```

**最终 Redis 写入格式：**

```
Key:   wedo_video_online_grade_{groupByKey}_sort
Type:  String
Value: "online_room1:8800;online_room2:7600;..."  （按 max_showrankscore 降序，分号分隔，top-500）
```

---

### ComputerUserWinPercent —— 玩家胜率计算

消费视频上传流水，**延迟 4 秒**后读取对局详情，计算参战玩家的历史胜率，并将**本场所有玩家中的最低胜率**写入 ES 和 Redis（用于标注视频的"最低竞技强度"）。

**执行拓扑**

```
[Source] uploadvideoflow (master) ─┐
[Source] uploadvideoflow (slave)  ─┼─► union ─► ProcessFunction ① ─► keyBy(roomId) ─►
                                                  (过滤 isAI)
                                                               KeyedProcessFunction ②
                                                               (ValueState + 4s Timer)
                                                                        │
                                                               AsyncRedisPvpDetail ③
                                                               (读 Redis 对局缓存, param "1")
                                                                        │
                                                               ProcessFunction ④
                                                               (计算每个玩家的胜负, 展开为多条)
                                                                        │
                                                               AsyncRedisPlayerRecordHGet ⑤
                                                               (读玩家历史胜负 Hash)
                                                                        │
                                                               ProcessFunction ⑥
                                                               (滑动窗口胜率计算 + 写 Redis Hash)
                                                               → emit Tuple8(roomId, ruid, 6档胜率...)
                                                                        │
                                                               keyBy(roomId) ─►
                                                               TimeWindow(2s) ─►
                                                               ProcessWindowFunction ⑦
                                                               (取所有玩家最低胜率)
                                                               → ES Sink + Redis Sink
```

**② KeyedProcessFunction —— 4 秒 Timer**

```java
public void processElement(String s, Context context, Collector<String> out) {
    roomId.update(s);
    context.timerService().registerProcessingTimeTimer(
        context.timerService().currentProcessingTime() + 4000
    );
}
public void onTimer(long timestamp, OnTimerContext ctx, Collector<String> out) {
    out.collect(roomId.value());
    roomId.clear();
}
```

**④ ProcessFunction —— 计算每个玩家的胜负**

从对局详情中提取阵营信息，判断每个玩家的胜负：

```java
String winCampId = detail.get("wincampid").toString();
String campNum   = detail.get("camp_num").toString();
String[] uidList = gson.fromJson(detail.getAsJsonArray("uids"), String[].class);

// 胜利阵营的 uid 范围：
// startIndex = (uidList.length / campNum) * winCampId
// endIndex   = startIndex + (uidList.length / campNum) - 1
int startIndex = (uidList.length / parseInt(campNum)) * parseInt(winCampId);
int endIndex   = startIndex + uidList.length / parseInt(campNum) - 1;

String commonField = env + "_" + version;
for (int i = 0; i < uidList.length; i++) {
    boolean isWin = (i >= startIndex && i <= endIndex);
    // 输出: Tuple3(env_version_ruid, isWin, eventTime_roomId)
    out.collect(Tuple3.of(commonField + "_" + uidList[i], isWin, eventTime + "_" + roomId));
}
```

**⑥ ProcessFunction —— 滑动窗口胜率计算**

读取玩家历史 Hash 后，维护最近 100 场 + 7 天内的滑动窗口，计算 6 个粒度的胜率：

```java
// 清理过期数据：
// 1. 删除超过 7 天的记录
// 2. 若记录 >= 100 条，再删除最旧一条
ArrayList<String> delRoleList = new ArrayList<>();
for (Map.Entry<String, String> entry : lastRecord.entrySet()) {
    long lastTime = parseLong(entry.getKey());  // field=eventTime（秒级时间戳）
    if (eventTime - lastTime > 7 * 24 * 60 * 60) {
        delRoleList.add(entry.getKey());
    }
}
if (remainRecord.size() >= 100) {
    delRoleList.add(String.valueOf(minTime));  // 额外删除最旧一条
}

// 更新 Redis Hash（field=eventTime, value="1"表示胜/"0"表示负）
asyncCommands.hdel(rankPrefix + ruid, delRoleList.toArray(new String[0]));
asyncCommands.hset(rankPrefix + ruid, String.valueOf(eventTime), isWin ? "1" : "0");

// 按时间降序遍历，计算各粒度胜率
SortedSet<String> descendingKeys = remainRecord.descendingKeySet();
int winNum = 0, failNum = 0;
double tenWinPercent = 0.0, twentyWinPercent = 0.0, thirtyWinPercent = 0.0;
double fortyWinPercent = 0.0, fiftyWinPercent = 0.0, allWinPercent = 0.0;
for (String key : descendingKeys) {
    if (remainRecord.get(key).equals("1")) winNum++; else failNum++;
    int total = winNum + failNum;
    if      (total == 10)  tenWinPercent    = winNum / 10.0;
    else if (total == 20)  twentyWinPercent = winNum / 20.0;
    else if (total == 30)  thirtyWinPercent = winNum / 30.0;
    else if (total == 40)  fortyWinPercent  = winNum / 40.0;
    else if (total == 50)  fiftyWinPercent  = winNum / 50.0;
    else                   allWinPercent    = winNum / (double)remainRecord.size();
}

// 输出: Tuple8(roomId, ruid, 10场胜率, 20场胜率, 30场胜率, 40场胜率, 50场胜率, 全部胜率)
collector.collect(Tuple8.of(roomId, ruid, ten, twenty, thirty, forty, fifty, all));
```

**⑦ ProcessWindowFunction（2s 窗口）—— 取最低胜率 + 双写输出**

2s 窗口聚合同一房间所有玩家的胜率，取每个粒度的**最小值**（代表本场视频的"最低竞技水平"），写入 ES 和 Redis：

```java
// 初始值 2.0（大于合法范围 [0,1]），对所有玩家取 min
double allWinPercent = 2.0;
for (Tuple8<...> element : elements) {
    allWinPercent = Math.min(allWinPercent, element.f7);
    // 同理计算其余 5 个粒度
}

// Redis: 取一位小数，最大值 0.9
double resultPercent = Math.floor(allWinPercent * 10) / 10;
if (resultPercent > 1.0) resultPercent = 0.9;
asyncCommands.set("wedo_video_online_win_min_percent_" + roomId, String.valueOf(resultPercent));

// ES: 写入 6 个胜率字段
HashMap<String, Object> playerWinPercent = new HashMap<>();
playerWinPercent.put("key", roomId);
playerWinPercent.put("10_win_percent",  minTenWinPercent);
playerWinPercent.put("20_win_percent",  minTwentyWinPercent);
playerWinPercent.put("30_win_percent",  minThirtyWinPercent);
playerWinPercent.put("40_win_percent",  minFortyWinPercent);
playerWinPercent.put("50_win_percent",  minFiftyWinPercent);
playerWinPercent.put("all_win_percent", allWinPercent);
out.collect(playerWinPercent);
```

**最终 Redis 写入格式：**

```
// 玩家历史胜负记录 Hash
Key:   wedo_video_online_win_percent_{env_version_ruid}
Type:  Hash
Field: {eventTime（秒级时间戳）}
Value: "1"（胜）或 "0"（负）
容量:  最多 100 条；自动清理 7 天前的数据

// 本场最低胜率（供推荐侧使用）
Key:   wedo_video_online_win_min_percent_{env_roomId}
Type:  String
Value: "0.6"（1 位小数，最大 0.9）
```

---

### 状态管理汇总

| Job | 状态描述符名称 | 类型 | TTL | 用途 |
|-----|--------------|------|-----|------|
| PvpPlayerToES | `wedo_video_online_members_user_battle` | `ListState<PvpBatterInfo>` | 无 | 聚合同一房间玩家流水，Timer 后 `clear()` |
| PvpPlayerToES | `wedo_video_online_members_is_register_new` | `ValueState<Boolean>` | 无 | 防 Timer 重复注册 |
| PvpRoleDamageToES | *(同 Player 类似)* | `ListState<Tuple4>` + `ValueState<Boolean>` | 无 | 聚合同一房间角色伤害流水 |
| BattleDetailToRedis | `(prefix)_room_id_state` | `ValueState<String>` | 无 | 3s Timer 防抖，存储 env_roomId |
| PvpRankSort | `(prefix)_room_id_state` | `ValueState<String>` | 无 | 8s Timer 防抖 |
| ComputerUserWinPercent | `(rankPrefix)_temp_data` | `ValueState<String>` | 无 | 4s Timer 防抖 |
| VideoLoveToES | *(无 Flink 托管状态)* | — | — | 使用 Processing Time Window 聚合 |

> **所有 ValueState 和 ListState 均在 Timer 触发后主动 `clear()`**，正常情况下占用时间短暂（2~8s）。若 Job 异常重启导致 Timer 丢失，State 可能积压，需关注 TaskManager 堆内存水位。

---

### 双 Kafka 容灾模式

所有 Job 均采用相同的容灾模式：

```java
FlinkKafkaConsumer<String> consumer1 = KafkaUtil.initKafka(kafkaMaster, topic, groupID);
FlinkKafkaConsumer<String> consumer2 = KafkaUtil.initKafka(kafkaSlave,  topic, groupID);

env.addSource(consumer1)
   .union(env.addSource(consumer2))
   ...
```

| Kafka 集群 | 地址 | 角色 |
|-----------|------|------|
| Master | `199-kafka-clb.kpmq.woa.com:9092` | 主集群（CLB 地址） |
| Slave  | `204-kafka-clb.kpmq.woa.com:9092` | 备集群（CLB 地址） |

两路 Consumer 使用**相同的 `groupID`**，重复消费时通过以下机制保证幂等性：

| Job | 幂等保障机制 |
|-----|------------|
| PvpToES | ES 按 `key` upsert；多余字段覆盖但不丢失原有字段 |
| PvpPlayerToES | ES upsert；同一 roomId 在 2s 窗口内幂等（Timer 只注册一次） |
| PvpRoleDamageToES | ES upsert；同一 roomId 在 2s 窗口内幂等 |
| BattleDetailToRedis | Redis `setex` 覆盖写；同一 key 多次写入结果一致 |
| PvpRankSort / ComputerUserWinPercent | Redis PriorityQueue / Hash 均为幂等操作；同一 roomId 多次写入结果稳定 |
| VideoLoveToES | 窗口内 int 累加；Redis `set` 覆盖写 |

---

## 存储设计

### ES 索引分阶段写入机制

`wedo_video_rec_online` 索引由 4 个 Job 分阶段写入，最终组成完整文档：

| 写入阶段 | Job | 写入字段 | 就绪标志 |
|---------|-----|---------|---------|
| 阶段 1 | PvpToES（pvpendflow） | roomid, fubenid, score, wincampid, mmr, version, game_play_record_tags, … | `pvp_end_success="1"` |
| 阶段 1 | PvpToES（uploadvideoflow） | is_success, is_ai, rand | `is_success="1"` |
| 阶段 2 | PvpPlayerToES | uids[], nicknames[], roleids[], max/min_showrankscore, player_ranks[], camp_num | `pvp_player_success="1"` |
| 阶段 3 | PvpRoleDamageToES | roleperformance | — |
| 阶段 4 | ComputerUserWinPercent | 10/20/30/40/50/all_win_percent | — |
| 实时更新 | VideoLoveToES | total_like_num | — |

> **完整性校验逻辑**：`BattleDetailToRedis`、`PvpRankSort`、`ComputerUserWinPercent`、`VideoLoveToES` 均会检查三个 success 标志（`pvp_end_success + pvp_player_success + is_success`）是否都存在，确保只处理数据完整的对局。

### Redis 键空间设计

| 键模板 | 类型 | 说明 | TTL |
|--------|------|------|-----|
| `wedo_video_online_details_{env_roomId}` | String | ES 文档完整 JSON 缓存，供 PvpRankSort/ComputerUserWinPercent 快速读取 | 6 个月 |
| `wedo_video_online_like_{env_roomId}` | String | 对局视频总点赞数 | 无（持久）|
| `wedo_video_online_like_hash_{uid}` | Hash | 玩家点赞过的视频 top-50，Field=env_roomId，Value=likeNum | 无（持久）|
| `wedo_video_online_{groupByKey}_sort` | String | 点赞排行榜，分号分隔，按点赞数降序，top-200 | 无（持久）|
| `wedo_video_online_grade_{groupByKey}_sort` | String | 排名排行榜，按 max_showrankscore 降序，top-500 | 无（持久）|
| `wedo_video_online_win_percent_{env_version_ruid}` | Hash | 玩家历史胜负记录，Field=时间戳，Value="1"/"0" | 无（代码层面维护 100条/7天窗口）|
| `wedo_video_online_win_min_percent_{env_roomId}` | String | 本场所有玩家的最低胜率，1 位小数，上限 0.9 | 无（持久）|

### Elasticsearch 索引设计

```
索引: wedo_video_rec_online
主键: key（= env_roomId）

关键字段:
  - key (keyword)              ← 唯一主键（env_roomId）
  - env (keyword)              ← 环境标识
  - version (keyword)          ← 版本号
  - roomid (keyword)           ← 对局房间 ID
  - fubenid (keyword)          ← 关卡 ID
  - fubenhardlv (integer)      ← 关卡难度
  - score (long)               ← 胜方分数
  - score_sum (long)           ← 总分数
  - wincampid (keyword)        ← 胜利阵营
  - mmr (integer)              ← 对局 MMR
  - start_time (keyword)       ← 对局开始时间
  - game_play_record_tags (keyword[]) ← 标签数组
  - pvp_end_success (keyword)  ← 对局结束数据就绪标志
  - is_success (keyword)       ← 视频上传成功标志
  - is_ai (keyword)            ← 是否 AI 对局
  - rand (integer)             ← 随机分桶 [0,20)
  - uids (keyword[])           ← 参赛玩家 ruid 列表（按阵营顺序）
  - nicknames (keyword[])      ← 昵称列表
  - roleids (integer[])        ← 角色 ID 列表
  - max_showrankscore (integer) ← 本场最高排名分
  - min_showrankscore (integer) ← 本场最低排名分
  - player_ranks (integer[])   ← 各玩家排名分列表
  - rank_orders (integer[])    ← 各玩家本场排名
  - camp_num (integer)         ← 阵营数量
  - pvp_player_success (keyword) ← 玩家数据就绪标志
  - roleperformance (keyword)  ← 角色表现，格式 "roleId:score,damage+..."
  - 10_win_percent (double)    ← 近 10 场最低胜率
  - 20_win_percent (double)    ← 近 20 场最低胜率
  - 30_win_percent (double)    ← 近 30 场最低胜率
  - 40_win_percent (double)    ← 近 40 场最低胜率
  - 50_win_percent (double)    ← 近 50 场最低胜率
  - all_win_percent (double)   ← 全部场次最低胜率
  - total_like_num (integer)   ← 总点赞数
```

---

## 部署与运维

### Flink 任务列表

| Job 名称 | 入口类 | 消费 Topic |
|---------|--------|----------|
| `wedo_pvp_end_to_es` | PvpToES | `billow_ex_wedo_pvpendflow` + `billow_ex_wedo_uploadvideoflow` |
| `wedo_pvp_player_to_es` | PvpPlayerToES | `billow_ex_wedo_pvpplayerendflow` |
| `wedo_pvp_roleid_damage_to_es` | PvpRoleDamageToES | `billow_ex_wedo_pvproleendflow` |
| `wedo_pvp_battle_detail_redis` | BattleDetailToRedis | `billow_ex_wedo_uploadvideoflow` |
| `wedo_pvp_win_percent` | ComputerUserWinPercent | `billow_ex_wedo_uploadvideoflow` |
| `wedo_pvp_rank_sort` | PvpRankSort | `billow_ex_wedo_uploadvideoflow` |
| `wedo_pvp_like_num_sort` | VideoLoveToES | `billow_ex_wedo_gameplayrecordlikeflow` |

### Kafka Topic 速查

| Topic | 用途 | 消费 Flink 任务 |
|-------|------|----------------|
| `billow_ex_wedo_pvpendflow` | 对局结束事件（含对局元信息） | PvpToES |
| `billow_ex_wedo_pvpplayerendflow` | 对局玩家结算事件 | PvpPlayerToES |
| `billow_ex_wedo_pvproleendflow` | 对局角色结算事件 | PvpRoleDamageToES |
| `billow_ex_wedo_uploadvideoflow` | 视频上传成功事件 | PvpToES / BattleDetailToRedis / ComputerUserWinPercent / PvpRankSort |
| `billow_ex_wedo_gameplayrecordlikeflow` | 视频点赞/取消点赞事件 | VideoLoveToES |

### Flink 算子速查

| 算子 | 使用 Job | 用途 |
|------|---------|------|
| `addSource().union()` | 全部 | 双 Kafka 容灾 |
| `process(ProcessFunction)` | 全部 | TLog 管道符字段解析，提取关键字段 |
| `keyBy(KeySelector)` | 全部 | 按 env_roomId / uid / groupByKey 分区 |
| `process(KeyedProcessFunction)` + `ValueState` + `Timer(2/3/4/8s)` | PvpPlayerToES / PvpRoleDamageToES / BattleDetailToRedis / ComputerUserWinPercent / PvpRankSort | 防抖 Timer 引入固定延迟 |
| `process(KeyedProcessFunction)` + `ListState` | PvpPlayerToES / PvpRoleDamageToES | 聚合同一房间多玩家/多角色数据 |
| `AsyncDataStream.unorderedWait(AsyncESQuery)` | BattleDetailToRedis | 异步读取 ES 文档，缓存到 Redis |
| `AsyncDataStream.unorderedWait(AsyncRedisPvpDetail)` | ComputerUserWinPercent / PvpRankSort | 异步读取 Redis 对局缓存 |
| `AsyncDataStream.unorderedWait(AsyncRedisPlayerRecordHGet)` | ComputerUserWinPercent | 异步读取玩家历史胜负 Hash |
| `AsyncDataStream.unorderedWait(AsyncRedisMGet)` | VideoLoveToES | 并行读取对局详情 + 现有点赞数 |
| `timeWindow(5s, ProcessingTime)` + `ProcessWindowFunction` | VideoLoveToES / PvpRankSort | 批量聚合点赞变更；多维度排行榜更新 |
| `timeWindow(2s, ProcessingTime)` + `ProcessWindowFunction` | ComputerUserWinPercent | 合并同房间所有玩家胜率，取最小值 |
| `timeWindow(500ms, ProcessingTime)` + `ProcessWindowFunction` | VideoLoveToES | 合并同 uid 的点赞变更 |
| `addSink(EsSink)` | PvpToES / PvpPlayerToES / PvpRoleDamageToES / ComputerUserWinPercent / VideoLoveToES | ES upsert 写入各阶段字段 |
