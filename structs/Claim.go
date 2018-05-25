package structs

type Claim struct {
	Iss string
	Email string
	Iat int64
	Exp int64
}