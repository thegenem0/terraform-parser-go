package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type IDBService interface {
	Close() error
	AutoMigrate() error
	Connection() *gorm.DB
}

type DBService struct {
	connection *gorm.DB
}

func NewDBService() (*DBService, error) {
	connString := "host=localhost user=terraspect_root password=SuperSecretPassword dbname=terraspect_db port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  connString,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &DBService{
		connection: db,
	}, nil
}

func (dbs *DBService) Close() error {
	sqlDB, err := dbs.connection.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (dbs *DBService) AutoMigrate() error {
	for _, model := range []interface{}{&Plan{}, &State{}} {
		err := dbs.connection.AutoMigrate(model)
		if err != nil {
			return err
		}
	}
	return nil
}

func (dbs *DBService) Connection() *gorm.DB {
	return dbs.connection
}
