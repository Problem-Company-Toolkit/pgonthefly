package pgonthefly

import (
	"fmt"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type DatabaseOptions struct {
	DbName          string
	DbHost          string
	DbPort          string
	DbUser          string
	DbPassword      string
	DbSchema        string
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
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable&TimeZone=UTC", user, password, host, port, dbName)
}

func CreateDatabase(opts DatabaseOptions) (*DB, error) {
	dsn := GenerateDSN(opts.DbHost, opts.DbPort, opts.DbName, opts.DbUser, opts.DbPassword)

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

	dsn = GenerateDSN(opts.DbHost, opts.DbPort, databaseName, opts.DbUser, opts.DbPassword)

	dbSchema := opts.DbSchema

	if dbSchema == "" {
		dbSchema = "public"
	}

	tablePrefix := fmt.Sprintf("%s.", dbSchema)

	testConn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: tablePrefix,
		},
	})
	if err != nil {
		return nil, err
	}

	db := NewDB(testConn, databaseName, dbSchema)

	if opts.AutomigrateFunc != nil {
		err := opts.AutomigrateFunc(db)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

type DeleteDatabaseOpts struct {
	DbName     string
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassword string
	Target     string
}

func DeleteDatabase(opts DeleteDatabaseOpts) error {
	dsn := GenerateDSN(opts.DbHost, opts.DbPort, opts.DbName, opts.DbUser, opts.DbPassword)

	defaultDB, err := GetSchemaConnection(dsn)
	if err != nil {
		return err
	}

	if err := defaultDB.Exec(
		fmt.Sprintf(`DROP DATABASE "%s" WITH (FORCE)`, opts.Target),
	).Error; err != nil {
		return err
	}

	sqlDB, err := defaultDB.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}
