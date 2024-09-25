package migrator

import (
	"context"
	"github.com/IgorViskov/go_33_shortener/internal/storage"
	"github.com/IgorViskov/go_33_shortener/internal/storage/db"
)

func AutoMigrate(connector db.Connector) error {
	session := connector.GetConnection(context.Background())
	err := session.AutoMigrate(&storage.User{}, &storage.Record{})
	if err != nil {
		return err
	}

	return session.Exec(CreatePartialIndex).Error
}
