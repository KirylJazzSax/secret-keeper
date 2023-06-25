package encryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

type Encryptor interface {
	Encrypt(secret string, encTo *string) error
	Decrypt(encrypted string, decTo *string) error
}

type SimpleEncryptor struct {
	secretKey string
	iv        string
}

func (e *SimpleEncryptor) Encrypt(secret string, encTo *string) error {
	block, err := aes.NewCipher([]byte(e.secretKey))
	if err != nil {
		return err
	}
	plainText := []byte(secret)
	cfb := cipher.NewCFBEncrypter(block, []byte(e.iv))
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	*encTo = base64.StdEncoding.EncodeToString(cipherText)
	return nil
}

func (e *SimpleEncryptor) Decrypt(encrypted string, decTo *string) error {
	block, err := aes.NewCipher([]byte(e.secretKey))
	if err != nil {
		return err
	}

	cipherText, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return err
	}
	cfb := cipher.NewCFBDecrypter(block, []byte(e.iv))
	secret := make([]byte, len(cipherText))
	cfb.XORKeyStream(secret, cipherText)
	*decTo = string(secret)
	return nil
}

func NewSimpleEncryptor(key, iv string) *SimpleEncryptor {
	return &SimpleEncryptor{
		secretKey: key,
		iv:        iv,
	}
}
