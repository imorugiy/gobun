package pgdriver

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"os"
	"time"
)

type Config struct {
	Network     string
	Addr        string
	DialTimeout time.Duration

	Dialer func(ctx context.Context, network, addr string) (net.Conn, error)

	TLSConfig *tls.Config

	User     string
	Password string
	Database string
	AppName  string

	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func newDefaultConfig() *Config {
	host := env("PGHOST", "localhost")
	port := env("PGPORT", "5432")

	cfg := &Config{
		Network:     "tcp",
		Addr:        net.JoinHostPort(host, port),
		DialTimeout: 5 * time.Second,
		TLSConfig:   nil, //&tls.Config{InsecureSkipVerify: true},

		User:     env("PGUSER", "postgres"),
		Database: env("PGDATABASE", "postgres"),

		ReadTimeout:  10 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	cfg.Dialer = func(ctx context.Context, network, addr string) (net.Conn, error) {
		netDialer := &net.Dialer{
			Timeout:   cfg.DialTimeout,
			KeepAlive: 5 * time.Minute,
		}
		return netDialer.DialContext(ctx, network, addr)
	}

	return cfg
}

func (c *Config) verify() error {
	if c.User == "" {
		return errors.New("pgdriver: User option is empty (to configure, use WithUser)")
	}
	return nil
}

func env(key, defaultValue string) string {
	if s := os.Getenv(key); s != "" {
		return s
	}
	return defaultValue
}
