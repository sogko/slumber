package sessions_test

import (
	"errors"
	"fmt"
	"github.com/modocache/gory"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	"github.com/sogko/golang-rest-api-server-example/domain"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"testing"
	"time"
)

var TestDatabaseServerName = "localhost"
var TestDatabaseName = "test_db"

var privateSigningKey []byte
var publicSigningKey []byte

func TestSessions(t *testing.T) {
	defineFactories()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sessions Suite")
}

var _ = BeforeSuite(func() {
	// each test node uses its own database
	TestDatabaseName = fmt.Sprintf("test_db_node%v", config.GinkgoConfig.ParallelNode)

	var err error
	privateSigningKey, err = ioutil.ReadFile("../keys/demo.rsa")
	if err != nil {
		panic(errors.New(fmt.Sprintf("Error loading private signing key: %v", err.Error())))
	}
	publicSigningKey, err = ioutil.ReadFile("../keys/demo.rsa.pub")
	if err != nil {
		panic(errors.New(fmt.Sprintf("Error loading public signing key: %v", err.Error())))
	}
})

var _ = AfterSuite(func() {
})

func defineFactories() {
	gory.Define("user", domain.User{}, func(factory gory.Factory) {
		factory["ID"] = gory.Sequence(func(n int) interface{} {
			return bson.NewObjectId()
		})
		factory["Username"] = gory.Sequence(func(n int) interface{} {
			return fmt.Sprintf("johndoe%d", n)
		})
		factory["Email"] = gory.Sequence(func(n int) interface{} {
			return fmt.Sprintf("john-doe-%d@example.com", n)
		})
		factory["Status"] = domain.StatusActive
		factory["Roles"] = domain.Roles{domain.RoleUser}
		factory["LastModifiedDate"] = time.Now()
		factory["CreatedDate"] = time.Now()
	})
}
