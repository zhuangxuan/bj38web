package conf

import (
	"fmt"

	"github.com/fsnotify/fsnotify"

	"github.com/spf13/viper"
)

// Conf 全局配置变量
var Conf = new(AppConfig)

type AppConfig struct {
	Name         string `mapstructure:"name"`
	Mode         string `mapstructure:"mode"`
	Port         int    `mapstructure:"port"`
	Version      string `mapstructure:"version"`
	StartTime    string `mapstructure:"start_time"`
	MachineID    int64  `mapstructure:"machine_id"`
	*Service     `mapstructure:"service"`
	*Etcd        `mapstructure:"etcd"`
	*Captcha     `mapstructure:"captcha"`
	*RedisConfig `mapstructure:"redis"`
	*MySQLConfig `mapstructure:"mysql"`
	*LogConfig   `mapstructure:"log"`
	*AuthConfig  `mapstructure:"auth"`
}

type Service struct {
	GetCaptcha string `mapstructure:"GetCaptcha"`
	User       string `mapstructure:"User"`
	GetArea    string `mapstructure:"GetArea"`
	House      string `mapstructure:"House"`
	Order      string `mapstructure:"Order"`
}

type Etcd struct {
	Address string `mapstructure:"address"`
}

type Captcha struct {
	Num     int `mapstructure:"num"`
	StrType int `mapstructure:"strType"`
}
type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
	MaxConnAge   int    `mapstructure:"max_conn_age"`
	IdleTimeout  int    `mapstructure:"idle_timeout"`
}
type MySQLConfig struct {
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	DB              string `mapstructure:"database"`
	Charset         string `mapstructure:"charset"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	MaxConnLifetime int    `mapstructure:"max_conn_lifetime"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type AuthConfig struct {
	JwtExpire int `mapstructure:"jwt_expire"`
}

func InitConfig() (err error) {
	viper.SetConfigFile("./conf/config.yml")

	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("ReadInConfig failed, err: %v", err))
	}

	// 将配置文件解析到对应的结构体中
	err = viper.Unmarshal(&Conf)
	if err != nil {
		panic(fmt.Errorf("unmarshal to Conf failed, err:%v", err))
	}

	// 监控配置文件是否修改
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("夭寿啦~配置文件被人修改啦...")
		if err := viper.Unmarshal(&Conf); err != nil {
			fmt.Printf("ReadInConfig failed, err: %v", err)
		}
	})
	return err
}
