module github.com/meateam/permission-service

go 1.13

require (
	github.com/golang/protobuf v1.4.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0
	github.com/meateam/elasticsearch-logger v1.1.3-0.20190901111807-4e8b84fb9fda
	github.com/olivere/elastic/v7 v7.0.22 // indirect
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/viper v1.6.1
	go.elastic.co/apm/module/apmgrpc v1.6.0 // indirect
	go.elastic.co/apm/module/apmmongo v1.6.0
	go.mongodb.org/mongo-driver v1.4.2
	golang.org/x/net v0.0.0-20201021035429-f5854403a974 // indirect
	golang.org/x/sys v0.0.0-20201020230747-6e5568b54d1a // indirect
	google.golang.org/genproto v0.0.0-20201021134325-0d71844de594 // indirect
	google.golang.org/grpc v1.33.1
	google.golang.org/grpc/examples v0.0.0-20201021230544-4e8458e5c638 // indirect
)

replace github.com/meateam/permission-service/server => ./server

replace github.com/meateam/permission-service/service => ./service

replace github.com/meateam/permission-service/proto => ./proto
