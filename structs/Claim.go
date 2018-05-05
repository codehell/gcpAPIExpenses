package structs

type Claim struct {
	iss string
	email string
	iat int64
	exp int64
}