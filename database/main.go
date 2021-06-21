package database

import (
	"fmt"
	_ "github.com/ppdraga/go-shortener/settings"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"time"
)

type R struct {
	DB *gorm.DB
	//conn   *sql.DB
}

func InitDB() (*R, error) {
	host := os.Getenv("PG_HOST")
	user := os.Getenv("PG_USER")
	password := os.Getenv("PG_PASSWD")
	dbname := os.Getenv("PG_DB")
	port := os.Getenv("PG_PORT")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Moscow",
		host, user, password, dbname, port)
	log.Print(dsn)

	var dbcon *gorm.DB
	var err error
	for _, attempt := range []int{1, 2, 3} {
		log.Info().Msgf("Connecting to DB, attempt %d...", attempt)
		dbcon, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Info().Msg("Connected!!!")
			break
		}
		time.Sleep(3 * time.Second)
	}
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
