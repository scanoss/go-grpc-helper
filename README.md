# SCANOSS Platform 2.0 gRPC Helper Package
Welcome to the SCANOSS Platform 2.0 gRPC server helper package.

This package contains helper functions to make development of Go gRPC services easier to configure.

## Repository Structure
This repository is made up of the following components:
* [File manipulation](pkg/files/files.go)
* [gRPC server helpers](pkg/grpc/server/server.go)
* [REST/gRPC gateway setup](pkg/grpc/gateway/gateway.go)
* [Database helpers](pkg/grpc/database/database.go)
* [Utilities](pkg/grpc/utils/utils.go)

## Usage
### File Manipulation
The [files](pkg/files) package provides the following helpers:
* Check if TLS should be enabled or not
* Load filter config files

#### TLS
```go
startTLS, err = CheckTLS("server.crt", "server.key")
```

#### Filtering
```go
allowedIPs, deniedIPs, err = LoadFiltering("allow_list.txt", "deny_list.txt")
```
### gRPC Server
The [server](pkg/grpc/server) package provides the following helpers:
* Configure
* Start
* Stop

#### Configure
```go
listen, server, err := SetupGrpcServer(":0", "server.crt", "server.key", allowedIPs, deniedIPs, true, true, false)
```

#### Start
```go
StartGrpcServer(listen, server, true)
```

#### Stop
```go
err = WaitServerComplete(srv, server)
```

### REST/gRPC Gateway
The [gateway](pkg/grpc/gateway) package provides the following helpers:
* Configure
* Start

#### Configure
```go
allowedIPs := []string{"127.0.0.1"}
deniedIPs := []string{"192.168.0.1"}
srv, mux, gateway, opts, err := SetupGateway("9443", "8443", "server.crt", allowedIPs, deniedIPs, true, false, true)
```
#### Start
```go
StartGateway(srv, "server.crt", "server.key", true)
```

### Database
The [database](pkg/grpc/database) package provide the following helpers:
* Open DB connection
* Setup DB options
* Close DB connection
* Close SQL connection

#### Open
```go
db, err := OpenDBConnection(":memory:", "sqlite3", "", "", "", "", "")
```

#### Setup
```go
err = SetDBOptionsAndPing(db)
```

#### Close DB Connection
```go
CloseDBConnection(db)
```

#### Close SQL Connection
```go
CloseSQLConnection(conn)
```

## Bugs/Features
To request features or alert about bugs, please do so [here](https://github.com/scanoss/go-grpc-helper/issues).

## Changelog
Details of major changes to the library can be found in [CHANGELOG.md](CHANGELOG.md).
