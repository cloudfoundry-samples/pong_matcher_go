package main_test

import (
	. "pong_matcher_go"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"strings"
)

func stubbedMatchRequestRetrieval(success bool) func(string) (bool, MatchRequest) {
	return func(uuid string) (bool, MatchRequest) {
		mr := MatchRequest{}
		return success, mr
	}
}

func stubbedMatchRetrieval(success bool) func(string) (bool, Match) {
	return func(uuid string) (bool, Match) {
		m := Match{}
		return success, m
	}
}

var _ = Describe("Request handlers", func() {
	Describe("getting a match request", func() {
		Context("when a match request is found", func() {
			It("responds with 200", func() {
				handle := GetMatchRequestHandler(stubbedMatchRequestRetrieval(true))

				resp := httptest.NewRecorder()
				req, err := http.NewRequest(
					"GET",
					"/match_requests/foo",
					nil,
				)

				handle(resp, req)

				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Code).To(Equal(200))
			})
		})

		Context("when a match request is not found", func() {
			It("responds with 404", func() {
				handle := GetMatchRequestHandler(stubbedMatchRequestRetrieval(false))

				resp := httptest.NewRecorder()
				req, err := http.NewRequest(
					"GET",
					"/match_requests/foo",
					nil,
				)

				handle(resp, req)

				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Code).To(Equal(404))
			})
		})
	})

	Describe("creating a match request", func() {
		nullPersist := func(mr MatchRequest) {}
		handle := CreateMatchRequestHandler(nullPersist)

		Context("with a valid body", func() {
			It("responds with 200", func() {
				resp := httptest.NewRecorder()
				req, err := http.NewRequest(
					"PUT",
					"/match_requests/foo",
					strings.NewReader(`{"player": "some-player"}`),
				)

				handle(resp, req)

				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Code).To(Equal(200))
			})
		})

		Context("without a body", func() {
			It("responds with 400", func() {
				resp := httptest.NewRecorder()
				req, err := http.NewRequest(
					"PUT",
					"/match_requests/foo",
					strings.NewReader(""),
				)

				handle(resp, req)

				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Code).To(Equal(400))
			})
		})
	})

	Describe("getting a match", func() {
		Context("when a match is found", func() {
			It("responds with 200", func() {
				handle := MatchHandler(stubbedMatchRetrieval(true))

				resp := httptest.NewRecorder()
				req, err := http.NewRequest(
					"GET",
					"/matches/foo",
					nil,
				)

				handle(resp, req)

				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Code).To(Equal(200))
			})
		})

	})

	Describe("posting results of a match", func() {
		nullPersist := func(r Result) {}
		handle := ResultsHandler(nullPersist)

		Context("without a body", func() {
			It("responds with 400", func() {
				resp := httptest.NewRecorder()
				req, err := http.NewRequest(
					"POST",
					"/results",
					strings.NewReader(""),
				)

				handle(resp, req)

				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Code).To(Equal(400))
			})
		})
	})
})
