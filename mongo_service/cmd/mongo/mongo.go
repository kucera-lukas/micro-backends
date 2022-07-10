package main

import (
	"github.com/kucera-lukas/micro-backends/mongo-service/pkg/infrastructure/env"
	"github.com/kucera-lukas/micro-backends/mongo-service/pkg/infrastructure/server"
)

func main() {
	config := env.MustLoad()
	server.Run(config)
}
