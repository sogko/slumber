package domain_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sogko/slumber/domain"
)

var _ = Describe("Users Tests", func() {
	Describe("User Struct", func() {
		Describe("IsValid()", func() {
			Context("when user is valid", func() {
				user := domain.User{
					Username: "username",
					Email:    "email@gmail.com",
				}
				It("should return true", func() {
					Expect(user.IsValid()).To(BeTrue())
				})
			})
			Context("when username is empty", func() {
				user := domain.User{
					Username: "",
					Email:    "email@gmail.com",
				}
				It("should return false", func() {
					Expect(user.IsValid()).To(BeFalse())
				})
			})
			Context("when email is empty", func() {
				user := domain.User{
					Username: "username",
					Email:    "",
				}
				It("should return false", func() {
					Expect(user.IsValid()).To(BeFalse())
				})
			})
			Context("when email format is invalid", func() {
				user := domain.User{
					Username: "username",
					Email:    "email-asd.com",
				}
				It("should return false", func() {
					Expect(user.IsValid()).To(BeFalse())
				})
			})
		})

		Describe("IsCodeVerified()", func() {
			user := domain.User{
				Username: "username",
				Email:    "email@gmail.com",
			}
			user.GenerateConfirmationCode()
			Context("when code is correct", func() {
				It("should return true", func() {
					Expect(user.IsCodeVerified(user.ConfirmationCode)).To(BeTrue())
				})
			})
			Context("when code is empty", func() {
				It("should return false", func() {
					Expect(user.IsCodeVerified("")).To(BeFalse())
				})
			})
			Context("when code is wrong", func() {
				It("should return false", func() {
					Expect(user.IsCodeVerified("nottherightcode")).To(BeFalse())
				})
			})
		})

		Describe("IsCredentialsVerified()", func() {
			user := domain.User{
				Username: "username",
				Email:    "email@gmail.com",
			}
			password := "password"
			user.SetPassword(password)
			Context("when password is correct", func() {
				It("should return true", func() {
					Expect(user.IsCredentialsVerified(password)).To(BeTrue())
				})
			})
			Context("when password is empty", func() {
				It("should return false", func() {
					Expect(user.IsCredentialsVerified("")).To(BeFalse())
				})
			})
			Context("when password is wrong", func() {
				It("should return false", func() {
					Expect(user.IsCredentialsVerified("nottherightpassword")).To(BeFalse())
				})
			})
		})

		Describe("HasRole()", func() {
			Context("when user has one role", func() {
				user := domain.User{
					Username: "username",
					Email:    "email@gmail.com",
					Roles: domain.Roles{
						domain.RoleUser,
					},
				}
				Context("when user has role", func() {
					It("should return true", func() {
						Expect(user.HasRole(domain.RoleUser)).To(BeTrue())
					})
				})
				Context("when user does not have role", func() {
					It("should return false", func() {
						Expect(user.HasRole(domain.RoleAdmin)).To(BeFalse())
					})
				})
			})
			Context("when user has two roles", func() {
				user := domain.User{
					Username: "username",
					Email:    "email@gmail.com",
					Roles: domain.Roles{
						domain.RoleUser,
						domain.RoleAdmin,
					},
				}
				Context("when user has role", func() {
					It("should return true", func() {
						Expect(user.HasRole(domain.RoleUser)).To(BeTrue())
					})
				})
				Context("when user does not have role", func() {
					It("should return false", func() {
						Expect(user.HasRole("")).To(BeFalse())
					})
				})
			})
			Context("when user has no roles", func() {
				user := domain.User{
					Username: "username",
					Email:    "email@gmail.com",
				}
				Context("when user does not have role", func() {
					It("should return false", func() {
						Expect(user.HasRole(domain.RoleUser)).To(BeFalse())
					})
				})
			})
		})

		Describe("SetPassword()", func() {
			user := domain.User{
				Username: "username",
				Email:    "email@gmail.com",
				Roles: domain.Roles{
					domain.RoleUser,
				},
			}
			password := "password"
			user.SetPassword(password)
			It("should have HashedPassword", func() {
				Expect(user.HashedPassword).ToNot(BeNil())
			})
			It("should have HashedPassword that is not equal to plaintext", func() {
				Expect(user.HashedPassword).ToNot(Equal(password))
			})
		})

		Describe("GenerateConfirmationCode()", func() {
			user := domain.User{
				Username: "username",
				Email:    "email@gmail.com",
				Roles: domain.Roles{
					domain.RoleUser,
				},
			}
			user.GenerateConfirmationCode()
			It("should have ConfirmationCode", func() {
				Expect(user.ConfirmationCode).ToNot(BeNil())
			})
		})
	})
})
