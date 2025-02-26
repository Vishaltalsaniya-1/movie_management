package cnf

import (
	"log"

	"github.com/caarlos0/env"
)

type Mysql struct {
	DB_HOST     string `env:"DB_HOST" envDefault:"localhost"`
	DB_USER     string `env:"DB_USER" envDefault:"vishal"`
	DB_PASSWORD string `env:"DB_PASSWORD" envDefault:"Vishal@123"`
	DB_NAME     string `env:"DB_NAME" envDefault:"movies"`
	DB_PORT     string `env:"DB_PORT" envDefault:"3306"`
}

type ConsumerConfig struct {
	Url                string `env:"URL" validate:"required" envDefault:"amqp://guest:guest@localhost:5672/"`
	Exchange           string `env:"EXCHANGE_NAME"  envDefault:"movie_add_exchange"`
	ExchangeType       string `env:"EXCHANGE_TYPE"  envDefault:"direct"`
	PrefetchCount      int    `env:"PREFETCH_COUNT"  envDefault:"100"`
	ConnectionPoolSize int    `env:"CONNECTIONPOOL_SIZE"  envDefault:"10"`
	QueueName          string `env:"QUEUE_NAME" envDefault:"movie_add"`
	BindingKeyName     string `env:"BINDING_KEY_NAME" envDefault:"movie_add_bindkey"`
	DelayedQueueName   string `env:"DELAYED_QUEUE_NAME" envDefault:"movie_add_delay_queue"`
	QueueTaskName      string `env:"MOVIE_QUEUE_TASK"  envDefault:"MovieCreatedTask"`
	RunConsumer        bool   `env:"RUN_CONSUMER" envDefault:"true"`
	RunProducer        bool   `env:"RUN_PRODUCER" envDefault:"true"`
}

var Consumerconfig ConsumerConfig

func Loadcosumer() {
	if err := env.Parse(&Consumerconfig); err != nil {
		log.Fatalf("Failed to parse environment variables: %v", err)
	}
}

type Config struct {
	Mysql     Mysql
	JwtSecret string `env:"JWT_SECRET_KEY" envDefault:"vishal"` // Secret key for JWT

}

func LoadConfig() (*Config, error) {
	var cfg Config
	if err := env.Parse(&cfg.Mysql); err != nil {
		log.Printf("Failed to load config: %v", err)
		return nil, err
	}
	return &cfg, nil
}

// type Rabbitmq struct {
// 	RabbitmqValue string `env:"RabbitmqValue"  envDefault:"TRUE"`
// }

// func RabbitmqMovie() (*Rabbitmq, error) {
// 	var RabbitmqConfig Rabbitmq

// 	if err := env.Parse(&RabbitmqConfig); err != nil {
// 		log.Println("failed to Rabbitmq", err)

// 	}
// 	return &RabbitmqConfig, nil

// }
