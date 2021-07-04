package database

import (
	"fmt"
	_ "github.com/ppdraga/go-shortener/settings"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"time"
)

type R struct {
	DB *gorm.DB
	//conn   *sql.DB
}

func InitDB(logger *logrus.Logger) (*R, error) {
	host := os.Getenv("PG_HOST")
	user := os.Getenv("PG_USER")
	password := os.Getenv("PG_PASSWD")
	dbname := os.Getenv("PG_DB")
	port := os.Getenv("PG_PORT")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Moscow",
		host, user, password, dbname, port)
	logger.Info(dsn)

	var dbcon *gorm.DB
	var err error
	for _, attempt := range []int{1, 2, 3} {
		logger.Infof("Connecting to DB, attempt %d...", attempt)
		dbcon, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			logger.Info("Connected!!!")
			break
		}
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		logger.Error("Can't connect to DB")
		return nil, err
	}

	dbcon.Exec(`
		-- DROP TABLE IF EXISTS links;
		CREATE TABLE IF NOT EXISTS links (
            id serial PRIMARY KEY,
            created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
            resource TEXT NOT NULL,
            short_link VARCHAR(255) NOT NULL,
            short_link_num BIGINT NOT NULL DEFAULT 0,
            custom_name VARCHAR(255) NULL
        );
        DROP INDEX IF EXISTS short_link_idx;
        DROP INDEX IF EXISTS short_link_num_idx;
        CREATE INDEX short_link_idx ON links (short_link);
        CREATE INDEX short_link_num_idx ON links (short_link_num);

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
