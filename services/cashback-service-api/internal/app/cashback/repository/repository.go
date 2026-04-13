//go:generate mockery --all --case=snake --disable-version-string --with-expecter

package repository

import (
	"gorm.io/gorm"
)

type (
	OutboxWriter interface {
		CreateWithTx(tx *gorm.DB, eventType, aggregateType string, aggregateID int64, payload []byte) error
	}

	Repository struct {
		db           *gorm.DB
		outboxWriter OutboxWriter
	}
)

func New(db *gorm.DB, ow OutboxWriter) Repository {
	return Repository{
		db:           db,
		outboxWriter: ow,
	}
}
