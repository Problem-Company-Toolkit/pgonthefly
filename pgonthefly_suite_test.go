package pgonthefly_test

import (
	"os"
	"testing"

	pgonthefly "github.com/problem-company-toolkit/pgonthefly"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPgonthefly(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pgonthefly Suite")
}

var (
	dbName   string
	host     string
	port     string
	user     string
	password string
	opts     pgonthefly.DatabaseOptions
)

var _ = BeforeEach(func() {
	dbName = os.Getenv("DATABASE_NAME")
	host = os.Getenv("DATABASE_HOST")
	port = os.Getenv("DATABASE_PORT")
	user = os.Getenv("DATABASE_USER")
	password = os.Getenv("DATABASE_PASSWORD")

	automigrateAll := func(db *pgonthefly.DB) error {
		// This is just a dummy function. Replace this with your own logic.
		return nil
	}

	opts = pgonthefly.DatabaseOptions{
		AutomigrateFunc: automigrateAll,
	}
})
