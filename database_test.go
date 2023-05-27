package pgonthefly_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	pgonthefly "github.com/problem-company-toolkit/pgonthefly"
)

var _ = Describe("pgonthefly", func() {
	Describe("Creating and deleting a database", func() {
		It("creates a database successfully", func() {
			db, err := pgonthefly.CreateDatabase(dbName, host, port, user, password, opts)
			Expect(err).NotTo(HaveOccurred())
			Expect(db).NotTo(BeNil())

			// Test that the database connection is working
			err = db.Conn.Transaction(func(tx *gorm.DB) error {
				// This is just a dummy transaction for testing purposes
				return nil
			})
			Expect(err).NotTo(HaveOccurred())

			// Test database deletion
			err = pgonthefly.DeleteDatabase(dbName, host, port, user, password, db.Name)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
