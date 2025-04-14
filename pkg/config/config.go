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

type SlackConfig struct {
	BotToken string
}

type EmailConfig struct {
	SendGridAPIKey string
	FromAddress    string
	FromName       string
	DefaultSubject string
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
	RabbitMQ RabbitMQConfig
	Twilio   TwilioConfig
	Slack    SlackConfig
	Email    EmailConfig
	Retry    RetryConfig
	UseMockProviders bool
}

func LoadConfigFromFile(filename string) (*Config, error) {
	err := godotenv.Load(filename)
	if err != nil {
		log.Fatal("Error loading .env file")
		return nil, err
	}

	return InitConfigValue(), nil
}


func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		return nil, err
	}

	return InitConfigValue(), nil
}

func InitConfigValue()  *Config  {
	maxRetries, _ := strconv.Atoi(os.Getenv("MAX_RETRY_ATTEMPTS"))
	initialDelayMs, _ := strconv.Atoi(os.Getenv("INITIAL_RETRY_DELAY_MS"))
	maxDelayMs, _ := strconv.Atoi(os.Getenv("MAX_RETRY_DELAY_MS"))
	processTimeout, _ := strconv.Atoi(os.Getenv("PROCESS_TIMEOUT_SECONDS"))

	maxOpenConns, _ := strconv.Atoi(os.Getenv("DB_MAX_OPEN_CONNS"))
	maxIdleConns, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNS"))
	connMaxLifetime, _ := strconv.Atoi(os.Getenv("DB_CONN_MAX_LIFETIME_MINUTES"))
	queryTimeout, _ := strconv.Atoi(os.Getenv("DB_QUERY_TIMEOUT_SECONDS"))

	requestTimeout, _ := strconv.Atoi(os.Getenv("REQUEST_TIMEOUT_SECONDS"))

	serverConfig := ServerConfig{
		Port:           os.Getenv("SERVER_PORT"),
		Host:           os.Getenv("SERVER_HOST"),
		RequestTimeout: requestTimeout,
	}

	dbConfig := DatabaseConfig{
		Host:            os.Getenv("DB_HOST"),
		Port:            os.Getenv("DB_PORT"),
		User:            os.Getenv("DB_USER"),
		Password:        os.Getenv("DB_PASSWORD"),
		Name:            os.Getenv("DB_NAME"),
		MaxOpenConns:    maxOpenConns,
		MaxIdleConns:    maxIdleConns,
		ConnMaxLifetime: connMaxLifetime,
		QueryTimeout:    queryTimeout,
	}

	rabbitMQConfig := RabbitMQConfig{
		Host:     os.Getenv("RABBITMQ_HOST"),
		Port:     os.Getenv("RABBITMQ_PORT"),
		User:     os.Getenv("RABBITMQ_USER"),
		Password: os.Getenv("RABBITMQ_PASS"),
		ChannelQueues: map[model.NotificationChannel]string {
			model.ChannelSMS: os.Getenv("RABBITMQ_SMS_QUEUE"),
			model.ChannelEmail: os.Getenv("RABBITMQ_EMAIL_QUEUE"),
			model.ChannelSlack: os.Getenv("RABBITMQ_SLACK_QUEUE"),
		},
		DLQPrefix: os.Getenv("RABBITMQ_DLQ_PREFIX"),
	}

	twilioConfig := TwilioConfig{
		AccountSID: os.Getenv("TWILIO_ACCOUNT_SID"),
		AuthToken:  os.Getenv("TWILIO_AUTH_TOKEN"),
		FromNumber: os.Getenv("TWILIO_FROM_NUMBER"),
	}

	slackConfig := SlackConfig{
		BotToken: os.Getenv("SLACK_BOT_TOKEN"),
	}

	emailConfig := EmailConfig{
		SendGridAPIKey: os.Getenv("SENDGRID_API_KEY"),
		FromAddress:    os.Getenv("SENDGRID_FROM_ADDRESS"),
		FromName:       os.Getenv("SENDGRID_FROM_NAME"),
		DefaultSubject: os.Getenv("EMAIL_DEFAULT_SUBJECT"),
	}

	retryConfig := RetryConfig{
		MaxRetries:     maxRetries,
		InitialDelayMs: initialDelayMs,
		MaxDelayMs:     maxDelayMs,
		ProcessTimeout: processTimeout,
	}

	useMockProviders, _ := strconv.ParseBool(os.Getenv("USE_MOCK_PROVIDERS"))

	return &Config{
		Server:   serverConfig,
		Database: dbConfig,
		RabbitMQ: rabbitMQConfig,
		Twilio:   twilioConfig,
		Slack:    slackConfig,
		Email:    emailConfig,
		Retry:    retryConfig,
		UseMockProviders: useMockProviders,
	}
}