package config

type ServerConfigProvider interface {
	GetServerEnv() string
	GetServerName() string
	GetServerHost() string
	GetServerPort() string
	GetServerBaseApiPrefix() string
}

type DatabaseConfigProvider interface {
	GetDBHost() string
	GetDBPort() string
	GetDBUser() string
	GetDBPassword() string
	GetDBName() string
	GetDBSSLMode() string
	GetDBTimezone() string
}

type AppConfigProvider interface {
	ServerConfigProvider
	DatabaseConfigProvider
}
