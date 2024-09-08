package db

import (
	"context"
	"errors"
	"github.com/IgorViskov/go_33_shortener/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"sync"
)

type Connector interface {
	IsConnected() bool
	TryConnected() bool
	GetConnection() *pgxpool.Pool
	GetError() error
	GetContext() context.Context
	Close() error
}

type connector struct {
	conn             *pgxpool.Pool
	isConnected      bool
	err              error
	mutex            sync.Mutex
	bdContext        context.Context
	connectionString string
}

func NewConnector(conf *config.AppConfig, contexts ...context.Context) Connector {
	bdContext := context.Background()
	if len(contexts) > 0 {
		bdContext = contexts[0]
	}
	return &connector{
		bdContext:        bdContext,
		connectionString: conf.ConnectionString,
	}
}

func (c *connector) IsConnected() bool {
	return c.isConnected
}
func (c *connector) GetConnection() *pgxpool.Pool {
	if c.conn == nil {
		c.mutex.Lock()
		if c.conn == nil {
			c.conn, c.err = c.connect()
			if c.err == nil {
				c.isConnected = true
			}
		}
		c.mutex.Unlock()
	}
	return c.conn
}

func (c *connector) GetError() error {
	return c.err
}

func (c *connector) Close() error {
	c.conn.Close()
	return nil
}

func (c *connector) connect() (*pgxpool.Pool, error) {
	if c.connectionString == "" {
		return nil, errors.New("no connection string provided")
	}
	return pgxpool.New(c.bdContext, c.connectionString)
}

func (c *connector) GetContext() context.Context {
	return c.bdContext
}

func (c *connector) TryConnected() bool {
	c.GetConnection()
	return c.IsConnected()
}
