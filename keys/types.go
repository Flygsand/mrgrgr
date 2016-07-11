package keys

type PublicKey struct {
	Body string
}

type KeySource interface {
	PublicKeys() ([]PublicKey, error)
}
