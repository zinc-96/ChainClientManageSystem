package main

import (
	"ChainClientManageSystem/config"
	"ChainClientManageSystem/internal/router"
)

func main() {
	config.InitConfig()
	router.InitRouterAndServe()
}
