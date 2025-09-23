package user

type UserRepository interface {
	Create(user User, credentials Credentials) (User, error)
	findByID(id int) (User, error)
	findByEmail(email string) (User, Credentials, error)
}