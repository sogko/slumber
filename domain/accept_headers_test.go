package domain_test

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sogko/slumber/domain"
)

var _ = Describe("AcceptHeaders Tests", func() {
	var _ = Describe("NewAcceptHeadersFromString", func() {

		type testMap struct {
			TestValue       string
			ExpectedLen     int
			ExpectedResults domain.AcceptHeaders
		}
		type testMaps []testMap

		var tests = testMaps{
			testMap{
				TestValue: "",
				ExpectedResults: domain.AcceptHeaders{
					domain.AcceptHeader{domain.MediaType{}, 1},
				},
			},
			testMap{
				TestValue: ";",
				ExpectedResults: domain.AcceptHeaders{
					domain.AcceptHeader{domain.MediaType{}, 1},
				},
			},
			testMap{
				TestValue: ";q=",
				ExpectedResults: domain.AcceptHeaders{
					domain.AcceptHeader{domain.MediaType{
						Parameters: domain.MediaTypeParams{"q": ""},
					}, 1},
				},
			},
			testMap{
				TestValue: "application/json",
				ExpectedResults: domain.AcceptHeaders{
					domain.AcceptHeader{domain.MediaType{
						String:  "application/json",
						Type:    "application",
						Tree:    "",
						SubType: "json",
						Suffix:  "",
					}, 1},
				},
			},
			testMap{
				TestValue: "application/json;q=",
				ExpectedResults: domain.AcceptHeaders{
					domain.AcceptHeader{domain.MediaType{
						String:     "application/json",
						Type:       "application",
						Tree:       "",
						SubType:    "json",
						Suffix:     "",
						Parameters: domain.MediaTypeParams{"q": ""},
					}, 1},
				},
			},
			testMap{
				TestValue: "application/json;q",
				ExpectedResults: domain.AcceptHeaders{
					domain.AcceptHeader{domain.MediaType{
						String:     "application/json",
						Type:       "application",
						Tree:       "",
						SubType:    "json",
						Suffix:     "",
						Parameters: domain.MediaTypeParams{"q": ""},
					}, 1},
				},
			},
			testMap{
				TestValue: "application/json;  q=0.9 ",
				ExpectedResults: domain.AcceptHeaders{
					domain.AcceptHeader{domain.MediaType{
						String:     "application/json",
						Type:       "application",
						Tree:       "",
						SubType:    "json",
						Suffix:     "",
						Parameters: domain.MediaTypeParams{"q": "0.9"},
					}, 0.9},
				},
			},
			testMap{
				TestValue: "application/vnd.sgk.v1+json ",
				ExpectedResults: domain.AcceptHeaders{
					domain.AcceptHeader{domain.MediaType{
						String:  "application/vnd.sgk.v1+json",
						Type:    "application",
						Tree:    "vnd",
						SubType: "sgk.v1",
						Suffix:  "json",
					}, 1},
				},
			},
			testMap{
				TestValue: "application/vnd.sgk.v1+json;q=0.8",
				ExpectedResults: domain.AcceptHeaders{
					domain.AcceptHeader{domain.MediaType{
						String:  "application/vnd.sgk.v1+json",
						Type:    "application",
						Tree:    "vnd",
						SubType: "sgk.v1",
						Suffix:  "json",
						Parameters: domain.MediaTypeParams{
							"q": "0.8",
						},
					}, 0.8},
				},
			},
			testMap{
				TestValue: "application/vnd.sgk+json; q=0.8 ;version=1.0",
				ExpectedResults: domain.AcceptHeaders{
					domain.AcceptHeader{domain.MediaType{
						String:  "application/vnd.sgk+json",
						Type:    "application",
						Tree:    "vnd",
						SubType: "sgk",
						Suffix:  "json",
						Parameters: domain.MediaTypeParams{
							"q":       "0.8",
							"version": "1.0",
						},
					}, 0.8},
				},
			},
			testMap{
				TestValue: "application/vnd.sgk.rest-api-server.v1+json; q=0.8 ;version=1.0,*/*",
				ExpectedResults: domain.AcceptHeaders{
					domain.AcceptHeader{domain.MediaType{
						String:  "application/vnd.sgk.rest-api-server.v1+json",
						Type:    "application",
						Tree:    "vnd",
						SubType: "sgk.rest-api-server.v1",
						Suffix:  "json",
						Parameters: domain.MediaTypeParams{
							"q":       "0.8",
							"version": "1.0",
						},
					}, 0.8},
					domain.AcceptHeader{domain.MediaType{
						String:  "*/*",
						Type:    "*",
						Tree:    "",
						SubType: "*",
						Suffix:  "",
					}, 1},
				},
			},
			testMap{
				TestValue: "application/vnd.sgk+json; q=0.8 ;version=1.0,application/json , */*;q=noninteger",
				ExpectedResults: domain.AcceptHeaders{
					domain.AcceptHeader{domain.MediaType{
						String:  "application/vnd.sgk+json",
						Type:    "application",
						Tree:    "vnd",
						SubType: "sgk",
						Suffix:  "json",
						Parameters: domain.MediaTypeParams{
							"q":       "0.8",
							"version": "1.0",
						},
					}, 0.8},
					domain.AcceptHeader{domain.MediaType{
						String:  "application/json",
						Type:    "application",
						Tree:    "",
						SubType: "json",
						Suffix:  "",
					}, 1},
					domain.AcceptHeader{domain.MediaType{
						String:  "*/*",
						Type:    "*",
						Tree:    "",
						SubType: "*",
						Suffix:  "",
						Parameters: domain.MediaTypeParams{
							"q": "noninteger",
						},
					}, 1},
				},
			},
		}

		for _, test := range tests {
			testValue := test.TestValue
			expectedResults := test.ExpectedResults
			Context(fmt.Sprintf("when `Accept=%v`", testValue), func() {
				result := domain.NewAcceptHeadersFromString(testValue)
				It("parses OK", func() {
					Expect(len(result)).To(Equal(len(expectedResults)))
					for i, _ := range expectedResults {
						Expect(result[i].MediaType).To(Equal(expectedResults[i].MediaType))
						Expect(result[i].QualityFactor).To(Equal(expectedResults[i].QualityFactor))
					}
				})
			})
		}

	})
})
