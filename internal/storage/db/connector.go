package db

import (
	"context"
	"github.com/IgorViskov/go_33_shortener/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
)

type Connector interface {
	IsConnected() bool
	GetConnection(context context.Context) *gorm.DB
	GetError() error
	Close() error
}

type connector struct {
	db               *gorm.DB
	err              error
	mutex            sync.Mutex
	connectionString string
	state            ConnectionState
}

func NewConnector(conf *config.AppConfig) Connector {
	return &connector{
		connectionString: conf.ConnectionString,
	}
}

func (c *connector) IsConnected() bool {
	switch c.state {
	case NotConnected:
		c.GetConnection(context.Background())
		return c.IsConnected()
	case RefusedConnection:
	case InvalidConnectionString:
		return false
	case Connected:
		return true
	}
	return false
}
func (c *connector) GetConnection(context context.Context) *gorm.DB {
	if c.db == nil {
		c.mutex.Lock()
		if c.db == nil {
			c.db, c.err = c.connect()
			if c.err == nil {
				c.state = Connected
			} else if c.state != InvalidConnectionString {
				c.state = RefusedConnection
			}
		}
		c.mutex.Unlock()
	}
	return c.db.Session(&gorm.Session{
		Context: context,
	})
}

func (c *connector) GetError() error {
	return c.err
}

func (c *connector) Close() error {
	db, err := c.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func (c *connector) connect() (*gorm.DB, error) {
	_, err := pgxpool.ParseConfig(c.connectionString)
	if err != nil {
		c.err = err
		c.state = InvalidConnectionString
		return nil, err
	}
	return gorm.Open(postgres.Open(c.connectionString))
}

func (c *connector) getContext(contexts []context.Context) context.Context {
	if len(contexts) > 0 {
		return contexts[0]
	}
	return context.Background()
}
