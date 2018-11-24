package keyLibrary

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)



func GeneratePrivPubKey() (*rsa.PrivateKey, error) {

	key, _ := rsa.GenerateKey(rand.Reader, 2048)

	return key, nil
}

func PubKeyEncrypt(pubKey *rsa.PublicKey, message []byte) ([]byte, error) {

	ciphertext, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		pubKey,
		message,
		nil,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return ciphertext, nil
}

func PrivKeyDecrypt(privKey *rsa.PrivateKey, cipherText []byte) ([]byte, error) {
	plainText, err := rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		privKey,
		cipherText,
		nil,
	)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return plainText, nil
}
