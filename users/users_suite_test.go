package users_test

import (
	"fmt"
	"github.com/modocache/gory"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
	. "github.com/sogko/golang-rest-api-server-example/users"
	"gopkg.in/mgo.v2/bson"
	"testing"
	"time"
)

var TestDatabaseServerName = "localhost"
var TestDatabaseName = "test_db"

func TestUsers(t *testing.T) {
	defineFactories()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Users Suite")
}

var _ = BeforeSuite(func() {
	// each test node uses its own database
	TestDatabaseName = fmt.Sprintf("test_db_node%v", config.GinkgoConfig.ParallelNode)
})

var _ = AfterSuite(func() {
})

func defineFactories() {
	gory.Define("user", User{}, func(factory gory.Factory) {
		factory["ID"] = gory.Sequence(func(n int) interface{} {
			return bson.NewObjectId()
		})
		factory["Username"] = gory.Sequence(func(n int) interface{} {
			return fmt.Sprintf("johndoe%d", n)
		})
		factory["Email"] = gory.Sequence(func(n int) interface{} {
			return fmt.Sprintf("john-doe-%d@example.com", n)
		})
		factory["Status"] = StatusActive
		factory["LastModifiedDate"] = time.Now()
		factory["CreatedDate"] = time.Now()
		factory["Roles"] = Roles{RoleUser}
	})
	gory.Define("userInvalidEmail", User{}, func(factory gory.Factory) {
		factory["ID"] = gory.Sequence(func(n int) interface{} {
			return bson.NewObjectId()
		})
		factory["Username"] = gory.Sequence(func(n int) interface{} {
			return fmt.Sprintf("johndoe%d", n)
		})
		factory["Email"] = gory.Sequence(func(n int) interface{} {
			return fmt.Sprintf("john-doe-%d", n)
		})
	})
	gory.Define("userUnconfirmed", User{}, func(factory gory.Factory) {
		factory["ID"] = gory.Sequence(func(n int) interface{} {
			return bson.NewObjectId()
		})
		factory["Username"] = gory.Sequence(func(n int) interface{} {
			return fmt.Sprintf("johndoe%d", n)
		})
		factory["Email"] = gory.Sequence(func(n int) interface{} {
			return fmt.Sprintf("john-doe-%d", n)
		})
		factory["Status"] = StatusPending
		factory["LastModifiedDate"] = time.Now()
		factory["CreatedDate"] = time.Now()
	})
}
