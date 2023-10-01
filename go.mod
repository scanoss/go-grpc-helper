module github.com/scanoss/go-grpc-helper

go 1.19

require (
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.17.1
	github.com/jmoiron/sqlx v1.3.5
	github.com/lib/pq v1.10.9
	github.com/mattn/go-sqlite3 v1.14.17
	github.com/scanoss/ipfilter/v2 v2.0.2
	github.com/scanoss/zap-logging-helper v0.2.1
	github.com/stretchr/testify v1.8.4
	golang.org/x/net v0.15.0
	google.golang.org/grpc v1.57.0
)

//replace github.com/scanoss/ipfilter => ../../../oss/ipfilter

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/uuid v1.3.1 // indirect
	github.com/phuslu/iploc v1.0.20230201 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rogpeppe/go-internal v1.10.0 // indirect
	github.com/tomasen/realip v0.0.0-20180522021738-f0c99a92ddce // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.25.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230822172742-b8732ec3820d // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230822172742-b8732ec3820d // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
