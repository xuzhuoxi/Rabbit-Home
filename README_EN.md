# Rabbit-Home

English | [简体中文](./README.md)

## Introduction
Rabbit-Home is a project used for managing and monitoring server instances. It supports managing server instances via command-line and HTTP requests, including functions such as connecting, disconnecting, updating status, querying information, and kicking out instances.

## Project Structure
```
core/               // Core part of the Rabbit-Home project, contains client, command-line, configuration management, and server-related logic.
├── client/         // Contains logic related to the client, for connecting, disconnecting, and updating server instance statuses.
├── cmd/            // Contains command-line related logic, for managing server instances via the command line.
├── conf/           // Contains configuration-related logic, for parsing and managing parameters in the configuration file.
├── home/           // Contains server-related logic, for starting the server, handling HTTP requests, and managing server instances.
└── core.go         // core package exposes interfaces and structure definitions
res/
└── config.yaml     // Configuration file used for configuring the server's listening address, internal and external IP control, timeout parameters, and logging parameters.
src/
├── main.go         // Entry file for the project, starts the server and begins command-line listening.
└── main_test.go    // Test case for testing connection functionality.
```

- **core**: Core package, containing client, command-line, configuration management, and server-related logic.
  + **client**: Client-related logic, such as connecting, disconnecting, and updating statuses.
  + **cmd**: Command-line related logic, such as listing, querying information, and kicking out instances.
  + **conf**: Configuration-related logic, such as IP validators and configuration parsing.
  + **home**: Server-related logic, such as server startup, request handling, and entity management.
  + **core.go**: Defines interfaces and structure for the core package.
- **res**: Contains configuration files for configuring the server's listening address, internal and external IP control, timeout parameters, and logging parameters.
  + **config.yaml**: Configuration file for configuring the server's listening address, internal and external IP control, timeout parameters, and logging parameters.
- **src**: Contains the project's entry file, which starts the server and begins command-line listening.
  + **main.go**: The entry file for the project, starts the server and begins command-line listening.
  + **main_test.go**: Test case for testing connection functionality.

## Installation and Running
You can use Rabbit-Home in your project by using `go mod` or cloning the entire repository into the `gopath`.

### Installation
Supports using `go.mod` or `gopath` for managing the repository.

- Load the repository via `gopath`:
```bash
go get -u github.com/xuzhuoxi/Rabbit-Home
```   

- Load the repository via `go.mod`:
  Clone the project locally and install dependencies:
```bash
git clone github.com/xuzhuoxi/Rabbit-Home 
cd Rabbit-Home
go mod tidy
```

### Run the Project
```bash
go run src/main.go
```

## Configuration
The configuration file `res/config.yaml` is used to configure the server's listening address, internal and external IP control, timeout parameters, and logging parameters.

### 1. HTTP Service Configuration
```yaml
http:
  addr: "127.0.0.1:9000"
```
- **addr**: The server's listening address and port.
  + Value: 127.0.0.1:9000
  + Description: The server will listen for HTTP requests on the local address 127.0.0.1 at port 9000.

### 2. Internal IP Control
```yaml
internal:
  post: false
  allows_on: false
  allows:
  blocks_on: true
  blocks:
    - "192.168.0.1"
    - "10.0.0.1-20"
```
- **post**: Whether POST requests are required.
- **allows_on**: Whether the internal IP whitelist is enabled.
  + Value: false/true
  + Description: If true, the `allows` configuration takes effect.
- **allows**: List of internal IP whitelist.
  + Value: Array containing IP addresses or IP ranges (only supports ranges in the last segment using '-').
  + Description: If `allows_on` is false, IP addresses in the whitelist won't take effect.
- **blocks_on**: Whether the internal IP blacklist is enabled.
  + Value: false/true
  + Description: If true, the `blocks` configuration takes effect.
- **blocks**: List of internal IP blacklist.
  + Value: Array containing IP addresses or IP ranges (only supports ranges in the last segment using '-').
  + Description: If `blocks_on` is false, IP addresses in the blacklist won't take effect.
- **Note**:
  + Priority: Blacklist > Whitelist

### 3. External IP Control
```yaml
  post: true
  allows_on: false
  allows:
  blocks_on: true
  blocks:
    - "8.8.8.8"
```
- **post**: Whether POST requests are required.
- **allows_on**: Whether the external IP whitelist is enabled.
  + Value: false/true
  + Description: If true, the `allows` configuration takes effect.
- **allows**: List of external IP whitelist.
  + Value: Array containing IP addresses or IP ranges (only supports ranges in the last segment using '-').
  + Description: If `allows_on` is false, IP addresses in the whitelist won't take effect.
- **blocks_on**: Whether the external IP blacklist is enabled.
  + Value: false/true
  + Description: If true, the `blocks` configuration takes effect.
- **blocks**: List of external IP blacklist.
  + Value: Array containing IP addresses or IP ranges (only supports ranges in the last segment using '-').
  + Description: If `blocks_on` is false, IP addresses in the blacklist won't take effect.
- **Note**:
  + Priority: Blacklist > Whitelist

### 4. Timeout Settings
```yaml
timeout: 300000000000
```
- **timeout**: Timeout duration for instances.
  + Value: In nanoseconds (i.e., 5 minutes).
  + Description: If an instance doesn't send an information update to the server within the timeout period, it is considered disconnected.

### 5. Log Configuration
```yaml
log:
  type: 0   # 0:Console 1:RollingFile 2:DailyFile 3:DailyRollingFile
  level: 2  # 0:All 1:Trace 2:Debug 3:Info 4:Warn 5:Error 6:Fatal 7:Off
  path: "RabbitHome.log"
  size: '1MB' # 1MB
```
- **type**: Log record type.
  + Value: 
    - 0: Console
    - 1: RollingFile
    - 2: DailyFile
    - 3: DailyRollingFile
  + Description: 
    - Console: Outputs to the console. Ignores `path` and `size` configurations.
    - RollingFile: Rolling file logging, files are named with sequence numbers.
    - DailyFile: Logs in daily files, named with the date.
    - DailyRollingFile: Rolling file logging, with daily logs and sequence numbers.

- **path**: Log file path.
  + Value: Path for saving log files, supports relative and absolute paths.
  + Description: Ignored in Console mode.

- **level**: Log level.
  + Value: 
    - 0: All
    - 1: Trace
    - 2: Debug
    - 3: Info
    - 4: Warn
    - 5: Error
    - 6: Fatal
    - 7: Off

- **size**: Log file size.
  + Format: Value[unit]
  + Description: In RollingFile and DailyRollingFile modes, this defines the file size limit for rolling logs.

## Instance Information
Rabbit-Home records instance information as `RegisteredEntity`.

```go
type RegisteredEntity struct {
    core.LinkEntity
    State  core.EntityStatus       // Instance simple status
    Detail core.EntityDetailStatus // Instance detailed status
    lastUpdateNano int64 // Last update timestamp
    hit            int   // Hit count
}
```

## Command-Line Usage
After starting the project, the following operations can be performed via the command line:

- List instance list: 
  + Example: `list -name=Name -on=[true|false] -pid=PID`
- Query instance info:
  + Example: `info -id=Id`
- Kick instance:
  + Example: `kick -id=Id`

## Server Instance Connection
Instances can communicate with Rabbit-Home via the `core/client/` API.

## Dependencies
- **infra-go** [https://github.com/xuzhuoxi/infra-go](https://github.com/xuzhuoxi/infra-go): Basic library support.
- **goxc** [https://github.com/laher/goxc](https://github.com/laher/goxc): Packaging dependencies, mainly for cross-compiling.
- **json-iterator** [https://github.com/json-iterator/go](https://github.com/json-iterator/go): JSON parser with corresponding structure.

## Contact
xuzhuoxi  
<xuzhuoxi@gmail.com> or <mailxuzhuoxi@163.com> or <m_xuzhuoxi@outlook.com>

## License
Rabbit-Home source code is available under the MIT [License](/LICENSE).
