package repository

import (
	"gorm.io/gorm"
)

type (
	outboxWriter interface {
		CreateWithTx(tx *gorm.DB, eventType, aggregateType string, aggregateID int64, payload []byte) error
	}

	Repository struct {
		db           *gorm.DB
		outboxWriter outboxWriter
	}
)

func New(db *gorm.DB, ow outboxWriter) Repository {
	return Repository{
		db:           db,
		outboxWriter: ow,
	}
}
