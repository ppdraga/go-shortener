package database

import (
	_ "github.com/ppdraga/go-shortener/settings"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type R struct {
	DB *gorm.DB
	//conn   *sql.DB
}

func InitDB() (*R, error) {
	dsn := "host=localhost user=postgres password=postgres dbname=postgres port=5432 sslmode=disable TimeZone=Europe/Moscow"
	dbcon, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("Can't connect to DB")
		return nil, err
	}

	dbcon.Exec(`
		CREATE TABLE IF NOT EXISTS links (
            id serial PRIMARY KEY,
            created_at date NOT NULL DEFAULT CURRENT_DATE,
            resource TEXT NOT NULL,
            short_link VARCHAR(255) NOT NULL,
            custom_name VARCHAR(255) NULL
        );
        DROP INDEX IF EXISTS short_link_idx;
        CREATE INDEX short_link_idx ON links (short_link);

`)
	return &R{DB: dbcon}, nil
}

func (r *R) Release() error {
	sqlDB, err := r.DB.DB()
	if err != nil {
		return err
	}
	err = sqlDB.Close()
	if err != nil {
		return err
	}
	return nil
}
