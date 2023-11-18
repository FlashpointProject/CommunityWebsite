package config

import (
	"fmt"
	"os"
	"strconv"
)

type OauthConfig struct {
	ClientID           string `json:"client_id"`
	ClientSecret       string `json:"client_secret"`
	AuthorizeEndpoint  string `json:"authorize_endpoint"`
	ProfileEndpoint    string `json:"profile_endpoint"`
	TokenEndpoint      string `json:"token_endpoint"`
	Callback           string `json:"callback"`
	Scope              string `json:"scope"`
	FpfssClientID      string `json:"fpfss_client_id"`
	FpfssClientSecret  string `json:"fpfss_client_secret"`
	FpfssClientScope   string `json:"fpfss_client_scope"`
	FpfssTokenEndpoint string `json:"fpfss_token_endpoint"`
}

type AppConfig struct {
	Name                         string
	Port                         int64
	FpfssApiUrl                  string
	Version                      string
	OauthConfig                  *OauthConfig
	SessionExpirationSeconds     int64
	PostgresUser                 string
	PostgresPassword             string
	PostgresHost                 string
	PostgresPort                 int64
	HostBaseUrl                  string
	SecurecookieHashKeyPrevious  string
	SecurecookieBlockKeyPrevious string
	SecurecookieHashKeyCurrent   string
	SecurecookieBlockKeyCurrent  string
}

func GetConfig() (*AppConfig, error) {
	return &AppConfig{
		Name:        EnvString("APP_NAME"),
		Port:        EnvInt("APP_PORT"),
		FpfssApiUrl: EnvString("FPFSS_API_URL"),
		OauthConfig: &OauthConfig{
			ClientID:           EnvString("OAUTH_CLIENT_ID"),
			ClientSecret:       EnvString("OAUTH_CLIENT_SECRET"),
			AuthorizeEndpoint:  EnvString("OAUTH_AUTHORIZE_ENDPOINT"),
			ProfileEndpoint:    EnvString("OAUTH_PROFILE_ENDPOINT"),
			TokenEndpoint:      EnvString("OAUTH_TOKEN_ENDPOINT"),
			Callback:           EnvString("OAUTH_CALLBACK"),
			Scope:              EnvString("OAUTH_SCOPE"),
			FpfssClientID:      EnvString("OAUTH_FPFSS_CLIENT_ID"),
			FpfssClientSecret:  EnvString("OAUTH_FPFSS_CLIENT_SECRET"),
			FpfssClientScope:   EnvString("OAUTH_FPFSS_CLIENT_SCOPE"),
			FpfssTokenEndpoint: EnvString("OAUTH_FPFSS_TOKEN_ENDPOINT"),
		},
		SessionExpirationSeconds:     EnvInt("SESSION_EXPIRATION_SECONDS"),
		PostgresUser:                 EnvString("POSTGRES_USER"),
		PostgresPassword:             EnvString("POSTGRES_PASSWORD"),
		PostgresHost:                 EnvString("POSTGRES_HOST"),
		PostgresPort:                 EnvInt("POSTGRES_PORT"),
		HostBaseUrl:                  EnvString("HOST_BASE_URL"),
		SecurecookieHashKeyPrevious:  EnvString("SECURECOOKIE_HASH_KEY_PREVIOUS"),
		SecurecookieBlockKeyPrevious: EnvString("SECURECOOKIE_BLOCK_KEY_PREVIOUS"),
		SecurecookieHashKeyCurrent:   EnvString("SECURECOOKIE_HASH_KEY_CURRENT"),
		SecurecookieBlockKeyCurrent:  EnvString("SECURECOOKIE_BLOCK_KEY_CURRENT"),
	}, nil
}

func EnvString(name string) string {
	s := os.Getenv(name)
	if s == "" {
		panic(fmt.Sprintf("env variable '%s' is not set", name))
	}
	return s
}

func EnvInt(name string) int64 {
	s := os.Getenv(name)
	if s == "" {
		panic(fmt.Sprintf("env variable '%s' is not set", name))
	}
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return i
}

func EnvBool(name string) bool {
	s := os.Getenv(name)
	if s == "" {
		panic(fmt.Sprintf("env variable '%s' is not set", name))
	} else if s == "True" {
		return true
	} else if s == "False" {
		return false
	}
	panic(fmt.Sprintf("invalid value of env variable '%s'", name))
}
