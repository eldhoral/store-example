package security

import (
    "errors"
    "time"

    "github.com/go-playground/validator/v10"
)

type (
    Session struct {
        UserId   string `json:"user_id" validate:"required"`
        Username string `json:"username" validate:"required"`
        Name     string `json:"name" validate:"required"`
        Role     string `json:"role" validate:"required"`
        Iat      int64  `json:"iat" validate:"required"`
        Expired  int64  `json:"exp" validate:"required"`
    }
)

func (ss *Session) IsSessionExpired() error {
    if time.Now().After(time.Unix(ss.Expired, 0)) {
        return errors.New("Expired Session")
    }

    return nil
}

func (ss *Session) ExtendSession(cr Crypto, duration int64) (string, error) {
    ss.Expired = time.Now().Add(time.Duration(duration) * time.Second).Unix()

    return ss.Encrypt(cr)
}

func (ss *Session) Encrypt(cr Crypto) (string, error) {
    enc, _ := cr.Encrypt(ss)

    return string(enc), nil
}

func (ss *Session) Valid() error {
    return nil
}

func NewSession(cr Crypto, session string) (*Session, error) {
    var (
        ss       = &Session{}
        dec, err = cr.Decrypt(ss, session)
    )

    if err != nil {
        return nil, errors.New("Error Session " + err.Error())
    }

    if err := validator.New().Struct(dec); err != nil {
        return nil, errors.New("Error Session " + err.Error())
    }

    if err := ss.IsSessionExpired(); err != nil {
        return nil, err
    }

    return ss, nil
}
