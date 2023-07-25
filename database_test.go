package pgonthefly_test

import (
	"fmt"
	"strings"

	"github.com/brianvoe/gofakeit/v6"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	pgonthefly "github.com/problem-company-toolkit/pgonthefly"
)

var _ = Describe("pgonthefly", func() {
	Describe("Creating and deleting a database", func() {
		It("creates a database successfully", func() {
			db, err := pgonthefly.CreateDatabase(pgonthefly.DatabaseOptions{
				DbName:     dbName,
				DbHost:     host,
				DbPort:     port,
				DbUser:     user,
				DbPassword: password,
				DbSchema:   "public",
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(db).NotTo(BeNil())

			// Test that the database connection is working
			err = db.Conn.Transaction(func(tx *gorm.DB) error {
				// This is just a dummy transaction for testing purposes
				return nil
			})
			Expect(err).NotTo(HaveOccurred())

			// Test database deletion
			err = pgonthefly.DeleteDatabase(pgonthefly.DeleteDatabaseOpts{
				DbName:     dbName,
				DbHost:     host,
				DbPort:     port,
				DbUser:     user,
				DbPassword: password,
				Target:     db.Name,
			})
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("Using a custom database schema", func() {
		type Example struct {
			Value string `gorm:"not null"`
		}

		var (
			db           *pgonthefly.DB
			randomSchema string
		)

		BeforeEach(func() {
			var err error
			randomSchema = strings.ToLower(gofakeit.Word())
			db, err = pgonthefly.CreateDatabase(pgonthefly.DatabaseOptions{
				DbName:     dbName,
				DbHost:     host,
				DbPort:     port,
				DbUser:     user,
				DbPassword: password,
				DbSchema:   randomSchema,
			})

			if err != nil {
				Fail(err.Error())
				return
			}

			if err := db.Conn.Exec(fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS "%s"`, randomSchema)).Error; err != nil {
				Fail(err.Error())
				return
			}

			if err := db.Conn.AutoMigrate(&Example{}); err != nil {
				Fail(err.Error())
				return
			}
		})

		AfterEach(func() {
			err := pgonthefly.DeleteDatabase(pgonthefly.DeleteDatabaseOpts{
				DbName:     dbName,
				DbHost:     host,
				DbPort:     port,
				DbUser:     user,
				DbPassword: password,
				Target:     db.Name,
			})

			if err != nil {
				Fail(err.Error())
				return
			}
		})

		It("access an table inside the database with the correct schema", func() {
			err := db.Conn.Exec(fmt.Sprintf(`SELECT FROM %s.examples`, randomSchema)).Error

			Expect(err).ShouldNot(HaveOccurred())
		})
	})
})
