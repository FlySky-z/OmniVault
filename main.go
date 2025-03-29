package main

import (
	"log"
	"omnivault/config"
	"omnivault/routes"
)

func main() {
	// 初始化配置
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("LoadConfig failed! %s", err)
	}

	r := routes.SetupRouter()
	if err := r.Run(); err != nil {
		log.Fatalf("router start failed! %s", err)
	}
}
