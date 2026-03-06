package config

import (
	"fmt"
	"net/url"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	URL      string     `env:"URL" help:"The Postgres Connection Info in postgres://(user):(password)@host:[port]j/dbname[?...] form"`
	Host     string     `env:"HOST" default:"localhost" help:"The IP or FQDN of the Postgres host"`
	Port     uint16     `env:"PORT" default:"5432" help:"The Postgres Port"`
	Name     string     `env:"NAME" default:"coredns" help:"The DB Name"`
	User     string     `env:"USER" help:"The DB Username"`
	Password string     `env:"PASSWORD" help:"The DB Password"`
	SSLMode  string     `env:"SSL_MODE" enum:"disable,allow,prefer,require,verify-ca,verify-full" default:"prefer" help:"The SSL Connection Mode"`
	Extra    url.Values `env:"EXTRA_CONNECTION_PARAMS" help:"Extra connection parameters"`
}

func (d DatabaseConfig) ConnectionString() string {
	if d.URL != "" {
		return d.URL
	}

	user := url.QueryEscape(d.User)
	pass := url.QueryEscape(d.Password)
	name := url.QueryEscape(d.Name)
	if d.Extra == nil {
		d.Extra = make(url.Values)
	}
	d.Extra.Set("ssl-mode", d.SSLMode)

	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s", user, pass, d.Host, d.Port, name, d.Extra.Encode())
}

func (d DatabaseConfig) OpenDB() (*gorm.DB, error) {
	return gorm.Open(postgres.Open(d.ConnectionString()))
}
