package withdb

import (
	"github.com/ppdraga/go-shortener/internal/shortener/link/datatype"
	"gorm.io/gorm"
)

type WithDB struct {
	db *gorm.DB
}

func New(db *gorm.DB) *WithDB {
	return &WithDB{
		db: db,
	}
}

func (wdb *WithDB) ReadLink(id int64) (*datatype.Link, error) {

	return &datatype.Link{}, nil
}

func (wdb *WithDB) WriteLink(external *datatype.Link) error {

	return nil
}
