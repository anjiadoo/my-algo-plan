# 互娱场景通用抽奖系统设计
> 互娱场景通用抽奖系统，采用插件化架构，支持有放回/无放回抽奖、分层奖池、里程碑保底、多维度道具过滤、连抽、概率缩放、异步发货、A/B 实验等能力，适用于腾讯游戏旗下多款游戏的各类抽奖活动场景。

---

## 10个关键技术决策

| 决策 | 选择 | 核心理由 |
|------|------|---------|
| **加权随机采样 + Alias Method** | 预构建 Alias 表，每次抽奖 O(1) | 概率组道具数量可能达数百个，线性扫描 O(n) 无法满足百万 QPS；Alias Method 预处理 O(n) 后每次采样 O(1)，内存换时间 |
| **Redis Lua 原子扣减限量库存** | DECRBY + 判断 ≥ 0，失败立即 INCRBY 回补 | 限量道具全局并发竞争，分布式锁性能差；Lua 脚本单线程原子执行天然防超发，单分片承载 10万+ QPS |
| **RocketMQ 事务消息异步发货** | 半消息 + 本地事务 + 回查 | 发货涉及多个游戏服务器 RPC，同步调用会拖垮抽奖 RT；事务消息保证"本地扣减成功 ↔ 发货消息必达"的原子性 |
| **主流程插件化（Hook 机制）** | 核心流程固定，15+ Hook 点可被业务插件覆盖 | 多款游戏玩法差异大（LOL/王者/原神式），硬编码无法维护；插件后注册覆盖先注册，新业务只需实现差异插件 |
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

**实体（Entity，写模型）**

| 模型 | 职责 | 核心属性 | 存储位置 |
|------|------|---------|---------|
| **LotteryPool** 奖池 | 奖池全生命周期：配置→生效→下线 | 奖池ID、游戏ID、版本号、类型、时间范围、层级配置、状态 | MySQL 配置表 + Redis 缓存 |
| **UserLotteryState** 用户抽奖状态 | 用户在某奖池的完整进度 | 用户ID、奖池ID、Layer、已抽道具、幸运值、保底计数、当前概率组ID、奖池道具集合 | Redis（权威源）+ MySQL 异步持久化 |
| **DrawRecord** 抽奖记录 | 每次抽奖的完整流水 | 抽奖ID(雪花)、批次ID、用户ID、奖池ID、道具ID、稀有度、是否保底、概率快照 | MySQL 按月分表 |
| **Inventory** 限量库存 | 限量道具的全局库存管控 | 奖池ID、道具ID、总库存、已发放、Redis 实时值 | Redis 原子计数（权威）+ MySQL 对账兜底 |
| **DeliveryRecord** 发货记录 | 道具发放追踪 | 发货ID、抽奖ID、用户ID、道具ID、渠道、状态、重试次数 | MySQL |

**事件（Event，事件流）**

| 事件 | 触发时机 | 下游消费 |
|------|---------|---------|
| **DrawCompleted** 抽奖完成 | 本地事务提交成功 | ① 异步发货 ② 写抽奖流水 ③ 统计上报 |
| **DeliverySuccess** 发货成功 | 游戏服务器确认到账 | ① 更新发货状态 ② 推送到账通知 ③ 大奖广播 |
| **DeliveryFailed** 发货失败 | 重试耗尽 | ① 人工工单 ② 补偿触发 |
| **StockDepleted** 库存耗尽 | Redis 库存 ≤ 0 | ① 运营告警 ② 自动切换兜底道具 |

**模型关系图**

```
  [写路径（抽奖核心）]              [事件流]                    [读路径]
  ┌──────────────────┐                                    ┌──────────────────┐
  │ UserLotteryState │──DrawCompleted─────┐               │   DrawRecord     │ ← 抽奖流水
  │ (Redis 状态机)    │                    │               │  (MySQL 按月分表) │
  │ ├─ 幸运值/保底     │                    │              └───────────────────┘
  │ ├─ 已抽道具       │                    │               ┌──────────────────┐
  │ └─ 奖池道具集合    │                    ├─ RocketMQ──→  │  DeliveryRecord  │ ← 发货记录
  └──────────────────┘                    │               └──────────────────┘
         │                                │               ┌──────────────────┐
         │ 库存扣减                        │                │  到账通知/广播    │
         ▼                                │               └──────────────────┘
  ┌──────────────────┐                    │
  │   Inventory      │──StockDepleted─────┘
  │ (Redis 原子计数)  │
  │ MySQL 对账兜底    │
  └──────────────────┘

  ┌──────────────────┐
  │  LotteryPool     │ ← 运营后台配置，热加载到 Redis
  │ (MySQL + Redis)  │
  └──────────────────┘
```

**设计原则：**
- **写路径极简**：抽奖核心只有"Redis 读状态 + Lua 扣库存 + Redis 写状态 + 事务消息"，不同步写 DB
- **事件必然对应记录**：每次抽奖成功事件必落 DrawRecord + DeliveryRecord，通过事务消息保证不丢
- **Redis 权威 + MySQL 兜底**：用户状态/限量库存热路径靠 Redis（50万 QPS），MySQL 只做持久化兜底（对账/恢复）
- **幂等最终兜底**：DeliveryRecord 的 `uk_draw_item(draw_id, item_index)` 唯一索引是"一次抽奖一个道具只发一次"的最终防线

### 完整库表设计

```sql
-- =====================================================
-- 奖池配置表（单库，读多写少，运营后台管理）
-- =====================================================
CREATE TABLE t_pool_config (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    pool_id     VARCHAR(64) NOT NULL  COMMENT '奖池唯一标识',
    version     INT NOT NULL DEFAULT 1 COMMENT '配置版本号，每次修改+1',
    pool_name   VARCHAR(128) NOT NULL,
    pool_type   VARCHAR(32) NOT NULL  COMMENT 'layered/static/dynamic/collection',
    game_id     VARCHAR(32) NOT NULL  COMMENT '所属游戏ID',
    config_json JSON NOT NULL         COMMENT '完整奖池配置（层级/道具/概率组/里程碑/过滤规则）',
    status      TINYINT NOT NULL DEFAULT 0 COMMENT '0草稿 1审核中 2生效 3下线',
    start_time  DATETIME NOT NULL,
    end_time    DATETIME NOT NULL,
    created_by  VARCHAR(64) NOT NULL,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_pool_version (pool_id, version),
    KEY idx_game_status (game_id, status),
    KEY idx_status_time (status, start_time, end_time) COMMENT '活动时间查询索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='奖池配置表';


-- =====================================================
-- 抽奖记录表（按 user_id % 16 分16库，按月分表）
-- 核心：记录每一次抽奖的完整快照，用于审计/查询/对账
-- =====================================================
CREATE TABLE t_draw_record_202401 (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    draw_id         VARCHAR(64) NOT NULL  COMMENT '雪花ID，全局唯一',
    batch_id        VARCHAR(64) NOT NULL  COMMENT '连抽批次ID，同一次N连抽共享',
    user_id         BIGINT UNSIGNED NOT NULL,
    game_id         VARCHAR(32) NOT NULL,
    pool_id         VARCHAR(64) NOT NULL,
    pool_version    INT NOT NULL          COMMENT '抽奖时的奖池版本（可追溯配置）',
    draw_index      SMALLINT NOT NULL     COMMENT '连抽中的序号 0~N-1',
    draw_count      SMALLINT NOT NULL     COMMENT '本批次连抽总数',
    item_id         VARCHAR(64) NOT NULL,
    item_name       VARCHAR(128) NOT NULL,
    rarity          VARCHAR(16) NOT NULL  COMMENT 'SSR/SR/R/N',
    prob_id         VARCHAR(64) NOT NULL  COMMENT '使用的概率组ID',
    is_pity         TINYINT NOT NULL DEFAULT 0 COMMENT '是否保底触发 0否 1是',
    is_up           TINYINT NOT NULL DEFAULT 0 COMMENT '是否UP道具',
    cost_type       VARCHAR(32) NOT NULL  COMMENT '消耗类型：diamond/ticket/coin',
    cost_amount     INT NOT NULL          COMMENT '消耗数量（整数）',
    pity_counter    INT NOT NULL DEFAULT 0 COMMENT '抽奖前的保底计数（用于分析）',
    draw_time       DATETIME(3) NOT NULL,
    extra_info      JSON                  COMMENT '扩展信息（AB分组/客户端版本等）',
    UNIQUE KEY uk_draw_id (draw_id),
    KEY idx_user_time (user_id, draw_time),
    KEY idx_pool_time (pool_id, draw_time),
    KEY idx_batch (batch_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='抽奖记录表（按月分表）';


-- =====================================================
-- 发货记录表（按 draw_id % 16 分16库）
-- 核心：uk_draw_item 联合唯一索引是防重复发货的数据库最终兜底
-- =====================================================
CREATE TABLE t_delivery_record (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    delivery_id     VARCHAR(64) NOT NULL  COMMENT '发货单雪花ID',
    draw_id         VARCHAR(64) NOT NULL  COMMENT '关联抽奖记录ID',
    item_index      SMALLINT NOT NULL     COMMENT '连抽中的道具序号',
    user_id         BIGINT UNSIGNED NOT NULL,
    item_id         VARCHAR(64) NOT NULL,
    item_type       VARCHAR(32) NOT NULL  COMMENT 'character/skin/currency/coupon/physical',
    channel         VARCHAR(32) NOT NULL  COMMENT '发货渠道：game_server_rpc/wallet/coupon_center',
    status          TINYINT NOT NULL DEFAULT 0
                    COMMENT '0待发 1发送中 2成功 3失败 4补偿中 5补偿成功',
    retry_count     TINYINT NOT NULL DEFAULT 0,
    channel_order_id VARCHAR(128)         COMMENT '渠道方返回的单号（对账用）',
    error_msg       VARCHAR(512),
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_draw_item (draw_id, item_index) COMMENT '幂等唯一索引（防重复发货最终兜底）',
    KEY idx_user_status (user_id, status),
    KEY idx_status_time (status, updated_at) COMMENT '发货Worker重试/补偿扫描索引'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='发货记录表';


-- =====================================================
-- 库存表（单库，限量道具不多，写入量可控）
-- Redis 原子计数为实时权威，此表为对账基准+恢复来源
-- =====================================================
CREATE TABLE t_inventory (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    pool_id     VARCHAR(64) NOT NULL,
    item_id     VARCHAR(64) NOT NULL,
    total_stock INT NOT NULL DEFAULT -1  COMMENT '-1表示无限库存',
    used_stock  INT NOT NULL DEFAULT 0   COMMENT '已发放数量（异步更新，对账用）',
    frozen_stock INT NOT NULL DEFAULT 0  COMMENT '冻结中数量（发货中未确认）',
    version     INT NOT NULL DEFAULT 0   COMMENT '乐观锁（DB降级模式用）',
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_pool_item (pool_id, item_id),
    KEY idx_pool_id (pool_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='限量库存表';


-- =====================================================
-- 用户状态持久化表（按 user_id % 16 分库）
-- 异步落盘：每次抽奖后MQ异步写入，Redis宕机时恢复用
-- =====================================================
CREATE TABLE t_user_state (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id     BIGINT UNSIGNED NOT NULL,
    game_id     VARCHAR(32) NOT NULL,
    pool_type   VARCHAR(32) NOT NULL     COMMENT '按池类型隔离（跨池保底继承）',
    layer       VARCHAR(16) NOT NULL DEFAULT '' COMMENT 'Layer层级标识',
    state_json  MEDIUMTEXT NOT NULL       COMMENT 'UserRealTimeState JSON（可选zlib压缩）',
    version     INT NOT NULL DEFAULT 0   COMMENT '状态版本号（每次更新+1）',
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_pool_layer (user_id, game_id, pool_type, layer),
    KEY idx_user_game (user_id, game_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户抽奖状态持久化表（Redis兜底）';


-- =====================================================
-- 发货事务补偿表（幂等 + MQ 回查兜底）
-- 类比红包系统的 send_transaction，用于 RocketMQ 事务回查
-- =====================================================
CREATE TABLE t_draw_transaction (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    draw_id         VARCHAR(64) NOT NULL  COMMENT '抽奖ID = 事务消息的唯一标识',
    batch_id        VARCHAR(64) NOT NULL,
    user_id         BIGINT UNSIGNED NOT NULL,
    pool_id         VARCHAR(64) NOT NULL,
    draw_count      SMALLINT NOT NULL     COMMENT '本次连抽次数',
    cost_type       VARCHAR(32) NOT NULL,
    cost_amount     INT NOT NULL          COMMENT '总消耗金额',
    status          TINYINT NOT NULL DEFAULT 0
                    COMMENT '0处理中(半消息已发) 1已提交(本地事务成功) 2已回滚(本地事务失败)',
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_draw_id (draw_id) COMMENT '事务回查时按此字段定位',
    KEY idx_status_create (status, created_at) COMMENT '定时补偿任务扫描（status=0超时未确认的）'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='抽奖事务补偿表（RocketMQ回查用）';


-- =====================================================
-- 日级对账任务表（兜底）
-- =====================================================
CREATE TABLE t_reconciliation_task (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    task_date       DATE NOT NULL         COMMENT '对账日期',
    pool_id         VARCHAR(64) NOT NULL,
    item_id         VARCHAR(64) NOT NULL,
    redis_stock     INT NOT NULL          COMMENT '对账时Redis库存值',
    mysql_used      INT NOT NULL          COMMENT 'MySQL统计已发放数',
    expected_remain INT NOT NULL          COMMENT '预期剩余 = total - mysql_used',
    diff            INT NOT NULL          COMMENT '差异 = redis_stock - expected_remain',
    status          TINYINT NOT NULL DEFAULT 0 COMMENT '0待处理 1已修正 2已忽略 3人工处理',
    fix_action      VARCHAR(256)          COMMENT '修正动作描述',
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_date_pool_item (task_date, pool_id, item_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='库存对账任务表';
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
│  │限流熔断 │ │身份认证    │ │参数校验     │ │灰度路由   │  │协议转换   │  │
│  └────────┘ └──────────┘ └────────────┘ └──────────┘ └──────────┘  │
└─────────────────────────────┬───────────────────────────────────────┘
                              │
┌─────────────────────────────▼───────────────────────────────────────┐
│                   抽奖核心服务 (LotteryTemplate)                       │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │              抽奖主流程 + Hook 插件链                             │   │
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
│  │ Redis Cluster│  │  RocketMQ    │  │    ETCD      │              │
│  │ 32分片1主2从  │   │ 16主16从     │  │  配置中心     │              │
│  └──────────────┘  └──────────────┘  └──────────────┘              │
└─────────────────────────────┬───────────────────────────────────────┘
                              │
┌─────────────────────────────▼───────────────────────────────────────┐
│                        存储层                                         │
│  ┌──────────┐  ┌──────────┐ ┌──────────┐ ┌──────────┐              │
│  │MySQL 16库 │ │ ES 日志   │ │Prometheus│ │  HDFS    │              │
│  │按月分表   │  │ 审计查询   │ │+ Grafana │ │ 冷数据   │              │
│  └──────────┘  └──────────┘ └──────────┘ └──────────┘              │
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

### 7.1 概率组嵌套（子概率组）

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

### 7.2 随机算法选型：业界实践与决策

**结论先行**：本系统采用 **CSPRNG 安全伪随机（crypto/rand seed + xoshiro256**）+ Alias Method 加权采样**，不使用真随机、不使用简单伪随机（math/rand 默认种子）。

**业界主流方案对比**：

| 方案 | 代表产品 | 原理 | 性能 | 公平性 | 本系统是否采用 |
|------|---------|------|------|--------|--------------|
| **CSPRNG 种子 + 快速 PRNG** | 原神/王者/LOL/绝大多数商业游戏 | 启动时用 crypto/rand 获取高质量种子，运行时用 xoshiro256**/xorshift128+ 快速生成 | **极高**（~1ns/次） | 统计均匀，不可预测 | ✅ **采用** |
| 硬件真随机 (HRNG) | 彩票/博彩（法律要求） | 热噪声/放射性衰变物理源 | 极低（需硬件设备） | 物理随机 | ❌ 不采用 |
| /dev/urandom 每次调用 | 部分安全敏感场景 | 内核熵池 | 低（系统调用开销） | 高 | ❌ 不采用 |
| math/rand 默认种子 | 教学Demo | 固定种子/时间种子 | 高 | **可预测/可复现** | ❌ 不采用 |
| 区块链可验证随机 (VRF) | Web3 链游 | 链上验证随机数生成过程 | 极低（需上链） | 可验证 | ❌ 不采用 |

**为什么选 CSPRNG 种子 + xoshiro256**（业界最佳实践）**：

1. **不可预测性**：种子来自 `crypto/rand`（Linux 下读 `/dev/urandom` + getrandom 系统调用），攻击者无法通过观察输出反推后续序列。即使玩家抓包获取历史结果，也无法预测下一次。

2. **极高性能**：xoshiro256** 单次生成仅需 ~1ns（几条位运算指令），50万 QPS × 每次抽奖平均 3 次采样 = 150万次/s 随机数生成，CPU 开销可忽略。

3. **统计质量**：通过 BigCrush/PractRand 等统计测试套件，周期 2^256，均匀性远超业务需要。

4. **可审计**：服务端种子可存档（draw_id 关联），监管审计时可"重放验证"——用相同种子+相同状态复现抽奖结果，证明未篡改。

**为什么不选其他方案**：

| 不选方案 | 核心理由 |
|---------|---------|
| **硬件真随机** | ① 游戏不是彩票，不需要物理随机级别的"真"；② 硬件设备增加运维复杂度；③ 性能差 1000 倍以上，无法支撑 50万 QPS |
| **每次调用 /dev/urandom** | ① 系统调用开销 ~500ns/次，是 xoshiro 的 500 倍；② 50万QPS × 3 = 150万次系统调用/s，造成不必要的内核切换；③ 安全性与"CSPRNG seed + fast PRNG"方案相同，但性能远差 |
| **math/rand 默认种子** | ① Go 1.20 之前默认种子固定=1，输出完全可预测；② 即使用 time.Now().UnixNano() 作种子，同一毫秒启动的多个 Goroutine 种子相同；③ 玩家可通过已知结果反推种子（逆向工程） |
| **区块链 VRF** | ① 延迟 1-30s（等区块确认），RT 目标 50ms 完全不可能；② 上链成本高，每次抽奖都要 Gas fee；③ 游戏场景不需要"对外公开可验证"，内部审计即可 |

**实现细节**：

```go
import (
    crand "crypto/rand"
    "encoding/binary"
)

// 服务启动时：从 CSPRNG 获取高质量种子
func newSecureSeed() uint64 {
    var seed [8]byte
    crand.Read(seed[:])  // 从 /dev/urandom 读取
    return binary.LittleEndian.Uint64(seed[:])
}

// xoshiro256** 快速伪随机生成器（每个 Goroutine 独立实例，无锁）
type Xoshiro256ss struct {
    s [4]uint64
}

func NewXoshiro256ss() *Xoshiro256ss {
    return &Xoshiro256ss{s: [4]uint64{
        newSecureSeed(), newSecureSeed(), newSecureSeed(), newSecureSeed(),
    }}
}

func (x *Xoshiro256ss) Uint64() uint64 {
    // xoshiro256** 核心：5条位运算，~1ns
    result := rotl(x.s[1]*5, 7) * 9
    t := x.s[1] << 17
    x.s[2] ^= x.s[0]; x.s[3] ^= x.s[1]
    x.s[1] ^= x.s[2]; x.s[0] ^= x.s[3]
    x.s[2] ^= t
    x.s[3] = rotl(x.s[3], 45)
    return result
}

func rotl(x uint64, k int) uint64 { return (x << k) | (x >> (64 - k)) }

// 生成 [0, n) 的均匀随机数（无模偏差）
func (x *Xoshiro256ss) Intn(n int) int {
    // Lemire's nearly divisionless method（消除模偏差）
    // 比 rand.Intn 更快且无偏
    v := x.Uint64()
    hi, lo := bits.Mul64(v, uint64(n))
    if lo < uint64(n) {
        thresh := -uint64(n) % uint64(n)
        for lo < thresh {
            v = x.Uint64()
            hi, lo = bits.Mul64(v, uint64(n))
        }
    }
    return int(hi)
}
```

**关键设计决策**：
- **每个 Goroutine 独立 RNG 实例**：避免全局 `math/rand` 的 Mutex 锁竞争，50万 QPS 下全局锁是瓶颈
- **Lemire's method 消除模偏差**：`rand() % n` 在 n 不是 2 的幂时有微小偏差，Lemire 方法通过乘法+判断彻底消除
- **种子不可外泄**：种子存在进程内存，不写入日志/Redis/DB，外部无法获取

### 7.3 采样算法对比与选型

#### 7.3.1 游戏抽奖场景特点

在本系统中，采样发生在**多层过滤链之后**——概率组已经是一个纯粹的、过滤好的子集：

```
原始概率组（上千道具）
    │
    ├─ 有效期过滤
    ├─ 已有道具过滤（大R可能过滤掉80%+）
    ├─ 背包过滤
    ├─ 限量过滤
    ├─ 产出优先级过滤
    │
    ▼
过滤后概率组（纯粹的 items + weights）
    │
    ▼
采样算法：从这个干净的概率组里抽一个
```

**关键事实**：
- 过滤后的概率组 **per-user 不同**、**per-request 可能变化**（连抽中抽中一个新道具后下一次就少一个）
- 平均池子大小 ~500，过滤后通常更小
- 性能不是核心考量：n=500 时三种算法耗时都在微秒级，相比整个请求的网络IO（Redis查背包/库存）完全可忽略
- 核心考量：**实现简洁性、可维护性、与动态过滤的契合度**

#### 7.3.2 三种采样算法对比

| 算法 | 预处理 | 单次采样 | 实现复杂度 | 动态过滤友好度 |
|------|--------|---------|-----------|--------------|
| **轮盘赌（线性扫描）** | 无 | O(n) | ⭐ 极简 | ⭐⭐⭐ 天然适配 |
| **前缀和 + 二分查找** | O(n) 构建前缀和 | O(log n) | ⭐⭐ 简单 | ⭐⭐⭐ 适配 |
| **Alias Method** | O(n) 构建别名表 | O(1) | ⭐⭐⭐ 较复杂 | ⭐ 需重建 |

**n=500 时实际耗时对比**（仅供参考，均可忽略）：

| 算法 | 耗时 | 占整个请求(P99=50ms)比例 |
|------|------|------------------------|
| 轮盘赌 | ~0.5μs | 0.001% |
| 前缀和+二分 | ~0.1μs | 0.0002% |
| Alias Method | ~0.02μs | 0.00004% |

---

#### 7.3.3 算法一：轮盘赌（线性扫描）

**原理**：生成 [0, totalWeight) 的随机数，从头到尾逐个累减权重，减到负数时命中。

```
概率轴：
[0 ═══ w0 ═══ w0+w1 ═══ w0+w1+w2 ═══ ... ═══ totalWeight]

rand = 0.55 × totalWeight
从左往右扫，逐个减掉每个道具的权重，哪个让 remain < 0 就选中谁
```

**优点**：
- 零预处理，拿到过滤后的概率组直接就能抽
- 代码极少，不容易出 bug
- 支持 Fixed/Unfixed 两层概率的灵活处理

**完整可运行 Demo**（参考本系统 ProbConf.WeightSampling 简化）：

```go
package main

import (
    crand "crypto/rand"
    "encoding/binary"
    "fmt"
    "math/bits"
)

// ============ 随机数生成器（同 7.5 节） ============

type Xoshiro256ss struct{ s [4]uint64 }

func NewRNG() *Xoshiro256ss {
    var buf [32]byte
    crand.Read(buf[:])
    rng := &Xoshiro256ss{}
    for i := 0; i < 4; i++ {
        rng.s[i] = binary.LittleEndian.Uint64(buf[i*8:])
    }
    return rng
}

func (x *Xoshiro256ss) Uint64() uint64 {
    result := bits.RotateLeft64(x.s[1]*5, 7) * 9
    t := x.s[1] << 17
    x.s[2] ^= x.s[0]; x.s[3] ^= x.s[1]
    x.s[1] ^= x.s[2]; x.s[0] ^= x.s[3]
    x.s[2] ^= t; x.s[3] = bits.RotateLeft64(x.s[3], 45)
    return result
}

func (x *Xoshiro256ss) Float64() float64 {
    return float64(x.Uint64()>>11) / (1 << 53)
}

// ============ 轮盘赌采样（线性扫描） ============

type ProbConf struct {
    items   []string    // 道具ID列表（过滤后）
    weights []float64   // 对应权重（过滤后）
}

// 采样：O(n)，无预处理
func (p *ProbConf) WeightSampling(rng *Xoshiro256ss) string {
    // 计算总权重（过滤后的子集权重和不一定为1）
    var totalWeight float64
    for _, w := range p.weights {
        totalWeight += w
    }

    // 生成 [0, totalWeight) 随机数，线性扫描
    remain := rng.Float64() * totalWeight
    for i, w := range p.weights {
        remain -= w
        if remain < 0 {
            return p.items[i]
        }
    }
    return p.items[len(p.items)-1] // 浮点兜底
}

// ============ 模拟抽奖 ============

func main() {
    // 模拟过滤后的概率组（已过滤掉用户已拥有的道具）
    prob := &ProbConf{
        items:   []string{"SSR龙王", "SR剑姬", "R铁剑", "N金币"},
        weights: []float64{6, 51, 300, 643},
    }

    rng := NewRNG()
    stats := map[string]int{}
    for i := 0; i < 100000; i++ {
        stats[prob.WeightSampling(rng)]++
    }

    fmt.Println("=== 轮盘赌采样 10万次统计 ===")
    for _, item := range prob.items {
        fmt.Printf("%-10s: %5d 次 (%.2f%%)\n", item, stats[item], float64(stats[item])/1000)
    }
    // 输出示例：
    // SSR龙王   :   612 次 (0.61%)
    // SR剑姬    :  5089 次 (5.09%)
    // R铁剑     : 30021 次 (30.02%)
    // N金币     : 64278 次 (64.28%)
}
```

---

#### 7.3.4 算法二：前缀和 + 二分查找

**原理**：先对权重做前缀和（cumulative sum），采样时生成随机数后用二分查找定位落在哪个区间。

```
weights:   [6,    51,   300,  643 ]
prefixSum: [6,    57,   357,  1000]
             │     │      │      │
           idx=0  idx=1  idx=2  idx=3

rand = 0.55 × 1000 = 550
二分查找第一个 prefixSum[i] > 550 → prefixSum[2]=357 ≤ 550, prefixSum[3]=1000 > 550
→ idx=3, 选中 "N金币"
```

**优点**：
- 采样 O(log n)，比线性扫描快（n 大时差异明显）
- 预处理就是一次简单的累加，几乎无开销
- 逻辑清晰，二分查找是标准库函数

**完整可运行 Demo**：

```go
package main

import (
    crand "crypto/rand"
    "encoding/binary"
    "fmt"
    "math/bits"
    "sort"
)

// ============ 随机数生成器（同 7.5 节） ============

type Xoshiro256ss struct{ s [4]uint64 }

func NewRNG() *Xoshiro256ss {
    var buf [32]byte
    crand.Read(buf[:])
    rng := &Xoshiro256ss{}
    for i := 0; i < 4; i++ {
        rng.s[i] = binary.LittleEndian.Uint64(buf[i*8:])
    }
    return rng
}

func (x *Xoshiro256ss) Uint64() uint64 {
    result := bits.RotateLeft64(x.s[1]*5, 7) * 9
    t := x.s[1] << 17
    x.s[2] ^= x.s[0]; x.s[3] ^= x.s[1]
    x.s[1] ^= x.s[2]; x.s[0] ^= x.s[3]
    x.s[2] ^= t; x.s[3] = bits.RotateLeft64(x.s[3], 45)
    return result
}

func (x *Xoshiro256ss) Float64() float64 {
    return float64(x.Uint64()>>11) / (1 << 53)
}

// ============ 前缀和 + 二分查找采样 ============

type PrefixSumSampler struct {
    items     []string
    prefixSum []float64   // prefixSum[i] = sum(weights[0..i])
    total     float64     // 总权重
}

// 预处理：O(n) 构建前缀和
func BuildPrefixSum(items []string, weights []float64) *PrefixSumSampler {
    n := len(items)
    prefixSum := make([]float64, n)
    prefixSum[0] = weights[0]
    for i := 1; i < n; i++ {
        prefixSum[i] = prefixSum[i-1] + weights[i]
    }
    return &PrefixSumSampler{
        items:     items,
        prefixSum: prefixSum,
        total:     prefixSum[n-1],
    }
}

// 采样：O(log n) 二分查找
func (s *PrefixSumSampler) Draw(rng *Xoshiro256ss) string {
    r := rng.Float64() * s.total
    // 找第一个 prefixSum[i] > r 的位置
    idx := sort.Search(len(s.prefixSum), func(i int) bool {
        return s.prefixSum[i] > r
    })
    return s.items[idx]
}

// ============ 模拟抽奖 ============

func main() {
    // 模拟过滤后的概率组
    items := []string{"SSR龙王", "SR剑姬", "R铁剑", "N金币"}
    weights := []float64{6, 51, 300, 643}

    sampler := BuildPrefixSum(items, weights) // O(n) 预处理
    rng := NewRNG()

    stats := map[string]int{}
    for i := 0; i < 100000; i++ {
        stats[sampler.Draw(rng)]++ // 每次 O(log n)
    }

    fmt.Println("=== 前缀和+二分查找 10万次统计 ===")
    for _, item := range items {
        fmt.Printf("%-10s: %5d 次 (%.2f%%)\n", item, stats[item], float64(stats[item])/1000)
    }
    // 输出示例：
    // SSR龙王   :   612 次 (0.61%)
    // SR剑姬    :  5089 次 (5.09%)
    // R铁剑     : 30021 次 (30.02%)
    // N金币     : 64278 次 (64.28%)
}
```

---

#### 7.3.5 算法三：Alias Method（别名法）

**原理**：将 n 个不等概率的桶重新分配为 n 个等概率的桶，每个桶最多装两个道具。采样时先均匀选桶，再按阈值决定取哪个道具。

```
原始权重 → 归一化为"每桶面积=1" → 大桶填补小桶

桶0 [SSR龙王 0.024 | N金币 0.976]   ← 小概率道具+别名
桶1 [SR剑姬  0.204 | N金币 0.796]   ← 小概率道具+别名
桶2 [R铁剑   1.0              ]     ← 刚好满桶
桶3 [N金币   0.772 | R铁剑 0.228]   ← 大概率道具也可能被切分

采样：随机数1选桶(如桶1) → 随机数2=0.15 < 0.204 → 返回"SR剑姬"
```

**优点**：
- 采样 O(1)，恒定时间，与道具数量无关
- 概率精确，无浮点累积误差

**缺点**：
- 预处理逻辑较复杂（small/large 双队列）
- 不支持动态增删，过滤后需完整重建
- 浮点精度问题可能导致 small/large 队列不平衡，需要兜底处理

**完整可运行 Demo**：

```go
package main

import (
    crand "crypto/rand"
    "encoding/binary"
    "fmt"
    "math/bits"
)

// ============ 随机数生成器（同 7.5 节） ============

type Xoshiro256ss struct{ s [4]uint64 }

func NewRNG() *Xoshiro256ss {
    var buf [32]byte
    crand.Read(buf[:])
    rng := &Xoshiro256ss{}
    for i := 0; i < 4; i++ {
        rng.s[i] = binary.LittleEndian.Uint64(buf[i*8:])
    }
    return rng
}

func (x *Xoshiro256ss) Uint64() uint64 {
    result := bits.RotateLeft64(x.s[1]*5, 7) * 9
    t := x.s[1] << 17
    x.s[2] ^= x.s[0]; x.s[3] ^= x.s[1]
    x.s[1] ^= x.s[2]; x.s[0] ^= x.s[3]
    x.s[2] ^= t; x.s[3] = bits.RotateLeft64(x.s[3], 45)
    return result
}

func (x *Xoshiro256ss) Intn(n int) int {
    hi, _ := bits.Mul64(x.Uint64(), uint64(n))
    return int(hi)
}

func (x *Xoshiro256ss) Float64() float64 {
    return float64(x.Uint64()>>11) / (1 << 53)
}

// ============ Alias Method ============

type AliasTable struct {
    n     int
    prob  []float64 // 每个桶的保留概率
    alias []int     // 别名跳转目标
    items []string  // 道具ID列表
}

// 预处理：O(n)，构建别名表
func BuildAlias(items []string, weights []int) *AliasTable {
    n := len(items)
    total := 0
    for _, w := range weights { total += w }

    t := &AliasTable{n: n, prob: make([]float64, n), alias: make([]int, n), items: items}
    scaled := make([]float64, n)
    small, large := []int{}, []int{}

    for i, w := range weights {
        scaled[i] = float64(w) * float64(n) / float64(total)
        if scaled[i] < 1.0 { small = append(small, i) } else { large = append(large, i) }
    }
    for len(small) > 0 && len(large) > 0 {
        s, l := small[len(small)-1], large[len(large)-1]
        small, large = small[:len(small)-1], large[:len(large)-1]
        t.prob[s] = scaled[s]
        t.alias[s] = l
        scaled[l] -= (1.0 - scaled[s])
        if scaled[l] < 1.0 { small = append(small, l) } else { large = append(large, l) }
    }
    for _, i := range large { t.prob[i] = 1.0 }
    for _, i := range small { t.prob[i] = 1.0 } // 浮点兜底
    return t
}

// 采样：O(1)
func (t *AliasTable) Draw(rng *Xoshiro256ss) string {
    i := rng.Intn(t.n)            // 随机数1：均匀选桶
    if rng.Float64() < t.prob[i] { // 随机数2：和桶阈值比较
        return t.items[i]          //   小于阈值 → 本桶道具
    }
    return t.items[t.alias[i]]     //   大于阈值 → 跳到别名道具
}

// ============ 模拟抽奖 ============

func main() {
    // 模拟过滤后的概率组
    items := []string{"SSR龙王", "SR剑姬", "R铁剑", "N金币"}
    weights := []int{6, 51, 300, 643}

    table := BuildAlias(items, weights) // O(n) 预处理
    rng := NewRNG()

    stats := map[string]int{}
    for i := 0; i < 100000; i++ {
        stats[table.Draw(rng)]++ // 每次 O(1)
    }

    fmt.Println("=== Alias Method 10万次统计 ===")
    for _, item := range items {
        fmt.Printf("%-10s: %5d 次 (%.2f%%)\n", item, stats[item], float64(stats[item])/1000)
    }
    // 输出示例：
    // SSR龙王   :   612 次 (0.61%)
    // SR剑姬    :  5089 次 (5.09%)
    // R铁剑     : 30021 次 (30.02%)
    // N金币     : 64278 次 (64.28%)
}
```

---

#### 7.3.6 选型结论：游戏商业化场景推荐轮盘赌

**在"过滤链前置 → 最终得到纯粹概率组 → 采样"的架构下**：

| 维度 | 轮盘赌 | 前缀和+二分 | Alias Method |
|------|--------|------------|--------------|
| 实现复杂度 | 10 行核心代码 | 20 行 | 40+ 行 |
| 可读性/可维护性 | 新人一眼看懂 | 需理解二分语义 | 需理解别名构建逻辑 |
| 动态过滤适配 | 拿到过滤后列表直接抽 | 需先构建前缀和 | 需完整重建别名表 |
| 连抽中道具集变化 | 无影响（每次独立扫描） | 需重建前缀和 | 需重建别名表 |
| 调试/审计 | 流程线性，容易复现 | 中等 | 桶+别名映射不直观 |
| Fixed/Unfixed 支持 | 天然支持（扫描时分别处理） | 需额外处理 | 需额外处理 |
| 性能（n=500） | ~0.5μs | ~0.1μs | ~0.02μs |

**最终选择：轮盘赌（线性扫描）**

理由：
1. **过滤后的概率组是 per-request 动态生成的**，Alias Method 每次都要 O(n) 重建，而轮盘赌省掉了这个预处理步骤——总开销反而更小
2. **连抽场景**中每抽中一个新道具就要从概率组移除，轮盘赌只需 `continue` 跳过或从 slice 中删除，Alias 和前缀和都需要重建
3. **实现简洁**，符合"代码即文档"的原则，降低后续维护和审计的心智负担
4. **性能差异在微秒级**（0.5μs vs 0.02μs），相比整个请求链路中的 Redis IO（~1ms）和过滤链逻辑（~0.5ms），采样算法本身的性能差异对端到端延迟没有可观测影响

```
实际请求耗时分布（P99=50ms）：
├─ Redis 查用户背包/库存    ~3ms   ██████████████████████████ 60%
├─ 过滤链执行              ~1ms   █████████ 20%
├─ 限量扣减(Redis Lua)     ~1ms   █████████ 15%
├─ 采样算法                ~0.5μs  ▏ 0.001% ← 三种算法的差异在这里，完全无意义
└─ 其他(序列化/日志等)      ~2ms   █████ 5%
```

**Alias Method 适用场景**（本系统中不适用但作为参考）：
- 概率组固定不变、无动态过滤（如稀有度决定层：SSR 1% / SR 10% / R 89%）
- 同一张表被采样数十万次才变更一次
- 对单次采样延迟极度敏感（如实时物理引擎中的粒子生成）

## 8. 里程碑核心机制

```
每次抽奖后 → LuckyValuePlugin.UpdateLuckyValue 更新幸运值
每次抽奖前 → ProbIdGenPlugin.GenProbId → LuckyValuePlugin.Match 匹配里程碑
命中里程碑 → 切换到高稀有度概率组 → 保证产出大奖
抽中大奖后 → 清零对应优先级组的幸运值（重新累积）
```

### 8.1 里程碑配置

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

### 8.2 支持的保底玩法

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

**分区数设计基准：**
- RocketMQ 单分区单消费者线程顺序消费，单分区安全吞吐约 **5000 条/s**（8核16G Broker，同步刷盘约 3000/s，异步刷盘约 8000/s，保守取中间值）
- 所需分区数 = 峰值消息速率 ÷ (单分区吞吐 × 冗余系数0.7)，再取 2 的幂次
- 冗余系数 0.7：分区实际负载不超过 70%，预留消息堆积消化、消费者扩容窗口

| Topic | 峰值消息速率 | 分区数计算 | 分区数 | 刷盘策略 | 用途 | 消费者 |
|-------|------------|-----------|--------|---------|------|--------|
| `LOTTERY_DELIVERY` | 50万条/s | 50万÷(5000×0.7)≈143, 取2的幂次 | **64** | 同步刷盘 | 抽奖发货（道具不能丢） | 发货服务 |
| `LOTTERY_RECORD` | 50万条/s | 50万÷(8000×0.7)≈89, 取2的幂次 | **32** | 异步刷盘 | 抽奖流水异步落库 | 记录服务 |
| `LOTTERY_BROADCAST` | 1万条/s | 远低于单分区吞吐, 保障高可用取 | **8** | 异步刷盘 | 大奖全服广播 | 通知服务 |
| `LOTTERY_STAT` | 50万条/s | 允许堆积消费, 降配 | **16** | 异步刷盘 | 概率统计/产出监控 | 统计服务 |
| `LOTTERY_DLQ` | 极低 | 保障多节点消费 | **4** | 同步刷盘 | 死信兜底 | 告警+人工 |

> **为何 LOTTERY_DELIVERY 用同步刷盘？** 道具=用户资产，消息丢失=用户花了钱没拿到东西（P0事故），牺牲20%写性能换资产安全。
> **为何 LOTTERY_STAT 分区数少于理论值？** 统计允许秒级延迟，可堆积消费；紧急时暂停统计消费，优先保障发货 Topic 消费算力。

### 13.2 消息可靠性

**生产者端（事务消息 + 回查）：**
```
① 抽奖服务发送半消息 → MQ 持久化 ACK
② 执行本地事务：MySQL（扣货币 + 写 draw_record + 写 draw_transaction）
   ├─ 成功 → Commit 半消息（消费者可见）
   └─ 失败 → Rollback + Redis INCRBY 回补库存
③ 事务回查（Producer宕机/超时时MQ主动回查）：
   ├─ 查 draw_transaction.status=1 → Commit
   ├─ 不存在 → Rollback
   └─ status=0 且超过5分钟 → 补偿重试（最多3次后标记失败告警）
```

**消费者端（幂等消费 + 发货路由）：**
```go
func consumeDelivery(msg *DeliveryMsg) error {
    // Step1: 幂等检查（delivery_record 唯一索引 uk_draw_item）
    exists, err := checkDeliveryRecord(msg.DrawID, msg.ItemIndex)
    if err != nil { return err }  // 查询失败 → 触发重试
    if exists { return nil }      // 已处理 → 直接 ACK

    // Step2: 路由发货渠道
    channel := routeChannel(msg.ItemType)
    // character/skin → game_server_rpc
    // currency/coin  → wallet_service
    // coupon         → coupon_center
    // physical       → logistics_service（需收货地址）

    // Step3: 调用渠道方发货API（超时5s）
    result, err := channel.Deliver(msg.UserID, msg.ItemID, msg.DeliveryID)
    if err != nil {
        return err  // 失败 → RocketMQ 自动重试（指数退避，最多16次）
    }

    // Step4: 写发货记录（uk_draw_item 唯一索引兜底防重）
    if err := insertDeliveryRecord(msg, result); err != nil {
        if isDuplicateKeyErr(err) { return nil }  // 并发重复，忽略
        return err
    }

    // Step5: 后续通知（异步，不影响ACK）
    go func() {
        pushArrivalNotification(msg.UserID, msg.ItemID)  // 到账推送
        if isRareItem(msg.ItemID) {
            publishBroadcast(msg)  // 大奖全服广播
        }
    }()
    return nil  // ACK
}
```

**消息堆积处理：**
- `LOTTERY_DELIVERY` 堆积 > 1万 → P1 告警；> 5万 → P0 告警
- 紧急处理：动态扩容消费者线程池（16 → 128 线程），开启批量消费（每批 50 条）
- 优先级：暂停 `LOTTERY_STAT` 消费（统计允许延迟），优先保障 `LOTTERY_DELIVERY`
- 兜底：堆积超30分钟未消费的消息 → 转入 `LOTTERY_DLQ` 死信 + 自动触发补偿流程

---

## 14. 缓存架构与一致性

### 14.1 多级缓存设计

```
L1 本地缓存（进程内 go-cache，命中率目标 90%）:
   ├── pool_config_{pool_id}_{version}    奖池配置，TTL=5min
   ├── alias_table_{prob_id}_{version}    Alias表，TTL=5min
   ├── item_config_{item_id}              道具配置，TTL=10min
   └── 特点：读多写少数据，版本号变更时主动失效

L2 Redis Cluster（32分片，命中率目标 99%+）:
   ├── state:{game_id}:{user_key}[_{layer}]  用户实时状态（权威源）
   ├── stock:{pool_id}:{item_id}              限量库存计数
   ├── limit:{item_id}:{scope}                限量累计计数
   ├── rsp_cache:{user_key}[_{layer}]         响应缓存
   ├── lock:{user_key}                        分布式锁
   └── alias_table:{pool_id}:{prob_id}:{ver}  概率Alias表缓存

L3 MySQL（最终持久化）:
   └── 抽奖记录、发货记录、配置表 —— 核心链路不直连
```

### 14.2 Redis Key 设计与内存估算

| Key 模式 | 类型 | 单条大小 | 数量 | 总内存 | TTL |
|---------|------|---------|------|--------|-----|
| `state:{game}:{user}[_{layer}]` | String(JSON/zlib) | ~2KB | 2000万 | **40GB** | 永久(常驻)/活动后30天 |
| `rsp_cache:{user}[_{layer}]` | String(JSON) | ~1KB | 2000万 | **20GB** | 常驻=10年/普通=1年 |
| `stock:{pool}:{item}` | String(int) | 64B | 1万 | **640KB** | 永久 |
| `limit:{item}:{scope}` | String(int) | 64B | 5万 | **3.2MB** | date:86400s/global:永久 |
| `lock:{user}` | String | 50B | ~50万(并发在线) | **25MB** | 5s |
| `alias:{pool}:{prob}:{ver}` | String(Protobuf) | ~50KB | 500 | **25MB** | 3600s |
| `dbl_state:{game}:{user}` | String(JSON) | ~4KB | 2000万 | 仅非PlayBack模式 | 同state |
| **合计** | — | — | — | **~60GB** | — |

> 60GB ÷ 32分片 = 每分片约 2GB，充裕。

### 14.3 缓存一致性策略

| 数据类型 | 一致性策略 | 失效方式 |
|---------|-----------|---------|
| 用户状态 | Redis 为权威源，MySQL 异步持久化 | 不存在"一致性"问题，Redis即真相 |
| 奖池配置 | MySQL 写 → ETCD Watch 通知 → 本地版本号校验失效 | 版本号+1 触发全节点重加载 |
| Alias 表 | 概率组变更时重建写 Redis，本地按版本号失效 | `alias:{pool}:{prob}:{ver}` 新版本覆盖旧版本 |
| 限量库存 | Redis 原子计数为权威，MySQL 定时对账修正 | 不依赖失效，靠对账保持一致 |
| 响应缓存 | stateKey(drawType_drawRound) 自然失效 | drawRound 变化 → 旧缓存 stateKey 不匹配 → 自动跳过 |

**配置变更通知流程**：
```
运营后台修改奖池配置 → MySQL 写入新版本(version+1) → 
写入 ETCD Key "/lottery/config/{pool_id}" = {new_version} →
各服务节点 Watch ETCD 前缀 "/lottery/config/" →
收到变更事件后清除本地 go-cache 对应 Key →
下次请求时重新从 Redis 加载最新版本配置
延迟：< 500ms（ETCD Watch 长轮询实时推送，比 Redis Pub/Sub 更可靠——不丢消息）
```

### 14.4 缓存穿透/击穿/雪崩防护

| 问题 | 场景 | 解决方案 |
|------|------|---------|
| **穿透** | 查询不存在的 pool_id/user_id | 空值缓存 TTL=60s + 布隆过滤器（活跃奖池集合）|
| **击穿** | 热点奖池 Alias 表过期瞬间大量并发重建 | singleflight 合并：同一 Key 只有1个请求穿透到 Redis 重建 |
| **雪崩** | 大量 Key 同时过期 | TTL 加随机偏移（±30s）；奖池配置不设TTL，靠版本号主动失效 |
| **热点Key** | 限量道具 stock Key 被百万QPS打 | 分桶打散：stock:{pool}:{item}:{shard_id}，按user_id%N路由 |

### 14.5 Redis 宕机恢复

```
Redis 主从切换（<30s）：
  ├─ 哨兵自动提升从为主，大多数情况数据不丢
  └─ 切换期间该分片请求失败 → 客户端重试即可

Redis 数据丢失（极端场景）：
  ├─ 用户状态：从 t_user_state 表加载（可能丢最后1-2次抽奖进度）
  ├─ 限量库存：从 delivery_record COUNT 重建：SET stock = total - COUNT(success)
  ├─ 响应缓存：丢失无害，下次请求走正常流程即可
  └─ 原则："宁可少扣保底不可多扣"——恢复后用户保底可能多给1-2次，可接受
```

---

## 15. 容错性设计

### 15.1 限流（分层精细化）

| 层次 | 维度 | 算法 | 阈值 | 动作 |
|------|------|------|------|------|
| 网关全局 | 总流量 | 令牌桶 | 200万 QPS | 返回 503 |
| 游戏维度 | 单游戏 | 令牌桶 | 50万 QPS（可配） | 排队等待 |
| 用户维度 | 单 uid | 滑动窗口 | 10次/s | 返回频率限制 |
| 奖池维度 | 单奖池 | 令牌桶 | 配置上限 | 排队等待 |
| IP 维度 | 单 IP | 滑动窗口 | 50次/s | 超出拉黑10min |

**限流 Redis Lua 实现（滑动窗口）**：
```lua
-- KEYS[1] = rate_limit:{uid}
-- ARGV[1] = 窗口大小(ms)  ARGV[2] = 最大请求数  ARGV[3] = 当前时间(ms)
local key = KEYS[1]
local window = tonumber(ARGV[1])
local max_requests = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
-- 移除窗口外的旧请求
redis.call('ZREMRANGEBYSCORE', key, 0, now - window)
-- 计算当前窗口内请求数
local count = redis.call('ZCARD', key)
if count >= max_requests then
    return 0  -- 限流
end
-- 记录本次请求
redis.call('ZADD', key, now, now .. math.random())
redis.call('PEXPIRE', key, window)
return 1  -- 放行
```

### 15.2 熔断与降级

```
异常触发条件（满足任一项自动触发）：
├─ Redis P99 > 30ms（正常 < 5ms）
├─ MQ 发货消息堆积 > 5万条
├─ 核心接口错误率 > 0.5%（1分钟窗口）
├─ 概率偏离 > 30%（实际 vs 配置）
└─ 上游游戏服务器超时率 > 10%
        ↓
熔断行为：
├─ 半开探测：熔断 10s 后放行 5% 流量试探
├─ 快速失败：返回"系统繁忙，请稍后再试"
└─ 恢复判定：连续 10 次探测成功 + 错误率 < 0.1%
        ↓
分级降级（lottery.degrade_level: 0~4）：
├─ Level 0: 正常运行
├─ Level 1: 关闭全服广播、统计上报、概率监控（非核心旁路）
├─ Level 2: 禁止新池开启，禁止连抽（仅允许单抽），降低限流阈值
├─ Level 3: 发货切同步RPC（RT升高但不丢），禁止高消耗操作
└─ Level 4: Redis全挂 → 切MySQL乐观锁模式
```

**Level 4 DB 降级模式详细 SQL**：
```sql
-- 降级模式：MySQL CAS 乐观锁替代 Redis 原子操作
-- 性能：从50万QPS降至5万TPS，短期可控

-- 1. 读取用户状态（SELECT FOR UPDATE 加行锁）
SELECT state_json, version FROM t_user_state
WHERE user_id=? AND game_id=? AND pool_type=? AND layer=?
FOR UPDATE;

-- 2. 应用层执行抽奖逻辑（内存中处理）

-- 3. CAS 更新状态（version 乐观锁）
UPDATE t_user_state SET state_json=?, version=version+1
WHERE user_id=? AND game_id=? AND pool_type=? AND layer=? AND version=?;
-- 受影响行数=0 → 被其他请求抢占，重试

-- 4. 限量库存降级扣减
UPDATE t_inventory SET used_stock=used_stock+1, version=version+1
WHERE pool_id=? AND item_id=? AND used_stock < total_stock AND version=?;
-- 受影响行数=0 → 库存不足或并发冲突，重试/拒绝
```

### 15.3 兜底方案矩阵

| 故障场景 | 兜底策略 | 用户感知 | 恢复方式 |
|---------|---------|---------|---------|
| Redis 单分片宕机 | 哨兵自动切主从（<30s） | 30s内偶发失败 | 自动恢复 |
| Redis 集群全挂 | 切 MySQL FOR UPDATE 模式 | RT 升至 200ms | 手动恢复 |
| RocketMQ 宕机 | 发货切同步 RPC + 本地 WAL | 发货稍慢 | 手动恢复 |
| 游戏服务器不可达 | 消息积压MQ，恢复后自动消费 | 延迟到账 | 自动恢复 |
| 限量库存 Redis 丢失 | 从 delivery_record COUNT 重建 | 无感知 | 自动恢复 |
| 概率异常/Bug | 一键下线奖池 + 版本回滚 | 活动暂停 | 人工处理 |
| 奖池全部道具被过滤 | complement 兜底概率组 | 获得通用奖励 | 自动 |

### 15.4 动态配置开关（ETCD，秒级生效）

```yaml
lottery.switch.global: true           # 全局抽奖开关
lottery.switch.delivery_async: true   # 发货异步开关（故障时切同步）
lottery.switch.db_fallback: false     # DB 降级模式开关
lottery.switch.broadcast: true        # 全服广播开关
lottery.switch.stat_report: true      # 统计上报开关
lottery.switch.new_pool: true         # 新池开启开关
lottery.limit.draw_qps: 500000        # 抽奖总 QPS 上限
lottery.limit.user_rate: 10           # 单用户每秒上限
lottery.degrade_level: 0              # 降级级别 0~4
lottery.hotspot.shard_count: 1        # 热点道具分桶数（默认1不分桶）
lottery.delivery.timeout_ms: 5000     # 发货超时时间
lottery.delivery.max_retry: 16        # 最大重试次数
lottery.reconcile.interval_sec: 300   # 对账间隔
```

### 15.5 水平扩展方案

**服务层扩展（无状态，K8s 管理）**：
```
所有服务无状态 → K8s Deployment 管理
├─ HPA 策略：CPU > 60% 自动扩容，CPU < 30% 自动缩容
├─ 大型活动预扩容：提前3天按3倍峰值扩容
└─ 节点故障：健康检查失败自动摘流，新 Pod 30s 内就绪
```

**Redis 在线扩容（32分片 → 64分片）**：
```
① 新增32个分片节点并入集群
② redis-cli --cluster reshard 在线迁移 slot
③ 扩容期间：
   ├─ 用户状态 Key 随 slot 迁移，无需应用感知
   ├─ 限量库存 Key 需暂停写入 → 迁移完成后恢复（秒级）
   └─ 临时暂停热点道具抽奖，规避迁移中数据不完整
④ 验证：全量 Key 扫描确认分布均匀
```

**DB 分库扩容（16库 → 32库，双写迁移）**：
```
① 新建目标集群：32库256表
② 开启双写：业务写入同时落旧16库 + 新32库
③ 后台异步迁移历史数据（按 user_id 重算分片路由）
④ 全量校验：新旧集群数据一致
⑤ 路由切换：指向32库新集群，关闭双写
⑥ 旧集群只读保留7天后下线
```

**冷热数据分层**：
```
热数据（0~7天）  : Redis + MySQL（在线实时查询）
温数据（7~90天） : MySQL 归档库（读写分离，按需查询）
冷数据（>90天）  : 导入 HDFS / TiDB（海量历史离线分析）
迁移策略：定时任务每日凌晨将超过7天的 draw_record 搬运到归档库
```

---

## 16. 监控与合规

### 16.1 监控指标体系

**业务指标（Prometheus + Grafana）**：
```yaml
business_metrics:
  - name: draw_qps
    type: counter
    labels: [game_id, pool_id, draw_type]
    description: "抽奖QPS"
    
  - name: draw_latency_ms
    type: histogram
    buckets: [5, 10, 25, 50, 100, 250, 500]
    labels: [game_id, pool_id]
    description: "抽奖接口延迟（目标P99<50ms）"
    
  - name: delivery_success_rate
    type: gauge
    labels: [channel, item_type]
    description: "发货成功率（目标>99.9%）"
    
  - name: item_output_rate
    type: gauge
    labels: [pool_id, item_id, rarity]
    description: "各道具实际产出率（对比配置概率）"
    
  - name: pity_trigger_count
    type: counter
    labels: [pool_id, pity_type]
    description: "保底触发次数"
    
  - name: stock_remaining
    type: gauge
    labels: [pool_id, item_id]
    description: "限量道具剩余库存"
    
  - name: revenue_amount
    type: counter
    labels: [game_id, pool_id, cost_type]
    description: "实时收入（分）"
```

### 16.2 告警分级与响应 SLA

| 级别 | 条件 | 响应SLA | 通知方式 | 示例 |
|------|------|---------|---------|------|
| **P0** | 资产损失/核心不可用 | 5分钟内响应 | 电话+短信+企微群 | 超发/概率偏离>30%/发货堆积>5万/错误率>0.5% |
| **P1** | 功能受损/性能劣化 | 15分钟内响应 | 短信+企微群 | 库存<5%/P99>100ms/单分片宕机 |
| **P2** | 预警/非紧急异常 | 1小时内响应 | 企微群 | 对账差异<5/消费延迟>10s |
| **INFO** | 运营关注 | 下一工作日 | 日报邮件 | 新池上线通知/大奖产出/配置变更 |

### 16.3 概率监控大盘

```
┌────────────────────────────────────────────────────────────────────┐
│                      概率监控大盘 (Grafana)                          │
├────────────────────────────────────────────────────────────────────┤
│                                                                    │
│  Panel 1: 实时产出率 vs 配置概率（时序图）                            │
│  ├─ 曲线1: actual_ssr_rate（蓝色）                                  │
│  ├─ 曲线2: configured_ssr_rate（绿色虚线）                          │
│  ├─ 告警带: ±20% 偏差区域（红色阴影）                                │
│  └─ 触发: 曲线1 超出告警带持续5分钟 → P0                            │
│                                                                    │
│  Panel 2: 保底触发分布（直方图）                                     │
│  ├─ X轴: 触发时的抽奖次数（1~90）                                   │
│  ├─ Y轴: 触发次数                                                  │
│  └─ 预期: 几何分布 Geo(0.006)，偏离则概率配置有误                    │
│                                                                    │
│  Panel 3: 用户欧气分布（百分位图）                                   │
│  ├─ P50/P90/P99 用户的SSR率                                        │
│  └─ 异常: 有用户SSR率 > 3σ → 标记可疑                               │
│                                                                    │
│  Panel 4: 限量道具库存消耗进度（仪表盘）                             │
│  ├─ 各限量道具：已发放/总量 百分比                                   │
│  └─ 触发: 剩余 < 5% → P1 告警                                      │
│                                                                    │
│  Panel 5: 异常用户 TOP10（表格）                                     │
│  ├─ 列: user_id, draw_count, ssr_count, ssr_rate, device_risk     │
│  └─ 筛选: ssr_rate > 3 × expected_rate                             │
│                                                                    │
└────────────────────────────────────────────────────────────────────┘
```

### 16.4 全链路追踪

```
draw_id（雪花ID）贯穿全流程，可通过 draw_id 追踪：

抽奖请求 → draw_id 生成
    │
    ├─ draw_record 表：draw_id → 完整抽奖快照
    ├─ draw_transaction 表：draw_id → 事务状态
    ├─ RocketMQ 消息 Keys：draw_id → 消息轨迹
    ├─ delivery_record 表：draw_id → 发货状态
    └─ ES 日志：draw_id → 全链路日志（请求/响应/中间状态）

排查工具：
  输入 draw_id → 一键查询全链路状态（抽奖结果/事务状态/MQ投递/发货结果）
  输入 user_id + time_range → 查询该用户所有抽奖记录
```

### 16.5 日级对账

```
每日凌晨3点执行全量对账：

① 库存对账:
   FOR EACH item IN limited_items:
     expected = total_stock - COUNT(delivery WHERE item_id=item AND status=SUCCESS)
     actual = Redis GET stock:{pool}:{item}
     IF actual != expected:
       INSERT t_reconciliation_task (diff=actual-expected)
       IF abs(diff) > 5: ALERT P1
       SET stock:{pool}:{item} = expected  -- 以MySQL为准修正

② 发货对账:
   undelivered = SELECT * FROM draw_record WHERE draw_id NOT IN 
                 (SELECT draw_id FROM delivery_record WHERE status IN (SUCCESS, COMPENSATED))
                 AND draw_time < NOW() - INTERVAL 30 MINUTE
   FOR EACH record IN undelivered:
     自动触发补发 OR 创建人工工单

③ 金额对账（金融等式）:
   总消耗 = SUM(draw_record.cost_amount)
   总产出价值 = SUM(道具配置价值 × 产出数量)
   产出率 = 总产出价值 / 总消耗  -- 应接近配置的期望产出率
   偏离 > 10% → P1 告警（可能概率配置错误）
```

### 16.6 合规要求

```yaml
compliance:
  # 概率公示（国家新闻出版署要求）
  probability_disclosure:
    public_api: "/api/v1/pool/{pool_id}/probability"
    precision: 4                    # 小数点后4位
    update_on_change: true          # 概率变化时更新公示
    display_all_items: true         # 展示所有可获得道具及其概率
    
  # 未成年人保护（防沉迷+消费限制）
  minor_protection:
    age_0_8: "ban"                  # 完全禁止抽奖
    age_8_16: "limit_200_monthly"   # 月消费限额200元
    age_16_18: "limit_400_monthly"  # 月消费限额400元
    require_real_name: true         # 必须实名认证
    
  # 审计日志（不可篡改）
  audit_log:
    retention_days: 365             # 保留1年
    immutable: true                 # 写入后不可修改
    storage: "append_only_log"      # 追加写入，不允许UPDATE/DELETE
    fields: [operator, action, target, before_state, after_state, timestamp, ip]
    
  # 消费确认
  purchase_confirmation:
    threshold: 648                  # 单次消费超过648元需二次确认
    daily_alert: 2000               # 日消费超2000元发送提醒
```

---

## 17. 面试必考专题

### 专题一：高并发流量削峰

**Q: 新池开启瞬间百万QPS涌入，系统如何扛住？**

**答**：四层削峰架构，层层拦截递减流量：

```
用户请求 100万QPS
    │
    ├─ ① 网关层限流（令牌桶 + 滑动窗口）
    │   ├─ 全局硬限：200万 QPS（超出直接 503）
    │   ├─ 用户维度：单 uid 10次/s（滑动窗口）
    │   ├─ IP 维度：单 IP 50次/s
    │   └─ 奖池维度：单池可配上限
    │   拦截效果：约20%无效/超频请求被拦
    │
    ├─ ② 响应缓存层（本地 + Redis）
    │   ├─ 相同 drawType+drawRound 重复请求 → 命中缓存直接返回
    │   ├─ 缓存检查在分布式锁之前，无需任何竞争
    │   └─ 拦截效果：80%重试/重复请求被拦
    │   穿透到业务逻辑：约 50万 QPS
    │
    ├─ ③ 抽奖核心（纯内存计算 + Redis）
    │   ├─ Alias Method O(1) 采样，不随道具数增长
    │   ├─ 用户状态全量在 Redis，不查 MySQL
    │   ├─ 本地缓存：奖池配置/Alias表/道具表，命中率>90%
    │   └─ 单次抽奖：1次Redis Pipeline（读状态+扣库存+写状态）
    │   RT: P99 < 50ms
    │
    └─ ④ 异步发货（MQ 削峰）
        ├─ 抽奖结果确定后立即返回客户端，不等发货
        ├─ RocketMQ 事务消息投递，Worker 按 DB 承受能力消费
        ├─ 50万/s 抽奖产生的发货消息，Worker 以 5万 TPS 持续消费
        └─ MQ 充当"蓄水池"，将秒级洪峰平摊为分钟级稳定写入
```

**预热策略**：
- **活动前3天**：按3倍峰值预扩容服务节点（K8s Deployment 调 replicas）
- **活动前1天**：全链路压测，验证扩容后水位
- **Alias 表预构建**：奖池配置发布时即构建所有概率组的 Alias 表写入 Redis，不等首次请求触发
- **Redis 预热**：热门奖池的用户状态提前加载（18201 初始化批量触发）
- **本地缓存预加载**：服务启动时主动拉取所有活跃奖池配置到 go-cache

**限流算法选型对比**：

| 算法 | 适用层 | 优势 | 本系统使用场景 |
|------|--------|------|--------------|
| 令牌桶 | 网关全局 | 允许突发流量，平滑限制 | 全局 QPS 上限 |
| 滑动窗口 | 用户维度 | 精确计数，无突刺 | 单用户 10次/s |
| 漏桶 | 发货消费 | 固定速率输出，天然削峰 | Worker 消费速率控制 |

---

### 专题二：热点 & 超卖防护

**Q: 限量道具（如全服仅10个）如何在50万QPS下保证绝不超发？**

**答**：三层防超发闭环（参考红包防超抢设计）：

```
  【限量道具防超发 三层防护体系】
        ↓
① 第一层：Redis Lua 原子扣减（实时拦截，承载10万+ QPS/分片）
        ├─ Lua 脚本保证"判断+扣减"原子执行，无竞态窗口
        ├─ 返回 ≥ 0 → 扣减成功
        ├─ 返回 -1 → 库存不足，立即拒绝（不进入后续流程）
        ├─ DB 事务失败 → INCRBY 立即回补
        └─ 天然防超发：原子操作 + 失败回补 = 闭环
        ↓
② 第二层：发货 Worker 二次校验 MySQL（消费端兜底）
        ├─ SELECT used_stock FROM t_inventory WHERE item_id=? FOR UPDATE
        ├─ 若 used_stock >= total_stock → 拒绝发货 + 回补 Redis + 告警
        └─ 防的是：Redis 主从切换丢数据导致多弹出的极端场景
        ↓
③ 第三层：定时对账修正（最终兜底）
        ├─ 每5分钟：redis_stock vs (total - COUNT(delivery WHERE success))
        ├─ 差异修正：以 MySQL 为准纠正 Redis（SET stock = expected_remain）
        ├─ 差异 > 5 → P1 告警
        └─ 每日凌晨全量对账
```

**Redis Lua 原子扣减脚本**：
```lua
-- KEYS[1] = stock:{pool_id}:{item_id}
-- ARGV[1] = 扣减数量（通常为1）
local stock = redis.call('GET', KEYS[1])
if stock == false then return -1 end
local remain = tonumber(stock) - tonumber(ARGV[1])
if remain < 0 then return -1 end          -- 库存不足，立即拒绝
redis.call('DECRBY', KEYS[1], ARGV[1])
return remain                              -- 返回剩余库存
```

**热点道具分桶**（极端场景：全服1000万人抢同一个限量道具）：
```
正常场景：单 Key stock:{pool}:{item} 单分片承载
极端热点：按 user_id % N 路由到 N 个分桶
  ├─ stock:{pool}:{item}:0 → 库存 2
  ├─ stock:{pool}:{item}:1 → 库存 2
  ├─ ...
  └─ stock:{pool}:{item}:4 → 库存 2（共10个，分5桶）
  
  本桶为空 → 向相邻桶"借"（LPOP/尝试扣减），最多轮询3次
```

---

### 专题三：一致性 & 幂等

**Q: 如何保证"扣了钱一定发货"且"不会多扣/多发"？**

**答**：RocketMQ 事务消息 + 三级幂等 + 唯一索引兜底：

**一致性保证（不丢）**：
```
① 发送半消息 → MQ 返回 ACK（消息已持久化）
② 执行本地事务：MySQL 事务（扣货币 + 写 draw_record + 写 draw_transaction）
   ├─ 成功 → Commit 半消息（消费者可见）
   └─ 失败 → Rollback 半消息 + Redis INCRBY 回补
③ Producer 宕机 → RocketMQ 定时回查
   ├─ 查 draw_transaction.status=1 → Commit
   ├─ 查 draw_transaction 不存在 → Rollback
   └─ 不确定 → UNKNOW（继续回查，最多15次，间隔60s）
```

**幂等保证（不重）**：

| 层级 | 机制 | 防护场景 | 实现 |
|------|------|---------|------|
| L1 响应缓存 | stateKey(drawType_drawRound) 匹配 | 客户端网络超时重试 | 命中直接返回，不执行任何逻辑 |
| L2 状态回退/回放 | 双缓存 PrevState / PlayBackRewardInfo | 客户端动画重播 | 回退状态重抽 or 按序回放历史 |
| L3 分布式锁 | Redis SET NX (TTL=5s) | 同用户并发请求 | 同时只有1个请求执行 |
| L4 唯一索引 | delivery_record.uk_draw_item | MQ 重复投递 | INSERT 冲突则 ACK |

**消息不乱（顺序性）**：
- 同一用户的发货消息按 user_id 哈希到固定分区 → 保证单用户串行消费
- 连抽内多个道具按 item_index 排序 → 发货顺序可预测

**唯一抽奖凭证设计**：
```
draw_id = 雪花ID（全局唯一，含时间戳+机器ID+序列号）
batch_id = draw_id 的前缀部分（同一次连抽共享）
idempotent_key = draw_id + "_" + item_index（单道具粒度幂等）
```

---

### 专题四：防刷 & 风控

**Q: 如何防止脚本刷单、概率操纵、设备作弊？**

**答**：五层风控体系，从流量入口到结果审计全覆盖：

```
┌─────────────────────────────────────────────────────────────────┐
│                    五层风控体系                                    │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  L1: 频率控制（网关层，<1ms）                                     │
│  ├─ 单用户：10次/s（滑动窗口），超出直接拒绝                       │
│  ├─ 单 IP：50次/s                                               │
│  ├─ 单设备：200次/min                                            │
│  └─ 拦截率：~30% 恶意流量                                        │
│                                                                 │
│  L2: 请求验签（服务层，<3ms）                                     │
│  ├─ HMAC-SHA256(user_id, pool_id, count, timestamp, nonce)      │
│  ├─ timestamp 容忍30s偏差（防重放）                               │
│  ├─ nonce 60s 内 Redis 去重（防重复提交）                         │
│  └─ 签名不匹配 → 拒绝 + 标记可疑                                │
│                                                                 │
│  L3: 行为检测（异步分析，不阻塞主流程）                            │
│  ├─ 固定间隔检测：50次时间间隔 σ/μ < 5% → 脚本嫌疑               │
│  │   → 强制增加5s冷却 + 标记账号                                  │
│  ├─ 短时爆发：5分钟内 > 500次 → 暂停抽奖 + 人工审核              │
│  ├─ 多账号关联：同设备 > 5账号抽奖 → 标记设备 + 全部关联账号      │
│  └─ 异常时段：凌晨2-6点高频抽奖 → 提高风控等级                    │
│                                                                 │
│  L4: 设备指纹 + 环境检测                                         │
│  ├─ 采集：device_model, os_version, screen, sensors_hash        │
│  ├─ Root/越狱检测 → 风控等级提升                                  │
│  ├─ 模拟器检测 → 限制抽奖次数                                    │
│  └─ 指纹碰撞（不同账号相同指纹）→ 标记养号设备                    │
│                                                                 │
│  L5: 结果审计（事后分析）                                         │
│  ├─ 个人 SSR 率 > 3σ → 标记异常（可能是Bug而非作弊）             │
│  ├─ 连续大奖检测：连续3个SSR → 记录告警（概率约1/150万）          │
│  ├─ 全服产出率实时监控：偏离配置 > 20% → P0 告警                  │
│  └─ 概率种子回溯：可通过 draw_id 重放验证结果是否被篡改           │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

**黑/白/灰名单**：
- **白名单**：测试账号，跳过风控+无需消耗，busi.yaml 配置
- **黑名单**：Redis Set 存储，命中直接拒绝，支持 uid/device_id/IP 三维度
- **灰度名单**：新功能灰度，按 user_id 哈希命中百分比

**核心原则：服务端决定论**
- 所有随机数服务端生成（xorshift128+），客户端无法影响
- Alias 表不下发客户端，客户端只知道结果不知道概率表
- 概率组版本号 + 服务端种子 + 用户状态 → 共同决定结果，抓包无法预判

---

### 专题五：高可用（降级 / 熔断 / 兜底）

**Q: 核心依赖挂了（Redis/MQ/游戏服务器），抽奖系统怎么办？**

**答**：分级降级 + 自动熔断 + 兜底策略矩阵：

**分级降级（ETCD 动态配置，秒级生效）**：

```
lottery.degrade_level: 0~4（动态调整）

Level 0: 正常运行
Level 1: 关闭全服广播/统计上报/概率监控（非核心旁路）
Level 2: 禁止新池开启，禁止连抽（仅允许单抽），降低限流阈值
Level 3: 发货切同步RPC（RT升高但不丢），禁止高消耗操作
Level 4: Redis全挂 → 切MySQL乐观锁模式，QPS从50万降至5万
```

**熔断策略**：
```
触发条件（满足任一自动触发）：
├─ Redis P99 > 30ms（正常 < 5ms）
├─ MQ 发货消息堆积 > 5万
├─ 核心接口错误率 > 0.5%（1分钟窗口）
└─ 上游游戏服务器超时率 > 10%
        ↓
熔断行为：
├─ 半开探测：熔断 10s 后放行 5% 流量试探
├─ 快速失败：抽奖接口返回"系统繁忙，请稍后再试"
└─ 异步积压：已成功的抽奖结果积压在 MQ，服务恢复后自动消费
        ↓
恢复判定：
└─ 连续 10 次探测成功 + 错误率 < 0.1% → 自动关闭熔断
```

**兜底策略矩阵**：

| 故障场景 | 降级策略 | 用户感知 | 恢复方式 |
|---------|---------|---------|---------|
| Redis 单分片宕机 | 哨兵自动切主从（<30s），期间该分片暂停 | 30s内偶发失败 | 自动恢复 |
| Redis 集群全挂 | 切 MySQL FOR UPDATE 乐观锁 | RT 升至 200ms | 手动恢复 |
| RocketMQ 宕机 | 发货切同步 RPC + 本地 WAL | 发货稍慢但不丢 | 手动恢复 |
| 游戏服务器不可达 | 消息积压在 MQ，道具暂时"在路上" | 延迟到账 | 自动恢复 |
| 限量库存 Redis 丢失 | 从 delivery_record COUNT 重建 | 可能少发1-2个 | 自动恢复 |
| 概率异常/Bug | 一键下线奖池 + 通知运营 | 活动暂停 | 人工处理 |

**奖池降级 & 概率降级**：
- **奖池降级**：活跃奖池下线时自动切换到"默认保底池"（只含通用奖励），用户仍可抽奖
- **概率降级**：概率监控异常时，自动切回上一个正常版本的概率表（版本回滚）
- **保底奖/安慰奖**：所有过滤器过滤完后池空 → `complement` 兜底概率组（金币/经验），保证用户不会"抽了个寂寞"

---

### 专题六：数据结构选型 & 随机算法

**Q: 奖池用什么数据结构？随机算法怎么选？为什么不用洗牌？**

**答**：

**奖池数据结构选型**：

| 数据结构 | 适用场景 | 时间复杂度 | 本系统使用 |
|---------|---------|-----------|-----------|
| **Alias Table（权重数组）** | 有放回加权随机 | 预处理O(n), 采样O(1) | ✅ 核心采样算法 |
| **Sorted Array + 二分** | 有放回加权随机 | 采样O(logN) | ❌ 大道具数性能不够 |
| **Redis Set** | 无放回（集合差集） | O(1) 判断/移除 | ✅ 已抽道具集合 |
| **BitSet** | 无放回（位图标记） | O(1) 判断 | ✅ 道具是否已有 |
| **Redis ZSet** | 排行榜/产出排序 | O(logN) 插入/排名 | ✅ 全服大奖排行 |
| **HashMap** | 道具配置/用户状态 | O(1) 读写 | ✅ ProbConfMap/PoolItemsMap |

**Alias Method 详解（为什么不用轮盘赌/二分？）**：

```go
// 轮盘赌法：O(n) —— 道具多时性能差
func roulette(weights []int) int {
    sum := sumAll(weights)
    r := rand.Intn(sum)
    for i, w := range weights {
        r -= w
        if r < 0 { return i }
    }
    return len(weights) - 1
}

// 二分查找法：O(logN) —— 需要维护前缀和数组
func binarySearch(cumWeights []int) int {
    r := rand.Intn(cumWeights[len(cumWeights)-1])
    return sort.SearchInts(cumWeights, r+1)
}

// Alias Method：O(1) —— 预处理后每次采样恒定时间
type AliasTable struct {
    n     int
    prob  []float64  // 每个桶的概率阈值
    alias []int      // 别名跳转目标
}
func (t *AliasTable) Sample() int {
    i := rand.Intn(t.n)            // 随机选桶：O(1)
    if rand.Float64() < t.prob[i] {
        return i                    // 留在本桶
    }
    return t.alias[i]              // 跳转别名
}
```

**为什么不用洗牌算法（Fisher-Yates）？**
- 洗牌适合**无放回 + 全部取出**的场景（如发牌）
- 抽奖是**有放回 + 每次取一个**的场景（每次独立采样）
- 洗牌后取第一个 = O(n) 重排只取一个，浪费计算
- 且洗牌无法处理"权重不等"的场景（需要额外展开）

**无放回场景如何处理？**
- 不用洗牌，而是通过**过滤链**实现：每次抽奖前 `FilterItemByHaveItems` + `FilterItemByDrawn` 移除已有道具，然后对剩余道具重建 Alias 表（或权重缩放采样）
- 性能保证：过滤后的道具数通常远小于总数，重建 Alias 成本低

**Redis ZSet 在抽奖系统中的使用**：
```
场景1：全服大奖排行榜（欧气排行）
  Key: rank:ssr:{pool_id}
  Score: ssr_count（SSR数量）
  Member: user_id
  操作: ZINCRBY +1 on SSR hit; ZREVRANGE 0 99 展示TOP100

场景2：限量道具全服产出进度
  Key: progress:{pool_id}:{item_id}
  Score: timestamp
  Member: user_id
  操作: ZADD 记录谁在什么时间获得; ZCARD 获取已产出数

场景3：概率实时监控（按小时桶）
  Key: stat:{pool_id}:{hour}
  HINCRBY total_draws 1
  HINCRBY ssr_hits 1
  实际SSR率 = ssr_hits / total_draws（对比配置值告警）
```

**权重随机的精度保证**：
- 所有权重使用**整数**表示（如 weight=6 表示 0.6%，总权重为1000）
- 概率计算通过整数除法：`P = weight / total_weight`
- 不使用浮点数运算，避免 `0.1 + 0.2 ≠ 0.3` 的精度问题
- 概率公示时再转为百分比字符串：`fmt.Sprintf("%.4f%%", float64(weight)/float64(total)*100)`

---

### 综合追问：串联全链路

**Q: 从用户点击"十连抽"到看到结果，中间经过了什么？完整链路各环节的延迟是多少？**

**答**：

```
用户点击"十连抽" (客户端)
    │ 0ms
    ▼
① 网关层：限流校验 + 签名验证 + 路由转发
    │ ~3ms（令牌桶判断+HMAC验签+转发）
    ▼
② 响应缓存检查：Redis GET rsp_cache:{user_key}
    │ ~2ms（命中直接返回，80%请求到此结束）
    │ [未命中] ↓
    ▼
③ 分布式锁：Redis SET NX lock:{user_key} TTL=5s
    │ ~1ms
    ▼
④ 读取用户状态：Redis GET state:{game_id}:{user_key}
    │ ~2ms（含 zlib 解压）
    ▼
⑤ 循环10次单轮抽取：
    │ ├─ LuckyValue.Match → 确定概率组 (~0.1ms)
    │ ├─ 加载 Alias 表（本地缓存命中 ~0.01ms）
    │ ├─ 6层过滤链 (~0.5ms，含背包查询缓存)
    │ ├─ Alias Method 采样 (~0.001ms)
    │ ├─ 限量检查 Redis Lua (~1ms，仅限量道具触发)
    │ └─ 更新幸运值 (内存操作 ~0.001ms)
    │ 10轮合计: ~5-15ms（取决于是否命中限量道具重抽）
    ▼
⑥ 本地事务：MySQL（扣货币+写draw_transaction）
    │ ~5ms
    ▼
⑦ RocketMQ 事务消息 Commit
    │ ~3ms
    ▼
⑧ 写回 Redis：更新 UserRealTimeState + 写响应缓存
    │ ~2ms（Pipeline 批量）
    ▼
⑨ 构建响应 + 返回客户端
    │ ~1ms（JSON序列化 + 网络传输）
    ▼
用户看到抽奖结果展示
    总计: P50 ~25ms, P99 ~45ms ✅ < 50ms 目标

===== 以下异步执行，用户已看到结果 =====

⑩ 发货 Worker 消费 MQ → 调用游戏服务器发货
    │ ~200ms（含游戏服务器 gRPC 调用）
    ▼
⑪ 发货成功 → 推送到账通知 → 用户收到"获得xxx"
    │ ~1-3s（端到端延迟）
```

**延迟分布**：
- **用户感知延迟**（点击到看到结果）：P99 < 50ms
- **道具到账延迟**（结果到实际拥有）：P99 < 5s
- **异常补偿延迟**（发货失败到自动补偿）：< 30min
