package conf

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// 软件信息
type SoftwareInfo struct {
	*AppConfig     `mapstructure:"app"`     // 软件信息
	*DatasetConfig `mapstructure:"dataset"` // 数据集配置
	*MySQLConfig   `mapstructure:"mysql"`   // MySQL配置
}

type AppConfig struct {
	AppName        string `mapstructure:"name"`        // 软件名称
	AppVersion     string `mapstructure:"version"`     // 软件版本
	AppDescription string `mapstructure:"description"` // 软件描述
	AppAuthor      string `mapstructure:"author"`      // 软件作者
	AppStartup     string `mapstructure:"startup"`     // 软件启动方式
}

type DatasetConfig struct {
	TmpDir        string `mapstructure:"tmp_dir"`
	SaveDir       string `mapstructure:"save_dir"`
	AutoRename    bool   `mapstructure:"auto_rename"`
	AutoRenameKey string `mapstructure:"auto_rename_key"`
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
	// 设置 viper 配置
	viper.SetConfigFile(confpath)
	viper.SetConfigType("yaml") // 明确指定配置文件类型

	// 读取配置文件
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
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("反序列化配置文件失败, err:%v\n", err)
			return
		}
	})
	return
}

// SaveConfig 保存配置文件
func SaveConfig() error {
	// 检查配置字段并逐个更新viper -- 通过设置可更改的字段
	// Conf.DatasetConfig
	if Conf.DatasetConfig != nil {
		// 不要设置整个结构体，而是逐个设置字段
		viper.Set("dataset.tmp_dir", Conf.DatasetConfig.TmpDir)
		viper.Set("dataset.save_dir", Conf.DatasetConfig.SaveDir)
		viper.Set("dataset.auto_rename", Conf.DatasetConfig.AutoRename)
		viper.Set("dataset.auto_rename_key", Conf.DatasetConfig.AutoRenameKey)
	}

	// Conf.MySQLConfig
	if Conf.MySQLConfig != nil {
		viper.Set("mysql.host", Conf.MySQLConfig.Host)
		viper.Set("mysql.user", Conf.MySQLConfig.User)
		viper.Set("mysql.password", Conf.MySQLConfig.Password)
		viper.Set("mysql.dbname", Conf.MySQLConfig.DB)
		viper.Set("mysql.port", Conf.MySQLConfig.Port)
		viper.Set("mysql.max_open_conns", Conf.MySQLConfig.MaxOpenConns)
		viper.Set("mysql.max_idle_conns", Conf.MySQLConfig.MaxIdleConns)
	}

	// 写入配置文件
	return viper.WriteConfig()
}
