package main

import (
	"github.com/meateam/permission-service/server"
)

func main() {
	server.NewServer(nil).Serve(nil)
}
