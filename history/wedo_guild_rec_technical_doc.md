# WeDo 战队推荐系统技术实现文档（实时流计算部分）


---

## 需求背景

### 业务场景

WeDo 平台的**战队（Guild）** 是另一类核心社交单元，玩家可以创建或加入战队，与战友长期组队参与 PvP 对局。与车队（临时组队）不同，战队具有固定成员关系、等级体系和荣誉积分，社交黏性更强。

实时流计算层的目标是：**及时将游戏服务器产生的战队事件（创建、加入、退出、销毁、改名等）和玩家状态变化（登录/登出）同步到存储系统**，为推荐服务提供准确的战队信息底层数据。

### 核心挑战

| 挑战 | 描述 |
|------|------|
| **事件类型多** | 战队变更事件包含 8 种操作类型（加入/退出/职务变更/信息变更/升级/踢出/销毁/改名），需分类处理 |
| **在线人数实时性** | 战队在线成员数依赖登录/登出/加入/退出四路事件的联合计算，需事件时间对齐 |
| **屏蔽列表增量维护** | 用户对战队推荐的屏蔽/恢复操作分散在多条流水中，需在时间窗口内合并后增量更新 |
| **好友关系去重** | 登录时同一用户可能上报游戏好友和平台好友两条流水，需在 2s 内合并去重 |
| **数据一致性** | 退出战队时需先查 Redis 确认用户当前所在战队，再修改在线成员集合 |

### 技术选型

与车队推荐相同，战队推荐同样采用：

- **Flink**：处理四路 Kafka 流水，构建好友关系链、战队 ES 索引和在线成员状态
- **Elasticsearch**：存储战队结构化信息，索引名 `wedo_guild_online`，主键 `guild_id`
- **Redis (TendisPlus)**：缓存好友关系哈希表、用户战队归属、屏蔽列表

---

## Flink 实时流计算详解

### 整体 Job 架构

战队推荐涉及三个核心 Flink Job，分工明确：

| Job 名称 | 入口类 | 输出目标 | 作用 |
|---------|--------|---------|------|
| `wedo_friend_online` | FriendsOperator | Redis Hash | 实时维护用户好友关系链，为推荐服务提供"哪些好友在哪支战队"的基础数据 |
| `wedo_guild_online` | LogOperator | Redis String + Elasticsearch | 实时同步战队信息到 ES，并维护在线成员状态 |
| `wedo_guild_test_audit_list` | GuildAudit | Redis String | 实时维护用户对战队推荐的屏蔽列表 |

代码模块结构：

```
bus/wedo/
└── guild/
    ├── online/
    │   ├── FriendsOperator.java    ← 好友关系链（核心）
    │   ├── LogOperator.java        ← 战队信息同步 + 在线成员统计（核心）
    │   └── GuildAudit.java         ← 战队推荐屏蔽列表维护
    └── redis/
        ├── AsyncRedisGet.java      ← 异步 GET 封装（登出时查用户战队归属）
        └── AsyncRedisGetFun.java   ← 异步 GET + 合并封装（屏蔽列表读取）
```

所有 Job 均采用**主从双 Kafka + `.union()` 合并**容灾模式，任意一路断流不影响数据处理。

---

### FriendsOperator —— 好友关系链维护

消费玩家每日首次登录上报的好友列表流水，将其写入 Redis Hash，为推荐服务提供"用户有哪些好友"的数据基础。

**执行拓扑**

```
[Source] playerfriendslistflow (master) ─┐
[Source] playerfriendslistflow (slave)  ─┼─► union ─► flatMap ① ─► keyBy ②
                                                       (TLog解析)  (prefix+rUid+friendType)
                                                                        │
                                                                        ▼
                                                         KeyedProcessFunction ③
                                                         (ListState + ValueState)
                                                         (2s Processing Time Timer)
                                                                        │
                                                                        ▼
                                                         asyncCommands.hset() / hdel()
```

**算子逐一详解**

**① flatMap（RichFlatMapFunction）—— TLog 流水解析**

TLog 流水以 `|` 分隔，固定字段索引提取关键信息。好友列表为 `;` 分隔的条目列表，每个条目以 `,` 分隔，其中第二个元素（索引 1）为好友 ruid：

```java
.flatMap(new RichFlatMapFunction<String, Tuple3<String, String, String>>() {
    public void flatMap(String s, Collector<Tuple3<...>> collector) {
        String[] info = s.split("\\|");
        if (info.length > 24) {
            String rUid       = info[12];
            String friendType = info[21];  // "1"=游戏好友, "2"=平台好友
            String friends    = info[24].replaceAll("\n|\r", "");
            // 只处理游戏好友(1)和平台好友(2)，过滤空列表
            if (!CommonUtil.isEmpty(friends) && (friendType.equals("1") || friendType.equals("2"))) {
                StringBuilder friendBuilder = new StringBuilder();
                for (String friendInfo : friends.split(";")) {
                    friendBuilder.append(friendInfo.split(",")[1]).append(","); // 提取 friendRuid
                }
                friendBuilder.deleteCharAt(friendBuilder.length() - 1);
                collector.collect(Tuple3.of(friendListPrefix + rUid, friendBuilder.toString(), friendType));
            }
        }
    }
})
```

> **与车队推荐的差异**：车队推荐的好友链包含亲密度、组队时间等多维信号；战队推荐仅保存好友 ruid 集合，数据结构更简洁，因为战队推荐的排序信号主要来源于在线成员状态，而非亲密度。

**② keyBy —— 按 `prefix+rUid+"|"+friendType` 分区**

与车队推荐不同，战队推荐将**同一用户的不同好友类型路由到不同子任务**，使 Timer 聚合逻辑天然按类型隔离：

```java
.keyBy(data -> data.f0 + "|" + data.f2)
// 例："wedo_guild_online_friendlist_576465154108032693|1"（游戏好友）
// 例："wedo_guild_online_friendlist_576465154108032693|2"（平台好友）
```

**③ KeyedProcessFunction + Timer —— 2 秒微批聚合 + 写 Redis**

同一用户同一好友类型的多条流水在 2s 内合并，用 `HashSet` 自动去重后写入 Redis：

```java
// 状态定义（无显式 TTL）
ValueState<Boolean> isRegister;       // Timer 防重复注册
ListState<String>   friendListState;  // 累积同次登录的多条好友数据

public void processElement(Tuple3<...> data, Context context, Collector<Object> collector) {
    if (isRegister.value() == null) {
        context.timerService().registerProcessingTimeTimer(
            context.timerService().currentProcessingTime() + 2000
        );
        isRegister.update(true);
    }
    friendListState.add(data.f1);  // 追加好友 ruid 列表（逗号分隔字符串）
}

public void onTimer(long timestamp, OnTimerContext ctx, Collector<Object> out) {
    HashSet<String> uidList = new HashSet<>();
    for (String element : friendListState.get()) {
        if (!CommonUtil.isEmpty(element)) {
            uidList.addAll(Arrays.asList(element.split(",")));
        }
    }
    String[] info = ctx.getCurrentKey().split("\\|");
    String key        = info[0];  // wedo_guild_online_friendlist_{rUid}
    String friendType = info[1];  // "1" 或 "2"

    if (uidList.size() == 0) {
        asyncCommands.hdel(key, friendType);   // 无好友则删除该类型字段
    } else {
        asyncCommands.hset(key, friendType, Strings.join(uidList, ","));
    }
    isRegister.clear();
    friendListState.clear();
}
```

**最终 Redis 写入格式：**

```
Key:   wedo_guild_online_friendlist_{rUid}
Type:  Hash
Field: {friendType}    ("1" 或 "2")
Value: {ruid1},{ruid2},{ruid3},...

示例：
HGETALL wedo_guild_online_friendlist_576465154108032693
→ "1" → "123456789,987654321,111222333"   （游戏好友）
→ "2" → "444555666,777888999"              （平台好友）
```

> **与车队推荐对比**：车队好友链的 Hash Field 是 `friendRuid`（每个好友一条记录），战队好友链的 Hash Field 是 `friendType`（每种类型一条记录，value 为逗号分隔的 ruid 列表）。车队方案适合精细化亲密度查询；战队方案结构更紧凑，适合快速判断"有哪些好友在战队中"。

---

### LogOperator —— 战队信息同步 + 在线成员统计

这是战队推荐最复杂的 Flink Job。它同时消费**四路 Kafka 数据源**，负责两件事：

1. **战队信息同步到 ES**：处理战队创建和变更流水，维护 `wedo_guild_online` ES 索引
2. **在线成员实时统计**：联合登录/登出/加入/退出四路事件，实时更新每支战队的在线人数

**执行拓扑**

```
[Source] playerlogin  (master/slave) ─► ProcessFunction ①
                                        (写 Redis: rUid→guildId)
                                        (emit: rUid, guildId, 1=上线, ts)
                                                    │
[Source] playerlogout (master/slave) ─► ProcessFunction ②          ──────────────────────────┐
                                        (emit: rUid, key)                                     │
                                              │                                               │
                                      AsyncDataStream.unorderedWait(AsyncRedisGet)            │
                                        (查询 Redis: rUid → guildId)                          │
                                      ProcessFunction ③                                       │
                                        (emit: rUid, guildId, 0=下线, ts)                     │
                                                    │                                         │
[Source] modifyguild  (master/slave) ─► keyBy(guildId)                                       │
                                        KeyedProcessFunction ④                                │
                                          (事件类型分发)                                        │
                                          ├── type=1(加入)  → Redis: rUid→guildId            │
                                          │                 → SideOutput(成员变更)             │
                                          ├── type=2/6(退出/踢出) → Redis: rUid→"0"          │
                                          │                       → SideOutput(成员变更)       │
                                          ├── type=7(销毁) → 写 ES: member_size=0             │
                                          └── 其他变更   → 写 ES: 全量战队字段               │
                                                    │                                         │
[Source] createguild  (master/slave) ─► ProcessFunction ⑤                                    │
                                          (写 ES: 全量战队字段)                                │
                                          (写 Redis: rUid→guildId)                            │
                                          (SideOutput: 创建者上线)                            │
                                                    │                                         │
                                    ┌───────────────┼──────────────────────────────────────── ┘
                                    │  union 四路事件流（登录/登出/成员变更/战队创建）
                                    │  assignTimestampsAndWatermarks (lag=3s)
                                    │  keyBy(guildId)
                                    │  timeWindow(5s, EventTime)
                                    │  ProcessWindowFunction ⑥
                                    │  (MapState 维护在线成员集合)
                                    │  → 写 ES: {guild_id, online_num}
                                    ▼
                              ES Sink (wedo_guild_online)
```

**算子逐一详解**

**① ProcessFunction —— 登录事件处理**

玩家登录时，将其战队归属写入 Redis，并向在线成员流 emit 上线信号：

```java
// 过滤条件：info[34] == "0" 表示非重复登录（首次有效登录）
if (!info[34].equals("0")) return;

String rUid    = info[14];
String guildId = info[89];  // 当前所属战队 ID，"0" 表示未加入战队

// 无论是否在战队，都更新 Redis（确保后续登出时能查到最新归属）
asyncCommands.set(guildIdPrefix + rUid, guildId);

// 只有在战队中的玩家才向下游 emit 在线信号
if (!guildId.equals("0")) {
    out.collect(Tuple4.of(rUid, guildId, 1 /*上线*/, timeStamp));
}
```

**② → ③ 登出事件 + 异步 Redis 查询**

登出流水中不包含玩家当前所在的战队 ID，需先查 Redis 获取归属关系：

```java
// ② 登出流：提取 rUid，查询 Redis Key
DataStream<Tuple2<String, String>> logoutSource = ...process(s -> {
    String rUid = info[14];
    out.collect(Tuple2.of(guildIdPrefix + rUid, s));  // f0=Redis Key, f1=原始日志
}).rescale();

// 异步查询 Redis（超时 400s，最大并发 100）
AsyncDataStream.unorderedWait(logoutSource, new AsyncRedisGet(...), 400, TimeUnit.SECONDS, 100)

// ③ 拿到 Redis 结果后：用 value.f0（Redis GET 返回值）作为 guildId
.process(value -> {
    String guildId = value.f0;  // Redis 返回的 guildId
    if (!isEmpty(guildId) && !guildId.equals("0")) {
        out.collect(Tuple4.of(rUid, guildId, 0 /*下线*/, timeStamp));
    }
});
```

> **为什么登出需要查 Redis？** 登出流水中不携带战队 ID，但计算"战队在线人数"时需要知道"这个人离开了哪个战队"。通过 Redis 的 `wedo_guild_online_join_{rUid}` 维护用户当前战队归属，实现登出事件与战队的关联。

**④ KeyedProcessFunction —— 战队信息变更分类处理**

消费 `modifyguildflow`，按 `guildId`（字段索引 13）keyBy 后，按事件类型分发处理：

```java
// 战队等级 → 最大成员数映射（level i 对应上限 i+20 人）
// level 1 → 21人, level 2 → 22人, ..., level 10 → 30人
HashMap<Integer, Integer> guildLevelToNum;

public void processElement(String s, Context ctx, Collector<HashMap<String, Object>> out) {
    // 过滤无效事件：非成功（data[11] != "0"）或战队等级非法（<1 或 >10）
    if (!data[11].equals("0") || guildLevel < 1 || guildLevel > 10) return;

    String type = data[12]; // 事件类型
    // type=1: 成员加入; 2: 成员退出; 3: 职务变更; 4: 信息变更
    // type=5: 战队升级; 6: 踢出成员; 7: 战队销毁; 8: 修改名字

    if (type.equals("1")) {
        // 成员加入：更新 Redis 归属 + SideOutput 触发在线成员更新
        asyncCommands.set(guildIdPrefix + desRUid, guildId);
        ctx.output(memberChangeOutputTag, Tuple4.of(desRUid, guildId, 1, timeStamp));
    }
    if (type.equals("2") || type.equals("6")) {
        // 成员退出/踢出：清空 Redis 归属 + SideOutput 触发在线成员更新
        asyncCommands.set(guildIdPrefix + desRUid, "0");
        ctx.output(memberChangeOutputTag, Tuple4.of(desRUid, guildId, 2/*退出*/, timeStamp));
    }

    if (!type.equals("7")) {
        // 非销毁事件：写全量战队信息到 ES
        int isFull = (guildLevelToNum.get(guildLevel) == memberInfoList.length) ? 1 : 0;
        json.put("guild_id", guildId);
        json.put("guild_level", guildLevel);
        json.put("guild_name", data[15]);
        json.put("member_size", memberInfoList.length);
        json.put("member_list", rUidBuilder.toString());  // 逗号分隔的成员 ruid 串
        json.put("is_full", isFull);
        json.put("join_limit_type", ...);
        json.put("join_limit_value", ...);
        json.put("guild_combat_power", guildCombatPower);
        json.put("guild_audit_type", guildAuditType);
        json.put("env_id", envId);
        json.put("guild_br_rank_score", ...);
        json.put("guild_br_rank_season_id", ...);
        // ... 其他字段
        out.collect(json);
    } else {
        // 销毁事件：仅更新 member_size=0
        json.put("guild_id", guildId);
        json.put("member_size", 0);
        out.collect(json);
    }
}
```

**⑤ ProcessFunction —— 战队创建处理**

战队创建事件写入全量字段（含 `create_time`、`partition` 随机分桶等），并将创建者的战队归属记录到 Redis：

```java
// 创建成功（data[11] == "0"）才处理
int random = new Random().nextInt(50);  // [0,50) 随机分桶，用于推荐侧负载均衡
json.put("guild_id", guildId);
json.put("create_time", createTime);
json.put("partition", random);          // 战队特有字段，变更事件中无此字段
json.put("member_size", 1);
json.put("online_num", 1);
// ... 其他字段与变更事件一致

asyncCommands.set(guildIdPrefix + rUid, guildId);  // 创建者战队归属
ctx.output(createGuildOutputTag, Tuple4.of(rUid, guildId, 1, createTime));  // 创建者上线
```

**⑥ ProcessWindowFunction —— 5 秒事件时间窗口统计在线成员**

将四路事件流（登录、登出、成员加入/退出、战队创建）合并，以事件时间对齐，在 5s 窗口内统计每支战队的在线成员集合：

```java
// Event Time 水位线配置：3s 容忍乱序延迟
.assignTimestampsAndWatermarks(new AssignerWithPeriodicWatermarks<...>() {
    private final long maxTimeLag = 3000;
    public Watermark getCurrentWatermark() {
        return new Watermark(curMaxTime - maxTimeLag);
    }
})

// keyBy guildId → 5s EventTime 窗口
.keyBy(data -> data.f1)
.timeWindow(Time.seconds(5))
.process(new ProcessWindowFunction<...>() {
    private MapState<String, Integer> guildIdMap;  // rUid → 是否在线（1/移除）

    public void process(String guildId, Context ctx,
                        Iterable<Tuple4<...>> iterable,
                        Collector<HashMap<String, Object>> out) {
        for (Tuple4<...> data : iterable) {
            if (data.f2 == 1) {        // 上线 / 加入
                guildIdMap.put(data.f0, 1);
            } else {                   // 下线 / 退出
                guildIdMap.remove(data.f0);
            }
        }
        // 统计在线成员数并写 ES（partial update）
        int onlineNum = 0;
        for (String uid : guildIdMap.keys()) onlineNum++;
        HashMap<String, Object> value = new HashMap<>();
        value.put("guild_id", guildId);
        value.put("online_num", onlineNum);
        out.collect(value);
    }
})
.addSink(ESUtil.getESSinkBuilder(..., "guild_id").build());
```

> **在线成员统计使用 EventTime 而非 ProcessingTime 的原因**：四路事件来自不同 Kafka Topic，网络延迟和消费速度可能不同，使用 EventTime + 3s Watermark 可以容忍轻微乱序，保证同一时间段内的上线/下线事件被正确聚合到同一个窗口中。

**ES 文档结构示例：**

```json
{
  "guild_id": "guild_456",
  "guild_level": 5,
  "guild_name": "胜利战队",
  "plat": 1,
  "env_id": 34,
  "world_id": 1001,
  "member_size": 18,
  "member_list": "ruid_001,ruid_002,ruid_003",
  "is_full": 0,
  "online_num": 7,
  "join_limit_type": 1,
  "join_limit_value": 1000,
  "head_ruid": "ruid_001",
  "guild_slogan": "战无不胜",
  "guild_logo": "logo_001",
  "guild_combat_power": 88000,
  "guild_audit_type": 0,
  "guild_br_rank_score": 12500,
  "guild_br_rank_season_id": 3,
  "create_time": 1718250000,
  "partition": 23
}
```

---

### GuildAudit —— 战队推荐屏蔽列表维护

消费玩家对战队推荐的屏蔽/恢复流水，将每个用户的屏蔽战队列表维护到 Redis。推荐服务在召回阶段会过滤掉用户已屏蔽的战队。

**执行拓扑**

```
[Source] blockguildrecommendflow (master) ─┐
[Source] blockguildrecommendflow (slave)  ─┼─► union ─► keyBy ① ─► TimeWindow(5s) ② ─► AsyncDataStream ③ ─► ProcessFunction ④
                                                        (rUid)      (批量合并增删)       (异步 Redis GET         (合并历史数据
                                                                                          读取现有屏蔽列表)         → 写 Redis)
```

**算子逐一详解**

**① keyBy —— 按 rUid（字段索引 14）分区**

将同一用户的屏蔽操作路由到同一子任务，确保同一窗口内可见：

```java
.keyBy(data -> data.split("\\|")[14])
```

**② ProcessWindowFunction（5s ProcessingTime 窗口）—— 增删集合计算**

在 5s 窗口内聚合同一用户的多次屏蔽/恢复/清空操作，输出增量变化（`addSet`、`delSet`）和是否触发了清空（`delAllTime`）：

```java
// 事件类型
// type="0"：屏蔽某个战队（加入屏蔽列表）
// type="1"：取消屏蔽（从屏蔽列表移除）
// type="2"：清空所有屏蔽

.process(new ProcessWindowFunction<...>() {
    public void process(String rUid, Context ctx, Iterable<String> elements,
                        Collector<Tuple2<String, String>> out) {
        HashSet<String> addSet = new HashSet<>();
        HashSet<String> delSet = new HashSet<>();
        long delAllTime = 0L; // 最新一次清空操作的时间戳

        for (String element : elements) {
            String type      = data[17];
            long eventTime   = Long.parseLong(data[19]);
            String guildId   = data[18];
            if (type.equals("2")) {
                delAllTime = Math.max(delAllTime, eventTime);  // 记录最新清空时间
            }
            allData.add(Tuple3.of(guildId, type, eventTime));
        }

        // 只处理清空操作之后发生的增删事件（忽略早于清空的历史记录）
        for (Tuple3<...> temp : allData) {
            if (temp.f2 > delAllTime) {
                if (temp.f1.equals("0")) {       // 屏蔽：加入 addSet（若在 delSet 中则移除）
                    if (delSet.contains(temp.f0)) delSet.remove(temp.f0);
                    else addSet.add(temp.f0);
                } else if (temp.f1.equals("1")) { // 取消屏蔽：加入 delSet（若在 addSet 中则移除）
                    if (addSet.contains(temp.f0)) addSet.remove(temp.f0);
                    else delSet.add(temp.f0);
                }
            }
        }

        // 下发格式："{adds}|{dels}|{isClearAll}"
        String res = join(addSet) + "|" + join(delSet) + "|" + (delAllTime != 0 ? "1" : "0");
        out.collect(Tuple2.of(guildAuditPrefix + rUid, res));
    }
})
```

> **清空操作的时序处理**：如果窗口内同时存在"清空"和后续的"屏蔽新战队"，`delAllTime` 机制确保只保留清空之后的增量操作，避免已清空的战队重新出现在屏蔽列表中。

**③ AsyncDataStream.unorderedWait —— 异步读取现有屏蔽列表**

写入前需先读取 Redis 中用户的现有屏蔽列表，以便做合并计算（超时 50s，最大并发 20）：

```java
AsyncDataStream.unorderedWait(
    guildAuditStream,
    new AsyncRedisGetFun(host, port, password),  // 异步 GET 当前列表
    50, TimeUnit.SECONDS,
    20
)
```

**④ ProcessFunction —— 合并历史数据 + 写 Redis**

拿到 Redis 现有数据后，与窗口计算结果合并，最终写回：

```java
// data.f0 = Redis Key
// data.f1 = "{adds}|{dels}|{isClearAll}"
// data.f2 = Redis 现有值（AsyncRedisGetFun 回填）
String lasts = data.f2;  // 现有屏蔽列表
if (operators[2].equals("1")) lasts = "";  // 若有清空操作，忽略历史数据

HashSet<String> guilds = new HashSet<>(Arrays.asList(lasts.split(",")));
guilds.addAll(Arrays.asList(adds.split(",")));    // 合并新增
for (String del : dels.split(",")) guilds.remove(del); // 移除已取消

if (guilds.size() > 0) {
    asyncCommands.set(key, Strings.join(guilds, ","));
} else {
    asyncCommands.del(key);  // 列表为空时删除 Key
}
```

**最终 Redis 写入格式：**

```
Key:   wedo_guild_online_audit_list_{rUid}
Type:  String
Value: {guildId1},{guildId2},{guildId3},...

示例：
GET wedo_guild_online_audit_list_576465154108032693
→ "guild_456,guild_789,guild_101"
```

---

### 状态管理汇总

| Job | 状态描述符名称 | 类型 | TTL | 用途 |
|-----|--------------|------|-----|------|
| FriendsOperator | `wedo_is_register_new` | `ValueState<Boolean>` | 无 | 防 Timer 重复注册 |
| FriendsOperator | `wedo_friend` | `ListState<String>` | 无 | 聚合同次登录的多条好友数据，Timer 触发后主动 `clear()` |
| LogOperator | `guildId_number_online` | `MapState<String, Integer>` | 无 | 维护每支战队的在线成员集合（rUid → 1） |
| GuildAudit | *(无 Flink 托管状态)* | — | — | 使用 Time Window 聚合，无 Keyed State |

> **FriendsOperator 状态无显式 TTL 的说明**：ListState 在 2s Timer 触发时被主动 `clear()`，正常情况下实际占用时间极短。与车队推荐的 10s TTL 方案相比，战队版本依赖代码层面的 `clear()` 保障，若 Job 异常重启导致 Timer 丢失，ListState 可能出现积压，需关注 TaskManager 内存水位。

---

### 双 Kafka 容灾模式

三个 Job 均采用相同的容灾模式：

```java
FlinkKafkaConsumer<String> consumer1 = KafkaUtil.initKafka(kafkaMaster, topic, groupID);
FlinkKafkaConsumer<String> consumer2 = KafkaUtil.initKafka(kafkaSlave,  topic, groupID);

env.addSource(consumer1)
   .union(env.addSource(consumer2))
   ...
```

| Kafka 集群 | 地址 | 角色 |
|-----------|------|------|
| Master | `njpub1/2/3-199.kpmq.tencent.net:9092` | 主集群 |
| Slave  | `njpub1/2/3-204.kpmq.tencent.net:9092` | 备集群 |

> **注意**：GuildAudit 使用的是 CLB 地址（`199-kafka-clb.kpmq.woa.com:9092` / `204-kafka-clb.kpmq.woa.com:9092`），其余两个 Job 使用直连地址，接入方式有所不同。

两路 Consumer 使用**相同的 `groupID`**，重复消费时通过以下机制保证幂等性：

| Job | 幂等保障机制 |
|-----|------------|
| FriendsOperator | `hset` 覆盖写，同样的好友集合多次写入结果不变 |
| LogOperator | ES 按 `guild_id` upsert；Redis `set` 覆盖写 |
| GuildAudit | 窗口内 HashSet 去重；`set` 覆盖写 |

---

## 存储设计

### Redis 键空间设计

| 键模板 | 类型 | 说明 | TTL |
|--------|------|------|-----|
| `wedo_guild_online_friendlist_{rUid}` | Hash | 好友关系链，Field=friendType（"1"/"2"），Value=逗号分隔的好友 ruid 列表 | 无（持久）|
| `wedo_guild_online_join_{rUid}` | String | 玩家当前所属战队 ID，"0" 表示未加入战队 | 无（持久）|
| `wedo_guild_online_audit_list_{rUid}` | String | 用户屏蔽的战队 ID 列表，逗号分隔 | 无（持久）|

### Elasticsearch 索引设计

```
索引: wedo_guild_online
主键: guild_id

关键字段:
  - guild_id (keyword)         ← 唯一主键
  - env_id (integer)           ← 必过滤
  - world_id (integer)         ← 区服过滤
  - plat (integer)             ← 平台过滤
  - guild_level (integer)
  - guild_name (keyword)
  - member_size (integer)      ← 当前成员数
  - member_list (keyword)      ← 成员 ruid 逗号串
  - is_full (integer)          ← 是否满员（0/1）
  - online_num (integer)       ← 当前在线成员数（实时更新）
  - join_limit_type (integer)  ← 加入限制类型
  - join_limit_value (integer) ← 加入限制值（如最低战力）
  - guild_audit_type (integer) ← 审批模式（0=自由加入 1=需审批）
  - guild_combat_power (long)  ← 战队战力
  - guild_br_rank_score (long) ← 吃鸡排名积分
  - create_time (long)
  - partition (integer)        ← 随机分桶 [0,50)，创建时写入
```

---

## 部署与运维

### Flink 任务列表

| Job 名称 | 入口类 | 消费 Topic |
|---------|--------|----------|
| `wedo_friend_online` | FriendsOperator | `billow_ex_wedo_playerfriendslistflow` |
| `wedo_guild_online` | LogOperator | `billow_ex_wedo_playerlogin` / `billow_ex_wedo_playerlogout` / `billow_ex_wedo_modifyguildflow` / `billow_ex_wedo_createguildflow` |
| `wedo_guild_test_audit_list` | GuildAudit | `billow_ex_wedo_blockguildrecommendflow` |

### Kafka Topic 速查

| Topic | 用途 | 消费 Flink 任务 |
|-------|------|----------------|
| `billow_ex_wedo_playerfriendslistflow` | 玩家首次登录上报好友列表 | FriendsOperator |
| `billow_ex_wedo_playerlogin` | 玩家登录事件 | LogOperator |
| `billow_ex_wedo_playerlogout` | 玩家登出事件 | LogOperator |
| `billow_ex_wedo_modifyguildflow` | 战队信息变更（加入/退出/改名/销毁等） | LogOperator |
| `billow_ex_wedo_createguildflow` | 战队创建事件 | LogOperator |
| `billow_ex_wedo_blockguildrecommendflow` | 用户屏蔽/恢复战队推荐 | GuildAudit |

### Flink 算子速查

| 算子 | 使用 Job | 用途 |
|------|---------|------|
| `addSource().union()` | 全部 | 双 Kafka 容灾 |
| `flatMap(RichFlatMapFunction)` | FriendsOperator | TLog 管道符字段解析，提取 ruid 和好友列表 |
| `keyBy(KeySelector)` | 全部 | 按 rUid/guildId 分区 |
| `process(KeyedProcessFunction)` + `ListState` + `ValueState` + `Timer` | FriendsOperator | 2s 微批聚合同次登录好友流水，HashSet 去重后写 Redis Hash |
| `process(KeyedProcessFunction)` | LogOperator | 战队变更事件类型分发，写 ES + SideOutput |
| `process(ProcessFunction)` | LogOperator | 登录/登出/创建事件处理，写 Redis |
| `AsyncDataStream.unorderedWait(AsyncRedisGet)` | LogOperator | 异步查询用户战队归属，供登出事件使用 |
| `assignTimestampsAndWatermarks` + `timeWindow(5s, EventTime)` + `ProcessWindowFunction` + `MapState` | LogOperator | 4路事件流合并，5s EventTime 窗口统计战队在线成员 |
| `timeWindow(5s, ProcessingTime)` + `ProcessWindowFunction` | GuildAudit | 5s 窗口聚合屏蔽/恢复操作，含清空时序处理 |
| `AsyncDataStream.unorderedWait(AsyncRedisGetFun)` | GuildAudit | 异步读取现有屏蔽列表，与窗口结果合并后写回 Redis |
| `addSink(EsSink)` | LogOperator | ES upsert 写入战队信息和在线成员数 |
