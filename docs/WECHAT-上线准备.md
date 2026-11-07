# 微信上线准备清单（微信支付 + 小程序）

> 结论先行：**微信支付必须「企业 / 个体工商户」主体**，个人主体不能开通支付。
> 本应用走「H5 + 微信支付 JSAPI」最快；小程序可选「web-view 壳」复用现有 H5。

## 两条可选路线

| | 路线 A：公众号 H5 | 路线 B：微信小程序（web-view 壳） |
|---|---|---|
| 前端改动 | 几乎不用改 | 用 `miniprogram/` 壳加载 H5 |
| 支付 | JSAPI 支付（微信内浏览器） | 小程序原生支付（需跳原生页） |
| 主体 | 可个人（但个人不能支付） | 必须企业/个体户 |
| 需求 | 公众号 AppID + 商户号 + 备案域名 + HTTPS | 小程序 AppID + 商户号 + 备案域名 + HTTPS + 业务域名 |

## 一、注册与认证（你去办）

1. 官网注册：`https://mp.weixin.qq.com`（邮箱即可，邮箱未注册过公众号/小程序）。
2. 主体类型选 **企业或个体工商户**（要营业执照）。
3. **微信认证**：企业认证 ¥300/年，需对公账户打款或微信认证打款验证。
4. 完成后到「开发 → 开发设置」拿：
   - **AppID**
   - **AppSecret**

## 二、开通微信支付（你去办）

1. 后台 →「微信支付」→ 申请开通商户号（自动关联）。
2. 登录商户平台 `https://pay.weixin.qq.com`：
   - 记下 **商户号 mchid**
   - 「账户中心 → API安全」：设置 **APIv3 密钥**（32 位，自己定）
   - 下载 **商户 API 证书**（`apiclient_key.pem`）及 **证书序列号**
   - 下载 **平台证书**（`wechatpay_xxx.pem`）或开启自动更新
3. 「产品中心」开通：**JSAPI 支付**（公众号）或 **小程序支付**（小程序）。
4. JSAPI 支付需在「产品中心 → JSAPI 支付 → 授权目录/回调」配置。

## 三、域名 + 备案 + HTTPS（你去办，周期 1~2 周）

1. 购买域名（腾讯云/阿里云）。
2. **ICP 备案**：在云厂商备案中心提交（需身份证、人脸、域名、服务器信息）。
3. 申请 **HTTPS 证书**（云厂商免费证书或 Let's Encrypt）。
4. 白名单：
   - 公众号：「设置 → 公众号设置 → 功能设置」配置「网页授权域名」「JS接口安全域名」。
   - 小程序：「开发 → 开发设置 → 业务域名」添加 https 域名（需上传校验文件到网站根目录）。

## 四、服务器（腾讯云）

- 轻量应用服务器 / CVM，安装 Docker 或直接部署本项目（见 README「部署到腾讯云」）。
- Nginx 反代 + HTTPS；把备案/微信校验文件放到站点根目录。

## 五、办完把这些信息给我（填配置）

| 项目 | 填入位置 |
|---|---|
| 小程序/公众号 AppID | `config.yaml` payment.wechat.appid |
| AppSecret | `config.yaml` payment.wechat.app_secret |
| 商户号 mchid | payment.wechat.mchid |
| 商户证书序列号 | payment.wechat.mch_serial_no |
| apiclient_key.pem 路径 | payment.wechat.mch_private_key_path |
| 平台证书路径 | payment.wechat.platform_cert_path |
| APIv3 密钥 | payment.wechat.api_v3_key |
| 备案域名 | app.base_url + notify_url |

## 代码侧已就绪（无需资质即可开发）

- [x] 微信支付 v3：JSAPI / Native / H5 下单 + 回调验签解密（`internal/payment/wechat.go`）
- [x] 微信登录拿 openid（`POST /api/v1/wechat/login`，内部封装 `jscode2session`）
- [x] 下单接口支持传入 openid（JSAPI 支付用）
- [x] 小程序 web-view 壳（`miniprogram/`，用微信开发者工具导入即可）
- [x] 模拟支付可在本地完整跑通全流程

> 上线前务必把 `payment.provider` 从 `mock` 改为 `wechat` 并填全上面的密钥。