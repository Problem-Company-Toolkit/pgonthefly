# pgonthefly

`pgonthefly` is a Golang package providing an interface for creating and deleting PostgreSQL databases on-the-fly. This tool is intended for use in integration testing, allowing developers to verify that database migrations and interfaces are functioning correctly in a controlled environment.

## Use Cases

This package is useful for:

- Testing if migrations work.
- Testing if the database interfaces work correctly with the database (integration testing).

## Limitations

- This package is designed for integration testing and should **not** be used for unit testing.
- Databases created by this package require manual deletion.

## Usage

Here is an example of how to use `pgonthefly`:

```go
package mypackage_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	pgonthefly "github.com/problem-company-toolkit/pgonthefly"
)

var _ = Describe("My DAO tests", func() {
	var db *pgonthefly.DB

	AutomigrateAll := func(db *pgonthefly.DB) error {
		// Implement your logic here
		return nil
	}

	BeforeEach(func() {
		var err error
		db, err = pgonthefly.CreateDatabase("my_db", "localhost", "5432", "my_user", "my_password", pgonthefly.DatabaseOptions{
			AutomigrateFunc: AutomigrateAll,
		})
		Expect(err).NotTo(HaveOccurred())

		// Initiate your DAO here using the db.Conn
	})

	AfterEach(func() {
		err := pgonthefly.DeleteDatabase("my_db", "localhost", "5432", "my_user", "my_password", db.Name)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should interact with the real database correctly", func() {
		// Your tests here
	})
})
```

The example illustrates a typical testing scenario where:

- In a top-level `BeforeEach`, we create a test database and initiate the DAO.
- In the tests, we interact with the real database.
- In a top-level `AfterEach`, we delete the test database to clean up.

**IMPORTANT:** Databases created by this package need to be deleted manually. Ensure you delete any test databases created during your testing to avoid resource clutter.