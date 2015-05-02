package server_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/modocache/gory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sogko/golang-rest-api-server-example/server"
	. "github.com/sogko/golang-rest-api-server-example/server/models"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Routes * /api/v1/customers", func() {
	var server *Server
	var session *DatabaseSession
	var renderer *Renderer
	var request *http.Request
	var recorder *httptest.ResponseRecorder
	var bodyJSON map[string]interface{}

	BeforeEach(func() {

		// set up server with test components
		session = NewSession(DatabaseOptions{
			ServerName:   TestDatabaseServerName,
			DatabaseName: TestDatabaseName,
		})

		renderer = NewRenderer(RendererOptions{
			IndentJSON: true,
		})
		components := Components{
			DatabaseSession: session,
			Renderer:        renderer,
		}
		server = NewServer(&components)

		// record HTTP responses
		recorder = httptest.NewRecorder()
	})

	AfterEach(func() {
		// drop database after each test
		session.DB(session.DatabaseName).DropDatabase()
	})

	Describe("GET /api/v1/customers", func() {

		Context("when no customers exist", func() {

			BeforeEach(func() {
				// serve request
				request, _ = http.NewRequest("GET", "/api/v1/customers", nil)
				server.ServeHTTP(recorder, request)
				bodyJSON = MapFromJSON(recorder.Body.Bytes())
			})

			It("returns status code of StatusOK (200)", func() {
				Expect(recorder.Code).To(Equal(http.StatusOK))
			})
			It("returns zero customers", func() {
				customers := bodyJSON["customers"].([]interface{})

				Expect(bodyJSON["success"]).To(Equal(true))
				Expect(len(customers)).To(Equal(0))
			})
		})

		Context("when two customers exist", func() {

			BeforeEach(func() {
				// insert two valid customers into database
				collection := session.DB(TestDatabaseName).C(CustomersCollection)
				collection.Insert(gory.Build("customer"))
				collection.Insert(gory.Build("customer"))

				// serve request
				request, _ = http.NewRequest("GET", "/api/v1/customers", nil)
				server.ServeHTTP(recorder, request)
				bodyJSON = MapFromJSON(recorder.Body.Bytes())
			})

			It("returns status code of StatusOK (200)", func() {
				Expect(recorder.Code).To(Equal(http.StatusOK))
			})
			It("returns two customers", func() {
				customers := bodyJSON["customers"].([]interface{})

				Expect(bodyJSON["success"]).To(Equal(true))
				Expect(len(customers)).To(Equal(2))
			})

		})
	})

	Describe("POST /api/v1/customers", func() {
		Context("when adding one valid customer", func() {

			var newCustomer *Customer
			BeforeEach(func() {

				newCustomer = gory.Build("customer").(*Customer)
				body, _ := json.Marshal(newCustomer)

				request, _ = http.NewRequest("POST", "/api/v1/customers", bytes.NewReader(body))
				server.ServeHTTP(recorder, request)
				bodyJSON = MapFromJSON(recorder.Body.Bytes())
			})

			It("returns status code of StatusCreated (201)", func() {
				Expect(recorder.Code).To(Equal(http.StatusCreated))
			})
			It("returns newly-created customer", func() {
				customer := bodyJSON["customer"].(map[string]interface{})

				Expect(bodyJSON["success"]).To(Equal(true))
				Expect(customer["firstName"]).To(Equal(newCustomer.FirstName))
			})
		})

		Context("when POSTing a malformed JSON body", func() {

			BeforeEach(func() {

				request, _ = http.NewRequest("POST", "/api/v1/customers", bytes.NewReader([]byte("Bad JSON")))
				server.ServeHTTP(recorder, request)
				bodyJSON = MapFromJSON(recorder.Body.Bytes())
			})

			It("returns status code of StatusBadRequest (400)", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
			It("returns nil customer", func() {
				Expect(bodyJSON["success"]).To(Equal(false))
				Expect(bodyJSON["customer"]).To(BeNil())
			})
		})

		Context("when adding one invalid customer (missing `firstName`)", func() {

			var newCustomer *Customer
			BeforeEach(func() {

				newCustomer = gory.Build("customerMissingFirstName").(*Customer)
				body, _ := json.Marshal(newCustomer)

				request, _ = http.NewRequest("POST", "/api/v1/customers", bytes.NewReader(body))
				server.ServeHTTP(recorder, request)
				bodyJSON = MapFromJSON(recorder.Body.Bytes())
			})

			It("returns status code of StatusBadRequest (400)", func() {
				Expect(recorder.Code).To(Equal(http.StatusBadRequest))
			})
			It("returns nil customer", func() {
				Expect(bodyJSON["success"]).To(Equal(false))
				Expect(bodyJSON["customer"]).To(BeNil())
			})
		})

	})

	Describe("GET /api/v1/customers/{id}", func() {
		Context("when customer exists", func() {

			var customer *Customer
			BeforeEach(func() {

				// insert a customer into database
				customer = gory.Build("customer").(*Customer)
				collection := session.DB(TestDatabaseName).C(CustomersCollection)
				customer.ID = bson.NewObjectId()
				collection.Insert(customer)

				// serve request
				request, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/customers/%v", customer.ID.Hex()), nil)
				server.ServeHTTP(recorder, request)
				bodyJSON = MapFromJSON(recorder.Body.Bytes())
			})

			It("returns status code of StatusOK (200)", func() {
				Expect(recorder.Code).To(Equal(http.StatusOK))
			})
			It("returns specified customer", func() {
				Expect(bodyJSON["success"]).To(Equal(true))

				retCustomer := bodyJSON["customer"].(map[string]interface{})

				Expect(retCustomer["_id"]).To(Equal(customer.ID.Hex()))
				Expect(retCustomer["firstName"]).To(Equal(customer.FirstName))
			})
		})

		Context("when customer does not exists", func() {

			BeforeEach(func() {

				// serve request
				request, _ = http.NewRequest("GET", fmt.Sprintf("/api/v1/customers/%v", bson.NewObjectId().Hex()), nil)
				server.ServeHTTP(recorder, request)
				bodyJSON = MapFromJSON(recorder.Body.Bytes())
			})

			It("returns status code of StatusOK (200)", func() {
				Expect(recorder.Code).To(Equal(http.StatusOK))
			})
			It("returns unsuccessful response", func() {
				Expect(bodyJSON["success"]).To(Equal(false))
				Expect(bodyJSON["customer"]).To(BeNil())
			})
		})

		Context("when customer `id` is invalid", func() {

			BeforeEach(func() {

				// serve request
				request, _ = http.NewRequest("GET", "/api/v1/customers/INVALIDID", nil)
				server.ServeHTTP(recorder, request)
				bodyJSON = MapFromJSON(recorder.Body.Bytes())
			})

			It("returns status code of StatusOK (200)", func() {
				Expect(recorder.Code).To(Equal(http.StatusOK))
			})
			It("returns unsuccessful response", func() {
				Expect(bodyJSON["success"]).To(Equal(false))
				Expect(bodyJSON["customer"]).To(BeNil())
			})
		})

	})
})
