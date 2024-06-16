package security

import (
    "crypto/aes"
    "crypto/cipher"
    "crypto/rand"
    "encoding/hex"
    "io"

    "github.com/golang-jwt/jwt"
)

type (
    Crypto interface {
        Encrypt(claims jwt.Claims) ([]byte, error)
        Decrypt(claims jwt.Claims, text string) (jwt.Claims, error)

        EncryptAes(text string) ([]byte, error)
        DecryptAes(text string) ([]byte, error)
    }

    impl struct {
        secret []byte
    }
)

func (i *impl) Encrypt(claims jwt.Claims) ([]byte, error) {
    token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), claims)
    ciphertext, err := token.SignedString(i.secret)
    if err != nil {
        return nil, err
    }

    return []byte(ciphertext), nil
}

func (i *impl) Decrypt(claims jwt.Claims, tokenString string) (jwt.Claims, error) {
    keyFunc := func(token *jwt.Token) (interface{}, error) {
        return i.secret, nil
    }

    token, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)
    if err != nil {
        return nil, err
    }

    return token.Claims, nil
}

func (i *impl) EncryptAes(text string) ([]byte, error) {
    block, _ := aes.NewCipher(i.secret)
    gcm, err := cipher.NewGCM(block)

    if err != nil {
        return nil, err
    }

    nonce := make([]byte, gcm.NonceSize())
    if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }

    ciphertext := gcm.Seal(nonce, nonce, []byte(text), nil)
    return ciphertext, nil
}

func (i *impl) DecryptAes(text string) ([]byte, error) {
    data, err := hex.DecodeString(text)

    if err != nil {
        return nil, err
    }

    block, err := aes.NewCipher(i.secret)

    if err != nil {
        return nil, err
    }
    gcm, err := cipher.NewGCM(block)

    if err != nil {
        return nil, err
    }

    nonceSize := gcm.NonceSize()
    nonce, ciphertext := data[:nonceSize], data[nonceSize:]

    plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
    if err != nil {
        return nil, err
    }

    return []byte(plaintext), nil
}

func New(secret string) (Crypto, error) {
    return &impl{[]byte(secret)}, nil
}
