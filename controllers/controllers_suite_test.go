package controllers_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/ginkgo/config"
	"testing"
	"fmt"
)

var TestDatabaseServerName = "localhost"
var TestDatabaseName = "test_db"

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Controllers Suite")
}

var _ = BeforeSuite(func() {
	// each test node uses its own database
	TestDatabaseName = fmt.Sprintf("test_db_node%v", GinkgoConfig.ParallelNode)
})

var _ = AfterSuite(func() {
})
