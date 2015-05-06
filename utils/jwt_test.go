package utils_test

import (
	"github.com/dgrijalva/jwt-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sogko/golang-rest-api-server-example/utils"
	"io/ioutil"
	"time"
)

var hmacTestKey, _ = ioutil.ReadFile("../keys/hmacTestKey")

var _ = Describe("JWT", func() {

	Describe("CreateNewToken()", func() {
		Context("when all claims are set", func() {
			tokenString, err := utils.CreateNewToken(utils.TokenClaims{
				ID:     "id1234",
				Status: "active",
				Roles:  []string{"user"},
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns a valid tokenString", func() {
				Expect(tokenString).NotTo(BeNil())
			})

		})
		Context("when only some claims are set", func() {
			tokenString, err := utils.CreateNewToken(utils.TokenClaims{
				ID: "id1234",
			})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns a valid tokenString", func() {
				Expect(tokenString).NotTo(BeNil())
			})

		})
		Context("when claims in nil", func() {
			tokenString, err := utils.CreateNewToken(utils.TokenClaims{})
			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns a valid tokenString", func() {
				Expect(tokenString).NotTo(BeNil())
			})

		})
	})
	Describe("ParseAndVerifyTokenString()", func() {
		var originalClaims = utils.TokenClaims{
			ID:     "id1234",
			Status: "active",
			Roles:  []string{"user"},
		}

		Context("when token string is a valid token", func() {

			var err error
			var tokenString string

			var token *jwt.Token
			var claims *utils.TokenClaims

			BeforeEach(func() {
				tokenString, _ = utils.CreateNewToken(originalClaims)
				token, claims, err = utils.ParseAndVerifyTokenString(tokenString)
			})

			It("returns no error", func() {
				Expect(err).To(BeNil())
			})
			It("returns a valid token", func() {
				Expect(token.Valid).To(BeTrue())
			})
			It("returns equal claims", func() {
				Expect(claims.ID).To(Equal(originalClaims.ID))
				Expect(claims.Status).To(Equal(originalClaims.Status))
				Expect(claims.Roles).To(Equal(originalClaims.Roles))
			})
		})

		Context("when token string is not a valid token", func() {
			var err error

			var token *jwt.Token
			var claims *utils.TokenClaims

			BeforeEach(func() {
				token, claims, err = utils.ParseAndVerifyTokenString("RANDOMTOKENSTRING")
			})
			It("returns error", func() {
				Expect(err).NotTo(BeNil())
			})
			It("return a nil token", func() {
				Expect(token).To(BeNil())
			})
			It("returns a nil claims", func() {
				Expect(claims).To(BeNil())
			})

		})

		Context("when token string was signed by another algorithm", func() {
			var err error
			var tokenString string

			var token *jwt.Token
			var claims *utils.TokenClaims

			BeforeEach(func() {
				// a different algorithm
				token := jwt.New(jwt.SigningMethodHS256)
				token.Claims = map[string]interface{}{
					"id":     originalClaims.ID,
					"status": originalClaims.Status,
					"roles":  originalClaims.Roles,
					"exp":    time.Now().Add(time.Hour * 72).Unix(), // 3 days
					"iat":    time.Now(),
					"jti":    "nonce",
				}
				tokenString, err = token.SignedString(hmacTestKey)
				token, claims, err = utils.ParseAndVerifyTokenString(tokenString)
			})
			It("returns error", func() {
				Expect(err).NotTo(BeNil())
			})
			It("return a nil token", func() {
				Expect(token).To(BeNil())
			})
			It("returns a nil claims", func() {
				Expect(claims).To(BeNil())
			})

		})
	})
})
