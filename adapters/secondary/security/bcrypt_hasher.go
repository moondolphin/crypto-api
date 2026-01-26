package security

import "golang.org/x/crypto/bcrypt"

type BcryptHasher struct {
	Cost int
}

func NewBcryptHasher(cost int) BcryptHasher {
	if cost <= 0 {
		cost = bcrypt.DefaultCost
	}
	return BcryptHasher{Cost: cost}
}

func (h BcryptHasher) Hash(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), h.Cost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (h BcryptHasher) Compare(hash, password string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {
		return true, nil
	}
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}
	return false, err
}
