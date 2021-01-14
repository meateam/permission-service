module github.com/meateam/permission-service

go 1.13

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/golang/protobuf v1.4.3
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0
	github.com/meateam/api-gateway v0.0.0-20210113140950-b2c25acb408f
	github.com/meateam/download-service v0.0.0-20191216103739-80620a5c7311
	github.com/meateam/elasticsearch-logger v1.1.3-0.20190901111807-4e8b84fb9fda
	github.com/olivere/elastic/v7 v7.0.22 // indirect
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/viper v1.6.1
	go.elastic.co/apm/module/apmgrpc v1.6.0 // indirect
	go.elastic.co/apm/module/apmmongo v1.6.0
	go.mongodb.org/mongo-driver v1.4.2
	google.golang.org/grpc v1.34.0
)

replace github.com/meateam/permission-service/server => ./server

replace github.com/meateam/permission-service/service => ./service

replace github.com/meateam/permission-service/proto => ./proto
