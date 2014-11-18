package main_test

import (
	. "pong_matcher_go"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"strings"
)

func stubbedRetrieval(success bool) func(string) (bool, MatchRequest) {
	return func(uuid string) (bool, MatchRequest) {
		mr := MatchRequest{}
		return success, mr
	}
}

var _ = Describe("Request handlers", func() {
	Describe("GetMatchRequestHandler", func() {
		Context("when a match request is found", func() {
			It("responds with 200", func() {
				handle := GetMatchRequestHandler(stubbedRetrieval(true))

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
				handle := GetMatchRequestHandler(stubbedRetrieval(false))

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

	Describe("CreateMatchRequestHandler", func() {
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
})
