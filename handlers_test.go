package main_test

import (
	. "pong_matcher_go"
	. "pong_matcher_go/domain"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
)

func stubbedMatchRequestRetrieval(success bool) MatchRequestRetriever {
	return func(uuid string) (bool, MatchRequest, error) {
		mr := MatchRequest{}
		return success, mr, nil
	}
}

func stubbedMatchRetrieval(success bool) MatchRetriever {
	return func(uuid string) (bool, Match) {
		m := Match{}
		return success, m
	}
}

var _ = Describe("Request handlers", func() {
	wellFormedJson := func() io.Reader { return strings.NewReader(`{"player": "some-player"}`) }
	blankBody := func() io.Reader { return strings.NewReader("") }

	Describe("getting a match request", func() {
		Context("when a match request is found", func() {
			It("responds with 200", func() {
				handle := GetMatchRequestHandler(stubbedMatchRequestRetrieval(true))

				resp := httptest.NewRecorder()
				req, err := http.NewRequest(
					"GET", "/match_requests/foo", nil,
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
					"GET", "/match_requests/foo", nil,
				)

				handle(resp, req)

				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Code).To(Equal(404))
			})
		})

		Context("when selecting the match request produces an error", func() {
			errorFetch := func(uuid string) (bool, MatchRequest, error) {
				return false, MatchRequest{}, errors.New("Bad stuff")
			}

			handle := GetMatchRequestHandler(errorFetch)

			It("responds with 500", func() {
				resp := httptest.NewRecorder()
				req, err := http.NewRequest(
					"GET", "/match_requests/foo", nil,
				)

				handle(resp, req)

				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Code).To(Equal(500))
			})
		})
	})

	Describe("creating a match request", func() {
		Context("when persistence works", func() {
			nullPersist := func(mr MatchRequest) error { return nil }
			handle := CreateMatchRequestHandler(nullPersist)

			Context("with a well-formed body", func() {
				It("responds with 200", func() {
					resp := httptest.NewRecorder()
					req, err := http.NewRequest(
						"PUT", "/match_requests/foo", wellFormedJson(),
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
						"PUT", "/match_requests/foo", blankBody(),
					)

					handle(resp, req)

					Expect(err).NotTo(HaveOccurred())
					Expect(resp.Code).To(Equal(400))
				})
			})
		})

		Context("when persistence fails", func() {
			errorPersist := func(mr MatchRequest) error {
				return errors.New("Bad stuff")
			}
			handle := CreateMatchRequestHandler(errorPersist)

			It("responds with 500", func() {
				resp := httptest.NewRecorder()
				req, err := http.NewRequest(
					"PUT", "/match_requests/foo", wellFormedJson(),
				)

				handle(resp, req)

				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Code).To(Equal(500))
			})
		})
	})

	Describe("getting a match", func() {
		Context("when a match is found", func() {
			It("responds with 200", func() {
				handle := MatchHandler(stubbedMatchRetrieval(true))

				resp := httptest.NewRecorder()
				req, err := http.NewRequest(
					"GET", "/matches/foo", nil,
				)

				handle(resp, req)

				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Code).To(Equal(200))
			})
		})

	})

	Describe("posting results of a match", func() {
		nullPersist := func(r Result) error { return nil }

		Context("with a well-formed body", func() {
			handle := ResultsHandler(nullPersist)

			It("responds with 201", func() {
				resp := httptest.NewRecorder()
				req, err := http.NewRequest(
					"POST", "/results", wellFormedJson(),
				)

				handle(resp, req)

				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Code).To(Equal(201))
			})
		})

		Context("without a body", func() {
			handle := ResultsHandler(nullPersist)

			It("responds with 400", func() {
				resp := httptest.NewRecorder()
				req, err := http.NewRequest(
					"POST", "/results", blankBody(),
				)

				handle(resp, req)

				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Code).To(Equal(400))
			})
		})

		Context("when persistence produces an error", func() {
			errorPersist := func(r Result) error {
				return errors.New("Bad stuff")
			}
			handle := ResultsHandler(errorPersist)

			It("responds with 500", func() {
				resp := httptest.NewRecorder()
				req, err := http.NewRequest(
					"POST", "/results", wellFormedJson(),
				)

				handle(resp, req)

				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Code).To(Equal(500))
			})
		})
	})

	Describe("wiping the database", func() {
		Context("when deletion produces an error", func() {
			errorDelete := func() error {
				return errors.New("Bad stuff")
			}
			handle := AllHandler(errorDelete)

			It("responds with 500", func() {
				resp := httptest.NewRecorder()
				req, err := http.NewRequest("DELETE", "/all", nil)

				handle(resp, req)

				Expect(err).NotTo(HaveOccurred())
				Expect(resp.Code).To(Equal(500))
			})
		})
	})
})
