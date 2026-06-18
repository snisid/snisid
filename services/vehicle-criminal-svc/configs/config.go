package configs

import (
	"os"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string

	RedisAddr     string
	RedisPassword string

	Neo4jURI      string
	Neo4jUser     string
	Neo4jPassword string

	ClickHouseAddr string
	ClickHouseDB   string

	KafkaBrokers string
	KafkaGroupID string

	InterpolGatewayURL string
	InterpolAPIKey     string
	InterpolNCBCode    string

	FovesServiceURL   string
	FovesServiceToken string

	ServicePort string
	LogLevel    string
	Env         string
}

func Load() *Config {
	return &Config{
		DBHost:     getEnv("SIVC_DB_HOST", "localhost"),
		DBPort:     getEnv("SIVC_DB_PORT", "5432"),
		DBName:     getEnv("SIVC_DB_NAME", "snisid_sivc"),
		DBUser:     getEnv("SIVC_DB_USER", "sivc_svc"),
		DBPassword: getEnv("SIVC_DB_PASSWORD", ""),

		RedisAddr:     getEnv("SIVC_REDIS_ADDR", "redis-master:6379"),
		RedisPassword: getEnv("SIVC_REDIS_PASSWORD", ""),

		Neo4jURI:      getEnv("SIVC_NEO4J_URI", "bolt://neo4j:7687"),
		Neo4jUser:     getEnv("SIVC_NEO4J_USER", "neo4j"),
		Neo4jPassword: getEnv("SIVC_NEO4J_PASSWORD", ""),

		ClickHouseAddr: getEnv("SIVC_CLICKHOUSE_ADDR", "clickhouse:9000"),
		ClickHouseDB:   getEnv("SIVC_CLICKHOUSE_DB", "snisid_analytics"),

		KafkaBrokers: getEnv("SIVC_KAFKA_BROKERS", "kafka:9092"),
		KafkaGroupID: getEnv("SIVC_KAFKA_GROUP_ID", "sivc-consumer-group"),

		InterpolGatewayURL: getEnv("INTERPOL_GATEWAY_URL", "https://i247-gateway.pnh.gov.ht/api"),
		InterpolAPIKey:     getEnv("INTERPOL_API_KEY", ""),
		InterpolNCBCode:    getEnv("INTERPOL_NCB_CODE", "HTI"),

		FovesServiceURL:   getEnv("FOVES_SERVICE_URL", "http://foves-svc:8080"),
		FovesServiceToken: getEnv("FOVES_SERVICE_TOKEN", ""),

		ServicePort: getEnv("SIVC_SERVICE_PORT", "8090"),
		LogLevel:    getEnv("SIVC_LOG_LEVEL", "info"),
		Env:         getEnv("SIVC_ENV", "production"),
	}
}

func (c *Config) DatabaseDSN() string {
	return "host=" + c.DBHost + " port=" + c.DBPort + " dbname=" + c.DBName +
		" user=" + c.DBUser + " password=" + c.DBPassword + " sslmode=disable"
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
