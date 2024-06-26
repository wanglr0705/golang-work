package utils

import (
	"go_xorm_mysql_redis/config"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

func GetYaml(cfg *config.Config) *config.Config {
	// 读取YAML配置文件
	data, err := os.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// 将读取的YAML数据解码到结构体中
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}
