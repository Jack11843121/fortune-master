# 算命大师 · fortune-master

> 一款可自部署的开源**八字算命应用**：专业排盘 + 多场景命理解读 + 算乐老年游戏 + 会员/AI/社区/后台，单二进制即可上线。

[![Go](https://img.shields.io/badge/Go-1.27-00ADD8?logo=go&logoColor=white)](go.mod)
[![Vue](https://img.shields.io/badge/Vue-3-4FC08D?logo=vue.js&logoColor=white)](web/package.json)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)

**定位**：面向老人及帮他人算命的用户，界面通俗、大字大按钮、好上手；后端排盘用权威历法库，结果严谨，文案积极吉祥。

---

## 这是什么

`算命大师` 把「排盘、解读、游戏、支付、会员、AI、社区、后台」做成一个完整的应用，**前后端打包进单个 Go 二进制**，拿一台服务器（或本地电脑）即可运行：

```
浏览器 ──► fortune（Go 单二进制，内嵌 Vue3 前端）
              ├── 八字排盘引擎（lunar-go）：四柱/十神/藏干/纳音/大运/流年
              ├── 多场景解读：个人命理/合婚/招聘/财运/起名
              ├── 算乐老年游戏：消消乐/五子棋/象棋/围棋/摇色子/我画你猜/五行音乐
              ├── Freemium：免费预览 + 解锁全文（mock / 微信 / 支付宝）
              ├── 会员体系 + AI 算命师（DeepSeek 等 OpenAI 兼容接口）
              ├── 交流社区 + 管理后台（看板/会员/订单）
              └── 存储：SQLite（本地零安装）/ MySQL（生产）
```

---

## 功能特性

### 专业八字排盘
- 基于 [lunar-go](https://github.com/6tail/lunar-go)，精确节气/干支/藏干/纳音/大运/流年，不手搓天文算法。
- 排四柱、十神、五行分布、大运、流年。

### 多场景命理解读
| 场景 | 说明 |
|---|---|
| ☯️ 个人命理 | 自己 / 家人 / 老人，8 个维度：总览、性格、事业财运、情感姻缘、健康、大运、流年、趋吉避凶 |
| 💞 合婚匹配 | 两人八字 + 生肖合冲 + 五行互补，0–100 评分 |
| 💼 招聘看八字 | 候选人 vs 岗位五行、性格画像、与老板/团队匹配度 |
| 💰 财运 / 股市 | 财运走向、投资性格、板块倾向、求财开运 |
| 👶 宝宝起名 | 结合八字喜用与姓名学起名推荐 |

### 算乐 · 老年游戏（通关/对弈可抽「灵签」）
- 💎 五行消消乐（三消五行、连锁连消）
- ⚫ 五子棋（双人 / 简单人机 AI）
- 🐴 象棋（完整规则：马脚象眼、炮架、将军/将死、将帅对脸）
- ⚪ 围棋（9 路对弈）
- 🎲 摇色子、🎨 我画你猜、🎵 五行音乐
- 🔮 每日一签 + 今日黄历宜忌

### 先看后买（Freemium）
- 命盘 + 前几段解读免费看（`unlock.free_preview_ratio` 控制比例），全文解锁后永久可看本次结果。
- 本地默认「**模拟支付**」零配置跑通全流程，拿到商户号填配置即切真实支付。

### 会员 · AI · 社区 · 后台
- **会员体系**：月付开通 VIP、成长值/等级（Lv1–9）、AI 调用额度。
- **AI 算命师**：结合八字对话答疑，按次扣额度；未配模型时自动本地兜底。
- **交流社区**：发帖 / 评论 / 点赞 / 分类；内容分享累计成长值。
- **管理后台**：数据看板、会员管理、支付订单、日常咨询、社区审核、会员流水。

### 支付
- 微信支付 v3（JSAPI / Native / H5）+ 支付宝（RSA2），回调验签解密、幂等发货。
- 默认 `mock`，本地无需任何资质即可体验「下单 → 支付 → 解锁」完整链路。

### 部署体验
- **单二进制**：前端 Vue3 产物经 `go:embed` 内嵌，一个二进制 + 一个配置即可运行。
- **SEO**：`/robots.txt` + `/sitemap.xml`、首页完整 meta/OG/JSON-LD。

---

## 技术栈

- **后端**：Go 1.27 · [chi](https://github.com/go-chi/chi) 路由 · [GORM](https://gorm.io)（SQLite / MySQL）
- **前端**：Vue 3 · Vite（移动端 H5，构建产物内嵌）
- **排盘**：[lunar-go](https://github.com/6tail/lunar-go)
- **存储**：SQLite（本地默认）/ MySQL 8.0（生产）

## 目录结构

```
fortune-master/
├── cmd/server/            # 程序入口
├── internal/
│   ├── config/            # 配置加载与默认值
│   ├── store/             # 数据模型 + SQLite/MySQL 存取
│   ├── fortune/           # 排盘引擎 + 解读内容 + 合婚/招聘/起名/财运匹配
│   ├── payment/           # 支付抽象（mock + 微信 v3 + 支付宝 RSA2）
│   ├── handler/           # HTTP API（命理/订单/后台/AI/社区/会员/登录）
│   ├── auth/              # 口令哈希与令牌
│   ├── ai/                # AI 算命师（OpenAI 兼容 + 本地兜底）
│   ├── ratelimit/         # 接口限流
│   └── webutil/           # JSON 工具
├── web/                   # Vue3 前端（dist 经 go:embed 内嵌）
│   ├── src/               # 页面与组件
│   └── embed.go           # embed 声明
├── scripts/               # 构建 / e2e 回归 / 游戏棋盘测试
├── deploy/                # Docker、docker-compose、MySQL init、systemd
├── docs/                  # 会员与安全说明、微信上线准备、QA
├── config.yaml            # 主配置（本地示例）
├── Dockerfile
└── Makefile
```

---

## 快速开始（本地运行）

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

### 默认账号（首次启动自动写入，上线前务必修改）
- 管理后台：`http://localhost:8080/#/admin`，账号 `admin`，密码 `admin123`

### 本地体验支付流程
1. 任选场景算命 → 结果页显示命盘 + 部分解读，下方提示「全文尚未解锁」。
2. 点「解锁全文」→ 弹出「模拟支付 ¥9.90」→ 点击即完成（与真实支付同一套订单状态流转）。
3. 自动刷新为全文；刷新或重新进入，本次结果永久解锁。

---

## 配置说明（config.yaml）

| 配置 | 说明 | 默认 |
|---|---|---|
| `app.name` | 应用名 | 算命大师 |
| `app.base_url` | 对外地址（支付回调/二维码用），上线改为 `https://你的已备案域名` | `http://localhost:8080` |
| `app.price_cents` | 解锁全文价格（分） | `990`（¥9.9） |
| `server.addr` | 监听地址 | `:8080` |
| `database.driver` | `sqlite` / `mysql` | `sqlite` |
| `database.dsn` | SQLite 文件路径，或 MySQL DSN | `data/fortune.db` |
| `payment.provider` | `mock` / `wechat` / `alipay` | `mock` |
| `unlock.free_preview_ratio` | 免费预览比例（0~1） | `0.35` |
| `membership.*` | 月费 / 额度 / 充值包 | 10 元档 |
| `ai.enabled` / `ai.api_key` / `ai.model` | AI 算命师（OpenAI 兼容接口，留空走本地兜底） | 开启 |
| `admin.username` / `admin.password` | 管理后台种子账号（落库 bcrypt，之后以库为准） | `admin` / `admin123` |

> ⚠️ 配置文件可能包含数据库口令、支付/大模型密钥等，请勿提交到公开仓库；部署上线后建议 `chmod 600 config.yaml`。生产示例见 [`deploy/config.prod.yaml`](deploy/config.prod.yaml)。

---

## 支付接入（可选，拿到资质后）

### 微信支付 v3
1. 在[微信支付商户平台](https://pay.weixin.qq.com/)开通商户号，获取 `mchid`、APIv3 密钥、商户 API 证书（`apiclient_key.pem` + 序列号）、平台证书。
2. 修改 `config.yaml`：

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

3. 商户平台「支付结果通知地址」配成上面的 `notify_url`。

> 微信支付需要：**已备案域名 + HTTPS**；JSAPI/H5 需在公众平台配置「网页授权域名 / JS 安全域名」。

### 支付宝（RSA2）
`config.yaml` 的 `payment.alipay` 填入 `app_id`、应用私钥、支付宝公钥与回调地址，`payment.provider: alipay` 即可。详见 [`docs/WECHAT-上线准备.md`](docs/WECHAT-上线准备.md)。

---

## 部署

### 方式一：Docker Compose（推荐，含 MySQL）

```bash
cd deploy
# 先把 config.prod.yaml 改成你的域名/数据库密码，放到 keys/ 放支付证书
docker compose up -d --build
```

### 方式二：二进制 + systemd

```bash
make build-linux            # 交叉编译 fortune-linux（先 make build-web）
# 上传 fortune-linux + config.yaml 到服务器，例如 /opt/fortune/
sudo cp deploy/systemd/fortune.service /etc/systemd/system/
sudo systemctl daemon-reload && sudo systemctl enable --now fortune
```

### 方式三：轻量服务器 / 裸进程
上传 Linux 二进制 + `config.yaml`，`./fortune -config config.yaml &`，再配 Nginx 反向代理（80/443 → 8080）+ 免费 HTTPS 证书（推荐 acme.sh 自动续期）。

---

## 测试

```bash
go test ./...                                  # 单元测试（排盘引擎/支付验签/棋盘规则）
node scripts/e2e-remote.mjs http://127.0.0.1:8080   # 15 项端到端回归
# 覆盖：meta / 个人命理 / 合婚 / 招聘 / 起名 / 下单-支付-解锁 / 微信登录降级 /
#       静态资源 / 账户登录 / 财运 / 会员订阅-发货 / 会员免解锁 / AI / 社区 / 后台
```

---

## API 概览（节选）

| 方法 | 路径 | 说明 |
|---|---|---|
| GET | `/api/v1/meta` | 全局配置（价格、场景、角色） |
| POST | `/api/v1/readings` | 创建一次算命（`scenario`: personal/marriage/hiring/wealth/naming） |
| GET | `/api/v1/readings/{id}` | 读取结果（未付费=预览+锁定，已付费=全文） |
| POST | `/api/v1/orders` | 创建解锁订单 |
| GET | `/api/v1/orders/{orderNo}` | 查询订单状态 |
| POST | `/api/v1/pay/mock/confirm` | 模拟支付确认（仅 mock） |
| POST | `/api/v1/pay/notify/wechat` | 微信支付回调 |
| POST | `/api/v1/auth/dev` | 开发模式自动登录（免注册体验全功能） |
| POST | `/api/v1/ai/chat` | AI 算命师对话 |
| POST | `/api/v1/community/posts` | 社区发帖 |
| POST | `/api/v1/admin/login` | 管理后台登录 |
| GET | `/api/v1/admin/dashboard` | 后台看板 |

更完整的能力与安全说明见 [`docs/`](docs/)。

---

## 免责声明

本项目排盘计算基于权威历法库、逻辑严谨；命理解读属传统文化内容，请理性看待、结合自身情况判断，重大事项请咨询专业人士。内容不构成医疗、投资、婚恋、招聘等决策建议。

---

## 定制开发 / 联系方式

本项目 MIT 开源、永久免费使用。

如需**定制开发**（新增命理场景、私有化部署、企业级改造、小程序/公众号对接、AI Agent 落地等），或需要**打包 + 部署 + 使用指导**，欢迎联系：

📞 **电话：15919897624**
📧 **邮箱：11843121@qq.com**

> 可提供：命理/工具类应用定制、算乐游戏扩展、支付与登录对接、打包部署、运营指导等定制业务。

---

## 许可

MIT © fortune-master Contributors