package db

import (
	"context"
	"github.com/IgorViskov/go_33_shortener/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"sync"
)

type Connector interface {
	IsConnected() bool
	GetConnection() *gorm.DB
	GetError() error
	GetContext() context.Context
	Close() error
}

type connector struct {
	db               *gorm.DB
	err              error
	mutex            sync.Mutex
	dbContext        context.Context
	connectionString string
	state            ConnectionState
}

func NewConnector(conf *config.AppConfig, contexts ...context.Context) Connector {
	dbContext := context.Background()
	if len(contexts) > 0 {
		dbContext = contexts[0]
	}

	return &connector{
		dbContext:        dbContext,
		connectionString: conf.ConnectionString,
	}
}

func (c *connector) IsConnected() bool {
	switch c.state {
	case NotConnected:
		c.GetConnection()
		return c.IsConnected()
	case RefusedConnection:
	case InvalidConnectionString:
		return false
	case Connected:
		return true
	}
	return false
}
func (c *connector) GetConnection() *gorm.DB {
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
		Context: c.dbContext,
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
	return gorm.Open(postgres.Open(c.connectionString), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "main.",
			SingularTable: false,
		},
	})
}

func (c *connector) GetContext() context.Context {
	return c.dbContext
}
