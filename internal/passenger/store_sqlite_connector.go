package passenger

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"sync"
	"time"
)

const (
	maxActiveConnections = 1
	connectionsKeepAlive = 10 * time.Second
)

type Connector interface {
	Get() (*gorm.DB, error)
}

type connector struct {
	dbPath string
	db     *gorm.DB
	mu     sync.Mutex
}

func (c *connector) Get() (*gorm.DB, error) {
	if c.db == nil {
		c.mu.Lock()
		defer c.mu.Unlock()

		// initialize custom logger (for more info: https://gorm.io/docs/logger.html)
		sqliteLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold:             200 * time.Millisecond, // Slow SQL threshold
				LogLevel:                  logger.Error,           // Error Log level
				IgnoreRecordNotFoundError: true,                   // Ignore ErrRecordNotFound error for logger
				Colorful:                  true,                   // Enable colors
			},
		)

		db, err := gorm.Open(sqlite.Open(c.dbPath), &gorm.Config{Logger: sqliteLogger})
		if err != nil {
			return nil, err
		}

		// set connection pool
		sqlDB, err := db.DB()
		if err != nil {
			return nil, err
		}
		sqlDB.SetMaxOpenConns(maxActiveConnections)
		sqlDB.SetConnMaxLifetime(connectionsKeepAlive)

		c.db = db
	}

	return c.db, nil
}

func NewConnector(dbPath string) Connector {
	return &connector{dbPath: dbPath}
}
