package server

// GetRoutes Wire API routes to controllers (http.HandlerFunc)
func GetRoutes() *Routes {
	return &Routes{
		Route{"CustomersList", "GET", "/api/v1/customers", HandleCustomersGet},
		Route{"CustomerCreate", "POST", "/api/v1/customers", HandleCustomersPost},
		Route{"CustomerGet", "GET", "/api/v1/customers/{id}", HandleCustomerGet},
	}
}
