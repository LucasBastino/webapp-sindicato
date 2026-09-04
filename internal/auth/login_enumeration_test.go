package auth

import (
	"testing"

	"github.com/LucasBastino/webapp-sindicato/internal/security/password"
)

func TestLoginDummyHashRunsBcryptCompare(t *testing.T) {
	hasher := password.NewBcryptHasher()
	err := hasher.Compare(loginDummyPasswordHash, "not-the-dummy-password")
	if err == nil {
		t.Fatal("dummy hash compare should fail for wrong password")
	}
}
