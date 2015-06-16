package domain

type IRole interface{}
type IRoles []IRole

type IUser interface {
	GetID() string
	IsValid() bool
	IsCodeVerified(code string) bool
	IsCredentialsVerified(password string) bool
	SetPassword(password string) error
	GenerateConfirmationCode()
	HasRole(r IRole) bool
}

type IUsers interface{}
