package redis_distributed_lock

import (
	"github.com/google/uuid"
	"strings"
)

func GetToken() string {
	// 生成UUID
	uuidWithHyphen := uuid.New().String()

	// 去除UUID中的-
	uuidWithoutHyphen := strings.Replace(uuidWithHyphen, "-", "", -1)

	return uuidWithoutHyphen
}
