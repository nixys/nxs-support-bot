package primedb

import (
	"fmt"

	gmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB it is a DB module context structure
type DB struct {
	client *gorm.DB
}

// Settings contains settings for DB
type Settings struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

// Connect connects to DB
func Connect(s Settings) (DB, error) {

	client, err := gorm.Open(gmysql.Open(fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		s.User,
		s.Password,
		s.Host,
		s.Port,
		s.Database)), &gorm.Config{})
	if err != nil {
		return DB{}, err
	}

	return DB{
		client: client,
	}, nil
}

// Close closes DB connection
func (db *DB) Close() error {
	c, _ := db.client.DB()
	return c.Close()
}
