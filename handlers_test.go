package main_test

import (
	. "pong_matcher_go"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"strings"
)

var _ = Describe("Request handlers", func() {
	Describe("CreateMatchRequestHandler", func() {
		var nullPersist = func(mr MatchRequest) {}
		var handle = CreateMatchRequestHandler(nullPersist)

		Describe("with a valid body", func() {
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

		Describe("without a body", func() {
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
