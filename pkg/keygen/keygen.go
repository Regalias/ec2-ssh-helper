package keygen

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path"

	"golang.org/x/crypto/ssh"
)

func (km *KeyManager) GenerateKey() (privateKeyPath string, publicKeyData string, err error) {

	// Generate temp key name
	randString, err := getRandomString()
	if err != nil {
		return "", "", fmt.Errorf("failed to get random file name for private key: %v", err)
	}

	privateKeyPath = path.Join(km.privateKeyPath, randString)

	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, km.keyBitSize)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate private key: %v", err)
	}
	if err := privateKey.Validate(); err != nil {
		return "", "", fmt.Errorf("failed to validate generated private key: %v", err)
	}

	// Write key to disk in PEM
	file, err := os.Create(privateKeyPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to create file '%s': %v", privateKeyPath, err)
	}
	if err := os.Chmod(privateKeyPath, 0600); err != nil {
		return "", "", fmt.Errorf("failed to set permissions on created file: %v", err)
	}
	defer file.Close()

	block := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Bytes:   x509.MarshalPKCS1PrivateKey(privateKey),
		Headers: nil,
	}
	if err := pem.Encode(file, &block); err != nil {
		return "", "", fmt.Errorf("could not write private key file: %v", err)
	}

	// Derive public key material
	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", "", fmt.Errorf("failed to load public key data: %v", err)
	}

	publicKeyData = string(ssh.MarshalAuthorizedKey(publicKey))

	return privateKeyPath, publicKeyData, nil
}

func getRandomString() (string, error) {
	buffer := make([]byte, 10)
	if _, err := rand.Read(buffer); err != nil {
		return "", fmt.Errorf("failed to generate random data: %v", err)
	}
	return fmt.Sprintf("%X", buffer), nil
}
