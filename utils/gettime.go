package utils

import (
	"go_xorm_mysql_redis/item"
	"go_xorm_mysql_redis/types"
	"time"
)

func GetTime(region string) time.Time {
	var loc *time.Location
	var err error
	// 设置时区
	switch region {
	case "uk":
		//英国时间
		loc, err = time.LoadLocation("Europe/London")
		break
	case "jp":
		//日本时间
		loc, err = time.LoadLocation("Asia/Tokyo")
		break
	case "ur":
		//俄罗斯时间
		loc, err = time.LoadLocation("Europe/Moscow")
		break
	default:
		return time.Now()
	}

	if err != nil {
		//LogError(ch, err)
		panic(pojo.PanicData{Code: types.ErrMissingAppLocalHeader, Error: err})
		return time.Now()
	}
	// 获取当前时间
	return time.Now().In(loc)
}
