package server

func LoadRoutes() *Routes {
	return &Routes{
		Route{"CustomersList", "GET", "/customers", HandleCustomersGet},
		Route{"CustomerCreate", "POST", "/customers", HandleCustomersPost},
		Route{"CustomerGet", "GET", "/customers/{id}", HandleCustomerGet},
	}
}
