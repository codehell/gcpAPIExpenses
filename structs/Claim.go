package structs

type Claim struct {
	iss string
	nickname string
	iat int64
	exp int64
}

func NewClaim() *Claim {
	return new(Claim)
}
