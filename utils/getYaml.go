package utils

import (
	"go_xorm_mysql_redis/config"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

func GetYaml(cfg *config.Config) *config.Config {
	//logStorage := NewLogStorage()
	//logStorage.WriteLogFiler("正在读取yaml文件")
	// 读取YAML配置文件
	data, err := os.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	// 将读取的YAML数据解码到结构体中
	err = yaml.Unmarshal(data, cfg)
	if err != nil {
		//logStorage.WriteLogFiler("读取yaml文件失败，错误：" + err.Error())
	}
	//logStorage.WriteLogFiler("读取yaml文件成功")
	return cfg
}
