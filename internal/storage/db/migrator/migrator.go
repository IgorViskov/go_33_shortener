package migrator

import (
	"context"
	"github.com/IgorViskov/go_33_shortener/internal/storage"
	"github.com/IgorViskov/go_33_shortener/internal/storage/db"
)

func AutoMigrate(connector db.Connector) error {
	return connector.GetConnection(context.Background()).AutoMigrate(&storage.Record{})
}
