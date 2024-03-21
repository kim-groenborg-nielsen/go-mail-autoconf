package config

import (
	"os"
	"strconv"
)
import _ "github.com/joho/godotenv/autoload"

var cachedConfig *Config

type Config struct {
	ImapServer          string
	ImapPort            int
	ImapSsl             bool
	Pop3Server          string
	Pop3Port            int
	Pop3Ssl             bool
	SmtpServer          string
	SmtpPort            int
	SmtpSsl             bool
	SmtpGlobalPreferred bool
	EmailProvider       string
	ServerUrl           string
}

func GetEnv(key, fallback string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return value
}

func GetEnvInt(key string, fallback int) int {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return i
}

func GetEnvBool(key string, fallback bool) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	b, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return b
}

func NewConfig() *Config {
	if cachedConfig != nil {
		return cachedConfig
	}
	cachedConfig = &Config{
		ImapServer:          GetEnv("IMAP_SERVER", "imap.dummy.org"),
		ImapPort:            GetEnvInt("IMAP_PORT", 993),
		ImapSsl:             GetEnvBool("IMAP_SSL", true),
		Pop3Server:          GetEnv("POP3_SERVER", ""),
		Pop3Port:            GetEnvInt("POP3_PORT", 995),
		Pop3Ssl:             GetEnvBool("POP3_SSL", true),
		SmtpServer:          GetEnv("SMTP_SERVER", "smtp.dummy.org"),
		SmtpPort:            GetEnvInt("SMTP_PORT", 587),
		SmtpSsl:             GetEnvBool("SMTP_SSL", false),
		SmtpGlobalPreferred: GetEnvBool("SMTP_GLOBAL_PREFERRED", false),
		EmailProvider:       GetEnv("EMAIL_PROVIDER", "dummy.org"),
		ServerUrl:           GetEnv("SERVER_URL", "127.0.0.1:3000"),
	}
	return cachedConfig
}
