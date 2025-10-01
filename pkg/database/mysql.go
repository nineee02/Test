package database

import (
	"fmt"
	"sync"
	"time"

	"github.com/nineee02/gotest/pkg/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	mysqlDB     *gorm.DB
	mysqlDBOnce sync.Once
)

func NewMySQLDB(cfg *config.Configuration) (*gorm.DB, error) {
	var err error
	mysqlDBOnce.Do(func() {
		mysqlDB, err = connectMySQL(cfg)
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := mysqlDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql database: %w", err)
	}
	if err = sqlDB.Ping(); err != nil {
		mysqlDB, err = reconnectMySQL(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to reconnect to mysql database: %w", err)
		}
	}

	return mysqlDB, nil
}

func connectMySQL(cfg *config.Configuration) (*gorm.DB, error) {
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.MySQL.User,
		cfg.MySQL.Password,
		cfg.MySQL.Host,
		cfg.MySQL.Port,
		cfg.MySQL.DBname,
	)
	dial := mysql.Open(dataSourceName)
	option := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   cfg.MySQL.TablePrefix,
			SingularTable: true,
		},
	}

	db, err := gorm.Open(dial, option)
	if err != nil {
		return nil, fmt.Errorf("failed to open mysql database: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql database: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MySQL.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MySQL.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.MySQL.ConnMaxLifetime)

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping mysql database: %w", err)
	}

	return db, nil
}

func reconnectMySQL(cfg *config.Configuration) (*gorm.DB, error) {
	maxRetries := 5
	retryDelay := 2 * time.Second
	for i := 0; i < maxRetries; i++ {
		db, err := connectMySQL(cfg)
		if err == nil {
			return db, nil
		}
		time.Sleep(retryDelay)
		retryDelay *= 2
	}
	return nil, fmt.Errorf("failed to reconnect to MySQL after %d attempts", maxRetries)
}
