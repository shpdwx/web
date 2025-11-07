package conf

import (
	"github.com/spf13/viper"
)

type Conf struct {
	Minio    Minio    `mapstructure:"Minio"`
	CogView  CogView  `mapstructure:"CogView"`
	RabbitMQ RabbitMQ `mapstructure:"RabbitMQ"`
}

type Minio struct {
	Endpoint string `mapstructure:"endpoint"`
	Bucket   string `mapstructure:"bucket"`
	AK       string `mapstructure:"ak"`
	SK       string `mapstructure:"sk"`
	Dir      string `mapstructure:"dir"`
}

type CogView struct {
	Api    string `mapstructure:"api"`
	Token  string `mapstructure:"token"`
	Model  string `mapstructure:"model"`
	Origin string `mapstructure:"origin"`
}

type RabbitMQ struct {
	Endpoint     string `mapstructure:"endpoint"`
	User         string `mapstructure:"user"`
	Passwd       string `mapstructure:"passwd"`
	Vhost        string `mapstructure:"vhost"`
	Exchange     string `mapstructure:"exchange"`
	ExchangeType string `mapstructure:"exchangeType"`
	RouteKey     string `mapstructure:"routeKey"`
}

// conf.go conf.yml 是成对出现的
func NewConf() (Conf, error) {

	v := viper.New()
	v.SetConfigFile("conf/conf.yml")

	if err := v.ReadInConfig(); err != nil {
		return Conf{}, err
	}

	def(v)

	var conf Conf
	if err := v.Unmarshal(&conf); err != nil {
		return Conf{}, err
	}

	return conf, nil
}

func def(v *viper.Viper) {

	// minio def
	v.SetDefault("Minio.endpoint", "127.0.0.1:9000")

	// cogview def
	v.SetDefault("CogView.api", "https://open.bigmodel.cn/api/paas/v4/images/generations")

	// rabbitmq def
	v.SetDefault("RabbitMQ.endpoint", "127.0.0.1:5672")
	v.SetDefault("RabbitMQ.user", "guest")
	v.SetDefault("RabbitMQ.passwd", "guest")

}
