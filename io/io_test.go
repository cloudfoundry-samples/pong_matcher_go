package io_test

import (
	. "github.com/camelpunch/pong_matcher_go/io"
	. "github.com/camelpunch/pong_matcher_go/domain"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("IO:", func() {
	BeforeEach(func() {
		DeleteAll()
	})

	Describe("getting a match", func() {
		Context("with an empty database", func() {
			It("informs the caller that no match was found", func() {
				success, _ := GetMatch("nonexistentUUID")
				Expect(success).To(Equal(false))
			})
		})
	})

	Describe("entering a result", func() {
		It("returns a nil error on success", func() {
			err := PersistResult(Result{})
			Expect(err).To(BeNil())
		})

		It("removes match_id from match request (pending rename of match_id)", func() {
			err := PersistMatchRequest(MatchRequest{Uuid: "foo", RequesterId: "andrew"})
			Expect(err).To(BeNil())

			err = PersistMatchRequest(MatchRequest{Uuid: "bar", RequesterId: "india"})
			Expect(err).To(BeNil())

			_, matchRequest, err := GetMatchRequest("foo")

			err = PersistResult(Result{MatchId: matchRequest.MatchId.String})
			Expect(err).To(BeNil())

			_, matchRequest, err = GetMatchRequest("foo")
			Expect(err).To(BeNil())
			Expect(matchRequest.MatchId.Valid).To(BeFalse())
		})
	})
})
