package config

import (
	"fmt"
	"net/url"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ConnectionInfo struct {
	Host     string     `env:"HOST" default:"localhost" help:"The IP or FQDN of the Postgres host"`
	Port     uint16     `env:"PORT" default:"5432" help:"The Postgres Port"`
	Name     string     `env:"NAME" help:"The DB Name"`
	User     string     `env:"USER" help:"The DB Username"`
	Password string     `env:"PASSWORD" help:"The DB PASSWORD"`
	SSLMode  string     `env:"SSL_MODE" enum:"disable,allow,prefer,require,verify-ca,verify-full" default:"prefer" help:"The SSL Connection Mode"`
	Extra    url.Values `env:"EXTRA_CONNECTION_PARAMS" optional:"" help:"Extra connection parameters"`
}

func (c *ConnectionInfo) URL() string {
	user := url.QueryEscape(c.User)
	pass := url.QueryEscape(c.Password)
	name := url.QueryEscape(c.Name)
	if c.Extra == nil {
		c.Extra = make(url.Values)
	}
	c.Extra.Set("ssl-mode", c.SSLMode)

	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s", user, pass, c.Host, c.Port, name, c.Extra.Encode())
}

type DatabaseConfigFragment struct {
	DatabaseURL string          `required:"" group:"Connection Info" xor:"Connection Info" env:"POSTGRES_PLUGIN_DATABASE_URL" help:"The Postgres Connection Info in postgres://(user):(password)@host:[port]j/dbname[?...] form"`
	DSN         string          `required:"" group:"Connection Info" xor:"Connection Info" env:"POSTGRES_PLUGIN_DATABASE_DSN" help:"The Postgres Connection Info in 'user=... password=... ...' form"`
	DBInfo      *ConnectionInfo `required:"" group:"Connection Info" xor:"Connection Info" envprefix:"POSTGRES_PLUGIN_DATABASE_" prefix:"db-connection." embed:"" help:"The Postgres Connection Info in composite format"`
}

func (d DatabaseConfigFragment) OpenDB() (*gorm.DB, error) {
	connString := d.DatabaseURL
	if connString == "" {
		connString = d.DSN
	}
	if connString == "" {
		connString = d.DBInfo.URL()
	}

	return gorm.Open(postgres.Open(connString))
}
