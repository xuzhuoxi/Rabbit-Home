# Rabbit-Home

## 准备中的功能

### 检验
- 在Rabbit-Home中部署Rsa公钥
- 在Rabbit-Server部署Rsa私钥
- Rabbit-Server连接Rabbit-Home时进行RSA密钥校验
- Rabbit-Home校验成功后为Rabbit-Server生成Session Code
- 客户端向Rabbit-Home进行路由请求时，把Session Code 一起发送给客户端
- 客户端连接Rabbit-Server时，带上Session Code，Rabbit-Server对Session Code进行检验
- 客户端与Rabbit-Server通信时，都要带上Session Code，Rabbit-Server对Session Code进行校验

### Rabbit-Home中RSA公钥动态部署
- 识别密钥类型
- 启动时预加载 ./keys/rsa.key ./keys/pkcs.key
- 动态公钥文件存放在 ./keys/rsa/ ./keys/pkcs/
