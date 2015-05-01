package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

// Handler for GET /customers
func HandleCustomersGet(w http.ResponseWriter, req *http.Request) {
	r := RenderCtx(req)
	db := DbCtx(req)

	var message interface{} = nil
	customers, err := GetCustomers(db)
	success := (err == nil)

	r.JSON(w, http.StatusOK, map[string]interface{}{
		"customers": customers,
		"success":   success,
		"message":   message,
	})
}

// Handler for POST /customers
func HandleCustomersPost(w http.ResponseWriter, req *http.Request) {
	r := RenderCtx(req)
	db := DbCtx(req)

	var customer Customer

	// decode JSON body into Customer
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&customer)
	if err != nil {
		r.JSON(w, http.StatusBadRequest, map[string]interface{}{
			"message": "Malformed JSON body",
			"success": false,
		})
		return
	}

	// ensure that customer obj is valid
	if !customer.IsValid() {
		r.JSON(w, http.StatusBadRequest, map[string]interface{}{
			"message": "Invalid customer object",
			"success": false,
		})
		return
	}

	err = CreateCustomer(db, &customer)
	if err != nil {
		r.JSON(w, http.StatusBadRequest, map[string]interface{}{
			"message": "Failed to save customer object",
			"success": false,
		})
		return
	}

	r.JSON(w, http.StatusCreated, map[string]interface{}{
		"customer": customer,
		"success":  true,
	})
}

// Handler for GET /customer/{id}
func HandleCustomerGet(w http.ResponseWriter, req *http.Request) {
	r := RenderCtx(req)
	db := DbCtx(req)
	params := mux.Vars(req)
	id := params["id"]

	var message interface{} = nil
	customer, err := GetCustomer(db, id)
	if err != nil {
		message = err.Error()
		customer = nil
	}
	success := (err == nil)

	r.JSON(w, http.StatusOK, map[string]interface{}{
		"customer": customer,
		"success":  success,
		"message":  message,
	})
}
