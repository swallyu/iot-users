package main

import (
	"github.com/swallyu/iot-users/common"
	"github.com/swallyu/iot-users/config"
	"github.com/swallyu/iot-users/log"
	"github.com/swallyu/iot-users/users"
	"github.com/swallyu/iot-users/utils"
)

func main() {
	config.LoadConfig()
	log.Infoln(config.Conf.Database.Host)
	common.InitDbUtils()

	user := users.FindUser("admin")
	if user != nil {
		log.Println(user.UserName)
	}

	log.Infoln(config.Conf.Database.Host)
	utils.Sample()
	log.Info("SSS")
}
