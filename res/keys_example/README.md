# 密钥

## 1. OpenSSH

#### 生成
```bash
ssh-keygen -t rsa -b 2048 -f ~/.ssh/id_rsa -C "your_email@example.com"
```

#### 解释
- `-t rsa`：指定密钥的类型为 RSA。
- `-b 2048`：指定密钥的位数为 2048 位。
- `-f ~/.ssh/id_rsa`：指定私钥文件的保存路径。
- `-C "your_email@example.com"`：为公钥添加一个注释，通常是你的电子邮件地址或者其它标识信息。

#### 文件
- 私钥文件（如 id_rsa）不会包含邮箱信息，它只存储密钥本身。
- 公钥文件（如 id_rsa.pub）会包含邮箱信息，通常是密钥的最后一部分。例如：
```sql
ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAr+FfT7k7gOS6k4TfF5w5bKJvqx8X8gFdzg5x+Ys1sbE1... your_email@example.com

```

## 2. ppk

#### 生成
- 打开 PuTTYgen，点击 "Load" 按钮，选择你之前生成的私钥文件（id_rsa，不是 .pub 文件）。
- 点击 "Save private key"，保存为 .ppk 文件（例如，id_rsa_private.ppk）。
- 点击 "Save public key"，保存为 .ppk 文件（例如，id_rsa_public.ppk）。

#### 解释
- `.ppk` 是 **PuTTY Private Key** 文件的扩展名。PuTTY 是 Windows 上一个常用的 SSH 客户端，`.ppk` 文件是 PuTTY 使用的专有格式，用于存储私钥。

#### 文件
- `id_rsa_public.ppk`文件内容与`id_rsa.pub`文件内容**不一样**的。

## pem

#### 私钥生成
- 打开 PuTTYgen，点击 "Load" 按钮，选择你之前生成的私钥文件（id_rsa，不是 .pub 文件）。
- 在 PuTTYgen 窗口中，点击 Conversions 菜单，然后选择 Export OpenSSH key。
- 在文件保存对话框中，选择一个目录并输入文件名（例如 id_rsa_private.pem）。
- 保存时，文件扩展名应该是 .pem，这就是你需要的 PEM 格式私钥文件。

#### 公钥生成
- 提取公钥
使用以下 OpenSSL 命令来从私钥生成公钥：
```bash
openssl rsa -in id_rsa_private.pem -pubout -out id_rsa_public.pem
```
- 参数说明：
  + -in id_rsa_private.pem：指定输入的私钥文件。
  + -pubout：表示输出公钥。
  + -out id_rsa_public.pem：指定生成的公钥保存到 id_rsa_public.pem 文件。

### 文件
1. 编码格式：
- .pem 文件通常以 Base64 编码 存储数据，并且数据内容通常被包裹在标记中，比如 -----BEGIN CERTIFICATE----- 和 -----END CERTIFICATE-----。
- 示例内容：
```css
-----BEGIN CERTIFICATE-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA7vbVtYwPlwHfTjzoyWfa
...
-----END CERTIFICATE-----
```

2. 可以包含多种信息：
- .pem 文件不仅可以存储 公钥 或 私钥，也可以包含 SSL 证书、证书链等信息。它是非常通用的格式，可以用于很多加密协议，如 SSL/TLS、SSH。

3. 密钥存储：
- 对于 SSH 密钥，.pem 文件通常表示私钥（如通过 ssh-keygen 生成的文件），它可以用于身份验证。
- 公钥文件内容与id_rsa.pub完成**一致**。

4. 适用环境：
- .pem 格式是 OpenSSL 和 OpenSSH 的常用格式，适用于大多数 UNIX/Linux 系统、Web 服务器（如 Apache、Nginx）以及很多其他加密相关的应用。

