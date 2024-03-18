package common

import (
	"database/sql"
	goora "github.com/sijms/go-ora/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

type DatabaseConfig struct {
	Host               string
	Port               int
	Service            string
	User               string
	Pwd                string
	MaxIdleConnections int
	MaxOpenConnections int
	MaxLifetime        time.Duration
}

func NewSQLDB(cfg *DatabaseConfig) (*sql.DB, error) {
	connStr := goora.BuildUrl(cfg.Host, cfg.Port, cfg.Service, cfg.User, cfg.Pwd, nil)
	sqlDB, err := sql.Open("oracle", connStr)
	if err != nil {
		logrus.Error(err)
		return nil, err
	} else {
		err := sqlDB.Ping()
		if err != nil {
			logrus.Fatal("Failed to connect to Oracle: ", err)
		} else {
			logrus.Info("Connected to database")
		}
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConnections)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConnections)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(cfg.MaxLifetime)

	return sqlDB, nil
}

func ProvideDatabaseConfig() *DatabaseConfig {
	dbOption := &DatabaseConfig{
		Host:               viper.GetString("DB_HOST"),
		Port:               viper.GetInt("DB_PORT"),
		Service:            viper.GetString("DB_SERVICE"),
		User:               viper.GetString("DB_USER"),
		Pwd:                viper.GetString("DB_PWD"),
		MaxIdleConnections: viper.GetInt("MAX_IDLE_CONNECTION"),
		MaxOpenConnections: viper.GetInt("MAX_OPEN_CONNECTION"),
		MaxLifetime:        time.Duration(viper.GetInt("MAX_LIFETIME")) * time.Second,
	}
	return dbOption
}
