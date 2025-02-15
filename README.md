# Rabbit-Home 项目

## 简介
Rabbit-Home 是一个用于管理和监控服务器实例的项目。它支持通过命令行和HTTP请求来管理服务器实例，包括连接、断开连接、更新状态、查询信息和踢除实例等功能。

## 项目结构
```
core/               // Rabbit-Home 项目的核心部分，包含了客户端、命令行、配置管理和服务器相关的逻辑。
├── client/         // 包含了客户端相关的逻辑，用于连接、断开连接和更新服务器实例的状态。
├── cmd/            // 包含了命令行相关的逻辑，用于通过命令行管理服务器实例。
├── conf/           // 包含了配置相关的逻辑，用于解析和管理配置文件中的参数。
├── home/           // 包含了服务器相关的逻辑，用于启动服务器、处理 HTTP 请求和管理服务器实例。
└── core.go         // core包开放接口与结构体定义
res/
└── config.yaml     // 配置文件，用于配置服务器的监听地址、内部和外部IP控制、超时参数和日志记录参数。
src/
├── main.go         // 项目的入口文件，启动服务器并开始命令行监听。
└── main_test.go    // 测试用例，用于测试连接功能。
```

- **core**: 核心包，包含了客户端、命令行、配置管理和服务器相关的逻辑。
  + **client**: 客户端相关的逻辑，如连接、断开连接和更新状态等。
  + **cmd**: 命令行相关的逻辑，如列表、信息查询和踢除实例等。
  + **conf**: 配置相关的逻辑，如IP验证器和配置解析等。
  + **home**: 服务器相关的逻辑，如服务器启动、处理请求和管理实体等。
  + **core.go**: core包开放接口与结构体定义
- **res**: 包含了配置文件，用于配置服务器的监听地址、内部和外部IP控制、超时参数和日志记录参数。
  + **config.yaml**: 配置文件，用于配置服务器的监听地址、内部和外部IP控制、超时参数和日志记录参数。
- **src**: 包含了项目的入口文件，启动服务器并开始命令行监听。
  + **main.go**: 项目的入口文件，启动服务器并开始命令行监听。
  + **main_test.go**: 测试用例，用于测试连接功能。

## 安装和运行
通过 `go mod` 或 克隆整个仓库到 `gopath` 中便在项目中使用 Rabbit-Home.

### 安装
支持使用go.mod或gopath管理仓库

- 通过gopath加载仓库
``` bash
go get -u github.com/xuzhuoxi/Rabbit-Home
```   

- 通过go.mod加载仓库
  克隆项目到本地并安装依赖
``` bash
git clone github.com/xuzhuoxi/Rabbit-Home 
cd Rabbit-Home
go mod tidy
```

### 运行项目
``` bash
go run src/main.go
```

## 配置
配置文件 `res/config.yaml` 用于配置服务器的监听地址、内部和外部IP控制、超时参数和日志记录参数。

### 1. HTTP 服务配置
```yaml
http:
  addr: "127.0.0.1:9000"
```
- addr: 服务器的监听地址和端口。
  + 值: 127.0.0.1:9000
  + 说明: 服务器将在本地地址 127.0.0.1 的端口 9000 上监听 HTTP 请求

### 2. 内部 IP 控制
```yaml
internal:
  allows_on: false
  allows:
  blocks_on: true
  blocks:
    - "192.168.0.1"
    - "10.0.0.1-20"
```
- **allows_on**: 是否启用内部 IP 白名单。
  + 值: false/true
  + 说明: 为 true 时， allows 配置生效。
- **allows**: 内部 IP 白名单列表。
  + 值: 数组，内容为ip地址或i地址范围(仅支持最后一组使用'-'表示范围)，不包含接口
  + 说明: 当 allows_on 为 false，白名单列表中的 IP 地址不会生效。
- **blocks_on**: 是否启用内部 IP 黑名单。
  + 值: false/true
  + 说明: 为 true 时， blocks 配置生效。
- **blocks**: 内部 IP 黑名单列表。
  + 值: 数组，内容为ip地址或i地址范围(仅支持最后一组使用'-'表示范围)，不包含接口
  + 说明: 当 blocks_on 为 false，黑名单列表中的 IP 地址不会生效
- **注意：**
  + 优先级：黑名单 > 白名单

### 3. 外部 IP 控制
```yaml
  allows_on: false
  allows:
  blocks_on: true
  blocks:
    - "8.8.8.8"
```
- **allows_on**: 是否启用外部 IP 白名单。
  + 值: false/true
  + 说明: 为 true 时， allows 配置生效。
- **allows**: 外部 IP 白名单列表。
  + 值: 数组，内容为ip地址或i地址范围(仅支持最后一组使用'-'表示范围)，不包含接口
  + 说明: 当 allows_on 为 false，白名单列表中的 IP 地址不会生效。
- **blocks_on**: 是否启用外部 IP 黑名单。
  + 值: false/true
  + 说明: 为 true 时， blocks 配置生效。
- **blocks**: 内部 IP 黑名单列表。
  + 值: 数组，内容为ip地址或i地址范围(仅支持最后一组使用'-'表示范围)，不包含接口
  + 说明: 当 blocks_on 为 false，黑名单列表中的 IP 地址不会生效
- **注意：**
  + 优先级：黑名单 > 白名单

### 4. 超时设置
```yaml
timeout: 300000000000
```
- timeout: 实例的超时时间。
  + 值: 纳秒数（即 5 分钟）
  + 说明: 如果实例在 timeout时间内没有向服务器发送信息更新，则认为该实例已断开连接。 

### 5. 日志记录配置
```yaml
log:
  type: 0   # 0:Console 1:RollingFile 2:DailyFile 3:DailyRollingFile
  level: 2  # 0:All 1:Trace 2:Debug 3:Info 4:Warn 5:Error 6:Fatal 7:Off
  path: "RabbitHome.log"
  size: '1MB' # 1MB
```
- type: 日志记录类型。
  + 值: 
    - 支持 0:Console 1:RollingFile 2:DailyFile 3:DailyRollingFile
  + 说明: 
    - Console: 
      + 输出到控制台。
      + 忽略 path 和 size 配置。
    - RollingFile: 
      + 滚动文件记录，文件名以序号关联命名.
      + 文件大小到达 size 配置值，则创建新的文件，并继续记录日志。
      + 文件保存路径为：path中目录,最新日志文件名为path中文件名。
      + 文件名格式：
        - 最新文件名：(文件名).log, 文件名为path中配置文件名。
        - 历史文件名：(文件名)_(序号).log, 文件名为path中配置文件名，序号从0开始。
    - DailyFile: <br>
      + 每日文件记录，文件名以日期关联命名。
      + 忽略 size 配置。
      + 文件保存路径为：path中目录。
      + 文件名格式：(文件名)_(yyyyMMdd).log, 文件名为path中配置文件名，yyyyMMdd为日期。
    - DailyRollingFile:<br> 
      + 每日滚动文件记录，文件名以日期命名，文件名格式为：文件名_yyyyMMdd_序号.log。
      + 文件大小到达 size 配置值，则创建新的文件，并继续记录日志。
      + 文件保存路径为：path中目录, 最新日志文件名为path中文件名。
      + 文件名格式：
          - 最新文件名：(文件名)_(yyyyMMdd).log, 文件名为path中配置文件名, yyyyMMdd为日期。
          - 历史文件名：(文件名)_(yyyyMMdd)_(序号).log, 文件名为path中配置文件名，yyyyMMdd为日期, 序号从0开始。
- path: 日志文件路径。
  + 值: 
    - 日志文件保存的路径，支持相对路径和绝对路径。
  + 说明: 
    - Console 模式下忽略。
    - 文件名信息参与 RollingFile 和 DailyRollingFile 模式的文件名命名。
- level: 日志级别。
  + 值: 
    - 支持 0:All 1:Trace 2:Debug 3:Info 4:Warn 5:Error 6:Fatal 7:Off
  + 说明: 只记录 >= 值级别的日志信息。
- size: 日志文件大小。
  + 格式：数值[单位]
  + 说明: 
    - 数值: 支持小数
    - 单位: 支持 KB, MB, GB, TB, PB, EB。 不写表示为 Byte，但**不能写**Byte。
    - 在 RollingFile 和 DailyRollingFile 模式下生效，表示滚动文件大小上限。

## 实例信息
Rabbit-Home 记录的实例信息为 RegisteredEntity
```go
type RegisteredEntity struct {
	core.LinkEntity
	State  core.EntityStatus       // 实例简单状态
	Detail core.EntityDetailStatus // 实例详细状态

	lastUpdateNano int64 // 上一次更新时间戳
	hit            int   // 命中次数
}

type LinkEntity struct {
    Id         string `json:"id"`      // 实例Id(唯一)
    PlatformId string `json:"pid"`     // 平台Id
    Name       string `json:"name"`    // 实例名称(不唯一)
    Network    string `json:"network"` // 连接类型
    Addr       string `json:"addr"`    // 连接地址
}

type EntityStatus struct {
    Id     string  `json:"id"`     // 实例Id
    Weight float64 `json:"weight"` // 压力系数
}

type EntityDetailStatus struct {
    Id             string `json:"id"`           // 实例Id
    StartTimestamp int64  `json:"start"`        // 启动时间戳(纳秒)
    StatsInterval  int64  `json:"sta_interval"` // 统计间隔

    MaxLinks      uint64 `json:"max_links"`  // 最大连接数
    TotalReqCount int64  `json:"total_reg"`  // 总请求数
    TotalRespTime int64  `json:"total_resp"` // 总响应时间
    MaxRespTime   int64  `json:"max_resp"`   // 最大响应时间(纳秒)
    Links         uint64 `json:"links"`      // 连接数

    StatsTimestamp    int64 `json:"sta_start"` // 统计开始时间戳(纳秒)
    StatsReqCount     int64 `json:"sta_req"`   // 统计请求数
    StatsRespUnixNano int64 `json:"sta_resp"`  // 统计响应时间(纳称)

    Keys string `json:"sta_interval"` // 属性启用标记
}
```
- LinkEntity: 实例连接信息。
  + Id: 实例Id(唯一)
  + PlatformId: 平台Id
  + Name: 实例名称(不唯一)
  + Network: 连接类型
  + Addr: 连接地址
- State: 实例简单状态。
  + Id: 实例Id
  + Weight: 压力系数。
- Detail: 实例详细状态。
  + Id: 实例Id
  + StartTimestamp: 启动时间戳(纳秒)
  + StatsInterval: 统计间隔
  + MaxLinks: 最大连接数
  + TotalReqCount: 总请求数
  + TotalRespTime: 总响应时间
  + MaxRespTime: 最大响应时间(纳秒)
  + Links: 连接数
  + StatsTimestamp: 统计开始时间戳(纳秒)
  + StatsReqCount: 统计请求数
  + StatsRespUnixNano: 统计响应时间(纳秒)
  + Keys: 属性启用标记
- lastUpdateNano: 上一次更新时间戳。
- hit: 命中次数。

## 命令行使用
启动项目后，可以通过命令行进行以下操作:

- 列出实例
  + 示例: `list -name=Name -on=[true|false] -pid=PID`
  + -name: 实例名称。
  + -on: 实例是否在线。
  + -pid: 实例平台Id。
  + 以上参数如果不写，则列出所有符合的实例。
- 查询实例信息
  + 示例: `info -id=Id`
  + -id: 实例Id。
  + id参数为必填参数。
- 踢除实例
  + 示例: `kick -id=Id`
  + -id: 实例Id。
  + id参数为必填参数。
  + **功能未实现**。

## 服务器实例连接
实例可以通过 core/client/ 包下的函数，与 Rabbit-Home 进行通信。
以下以 {httpUrl} 代表config.yaml中配置 http.addr 关联的地址。
- 服务器连接到 Rabbit-Home：`{httpUrl}/link`
  + LinkWithGet 
  + LinkWithPost
- 服务器断开 Rabbit-Home：`{httpUrl}/unlink`
  + UnlinkWithGet
  + UnlinkWithPost
- 服务器更新信息到 Rabbit-Home：`{httpUrl}/update`
  + UnlinkWithGet
  + UpdateWithPost
- 查询实例路由：`{httpUrl}/route`

## 客户端请求
实例可以通过 core/client/ 包下的函数，与 Rabbit-Home 进行通信，获得得应该连接的服务器实例信息。
以下以 {httpUrl} 代表config.yaml中配置 http.addr 关联的地址。
- 客户端请求服务器实例：`{httpUrl}/route?data=xxxxxxxx`
- 

## 依赖库
- infra-go [https://github.com/xuzhuoxi/infra-go](https://github.com/xuzhuoxi/infra-go)<br>
  基础库支持库。
- goxc [https://github.com/laher/goxc](https://github.com/laher/goxc)<br>
  打包依赖库，主要用于交叉编译
- json-iterator [https://github.com/json-iterator/go](https://github.com/json-iterator/go)<br>
  带对应结构体的Json解释库

## Contact
xuzhuoxi<br>
<xuzhuoxi@gmail.com> or <mailxuzhuoxi@163.com> or <m_xuzhuoxi@outlook.com>

## License
Rabbit-Home source code is available under the MIT [License](/LICENSE).