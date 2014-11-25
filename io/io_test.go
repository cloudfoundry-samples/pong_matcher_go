package io_test

import (
	. "github.com/camelpunch/pong_matcher_go/io"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IO:", func() {
	Describe("getting a match", func() {
		Context("with an empty database", func() {
			BeforeEach(func() {
				DeleteAll()
			})

			It("informs the caller that no match was found", func() {
				success, _ := GetMatch("nonexistentUUID")
				Expect(success).To(Equal(false))
			})
		})
	})
})
