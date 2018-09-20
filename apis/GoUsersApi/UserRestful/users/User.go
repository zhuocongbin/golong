package users

type User struct {
	Id       string  `bson:"_id,omitempty" json:"_id"`
	Username string `bson:"username" json:"username"`
	Password string `bson:"password" json:"password"`
}

type UserService interface {
	CreateUser(u *User) error
	GetByUsername(username string) (*User, error)
}

type Hash interface {
	Generate(s string) (string, error)
	Compare(hash string, s string) error
}
