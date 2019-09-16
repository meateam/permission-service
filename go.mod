module github.com/meateam/permission-service

go 1.13

require (
	github.com/golang/protobuf v1.3.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0
	github.com/meateam/elasticsearch-logger v1.1.3-0.20190901111807-4e8b84fb9fda
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/viper v1.4.0
	go.elastic.co/apm/module/apmmongo v1.5.0
	go.mongodb.org/mongo-driver v1.1.0
	google.golang.org/grpc v1.23.1
)

replace github.com/meateam/permission-service/server => ./server

replace github.com/meateam/permission-service/service => ./service

replace github.com/meateam/permission-service/proto => ./proto
