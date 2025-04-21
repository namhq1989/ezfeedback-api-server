package config

import (
	"errors"
)

type (
	Server struct {
		RestPort string
		GRPCPort string

		AppName      string
		Environment  string
		IsEnvRelease bool
		SwaggerURL   string

		// Authentication
		AccessTokenSecret string
		AccessTokenTTL    int // seconds

		// Postgres
		PostgresConn string

		// Redis
		CachingRedisURL string

		// Queue
		QueueRedisURL    string
		QueueUsername    string
		QueuePassword    string
		QueueConcurrency int

		// R2
		R2AccessKey string
		R2SecretKey string
		R2Bucket    string
		R2Endpoint  string

		// Open Observe
		OpenObserveHttpEndpoint string
		OpenObserveStreamName   string
		OpenObserveToken        string

		// Sentry
		SentryDSN         string
		SentryMachineName string

		// Mailer
		MailerFromEmail string
		MailerService   string
		BrevoApiKey     string

		// Endpoint
		CDNEndpoint string
	}
)

func Init() Server {
	cfg := Server{
		RestPort: ":3000",
		GRPCPort: ":3001",

		AppName:     getEnvStr("APP_NAME"),
		Environment: getEnvStr("ENVIRONMENT"),
		SwaggerURL:  getEnvStr("SWAGGER_URL"),

		AccessTokenSecret: getEnvStr("ACCESS_TOKEN_SECRET"),
		AccessTokenTTL:    getEnvInt("ACCESS_TOKEN_TTL"),

		PostgresConn: getEnvStr("POSTGRES_CONN"),

		CachingRedisURL: getEnvStr("CACHING_REDIS_URL"),

		QueueRedisURL:    getEnvStr("QUEUE_REDIS_URL"),
		QueueUsername:    getEnvStr("QUEUE_USERNAME"),
		QueuePassword:    getEnvStr("QUEUE_PASSWORD"),
		QueueConcurrency: getEnvInt("QUEUE_CONCURRENCY"),

		R2AccessKey: getEnvStr("R2_ACCESS_KEY"),
		R2SecretKey: getEnvStr("R2_SECRET_KEY"),
		R2Bucket:    getEnvStr("R2_BUCKET"),
		R2Endpoint:  getEnvStr("R2_ENDPOINT"),

		OpenObserveHttpEndpoint: getEnvStr("OPEN_OBSERVE_HTTP_ENDPOINT"),
		OpenObserveStreamName:   getEnvStr("OPEN_OBSERVE_STREAM_NAME"),
		OpenObserveToken:        getEnvStr("OPEN_OBSERVE_TOKEN"),

		SentryDSN:         getEnvStr("SENTRY_DSN"),
		SentryMachineName: getEnvStr("SENTRY_MACHINE_NAME"),

		MailerFromEmail: getEnvStr("MAILER_FROM_EMAIL"),
		MailerService:   getEnvStr("MAILER_SERVICE"),
		BrevoApiKey:     getEnvStr("BREVO_API_KEY"),

		CDNEndpoint: getEnvStr("CDN_ENDPOINT"),
	}
	cfg.IsEnvRelease = cfg.Environment == "release"

	// validation
	if cfg.Environment == "" {
		panic(errors.New("missing ENVIRONMENT"))
	}

	if cfg.AccessTokenSecret == "" {
		panic(errors.New("missing ACCESS_TOKEN_SECRET"))
	}

	if cfg.PostgresConn == "" {
		panic(errors.New("missing POSTGRES_CONN"))
	}

	if cfg.CachingRedisURL == "" {
		panic(errors.New("missing CACHING_REDIS_URL"))
	}

	if cfg.QueueRedisURL == "" {
		panic(errors.New("missing QUEUE_REDIS_URL"))
	}

	if cfg.R2AccessKey == "" {
		panic(errors.New("missing R2_ACCESS_KEY"))
	}

	if cfg.MailerService == "" {
		panic(errors.New("missing MAILER_SERVICE"))
	}

	if cfg.BrevoApiKey == "" {
		panic(errors.New("missing BREVO_API_KEY"))
	}

	if cfg.CDNEndpoint == "" {
		panic(errors.New("missing CDN_ENDPOINT"))
	}

	return cfg
}
