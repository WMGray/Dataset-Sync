package conf

import (
	"fmt"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

// 软件信息
type SoftwareInfo struct {
	Name    string `mapstructure:"name"`    // 软件名称
	Version string `mapstructure:"version"` // 软件版本

	//*LogConfig   `mapstructure:"log"`   // 日志配置
	*MySQLConfig `mapstructure:"mysql"` // MySQL配置
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Password     string `mapstructure:"password"`
	DB           string `mapstructure:"dbname"`
	Port         int    `mapstructure:"port"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}

var Conf = new(SoftwareInfo)

func Init(confpath string) (err error) {
	// viper 读取配置文件
	viper.SetConfigFile(confpath)
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("读取配置文件失败, err:%v\n", err)
		return
	}
	// 把读取到的配置信息反序列化到 Conf 变量中
	if err = viper.Unmarshal(Conf); err != nil {
		fmt.Printf("反序列化配置文件失败, err:%v\n", err)
		return
	}
	fmt.Println("反序列化成功，输出配置")

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了...")
		if err = viper.Unmarshal(Conf); err != nil {
			fmt.Printf("反序列化配置文件失败, err:%v\n", err)
		}
	})
	return
}
