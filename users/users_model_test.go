package users_test

import (
	"fmt"
	"github.com/modocache/gory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/sogko/golang-rest-api-server-example/users"
)

const TestValidPassword = "PASSWORD124"
const TestInvalidPassword = "PASSWORDINVALID"

var _ = Describe("Users Model", func() {

	var user *User
	BeforeEach(func() {
		user = gory.Build("user").(*User)
	})

	Describe("user.IsValid()", func() {

		Context("when user is not valid", func() {

			BeforeEach(func() {
				user.Email = "invalid email"
			})
			It("returns false", func() {
				Expect(user.IsValid()).To(BeFalse())
			})
		})

		Context("when user is valid", func() {

			BeforeEach(func() {
				user.Email = "email@example"
			})
			It("returns false", func() {
				Expect(user.IsValid()).To(BeTrue())
			})
		})
	})
	Describe("user.IsCodeVerified()", func() {

		Context("when code is valid", func() {

			var code string
			BeforeEach(func() {
				user.GenerateConfirmationCode()
				code = user.ConfirmationCode
			})
			It("returns false", func() {
				Expect(user.IsCodeVerified(code)).To(BeTrue())
			})
		})

		Context("when code is not valid", func() {

			var code string
			BeforeEach(func() {
				user.GenerateConfirmationCode()
				code = fmt.Sprintf("%vWRONGCODE", user.ConfirmationCode)
			})
			It("returns false", func() {
				Expect(user.IsCodeVerified(code)).To(BeFalse())
			})
		})
	})
	Describe("user.IsCredentialsVerified()", func() {

		Context("when password is valid", func() {

			BeforeEach(func() {
				user.SetPassword(TestValidPassword)
			})
			It("returns false", func() {
				Expect(user.IsCredentialsVerified(TestValidPassword)).To(BeTrue())
			})
		})

		Context("when password is not valid", func() {

			BeforeEach(func() {
				user.SetPassword(TestValidPassword)
			})
			It("returns false", func() {
				Expect(user.IsCredentialsVerified(TestInvalidPassword)).To(BeFalse())
			})
		})
	})
})
