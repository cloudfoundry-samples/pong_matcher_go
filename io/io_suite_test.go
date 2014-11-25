package io_test

import (
	. "github.com/camelpunch/pong_matcher_go/io"

	"database/sql"
	"github.com/coopernurse/gorp"
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestPongMatcherGoIo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "PongMatcherGo IO Suite")
}

var _ = BeforeSuite(func() {
	createTestDb()
})

func createTestDb() {
	db, err := sql.Open("mysql", "gopong:gopong@tcp(127.0.0.1:3306)/pong_matcher_go_test")
	if err != nil {
		log.Fatalln("Failed to connect to test database:", err)
	}
	testDbMap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	InitDb(testDbMap)
	err = testDbMap.CreateTablesIfNotExists()
	if err != nil {
		log.Fatalln("Failed to create tables:", err)
	}
}
