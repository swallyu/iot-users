package users

import (
	"time"

	"github.com/swallyu/iot-users/common"
)

type User struct {
	Id               uint64
	UserName         string
	LoginName        string
	Password         string
	Salt             string
	Status           int
	Type             int
	CreatedTime      time.Time
	CurrentLoginTime time.Time
	LastLoginTime    time.Time
	Extension        map[string]interface{}
	Remark           string
}

func FindUser(loginName string) *User {
	h := common.DbUtils.Query("select * from user_config where login_name=?", loginName)
	users, _ := h.Rows(mapToUser)
	ret := users.([]User)
	if len(ret) > 0 {
		return &ret[0]
	}
	return nil
}

func mapToUser(rows []map[string]string) (interface{}, error) {
	users := make([]User, len(rows))
	for k, v := range rows {
		tmp := User{
			UserName: v["user_name"],
		}
		users[k] = tmp
	}
	return users, nil
}
