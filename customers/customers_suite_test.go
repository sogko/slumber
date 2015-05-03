package customers_test

import (
	"encoding/json"
	"fmt"
	"github.com/modocache/gory"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	. "github.com/sogko/golang-rest-api-server-example/customers"
	"testing"
)

var TestDatabaseServerName = "localhost"
var TestDatabaseName = "test_db"

func TestServer(t *testing.T) {
	defineFactories()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Customers Suite")
}

var _ = BeforeSuite(func() {

	// each test node uses its own database
	TestDatabaseName = fmt.Sprintf("test_db_node%v", config.GinkgoConfig.ParallelNode)
})

var _ = AfterSuite(func() {
})

func defineFactories() {
	gory.Define("customer", Customer{}, func(factory gory.Factory) {
		factory["FirstName"] = "John"
		factory["LastName"] = "Doe"
		factory["Email"] = gory.Sequence(func(n int) interface{} {
			return fmt.Sprintf("john-doe-%d@example.com", n)
		})
	})

	gory.Define("customerMissingFirstName", Customer{}, func(factory gory.Factory) {
		factory["LastName"] = "Doe"
		factory["Email"] = gory.Sequence(func(n int) interface{} {
			return fmt.Sprintf("john-doe-%d@example.com", n)
		})
	})
}

func MapFromJSON(data []byte) map[string]interface{} {
	var result interface{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		panic(fmt.Sprintf("MapFromJSON(): Not a valid JSON body\n%v", string(data)))
	}
	return result.(map[string]interface{})
}
