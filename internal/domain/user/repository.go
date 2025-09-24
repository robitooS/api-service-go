package user

type UserRepository interface {
	Create(user *User, credentials *Credentials) (*User, error)
	FindByID(id int) (*User, error)
	FindByEmail(email string) (*User, *Credentials, error)
}
