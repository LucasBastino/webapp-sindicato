package password

import "golang.org/x/crypto/bcrypt"

type BcryptHasher struct{}

func NewBcryptHasher() *BcryptHasher{
	return &BcryptHasher{}
}

func (b BcryptHasher) Generate(password string) (hash []byte, err error) {
	return bcrypt.GenerateFromPassword([]byte(password), 14)
}

func (b BcryptHasher) Compare(hash []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(password))
}
