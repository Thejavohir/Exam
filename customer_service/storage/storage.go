package storage

import (
	"github.com/exam/customer_service/storage/postgres"
	"github.com/exam/customer_service/storage/repo"

	"github.com/jmoiron/sqlx"
)

type IStorage interface {
	Customer() repo.CustomerStorageI
}

type storagePg struct {
	db           *sqlx.DB
	customerRepo repo.CustomerStorageI
}

func NewStoragePg(db *sqlx.DB) *storagePg {
	return &storagePg{
		db:           db,
		customerRepo: postgres.NewCustomerRepo(db),
	}
}

func (s storagePg) Customer() repo.CustomerStorageI {
	return s.customerRepo
}
