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
