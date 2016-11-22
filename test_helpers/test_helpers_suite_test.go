package test_helpers_test

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"
)

var TestDatabaseServerName = "localhost"
var TestDatabaseName = "test_db"

var privateSigningKey *rsa.PrivateKey
var publicSigningKey *rsa.PublicKey

func TestTestHelpers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Helpers Suite")
}

var _ = BeforeSuite(func() {
	// each test node uses its own database
	TestDatabaseName = fmt.Sprintf("test_db_node%v", config.GinkgoConfig.ParallelNode)

	var err error
	privateKey, err := ioutil.ReadFile("../keys/demo.rsa")
	if err != nil {
		panic(errors.New(fmt.Sprintf("Error loading private signing key: %v", err.Error())))
	}
	publicKey, err := ioutil.ReadFile("../keys/demo.rsa.pub")
	if err != nil {
		panic(errors.New(fmt.Sprintf("Error loading public signing key: %v", err.Error())))
	}

	// Casting keys loaded to proper type
	privateSigningKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		panic(errors.New(fmt.Sprintf("Error parsing private key: %v", err.Error())))
	}

	publicSigningKey, err = jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		panic(errors.New(fmt.Sprintf("Error parsing public key: %v", err.Error())))
	}
})
