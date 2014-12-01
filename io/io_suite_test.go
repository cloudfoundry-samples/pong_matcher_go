package io_test

import (
	. "github.com/camelpunch/pong_matcher_go/io"

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
	databaseUrl := "mysql2://gopong:gopong@127.0.0.1:3306/pong_matcher_go_test"
	InitDb(MigratedDbMap(databaseUrl, "../db/migrations"))
}
