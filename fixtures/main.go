package fixtures

import (
	"github.com/ppdraga/go-shortener/database"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitTestSQLite(logger *logrus.Logger) (*database.R, error) {
	dbcon, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		logger.Errorf("Can't connect to DB: %v", err)
	}
	ApplyMigrations(dbcon)
	return &database.R{DB: dbcon}, nil
}

func ApplyMigrations(db *gorm.DB) {
	db.Exec(`
		CREATE TABLE IF NOT EXISTS links (
            id serial PRIMARY KEY,
            created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
            resource TEXT NOT NULL,
            short_link VARCHAR(255) NOT NULL,
            short_link_num BIGINT NOT NULL DEFAULT 0,
            custom_name VARCHAR(255) NULL
        );
        CREATE INDEX short_link_idx ON links (short_link);
        CREATE INDEX short_link_num_idx ON links (short_link_num);
    `)

}
