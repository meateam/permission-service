module github.com/meateam/permission-service

go 1.13

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/golang/protobuf v1.4.1
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0
	github.com/meateam/api-gateway v0.0.0-20201118094700-aee405e404f7
	github.com/meateam/elasticsearch-logger v1.1.3-0.20190901111807-4e8b84fb9fda
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/viper v1.6.1
	go.elastic.co/apm/module/apmmongo v1.6.0
	go.mongodb.org/mongo-driver v1.3.0
	google.golang.org/grpc v1.27.0
)

replace github.com/meateam/permission-service/server => ./server

replace github.com/meateam/permission-service/service => ./service

replace github.com/meateam/permission-service/proto => ./proto
