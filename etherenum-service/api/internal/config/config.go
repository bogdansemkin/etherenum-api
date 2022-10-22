package config

type (
	Config struct {
		Postgres
		Log
		Mongo
		Etherscan
		HTTP
		RabbitMQ
	}

	Postgres struct {
		User     string `env:"POSTGRESQL_USER" env-default:"postgres"`
		Password string `env:"POSTGRESQL_PASSWORD" env-default:"postgres"`
		Host     string `env:"POSTGRESQL_HOST" env-default:"127.0.0.1"`
		Database string `env:"POSTGRESQL_DATABASE" env-default:"chat_api"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL" env-default:"debug"`
	}

	Etherscan struct {
		Key string `env:"ETHERSCAN_APIKEY" env-default:"ABIJ85IE2BUBKU24XGXPTKDZMJ7XU8ZJCU"`
	}

	Mongo struct {
		Name     string `env:"MONGODB_NAME" env-default:"mongodb"`
		User     string `env:"MONGODB_USER" env-default:"bogdansemkin"`
		Password string `env:"MONGODB_PASSWORD" env-default:"Omebem04"`
		DBname   string `env:"MONGODB_DBNAME" env-default:"etherenum-api"`
		//Host string `env:"MONGODB_HOST" env-default:"localhost"`
		//Port string `env:"MONGODB_PORT" env-default:"27017"`
	}

	HTTP struct {
		Port                       string `env:"HTTP_PORT" env-default:"8080"`
		SendDetailsOnInternalError bool   `env:"HTTP_SEND_DETAILS_ON_INTERNAL_ERROR" env-default:"true"`
	}

	RabbitMQ struct {
		User     string `env:"RABBITMQ_USER" env-default:"guest"`
		Password string `env:"RABBITMQ_PASSWORD" env-default:"guest"`
		Host     string `env:"RABBITMQ_HOST" env-default:"localhost"`
		Port     string `env:"RABBITMQ_USER" env-default:"5672"`
	}
)
