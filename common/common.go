package common

import (
	"github.com/swallyu/iot-users/config"
	"github.com/swallyu/iot-users/utils"
)

var DbUtils utils.DbUtils

func InitDbUtils() {
	dbconfig := utils.Config{
		Host:         config.Conf.Database.Host,
		Port:         config.Conf.Database.Port,
		Password:     config.Conf.Database.Pwd,
		UserName:     config.Conf.Database.User,
		DatabaseName: config.Conf.Database.DbName,
		Dbms:         "mysql",
	}

	DbUtils = utils.NewDbUtils(dbconfig)
}
