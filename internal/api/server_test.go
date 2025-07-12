package api

import (
	"log"
	"testing"

	"github.com/ffss92/webhookd/internal/postgres"
)

var (
	testDB *postgres.TestInstance
)

func TestMain(m *testing.M) {
	testDB = postgres.MustTestInstance()
	defer func() {
		if err := testDB.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	m.Run()
}
