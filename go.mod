module github.com/meateam/permission-service

go 1.12

require (
	github.com/golang/protobuf v1.3.1
	github.com/grpc-ecosystem/go-grpc-middleware v1.0.0
	github.com/meateam/elasticsearch-logger v1.1.3-0.20190901111807-4e8b84fb9fda
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/viper v1.4.0
	google.golang.org/grpc v1.21.0
)

replace git.apache.org/thrift.git => github.com/apache/thrift v0.12.0

replace github.com/meateam/permission-service/server => ./server

replace github.com/meateam/permission-service/permission => ./permission

replace github.com/meateam/permission-service/proto => ./proto
