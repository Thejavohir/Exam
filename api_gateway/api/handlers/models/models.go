package models

type Customer struct {
	FirstName   string
	LastName    string
	Bio         string
	Email       string
	PhoneNumber string
	Password    string
}

type CustomerRegister struct {
	FirstName   string
	LastName    string
	Bio         string
	Email       string
	Password    string
	PhoneNumber string
	Code        string
}

type Post struct {
	Name        string
	Description string
	CustomerId  string
}

type Review struct {
	Name        string
	Review      int
	Description string
	PostId      int
}

type Error struct {
	Code        int
	Error       error
	Description string
}

type VerifyResponse struct {
	Id           int64            `json:"id"`
	FirstName    string            `json:"first_name"`
	LastName     string            `json:"last_name"`
	Email        string            `json:"email"`
	Bio          string            `json:"bio"`
	PhoneNumber  string            `json:"phone_number"`
	JWT          string            `json:"jwt"`
	RefreshToken string            `json:"refresh"`
}


type CustomerData struct {
	FirstName string //`json:"first_name"`
	LastName  string //`json:"last_name"`
	Bio       string //`json:"username"`
	Password  string //`json:"password"`
	Email     string //`json:"email"`
	Code      string
}