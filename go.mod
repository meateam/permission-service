module github.com/meateam/permission-service

go 1.13

require (
	github.com/golang/protobuf v1.4.1
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0
	github.com/meateam/elasticsearch-logger v1.1.3-0.20190901111807-4e8b84fb9fda
	github.com/olivere/elastic/v7 v7.0.9 // indirect
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/viper v1.6.1
	go.elastic.co/apm/module/apmgrpc v1.6.0 // indirect
	go.elastic.co/apm/module/apmmongo v1.6.0
	go.mongodb.org/mongo-driver v1.3.0
	golang.org/x/net v0.0.0-20191002035440-2ec189313ef0 // indirect
	golang.org/x/sys v0.0.0-20200116001909-b77594299b42 // indirect
	google.golang.org/grpc v1.27.0
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

replace github.com/meateam/permission-service/server => ./server

replace github.com/meateam/permission-service/service => ./service

replace github.com/meateam/permission-service/proto => ./proto
