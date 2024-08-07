package setting

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var Conf = new(AppConfig)

type AppConfig struct {
	Name      string `mapstructure:"name"`
	Mode      string `mapstructure:"mode"`
	Version   string `mapstructure:"version"`
	StartTime string `mapstructure:"start_time"`
	MachineID int64  `mapstructure:"machine_id"`
	Port      int    `mapstructure:"port"`

	*LogConfig   `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
	*WsConfig    `mapstructure:"websocket"`
	*MQConfig    `mapstructure:"mq"`
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

type RedisConfig struct {
	Addr         string `mapstructure:"addr"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type WsConfig struct {
	ReadBufferSize  int              `mapstructure:"read_buffer_size"`
	WriteBufferSize int              `mapstructure:"write_buffer_size"`
	ChannelSize     int              `mapstructure:"channel_size"`
	MaxRetries      int              `mapstructure:"max_retries"`
	WSServers       []WSServerConfig `mapstructure:"servers"`
	TickerPeriod    time.Duration    `mapstructure:"ticker_period"`
	WriteDeadline   time.Duration    `mapstructure:"write_deadline"`
	ReadDeadline    time.Duration    `mapstructure:"read_deadline"`
}

type WSServerConfig struct {
	ID         int    `mapstructure:"server_id"`
	BucketSize int    `mapstructure:"bucket_size"`
	Addr       string `mapstructure:"addr"`
}

type MQConfig struct {
	URL                string `mapstructure:"url"`
	Exchange           string `mapstructure:"exchange"`
	Queue              string `mapstructure:"queue"`
	RoutingKey         string `mapstructure:"routing_key"`
	LogicPublishersNum int32  `mapstructure:"logic_publishers"`
	LogicConsumersNum  int    `mapstructure:"logic_consumers"`
}

func Init(filePath string) (err error) {
	viper.SetConfigFile(filePath)

	err = viper.ReadInConfig() // 读取配置信息
	if err != nil {
		// 读取配置信息失败
		fmt.Printf("viper.ReadInConfig failed, err:%v\n", err)
		return
	}

	// 把读取到的配置信息反序列化到 Conf 变量中
	if err := viper.Unmarshal(Conf); err != nil {
		fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件修改了...")
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("viper.Unmarshal failed, err:%v\n", err)
		}
	})
	return
}
