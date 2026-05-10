# 互娱场景通用抽奖系统设计
> 基于 LotteryTemplate 插件化架构，支持有放回/无放回抽奖、分层奖池、里程碑保底、多维度道具过滤、连抽、概率缩放、异步发货、A/B 实验等能力，适用于腾讯游戏旗下多款游戏的各类抽奖活动场景。

---

## 10个关键技术决策

| 决策 | 选择 | 核心理由 |
|------|------|---------|
| **加权随机采样 + Alias Method** | 预构建 Alias 表，每次抽奖 O(1) | 概率组道具数量可能达数百个，线性扫描 O(n) 无法满足百万 QPS；Alias Method 预处理 O(n) 后每次采样 O(1)，内存换时间 |
| **Redis Lua 原子扣减限量库存** | DECRBY + 判断 ≥ 0，失败立即 INCRBY 回补 | 限量道具全局并发竞争，分布式锁性能差；Lua 脚本单线程原子执行天然防超发，单分片承载 10万+ QPS |
| **RocketMQ 事务消息异步发货** | 半消息 + 本地事务 + 回查 | 发货涉及多个游戏服务器 RPC，同步调用会拖垮抽奖 RT；事务消息保证"本地扣减成功 ↔ 发货消息必达"的原子性 |
| **插件化编排引擎** | 核心流程固定，15+ 扩展点通过插件覆盖 | 多款游戏玩法差异大（LOL/王者/原神式），硬编码无法维护；插件后注册覆盖先注册，新业务只需实现差异插件 |
| **分层概率组 + Fixed/Unfixed 双轨** | Fixed 道具概率恒定，Unfixed 道具被过滤后按比例缩放 | 保证"保底道具/兜底道具"概率不受过滤影响，同时 Unfixed 道具概率自动重分配，无需手动调权重 |
| **三级幂等保障** | 响应缓存 → 双缓存回退/回放 → Redis 分布式锁 | 客户端网络重试、动画重播、并发请求三种场景各有应对，任何异常都不会导致多扣费或多发奖 |
| **里程碑幸运值 + 优先级分组** | 每次抽奖累加幸运值，达到阈值切换高稀有度概率组 | 通用保底机制：可配置多条规则多优先级，覆盖硬保底/软保底/大小保底/跨池继承等所有业界玩法 |
| **多维过滤链 + 兜底概率组** | 6层过滤器串行执行，全部过滤空时自动切换 complement 兜底 | 过滤后池可能为空（无放回抽完、限量耗尽），兜底机制保证每次抽奖必有结果，避免空产出 |
| **Layer 多层隔离** | 每层独立 Redis 状态/奖池/概率表/缓存 | 多轮活动、赛季制、分段抽奖天然隔离，切换 Layer 等于重开新局，无需清理旧状态 |
| **MySQL 对账 + 概率监控** | 定时 Redis vs MySQL 库存对账 + 实时产出率监控告警 | Redis 主从切换可能丢数据；概率偏离可能因 Bug/作弊导致——双重兜底保证数据正确、概率公平 |

---

## 1. 需求澄清与非功能性约束

### 功能性需求

**核心功能：**
- **基础抽奖**：支持单抽、N连抽（任意次数）、免费抽、折扣抽、阶梯消耗抽
- **概率模型**：分层概率组、概率缩放（UP活动）、动态概率（软保底递增）、伪随机分布
- **奖池管理**：静态/动态/分层/自选/轮转奖池，按 Layer 多层隔离，支持 A/B 实验驱动差异化奖池
- **里程碑保底**：硬保底、软保底、大小保底（50/50 then 100%）、跨池继承、里程碑奖励、天花板保底
- **有放回/无放回**：有放回（可重复获得）、无放回（抽到即移除）、部分无放回、数量限制
- **多维度过滤**：有效期→已有→背包→已抽→限量→顺序，过滤链可配置
- **异步发货**：RocketMQ 事务消息驱动，重试补偿，幂等入账
- **抽奖记录**：个人历史查询、保底进度查询、全服产出统计
- **到账提醒**：发货成功推送、大奖全服广播、保底接近提醒

**边界约束：**
- 单次连抽上限：200次
- 单用户单池每日抽奖上限：可配置（默认1000次）
- 限量道具全局库存：Redis 原子扣减，绝不超发
- 概率公示精度：小数点后4位
- 发货延迟：P99 < 5s，超过30分钟自动补偿

### 非功能性约束

| 维度 | 指标 |
|------|------|
| 可用性 | 抽奖核心链路 99.99%，发货链路 99.9% |
| 性能 | 抽奖接口 P99 < 50ms，初始化接口 P99 < 200ms |
| 一致性 | 限量道具绝不超发，货币扣减与发货最终一致 |
| 峰值 | 新池开启瞬间 100万 QPS（含重试），有效抽奖 50万 QPS |
| 合规 | 概率公示、未成年人保护、审计日志保留1年 |

### 明确禁行需求
- **禁止超发**：限量道具任何情况下实际发放 ≤ 配置总量
- **禁止同步发货阻塞抽奖**：发货走异步 MQ，抽奖接口立即返回结果
- **禁止实时查 DB 判断库存**：高并发下 DB 无法承载，库存判断只走 Redis
- **禁止概率计算浮点误差**：权重使用整数，概率通过整数权重比计算

---

## 2. 系统容量评估

### 核心指标定义

| 参数 | 数值 | 依据 |
|------|------|------|
| 接入游戏数 | **20+** | 腾讯游戏旗下 LOL/王者/CF/DNF 等 |
| 日均抽奖请求 | **5亿次** | 20款游戏 × 平均2500万次/天 |
| 平均抽奖 QPS | **5800 QPS** | 5亿 / 86400 |
| 峰值系数 | **× 100** | 新池开启瞬间/限时活动开启 |
| **峰值抽奖 QPS** | **50万 QPS** | 晚高峰 + 新池开启叠加 |
| 网关入口 QPS | **200万 QPS** | 含重试、轮询、初始化请求（4倍放大） |
| 有效 Redis 操作 | **150万 ops/s** | 每次抽奖约3次 Redis 操作（读状态+扣库存+写状态） |
| DB 写入（MQ 削峰后） | **5万 TPS** | 发货记录 + 抽奖流水异步落库 |

### 容量计算

**存储规划：**

| 数据 | 计算过程 | 估算结果 | 说明 |
|------|---------|---------|------|
| 用户状态 (Redis) | 2000万活跃用户 × 2KB/人 | **40 GB** | 保底计数 + 奖池状态 + 已抽道具 |
| 奖池配置 (Redis) | 500个活跃奖池 × 50KB/池 | **25 MB** | Alias 表 + 概率组缓存 |
| 限量库存 (Redis) | 1万种限量道具 × 64B/种 | **640 KB** | 库存计数 + 分桶 |
| 抽奖流水 (MySQL) | 5亿条/天 × 200B/条 | **100 GB/天** | 按月分表，保留1年 |
| 发货记录 (MySQL) | 5亿条/天 × 150B/条 | **75 GB/天** | 含重试记录 |

**Redis 集群：**
- 实际到达 Redis 的 QPS = 150万 ops/s
- Redis 单分片安全 QPS：**10万**
- 所需分片数：150万 ÷ 10万 = 15，取 **32分片**（2倍冗余）
- 内存规划：40GB × 1.5（碎片+复制缓冲）= **60GB**，每分片约 2GB

**DB 分库分表：**
- MySQL 单主库安全写入：**3000 TPS**
- MQ 削峰后 DB 写入：5万 TPS
- 所需分库数：5万 ÷ 3000 ≈ 17，取 **16库**
- 每库分 **256表**（按 user_id 哈希）

**服务节点（Go，8核16G）：**

| 服务 | 单机安全 QPS | 有效 QPS | 节点数 |
|------|-------------|---------|--------|
| 抽奖核心服务 | 2000 | 50万 | **360台** |
| 初始化服务 | 1000 | 10万 | **150台** |
| 发货消费服务 | 3000 | 5万 | **25台** |
| 查询服务 | 5000 | 20万 | **60台** |

**RocketMQ 集群：**
- 单 Broker 主节点安全吞吐：5万 TPS
- 所需主节点：50万（抽奖成功即发 MQ）÷ 5万 = 10，取 **16主节点**（冗余）
- 从节点：16从，每主配1从
- 发货 Topic 同步刷盘（道具不能丢）

---

## 3. 核心领域模型

### 实体 + 事件

#### 实体（Entity，写模型）

| 模型 | 职责 | 核心属性 | 存储位置 |
|------|------|---------|---------|
| **LotteryPool** 奖池 | 奖池全生命周期：配置→生效→下线 | 奖池ID、游戏ID、版本号、类型、时间范围、层级配置、状态 | MySQL 配置表 + Redis 缓存 |
| **UserLotteryState** 用户抽奖状态 | 用户在某奖池的完整进度 | 用户ID、奖池ID、Layer、已抽道具、幸运值、保底计数、当前概率组ID、奖池道具集合 | Redis（权威源）+ MySQL 异步持久化 |
| **DrawRecord** 抽奖记录 | 每次抽奖的完整流水 | 抽奖ID(雪花)、批次ID、用户ID、奖池ID、道具ID、稀有度、是否保底、概率快照 | MySQL 按月分表 |
| **Inventory** 限量库存 | 限量道具的全局库存管控 | 奖池ID、道具ID、总库存、已发放、Redis 实时值 | Redis 原子计数（权威）+ MySQL 对账兜底 |
| **DeliveryRecord** 发货记录 | 道具发放追踪 | 发货ID、抽奖ID、用户ID、道具ID、渠道、状态、重试次数 | MySQL |

#### 事件（Event，事件流）

| 事件 | 触发时机 | 下游消费 |
|------|---------|---------|
| **DrawCompleted** 抽奖完成 | 本地事务提交成功 | ① 异步发货 ② 写抽奖流水 ③ 统计上报 |
| **DeliverySuccess** 发货成功 | 游戏服务器确认到账 | ① 更新发货状态 ② 推送到账通知 ③ 大奖广播 |
| **DeliveryFailed** 发货失败 | 重试耗尽 | ① 人工工单 ② 补偿触发 |
| **MilestoneHit** 里程碑命中 | 幸运值达到阈值 | ① 日志记录 ② 保底统计 |
| **StockDepleted** 库存耗尽 | Redis 库存 ≤ 0 | ① 运营告警 ② 自动切换兜底道具 |

#### 模型关系图

```
  [写路径（抽奖核心）]              [事件流]                    [读路径]
  ┌──────────────────┐                                    ┌──────────────────┐
  │ UserLotteryState │──DrawCompleted─────┐               │   DrawRecord     │ ← 抽奖流水
  │ (Redis 状态机)    │                    │               │  (MySQL 按月分表) │
  │ ├─ 幸运值/保底   │                    │               └──────────────────┘
  │ ├─ 已抽道具      │                    │               ┌──────────────────┐
  │ └─ 奖池道具集合  │                    ├─RocketMQ──→  │  DeliveryRecord  │ ← 发货记录
  └──────────────────┘                    │               └──────────────────┘
         │                                │               ┌──────────────────┐
         │ 库存扣减                        │               │  到账通知/广播   │
         ▼                                │               └──────────────────┘
  ┌──────────────────┐                    │
  │   Inventory      │──StockDepleted─────┘
  │ (Redis 原子计数)  │
  │ MySQL 对账兜底   │
  └──────────────────┘

  ┌──────────────────┐
  │  LotteryPool     │ ← 运营后台配置，热加载到 Redis
  │ (MySQL + Redis)  │
  └──────────────────┘
```

---

## 4. 整体架构图

```
┌─────────────────────────────────────────────────────────────────────┐
│                        接入层 (Access Layer)                          │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────────┐    │
│  │ 游戏客户端│  │ H5/小程序 │  │  运营后台 │  │ 第三方平台(开平) │    │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────────┬─────────┘    │
└───────┼──────────────┼──────────────┼─────────────────┼─────────────┘
        │              │              │                  │
┌───────▼──────────────▼──────────────▼─────────────────▼─────────────┐
│                   网关层 (tRPC-Gateway)                               │
│  ┌────────┐ ┌──────────┐ ┌────────────┐ ┌──────────┐ ┌──────────┐  │
│  │限流熔断 │ │身份认证   │ │参数校验     │ │灰度路由   │ │协议转换   │  │
│  └────────┘ └──────────┘ └────────────┘ └──────────┘ └──────────┘  │
└─────────────────────────────┬───────────────────────────────────────┘
                              │
┌─────────────────────────────▼───────────────────────────────────────┐
│                   抽奖核心服务 (LotteryTemplate)                       │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │              插件化编排引擎 (Plugin Orchestrator)                │   │
│  │  ┌─────────┐┌──────────┐┌─────────┐┌──────────┐┌──────────┐ │   │
│  │  │ProbIdGen││FilterChain││DrawPlugin││LuckyValue││BuildRsp  │ │   │
│  │  └─────────┘└──────────┘└─────────┘└──────────┘└──────────┘ │   │
│  └──────────────────────────────────────────────────────────────┘   │
│  ┌─────────┐┌─────────┐┌─────────┐┌─────────┐┌─────────────────┐  │
│  │资格校验  ││概率模型  ││库存管理  ││状态管理  ││ 响应缓存/幂等    │  │
│  └─────────┘└─────────┘└─────────┘└─────────┘└─────────────────┘  │
└─────────────────────────────┬───────────────────────────────────────┘
                              │
┌─────────────────────────────▼───────────────────────────────────────┐
│                        中间件层                                       │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐              │
│  │ Redis Cluster │  │  RocketMQ    │  │    ETCD      │              │
│  │ 32分片1主2从  │  │ 16主16从     │  │  配置中心    │              │
│  └──────────────┘  └──────────────┘  └──────────────┘              │
└─────────────────────────────┬───────────────────────────────────────┘
                              │
┌─────────────────────────────▼───────────────────────────────────────┐
│                        存储层                                         │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐              │
│  │MySQL 16库 │ │ ES 日志  │ │Prometheus│ │  HDFS    │              │
│  │按月分表   │ │ 审计查询 │ │+ Grafana │ │ 冷数据   │              │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘              │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 5. 核心流程

### 5.1 抽奖主流程 (SceneID 18101)

```
请求进入
    │
    ▼
1. 响应缓存检查（stateKey = drawType_drawRound）
   ├─ 命中且 CheckRspCachePlugin 校验通过 → 直接返回缓存响应（无需加锁）
   └─ 未命中 ↓
    │
    ▼
2. Redis 分布式锁（SETNX，防并发抽奖）
    │
    ▼
3. 获取 A/B 实验分组
    │
    ▼
4. 获取 UserRealTimeState
   ├─ Redis 中存在 → 读取 + 处理回放/双缓存重抽
   └─ Redis 中不存在 → 调用 18201 自动初始化
    │
    ▼
5. 请求参数校验
   ├─ 必传字段校验 (partition/area/openid/roleid)
   ├─ drawType 校验（连抽次数 1~200）
   ├─ drawRound 校验（轮次连续性）
   ├─ haveItems 一致性校验
   └─ 消耗检查（货币/道具/次数）
    │
    ▼
6. 抽取道具 (DrawItemAndRecord) ─── 核心抽奖循环
   ├─ ExtRuleDrawPlugin 特殊规则（首抽必得等）
   ├─ 循环 drawType 次：
   │   ├─ ProbIdGenPlugin：确定概率组（里程碑匹配）
   │   ├─ BeforeRuleDrawPlugin：前置钩子
   │   ├─ RoundDrawItem 单轮抽取：
   │   │   ├─ GenProb：生成概率组（含 Alias 表）
   │   │   ├─ FilterProb：6层过滤链
   │   │   ├─ FixProbPlugin：概率修正（UP/衰减）
   │   │   ├─ BeforeDrawProbPlugin：采样前修正
   │   │   ├─ WeightSampling：加权随机采样（O(1) Alias）
   │   │   ├─ 限量检查：Redis Lua 原子扣减 → 触发则重抽
   │   │   ├─ ReplaceItemPlugin：道具替换
   │   │   └─ LuckyValuePlugin.Update：更新幸运值
   │   └─ AfterRuleDrawPlugin：后置钩子（可中断连抽）
   └─ AfterAllRoundRuleDrawPlugin：全轮次完成钩子
    │
    ▼
7. 本地事务提交（MySQL）
   ├─ 扣减用户货币
   ├─ 更新 UserRealTimeState → Redis
   └─ 写 RocketMQ 半消息 → Commit（发货指令）
    │
    ▼
8. 构建响应 (BuildRsp)
   ├─ BuildItemPlugin：道具列表
   ├─ BuildAlgoPlugin：算法信息（A/B 分组）
   ├─ 大奖产出监控检查
   └─ 持久化状态到 Redis + 写响应缓存
    │
    ▼
9. AlarmCheckPlugin：响应校验
    │
    ▼
10. PostRspPlugin → 返回客户端
```

### 5.2 单轮抽奖详细流程 (RoundDrawItem)

```
┌─────────────────────────────────────────────────────────────────────────┐
│                         RoundDrawItem 单轮抽取                            │
│                                                                         │
│  ① 确定概率组 ID                                                         │
│     ├─ 读取 UserRealTimeState.ProbID（默认 "base"）                       │
│     ├─ LuckyValuePlugin.Match() 匹配里程碑                               │
│     │   ├─ 幸运值 >= 阈值 → 切换到高稀有度概率组（如 "milestone_ssr"）      │
│     │   └─ 未命中 → 保持当前概率组                                         │
│     └─ 校验概率组在配置中存在                                              │
│                                                                         │
│  ② 生成概率组 (GenProb)                                                  │
│     ├─ 从 ProbConfMap 加载概率组配置                                       │
│     ├─ 区分 Fixed Items（固定概率）和 Unfixed Items（浮动概率）             │
│     └─ 支持概率组嵌套（子概率组递归，最大10层）                              │
│                                                                         │
│  ③ 道具过滤链 (FilterProb)                                               │
│     ├─ FilterItemByValidityPeriod ── 有效期过滤                           │
│     ├─ FilterItemByHaveItems     ── 已有道具过滤                          │
│     ├─ FilterItemByKnapsacks     ── 背包过滤（外部数据源）                 │
│     ├─ FilterItemByDrawn         ── 已抽中过滤（本次连抽去重）              │
│     ├─ FilterItemByLimiter       ── 限量过滤（全局/日期/用户维度）          │
│     └─ FilterItemByOrder         ── 产出顺序过滤                          │
│                                                                         │
│  ④ 概率修正 (FixProbPlugin)                                              │
│     ├─ UP 活动：目标道具权重翻倍，同层其他按比例缩小                        │
│     ├─ 连出衰减：连续命中高稀有度后概率衰减                                 │
│     └─ 自适应：根据全服产出率微调                                          │
│                                                                         │
│  ⑤ 概率重分配                                                            │
│     ├─ Fixed 道具：实际概率 = 配置概率（恒定不变）                          │
│     └─ Unfixed 道具：实际概率 = 配置概率 × (1-FixedSum) / UnfixedSum      │
│                                                                         │
│  ⑥ 加权随机采样 (WeightSampling / Alias Method)                          │
│     ├─ Alias 表预构建（概率组变更时重建）                                   │
│     ├─ 生成随机数 → O(1) 查表 → 得到道具 ID                               │
│     └─ 若命中子概率组 → 递归进入子概率组采样                               │
│                                                                         │
│  ⑦ 限量检查 (Redis Lua 原子操作)                                         │
│     ├─ 全局限量：DECRBY + 判断 ≥ 0                                       │
│     │   ├─ 成功 → 继续                                                   │
│     │   └─ 失败 → INCRBY 回补 → 标记该道具 → 重抽（最多3次）              │
│     ├─ 日期限量：Redis INCR + 当天 TTL                                    │
│     └─ 用户限量：本地 HaveItems 判断（无需 Redis）                         │
│                                                                         │
│  ⑧ 道具替换 (ReplaceItemPlugin)                                         │
│     └─ 默认不替换；自定义场景：道具升级/测试强制指定                        │
│                                                                         │
│  ⑨ 更新幸运值 (LuckyValuePlugin.UpdateLuckyValue)                        │
│     ├─ 抽中大奖 → 清零对应优先级组的幸运值                                 │
│     └─ 未中大奖 → 幸运值 +1（每轮累加）                                    │
│                                                                         │
│  ⑩ 池空处理                                                              │
│     ├─ Unfixed 全部被过滤 → 切换到兜底概率组 (complement)                  │
│     ├─ 支持自定义兜底 ({probId}_complement)                               │
│     └─ 兜底也为空 → 返回配置的 fallback_item（如金币）                     │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

### 5.3 异步发货流程（参考红包入账模式）

```
抽奖本地事务成功
        │
        ▼
① 发送 RocketMQ 半消息（LOTTERY_DELIVERY Topic）
   ├─ 消息体：draw_id, batch_id, user_id, items[], pool_id, timestamp
   └─ Keys: draw_id（用于消息轨迹查询）
        │
        ▼
② 执行本地事务
   ├─ MySQL 事务：扣货币 + 写抽奖记录 + 更新保底计数
   ├─ Redis：更新 UserRealTimeState
   └─ Redis：限量道具原子扣减（已在抽奖循环中完成）
        │
        ▼
③ 本地事务结果
   ├─ 成功 → Commit 半消息（消费者可见）
   └─ 失败 → Rollback 半消息 + Redis INCRBY 回补库存
        │
        ▼
④ 事务回查（Producer 宕机/超时时 RocketMQ 主动回查）
   ├─ 查询 draw_record 表：draw_id 存在 → Commit
   ├─ draw_id 不存在 → Rollback
   └─ 不确定 → UNKNOW（继续回查，最多15次）
        │
        ▼
⑤ 发货 Worker 消费（GID_LOTTERY_DELIVERY 集群消费）
   ├─ 幂等检查：delivery_record 表 uk_draw_item (draw_id, item_index)
   │   └─ 已存在 → 直接 ACK
   ├─ 路由发货渠道：
   │   ├─ character/skin → game_server_rpc（游戏服务器 gRPC）
   │   ├─ currency/coin → wallet_service（钱包服务）
   │   ├─ coupon → coupon_center（优惠券中心）
   │   └─ physical → logistics_service（物流，需地址）
   ├─ 调用渠道方发货 API
   │   ├─ 成功 → 更新 delivery_record.status=SUCCESS
   │   └─ 失败 → 重试（RocketMQ 自动重试，最多16次，指数退避）
   └─ 发货成功后：
       ├─ 推送到账通知（游戏内弹窗/Push）
       ├─ 大奖命中 → 全服广播 MQ（跑马灯）
       └─ 写 account_flow 流水
        │
        ▼
⑥ 死信处理（重试16次仍失败）
   ├─ 进入 %DLQ%GID_LOTTERY_DELIVERY
   ├─ 触发 P0 告警
   ├─ 创建人工工单
   └─ 30分钟后自动补偿（等价代币）

【日级对账（凌晨3点）】
全量校验：COUNT(draw_record) = COUNT(delivery_record WHERE status=SUCCESS) + COUNT(补偿记录)
差异 > 0 → 自动补发；差异 < 0（超发）→ P0 告警人工处理
```

### 5.4 库存管理流程（Redis原子扣减 + MySQL兜底）

```
┌──────────────────────────────────────────────────────────────────────┐
│                  限量道具库存管理                                       │
│                                                                      │
│  Redis Lua 原子扣减脚本:                                               │
│  ┌─────────────────────────────────────────────────────────────┐     │
│  │ local stock = redis.call('GET', KEYS[1])                    │     │
│  │ if stock == false then return -1 end                        │     │
│  │ local remain = tonumber(stock) - tonumber(ARGV[1])          │     │
│  │ if remain < 0 then return -1 end      -- 库存不足           │     │
│  │ redis.call('DECRBY', KEYS[1], ARGV[1])                     │     │
│  │ return remain                          -- 返回剩余库存      │     │
│  └─────────────────────────────────────────────────────────────┘     │
│                                                                      │
│  Key: stock:{pool_id}:{item_id}                                      │
│                                                                      │
│  ═══════════════════════════════════════════════════════════════════  │
│                                                                      │
│  操作流程:                                                            │
│  ① 抽奖命中限量道具 → 执行 Lua 扣减                                    │
│     ├─ 返回 ≥ 0 → 扣减成功，继续业务流程                               │
│     └─ 返回 -1  → 库存不足，标记该道具，重新抽取                        │
│  ② 本地 DB 事务失败 → Redis INCRBY 回补                               │
│  ③ 发货 Worker 消费前二次校验 MySQL（防 Redis 主从切换超发）             │
│                                                                      │
│  ═══════════════════════════════════════════════════════════════════  │
│                                                                      │
│  MySQL 兜底对账:                                                      │
│  ├─ 每5分钟: redis_stock vs (total - COUNT(delivery WHERE success))  │
│  ├─ 差异修正: 以 MySQL 发货记录为准纠正 Redis                           │
│  ├─ 每日凌晨: 全量对账                                                │
│  └─ 差异 > 5: 触发告警                                               │
│                                                                      │
└──────────────────────────────────────────────────────────────────────┘
```

---

## 6. 插件系统设计

### 6.1 插件架构

```
服务启动
  │
  ▼
plugin.InitPlugin(extPlugins)
  │
  ├─ 1. 注册所有默认插件 (Default*Plugin)
  │     提供空实现或基础逻辑
  │
  ├─ 2. customplugin.InitCustomPlugin()
  │     根据 busi.yaml 的 custom_plugins 配置
  │     从 CustomPlugins 预注册池中选择性覆盖默认实现
  │
  └─ 3. 注册外部传入的 extPlugins
        通过 NewLotteryServer(plugins...) 传入（优先级最高）
```

**关键设计**：后注册的同名插件覆盖先注册的——默认插件先注册，自定义插件后注册覆盖。

### 6.2 插件清单

| 插件类别 | 插件名 | 接口方法 | 职责 |
|---------|--------|---------|------|
| ID生成 | ProbIdGenPlugin | GenProbId / GenInitProbId | 确定概率组（里程碑匹配） |
| ID生成 | PoolIdGenPlugin | GenPoolId | 确定奖池（可按等级/A/B分） |
| 过滤 | HaveItemsFilterPlugin | GetFilterItemsByHaveItems | 已有道具过滤 |
| 过滤 | KnapsacksFilterPlugin | GetAllFilterItemsByKnapsacks | 背包数据过滤 |
| 过滤 | DrawnFilterPlugin | GetFilterItemsByDrawn | 连抽内已抽过滤 |
| 过滤 | LimitFilterPlugin | GetFilterItems / UpdateFilterItem | 四层限量过滤 |
| 抽奖流程 | ExtRuleDrawPlugin | ExtRuleDraw | 特殊规则（首抽必得） |
| 抽奖流程 | BeforeRuleDrawPlugin | BeforeRuleDraw | 单轮前置钩子 |
| 抽奖流程 | AfterRuleDrawPlugin | AfterRuleDraw | 单轮后置（可中断） |
| 抽奖流程 | AfterAllRoundRuleDrawPlugin | AfterAllRoundRuleDraw | 全轮次完成钩子 |
| 抽奖流程 | ReplaceItemPlugin | ReplaceItem | 道具替换 |
| 抽奖流程 | BeforeDrawProbPlugin | BeforeDrawProb | 采样前概率修正 |
| 里程碑 | LuckyValuePlugin | Match / UpdateLuckyValue / Disable | 保底幸运值管理 |
| 概率 | FixProbPlugin | FixProb | 概率动态修正 |
| 奖池 | InitPoolPlugin | InitPool | 奖池初始化（A/B） |
| 奖池 | InitSysPoolPlugin | InitSysPool | 系统默认奖池生成 |
| 响应 | BuildItemPlugin | BuildItem | 单道具响应构建 |
| 响应 | FilterItemsPlugin | FilterItems | 响应道具列表过滤 |
| 响应 | BuildAlgoPlugin | BuildAlgo | 算法/AB 信息 |
| 响应 | BuildRspPlugin | BuildExt1 / Build18201Ext1 | 扩展字段 |
| 响应 | PostRsp18101Plugin | PostRsp | 响应后处理+上报 |
| 校验 | CheckRspCachePlugin | CheckRspCache | 缓存有效性校验 |
| 校验 | AlarmCheckPlugin | CheckReq / CheckRsp | 告警校验（必须覆盖）|

### 6.3 插件扩展方式

```go
// 方式一：配置激活（推荐）
// 1. customplugin/ 目录创建文件
// 2. init() 中 Add(name, plugin) 预注册
// 3. busi.yaml 的 custom_plugins 添加插件名
custom_plugins: "del_empty_sub_prob_plugin,filter_complement_items_plugin"

// 方式二：外部传入（最高优先级）
server.NewLotteryServer(&MyPlugin1{}, &MyPlugin2{})

// 方式三：代码直接覆盖
plugindef.Add(&MyPlugin{})
```

---

## 7. 概率模型

### 7.1 概率组模型t

```yaml
prob_conf:
  prob_id: "base"
  items:
    # Fixed 道具：概率恒定，不受过滤影响
    - item_id: "ssr_char_001"
      weight: 6           # 0.6%
      fixed: true
    - item_id: "ssr_char_002"
      weight: 6
      fixed: true
      
    # Unfixed 道具：被过滤后剩余道具按比例缩放
    - item_id: "sr_char_001"
      weight: 51          # 5.1%
      fixed: false
    - item_id: "r_item_001"
      weight: 200
      fixed: false
    - item_id: "n_item_001"
      weight: 737
      fixed: false
      
  # 概率计算
  # Fixed 实际概率 = weight / total_weight（恒定）
  # Unfixed 实际概率 = weight × (1 - sum_fixed_prob) / sum_unfixed_weight
```

### 7.2 概率缩放（UP活动）

```yaml
probability_scaling:
  # 分层内缩放：目标道具权重翻倍，同层其他道具按比例缩小
  rules:
    - target: "ssr_char_up"
      mode: "within_layer"       # 不影响其他层概率
      scale_factor: 2.0          # UP 角色权重翻倍
      
    - target: "sr_weapon_up"
      mode: "within_layer"
      scale_factor: 1.5
```

### 7.3 概率组嵌套（子概率组）

```
道具 ID 可指向另一个概率组 → 命中时递归进入子概率组采样

示例：
  base 概率组:
    ├─ item_a (weight=500)
    ├─ item_b (weight=300)
    └─ sub_prob_rare (weight=200) ← 指向子概率组
        └─ rare 概率组:
            ├─ rare_item_1 (weight=50)
            ├─ rare_item_2 (weight=30)
            └─ rare_item_3 (weight=20)

最大递归深度：10层
```

### 7.4 概率线性滑动 (SlideProbLinear)

在两个概率组之间按轮次线性插值：

```
第1轮概率 = prob_start
第N轮概率 = prob_start + (prob_end - prob_start) × (N-1) / (total_rounds-1)
第total_rounds轮概率 = prob_end

适用场景：保底概率随轮次逐步提升（软保底）
```

### 7.5 Alias Method 高性能采样

```go
type AliasTable struct {
    n     int
    prob  []float64  // 每个桶的阈值
    alias []int      // 别名指向
}

// O(1) 采样
func (t *AliasTable) Sample() int {
    i := rand.Intn(t.n)           // 随机选桶
    if rand.Float64() < t.prob[i] {
        return i                   // 留在本桶
    }
    return t.alias[i]             // 跳转到别名
}

// 概率组变更时重建 Alias 表，缓存在 Redis
// Key: alias_table:{pool_id}:{prob_id}:{version}
```

---

## 8. 里程碑（保底）系统

### 8.1 核心机制

```
每次抽奖后 → LuckyValuePlugin.UpdateLuckyValue 更新幸运值
每次抽奖前 → ProbIdGenPlugin.GenProbId → LuckyValuePlugin.Match 匹配里程碑
命中里程碑 → 切换到高稀有度概率组 → 保证产出大奖
抽中大奖后 → 清零对应优先级组的幸运值（重新累积）
```

### 8.2 里程碑配置

```yaml
milestone_conf:
  rules:
    # 硬保底：90抽必出 SSR（优先级1，最高）
    - priority: 1
      priority_alias: "ssr_pity"
      lucky_value: 90
      prob_id: "milestone_ssr_100"    # 100%产出SSR的概率组
      
    # 软保底：73抽起概率递增
    - priority: 2
      priority_alias: "ssr_soft"
      lucky_value: 73
      prob_id: "milestone_ssr_soft"   # 概率提升的概率组
      use_slide_prob: true            # 启用线性滑动
      slide_end_lucky_value: 89       # 滑动到89抽
      
    # 小保底：10抽必出 SR
    - priority: 3
      priority_alias: "sr_pity"
      lucky_value: 10
      prob_id: "milestone_sr_100"
      
    # 大小保底（50/50 then 100%）
    - priority: 4
      priority_alias: "up_guarantee"
      lucky_value: 1                  # 首次命中SSR即触发
      prob_id: "milestone_up_check"   # 50%UP判定概率组
      condition: "last_ssr_not_up"    # 上次SSR非UP时下次必UP

  # 清零规则：抽中某道具后清零哪些优先级组
  clean_rules:
    - item_tag: "ssr"
      clean_groups: ["ssr_pity", "ssr_soft"]
    - item_tag: "sr"
      clean_groups: ["sr_pity"]
    - item_tag: "up_ssr"
      clean_groups: ["up_guarantee"]
```

### 8.3 支持的保底玩法

| 玩法 | 实现方式 |
|------|---------|
| 硬保底（90抽必出） | lucky_value=90, prob_id=100%SSR概率组 |
| 软保底（73抽起递增） | lucky_value=73, use_slide_prob=true |
| 大小保底 | 额外状态位 last_ssr_not_up，命中时强制UP |
| 跨池继承 | 幸运值存 UserRealTimeState，按 pool_type 隔离不按 pool_id |
| 累计消费保底 | lucky_value 按消耗量累加而非次数 |
| 里程碑奖励 | lucky_value=N 时切换到"额外奖励"概率组 |
| 天花板 | 限定道具在 max_draws 内必出（设为 Fixed + 递增） |
| 自选保底 | lucky_value=180 时触发自选UI，由 ExtRuleDrawPlugin 实现 |

---

## 9. 多维度道具过滤链

### 9.1 过滤执行顺序

```
原始概率组（全部道具）
    │
    ├─ ① FilterItemByValidityPeriod ── 有效期过滤
    │   └─ 当前时间不在道具 [start_time, end_time] 范围内 → 移除
    │
    ├─ ② FilterItemByHaveItems ── 已有道具过滤
    │   └─ 用户已拥有（HaveItems + 本次DrawnItems）且配置为唯一 → 移除
    │   └─ 支持 source→target 映射（拥有A则过滤B）
    │
    ├─ ③ FilterItemByKnapsacks ── 背包过滤
    │   └─ 查询外部背包数据源（IDIP/河图/DataMore）
    │   └─ 结果缓存在 CurrentLotteryState（同请求只查一次）
    │
    ├─ ④ FilterItemByDrawn ── 已抽中过滤（无放回）
    │   └─ 本次连抽已抽到的道具 → 移除（或达到阈值移除）
    │   └─ 典型：10连抽防同一SSR出现多次
    │
    ├─ ⑤ FilterItemByLimiter ── 限量过滤（四层）
    │   ├─ userLimit：用户累计获得数 ≥ 限制 → 移除
    │   ├─ dateLimit：全服当天产出达上限 → 移除
    │   ├─ globalLimit：全服永久限量达上限 → 移除
    │   └─ userDayLimit：用户当天获得数达上限 → 移除
    │
    └─ ⑥ FilterItemByOrder ── 产出顺序过滤
        └─ 按配置优先级保留最高优先级道具（教程引导池）
    │
    ▼
过滤后的概率组
    ├─ Unfixed 道具仍有剩余 → 正常采样
    └─ Unfixed 全部被过滤 → 切换 complement 兜底概率组
```

### 9.2 四层限量检查详细

| 检查层 | 数据来源 | Redis 操作 | 说明 |
|--------|---------|-----------|------|
| userLimit | 请求 HaveItems + 本次 DrawnItems | 无（本地计算） | 用户累计获得 ≥ 配置限制 |
| dateLimit | GlobalDateLimiter (Redis INCR) | `INCR limit:{item}:{date}` + TTL 86400 | 全服当天产出 |
| globalLimit | GlobalDateLimiter (Redis INCR) | `INCR limit:{item}:global` | 全服永久（支持时间段配置）|
| userDayLimit | 请求 HaveItemsToday + DrawnItems | 无（本地计算） | 用户当天获得 |

**两阶段限量检查：**
1. **抽奖前预过滤** (`GetFilterItemsByLimit`)：批量查询所有已达限量的道具集合，从概率组中移除
2. **抽奖后更新** (`UpdateFilterItemIdByLimit`)：采样结果确定后，Redis 原子递增。若超限返回 true → 主流程丢弃结果重抽

---

## 10. 奖池管理与 Layer 隔离

### 10.1 奖池初始化 (SceneID 18201)

```
初始化请求
    │
    ▼
1. PoolIdGenPlugin.GenPoolId() 确定奖池 ID
   ├─ 默认返回 "base"
   └─ 自定义：按服务器等级/A/B 分组选择不同奖池
    │
    ▼
2. InitPoolPlugin.InitPool() 生成奖池道具集合
   ├─ 从 LotteryItemConf 按 Layer 筛选道具
   ├─ A/B 实验驱动差异化道具集
   └─ 写入 UserRealTimeState.PoolItemsMap
    │
    ▼
3. ProbIdGenPlugin.GenInitProbId() 确定初始概率组
    │
    ▼
4. 初始化幸运值（全部归零或继承）
    │
    ▼
5. 持久化 UserRealTimeState 到 Redis
    │
    ▼
6. 构建 18201 响应（奖池道具列表）
```

### 10.2 Layer 多层隔离

```
设计目的：每个 Layer 是一个完整的抽奖大轮，抽完切换到下一 Layer，等于重开新局。

隔离维度：
├─ Redis 状态 Key：  userKey + "_" + layer
├─ 响应缓存 Key：   userKey + "_lottery_response_cache_" + layer
├─ 奖池/概率表：     base_{layer}（如 base_1, base_2）
└─ 道具 Layers 配置："all" / "1" / "1/3"（周期出现）

典型场景：
├─ 多轮活动：第一轮抽完 → layer+1 → 第二轮道具池重新装满
├─ 分段抽奖：不同 Layer 不同道具+概率
└─ 赛季制：每赛季一个 Layer
```

### 10.3 奖池切换模式

| 模式 | 说明 | 适用场景 |
|------|------|---------|
| NoSwitch | 固定奖池，永不切换 | 常规池 |
| NoSwitchAfterFirstLottery | 首次抽奖后锁定 | 防初始化后奖池配置变更 |
| SwitchAnytime | 可随时切换 | 动态奖池/A/B 实验 |

---

## 11. 并发安全与幂等保障

### 11.1 三级幂等体系

```
请求到达
    │
    ▼
第1级: 响应缓存幂等 (IsUsingResponseCache)
    │  ├─ 缓存 Key = {UserKey}_lottery_response_cache[_layer]
    │  ├─ stateKey = drawType + "_" + drawRound
    │  ├─ 命中且 stateKey 匹配 → CheckRspCachePlugin 校验
    │  │   ├─ 通过 → 直接返回缓存响应（不加锁，不碰 Redis 状态）
    │  │   └─ 不通过 → 继续走抽奖流程
    │  └─ 未命中 → 继续
    │  ※ 响应缓存检查在加锁之前，命中无需获取分布式锁
    ▼
第2级: 状态级重抽保护
    │  ├─ IsUsePlayBack=true（回放模式）：
    │  │   └─ 从 PlayBackRewardInfo 按序回放历史结果（动画重播场景）
    │  └─ IsUsePlayBack=false（双缓存模式）：
    │      └─ drawRound == LastDrawRound → 回退到 PrevState 重新随机
    ▼
第3级: Redis 分布式锁（防并发）
       └─ SET NX，确保同用户同时刻只有一个请求执行抽奖
```

### 11.2 回放模式 vs 双缓存模式

| 特性 | 回放模式 (IsUsePlayBack=true) | 双缓存模式 (IsUsePlayBack=false) |
|------|-----|------|
| 重抽方式 | 按顺序回放历史道具序列 | 回退 PrevState 重新随机 |
| 适用场景 | 客户端动画重播（确定性结果） | 客户端主动重试（允许新结果） |
| HaveItems | 回放时不更新（RecordWithoutUpdate） | 重抽时按新结果更新 |
| 回放范围 | 可跨多次请求累积 | 只能回退一次 |

---

## 12. 状态管理

### 12.1 UserRealTimeState（Redis 持久化）

| 字段 | 类型 | 说明 |
|------|------|------|
| HaveItemsMore | map[string]int | 历史已抽中道具及数量 |
| HaveItemsMoreToday | map[string]int | 今日已抽中（每日重置）|
| LuckyValue | map[string]int | 各优先级组的幸运值 |
| LastDrawRound | int | 上次抽奖轮次 |
| LastDrawType | int | 上次抽奖次数 |
| ProbID | string | 当前概率组 ID |
| PoolID | string | 当前奖池 ID |
| PoolItemsMap | map[string]bool | 奖池道具集合（无放回时动态变化）|
| RewardInfo | []RewardRecord | 最近中奖记录（回放用）|
| TopRewardNum | int | 大奖产出计数（监控用）|

### 12.2 双缓存结构

```go
type DoubleCacheUserRealTimeState struct {
    CurrentState string  // 当前状态 JSON
    PrevState    string  // 上一次状态 JSON（用于回退重抽）
}
```

---

## 13. 消息队列设计

### 13.1 Topic 设计

| Topic | 峰值消息速率 | 分区数 | 刷盘策略 | 用途 | 消费者 |
|-------|------------|--------|---------|------|--------|
| `LOTTERY_DELIVERY` | 50万条/s | **64** | 同步刷盘 | 抽奖发货（道具不能丢） | 发货服务 |
| `LOTTERY_RECORD` | 50万条/s | **32** | 异步刷盘 | 抽奖流水异步落库 | 记录服务 |
| `LOTTERY_BROADCAST` | 1万条/s | **8** | 异步刷盘 | 大奖全服广播 | 通知服务 |
| `LOTTERY_STAT` | 50万条/s | **16** | 异步刷盘 | 概率统计/产出监控 | 统计服务 |
| `LOTTERY_DLQ` | 极低 | **4** | 同步刷盘 | 死信兜底 | 告警+人工 |

### 13.2 消息可靠性

**生产者端（事务消息）：**
```
半消息 → 本地事务（扣货币+写记录+更状态） → Commit/Rollback
宕机时 RocketMQ 自动回查 → 查 draw_record 表判断事务结果
```

**消费者端（幂等消费）：**
```go
func consumeDelivery(msg *DeliveryMsg) error {
    // 幂等检查：delivery_record 唯一索引 uk_draw_item
    exists, err := checkDeliveryRecord(msg.DrawID, msg.ItemIndex)
    if err != nil { return err }
    if exists { return nil }  // 已处理，ACK

    // 执行发货
    result, err := deliverToGameServer(msg)
    if err != nil { return err }  // 失败，触发重试

    // 写发货记录
    insertDeliveryRecord(msg, result)
    return nil
}
```

---

## 14. 缓存架构

```
L1 本地缓存（进程内 go-cache）:
   ├── pool_config_{pool_id}_{version}    奖池配置，TTL=5min
   ├── alias_table_{prob_id}_{version}    Alias表，TTL=5min
   ├── item_config_{item_id}              道具配置，TTL=10min
   └── 命中率目标：90%（配置类数据，变更频率低）

L2 Redis Cluster（32分片）:
   ├── state:{game_id}:{user_key}[_{layer}]  用户实时状态（权威源）
   ├── stock:{pool_id}:{item_id}              限量库存计数
   ├── limit:{item_id}:{scope}                限量累计计数
   ├── rsp_cache:{user_key}[_{layer}]         响应缓存
   ├── lock:{user_key}                        分布式锁
   └── alias_table:{pool_id}:{prob_id}:{ver}  概率Alias表缓存

L3 MySQL（最终持久化）:
   └── 抽奖记录、发货记录、配置表 —— 核心链路不直连
```

---

## 15. 容错性设计

### 15.1 限流

| 层次 | 维度 | 阈值 | 动作 |
|------|------|------|------|
| 网关全局 | 总流量 | 200万 QPS | 返回 503 |
| 游戏维度 | 单游戏 | 50万 QPS | 排队等待 |
| 用户维度 | 单 uid | 10次/s | 返回频率限制 |
| 奖池维度 | 单奖池 | 配置上限 | 排队等待 |

### 15.2 熔断与降级

```
异常触发（满足任一）:
├─ Redis P99 > 30ms
├─ MQ 发货消息堆积 > 5万
├─ 核心接口错误率 > 0.5%
└─ 概率偏离 > 30%
        ↓
分级降级:
├─ 一级：关闭全服广播、统计上报，核心抽奖不受影响
├─ 二级：关闭新池开启，集中资源保存量池
├─ 三级：发货切同步模式（RT升高但不丢），禁止连抽（仅允许单抽）
└─ 四级：Redis 全挂 → 切 MySQL 乐观锁模式，QPS 从50万降至5万
```

### 15.3 兜底方案矩阵

| 故障场景 | 兜底策略 | 恢复时序 |
|---------|---------|---------|
| Redis 单分片宕机 | 哨兵自动切主从，<30s 该分片暂停抽奖 | 自动恢复 |
| Redis 集群全挂 | 切 MySQL 乐观锁 + 本地缓存状态 | 手动恢复 |
| RocketMQ 宕机 | 发货切同步 RPC，抽奖流水写本地 WAL | 手动恢复 |
| 游戏服务器不可达 | 发货消息积压，服务恢复后自动消费 | 自动恢复 |
| 限量库存 Redis 丢失 | 从 MySQL delivery_record COUNT 重建 | 自动恢复 |
| 概率异常 | 实时监控告警 + 一键下线奖池 | 人工处理 |

### 15.4 动态配置开关（ETCD，秒级生效）

```yaml
lottery.switch.global: true           # 全局抽奖开关
lottery.switch.delivery_async: true   # 发货异步开关（故障时切同步）
lottery.switch.db_fallback: false     # DB 降级模式
lottery.limit.draw_qps: 500000        # 抽奖总 QPS 上限
lottery.degrade_level: 0              # 降级级别 0~4
```

---

## 16. 监控与合规

### 16.1 监控指标

```yaml
business_metrics:
  - draw_qps: "抽奖QPS（按游戏/奖池/类型）"
  - draw_latency_p99: "抽奖延迟（目标<50ms）"
  - delivery_success_rate: "发货成功率（目标>99.9%）"
  - item_output_rate: "各道具实际产出率（对比配置概率）"
  - pity_trigger_rate: "保底触发率"
  - stock_remaining: "限量道具剩余库存"
  - revenue_realtime: "实时收入"
  
alerts:
  - "概率偏离：|实际-配置|/配置 > 20%，5分钟窗口 → P0"
  - "库存不足：剩余 < 5% → P1"
  - "发货积压：堆积 > 5万 → P0"
  - "抽奖延迟：P99 > 100ms → P1"
  - "错误率：> 0.5% → P0"
```

### 16.2 合规要求

```yaml
compliance:
  probability_disclosure:
    public_api: "/api/v1/pool/{pool_id}/probability"
    precision: 4
    update_on_change: true
    
  minor_protection:
    age_0_8: "ban"
    age_8_16: "limit_200_monthly"
    age_16_18: "limit_400_monthly"
    
  audit_log:
    retention_days: 365
    immutable: true
```

---

## 17. 数据库设计

```sql
-- 奖池配置表
CREATE TABLE t_pool_config (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    pool_id     VARCHAR(64) NOT NULL,
    version     INT NOT NULL DEFAULT 1,
    pool_name   VARCHAR(128) NOT NULL,
    pool_type   VARCHAR(32) NOT NULL,
    game_id     VARCHAR(32) NOT NULL,
    config_json JSON NOT NULL COMMENT '完整奖池配置',
    status      TINYINT NOT NULL DEFAULT 0 COMMENT '0草稿 1审核中 2生效 3下线',
    start_time  DATETIME NOT NULL,
    end_time    DATETIME NOT NULL,
    created_by  VARCHAR(64) NOT NULL,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_pool_version (pool_id, version),
    KEY idx_game_status (game_id, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- 抽奖记录表（按月分表，按 user_id % 16 分库）
CREATE TABLE t_draw_record_202401 (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    draw_id         VARCHAR(64) NOT NULL COMMENT '雪花ID',
    batch_id        VARCHAR(64) NOT NULL COMMENT '连抽批次ID',
    user_id         BIGINT UNSIGNED NOT NULL,
    game_id         VARCHAR(32) NOT NULL,
    pool_id         VARCHAR(64) NOT NULL,
    pool_version    INT NOT NULL,
    draw_index      SMALLINT NOT NULL COMMENT '连抽序号',
    draw_count      SMALLINT NOT NULL COMMENT '批次总数',
    item_id         VARCHAR(64) NOT NULL,
    rarity          VARCHAR(16) NOT NULL,
    prob_id         VARCHAR(64) NOT NULL,
    is_pity         TINYINT NOT NULL DEFAULT 0,
    cost_type       VARCHAR(32) NOT NULL,
    cost_amount     INT NOT NULL,
    pity_counter    INT NOT NULL DEFAULT 0,
    draw_time       DATETIME(3) NOT NULL,
    UNIQUE KEY uk_draw_id (draw_id),
    KEY idx_user_time (user_id, draw_time),
    KEY idx_pool_time (pool_id, draw_time),
    KEY idx_batch (batch_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- 发货记录表（按 draw_id % 16 分库）
CREATE TABLE t_delivery_record (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    delivery_id     VARCHAR(64) NOT NULL,
    draw_id         VARCHAR(64) NOT NULL,
    item_index      SMALLINT NOT NULL,
    user_id         BIGINT UNSIGNED NOT NULL,
    item_id         VARCHAR(64) NOT NULL,
    item_type       VARCHAR(32) NOT NULL,
    channel         VARCHAR(32) NOT NULL,
    status          TINYINT NOT NULL DEFAULT 0 COMMENT '0待发 1发送中 2成功 3失败 4补偿',
    retry_count     TINYINT NOT NULL DEFAULT 0,
    channel_order_id VARCHAR(128),
    error_msg       VARCHAR(512),
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_draw_item (draw_id, item_index) COMMENT '幂等唯一索引',
    KEY idx_user_status (user_id, status),
    KEY idx_status_time (status, updated_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- 库存表
CREATE TABLE t_inventory (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    pool_id     VARCHAR(64) NOT NULL,
    item_id     VARCHAR(64) NOT NULL,
    total_stock INT NOT NULL DEFAULT -1 COMMENT '-1无限',
    used_stock  INT NOT NULL DEFAULT 0,
    version     INT NOT NULL DEFAULT 0,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_pool_item (pool_id, item_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


-- 用户状态持久化（Redis 宕机恢复用）
CREATE TABLE t_user_state (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id     BIGINT UNSIGNED NOT NULL,
    game_id     VARCHAR(32) NOT NULL,
    pool_type   VARCHAR(32) NOT NULL,
    layer       VARCHAR(16) NOT NULL DEFAULT '',
    state_json  MEDIUMTEXT NOT NULL COMMENT 'zlib压缩后的状态JSON',
    version     INT NOT NULL DEFAULT 0,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_pool_layer (user_id, game_id, pool_type, layer)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

---

## 18. 典型业务场景配置

### 原神式角色UP池
```yaml
milestone: [{priority:1, lucky_value:90, prob_id:"ssr_100"}, {priority:2, lucky_value:73, slide:true}]
up_guarantee: "50_50_then_100"
carry_over: true
```

### 王者荣耀荣耀水晶（361保底）
```yaml
milestone: [{priority:1, lucky_value:361, prob_id:"self_select_crystal"}]
no_replacement: { scope: "crystal_items" }
```

### 集卡活动（无放回）
```yaml
lottery_type: "NoReplacement"
# FilterItemByHaveItems + FilterItemByDrawn 联合去重
# 拥有8张后 FixProbPlugin 提高未有卡概率×3
```

### 转盘/九宫格
```yaml
pool_type: "static_weighted"
# 8个道具直接配权重，Alias Method 采样
# 近Miss 由客户端本地实现
```

---

## 19. 技术栈

| 类别 | 技术选型 |
|------|---------|
| 语言 | Go (tRPC-Go 微服务框架) |
| 存储 | Redis Cluster + MySQL |
| 消息队列 | RocketMQ (事务消息) |
| 配置中心 | ETCD |
| 序列化 | json-iterator + Protobuf |
| 压缩 | zlib (用户状态) |
| 监控 | Prometheus + Grafana + TDBank |
| 日志 | ELK |

---

## 20. 面试常见问题

### Q1: 如何保证限量道具绝对不超发？

**答**：三层防超发闭环：

1. **Redis Lua 原子扣减（第一层，实时拦截）**：抽奖命中限量道具时执行 Lua 脚本 `DECRBY + 判断 ≥ 0`，整个操作在 Redis 单线程中原子执行。返回 -1 立即拒绝，不进入后续流程——从源头杜绝超发。
2. **发货 Worker 二次校验 MySQL（第二层，消费端兜底）**：发货消费者在实际调用游戏服务器前，执行 `SELECT used_stock FROM t_inventory WHERE item_id=? FOR UPDATE`，若已达上限则拒绝发货并回补 Redis。这一层防的是 Redis 主从切换丢数据的极端场景。
3. **定时对账修正（第三层，最终兜底）**：每5分钟对比 `redis_stock` 与 `total_stock - COUNT(delivery WHERE success)`，以 MySQL 为准修正 Redis 漂移。每日凌晨全量对账。

**关键设计**：扣减失败时立即 `INCRBY` 回补，保证 Redis 库存值的最终准确。DB 事务失败也回补，形成闭环。

---

### Q2: 抽奖接口如何做到 P99 < 50ms？

**答**：五个关键优化点：

1. **O(1) 采样**：Alias Method 预构建概率表，每次抽奖只需生成一个随机数 + 一次数组查找，不随道具数量增长。
2. **发货异步化**：抽奖接口只负责"决定结果"，不等发货完成。结果确定后立即返回客户端，发货走 RocketMQ 事务消息异步处理。
3. **响应缓存前置**：响应缓存检查在**分布式锁之前**执行，命中缓存直接返回，连 Redis SETNX 都不需要。80% 重试请求在此拦截。
4. **状态全内存化**：UserRealTimeState 完整存储在 Redis（不查 MySQL），单次抽奖只需 3 次 Redis 操作（读状态、扣库存、写状态），Pipeline 可合并为 1 次 RTT。
5. **本地缓存配置**：奖池配置、Alias 表、道具表等读多写少数据全部缓存在进程内（go-cache），命中率 > 90%，避免远程调用。

---

### Q3: RocketMQ 事务消息如何保证"抽奖扣费成功 ↔ 道具一定发到"？

**答**：RocketMQ 事务消息的三阶段保证：

1. **发送半消息**：抽奖服务先发一条半消息（对消费者不可见），MQ 返回 ACK。
2. **执行本地事务**：在同一个 MySQL 事务中完成"扣货币 + 写抽奖记录 + 更新保底计数"。
   - 事务成功 → Commit 半消息（消费者可见，发货 Worker 开始消费）
   - 事务失败 → Rollback 半消息（消息丢弃，用户无感知）
3. **事务回查**：若 Producer 在 Commit/Rollback 前宕机，RocketMQ 会定时回查。Producer 查询 `draw_record` 表：记录存在 → Commit，不存在 → Rollback。

**与 TCC 的区别**：不需要 Try/Confirm/Cancel 三套代码，实现简单；发货天然是异步的，不需要资源预留；回查机制自动兜底宕机场景。

**消费端幂等**：`delivery_record` 表的 `uk_draw_item(draw_id, item_index)` 唯一索引做最终兜底，重复消费直接 ACK。

---

### Q4: 保底机制（里程碑）的技术实现原理是什么？

**答**：基于"幸运值 + 优先级分组 + 概率组切换"的通用框架：

1. **幸运值累加**：每次抽奖后，`LuckyValuePlugin.UpdateLuckyValue` 将用户对应优先级组的幸运值 +1。
2. **里程碑匹配**：每次抽奖前，`ProbIdGenPlugin.GenProbId` 调用 `LuckyValuePlugin.Match`，按优先级遍历里程碑规则。
3. **概率组切换**：当幸运值 ≥ 阈值时，切换到高稀有度概率组（如100% SSR）。
4. **清零重置**：抽中大奖后，`UpdateLuckyValue` 根据 `clean_rules` 清零对应优先级组。

**示例（原神式）**：
- 优先级1：幸运值=90 → 切到 `ssr_100`（硬保底，必出SSR）
- 优先级2：幸运值=73 → 切到 `ssr_soft`（软保底，启用线性滑动，概率递增）
- 优先级3：幸运值=10 → 切到 `sr_100`（小保底，必出SR）

**跨池继承**：幸运值按 `pool_type` 隔离而非 `pool_id`，换 UP 角色不清零保底计数。

---

### Q5: 多维度道具过滤后概率如何重分配？Fixed/Unfixed 的设计意图？

**答**：

**问题背景**：无放回抽奖中，道具被过滤后剩余道具的概率必须重新分配。如果简单等比放大，会导致保底道具/兜底道具的概率被意外提升。

**Fixed/Unfixed 双轨设计**：
- **Fixed 道具**：概率恒定，不受过滤影响。典型：保底概率组中的稀有道具。
  - 实际概率 = 配置权重 / 总权重（始终不变）
- **Unfixed 道具**：过滤后剩余道具按比例缩放瓜分"Unfixed 总概率份额"。
  - 实际概率 = 配置权重 × (1 - Fixed总概率) / 剩余Unfixed总权重

**举例**：概率组有 Fixed 道具A(1%) 和 Unfixed 道具B(30%)、C(30%)、D(39%)。当 C 被过滤：
- A 仍然是 1%（Fixed 恒定）
- B 的新概率 = 30% × (1-1%) / (30%+39%) = 30% × 99% / 69% ≈ 43%
- D 的新概率 = 39% × 99% / 69% ≈ 56%

**池空兜底**：当所有 Unfixed 道具都被过滤，自动切换到 `complement` 兜底概率组（通常是金币/经验等通用奖励），保证每次抽奖必有结果。

---

### Q6: 系统如何支撑新池开启瞬间100万QPS的流量洪峰？

**答**：分层防御 + 预计算 + 异步化：

1. **网关限流**：总流量硬上限 200万 QPS，单用户 10次/s，超出直接拒绝。
2. **响应缓存拦截**：同一用户同一轮次重复请求命中缓存直接返回，拦截 80% 重试流量。实际穿透到抽奖逻辑的约 50万 QPS。
3. **Alias 表预计算**：奖池上线时预构建所有概率组的 Alias 表缓存到 Redis，抽奖时直接加载，无需实时计算。
4. **Redis 分片打散**：32 分片集群，用户状态按 user_id 哈希均匀分布，无热点。
5. **发货完全异步**：抽奖结果确定后立即返回，发货消息积压在 MQ，由 Worker 按 DB 能承受的速率消费（MQ 削峰）。
6. **预扩容**：大型活动前 3 天按 3 倍峰值预扩容服务节点，活动后自动缩容。
7. **HPA 弹性**：K8s HPA 基于 CPU > 60% 自动扩容，应对未预期的流量突增。

---

### Q7: 连抽（十连抽）的设计有哪些需要注意的点？

**答**：

1. **原子性**：N 连抽在同一次请求中完成，共享同一个分布式锁，保证中间不会被其他请求插入。
2. **连抽内保底**：通过 `AfterRuleDrawPlugin` 实现"10连内必出SR"——每轮抽奖后检查本批次已抽列表，若剩余次数内未出SR，最后一轮强制切到 SR 概率组。
3. **连抽内去重**：`FilterItemByDrawn` 过滤器记录本次连抽已抽道具，配置 `FilterDrawnSource2Target` 映射规则，防止同一 SSR 连抽中出现多次。
4. **限量竞争**：连抽中每轮独立执行 Redis Lua 扣减，某轮限量失败只重抽该轮，不影响已完成的轮次。
5. **中断机制**：`AfterRuleDrawPlugin` 返回 `isStopDraw=true` 可提前终止连抽（如：抽到大奖后停止）。
6. **折扣计算**：连抽优惠在请求入口统一计算（如十连只收9抽的价格），不影响内部抽奖逻辑。
7. **回放兼容**：客户端崩溃后重进，回放模式按连抽序列逐个回放，不重新随机。

---

### Q8: Redis 宕机后如何恢复用户状态？会丢数据吗？

**答**：

**数据安全架构**：
- **Redis AOF + RDB 持久化**：AOF every second 策略，最多丢失1秒数据。
- **主从复制**：每分片 1主2从，哨兵自动切换，切换期间（<30s）该分片暂停抽奖。
- **MySQL 异步持久化**：UserRealTimeState 每次变更后异步写入 `t_user_state` 表。

**恢复流程**：
1. Redis 主从切换后自动恢复（大多数场景，数据不丢）。
2. 若 Redis 数据丢失（极端场景）：从 `t_user_state` 表加载最近一次持久化的状态。
3. 可能丢失的：最后1-2次抽奖的保底计数。但由于"宁可少扣保底不可多扣"，用户保底进度可能比实际少1-2次，下次抽奖时会被多给一些保底余量。
4. 限量库存恢复：从 `delivery_record` 表 COUNT 实际已成功发货数量，`redis_stock = total - used`。

**关键原则**：Redis 是用户状态权威源，MySQL 是兜底恢复源。正常运行时不查 MySQL，只在 Redis 不可用时降级。

---

### Q9: 如何防止用户通过外挂/脚本作弊刷概率？

**答**：多层风控体系：

1. **频率限制**：单用户 10次/s 硬限（网关层），单奖池每日上限（业务层）。超频直接拒绝。
2. **请求签名**：客户端对关键参数（user_id, pool_id, draw_count, timestamp, nonce）做 HMAC-SHA256 签名，服务端验签。timestamp 容忍30秒偏差，nonce 60秒内去重。
3. **行为检测**：
   - 固定间隔检测：50次抽奖时间间隔标准差 < 5% → 疑似脚本，强制增加5秒冷却。
   - 短时爆发：5分钟内 > 500次 → 暂停抽奖 + 人工审核。
4. **概率监控**：实时监控个人 SSR 率，超过 3σ 标记异常用户。
5. **设备指纹**：同设备5个以上账号抽奖 → 标记全部关联账号。
6. **服务端决定论**：所有随机数在服务端生成，客户端无法影响结果。Alias 表不下发客户端。
7. **结果不可预测**：概率组版本号 + 服务端种子 + 用户状态共同决定结果，无法通过抓包预判。

---

### Q10: 插件化架构如何做到"新游戏接入只需3天"？

**答**：

**接入流程**：
1. **配置 busi.yaml**（半天）：确定用户标识组合（openid/roleid/partition）、功能开关（响应缓存/回放/自定义奖池）。
2. **配置抽奖规则**（1天）：通过运营后台配置道具表、概率组、奖池、限量规则、里程碑——纯配置，无需代码。
3. **实现差异插件**（1天）：只覆盖与默认行为不同的插件。典型：
   - LOL 需要定制 `BuildItemPlugin`（道具格式不同）
   - 剑侠需要定制 `PoolIdGenPlugin`（按服务器等级选池）
   - 大部分游戏0个自定义插件（默认实现够用）
4. **启动服务**：`server.NewLotteryServer(&MyPlugin{})` 一行代码启动。

**插件覆盖机制**：
```go
// 默认插件先注册（提供基础能力）
plugindef.Add(&DefaultBuildItemPlugin{})
// 自定义插件后注册（同名覆盖）
plugindef.Add(&LOLBuildItemPlugin{})  // Name() 返回相同名称，自动覆盖
```

**核心设计**：15+ 扩展点覆盖了抽奖流程的所有环节（概率选择→过滤→采样→替换→响应构建），任何差异化需求都有对应插件点。默认实现覆盖80%场景，自定义插件覆盖剩余20%——新业务不需要理解整个系统，只需理解要覆盖的那1-2个插件接口。
