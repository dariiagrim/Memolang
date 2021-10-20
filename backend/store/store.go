package store

import (
	"fmt"
	"github.com/go-gormigrate/gormigrate/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"memolang/configuration"
	"memolang/migrations"
)

type Store struct {
	config *configuration.StoreConfig
	Db     *gorm.DB
}

func NewStore(config *configuration.StoreConfig) *Store {
	return &Store{
		config: config,
	}
}

func (s *Store) Open() error {
	db, err := gorm.Open(postgres.Open(fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v",  s.config.Host, s.config.User, s.config.Password, s.config.Name, s.config.Port)), &gorm.Config{})
	if err != nil {
		return err
	}
	m := gormigrate.New(db, gormigrate.DefaultOptions, migrations.GetMigrations())

	if err = m.Migrate(); err != nil {
		return err
	}
	s.Db = db
	return nil
}

