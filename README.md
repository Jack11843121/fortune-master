# 算命大师 · fortune-master

> **开源 · 可自部署 · 单二进制**的**八字算命 / 命理解读 / 排盘应用**，内置老年休闲游戏、微信/支付宝支付、会员体系、AI 算命师、社区与后台。前后端合一，一条命令即可在本地或云服务器跑起来。

[![Go](https://img.shields.io/badge/Go-1.27-00ADD8?logo=go&logoColor=white)](go.mod)
[![Vue](https://img.shields.io/badge/Vue-3-4FC08D?logo=vue.js&logoColor=white)](web/package.json)
[![SQLite](https://img.shields.io/badge/SQLite-MySQL-cb6c2e)](config.yaml)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)
[![Payments](https://img.shields.io/badge/Pay-微信%20v3%20%7C%20支付宝-07c160)](config.yaml)

**关键词**：八字 · 算命 · 命理 · 排盘 · 四柱八字 · 生辰八字 · 十神 · 大运 · 流年 · 合婚 · 起名 · 黄历 · 每日一签 · 五行 · 老年游戏 · bazi · fortune-telling · lunar calendar · Go · Vue3

---

## ✨ 这是什么

`算命大师` 是一套开箱即用的**八字命理应用**，把「专业排盘、多场景解读、算乐游戏、先看后买、会员、AI、社区、后台」全部集成进一个 **Go + Vue3 的单二进制**里：

```
浏览器 ──► fortune（Go 单二进制，内嵌 Vue3 前端）
              ├── 八字排盘引擎（lunar-go）：四柱 / 十神 / 藏干 / 纳音 / 五行 / 大运 / 流年
              ├── 多场景解读：个人命理 / 合婚 / 招聘 / 财运·股市 / 宝宝起名
              ├── 算乐老年游戏：消消乐 / 五子棋 / 象棋 / 围棋 / 摇色子 / 我画你猜 / 五行音乐
              ├── Freemium：免费预览 + 解锁全文（mock / 微信 v3 / 支付宝 RSA2）
              ├── 会员体系 + AI 算命师（DeepSeek 等 OpenAI 兼容接口）
              ├── 交流社区 + 管理后台（看板 / 会员 / 订单 / 咨询）
              └── 存储：SQLite（本地零安装）/ MySQL 8.0（生产）
```

**定位**：面向**老人**及**帮他人算命的用户**，界面通俗、大字大按钮、好上手；后端排盘走权威历法库，结果严谨；命理文案积极吉祥、不渲染宿命论。

---

## 🎯 谁适合用 / 使用场景

- **个人**：给自己、家人、朋友算八字，看一生运势、财运、婚姻、健康。
- **命理从业者 / 大师**：当线上工具用，快速排盘 + 生成通俗解读，帮客户解读。
- **开发者 / 独立开发者**：想做一个命理类小程序 / H5 / 网站的起点，或二次开发接单。
- **企业 / 门店**：养生馆、婚介、茶馆、文旅等场景，做「测一测 / 算一算」引流互动。
- **学习参考**：完整的 Go + Vue3 + GORM + 微信/支付宝支付 + 内嵌前端的工程范本。

---

## ⭐ 核心亮点

| 亮点 | 说明 |
|---|---|
| 🧮 **专业排盘** | 基于 [lunar-go](https://github.com/6tail/lunar-go)，精确节气/干支/藏干/纳音，不手搓天文算法 |
| 📦 **单二进制部署** | 前端产物经 `go:embed` 内嵌，一个二进制 + 一个配置即可上线，无多服务运维负担 |
| 💰 **支付开箱即用** | 默认 mock 零配置跑通「下单→支付→解锁」全链路；拿到商户号切微信 v3 / 支付宝 |
| 🎮 **算乐老年游戏** | 围绕算命玩法设计的消消乐/棋牌，通关抽「灵签」，留存与趣味兼顾 |
| 🤖 **AI 算命师** | 结合用户八字对话答疑；未配大模型自动本地兜底，配 DeepSeek 等即增强 |
| 🛡 **安全基线** | 口令 bcrypt、令牌哈希、参数化查询、限流、支付验签、越权校验（见 docs） |
| 🔍 **SEO 友好** | 内置 `robots.txt` / `sitemap.xml`、首页完整 meta/OG/JSON-LD |
| 🌐 **零依赖起步** | 本地 SQLite 零安装，一把跑起来就能用；生产无缝切 MySQL |

---

## 🧭 功能总览

| 模块 | 能力 |
|---|---|
| 排盘 | 四柱、日主、十神、藏干、纳音、五行分布、大运、流年、生肖、神煞 |
| 个人命理 | 总览 / 性格 / 事业财运 / 情感姻缘 / 健康 / 大运 / 流年 / 趋吉避凶 8 维度 |
| 合婚匹配 | 双人八字 + 生肖合冲 + 五行互补，0–100 匹配评分 |
| 招聘看八字 | 候选人 vs 岗位五行、性格画像、与老板/团队匹配度 |
| 财运 / 股市 | 财运走向、投资性格、板块倾向、求财开运 |
| 宝宝起名 | 结合八字喜用与姓名学起名推荐 |
| 算乐游戏 | 五行消消乐、五子棋、象棋、围棋(9路)、摇色子、我画你猜、五行音乐、每日一签+黄历 |
| 先看后买 | 命盘+前几段免费看，全文解锁；解锁后永久可看 |
| 支付 | 微信支付 v3（JSAPI/Native/H5）+ 支付宝（RSA2），回调验签解密、幂等发货 |
| 会员 | 月付开通 VIP、成长值/等级 Lv1–9、AI 调用额度、会员全站免解锁 |
| AI 算命师 | 按八字对话答疑、按次扣额度、历史记录 |
| 社区 | 发帖 / 评论 / 点赞 / 分类 / 内容分享 |
| 后台 | 数据看板、会员管理、支付订单、日常咨询、社区审核、会员流水 |
| 运营 | 开发模式自动登录（免注册体验全功能）、合作广告位、SEO |

---

## 🔮 详细功能

### 1. 专业八字排盘
- 支持**公历 / 农历**输入，精确到分钟，自动处理真太阳时与节气换月。
- 输出：**四柱（年/月/日/时）、日主、十神、藏干、纳音、五行分布、大运起运、流年**。
- 随排盘附「每日一签」与「今日黄历宜忌」，按日期轮换。

### 2. 命理解读（8 个维度，人话讲透）
`总览 · 性格 · 事业财运 · 情感姻缘 · 健康 · 大运 · 流年 · 趋吉避凶`，把专业术语翻译成通俗、积极、吉祥的表达。

### 3. 多场景
- ☯️ **个人命理**（自己/家人/老人）
- 💞 **合婚匹配**（两人八字 + 生肖合冲 + 五行互补，0–100 分）
- 💼 **招聘看八字**（候选人 vs 岗位，性格与团队匹配）
- 💰 **财运 / 股市**（财运走向、投资性格、板块倾向）
- 👶 **宝宝起名**（八字喜用 + 姓名学）

### 4. 算乐 · 老年游戏（通关/对弈可抽「灵签」）
- 💎 **五行消消乐**：三消五行、连锁连消、达标抽签
- ⚫ **五子棋**：双人对弈 / 简单人机
- 🐴 **象棋**：完整规则（马脚象眼、炮架、将军/将死、将帅对脸）
- ⚪ **围棋**：9 路对弈
- 🎲 **摇色子** · 🎨 **我画你猜** · 🎵 **五行音乐**
- 🔮 **每日一签 + 今日黄历宜忌**

### 5. 先看后买（Freemium）与支付
- `unlock.free_preview_ratio` 控制免费预览比例（默认 35%），全文需解锁一次。
- 支付通道：**mock（默认）/ 微信支付 v3（JSAPI·Native·H5）/ 支付宝（RSA2）**。
- 回调验签 + 解密 + 幂等发货，重复回调不重复到账。

### 6. 会员体系 + AI 算命师
- 月付 10 元开通 VIP，赠 1000 次 AI；成长值/等级 Lv1–9；会员全站免解锁。
- AI 算命师按八字对话、按次扣额度、保留历史；未配大模型走本地兜底。

### 7. 交流社区 + 管理后台
- 社区：发帖、评论、点赞、分类；内容分享累计成长值。
- 后台：数据看板、会员管理、支付订单、日常咨询、社区审核、会员流水。

### 8. SEO 与运维
- `/robots.txt` + `/sitemap.xml`；首页完整 title/description/keywords/OG/JSON-LD。
- 接口限流、panic 兜底、CORS、请求体限制等安全措施已内置。

---

## 🏗 技术栈

- **后端**：Go 1.27 · [chi](https://github.com/go-chi/chi) 路由 · [GORM](https://gorm.io)（SQLite / MySQL 8.0）
- **前端**：Vue 3 · Vite（移动端 H5，构建产物 `go:embed` 内嵌）
- **排盘**：[lunar-go](https://github.com/6tail/lunar-go)（权威历法库）
- **支付**：微信支付 v3（手写 RSA-SHA256 签名 + AES-GCM）、支付宝 RSA2
- **AI**：OpenAI 兼容接口（DeepSeek / 通义 / 任意兼容端点）

## 📁 目录结构

```
fortune-master/
├── cmd/server/            # 程序入口
├── internal/
│   ├── config/            # 配置加载与默认值
│   ├── store/             # 数据模型 + SQLite/MySQL 存取
│   ├── fortune/           # 排盘引擎 + 解读 + 合婚/招聘/起名/财运匹配
│   ├── payment/           # mock + 微信 v3 + 支付宝 RSA2
│   ├── handler/           # HTTP API（命理/订单/后台/AI/社区/会员/登录）
│   ├── auth/              # bcrypt 口令 + 令牌
│   ├── ai/                # AI 算命师（OpenAI 兼容 + 本地兜底）
│   ├── ratelimit/         # 接口限流
│   └── webutil/           # JSON 工具
├── web/                   # Vue3 前端（dist 经 go:embed 内嵌）
├── scripts/               # 构建 / e2e 回归 / 游戏棋盘测试
├── deploy/                # Docker、docker-compose、MySQL init、systemd
├── docs/                  # 会员与安全、微信上线准备、QA
├── config.yaml            # 主配置（本地示例）
└── Makefile
```

---

## 🚀 快速开始

### 前置
- [Go](https://go.dev/dl/) ≥ 1.27
- [Node.js](https://nodejs.org/) ≥ 20（构建前端用）

### 一键构建并启动

```bash
# macOS / Linux
make build          # 自动 cd web && npm install && npm run build，再 go build
./fortune -config config.yaml

# Windows（PowerShell）
powershell -File scripts/build.ps1
.\fortune.exe -config config.yaml
```

浏览器打开 **http://localhost:8080** 即可使用（默认「模拟支付」，本地 SQLite 零安装）。

### 默认账号（首次启动自动写入，**上线前务必修改**）
- 管理后台：`http://localhost:8080/#/admin`，账号 `admin`，密码 `admin123`

### 本地体验支付
1. 任选场景算命 → 结果页显示命盘 + 部分解读，下方提示「全文尚未解锁」。
2. 点「解锁全文」→「模拟支付 ¥9.90」→ 点击即完成（与真实支付同一套订单状态流转）。
3. 自动刷新为全文，本次结果永久解锁。

---

## ⚙️ 配置说明（config.yaml）

| 配置 | 说明 | 默认 |
|---|---|---|
| `app.name` | 应用名 | 算命大师 |
| `app.base_url` | 对外地址（支付回调/二维码用），上线改 `https://你的已备案域名` | `http://localhost:8080` |
| `app.price_cents` | 解锁全文价格（分） | `990`（¥9.9） |
| `server.addr` | 监听地址 | `:8080` |
| `database.driver` | `sqlite` / `mysql` | `sqlite` |
| `database.dsn` | SQLite 路径，或 MySQL DSN | `data/fortune.db` |
| `payment.provider` | `mock` / `wechat` / `alipay` | `mock` |
| `unlock.free_preview_ratio` | 免费预览比例（0~1） | `0.35` |
| `membership.*` | 月费 / 额度 / 充值包 | 10 元档 |
| `ai.enabled` / `ai.api_key` / `ai.model` | AI 算命师（OpenAI 兼容，留空走本地兜底） | 开启 |
| `admin.username` / `admin.password` | 后台种子账号（bcrypt 落库，之后以库为准） | `admin` / `admin123` |

> ⚠️ 配置可能含数据库口令、支付/大模型密钥，请勿提交到公开仓库；上线后 `chmod 600 config.yaml`。生产示例见 [`deploy/config.prod.yaml`](deploy/config.prod.yaml)。

---

## 💳 支付接入（可选）

### 微信支付 v3
1. 在[微信支付商户平台](https://pay.weixin.qq.com/)开通商户号，取 `mchid`、APIv3 密钥、商户 API 证书 `apiclient_key.pem`（+序列号）、平台证书。
2. 改 `config.yaml`：

```yaml
payment:
  provider: wechat
  wechat:
    appid: "wx..."
    mchid: "1234567890"
    mch_serial_no: "..."
    mch_private_key_path: "./keys/apiclient_key.pem"
    platform_cert_path: "./keys/wechatpay_xxx.pem"
    api_v3_key: "..."
    notify_url: "https://你的域名/api/v1/pay/notify/wechat"
    pay_type: jsapi      # jsapi / native / h5
```

3. 商户平台「支付结果通知地址」配成 `notify_url`。

> 微信支付需要：**已备案域名 + HTTPS**；JSAPI/H5 需配置「网页授权域名 / JS 安全域名」。

### 支付宝（RSA2）
`config.yaml` 的 `payment.alipay` 填 `app_id`、应用私钥、支付宝公钥、回调地址，`payment.provider: alipay`。详见 [`docs/WECHAT-上线准备.md`](docs/WECHAT-上线准备.md)。

---

## ☁️ 部署

### 方式一：Docker Compose（推荐，含 MySQL）
```bash
cd deploy
# 改好 config.prod.yaml（域名/数据库密码），证书放 keys/
docker compose up -d --build
```

### 方式二：二进制 + systemd
```bash
make build-linux            # 交叉编译 fortune-linux（先 make build-web）
# 上传 fortune-linux + config.yaml 到 /opt/fortune/
sudo cp deploy/systemd/fortune.service /etc/systemd/system/
sudo systemctl daemon-reload && sudo systemctl enable --now fortune
```

### 方式三：轻量服务器 / 裸进程
上传二进制 + `config.yaml`，`./fortune -config config.yaml &`，配 Nginx 反代（80/443 → 8080）+ HTTPS（推荐 acme.sh 自动续期）。

---

## 🧪 测试

```bash
go test ./...                                          # 单测（排盘/支付验签/棋盘）
node scripts/e2e-remote.mjs http://127.0.0.1:8080      # 15 项端到端回归
```

15 项覆盖：meta / 个人命理 / 合婚 / 招聘 / 起名 / 下单-支付-解锁 / 微信登录降级 / 静态资源 / 账户登录 / 财运 / 会员订阅-发货 / 会员免解锁 / AI / 社区 / 后台。

---

## 📡 API 概览（节选）

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/v1/meta` | 全局配置（价格、场景、角色） |
| POST | `/api/v1/readings` | 创建算命（scenario: personal/marriage/hiring/wealth/naming） |
| GET | `/api/v1/readings/{id}` | 读取结果（未付费=预览+锁定，已付费=全文） |
| POST | `/api/v1/orders` | 创建解锁订单 |
| GET | `/api/v1/orders/{orderNo}` | 查询订单状态 |
| POST | `/api/v1/pay/mock/confirm` | 模拟支付确认（仅 mock） |
| POST | `/api/v1/pay/notify/wechat` | 微信支付回调 |
| POST | `/api/v1/auth/dev` | 开发模式自动登录（免注册体验全功能） |
| POST | `/api/v1/ai/chat` | AI 算命师对话 |
| POST | `/api/v1/community/posts` | 社区发帖 |
| POST | `/api/v1/admin/login` | 后台登录 |
| GET | `/api/v1/admin/dashboard` | 后台看板 |

完整能力与安全说明见 [`docs/`](docs/)。

---

## ❓ FAQ

**Q：本地能直接跑吗？** 能。默认 SQLite + 模拟支付，`make build` 后 `./fortune` 即可，零外部依赖。

**Q：要上线收费怎么办？** 先备案域名 + 配 HTTPS，再把 `payment.provider` 切到 `wechat`/`alipay`，填入商户号与密钥即可。

**Q：想换掉默认后台账号密码？** 改 `config.yaml` 的 `admin.username/password`（首次启动种入数据库，之后以库为准）。

**Q：AI 算命师怎么接？** `ai.base_url` + `ai.api_key` 填任意 OpenAI 兼容端点（如 DeepSeek），不填则走本地兜底。

**Q：单二进制怎么实现的？** `web/dist` 由 `//go:embed all:dist` 内嵌进 Go 二进制，配合 `http.FileServer` 提供前端静态资源。

---

## 🗺 Roadmap

- [ ] 小程序 / 公众号版本
- [ ] SSR 预渲染，深化 SEO
- [ ] 更多命理模型（紫微斗数、手相、塔罗）
- [ ] 后台数据导出与报表
- [x] 微信支付 v3 + 支付宝（RSA2）

---

## 🤝 贡献

欢迎提 Issue / PR。约定：代码正确性优先（排盘走 lunar-go，勿手搓历法）；文案用积极吉祥的话术；提交前 `go test ./...` 通过。

## ⚖️ 免责声明

本项目排盘计算基于权威历法库、逻辑严谨；命理解读属传统文化内容，请理性看待、结合自身情况判断，重大事项请咨询专业人士。内容不构成医疗、投资、婚恋、招聘等决策建议。

---

## 📞 定制开发 / 联系方式

本项目 MIT 开源、永久免费使用。

如需**定制开发**（新增命理场景、私有化部署、企业级改造、小程序/公众号对接、AI Agent 落地等），或需要**打包 + 部署 + 使用指导**，欢迎联系：

📞 **电话：15919897624**
📧 **邮箱：11843121@qq.com**

> 可提供：命理/工具类应用定制、算乐游戏扩展、支付与登录对接、打包部署、运营指导等定制业务。

---

## 📄 许可

MIT © fortune-master Contributors