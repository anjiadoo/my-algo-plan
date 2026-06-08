> **定位**：AI 架构师「上帝视角」完全手册 · AI 工程师技术体系完整沉淀  
**核心逻辑**：每一次 AI 范式革命，都是「数据规模」「算力突破」「算法创新」三者共同驱动的结果。读懂因果链，才能在 Agent 时代分清「什么是真信号、什么是噪声」，而不只是追新闻。
>

---

## 序言：如何读懂一段 AI 史
AI 的发展，从来不是一条平滑的"技术进化曲线"，而是**一系列范式跃迁与寒冬交替**的剧烈震荡。每个纪元，都可以用同一个框架来解读：

```plain
数据规模 / 算力 / 场景需求 跃迁
    ↓ 触发
旧范式的核心局限暴露（不可解、不泛化、不可扩展）
    ↓ 逼出
新的理论矛盾（精度 vs 泛化 / 能力 vs 对齐 / 自主 vs 可控）
    ↓ 产生
理论突破 + 工程实践 + 算力配套
    ↓ 形成
新的 AI 范式
    ↓ 解锁
下一阶段更复杂的应用场景
```

整个 AI 技术史，本质是在回答同一个问题：**如何让一个系统——在没有显式规则的情况下——从经验中学习、泛化到未见数据、并最终在真实世界中可靠地行动？**

贯穿全文有三条主线：

1. **理论主线**：符号逻辑 → 统计学习 → 表征学习 → 预训练 + Scaling Law → 对齐与 RLHF → 测试时计算（Test-Time Compute）→ Agent 外循环
2. **架构主线**：感知机 → MLP → CNN / RNN → Transformer → MoE / 长上下文 → 推理模型 → 单 Agent → Multi-Agent + Harness
3. **工程主线**：手写规则 → 特征工程 → 模型训练 → Prompt Engineering → Context Engineering → Harness Engineering

---

## 第一纪元：符号主义与 AI 寒冬（1950s—1990s）
### "用规则编码智能，被现实击碎两次"
### 1.1 时代背景与核心矛盾
**业务特征**：学术探索期，没有"应用"概念——AI 是哲学家、数学家、心理学家共同回答"机器能否思考"的实验场  
**算力规模**：从电子管到 70 年代的 PDP-11，FLOPS 量级在 10⁴—10⁶  
**核心矛盾**：**真实世界的复杂性 vs 人能写出来的规则数量**——专家以为知识可以被穷举为 IF-THEN 规则，但真实场景的规则量是组合爆炸的

### 1.2 思想起点：图灵测试与达特茅斯会议
**1950，Alan Turing《Computing Machinery and Intelligence》**：提出"模仿游戏"（Imitation Game，即图灵测试）——如果机器能在文字对话中让人类无法判断对方是机器还是人，则可视为具有智能。这篇论文的真正贡献不是测试本身，而是把"智能"从形而上的哲学问题转化为**可被工程验证的行为问题**。

**1956，Dartmouth 会议**：John McCarthy、Marvin Minsky、Claude Shannon、Nathan Rochester 召集，正式提出 "Artificial Intelligence" 这个术语。会议上确立了 AI 的核心信念：**"学习的每一个方面，原则上都可以被精确描述，进而由机器模拟"**——这是符号主义（Symbolicism）的纲领。

**关键人物**：

+ **John McCarthy**：发明 Lisp 语言（1958），AI 第一个专用语言，成为后来 30 年专家系统的事实标准
+ **Marvin Minsky**：MIT AI Lab 创始人，发表《Perceptrons》（1969）——这本书也亲手把第一波神经网络送进了寒冬
+ **Herbert Simon & Allen Newell**：1957 年开发 Logic Theorist（首个 AI 程序，证明数学定理），1972 年提出 GPS（General Problem Solver）

### 1.3 第一次浪潮：感知机与连接主义的初次爆发（1957—1969）
**1943，McCulloch & Pitts**：发表《A Logical Calculus of the Ideas Immanent in Nervous Activity》，提出第一个**人工神经元数学模型**——把神经元抽象为带阈值的二值开关，这是所有后续神经网络的起点。

**1957，Frank Rosenblatt 感知机（Perceptron）**：在 IBM 704 上实现了一个能学习的单层神经网络，可以识别简单图形。媒体一度狂热宣称"机器即将拥有意识"。

**1969 年的当头一棒——Minsky & Papert《Perceptrons》**：

+ 严格数学证明：**单层感知机无法学习 XOR 函数**（线性不可分问题）
+ 暗示多层网络也面临训练困难（当时反向传播尚未提出）
+ 学术界、资金方信心崩塌，神经网络研究进入第一次寒冬

**寒冬启示**：理论的**严格性**比工程的**漂亮 demo** 更能决定一个流派的命运。Minsky 的批评在数学上完全正确，但他低估了 20 年后反向传播 + 算力提升能突破这个限制——这是 AI 史上最深刻的"过早盖棺定论"案例。

### 1.4 第二次浪潮：专家系统与符号主义的工业化（1970s—1980s）
**核心思想**：智能 = 知识 + 推理。把领域专家的知识编码为规则库（Knowledge Base），用推理引擎（Inference Engine）做前向 / 后向推理。

**代表系统**：

| 系统 | 年份 | 领域 | 商业化结果 |
| --- | --- | --- | --- |
| DENDRAL | 1965（Stanford） | 化学分子结构推断 | 学术成功，工业化有限 |
| MYCIN | 1972（Stanford） | 细菌感染诊断 | 准确率超过医生，但医院因责任问题不敢用 |
| XCON / R1 | 1980（DEC） | 配置 VAX 计算机订单 | 每年节省 4000 万美元，**首个商业大成功** |


**LISP 机的兴衰**：80 年代中期，专门运行 LISP 语言的硬件机器（Symbolics、LMI）成为商业热点，资金涌入。日本 1982 年启动「第五代计算机」十年计划（5 亿美元），目标做出能用自然语言交互的 AI 机器。

### 1.5 第二次寒冬（1987—1993）：符号主义的根本性挫败
专家系统在工业化中暴露了**致命缺陷**：

+ **知识获取瓶颈（Knowledge Acquisition Bottleneck）**：专家无法把所有"直觉"显式化；规则库一旦超过几千条就无法维护
+ **脆弱性**：规则之外的输入（比如稍微改写的问题）就会让系统瞬间崩溃，没有任何"泛化"能力
+ **不确定性处理弱**：现实世界充满概率，IF-THEN 难以表达"70% 可能是 X"
+ **LISP 机被通用工作站碾压**：Sun Workstation 性价比远超 Symbolics，专用硬件路线崩盘
+ **第五代计算机计划失败**：日本 10 年投入未达成目标，全球对符号 AI 信心崩塌

**1986 年的伏笔——反向传播复活神经网络**：在符号主义如日中天时，Rumelhart、Hinton、Williams 发表《Learning representations by back-propagating errors》，把反向传播算法系统化——这为后来 2012 年的深度学习革命埋下了 26 年的种子，但当时几乎没人重视。

**1989 年的伏笔——LeCun LeNet**：Yann LeCun 在贝尔实验室用卷积神经网络识别手写邮政编码，准确率达到商用水平。这是 CNN 的第一个工程成功，但同样被淹没在符号主义阵营对神经网络的怀疑中。

### 1.6 纪元总结
```plain
核心矛盾：真实世界的复杂性 vs 显式规则的有限性
核心理论：符号主义、谓词逻辑、推理引擎、知识表示
核心技术：LISP、Prolog、产生式规则系统、专家系统外壳（CLIPS）
架构范式：知识库 + 推理引擎（KB + IE）
时代局限：知识无法穷举，泛化能力为零，不能处理感知任务（视觉、语音）
两次寒冬的根本原因：低估了「学习」的重要性，过度迷信「规则」
留给后世：AI 必须从数据中学，而不是从专家嘴里抠
```

---

## 第二纪元：统计机器学习时代（1990s—2011）
### "放弃模拟人脑，用数学硬刚"
### 2.1 时代背景与核心矛盾
**触发事件**：互联网兴起（1995+）→ 数据规模从 MB 级跃升到 TB 级 → 第一次让"数据驱动学习"具备现实可行性  
**业务特征**：搜索引擎排序、垃圾邮件过滤、推荐系统、信用评分、机器翻译（统计模型）  
**算力规模**：CPU 集群，FLOPS 达到 10⁹—10¹¹  
**核心矛盾**：**模型容量 vs 过拟合**——数据多了，但怎么让模型既能拟合复杂模式又不在测试集上崩盘？这是统计学习理论（Statistical Learning Theory）要回答的问题

### 2.2 理论奠基：从经验风险到结构风险
**1971—1995，Vapnik & Chervonenkis 统计学习理论**：

+ **VC 维（VC Dimension）**：刻画模型「容量」的数学指标——VC 维越大，模型能拟合越复杂的函数，但也越容易过拟合
+ **结构风险最小化（SRM）**：不能只追求训练误差最小（经验风险），还要约束模型复杂度——这是正则化（Regularization）思想的理论根基
+ **核技巧（Kernel Trick）**：通过核函数把低维不可分数据隐式映射到高维空间，让线性模型能解决非线性问题

**1995，SVM（Support Vector Machine）**：Vapnik 提出，结合最大间隔分类器 + 核技巧，在小样本场景下达到当时最强性能。SVM 成为 90 年代后期到 2010 年前的**事实工业标准**——文本分类、人脸识别、生物信息全用 SVM。

### 2.3 集成学习的崛起：弱模型组合的智慧
**1996，Bagging（Breiman）+ 1995，Boosting（Schapire）**：理论突破后，集成学习（Ensemble Learning）成为提升精度的工程利器。

| 算法 | 年份 | 核心思想 | 工业地位 |
| --- | --- | --- | --- |
| **Random Forest** | 2001（Breiman） | 多棵决策树投票，特征 + 样本双随机 | 长期是 Kaggle 比赛 baseline 之一 |
| **AdaBoost** | 1995（Freund & Schapire） | 串行训练，每轮聚焦前一轮的错误样本 | 人脸检测（Viola-Jones）经典算法 |
| **GBDT** | 2001（Friedman） | 梯度提升决策树，每棵树拟合前面所有树的残差梯度 | 推荐、风控、广告 CTR 预估的核心模型 |
| **XGBoost** | 2014（陈天奇） | GBDT 工程极致优化（二阶导数、并行、稀疏感知） | 2015—2018 Kaggle 统治者 |
| **LightGBM** | 2016（微软） | 直方图算法 + Leaf-wise 生长，比 XGBoost 快 10x | 工业界 Tabular 数据事实标准 |
| **CatBoost** | 2017（Yandex） | 处理类别特征的 GBDT，对偏态数据鲁棒 | 类别特征多的场景首选 |


**为什么 GBDT 系是 Tabular 数据之王**：表格数据通常**没有强空间/序列结构**，深度学习的归纳偏置（Inductive Bias）——卷积的局部性、RNN 的时序性——在这里没用，反而是 GBDT 的"特征间非线性交互 + 树状条件分支"天然契合。这个事实直到 2026 年依然成立——**深度学习并没有统治一切**。

### 2.4 概率图模型与 NLP 的统计时代
**HMM（Hidden Markov Model）**：用于语音识别（声学模型）、词性标注、NER。统治 NLP 直到深度学习时代。

**CRF（Conditional Random Field，2001，Lafferty et al.）**：判别式模型，相比 HMM 能利用更丰富特征，是 2010 年代前序列标注任务的 SOTA。

**统计机器翻译（SMT）**：IBM 模型 1—5（1993）确立了"翻译 = 对齐 + 重排序"的统计框架，Google Translate 早期版本（2006—2016）即基于 SMT，2016 年才被神经机器翻译（NMT）替代。

**Latent Dirichlet Allocation（LDA，2003，Blei）**：主题模型，用于文档主题挖掘、推荐系统的语义建模。

### 2.5 深度学习的"地下复兴"（2006—2011）
主流学术界还在迷信 SVM 时，少数研究者坚持神经网络方向。

**2006，Hinton《A Fast Learning Algorithm for Deep Belief Nets》**：

+ 提出 **Deep Belief Network（DBN）+ 逐层预训练（Layer-wise Pre-training）**
+ 解决了深层网络梯度消失的训练难题
+ 论文里第一次正式使用 "Deep Learning" 这个术语
+ **历史意义**：把神经网络从"学界笑柄"重新拉回主流视野，开启深度学习复兴

**2009，ImageNet 数据集发布（李飞飞团队）**：

+ 1500 万张图片，22000 个类别
+ 配套的 ILSVRC 比赛（2010 起每年一届）成为视觉领域标准测试集
+ **战略意义**：李飞飞的判断是「算法已经够了，缺的是足够大的数据」——这个判断为 2012 年的爆发铺好了战场

**2010—2011，GPU 训练神经网络的早期实验**：

+ Andrew Ng + Jeff Dean 在 Google 用 16000 个 CPU 训练神经网络识别猫脸（"Cat Paper"，2012）
+ 多伦多大学的 Hinton 团队开始系统使用 NVIDIA GPU 训练 CNN
+ 这些工程实验直接催生了 2012 年的 AlexNet

### 2.6 纪元总结
```plain
核心矛盾：模型容量 vs 过拟合 / 特征工程 vs 自动化
核心理论：统计学习理论（VC 维、SRM）、贝叶斯框架、集成学习理论
核心技术：SVM、Random Forest、GBDT/XGBoost、HMM/CRF、LDA、DBN
架构范式：特征工程 + 浅层模型 + 集成
时代特征：算法工程师 = 特征工程师，70% 时间在做特征
关键转折：ImageNet 数据集就位（2009）+ GPU 训练成熟（2011）→ 为深度学习引爆做好物理准备
留给后世：「数据 + 算力 + 简单模型」可能比「巧思 + 小数据 + 复杂模型」走得更远
```

---

## 第三纪元：深度学习革命（2012—2017）
### "表征自动学习，端到端击碎特征工程"
### 3.1 时代背景与核心矛盾
**触发事件**：AlexNet 在 ImageNet 2012 上以 15.3% Top-5 错误率（亚军 26.2%）碾压传统方法，深度学习一战成名  
**业务特征**：图像识别、语音识别、机器翻译——感知类任务全面进入"超过人类"的时代  
**数据规模**：ImageNet 千万级标注图像，YouTube/Facebook 等互联网平台带来海量未标注数据  
**算力规模**：消费级 GPU（NVIDIA GTX 580 → K40 → V100）进入研究室，FLOPS 达到 10¹²  
**核心矛盾**：**表示学习（Representation Learning）能否取代手工特征工程**——如果模型能从原始像素 / 波形 / Token 直接学出有效表征，那"特征工程师"这个工种就被消灭了

### 3.2 引爆点：AlexNet 与 ImageNet 的命运一战
**2012.9，AlexNet（Krizhevsky / Sutskever / Hinton）**：

+ 8 层 CNN（5 卷积 + 3 全连接），6000 万参数
+ 关键创新：
    - **ReLU 激活函数**：替代 Sigmoid/Tanh，缓解梯度消失，训练速度提升 6 倍
    - **Dropout（2012，Hinton）**：训练时随机丢弃神经元，强力正则化
    - **Data Augmentation**：图像翻转、裁剪、PCA 颜色扰动
    - **GPU 并行训练**：在两块 GTX 580（3GB 显存）上训练 5—6 天
+ **战略意义**：
    - 学术界：Hinton 学派从边缘走向中心，深度学习成为主流
    - 工业界：Google 立刻收购 Hinton 创办的 DNNresearch（2013），微软、百度跟进
    - 算力：NVIDIA 股价从 12 美元起飞——这是 GPU 厂商命运转折点

### 3.3 CNN 黄金五年：视觉的范式标准化
| 模型 | 年份 | 核心创新 | 历史地位 |
| --- | --- | --- | --- |
| **AlexNet** | 2012 | 深 + ReLU + Dropout + GPU | 引爆深度学习 |
| **VGG** | 2014 | 全部 3×3 小卷积堆叠到 16/19 层 | 简洁优雅，至今仍是 backbone 教科书例子 |
| **GoogLeNet / Inception** | 2014 | Inception 模块（多尺度并行卷积），1×1 卷积降维 | 计算效率与精度的平衡 |
| **ResNet** | 2015（何恺明 / MSRA） | **残差连接（Skip Connection）**：让网络可以训练到 152 层甚至 1000 层 | 解决了"深层网络反而比浅层差"的退化问题，是深度学习史上最重要的架构创新之一 |
| **DenseNet** | 2017 | 每一层与之前所有层稠密连接，特征复用最大化 | 参数效率优于 ResNet |
| **MobileNet** | 2017 | Depthwise Separable Convolution，端侧推理友好 | 手机端 CV 应用的基础 |
| **EfficientNet** | 2019 | NAS 搜索 + 三维度复合缩放（深度/宽度/分辨率） | AutoML 时代的 SOTA |


**ResNet 的深远影响**：残差连接 `y = F(x) + x` 不仅解决了 CNN 的训练难题，后来被 **Transformer 完整继承**——今天的 GPT、Claude、Gemini 每一个 Transformer Block 中的 Add & Norm，本质就是 ResNet 思想的延续。这是 AI 架构史上最具复利的设计原语之一。

### 3.4 RNN/LSTM 与序列建模
**1997，LSTM（Hochreiter & Schmidhuber）**：通过门控机制（输入门、遗忘门、输出门）解决 RNN 的长程梯度消失问题。但要等到 2014 年才被广泛应用。

**2014，Seq2Seq（Sutskever et al.）**：Encoder—Decoder 架构，把一个序列映射到另一个序列，开启了**神经机器翻译（NMT）**时代。

**2014，Attention Mechanism（Bahdanau et al.）**：在 Seq2Seq 中引入注意力——Decoder 的每一步可以"看"Encoder 的所有位置，按相关性加权——解决了 Seq2Seq 中长序列信息压缩损失的问题。**这是后来 Transformer 的直接思想源头**。

**2016，Google NMT（GNMT）上线**：Google Translate 全面切换到 LSTM + Attention 架构，翻译质量跃升 60%，**统治传统统计机器翻译 23 年的 SMT 范式被一次性终结**。

### 3.5 生成模型的诞生：GAN 与 VAE
**2013，VAE（Variational Autoencoder，Kingma & Welling）**：变分自编码器，从概率角度建模数据生成，潜空间连续平滑，但生成质量偏模糊。

**2014，GAN（Generative Adversarial Network，Ian Goodfellow）**：

+ 两个网络博弈：生成器（Generator）造假，判别器（Discriminator）辨真
+ "训练一个会造假的网络去骗一个会辨假的网络"，纳什均衡时生成器学到真实分布
+ **历史意义**：开启了"图像生成"领域的 8 年黄金期（2014—2022），从 DCGAN → StyleGAN → BigGAN，直到 2022 年扩散模型（Diffusion）登顶

**2015，DCGAN**：把 CNN 引入 GAN，第一次让生成的图像在视觉上"像样"。

**2018—2019，StyleGAN（NVIDIA）**：风格迁移 + 渐进式生成，**ThisPersonDoesNotExist.com** 一夜爆红，公众第一次意识到"AI 可以无中生有创造图像"。

### 3.6 强化学习的高光：AlphaGo 与 DeepMind 的科学化路径
**2013，DeepMind《Playing Atari with Deep Reinforcement Learning》**：DQN（Deep Q-Network）算法，端到端学习 Atari 游戏，**首次证明深度学习 + 强化学习的可行性**。这篇论文让 Google 在 2014 年以 5 亿美元收购了仅 50 人的 DeepMind。

**2016.3，AlphaGo 击败李世石**：

+ 4:1 击败围棋世界冠军，第二局的"神之一手"（Move 37）让人类高手集体震撼
+ 核心技术：**蒙特卡洛树搜索（MCTS）+ 策略网络 + 价值网络 + 自我对弈**
+ **历史意义**：围棋被普遍认为是「20 年内 AI 攻不下来的圣杯」，AlphaGo 把这个时间表压缩到 0；公众认知里 AI 第一次有了"超人"色彩

**2017.10，AlphaGo Zero**：

+ 完全不用人类棋谱，从零自我对弈 40 天 → 100:0 击败 AlphaGo Lee
+ **战略启示**：人类经验有时是天花板而不是地板——纯自学习反而能突破人类水平

**2017.12，AlphaZero**：泛化 AlphaGo Zero 到围棋、国际象棋、将棋三个领域，4 小时学习就击败所有人类工程多年的引擎。

**这一系列工作奠定了现代强化学习的工程范式，也直接影响了 2024 年 OpenAI o1 推理模型的 Test-Time Search 思路**。

### 3.7 引爆点 2.0：Transformer 横空出世
**2017.6，Vaswani et al.《Attention is All You Need》**：

+ 论文标题就是宣言：抛弃 RNN/CNN，**完全用注意力机制**构建序列模型
+ 核心组件：
    - **Self-Attention（自注意力）**：每个 Token 计算与序列中所有其他 Token 的相关性
    - **Multi-Head Attention**：多组并行注意力，捕捉不同子空间的关系
    - **Positional Encoding**：因为没有循环结构，需要显式注入位置信息
    - **Encoder—Decoder 架构 + Layer Norm + Residual Connection**
+ 工程优势：**完全并行**（RNN 必须串行处理时间步），训练速度比 LSTM 快一个数量级
+ **历史评价**：这是过去 10 年最重要的一篇 AI 论文，没有之一——所有现代大模型（GPT、Claude、Gemini、LLaMA、DeepSeek）的底层都是 Transformer

**Transformer 为什么改变了一切**：

1. **可并行性**：训练规模从"被 GPU 数量限制"变成"被资本限制"，scale 成为可能
2. **长程依赖**：自注意力让任意两个位置直接交互，彻底解决 RNN 的长程依赖衰减
3. **统一性**：不仅适用于 NLP，后来在视觉（ViT）、语音、多模态、生物（AlphaFold 2）全面通吃——**Transformer 成了通用学习架构**

### 3.8 纪元总结
```plain
核心矛盾：手工特征工程 vs 端到端学习
核心理论：表示学习、反向传播、残差连接、注意力机制
核心技术：CNN（ResNet）、LSTM/Seq2Seq、GAN、DQN/AlphaGo、Transformer
架构范式：端到端深度网络，从数据直接学表征
关键工程突破：GPU 通用计算、CUDA 生态、PyTorch（2016）/ TensorFlow（2015）
中国里程碑：百度全面切换深度学习语音识别（2014）、商汤 / 旷视成立（2014）、阿里推荐系统深度学习化（2016+）
留给后世：「足够大的模型 + 足够多的数据 + Transformer」是下一个时代的入场券
```

---

## 第四纪元：预训练范式与 Scaling Law（2018—2022）
### "从手工训练每个任务，到一个通用模型解决所有任务"
### 4.1 时代背景与核心矛盾
**触发事件**：BERT（2018）+ GPT-2（2019）证明了"大规模无监督预训练 + 下游微调"的有效性  
**业务特征**：从"每个任务一个模型"转向"一个基座模型 + N 个 Adapter"  
**数据规模**：从 GB（标注数据）跃升到 TB（互联网文本爬取）  
**算力规模**：万卡 GPU 集群成为大模型训练标配，FLOPS 达到 10¹⁸—10²³  
**核心矛盾**：**模型规模放大到什么程度才会涌现新能力？投入和产出的关系是什么？**——这是 Scaling Law 要回答的问题

### 4.2 预训练范式的奠基
**2018.10，BERT（Bidirectional Encoder Representations from Transformers，Google）**：

+ 用 Transformer Encoder + 双向掩码语言模型（Masked Language Model，MLM）做预训练
+ 在 11 项 NLP 任务上刷新 SOTA，多数任务超过人类基线
+ **核心范式**：先在大规模无标注语料上预训练，再用少量标注数据微调（Fine-tuning）
+ 工业落地：2019.10 Google 宣布将 BERT 应用于搜索英文查询（首批覆盖约 10%，后续逐步扩展），是 Google 搜索 5 年来最大的算法升级

**2018.6，GPT-1（OpenAI）**：早于 BERT，但选择了**单向 Transformer Decoder + 自回归语言建模**（预测下一个 Token）的路线——这条路线在当时被 BERT 的双向建模碾压，但**约 2 年后用 GPT-3（2020.5）反向证明了它才是通用智能的正确路径**。

**两条路线的本质区别**：

+ BERT 的 MLM 善于理解（NLU），但生成能力天生弱
+ GPT 的自回归适合生成（NLG），且**涵盖理解作为副作用**——因为要生成连贯文本，模型必须理解语境
+ 历史告诉我们：**自回归 + 大规模 + Scaling = 通用智能的密钥**

### 4.3 GPT 系列与 Scaling Law 的发现
**2019.2，GPT-2（15 亿参数）**：

+ OpenAI 因"担心被滥用"延迟发布完整版本，引发学术界争议
+ 第一次展示了"零样本"（Zero-Shot）能力——不微调直接做任务
+ **范式启示**：模型大到一定程度，开始具备元能力（meta-capability）

**2020.5，GPT-3（1750 亿参数）**：

+ 比 GPT-2 大 100 倍，训练数据 570GB（清洗后）
+ 训练算力约 3640 PetaFLOPS-days，单次训练成本估计 1200 万美元
+ **里程碑能力**：
    - **Few-Shot Learning**：在 Prompt 中给几个示例，就能模仿做新任务
    - **In-Context Learning**：模型不更新参数，仅靠上下文就能"学会"——这是最反直觉也最重要的发现
+ **战略意义**：GPT-3 让 OpenAI 一举成为 AI 领军者，奠定了之后 5 年的 LLM 路线

**2020.1，Kaplan et al.《Scaling Laws for Neural Language Models》（OpenAI）**：

+ 系统研究"参数量 N、数据量 D、算力 C"与模型损失（Loss）之间的关系
+ **核心发现**：Loss 随 N、D、C 呈**幂律下降**——即 `Loss ∝ N^(-α)`，且三者比例需要平衡
+ **历史意义**：把"训练大模型"从一门玄学变成了**可预测的工程**——给定预算，可以算出最优的参数量和数据量

**2022.3，Chinchilla（DeepMind）**：

+ 发现 GPT-3 等模型其实**严重欠训练**——参数太大但数据太少
+ 提出 Chinchilla Optimal：**参数量与训练 Token 数应近似 1:20**
+ 用 70B 参数 + 1.4T tokens 击败 175B 参数的 GPT-3
+ **工程启示**：盲目堆参数是错的，**数据和参数要按比例同步增长**——这条规律直接指导了后来 LLaMA、Gemini、DeepSeek 等所有主流模型的训练配比

### 4.4 涌现能力（Emergent Abilities）：量变到质变
**2022.6，Wei et al.《Emergent Abilities of Large Language Models》（Google）**：

+ 系统观察到：**许多能力在模型规模较小时几乎为 0，但在某个临界规模突然出现**——这就是涌现
+ 涌现能力清单：算术、多步推理、code 生成、Chain-of-Thought、跨语言迁移……
+ **争议（2023+）**：Schaeffer et al. 后续指出部分"涌现"是评估指标不连续造成的假象，但**实践中确实存在某些能力的非线性跃迁**

**2022.1，Chain-of-Thought（CoT）Prompting（Google）**：

+ 在 Prompt 中加入"Let's think step by step"或者推理示例
+ 数学题、逻辑题准确率从 17.9% 跃升到 56.9%（PaLM-540B 在 GSM8K 上的原论文数据；后续叠加 Self-Consistency 可达 74.4%）
+ **本质**：把模型的"内部思考"显式化为输出 Token，让模型用更多算力做推理
+ **历史意义**：这是 2024 年 OpenAI o1 推理模型路线的思想原型

### 4.5 对齐与 RLHF：让模型"听话"
**2017，PPO（Proximal Policy Optimization，Schulman et al.）**：稳定的策略梯度强化学习算法，是后来 RLHF 的标准优化器。

**2022.1，InstructGPT（OpenAI，论文 2022.1.27 公布；同年部署到 API）**：

+ 用 **RLHF（Reinforcement Learning from Human Feedback）**对 GPT-3 做对齐
+ 三阶段：
    1. **SFT（Supervised Fine-Tuning）**：用人工标注的高质量回答微调
    2. **奖励模型（Reward Model）**：人工对模型的多个回答排序，训练 RM 预测"哪个回答更好"
    3. **PPO 强化学习**：用 RM 作为奖励信号，PPO 优化模型策略
+ **核心价值**：1.3B 参数的 InstructGPT 比 175B 的 GPT-3 更受用户喜爱——**对齐比规模更重要**

**Anthropic Constitutional AI（CAI，2022.12）**：

+ RLHF 需要大量人工标注，成本高且不一致
+ CAI 用"宪法"（一组高层原则）让模型自我批评、自我修改
+ **战略意义**：把人类反馈替换为 AI 反馈（RLAIF），让对齐过程可规模化——这是 Anthropic 后来 Claude 系列的核心方法论

### 4.6 多模态的雏形：CLIP 与 DALL·E
**2021.1，CLIP（Contrastive Language-Image Pre-training，OpenAI）**：

+ 4 亿对（图像, 文本）训练数据，图像编码器 + 文本编码器映射到同一向量空间
+ **零样本图像分类**：不需要微调，直接用文本描述新类别即可分类
+ **战略意义**：开启了"用自然语言操作视觉"的时代——后来的 Stable Diffusion、Flamingo、GPT-4V 全部基于这个思想

**2021.1，DALL·E**：基于 GPT-3 风格的自回归 Transformer + dVAE 架构（**先把图像离散化为 token，再像写文字一样自回归生成图像 token**）；CLIP 在此处仅作为**输出端的 reranker**（从多个候选中挑最匹配 prompt 的图），并不参与生成过程。直到 **2022.4 DALL·E 2 才真正将 CLIP 嵌入生成过程（unCLIP 架构 + Diffusion Decoder）**——这是文生图第一次商业级可用。

**2022，Diffusion Models 登顶**：

+ **DDPM（2020）**：Denoising Diffusion Probabilistic Models 理论奠基
+ **Stable Diffusion（2022.8，Stability AI 开源）**：在消费级 GPU 上即可运行，文生图全民可用
+ **DALL·E 2 / Midjourney v3 / Imagen**：商业化模型百花齐放
+ **GAN 时代的终结**：扩散模型生成质量、稳定性、可控性全面碾压 GAN，**统治图像生成至今**

### 4.7 中国大模型的起步（2021—2022）
| 模型 | 单位 | 时间 | 特点 |
| --- | --- | --- | --- |
| **悟道 2.0** | 智源研究院 | 2021.6 | **1.75 万亿参数 MoE 模型**，刷新当时全球最大模型总参数纪录（注：MoE 总参数大但单次激活稀疏，与稠密模型不可直接对比） |
| **盘古** | 华为 | 2021.4 | 千亿稠密参数，瞄准行业落地 |
| **M6** | 阿里达摩院 | 2021.3 / 2021.10 | 多模态大模型，初版千亿稠密；M6-10T (2021.10) 扩展到 10 万亿参数 MoE |
| **GLM-130B** | 清华 | 2022.10 | 双语开源，对标 GPT-3 |


**这个阶段的中国大模型**：参数堆得很大，但训练数据质量、对齐工程、生态构建均落后 OpenAI 一代——这是国内 AI 圈对 ChatGPT 冲击毫无心理准备的根本原因。

### 4.8 AI for Science 主线的开端：AlphaFold 2 与新的科学范式
**2020.7—2021.7，DeepMind AlphaFold 2**：

+ 在 CASP14 蛋白质结构预测大赛上以接近实验精度（中位 GDT_TS 92.4）碾压所有传统方法
+ 把困扰生物学家 50 年的"蛋白质折叠问题"事实上**解决到工程可用水平**
+ 2021.7 开源代码 + 公开人类全部 ~20000 种蛋白质结构数据库（与 EMBL-EBI 合作），后扩展至 2 亿+ 蛋白
+ **范式意义**：第一次证明深度学习不仅能做语言 / 视觉，还能**压缩自然规律本身**——AI 从生产力工具升级为科学加速器

**2022—2024，AI4S 全面铺开**：

+ **AlphaFold-Multimer / RoseTTAFold（David Baker 团队）**：复合物结构预测
+ **GraphCast（DeepMind 2023）**：10 天天气预报精度首次超过传统数值模式
+ **Pangu-Weather（华为 2023）**：3 秒完成传统超算 4 小时的气象预测，登 Nature
+ **GNoME（DeepMind 2023.11）**：发现 220 万个新晶体材料，扩大已知稳定材料 10 倍
+ **AlphaProof / AlphaGeometry 2（DeepMind 2024.7）**：IMO 数学奥赛达到银牌水平（接近金牌线）

**2024.5，AlphaFold 3**：扩展到蛋白-蛋白、蛋白-DNA/RNA、蛋白-小分子复合物，向"虚拟药物筛选"迈出关键一步；**Isomorphic Labs 把 AlphaFold 工业化为药物研发管线**。

**2024.10，AI 第一次拿下双诺奖（AI 史标志性时刻）**：

+ **物理学奖**：John Hopfield（Hopfield Network 1982）+ Geoffrey Hinton（玻尔兹曼机、反向传播推广）——表彰"用物理方法奠定现代神经网络基础"
+ **化学奖**：Demis Hassabis + John Jumper（AlphaFold）+ David Baker（计算蛋白设计 RoseTTA / RFdiffusion）——表彰"计算蛋白结构预测与设计"
+ **战略意义**：诺贝尔委员会以最高科学荣誉确认了"AI 是基础科学的核心方法"，**改变了 AI 的学科定位**

**留给后世的认知**：AI 不只是"能写代码 / 做客服 / 出图"，它的更长期价值是**作为科学发现的通用引擎**——Hassabis 把这一愿景表述为"用 AI 解决智能，再用智能解决一切"。这条 AI4S 主线在 2025—2026 进一步与 Agent 融合（科研 Agent、自动实验、闭环优化），成为不亚于 Agent 工程的并行主航道。

### 4.9 纪元总结
```plain
核心矛盾：模型规模与能力的关系 / 通用预训练 vs 专用任务 / 数据 vs 参数最优配比
核心理论：Scaling Law、Chinchilla 最优比、涌现能力、RLHF/Constitutional AI、AI4S 范式
核心技术：BERT/GPT 系列、Transformer Decoder-only、CLIP/Diffusion/DDPM、PPO 对齐、AlphaFold
架构范式：「大规模预训练 + 对齐微调」的双阶段范式正式确立；AI4S 主线开端
工程突破：万卡 GPU 训练、3D 并行（数据/张量/流水线）、ZeRO 优化器、混合精度训练
中国里程碑：万亿参数模型出现（2021），但工程对齐能力差距明显
留给后世：「Scaling 是通用智能的捷径，但不是终点；对齐和数据质量同等重要；AI 也开始压缩自然规律本身」
```

---

## 第五纪元：生成式 AI 爆发与多模态融合（2022.11—2024）
### "ChatGPT 把 AI 从论文带进了每个人的浏览器"
### 5.1 时代背景与核心矛盾
**触发事件**：2022.11.30，ChatGPT 上线——5 天破百万用户，2 个月达到 1 亿月活，**人类历史上最快用户增长的产品**  
**业务特征**：AI 从"工程师工具"变成"全民工具"，每个 SaaS 产品都在思考"怎么加 AI"  
**数据规模**：互联网级文本（数十 TB）+ 多模态（图像 / 视频 / 音频，PB 级）  
**算力规模**：单次训练成本破亿美元（GPT-4 估算 5000 万—1 亿美元），单家公司算力投入数百亿美元  
**核心矛盾**：**通用能力 vs 商业化场景**——模型什么都"能做一点"，但"做对、做稳、做出差异化"是另一个问题；同时**多模态融合**成为新战场

### 5.2 ChatGPT 现象级冲击与中国大模型补课
**ChatGPT 为什么是分水岭**：

+ **产品形态**：把 GPT-3.5 + RLHF 包装成对话界面，让普通人无门槛使用——产品工程的胜利远大于模型本身
+ **病毒传播**：5 天百万用户，靠的不是营销而是"你试试看"的口碑——产品力本身具备病毒性
+ **认知冲击**：2022 年 12 月起，全球科技圈、投资圈、政府圈集体陷入"AGI 焦虑"，Sam Altman 一夜成为最有影响力的科技领袖

**2023.3，GPT-4 发布**：

+ 多模态（接受图像输入）；上下文窗口默认 **8K Tokens**，同期推出 32K 版本（gpt-4-32k）；2023.11 GPT-4 Turbo 进一步扩到 128K
+ 律师资格考试（Bar Exam）排名前 10%（GPT-3.5 是后 10%）
+ **闭源**：OpenAI 不再公开训练细节、参数量、数据来源——AI 开源社区与闭源前沿的鸿沟正式形成

**中国大模型补课（2023）**：

| 模型 | 单位 | 发布 | 定位 |
| --- | --- | --- | --- |
| **文心一言** | 百度 | 2023.3 | 国内首个商用 ChatGPT 类产品 |
| **通义千问** | 阿里 | 2023.4，开源 7B/14B/72B（2023.9—） | 开源策略激进，国内开源生态领头 |
| **ChatGLM** | 智谱 AI | 2023.3—，开源 6B/130B | 学术派背景，长期开源 |
| **百川（Baichuan）** | 百川智能（王小川） | 2023.6—，开源 7B/13B | 开源 + 商业化并行 |
| **Moonshot Kimi** | 月之暗面（杨植麟） | 2023.10— | **长上下文（200 万字）差异化定位** |
| **豆包** | 字节跳动 | 2023.8— | 字节国内 C 端 AI 助手，2024 全面爆发 |
| **DeepSeek** | 幻方量化 | 2023.11—（V1 7B） | 当时不起眼，**两年后改写全球 AI 格局** |


**国内 AI 落地的核心特点**：

+ **开源策略激进**：与 OpenAI 闭源形成对比，阿里 Qwen、智谱 GLM、DeepSeek 都走开源路线，反而获得了全球开发者社区
+ **商业化更务实**：欧美在烧钱做基座，中国早期就在做应用（教育、客服、营销文案）
+ **算力受限**：H100 等出口管制，倒逼国内做更高效的训练算法（如 DeepSeek 的 FP8 训练）

### 5.3 国产模型关键节点（2024—2025 不容忽视的这条线）
**2024.5，DeepSeek-V2 引爆"价格战"**：

+ 首发 **MLA（Multi-head Latent Attention）+ DeepSeekMoE**，KV Cache 减少约 93%
+ API 定价 ¥1/百万输入 Token，是当时 GPT-4 价格的 1/100——**直接把国产 LLM API 价格拉到地板**
+ 字节豆包、阿里 Qwen、智谱、百度文心一周内集体跟进降价，**触发"百模大战 → 价格战"**
+ 战略意义：让"用 API 跑业务"在中国成本结构上变得可行，加速 to B 落地

**2024.6—2025.4，Qwen 系列演进（阿里）**：

+ **Qwen2 (2024.6)**：开源 7B—72B，多语言能力出色
+ **Qwen2.5 (2024.9)**：在数学、代码、多语言上全面对标 Llama 3，HuggingFace 开源下载量长期 Top
+ **Qwen2.5-Max / Qwen2.5-Coder / Qwen2.5-VL**：行业垂直版本矩阵
+ **Qwen3 (2025.4)**：混合推理架构（同模型可切换 thinking / non-thinking），开源 0.6B—235B 全尺寸，**事实上已成为全球开源 LLM 的最重要家族之一**

**2025，Kimi 与月之暗面的反击**：

+ **Kimi K1.5（2025.1，与 R1 几乎同期）**：推理模型，发表论文披露多模态推理 RL 训练方法
+ **Kimi K2（2025.7）**：开源万亿参数 MoE，主打 Agent 能力，**国产开源 Agent 基座的代表**

**其他主力**：

+ **MiniMax-01（2025.1）**：456B 总参 / 45.9B 激活 MoE，原生 4M 长上下文
+ **智谱 GLM-4 / GLM-4.5（2024—2025）**：稳定的开源中坚
+ **腾讯混元 Large / Hunyuan-T1**：腾讯系开源 MoE + 推理模型矩阵
+ **阶跃 Step / 百川 Baichuan / 零一万物 Yi**：各有特色的多模态 / 通用模型

> **2025 年的国产 AI 全景**：从"百模大战"沉淀为**"DeepSeek + Qwen + Kimi + 豆包/混元"几条主航道**，开源生态已经事实上与海外平分秋色，**部分维度（成本、长上下文、Agent 工程）甚至领先**。
>

### 5.4 Claude 系列：Anthropic 的差异化路线
**Anthropic 的诞生（2021）**：Dario Amodei 等 OpenAI 核心员工出走创立，主打"AI 安全"——这不只是公关定位，是真实的技术方法论分歧。

| 版本 | 时间 | 核心进步 |
| --- | --- | --- |
| **Claude 1（首发）** | 2023.3 | 9K 上下文，主打安全对话 |
| **Claude 1.3** | 2023.5 | **100K 上下文（当时业界最大）**，主打长文档处理 |
| **Claude 2** | 2023.7 | 编程能力首次接近 GPT-4 |
| **Claude 3 系列**（Haiku/Sonnet/Opus） | 2024.3 | 三档定价覆盖全场景，Opus 全面对标 GPT-4 |
| **Claude 3.5 Sonnet** | 2024.6 | 编程任务首次稳定超过 GPT-4，成为开发者首选 |
| **Claude 3.5 Sonnet (new) / 3.5 Haiku** | 2024.10 | 同期发布 Computer Use 能力 |
| **Claude 3.7 Sonnet** | 2025.2 | 引入 Extended Thinking（推理深度可控） |
| **Claude 4 / Sonnet 4 / Opus 4** | 2025.5 | Agentic 编程能力大幅提升 |
| **Claude Sonnet 4.5** | 2025.9 | 长任务自主执行能力（连续编程数十小时）成为标杆 |
| **Claude Opus 4.5 / Sonnet 4.6 等后续迭代** | 2025.11— | 1M 上下文成为旗舰标配，多模态与 Agent 任务持续提升 |


**Anthropic 的方法论积累**：

+ **Constitutional AI**：用宪法 + AI 反馈替代部分 RLHF
+ **Mechanistic Interpretability**：可解释性研究领先业界，2024 年发布 Sparse Autoencoders 解析 Claude 内部特征
+ **Computer Use（2024.10）**：Claude 直接控制屏幕、键盘、鼠标——是 Agent 时代的关键产品里程碑
+ **MCP（2024.11）**：开放协议标准让 LLM 与工具连接标准化

### 5.5 Gemini 与 Google 的反击
**Google 的处境（2022.12）**：内部 LaMDA 早就有，但因"声誉风险"不敢发布，被 ChatGPT 抢先后陷入"红色警报"状态。

**2023.12，Gemini 1.0**：原生多模态架构（不是事后拼接），但产品体验不如 GPT-4。

**2024.2，Gemini 1.5 Pro**：

+ **百万 Token 上下文**：第一个商用百万级长上下文模型
+ 多模态理解能力极强（视频、音频、图像、代码混合输入）
+ 一个工程奇迹：用 MoE + 改进的注意力实现 1M 上下文低成本推理

**2024.12，Gemini 2.0**：原生多模态生成（直接出图、出语音），并整合到 Google 搜索 / Workspace。

**2025—2026，Gemini 2.5 / 3**：在长上下文（已扩至 2M+）和原生多模态推理上继续领先；Google 凭借 TPU 自研芯片 + 海量数据（搜索、YouTube、Gmail）的飞轮，**重新追上前沿**。

### 5.5 开源大模型的崛起与生态
**2023.2，LLaMA（Meta）**：

+ 7B / 13B / 33B / 65B，原计划仅供研究使用，但权重在 4chan 泄露后全网扩散
+ **意外结果**：开源社区围绕 LLaMA 在 3 个月内做出了 Alpaca、Vicuna、Guanaco 等数十个变体——**催生了整个开源 LLM 生态**

**2023.7，LLaMA 2**：Meta 转向商业化开源（允许商用），成为开源旗舰直到 LLaMA 3。

**2024.4，LLaMA 3**：首发 8B / 70B；**2024.7 LLaMA 3.1 进一步发布 405B**（首个开源稠密模型对标 GPT-4），开源生态彻底反超。

**LLaMA 模型权重的扩散直接催生了**：

+ **HuggingFace**：从 2017 年 NLP 库公司变成全球开源 AI 模型托管中心，2024 年估值 45 亿美元
+ **vLLM / TGI / SGLang**：高性能推理框架百花齐放
+ **LoRA / QLoRA**：低秩微调让消费级 GPU 也能微调大模型，民主化训练
+ **本地部署生态**：Ollama、llama.cpp、LM Studio 让大模型运行在 MacBook 上成为日常

### 5.6 训练 / 推理基础设施的工程革命（2022—2024）
模型能 scale 到 GPT-4 / Claude 3 / DeepSeek-V3 这一代，**底层是一系列基础设施级的算法与系统创新支撑的**——这是大多数科普读物会跳过、但 AI 工程师必须理解的"隐形地基"。

**训练侧的关键突破**：

| 技术 | 时间 | 出品 | 解决的问题 | 影响量级 |
| --- | --- | --- | --- | --- |
| **FlashAttention v1/2/3** | 2022.5 / 2023.7 / 2024.7 | Tri Dao | Attention 显存 O(n²) 爆炸 | 训练速度 2—4×、显存大幅降低，**所有现代 LLM 训练标配** |
| **ZeRO-1/2/3 + DeepSpeed** | 2020—2022 | Microsoft | 单卡塞不下大模型 | 让万亿参数训练成为可能 |
| **3D 并行（DP+TP+PP）** | 2021— | NVIDIA Megatron-LM | 单一并行维度无法 scale | 万卡集群训练事实标准 |
| **混合精度（FP16/BF16/FP8）** | 2017— / 2023 FP8 落地 | NVIDIA | 算力与显存效率 | FP8 训练让 H100 算力实际可用，DeepSeek-V3 工程示范 |
| **MoE 路由演进**（DeepSeekMoE / Auxiliary-loss-free balancing） | 2024 | DeepSeek | 专家负载不均、训练不稳定 | DeepSeek-V2/V3 关键创新 |


**推理侧的关键突破**：

| 技术 | 时间 | 出品 | 解决的问题 | 影响量级 |
| --- | --- | --- | --- | --- |
| **PagedAttention / vLLM** | 2023.6 | UC Berkeley | KV Cache 内存碎片 | 推理吞吐 2—4×，**开源推理事实标准** |
| **Speculative Decoding** | 2023— | DeepMind / Google | 自回归生成串行慢 | 输出延迟降低 2—3× |
| **量化（GPTQ / AWQ / GGUF / GGML）** | 2022—2023 | 各家 | 大模型部署成本 | INT4/INT8 让 70B 模型跑在消费级 GPU、Mac M2 上 |
| **Continuous Batching** | 2022—2023 | Anyscale / vLLM | 推理 Batch 利用率低 | 吞吐量再提升 1.5—2× |
| **MLA（Multi-head Latent Attention）** | 2024.5 | DeepSeek-V2 | KV Cache 巨大导致长上下文成本爆炸 | KV Cache 减少约 93%，是 DeepSeek 性价比的核心秘密 |


**长上下文工程链路**：从 4K → 1M Tokens 不是模型变大变出来的，而是一组算法叠加：

+ **RoPE（Rotary Positional Embedding，2021 Su et al.）**：相对位置编码，长上下文外推的基础
+ **Position Interpolation（PI）/ NTK-aware scaling（2023）**：免训练扩展 RoPE 上下文范围
+ **YaRN（2023.9）**：Yet another RoPE extensioN，对 PI 的工程改进，是 LLaMA / Qwen 等扩长的标准方法
+ **Ring Attention（2023.10）/ Striped Attention**：跨 GPU 切分序列维度，支持百万级上下文训练
+ **GQA / MQA（Grouped/Multi-Query Attention）**：减少 KV Head 数，长上下文显存优化

> **核心认知**：**没有 FlashAttention 就没有 GPT-4 的训练规模，没有 vLLM / MLA / 量化就没有 LLM 推理的商业化落地**。基础设施层是 AI 工程师真正的护城河，模型新闻只是冰山一角。
>

### 5.7 训练后阶段的算法革命：从 PPO 到 DPO / GRPO
RLHF 的工程门槛极高（PPO 需要同时维护四个模型：Actor、Critic、Reward、Reference），2023 年起**训练后阶段（Post-training）**出现一系列简化与突破：

+ **DPO（Direct Preference Optimization，Rafailov et al. 2023.5）**：直接用偏好对（chosen/rejected）做对比学习，**去掉 RM + PPO 全过程**，公式上等价于隐式 Reward Model；2024 起成为**工业界 Post-training 的事实主流**（Llama 3、Qwen2、Mistral 全用 DPO 系）
+ **IPO / KTO / ORPO / SimPO（2024）**：DPO 家族的各种变体，分别针对噪声偏好、单点反馈、合一 SFT+DPO、长度偏置等问题
+ **RLAIF（Anthropic CAI 系）**：用 AI 生成偏好标签替代人工，规模化对齐
+ **GRPO（Group Relative Policy Optimization，DeepSeek 2024）**：去掉 Critic，用同一 prompt 多个采样的相对优势（Group-relative advantage）训练；**成为 R1 推理 RL 的核心算法**，是 2025 年开源复现潮的算法基础
+ **RLVR（RL with Verifiable Rewards）**：在数学、代码等可验证答案的任务上，用程序判分替代 RM——是推理模型（o1/R1）能用 RL 训出来的关键前提

> **方法论变化**：2024 年起，"做对齐"已经从"PPO + RM" 转向 **"SFT + DPO 系 + 任务可验证 RL（RLVR/GRPO）"** 三件套，这是工程师真正应该掌握的当代 Post-training 流水线。
>



**MoE（Mixture of Experts）历史**：

+ 1991 年 Jacobs et al. 提出 MoE 的雏形
+ 2017 年 Google《Outrageously Large Neural Networks》工程化
+ 2022 年 Switch Transformer / GLaM 开始在 LLM 上验证
+ **2023.12，Mistral Mixtral 8×7B 开源**：让全球开发者第一次跑起 MoE 模型
+ **2024，GPT-4 据 SemiAnalysis 等推测为 8×220B MoE（OpenAI 未官方证实）**：业界普遍认为头部模型都是 MoE
+ **2024.12，DeepSeek-V3（671B 总 / 37B 激活）开源**：以 1/10 的训练成本（约 600 万美元）追平 GPT-4 级别，**全球哗然**

**MoE 的工程价值**：

+ **总参数大但激活参数小**：模型容量高，但推理算力成本低
+ **专家分工**：不同 Token 路由到不同专家，实现"参数级别的专业化"
+ **扩展性**：千亿参数稠密模型训练几乎不可行，但万亿 MoE 是现实的

### 5.8 Transformer 之外：状态空间模型（SSM / Mamba）的另一条路
Transformer 并非没有挑战者。**2023.12，Mamba（Albert Gu & Tri Dao）**：

+ 基于状态空间模型（SSM）+ 选择性扫描机制，实现**线性时间复杂度**（Transformer 是 O(n²)）
+ 长序列推理速度优势显著，是百万级 Token 推理的另一种解法
+ **2024，Mamba-2 / Jamba（AI21）/ Zamba / Hymba（NVIDIA）等 Hybrid 架构涌现**：把 SSM 与 Attention 混合使用，取长补短

**为什么 2026 年 Transformer 仍占绝对主导**：

+ 生态（CUDA Kernel、训练数据、评测、工具链）全部围绕 Transformer 打造
+ 推理加速套件（FlashAttention、vLLM、Speculative Decoding）成熟到 SSM 难追
+ Mamba 的工程优势在"超长序列单纯吞吐"，但实际应用中 attention 的灵活性更关键

> **结论**：SSM/Mamba 是 Transformer 时代的"暗线候选"——短期不会颠覆，但长期可能在长序列、设备端、特定模态上占据一席之地。AI 架构师不应只懂 Transformer。
>

### 5.9 OpenAI 政变（2023.11）：AI 治理史的关键拐点
**事件回顾**：

+ 2023.11.17，OpenAI 董事会突然解雇 CEO Sam Altman，理由含糊（"沟通不一致"）
+ Greg Brockman 同步辞职以示支持
+ 微软 CEO Satya Nadella 宣布邀请 Altman 加入微软成立新 AI 研究部门
+ 数百名 OpenAI 员工签署联名信威胁集体跳槽
+ 11.21，董事会让步，Altman 复职；首席科学家 Ilya Sutskever 离开董事会

**深层影响**：

+ **AI 安全派 vs 加速派的路线之争走向公开化**：Ilya 代表的"对齐优先"派系在 OpenAI 内部失势
+ **2024.5，OpenAI Superalignment 团队解散**，Ilya 离职创立 SSI（Safe Superintelligence Inc.）
+ Jan Leike 等关键研究者出走至 Anthropic——**Anthropic 一夜成为"安全派"主流家园**
+ **资本与人才再分布**：Anthropic、Mistral、xAI、SSI 估值跃升，OpenAI 不再是唯一的"前沿研究中心"
+ **战略意义**：这是 AI 治理结构问题第一次以如此戏剧化方式呈现给公众——**模型越强，谁来决定怎么部署，成为不可回避的政治问题**

### 5.10 多模态的全面融合（2024）
**底层架构革命：DiT 与 Flow Matching**：

+ **DiT（Diffusion Transformer，Peebles & Xie 2023.2）**：把 Diffusion 模型的 U-Net 主干换成 Transformer——**Sora、Stable Diffusion 3、Flux、可灵等 2024 起新一代生成模型的统一架构**
+ **Flow Matching / Rectified Flow（2022—2023）**：替代 DDPM 的训练目标，采样步数大幅减少（几十步 → 4-8 步），SD3、Flux 全面切换
+ **意义**：图像/视频生成走出"U-Net + DDPM"的旧范式，与 LLM 共享"Transformer + 大规模训练"的方法论，**多模态在底层架构上首次实现统一**

**2024.5，GPT-4o（Omni）**：

+ 单模型原生处理文本、音频、图像、视频
+ 端到端语音对话延迟降到 320ms（接近人类对话节奏）
+ 不再是"语音 → 文本 → LLM → 文本 → 语音"的三段式拼接

**2024.2，Sora（OpenAI，demo 公开 / 2024.12 Sora Turbo 正式上线 ChatGPT Plus）**：

+ 文生视频，60 秒 1080p 连贯视频，物理一致性显著优于 Runway / Pika
+ 引爆"AI 视频"赛道，但也暴露了"幻觉"在视频生成中的放大效应（违反物理常识的镜头）

**视觉生成 / 世界模型**：

+ **Veo 2（Google，2024.12）/ Veo 3（2025）**：视频质量超越 Sora 1
+ **可灵（快手）/ Vidu（生数）/ 通义万相（阿里）/ MiniMax Hailuo / 腾讯混元视频（中国，2024—2025）**：视频生成百花齐放，**国产模型在 2025 年事实上与海外平分秋色**
+ **世界模型（World Model）概念兴起**：把视频生成视为对物理世界的隐式建模——LeCun 提出 **JEPA / V-JEPA（2024）** 作为非生成式替代路线，主张"预测嵌入而非像素"
+ **可控视频生成**：Runway Gen-3、Pika 1.5、字节即梦在控制（首尾帧、相机运动、角色一致性）方向卷起新一轮竞争

**音频与语音**：

+ **Suno v3 / v4（2024）、Udio（2024）**：文生音乐，AI 音乐进入大众消费级
+ **ElevenLabs、MiniMax T2A、字节豆包语音**：高质量语音克隆与表达
+ **GPT-4o / Gemini Live / 豆包实时语音**：端到端语音对话，把"语音 Agent"从概念变成产品

### 5.11 RAG 时代：连接 LLM 与企业知识
**RAG 的核心矛盾**：

+ LLM 训练数据有截止时间，且不包含企业私有数据
+ 全量微调成本高、不实时、易遗忘
+ **解法**：检索 + 生成——把外部知识在推理时拼接到 Prompt

**RAG 工程链条**：

```plain
原始文档（PDF/Word/网页/数据库）
    ↓ 切分（Chunking：固定大小 / 语义切分 / 层级切分）
文档片段（Chunks）
    ↓ Embedding（OpenAI text-embedding-3 / BGE / GTE）
向量（Vectors）
    ↓ 存储
向量数据库（Milvus / Qdrant / Weaviate / PGVector）
    ↓
[检索阶段]
用户查询 → Embedding → ANN 检索（HNSW/IVF）→ Top-K 片段
    ↓ Rerank（Cohere Rerank / BGE-Reranker，二次精排）
最相关 K 个片段
    ↓ 拼接到 Prompt
LLM 生成答案 + 引用来源
```

**RAG 的演进**：

+ **Naive RAG（2023）**：检索 + 拼接，简单粗暴，质量参差
+ **Advanced RAG（2024）**：查询改写、Hybrid Search（向量 + 关键词 BM25）、Rerank、上下文压缩
+ **Agentic RAG（2024+）**：让 Agent 决定何时检索、检索什么、是否要继续检索——RAG 成为 Agent 的工具之一，而不是固定流水线
+ **GraphRAG（Microsoft，2024.7）**：用知识图谱组织检索内容，解决多跳推理和全局问题

### 5.12 AI 编程的产品化
**2021.6，GitHub Copilot**：基于 Codex（GPT-3 微调），代码自动补全，是第一个大规模商业化 AI 编程产品。

**2023.3，Cursor 发布**：基于 GPT-4 的 IDE，把"AI 优先"融入编辑器交互，**打破了 VS Code 的垄断**——这是 AI 时代第一个真正颠覆传统工具的产品。

**2024，Claude 3.5 Sonnet 把 Cursor 推上高峰**：

+ Claude 3.5 在 SWE-Bench（GitHub 真实 Issue 修复）上首次达到 49%
+ Cursor 基于 Claude 切换默认模型后，开发者口碑爆发
+ 2024.8 Cursor 估值 25 亿美元，2025 年估值破百亿

**2024.3，Devin（Cognition AI）**："首个 AI 软件工程师"概念提出——能自主完成 GitHub Issue → 写代码 → 跑测试 → 提 PR 的完整流程。引发对"程序员是否会被取代"的全球讨论，但实际产品稳定性 2025 年才成熟。

### 5.13 纪元总结
```plain
核心矛盾：通用能力 vs 商业化落地 / 单模态 vs 多模态融合 / 闭源前沿 vs 开源追赶
核心理论：MoE 稀疏激活、长上下文工程（RoPE/YaRN/Ring Attn）、多模态对齐、DiT/Flow Matching、
        DPO 简化对齐、SSM/Mamba 替代路线
核心技术：GPT-4 / Claude 3.5 / Gemini 1.5 / LLaMA 3 / Mistral / DeepSeek-V2/V3、
        Stable Diffusion / Sora、向量数据库、RAG 工程、
        FlashAttention / vLLM / 量化 / MLA / Speculative Decoding
架构范式：「闭源前沿 + 开源追赶」并行；「单一模型 + 多模态原生」融合
中国里程碑：百模大战（2023）、DeepSeek-V2 价格战（2024.5）、DeepSeek-V3 以 1/10 成本追平 GPT-4（2024.12）、
          Kimi 长上下文创新、Qwen 系列成为全球开源主力
治理事件：2023.11 OpenAI 政变 → AI 安全 / 加速派分裂；Anthropic 成为安全派主流家园
留给后世：「ChatGPT 是产品工程胜利，不只是模型胜利；商业化的差距比技术差距更难追；
        没有 FlashAttention / vLLM / DPO 等工程地基，这一代 LLM 跑不起来」
```

---

## 第六纪元：推理模型与智能体元年（2024.9—2025）
### "模型学会思考，Agent 从概念走进生产"
### 6.1 时代背景与核心矛盾
**触发事件**：2024.9，OpenAI o1 发布——首个商用推理模型，把"测试时计算（Test-Time Compute）"作为新的 Scaling 维度  
**业务特征**：Agent 产品全面爆发（Claude Code、Cursor Composer、Devin、Manus、Replit Agent），从 demo 走向生产  
**算力规模**：训练成本继续指数增长，但**推理时算力（思考时间）**成为新的成本中心  
**核心矛盾**：**模型能力达标 vs 实际任务交付不可靠**——单次能力强不等于多步任务能完成；Agent 落地的瓶颈从模型转向工程

### 6.2 为什么是 2024 年？三件事同时撞上
理解 2024 年的转向，必须理解三件事**几乎同时发生**：

1. **Pre-train Scaling 边际收益开始下降**：从 GPT-3 → GPT-4 是 10× 算力换 ~30% 能力提升；从 GPT-4 → "GPT-4.5 / Orion" 出现首次"训练投入翻倍但能力进步不明显"的迹象（OpenAI 内部 2024 中开始显现，2025 初被多家媒体报道）。
2. **数据墙（Data Wall）问题**：Villalobos et al.《Will we run out of data?》（2022 / 更新到 2024）估计**互联网高质量公开文本约 2026—2028 年用尽**；Common Crawl 后续大量是 AI 生成的"模型回声"内容，质量持续下降。这迫使产业转向：
    - 高质量精炼数据（Phi 系列证明"教科书级数据 < 普通数据 + 体量"反例）
    - 合成数据 + 自我改进（Self-Reward、Self-Play）
    - **可验证 Reward 上做 RL（数学/代码）→ 不依赖人类数据**
3. **RL 训练成熟到可在 LLM 上稳定跑**：GRPO / RLVR / 长 CoT 训练 stack 在 2024 中成熟，让"用 RL 教模型思考"成为可工程化路径。

**这三件事的因果汇流**：

```plain
Pre-train Scaling 收益递减
        +
互联网数据用尽
        +
RL on LLM 成熟
        ↓
转向 Test-Time Scaling（思考时间换性能）
        ↓
o1 / R1 推理模型范式
```

**为什么是"必然"而不是"OpenAI 灵光一现"**：当数据耗尽时，唯一的算力扩张路径就是**让模型在推理时多算**——这与 AlphaGo 时代 MCTS 的逻辑同构：搜索深度可以无限增加，监督信号是稀疏但**可验证**的。

### 6.3 推理模型革命：Scaling 的第三个维度
**Pre-训练 Scaling（2020—2024）**：参数量 N + 数据量 D + 训练算力 C  
**Test-Time Scaling（2024.9—）**：在**推理时**用更多计算（思考链长度 / 多次采样 / 树搜索）换取更高质量

**2024.9，OpenAI o1**：

+ 用强化学习训练模型生成长思维链（CoT），生成时显式"思考"
+ 数学奥赛 AIME 2024 上从 GPT-4o 的 ~13% 跃升至 o1 的 83.3%
+ 物理 PhD 题准确率超过博士生平均水平
+ **代价**：单次推理时间从秒级变成分钟级，成本是 GPT-4 的 6—60 倍

**2024.12，OpenAI o3（公布 / 2025.1—4 分阶段开放）**：

+ AIME 96.7%，**ARC-AGI-1 上 high-compute 模式达到 87.5%（首次跨过 85% 人类阈值），但单题成本破千美元**
+ **重要警示**：在 2025 年新发布的 ARC-AGI-2 上，o3 分数大幅回落至个位数——揭示了"接近人类水平"的脆弱性，单一基准的突破不等于通用能力达标

**2025.1，DeepSeek-R1（开源）**：

+ **彻底改变全球 AI 格局**
+ 完全开源（MIT 协议），论文披露训练方法（GRPO 算法 + Cold Start SFT + 多阶段 RL）
+ 性能接近 o1，**训练成本极具性价比**：DeepSeek-V3 预训练 GPU 成本约 558 万美元（V3 论文披露），R1 在 V3 基础上做 RL 后训练，整体远低于 OpenAI o1 的数亿美元估算（R1 完整训练成本未公开披露）
+ **冲击波**：
    - 2025.1.27 美股 AI 板块单日蒸发 1 万亿美元市值，NVIDIA 跌 17%
    - 全球开发者 1 周内复现各种 R1 蒸馏版本（Qwen-R1、Llama-R1）
    - 让"美国闭源前沿垄断"叙事破产，**开源再次反超**
+ **方法论启示**：
    - 推理能力可以用 RL 从基础模型激发，不需要大量人工 CoT 标注
    - "Aha Moment"：训练中模型自发学会"等等，我重新检查一下"——元认知行为涌现

**2025，Reasoning 模型成为标配**：

| 模型 | 单位 | 时间 | 特点 |
| --- | --- | --- | --- |
| **GPT-5 / o3-pro** | OpenAI | 2025 | 推理 + 通用一体化 |
| **Claude 4 / 4.5 / Extended Thinking** | Anthropic | 2025 | 推理深度可控（用户选 Quick / Extended） |
| **Gemini 2.5 Thinking** | Google | 2025 | 深度集成搜索 + 推理 |
| **DeepSeek-R1 后续迭代**（R1-0528 等） | DeepSeek | 2025.5— | 持续开源迭代，巩固开源推理模型领头地位 |
| **Qwen3-Reasoning / Hunyuan-T1 / Kimi K1.5** | 阿里 / 腾讯 / Moonshot | 2025 | 国产推理模型矩阵 |


### 6.4 推理模型训练范式的三件套
理解推理模型不能只停在"长 CoT"层面，工程师必须懂训练侧：

| 概念 | 全称 | 含义 | 典型应用 |
| --- | --- | --- | --- |
| **ORM** | Outcome Reward Model | 只用最终答案对错评分 | DeepSeek-R1 主路径，配合 RLVR |
| **PRM** | Process Reward Model | 给中间每一步推理打分 | OpenAI Let's Verify Step-by-Step、Math-Shepherd |
| **RLVR** | RL with Verifiable Rewards | 数学 / 代码 / 形式化任务用程序判分替代 RM | R1-Zero 直接 RL 不需要 SFT 冷启动 |
| **Aha Moment / 自反思涌现** | — | 训练中模型自发学会"等等，让我重新检查"——元认知行为涌现 | R1 论文核心发现 |
| **长 CoT 蒸馏** | — | 把强推理模型的思维链作为 SFT 数据训练小模型 | R1-Distill-Qwen / Llama 系列 |


**关键认知**：

+ **R1 的最大方法论贡献不是 GRPO 算法本身，而是证明了"基础模型 + 可验证 RL → 推理涌现"是可复现的工程路径**——这之前被普遍认为是 OpenAI 的"秘方"。
+ **PRM 路线**（OpenAI 早期主推）成本高、对中间步骤标注依赖大；**ORM + RLVR**（DeepSeek 路线）更工程化，是 2025 年开源界主流。



### 6.5 Agent 架构的三个抽象层级
理解 Agent 工程，必须先理清**三层嵌套**的工程抽象（Prompt ⊂ Context ⊂ Harness）：

| 层级 | 解决核心问题 | 时间尺度 | 典型产物 | 价值规律 |
| --- | --- | --- | --- | --- |
| **Prompt Engineering** | 单次调用怎么让模型听懂意图 | 一次 API 调用 | 系统提示词、Few-shot 示例、CoT 引导 | 模型变强 → 技巧贬值 |
| **Context Engineering** | 每一步该给模型什么信息 | 单 Session 内 | RAG 管线、上下文压缩、记忆注入、动态工具集 | 模型变强 → 需求不减 |
| **Harness Engineering** | 整个多步任务怎么结构化推进、防崩、可治理 | 跨 Session、跨工具、跨人机边界 | Planner-Generator-Evaluator 三角色架构、控制矩阵、回退机制 | 模型变强 → 价值迁移但不消失 |


**核心洞察**：模型能力越强，**底层 Prompt 技巧越无用**（"你是一个专业的XX"已经几乎没用），但**上层 Context 和 Harness 的工程价值越突出**——因为模型能力越强，越能驾驭复杂的工具和长任务，那么"给它什么信息、怎么编排任务"的上层设计就越关键。

### 6.6 单 Agent 的核心循环：ReAct 范式
**2022.10，ReAct（Yao et al.）**：Reasoning + Acting 交替的 Prompt 框架——这是所有现代 Agent 的核心循环基础。

```plain
┌─────────────────────────────────────────┐
│              单 Agent 循环               │
└─────────────────────────────────────────┘

   [用户目标 / Goal]
        ↓
   ┌────────────┐
   │  Reason    │ ← 模型基于当前状态思考下一步
   │ (思考)     │
   └────────────┘
        ↓
   ┌────────────┐
   │  Act       │ ← 决定调用哪个工具，传什么参数
   │ (行动)     │
   └────────────┘
        ↓
   ┌────────────┐
   │  Observe   │ ← 工具返回结果（文件内容/搜索结果/命令输出）
   │ (观察)     │
   └────────────┘
        ↓
   [Loop until done]
        ↓
   [完成判断 / 输出]
```

**Function Calling 的标准化**：

+ 2023.6，OpenAI 推出 Function Calling，让 LLM 输出结构化的工具调用 JSON
+ 2024 年 Anthropic、Google、开源模型全部跟进
+ **本质**：把"LLM 调工具"从"用 Prompt 约束输出格式"升级为"模型原生支持的输出模式"

### 6.7 MCP 与 A2A：Agent 互操作的协议层
**2024.11，MCP（Model Context Protocol，Anthropic 主导）**：

+ **核心定位**：AI Agent 的 USB-C 接口——统一 LLM 与外部工具、数据源的协议
+ **核心抽象**：MCP Server（暴露工具/资源/Prompt 的服务）+ MCP Client（LLM 应用）
+ **解决的问题**：
    - 之前每个 Agent 应用都要重写"调用 GitHub / Slack / Database"的对接代码
    - MCP 让工具方做一次 Server，所有支持 MCP 的 Agent 都能用
+ **2025 年生态爆发**：
    - 截至 2026.6，MCP Server 注册数破万，覆盖几乎所有主流 SaaS、数据库、开发工具
    - OpenAI、Google、Anthropic、Microsoft 全部支持 MCP
    - **MCP 成为事实标准**，类比 Web 时代的 HTTP

**2025.4，A2A（Agent-to-Agent Protocol，Google 联合 Atlassian / Salesforce / SAP / MongoDB 等 50+ 厂商发布）**：

+ **核心定位**：Agent 之间互相通信的协议
+ 与 MCP 的关系存在争议：理论上 MCP 解决"Agent 调工具"、A2A 解决"Agent 调 Agent"，但 Anthropic 后续推出了自家的 **Sub-agent / Agent Skills / MCP-over-Agents** 路线，与 A2A 形成**竞合关系**
+ **战略意义**：多 Agent 协作生态的标准化基础正在形成，但**短期内并存多协议**是行业现实

**MCP + A2A 仍在演化中**——目前是行业共识的两条主线，但终极形态尚未定型。任何团队可以做一个专业 Agent，通过协议被其他系统调用、调用所有工具，**但跨厂商的信任、归责、计费体系还在早期**。

### 6.8 Agent 产品大爆发（2024.10—2025）
**2024.10，Anthropic Computer Use**：

+ Claude 直接控制电脑屏幕（截屏 → 视觉识别 → 鼠标键盘操作）
+ 第一次实现"AI 用人的方式用电脑"——突破"软件必须有 API"的限制

**2025.2，Claude Code（Anthropic 官方编程 Agent，2025.2 Research Preview / 2025.5 GA）**：

+ CLI 形态，深度集成终端、文件、Git、Bash
+ 主推 **Plan / Edit / Test 工作流**，强调"先计划再写代码"
+ **工程亮点**：通过 Hooks、Skills、SubAgents 等机制把 Harness 工程理念产品化

**2025.3，Manus（中国，Monica 团队）**：

+ "通用 Agent" 定位，能完成订机票、做研究、写报告等综合任务
+ 首发当日邀请码黑市炒到上万元
+ 验证了"通用 Agent 在某些垂直场景已商业可用"

**2025，Agent 产品全景**：

| 产品 | 公司 | 定位 | 特点 |
| --- | --- | --- | --- |
| **Claude Code** | Anthropic | 编程 Agent | CLI，深度 Harness 工程 |
| **Cursor / Cursor Composer** | Cursor | IDE Agent | IDE 体验，代码理解最强 |
| **Devin** | Cognition AI | 自主软件工程师 | 完整 PR 流程闭环 |
| **Replit Agent** | Replit | 应用开发 Agent | 一句话生成完整 web app |
| **Manus** | Monica | 通用任务 Agent | 跨领域综合能力 |
| **OpenAI Operator** | OpenAI | Web 浏览器 Agent | 打开网页帮你做事 |
| **ChatGPT Agent / Tasks** | OpenAI | 异步任务 Agent | 长时间后台执行 |


### 6.9 Multi-Agent 系统：协作与分工
**核心理念**：单 Agent 的上下文窗口和注意力是有限的，多个专精 Agent 分工 + 通信，可以处理远超单 Agent 复杂度的任务。

**典型架构模式**：

```plain
            ┌──────────────────┐
            │   Orchestrator   │ ← 主管，分解任务、协调子 Agent
            │   (Planner)      │
            └──────────────────┘
              ↓        ↓        ↓
         ┌──────┐  ┌──────┐  ┌──────┐
         │ Sub  │  │ Sub  │  │ Sub  │
         │Agent │  │Agent │  │Agent │  ← 专精执行
         │ A    │  │ B    │  │ C    │
         └──────┘  └──────┘  └──────┘
              ↓        ↓        ↓
            ┌──────────────────┐
            │   Evaluator      │ ← 验证、纠错、决定是否重做
            │   (Critic)       │
            └──────────────────┘
```

**主流多 Agent 框架**：

| 框架 | 出品 | 特点 | 适用场景 |
| --- | --- | --- | --- |
| **LangGraph** | LangChain | 图状态机驱动，可控性强 | 复杂业务流程编排 |
| **AutoGen** | Microsoft | 对话驱动的多 Agent，研究友好 | 学术原型、Multi-Agent 研究 |
| **CrewAI** | CrewAI | 角色 + 任务 + 流程的简洁抽象 | 中小规模业务 Agent |
| **OpenAI Swarm / Agents SDK** | OpenAI | 轻量、官方支持 | OpenAI 生态内开发 |
| **Anthropic 三角色架构** | Anthropic | Planner-Generator-Evaluator 模式 | 高可靠 Agent 工程 |


**Multi-Agent 的核心挑战（2025 年共识）**：

+ **级联错误**：一个子 Agent 输出错误，下游会放大
+ **协调开销**：Agent 之间通信本身消耗大量 Token
+ **可观测性**：调试一个 5 Agent 协作任务，比调试一个 50 万行单体应用还难
+ **结论**：**单 Agent 能解决的问题，永远优先单 Agent**——Multi-Agent 不是终点而是不得已

### 6.10 Agent 评测体系：没有评测就没有迭代
谈 Agent 不谈评测就是空中楼阁。**2024—2025 年形成的主流 Agent Benchmark**：

| 基准 | 出品 | 测什么 | 行业地位 |
| --- | --- | --- | --- |
| **SWE-Bench / SWE-Bench Verified** | Princeton / Anthropic | 真实 GitHub Issue 自动修复 | **编程 Agent 之王**，Anthropic / OpenAI 必报 |
| **Terminal-Bench** | Stanford / Anthropic | 终端命令行任务 | 衡量 Agent 真实操作系统能力 |
| **GAIA** | Meta / HuggingFace | 通用助理任务（多模态、工具调用） | 通用 Agent 能力黄金标准 |
| **WebArena / VisualWebArena** | CMU | 真实 Web 网站任务 | 浏览器 Agent 标准测试 |
| **OSWorld** | HKUST / Tsinghua | 真实操作系统 GUI 任务 | Computer-Use 类 Agent 评测 |
| **τ-Bench (tau-bench)** | Sierra | 多轮工具调用 + 业务规则遵循 | 客服 / 业务 Agent 评测 |
| **AgentBench** | 清华 | 8 个真实场景多 Agent 评测 | 学术界综合评测 |
| **MLE-Bench / RE-Bench** | OpenAI / METR | AI 自己做 ML 研究 | AI 自动化科研能力前沿 |


**2026 年的评测共识**：

+ **静态评测正在饱和**：SWE-Bench Verified 突破 70%、AIME 突破 95% 后，单一基准已无法区分 SOTA
+ **转向动态 / 长任务评测**：HCAST、RE-Bench 等开始评估"模型连续工作数小时甚至数天"的能力
+ **现实落地评测仍是空白**：用户真实任务的成功率、ROI、可控性，**至今没有标准基准**——这是 Harness 工程的关键痛点

### 6.11 Memory：长程记忆作为独立工程领域
**问题**：上下文窗口再大也是临时的，跨 Session 的记忆无处存放——而真正的 Agent 必须**记住用户偏好、过往交互、组织知识**。

**2024—2026 形成的 Memory 技术栈**：

| 类别 | 代表 | 工程定位 |
| --- | --- | --- |
| **产品级记忆** | ChatGPT Memory（2024.4）、Claude Projects/Memory tool（2024—2025）、Gemini 个性化 | 用户级长期记忆 |
| **开源记忆框架** | MemGPT / Letta、Mem0、Zep、Cognee | 自建 Agent 记忆基础设施 |
| **结构化记忆** | GraphRAG（Microsoft 2024.7）、LightRAG、HippoRAG | 用知识图谱组织长期记忆 |
| **分层记忆模型** | 短期（上下文）+ 中期（向量库）+ 长期（结构化文件 / KG） | Anthropic、Cognition、Letta 都是这套范式 |


**核心范式**：

```plain
[Working Memory] 短期：当前对话上下文
       ↓ 周期压缩 / 蒸馏
[Episodic Memory] 中期：会话历史向量库 + 摘要
       ↓ 周期巩固 / 抽取
[Semantic Memory] 长期：知识图谱 + 结构化文件 + 用户偏好库
```

**Memory 是 2026—2027 年的工程战场**：模型够强、Agent 框架够多、工具够丰富，**长程记忆的可控性和可治理性成为下一个差异化点**——这也是 ChatGPT、Claude、Gemini 在 2025 下半年争相投入的方向。

### 6.12 Context Engineering：上下文窗口的工程艺术
**为什么 Context Engineering 成为独立学科（2024 Anthropic 提出）**：

1. 模型上下文窗口虽扩大到 1M+ Tokens，但**注意力衰减**让模型实际有效利用的远小于窗口
2. **长上下文成本**：1M Tokens 推理一次成本数美元，必须精打细算
3. **上下文污染**：错误信息进入上下文后，会持续影响后续推理直到任务结束

**Context Engineering 的核心实践**：

| 实践 | 解决什么 | 工程手段 |
| --- | --- | --- |
| **Token Budget 管理** | 上下文不被"塞爆" | 优先级队列、TTL、压缩策略 |
| **Sliding Window** | 长任务中保留近期上下文 | 滑动窗口 + 摘要 |
| **Compaction（压缩）** | 接近上下文限制时不打断任务 | LLM 总结历史，保留关键信息（Claude Code 的 /compact） |
| **Context Retrieval** | 仅在需要时拉取信息，不预先全塞入 | RAG + 工具调用 |
| **40% 阈值法则** | 超过 40% 利用率开始幻觉增多 | 监控并触发压缩 |
| **Memory 分层** | 短期 vs 长期记忆分离 | 上下文（短期） + 向量库（中期） + 文件系统（长期） |


### 6.13 Harness Engineering 的初步形成
**Harness Engineering 的诞生背景（2025 年共识）**：  
2025 年 Agent 产品虽然能力达到「可自主工作数小时」，但落地中频繁翻车，根因是：

> 模型能力到了 2025 的水平，但**工程基础设施还停留在 2023 年的"单次对话"时代**
>

**Agent 的五大根本性挑战**（模型本身永远无法解决）：

| 挑战 | 本质 | 为什么模型解决不了 |
| --- | --- | --- |
| **状态持久性** | 跨时间 / Session 记住做过什么 | 模型无状态，上下文窗口有上限 |
| **目标一致性** | 长任务不漂移、不提前宣布完成 | 模型无外部锚点，无法校准"完成"标准 |
| **行动可验证性** | 区分"做了"和"做对了" | 模型自评会自我表扬、误判 |
| **熵增抑制** | 持续产出不累积冗余、不一致 | 模型会复制已有模式，哪怕模式本身劣质 |
| **人机边界** | 明确何时自主、何时交还人类 | 模型无可靠"不确定性自觉" |


**Harness 的核心定义（由 Mitchell Hashimoto、Simon Willison、swyx 等多位社区领袖在 2025 年逐步系统化）**：

> Harness = 让模型能作为 Agent 行动起来的**外循环系统**，包含计划分解、持久状态、工具编排、验证门控、反馈回路、回退机制、人机交接点、审计日志。
>

**类比**：模型是 CPU，Harness 是操作系统；**模型决定上限，Harness 决定底线**。

**2025 年 Harness 工程的关键经验数据**：

+ **同一模型仅更换工具调用格式与 scaffolding，SWE-Bench Verified 编码分数从 6.7% 跳到 68.3%**（Anthropic Engineering Blog 公开实验）
+ **LangChain 仅修改 Harness 不改模型，Terminal-Bench 排名从 Top30 升到 Top5**
+ **结论**：瓶颈不在模型而在 Harness——这是 2026 年初的全行业共识

### 6.14 纪元总结
```plain
核心矛盾：模型能力达标 vs 落地不可靠 / 推理成本爆炸 vs 实时性要求 / 数据墙 vs 持续 Scaling
核心理论：Test-Time Scaling、Reasoning RL（GRPO/RLVR/ORM-PRM）、ReAct、
        三层工程抽象（Prompt/Context/Harness）、Memory 分层、Agent 评测体系
核心技术：o1/o3/R1 推理模型、MCP/A2A 协议、Computer Use、Agent 产品全栈、
        Multi-Agent 框架、Memory 框架（MemGPT/Mem0/GraphRAG）、
        Agent Benchmark（SWE-Bench/Terminal-Bench/GAIA/OSWorld/τ-Bench）
架构范式：「单 Agent 循环 + 工具生态 + 长上下文 + Memory + Harness 外循环」
中国里程碑：DeepSeek-R1 引爆全球（2025.1）、Kimi K1.5 与 K2 接力（2025）、
          Manus 通用 Agent 出圈（2025.3）、字节豆包 / 阿里 Qwen3 / 腾讯混元 T1 推理矩阵
留给后世：「Agent 元年的爆发是模型 + 协议 + Harness 三者同步成熟的结果，缺一不可；
        转向 Test-Time Scaling 不是灵光一现，是数据墙下的工程必然」
```

---

## 第七纪元：Harness 工程与 AgentOS 萌芽（2026 至今）
### "重心从让 Agent 更能干，转向让 Agent 不翻车"
### 7.1 时代背景与核心矛盾
**触发事件**：2025 年 Agent 产品大量落地翻车，行业共识从「追模型能力」转向「建工程地基」  
**业务特征**：企业级 Agent 部署、Multi-Agent 协作生产化、AgentOS 概念走出实验室  
**算力规模**：推理算力反超训练算力，成为新的成本中心；专用推理芯片（Groq、Cerebras、华为昇腾推理优化卡）爆发  
**核心矛盾**：**自主性 vs 可控性**——给 Agent 越多自主权效率越高，但风险也越大；如何在工程上找到这个平衡点

### 7.2 Harness Engineering 的方法论体系
**2026 年共识：Harness 不是一个框架，而是一门工程实践**——类似 DevOps，是文化 + 方法的结合，不是某个具体工具。

**Harness 的六层工程架构**（业界共识的六层抽象，综合 Anthropic / Thoughtworks / LangChain 等多家方法论）：

```plain
┌──────────────────────────────────────────┐
│ L6: Governance Layer（治理层）            │  审计、合规、安全审查、人机交接
├──────────────────────────────────────────┤
│ L5: Observability Layer（可观测层）       │  Trace、Metric、Log、Eval、Replay
├──────────────────────────────────────────┤
│ L4: Verification Layer（验证层）          │  TDD、Pre/Post-condition、Critic Agent
├──────────────────────────────────────────┤
│ L3: Orchestration Layer（编排层）         │  Planner、SubAgent 调度、状态机
├──────────────────────────────────────────┤
│ L2: Context Layer（上下文层）             │  Memory、Compaction、RAG、动态工具集
├──────────────────────────────────────────┤
│ L1: Tool Layer（工具层）                  │  MCP Servers、Function Calling、API
└──────────────────────────────────────────┘
                  ↓
              [Model 模型]
```

**核心实践原则**：

1. **Plan First**：长任务先生成显式 Plan，模型 Self-evaluate Plan 质量后再执行——Anthropic 三角色架构的核心
2. **Verify Each Step**：每个工具调用结果都验证（类型检查 / 语义检查 / Critic Agent）
3. **Compact Before 40%**：上下文使用率超过 40% 触发压缩
4. **Idempotent Tools**：工具设计为幂等，重试不破坏状态
5. **Human Checkpoints**：高风险操作（删除、付款、生产环境变更）强制人工确认
6. **Audit Trail**：所有 Agent 决策可追溯、可回放、可问责

### 7.3 工程化落地的成熟路径
不需要一开始搭全量系统，按成熟度渐进：

| 阶段 | 成熟度 | 核心实践 | 适合团队 |
| --- | --- | --- | --- |
| **阶段 1：单 Agent + OpenSpec（规范先行）** | 入门 | 先写需求规格再写代码，需求文档化、接口契约化 | 小团队，单一场景试点 |
| **阶段 2：多 Agent + Superpowers（纪律加固）** | 中级 | 引入 TDD、代码审查、PR 流程等技能约束，强制质量门禁 | 中型团队，多场景并行 |
| **阶段 3：Agent 团队 + Harness（工业化协作）** | 高级 | 定义角色拓扑、并发调度、CI/CD 集成、可观测性平台 | 大型组织，企业级部署 |


**阶段误判的代价**：很多团队 2025 年试图直接跳到阶段 3，结果 Harness 工程量超过业务收益——**过度工程**也是 2026 年要警惕的反模式。

### 7.4 AgentOS：从用户态到内核态的下一跃迁
**研究起点（2024）**：Rutgers 大学 Yongfeng Zhang 团队发表《AIOS: LLM Agent Operating System》——把 Agent 系统类比操作系统，提出资源调度、上下文管理、隔离、日志的系统级抽象。同期国内（清华 / 上交 / 阿里 / 字节）也出现 AgentScope、Agent Hospital、AgentSims 等相关探索。

**2025—2026，操作系统社区接纳 Agent 议题**：OSDI、ASPLOS、SOSP 等系统顶会陆续出现 Agent / LLM serving 相关 Workshop 与论文，标志学术界开始将其作为系统问题对待。

**AgentOS 与 Harness 的分层关系**：

```plain
┌────────────────────────────────────────┐
│      Application（业务 Agent）          │ ← 业务团队
├────────────────────────────────────────┤
│      Harness（用户态层）                │ ← 应用工程师
│   任务分解 / 状态续航 / 验证反馈        │
├────────────────────────────────────────┤
│      AgentOS（内核态层）                │ ← 平台工程师
│   调度 / 隔离 / 资源 / 安全 / 审计      │
├────────────────────────────────────────┤
│      Model + Hardware                   │
└────────────────────────────────────────┘
```

**AgentOS 解决的核心系统级问题**：

+ **生命周期管理**：Agent 创建、休眠、唤醒、销毁
+ **上下文调度**：多 Agent 共享上下文窗口的预算分配（类比内存调度）
+ **隔离**：Agent A 的工具调用不能误伤 Agent B 的状态
+ **合规审计**：所有跨 Agent / 跨工具的调用统一记录、可审计
+ **资源配额**：Token 预算、API 调用预算、推理算力预算

### 7.5 工程师角色的彻底转型
**2026 年的核心趋势**：工程师从"代码生产者"升级为"自治系统设计者"——核心能力从"写出对的代码"变成"设计能让 AI 自动写对代码的环境"。

```plain
传统工程师工作模式：
   人写代码 → 人测试 → 人发布

Agent 早期模式（2024—2025）："In the Loop"
   AI 写代码 → 人改 → 人测试 → 人发布
                    ↑
                    手动介入每个 AI 产物

Harness 成熟模式（2026+）："On the Loop"
   AI 写代码 → AI 测试 → AI 修复 → 人审批关键节点
        ↑
        改 Harness（让系统下次自动做好）
```

**核心能力变化**：

+ **意图定义**：把模糊业务需求结构化为 Spec（规格化思维）
+ **环境设计**：设计 Agent 的工具集、知识库、约束条件
+ **反馈回路构建**：怎么让 Agent 学会"做对了"和"做错了"
+ **可治理性思维**：从一开始就考虑可观测、可回滚、可审计

### 7.6 Harnessability：系统的可驯化性成为核心指标
**2026 年新名词**：**Harnessability（可驯化性）** ——一个系统让 Agent 高效落地的难易程度。

**高 Harnessability 的系统特征**：

+ 强类型（TypeScript、Rust、Pydantic 而非纯动态语言）
+ 测试完备（Agent 改代码后能自动验证）
+ 边界清晰（模块化、接口契约清晰）
+ 文档版本化（与代码同步，机器可读）
+ 运行时可观测（结构化日志、Trace、Metric）

**低 Harnessability 的系统特征（典型棕地项目）**：

+ 知识散落在人脑、Slack、邮件、过时 Wiki
+ 弱类型 / 无 Schema / 接口随意
+ 测试覆盖率低
+ 大量隐式约定（"这个函数不能在周五调用"）

**核心结论**：

+ **绿地项目（从零开始）**：Harness 方法论已经比较成熟，落地顺利
+ **棕地项目（已有多年历史的技术债代码库）**：**Harness 改造是当前行业最大的空白**——也是 2026—2028 年最大的工程机会
+ 2026 年起，"评估系统的 Harnessability"成为架构 Review 的标准环节，类似 2010s 的"评估可扩展性"

### 7.7 安全模型的范式升级：从数据泄露到 Agency 操控
**传统安全威胁**：SQL 注入、XSS、数据泄露、DDoS——攻击目标是"数据"

**Agent 时代新威胁（2025—2026 集中爆发）**：

+ **MCP 工具投毒**：恶意 MCP Server 在工具描述中嵌入指令，诱导 Agent 做出危险操作
+ **跨工具数据外流**：Agent 在工具 A 读到敏感数据，被诱导调用工具 B 把数据发出去
+ **Prompt Injection 升级版**：网页内容、文件中嵌入隐藏指令劫持 Agent
+ **Agent 越权链**：Agent A 调用 Agent B，Agent B 滥用 A 的权限
+ **攻击目标变化**：从"窃取数据"转向"操控 Agency（行动权）"——让 Agent 替攻击者做事

**Harness 内置的安全机制（2026 行业标准）**：

+ **动态权限**：Agent 在每个时刻只持有最小必要权限，任务结束立即回收
+ **执行隔离**：Agent 调用工具在沙箱（容器、子进程、WASM）中执行
+ **人类审批介入点**：高风险操作（金钱、删除、生产变更）强制 HITL
+ **意图验证**：Agent 决策前 Critic Agent 评估"这个操作是否符合用户原始意图"
+ **审计回放**：所有动作可被完整回放、归因

### 7.8 AI 安全的独立技术线（2024—2026）
AI 安全已经从"科幻话题"变成可工程化的研究领域，主要分四条主线：

**1. 模型能力分级与发布前评估（Frontier Safety）**：

+ **Anthropic Responsible Scaling Policy（RSP，2023.9 / 2024 ASL-3）**：把模型按"灾难性风险能力"分 ASL-1~5 级，触发某级前必须通过对应安全评估
+ **OpenAI Preparedness Framework（2023.12）**：四类风险评分（Cybersecurity / CBRN / Persuasion / Model Autonomy）
+ **DeepMind Frontier Safety Framework（2024.5）**：定义关键能力阈值（CCL）
+ **独立第三方评测**：**METR（前 ARC Evals）** 评估 OpenAI / Anthropic 旗舰模型的自主复制 / 漂移风险，**事实上的"AI 安全审计公司"**

**2. 红队与对抗（Adversarial Robustness）**：

+ **Jailbreak 演进**：从 DAN / Grandma exploit → 多轮诱导 → 多模态注入（图像内嵌指令）→ 自动化攻击（PAIR、TAP）
+ **Prompt Injection 的工程级威胁**：Simon Willison 等 2023 起的反复警告在 2025 年成为现实——**间接 Prompt Injection（网页 / 邮件 / 文件中的恶意指令）**已是 Agent 时代头号攻击向量
+ **MCP 工具投毒、RAG 文档投毒**：成为 2025—2026 攻防焦点

**3. 可解释性（Mechanistic Interpretability）**：

+ **Anthropic Circuits 系列（2022—）**：从单个神经元到电路级理解模型计算
+ **Sparse Autoencoders（SAE，Anthropic 2024.5）"Mapping the Mind of Claude"**：在 Claude 3 Sonnet 中提取出可识别概念（金门大桥、bug、性别偏见等），**可解释性首次在前沿模型上规模化**
+ **意义**：从"黑盒输出 + 行为评测"转向"内部机制可观测"——这是 AI 安全的根本性升级

**4. 治理与监管**：

+ **EU AI Act（2024.3 通过 / 2025—2027 分阶段生效）**：全球首个综合性 AI 法规
+ **美国行政令（2023.10）/ 加州 SB-1047（2024 否决 / 2025 重提）/ AI Safety Institutes（UK / US / 日 / 中等多国设立）**
+ **2025—2026 的核心争议**：开源前沿模型权重是否应受管制？模型自主复制 / RSI 实验的红线在哪？这些尚无国际共识。

> **AI 安全已从"少数派理想主义"变成"前沿实验室的强制工程纪律"**——任何团队部署 Agent，都必须把红队评估、可解释性、人类监督作为标配。
>

### 7.9 与 Harness 并行的 2026 上半年主线
Harness 不是 2026 年的全部，至少还有四条同样关键的并行主线：

**1. 小模型革命（SLM, Small Language Models）**：

+ **Phi-3 / Phi-4（Microsoft，2024—2025）**：用"教科书级数据"训练，4B 参数能力接近 70B
+ **Gemma 2 / Gemma 3（Google）/ Qwen3-1.5B-4B / DeepSeek-V3-mini-distill**：开源小模型百花齐放
+ **Apple Intelligence（2024.6 公布 / 2025—2026 落地）**：3B 端侧模型 + 云端 PCC（Private Cloud Compute）混合架构，**让 AI 走入每一台 iPhone**
+ **意义**：不是所有任务都需要前沿大模型——**80% 的实际任务用 SLM 端侧 / 边缘部署性价比更高**

**2. 垂直 Agent 商业化的真实进展**：

+ **跑通的**：编程（Cursor、Claude Code）、客服（Sierra、Decagon）、销售（11x、Artisan）、研究（Glean、Perplexity Pro）
+ **半跑通的**：法律（Harvey）、医疗（开源 Health Agent）、设计（Figma AI）
+ **没跑通的**：完全自主的"AI 程序员"（Devin 仍以 Copilot 模式为主）、通用任务 Agent（Manus 等仍在演化）
+ **核心规律**：**领域知识 + 数据闭环 + 工程基础设施** > 单纯的模型能力

**3. AGI / ASI 时间表的公开分歧**：

+ **Dario Amodei（Anthropic）**：2026—2027 出现"国家级智能"（Powerful AI / Country of Geniuses）
+ **Sam Altman（OpenAI）**：2025 年内 Agent 可工作几小时；AGI 在"几年内"
+ **Demis Hassabis（DeepMind）**：5—10 年到 AGI，2030 是合理上限
+ **Yann LeCun（Meta）**：现有 LLM 路径是死路，需要 JEPA / 世界模型新范式才能到 AGI，可能需要数十年
+ **Yoshua Bengio**：相信短时间表，但更担心安全失控
+ **行业现实**：**没有共识。任何引用单一时间表的报告都应被打折扣**——这本身就是 AI 时代认知的重要部分

**4. 算力地缘格局**：

+ **NVIDIA Blackwell B200 / GB300 / Rubin（2024—2026）**：训练算力继续指数级
+ **Google TPU v5/v6、AWS Trainium 2、Meta MTIA**：超大规模公司自研芯片崛起
+ **国产算力**：华为昇腾 910C/910D、寒武纪、摩尔线程、燧原、海光——美国制裁倒逼的"备胎转正"，2025—2026 是关键决战年
+ **专用推理芯片**：Groq LPU、Cerebras WSE、SambaNova、Etched ASIC——为推理时代特化
+ **能源墙**：千卡集群单次训练耗电相当于一个小镇全年——**2026 年起电力比 GPU 更稀缺**，AI 数据中心的选址回归"挨着核电站"

### 7.10 当前未解决的核心问题（2026.6 行业无共识）
1. **怎么验证 Agent 真的「做对了事」**：用 AI 生成的测试验证 AI 生成的代码，本质是"用同一双眼睛检查自己作业"，可靠性存疑。
2. **AI 生成代码的长期可维护性**：LLM 倾向重新实现已有功能（"AI Slop"），长期运行的代码库熵增问题尚无成熟解决方案。
3. **Harness 的 ROI 度量**：怎么根据任务复杂度动态调整 Harness 深度，避免过度工程（Harness 开销超过质量提升）。
4. **单 Agent vs 多 Agent 的边界**：小项目单 Agent 足够，大项目需要多 Agent 分工，但具体规模阈值无明确标准。
5. **Agent 之间的信任传递**：A2A 协议解决了通信，但跨组织、跨厂商的 Agent 互相调用如何建立信任、归责？
6. **AGI 安全门槛**：随着 RSI（Recursive Self-Improvement）实验从研究走向小规模应用，怎么保证不失控？

### 7.11 纪元总结
```plain
核心矛盾：自主性 vs 可控性 / Agent 能力 vs 工程基础设施 / 绿地方法论 vs 棕地改造 /
        前沿大模型 vs 端侧小模型 / 闭源能力 vs 开源透明 / 加速派 vs 安全派
核心理论：Harness Engineering 六层架构、AgentOS、Harnessability、Agency 安全模型、
        Frontier Safety 分级、Mechanistic Interpretability、SLM 路线
核心技术：MCP/A2A 生态、Multi-Agent 编排、Continuous Verification、Sandbox 执行环境、
        SAE 可解释性、RSP/Preparedness 评估、端侧 SLM、专用推理芯片
架构范式：「Model + Harness」组合评估，Agent 作为一等公民设计系统；
        前沿大模型 + 端侧小模型 + 私有推理三层并存
中国里程碑：国产 Agent 产品矩阵成熟（豆包 / Qwen / 混元 / 智谱 / Kimi）、
          AgentOS 学术与工业并行（清华 / 上交 / 阿里 / 字节）、
          国产推理芯片备胎转正（昇腾 / 寒武纪 / 摩尔线程）
留给未来：「Harnessability 决定一个组织在 AI 时代的上限——不是模型能力，而是系统能不能被 AI 高效驯服；
        但同时也别忘了：AI 4S 让 AI 不只是生产力工具，而是文明加速器」
```

---

## 第八章：AI 架构师的终极认知框架
### 8.1 七个纪元的核心矛盾对照表
| 纪元 | 时间 | 业务规模 | 核心矛盾 | 理论突破 | 架构范式 |
| --- | --- | --- | --- | --- | --- |
| **符号主义** | 1950s—1990s | 实验室级 | 规则爆炸 vs 现实复杂性 | 谓词逻辑、推理引擎 | KB + IE |
| **统计机器学习** | 1990s—2011 | 互联网级 | 过拟合 vs 容量 | VC 维、SRM、集成学习 | 特征工程 + 浅层模型 |
| **深度学习革命** | 2012—2017 | 移动互联网级 | 手工特征 vs 端到端 | 反向传播、ResNet、Attention | CNN/RNN → Transformer |
| **预训练 + Scaling** | 2018—2022 | 全互联网级 | 单任务 vs 通用基座 | Scaling Law、RLHF、涌现 | 大模型预训练 + 微调 |
| **生成式 AI 爆发** | 2022.11—2024 | 全民级 | 闭源前沿 vs 开源追赶 | MoE、长上下文、多模态融合 | 闭源 + 开源并行 |
| **推理 + 智能体元年** | 2024.9—2025 | Agent 化 | 能力达标 vs 落地不可靠 | Test-Time Scaling、ReAct、MCP | 单/多 Agent + 工具生态 |
| **Harness + AgentOS** | 2026 至今 | 系统化 | 自主性 vs 可控性 | Harness、Harnessability、AgentOS | Model + Harness 组合 |


### 8.2 贯穿始终的三条技术主线
**主线一：智能从何而来——理论范式的演进**

```plain
符号逻辑（1950s）：智能 = 知识 + 推理
    ↓ 寒冬证明：现实规则无法穷举
统计学习（1990s）：智能 = 数据驱动的概率推理
    ↓ 浅层模型容量天花板
表征学习（2012）：智能 = 自动学习层次化表征
    ↓ Transformer 统一所有模态
预训练 Scaling（2018）：智能 = 大数据 + 大模型 + 大算力的涌现
    ↓ 模型能力强但不"听话"
对齐 RLHF（2022）：智能 = 能力 + 对齐
    ↓ 单次回答好不等于多步任务做对
推理 + Agent（2024）：智能 = 模型 + 工具 + 多步循环
    ↓ Agent 能跑但翻车
Harness（2026）：智能 = 模型 + 外循环工程系统
```

**主线二：架构抽象层级的持续上升**

```plain
神经元（生物学启发）
    → 网络结构（CNN/RNN/Transformer）
    → 模型（GPT/Claude）
    → 模型 + Prompt（API 应用）
    → Agent 循环（ReAct）
    → Multi-Agent 协作
    → Harness 系统（外循环）
    → AgentOS（内核层）
    → 自主智能体生态（MCP + A2A 标准化）
```

**主线三：工程实践的三层演进**

```plain
Prompt Engineering（2020—2023）：手工调提示词
   ⊂
Context Engineering（2023—2025）：每一步给模型什么信息
   ⊂
Harness Engineering（2025—）：整个任务怎么结构化推进、防崩、可治理
```

**记忆口诀**：`Prompt ⊂ Context ⊂ Harness`——三者是嵌套关系，不是替代关系。模型越强，**底层 Prompt 技巧越贬值，上层 Harness 价值越突出**。

### 8.3 Agent 系统设计决策框架：「四问一答」
面对 Agent 设计题，按以下框架作答：

```plain
第一问：任务的复杂度和时间尺度？
  → 单步问答（用 RAG 即可）
  → 多步任务（用 Agent 循环）
  → 长任务多 Session（必须做 Harness 状态持久化）

第二问：自主性需要多高？
  → 高（少打扰人）：On-the-Loop 模式 + 强 Critic + Audit
  → 中：HITL 关键节点确认
  → 低（每步都人审）：Copilot 模式

第三问：核心风险点在哪？
  → 行动可逆 → 自主性可以高
  → 行动不可逆（删除、付款、生产变更）→ 必须 HITL
  → 涉及隐私 / 合规 → 必须沙箱 + 审计

第四问：上下文压力？
  → 短任务（<30K Tokens）：直接做
  → 中等（30K—200K）：滑动窗口 + 摘要
  → 长任务（>200K）：必须做 Compaction、外存记忆、子任务交接

一答：给出具体方案：
  → 用什么模型（Reasoning vs Fast）
  → 单 Agent 还是 Multi-Agent
  → Context 策略（Memory 分层、压缩策略）
  → Harness 配置（Plan/Verify/Compact/HITL/Audit 各层是否启用）
  → 评估指标（成功率、Token 成本、人工介入率、回滚率）
```

### 8.4 模型选型决策树（2026 视角）
```plain
任务类型？
  │
  ├─ 1. 高频简单任务（分类、抽取、改写）
  │     → 小模型蒸馏（Qwen 7B / LLaMA 8B / Phi）
  │     → 成本极低，延迟毫秒级
  │
  ├─ 2. 通用对话 / 创作
  │     → Claude Sonnet / GPT-4o / Gemini Flash 级
  │     → 性价比好，能力够用
  │
  ├─ 3. 复杂推理（数学、科研、复杂代码）
  │     → o3 / Claude Opus + Extended Thinking / DeepSeek-R1 / Gemini Thinking
  │     → 单次成本高但准确率突破式提升
  │
  ├─ 4. 长上下文（整本书 / 大型代码库）
  │     → Gemini 1.5/2.5 Pro（2M+）/ Claude Opus 4.7（1M）
  │     → 注意上下文衰减，配合 Compaction
  │
  ├─ 5. 多模态（图文混合 / 视频）
  │     → GPT-4o / Gemini 2.5 / Claude 3.5+
  │
  ├─ 6. Agent 任务（多步工具调用）
  │     → Claude Sonnet（编程 Agent 之王）/ GPT-4o（通用）
  │     → 需要稳定的 Function Calling + 长指令跟随
  │
  └─ 7. 数据敏感 / 合规要求高
        → 私有化部署：Qwen3 / DeepSeek / LLaMA 4
        → 国产可控：豆包 / 混元 / 文心
```

### 8.5 RAG vs Fine-tuning vs Long Context：三大范式选型
| 维度 | RAG | Fine-tuning | Long Context |
| --- | --- | --- | --- |
| **新知识引入** | ✅ 实时 | ❌ 训练时一次性 | ✅ 推理时拼接 |
| **数据规模** | TB 级 | GB—TB 级 | 100K—2M Tokens |
| **延迟** | 中（检索 + 推理） | 低（直接推理） | 高（长上下文推理慢且贵） |
| **成本** | 检索 + 短推理 | 训练贵，推理便宜 | 单次推理极贵 |
| **可解释性** | 高（有引用来源） | 低（知识隐式编码） | 中（在上下文中可定位） |
| **适用场景** | 企业知识库、问答、客服 | 风格定制、特定 schema 输出 | 整本书分析、整库审计 |
| **2026 推荐** | **首选**：覆盖 80% 企业场景 | 仅在风格 / 输出格式严格时使用 | 仅在确实需要全局理解时使用 |


**反模式**：很多团队 2024—2025 年盲目尝试 Fine-tuning，结果发现 RAG + 强模型更好——因为模型能力越强，**用 Prompt + RAG 注入知识比改参数更高效**。

### 8.6 终极四句话
1. **符号主义告诉你：把所有规则写死，永远撞不出智能。**——AI 寒冬两次的根本教训。
2. **Scaling Law 告诉你：堆数据 + 堆参数 + 堆算力，能让"涌现"自然发生。**——但堆是有上限的（数据墙、能源墙），**对齐和数据质量同等重要**。
3. **Harness Engineering 告诉你：模型决定上限，Harness 决定底线。**——Agent 时代的胜负手不在模型本身，而在围绕模型的工程系统。
4. **AI for Science 告诉你：当模型能压缩自然规律，AI 就从生产力工具升级为文明加速器。**——AlphaFold 与 2024 年双诺奖把 AI 重新定义为基础科学的核心方法。

**最顶层的认知**：所有 AI 范式、模型架构、Agent 工程，都在回答同一个问题——在数据有限、算力有限、对齐困难、世界复杂的现实面前，如何让一个学习系统**既学得到、又泛化得动、还能在真实环境中可靠行动**。

**架构没有最优解，只有最适配的解。真正的 AI 架构师，知道什么场景必须用大模型，什么场景一行规则就够；知道什么操作可以让 Agent 自主，什么操作必须让人审批；知道什么时候该堆 Token，什么时候该 Compact；知道什么任务该上前沿模型，什么任务一个 4B 端侧模型就够。**

### 8.7 AI 工程化的常见反模式（2026 经验汇总）
| 反模式 | 表现 | 正确做法 |
| --- | --- | --- |
| **模型崇拜** | "我们换最强模型就行" | 90% 问题在 Context 和 Harness，不在模型 |
| **过度工程** | 简单任务也搭 Multi-Agent + 6 层 Harness | 单 Agent 能解决就用单 Agent |
| **盲目 Fine-tune** | 任何定制需求都先想训练 | 先尝试 Prompt + RAG，达到瓶颈再考虑 Fine-tune |
| **塞满上下文** | 把所有可能用到的信息都拼到 Prompt | 40% 阈值法则，超过则触发压缩 |
| **缺失评估** | 凭感觉觉得效果不错就上线 | Eval 数据集 + 持续监控 + Replay |
| **HITL 缺位** | 让 Agent 完全自主操作生产环境 | 高风险操作必须人工审批 |
| **MCP Server 不审计** | 装上别人写的 MCP Server 就用 | 工具调用是攻击面，必须沙箱 + 审计 |
| **多 Agent 万能论** | "拆成多个 Agent 就能解决任何问题" | 多 Agent 引入级联错误，先尝试单 Agent + 强 Harness |


---

## 附录 A：AI 关键论文 / 事件时间线
| 年份 | 论文 / 事件 | 意义 |
| --- | --- | --- |
| 1943 | McCulloch & Pitts《神经元数学模型》 | 人工神经元起点 |
| 1950 | Turing《Computing Machinery and Intelligence》 | 图灵测试，AI 哲学起点 |
| 1956 | Dartmouth 会议 | "AI" 一词正式诞生 |
| 1957 | Rosenblatt 感知机 | 第一个能学习的神经网络 |
| 1969 | Minsky & Papert《Perceptrons》 | 第一次神经网络寒冬的导火索 |
| 1986 | Rumelhart et al.《Back-Propagation》 | 反向传播系统化，神经网络复活 |
| 1989 | LeCun LeNet | CNN 首个工程成功 |
| 1995 | Vapnik SVM | 统计学习时代主力 |
| 1997 | LSTM（Hochreiter & Schmidhuber） | 长程序列建模的关键突破 |
| 2006 | Hinton《Deep Belief Nets》 | "Deep Learning" 一词首次正式使用 |
| 2009 | ImageNet 数据集（李飞飞） | 深度学习引爆的物理基础 |
| 2012 | AlexNet（ImageNet） | 深度学习革命引爆 |
| 2013 | Word2Vec（Mikolov） | 词向量经典 |
| 2014 | GAN（Goodfellow） | 生成对抗网络开山 |
| 2014 | Seq2Seq + Attention（Bahdanau） | 注意力机制起源 |
| 2015 | ResNet（何恺明） | 残差连接，深度学习里程碑 |
| 2016 | AlphaGo 击败李世石 | 公众认知的"超人"AI 元年 |
| 2017 | Vaswani et al.《Attention is All You Need》 | Transformer 诞生 |
| 2018 | BERT（Google） | 预训练 + 微调范式确立 |
| 2018 | GPT-1（OpenAI） | 自回归 LLM 路线起点 |
| 2020 | GPT-3（175B） | Few-shot / In-context Learning 涌现 |
| 2020 | Kaplan et al.《Scaling Laws》 | 把"训大模型"变成工程问题 |
| 2020 | Ho et al.《DDPM》 | 扩散模型理论奠基 |
| 2020.7 | **AlphaFold 2（CASP14）** | **AI 解决蛋白质折叠问题，AI4S 主线开端** |
| 2021 | CLIP / DALL·E（OpenAI） | 视觉-语言对齐 |
| 2021 | RoPE（Su et al.） | 旋转位置编码，长上下文基础 |
| 2022 | Chinchilla（DeepMind） | 数据 / 参数最优配比 |
| 2022 | InstructGPT / RLHF | 对齐范式确立 |
| 2022.5 | **FlashAttention（Tri Dao）** | **现代 LLM 训练的隐形地基** |
| 2022.8 | Stable Diffusion 开源 | 文生图全民化 |
| 2022.10 | ReAct（Yao et al.） | Agent 核心循环 |
| 2022.11 | ChatGPT 上线 | AI 进入大众视野，5 天百万用户 |
| 2022.12 | Constitutional AI（Anthropic） | RLAIF 路线确立 |
| 2023.2 | **DiT（Peebles & Xie）** | **Sora / SD3 等新一代生成模型架构基础** |
| 2023.3 | GPT-4 发布 | 多模态前沿、闭源化 |
| 2023.5 | **DPO（Rafailov et al.）** | **去掉 PPO 的简化对齐方法，工业界事实主流** |
| 2023.6 | **vLLM / PagedAttention** | **开源推理事实标准** |
| 2023.7 | LLaMA 2 开源（商用） | 开源 LLM 生态爆发 |
| 2023.9 | **Anthropic RSP（Responsible Scaling Policy）** | **前沿模型能力分级首次工程化** |
| 2023.11 | **Sam Altman 解雇与复职** | **AI 治理史拐点，安全/加速派分裂** |
| 2023.12 | **Mamba（Gu & Dao）** | **状态空间模型挑战 Transformer** |
| 2023.12 | Mixtral 8×7B 开源 | MoE 普及到全球开发者 |
| 2024.2 | Sora 公开 demo | 文生视频引爆 |
| 2024.5 | **Anthropic SAE / "Mapping the Mind of Claude"** | **可解释性首次在前沿模型规模化** |
| 2024.5 | **DeepSeek-V2 + MLA + 价格战** | **国产 LLM API 价格 1/100，引爆百模降价** |
| 2024.5 | GPT-4o | 端到端多模态 |
| 2024.5 | **AlphaFold 3** | **AI 药物发现进入工程化阶段** |
| 2024.6 | Claude 3.5 Sonnet | 编程能力首次稳超 GPT-4 |
| 2024.7 | **AlphaProof / AlphaGeometry 2 IMO 银牌** | **AI 接近数学奥赛金牌** |
| 2024.9 | OpenAI o1 | 推理模型范式开启 |
| 2024.10 | **2024 年诺贝尔物理学奖（Hinton, Hopfield）+ 化学奖（Hassabis, Jumper, Baker）** | **AI 首次同年拿下双诺奖** |
| 2024.10 | Anthropic Computer Use | Agent 直接控制电脑 |
| 2024.11 | MCP 协议发布（Anthropic） | Agent 工具协议标准化 |
| 2024.12 | DeepSeek-V3 开源（670B MoE） | 国产 MoE 追平 GPT-4 |
| 2025.1 | DeepSeek-R1 开源 | 全球 AI 格局重塑，开源反超 |
| 2025.1 | **Kimi K1.5 推理模型 + 论文披露** | **国产推理路径独立验证** |
| 2025.2 | Claude Code | 编程 Agent 工程化标杆 |
| 2025.3 | Manus（中国） | 通用 Agent 商业化样板 |
| 2025.4 | Qwen3 开源 | 全球开源 LLM 主力家族成形 |
| 2025.4 | A2A 协议（Google） | Agent 间通信标准化 |
| 2025.5 | Anthropic Claude 4 / Sonnet 4 | Agentic 编程基线大幅提升 |
| 2025.7 | Kimi K2 开源（万亿 MoE） | 国产开源 Agent 基座代表 |
| 2025.9 | Claude Sonnet 4.5 | 长任务自主执行能力跃升 |
| 2025—2026 | Harness Engineering / AgentOS | Agent 落地工程方法论成型 |
| 2025—2027 | **EU AI Act 分阶段生效** | **全球首个综合性 AI 法规** |


---

## 附录 B：AI 工程师必读书单 & 必读论文
### 经典书籍（按阅读优先级排序）
| 优先级 | 书名 | 作者 | 核心价值 | 适合阶段 |
| --- | --- | --- | --- | --- |
| ★★★★★ | **《Deep Learning》（花书）** | Goodfellow / Bengio / Courville | 深度学习理论圣经，覆盖基础到 RNN/CNN/生成模型 | 入门→中级 |
| ★★★★★ | **《Pattern Recognition and Machine Learning》（PRML）** | Bishop | 机器学习概率视角的奠基教材 | 中级→高级 |
| ★★★★★ | **《Hands-On Large Language Models》** | Jay Alammar et al. | LLM 时代实战导向教材，2024 出版 | 中级 |
| ★★★★☆ | **《Reinforcement Learning: An Introduction》** | Sutton & Barto | RL 经典教材，理解 RLHF/o1 必读 | 中级→高级 |
| ★★★★☆ | **《Building LLM-Powered Applications》** | 多种版本 | RAG / Agent 工程实战 | 中级 |
| ★★★★☆ | **《AI Engineering》** | Chip Huyen | 生产级 LLM 系统设计，2025 必读 | 中级→高级 |
| ★★★☆☆ | **《Designing Machine Learning Systems》** | Chip Huyen | ML 系统的工程实践 | 中级 |
| ★★★☆☆ | **《The Hundred-Page Machine Learning Book》** | Burkov | 快速建立 ML 全景 | 入门 |


### 必读论文（按时间轴排序）
| 年份 | 论文 | 核心贡献 | 读完能回答什么问题 |
| --- | --- | --- | --- |
| 1986 | Rumelhart et al.《Learning representations by back-propagating errors》 | 反向传播算法 | 神经网络如何训练？ |
| 1997 | Hochreiter & Schmidhuber《LSTM》 | 长程序列建模 | LSTM 的门控机制？ |
| 2012 | Krizhevsky et al.《ImageNet Classification with Deep CNNs》（AlexNet） | 深度学习实证起点 | 为什么 AlexNet 改变了一切？ |
| 2014 | Goodfellow et al.《Generative Adversarial Networks》 | GAN 起源 | GAN 的训练原理与崩溃模式？ |
| 2015 | He et al.《Deep Residual Learning》（ResNet） | 残差连接 | 为什么 152 层网络能训得动？ |
| 2017 | Vaswani et al.《Attention is All You Need》 | Transformer 诞生 | Self-Attention 的本质？多头注意力？ |
| 2018 | Devlin et al.《BERT》 | 双向预训练 | BERT 与 GPT 的本质差异？ |
| 2020 | Brown et al.《GPT-3 / Language Models are Few-Shot Learners》 | In-Context Learning | 为什么大模型不需要微调也能做新任务？ |
| 2020 | Kaplan et al.《Scaling Laws for Neural Language Models》 | Scaling Law | 给定预算，怎么决定模型规模和数据量？ |
| 2022 | Wei et al.《Chain-of-Thought Prompting》 | CoT 范式 | 为什么"Let's think step by step"有效？ |
| 2022 | Hoffmann et al.《Chinchilla》 | 数据/参数最优比 | GPT-3 为什么严重欠训练？ |
| 2022 | Ouyang et al.《InstructGPT》 | RLHF 三阶段 | 怎么对齐 LLM？ |
| 2022 | Yao et al.《ReAct: Synergizing Reasoning and Acting》 | Agent 核心循环 | Agent 的基本工作模式？ |
| 2022 | Dao et al.《FlashAttention》 | Attention IO 优化 | 为什么现代 LLM 训得起来？ |
| 2023 | Kwon et al.《Efficient Memory Management for LLM Serving》（vLLM/PagedAttention） | 推理内存优化 | 怎么提升推理吞吐？ |
| 2023 | Rafailov et al.《Direct Preference Optimization》 | DPO | 不要 PPO 怎么对齐？ |
| 2023 | Peebles & Xie《Scalable Diffusion Models with Transformers》（DiT） | 扩散 + Transformer | Sora / SD3 的架构基础？ |
| 2023 | Gu & Dao《Mamba: Linear-Time Sequence Modeling with Selective State Spaces》 | SSM 路线 | Transformer 之外还有什么？ |
| 2024 | Anthropic《Sparse Autoencoders Find Highly Interpretable Features》 / "Scaling Monosemanticity" | 可解释性突破 | 模型内部到底在算什么？ |
| 2024 | DeepSeek《DeepSeek-V3 Technical Report》 | 高效 MoE 训练 + MLA | 怎么用 1/10 成本训出顶级模型？ |
| 2024 | Jumper et al. / Abramson et al.《AlphaFold 3》 | 蛋白复合物预测 | AI 怎么进入药物发现？ |
| 2025 | DeepSeek《DeepSeek-R1: Incentivizing Reasoning via RL》 | 推理 RL 训练 / GRPO | 怎么从基础模型激发推理能力？ |
| 2024 | Yang et al.《AIOS: LLM Agent Operating System》 | AgentOS 起点 | Agent 系统级调度怎么做？ |
| 2024 | Anthropic《Building effective agents》 | 单 / 多 Agent 工程范式 | 什么时候用 Agent，什么时候不用？ |


> **读论文建议**：Transformer、Scaling Law、InstructGPT、DPO、ReAct、DeepSeek-R1 这六篇是 2026 年 LLM 工程师的"必修六经"——读完能建立从模型架构、对齐、Agent 到推理时代的完整脉络。如有余力，再补 FlashAttention、vLLM、AlphaFold 3 三篇了解工程地基与 AI4S 主线。
>

---

## 附录 C：术语对照表（中英对照）
| 中文 | 英文 | 简释 |
| --- | --- | --- |
| 通用人工智能 | AGI（Artificial General Intelligence） | 与人类水平相当的通用智能 |
| 大语言模型 | LLM（Large Language Model） | 基于 Transformer 的大规模语言模型 |
| 预训练 | Pre-training | 在大规模无标注数据上的自监督训练 |
| 微调 | Fine-tuning | 在小规模标注数据上的二次训练 |
| 监督微调 | SFT（Supervised Fine-Tuning） | 用人类标注的"指令-回答"对微调 |
| 人类反馈强化学习 | RLHF（Reinforcement Learning from Human Feedback） | 用人类偏好训练奖励模型，再 RL 优化 |
| 思维链 | CoT（Chain-of-Thought） | 显式输出推理步骤 |
| 上下文学习 | ICL（In-Context Learning） | 仅靠 Prompt 中的示例学新任务，不更新参数 |
| 涌现能力 | Emergent Abilities | 模型规模到达临界点后突然出现的能力 |
| 检索增强生成 | RAG（Retrieval-Augmented Generation） | 检索外部知识 + LLM 生成 |
| 混合专家 | MoE（Mixture of Experts） | 总参数大但每次只激活部分专家 |
| 测试时计算 | Test-Time Compute | 推理时用更多算力换更高质量 |
| 工具调用 | Function Calling / Tool Use | LLM 输出结构化的工具调用请求 |
| 模型上下文协议 | MCP（Model Context Protocol） | LLM 与工具/数据源的统一接口协议 |
| 智能体间协议 | A2A（Agent-to-Agent Protocol） | Agent 之间通信的协议 |
| 智能体操作系统 | AgentOS | Agent 系统级调度、隔离、审计的内核层 |
| 外循环系统 | Harness | 让模型作为 Agent 行动起来的工程基础设施 |
| 可驯化性 | Harnessability | 系统让 Agent 高效落地的难易程度 |
| 上下文工程 | Context Engineering | 决定每一步给模型什么信息的工程 |
| 提示词工程 | Prompt Engineering | 设计单次调用的提示词的技巧 |
| 人在回路 | HITL（Human-In-The-Loop） | 关键节点人工介入决策 |
| 人在回路上 | HOTL（Human-On-The-Loop） | 人监督但不干预每一步，仅在异常时介入 |
| 人在回路外 | HOOL（Human-Out-Of-The-Loop） | 完全自主，人事后审计 |


---

_文档版本：2026-06（v2 · 经 AI 技术专家团队事实校核与遗漏补全）_  
_覆盖范围：AI 完整发展脉络 · 符号主义→统计学习→深度学习→预训练→生成式爆发→推理与智能体→Harness 工程·AI4S 主线_  
_书写参考：《后端技术演进全景》同款叙事框架 · 业务规模驱动 · 矛盾跃迁 · 中国互联网视角_

