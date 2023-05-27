package pgonthefly

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseOptions struct {
	AutomigrateFunc func(*DB) error
}

type DB struct {
	Conn   *gorm.DB
	Name   string
	Schema string
}

func NewDB(conn *gorm.DB, dbName, schema string) *DB {
	return &DB{
		Conn:   conn,
		Name:   dbName,
		Schema: schema,
	}
}

func GetSchemaConnection(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GenerateDSN(host, port, dbName, user, password string) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, dbName, port)
}

func CreateDatabase(dbName, host, port, user, password string, opts DatabaseOptions) (*DB, error) {
	dsn := GenerateDSN(host, port, dbName, user, password)

	defaultDB, err := GetSchemaConnection(dsn)
	if err != nil {
		return nil, err
	}

	randomString := uuid.New().String()
	databaseName := fmt.Sprintf("test-%s", randomString)

	if err := defaultDB.Exec(
		fmt.Sprintf(`CREATE DATABASE "%s";`, databaseName),
	).Error; err != nil {
		return nil, err
	}

	dsn = GenerateDSN(host, port, databaseName, user, password)

	testConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db := NewDB(testConn, databaseName, "public")

	if opts.AutomigrateFunc != nil {
		err := opts.AutomigrateFunc(db)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

func DeleteDatabase(dbName, host, port, user, password, target string) error {
	dsn := GenerateDSN(host, port, dbName, user, password)

	defaultDB, err := GetSchemaConnection(dsn)
	if err != nil {
		return err
	}

	if err := defaultDB.Exec(
		fmt.Sprintf(`DROP DATABASE "%s" WITH (FORCE)`, target),
	).Error; err != nil {
		return err
	}

	sqlDB, err := defaultDB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
