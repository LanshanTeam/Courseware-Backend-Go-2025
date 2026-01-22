package main

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// 配置结构体，用于映射配置文件中的值
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	App      AppConfig
}

// 服务器配置
type ServerConfig struct {
	Host string
	Port int
}

// 数据库配置
type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     int
	Username string
	Password string
	DBName   string
}

// 应用配置
type AppConfig struct {
	Name        string
	Environment string
	Debug       bool
}

func main() {
	// 初始化配置
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 打印配置信息
	printConfig(config)
}

// loadConfig 加载配置文件
func loadConfig() (*Config, error) {
	// 创建一个新的viper实例
	v := viper.New()

	// 设置配置文件名称
	v.SetConfigName("config")
	// 设置配置文件类型
	v.SetConfigType("yaml")
	// 设置配置文件搜索路径
	v.AddConfigPath("./")
	v.AddConfigPath("./configs/")

	// 读取配置文件
	err := v.ReadInConfig()
	if err != nil {
		// 配置文件不存在时，使用默认值
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("配置文件未找到，使用默认值")
		} else {
			// 配置文件存在但解析错误
			return nil, fmt.Errorf("解析配置文件错误: %w", err)
		}
	}

	// 设置环境变量前缀
	v.SetEnvPrefix("APP")
	// 自动绑定环境变量
	v.AutomaticEnv()

	// 设置默认值
	v.SetDefault("server.host", "localhost")
	v.SetDefault("server.port", 8080)
	v.SetDefault("database.driver", "mysql")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 3306)
	v.SetDefault("database.username", "root")
	v.SetDefault("database.password", "")
	v.SetDefault("database.dbname", "app")
	v.SetDefault("app.name", "myapp")
	v.SetDefault("app.environment", "development")
	v.SetDefault("app.debug", true)

	// 解析配置到结构体
	var config Config
	err = v.Unmarshal(&config)
	if err != nil {
		return nil, fmt.Errorf("无法解析配置到结构体: %w", err)
	}

	return &config, nil
}

// printConfig 打印配置信息
func printConfig(config *Config) {
	fmt.Println("===== 配置信息 =====")
	fmt.Printf("服务器: %s:%d\n", config.Server.Host, config.Server.Port)
	fmt.Printf("数据库: %s@%s:%d/%s\n", 
		config.Database.Username, 
		config.Database.Host, 
		config.Database.Port, 
		config.Database.DBName)
	fmt.Printf("应用: %s (%s)\n", config.App.Name, config.App.Environment)
	fmt.Printf("调试模式: %v\n", config.App.Debug)
	fmt.Println("===================")
}