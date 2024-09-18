package main

import (
	"CengkeHelper/api"
	"CengkeHelper/logger"
)

func main() {
	err := api.HelperApp.Run(":8000")
	if err != nil {
		logger.Error(err)
		return
	} // 监听并在 0.0.0.0:8000 上启动服务
}
