package OnurTPIUser

import "golang.org/x/crypto/bcrypt"

const (
	//USER ROLE
	USER = 1
	// EDITOR ROLE
	EDITOR = 2
	// ADMIN ROLE
	ADMIN = 3
)

// CredentialsClient contains login data from client
type CredentialsClient struct {
	UserName string `json:"username"`
	Password string `json:"password,omitempty"`
}

// User contains user data
type User struct {
	ID        string `json:"id,omitempty"`
	UserName  string `json:"username"`
	RoleValue int    `json:"roleValue"`
}

// Account contains the users and password data
type Account struct {
	User
	Password string `json:"password, omitempty"`
}

// NewAccount creates an account with credential data
func (c *CredentialsClient) CreateAccount() *Account {
	a := new(Account)

	a.UserName = c.UserName
	a.RoleValue = 1

	// Hash Password
	a.Password = bcryptHash(c.Password)

	return a
}

// ************************************************************
// Password encryption
// ************************************************************
func bcryptHash(password string) string {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		panic(err)
	}

	return string(hashedPassword)
}

//CompareHashAndPassword With bcryptHash I compare the password
func CompareHashAndPassword(hashedPassword, password string) bool {
	// nil means it is a match
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	return err == nil
}
