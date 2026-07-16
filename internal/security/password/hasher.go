package password

type Hasher interface {
	Generate(password string) (hash []byte, err error)
	Compare(hash []byte, password string) error
}