package config

import (
	"log"
	"os"
	"strconv"

	"notification-system/pkg/model"

	"github.com/joho/godotenv"
)

type ServerConfig struct {
	Port            string
	Host            string
	RequestTimeout  int // in seconds
}

type DatabaseConfig struct {
	Host              string
	Port              string
	User              string
	Password          string
	Name              string
	MaxOpenConns      int
	MaxIdleConns      int
	ConnMaxLifetime   int // in minutes
	QueryTimeout      int // in seconds
}

type RabbitMQConfig struct {
	Host          string
	Port          string
	User          string
	Password      string
	Queue         string
	ChannelQueues map[model.NotificationChannel]string
	DLQPrefix     string
}

type TwilioConfig struct {
	AccountSID string
	AuthToken  string
	FromNumber string
}

type RetryConfig struct {
	MaxRetries      int
	InitialDelayMs  int
	MaxDelayMs      int
	ProcessTimeout  int // in seconds
}

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Twilio   TwilioConfig
	RabbitMQ RabbitMQConfig
	Retry    RetryConfig
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return nil, err
	}

	maxRetries, _ := strconv.Atoi(os.Getenv("MAX_RETRY_ATTEMPTS"))
	initialDelayMs, _ := strconv.Atoi(os.Getenv("INITIAL_RETRY_DELAY_MS"))
	maxDelayMs, _ := strconv.Atoi(os.Getenv("MAX_RETRY_DELAY_MS"))
	processTimeout, _ := strconv.Atoi(os.Getenv("PROCESS_TIMEOUT_SECONDS"))

	maxOpenConns, _ := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNS"))
	maxIdleConns, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNS"))
	connMaxLifetime, _ := strconv.Atoi(os.Getenv("DB_CONN_MAX_LIFETIME_MINUTES"))
	queryTimeout, _ := strconv.Atoi(os.Getenv("DB_QUERY_TIMEOUT_SECONDS"))

	requestTimeout, _ := strconv.Atoi(os.Getenv("REQUEST_TIMEOUT_SECONDS"))

	// Return the config struct populated with values from environment variables
	return &Config{
		Server: ServerConfig{
			Port:           os.Getenv("SERVER_PORT"),
			Host:           os.Getenv("SERVER_HOST"),
			RequestTimeout: requestTimeout,
		},
		Database: DatabaseConfig{
			Host:            os.Getenv("DB_HOST"),
			Port:            os.Getenv("DB_PORT"),
			User:            os.Getenv("DB_USER"),
			Password:        os.Getenv("DB_PASSWORD"),
			Name:            os.Getenv("DB_NAME"),
			MaxOpenConns:    maxOpenConns,
			MaxIdleConns:    maxIdleConns,
			ConnMaxLifetime: connMaxLifetime,
			QueryTimeout:    queryTimeout,
		},
		RabbitMQ: RabbitMQConfig{
			Host:     os.Getenv("RABBITMQ_HOST"),
			Port:     os.Getenv("RABBITMQ_PORT"),
			User:     os.Getenv("RABBITMQ_USER"),
			Password: os.Getenv("RABBITMQ_PASS"),
			ChannelQueues:    map[model.NotificationChannel]string {
				model.ChannelSMS: os.Getenv("RABBITMQ_SMS_QUEUE"),
				model.ChannelEmail: os.Getenv("RABBITMQ_EMAIL_QUEUE"),
				model.ChannelSlack: os.Getenv("RABBITMQ_SLACK_QUEUE"),
			},
			DLQPrefix: os.Getenv("RABBITMQ_DLQ_PREFIX"),
		},
		Twilio: TwilioConfig{
			AccountSID: os.Getenv("TWILIO_ACCOUNT_SID"),
			AuthToken:  os.Getenv("TWILIO_AUTH_TOKEN"),
			FromNumber: os.Getenv("TWILIO_FROM_NUMBER"),
		},
		Retry: RetryConfig{
			MaxRetries:     maxRetries,
			InitialDelayMs: initialDelayMs,
			MaxDelayMs:     maxDelayMs,
			ProcessTimeout: processTimeout,
		},
	}, nil
}