package otp

import (
    "bytes"
    "io/ioutil"
    "net/http"
    "time"

    "store-api/internal/application_form/presenter/req"

    jsoniter "github.com/json-iterator/go"
    "github.com/sirupsen/logrus"
)

func OTP(url string, target interface{}, request req.RegistrationPaylaterOTP) (int, interface{}, error) {
    b, err := jsoniter.Marshal(request)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
    if err != nil {
        logrus.Error(err)
        return http.StatusInternalServerError, req, err
    }
    req.Header.Set("Content-Type", "application/json")

    myClient := &http.Client{Timeout: 10 * time.Second}
    r, err := myClient.Do(req)
    if err != nil {
        logrus.Error(err)
        return http.StatusInternalServerError, r, err
    }

    defer r.Body.Close()

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        logrus.Error(err)
        return http.StatusInternalServerError, body, err
    }
    err = jsoniter.Unmarshal(body, &target)
    if err != nil {
        logrus.Error(err)
        return http.StatusInternalServerError, nil, err
    }
    return r.StatusCode, target, nil
}

func ValidateOTP(url string, target interface{}, request req.ValidateRegistrationOTP) (int, interface{}, error) {
    b, err := jsoniter.Marshal(request)

    req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
    if err != nil {
        logrus.Error(err)
        return http.StatusInternalServerError, req, err
    }

    req.Header.Set("Content-Type", "application/json")

    myClient := &http.Client{Timeout: 10 * time.Second}
    r, err := myClient.Do(req)
    if err != nil {
        logrus.Error(err)
        return http.StatusInternalServerError, r, err
    }

    defer r.Body.Close()

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        logrus.Error(err)
        return http.StatusInternalServerError, body, err
    }
    err = jsoniter.Unmarshal(body, &target)
    if err != nil {
        logrus.Error(err)
        return http.StatusInternalServerError, nil, err
    }
    return r.StatusCode, target, nil
}
