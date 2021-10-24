package keygen

type KeyManager struct {
	privateKeyPath string
	keyBitSize     int
}

func New(privateKeyPath string) *KeyManager {
	return &KeyManager{
		keyBitSize:     4096,
		privateKeyPath: privateKeyPath,
	}
}
