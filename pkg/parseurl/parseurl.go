package parseurl

import (
    "bytes"
    "io/ioutil"
    "net/http"
    "time"

    "store-api/internal/application_form/domain/access"
    "store-api/internal/application_form/presenter/req"
    "store-api/pkg/data/constant"
    "store-api/pkg/security"

    jsoniter "github.com/json-iterator/go"
    "github.com/sirupsen/logrus"
)

func GetJson(url string, target interface{}, access *access.Access) (int, interface{}, error) {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        logrus.Error(err)
        return http.StatusInternalServerError, req, err
    }

    requestID := security.GenerateRandomStringNumeric(13)

    req.Header.Set("Authorization", access.Token)
    req.Header.Set("AuthID", access.AuthID)
    req.Header.Set("X-CHANNEL", access.XChannel)
    req.Header.Set("X-REQUEST-ID", requestID)
    req.Header.Set("X-TIMESTAMP", access.XTimestamp)

    var myClient = &http.Client{Timeout: 10 * time.Second}
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

func EvaluateToken(url string, target interface{}, access access.Access) (int, interface{}, error) {
    req, err := http.NewRequest("POST", url, nil)
    if err != nil {
        logrus.Error(err)
        return http.StatusInternalServerError, req, err
    }
    requestID := security.GenerateRandomStringNumeric(13)

    req.Header.Set("Authorization", access.Token)
    req.Header.Set("AuthID", access.AuthID)
    req.Header.Set("X-CHANNEL", access.XChannel)
    req.Header.Set("X-REQUEST-ID", requestID)
    req.Header.Set("X-REQUEST-NAME", constant.XRequestNameLogin)
    req.Header.Set("X-SIGNATURE", access.Signature)
    req.Header.Set("X-TIMESTAMP", access.XTimestamp)

    var myClient = &http.Client{Timeout: 10 * time.Second}
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

func InquiryByNikAndChannel(url string, target interface{}, dataCustomer req.DataCustomer) (int, interface{}, error) {
    req, err := http.NewRequest("POST", url, nil)
    if err != nil {
        logrus.Error(err)
        return http.StatusInternalServerError, req, err
    }

    req.Header.Set("Cookie", constant.CookiePath)
    req.Header.Set("X-Device-ID", dataCustomer.XDeviceID)

    q := req.URL.Query()
    q.Set("channel", dataCustomer.Channel)
    q.Set("personalIdNumber", dataCustomer.PersonalIDNumber)
    q.Set("personalIdType", dataCustomer.PersonalIDType)
    req.URL.RawQuery = q.Encode()

    var myClient = &http.Client{Timeout: 10 * time.Second}
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

func MemberInfo(url string, target interface{}, access *access.Access) (int, interface{}, error) {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        logrus.Error(err)
        return http.StatusInternalServerError, req, err
    }
    requestID := security.GenerateRandomStringNumeric(13)

    req.Header.Set("Authorization", access.Token)
    req.Header.Set("AuthID", access.AuthID)
    req.Header.Set("X-CHANNEL", access.XChannel)
    req.Header.Set("X-REQUEST-ID", requestID)
    req.Header.Set("X-REQUEST-NAME", constant.GetMemberInfo)
    req.Header.Set("X-SIGNATURE", access.Signature)
    req.Header.Set("X-TIMESTAMP", access.XTimestamp)

    var myClient = &http.Client{Timeout: 10 * time.Second}
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

func MemberUpgradeInfo(url string, target interface{}, access *access.Access) (int, interface{}, error) {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        logrus.Error(err)
        return http.StatusInternalServerError, req, err
    }
    requestID := security.GenerateRandomStringNumeric(13)

    req.Header.Set("Authorization", access.Token)
    req.Header.Set("AuthID", access.AuthID)
    req.Header.Set("X-CHANNEL", access.XChannel)
    req.Header.Set("X-REQUEST-ID", requestID)
    req.Header.Set("X-REQUEST-NAME", constant.GetMemberUpgradeInfo)
    req.Header.Set("X-SIGNATURE", access.Signature)
    req.Header.Set("X-TIMESTAMP", access.XTimestamp)

    var myClient = &http.Client{Timeout: 10 * time.Second}
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

func DroppedInEmail(url string, target interface{}, request req.SendNotifDroppedInEmail) (int, interface{}, error) {
    b, err := jsoniter.Marshal(request)
    if err != nil {
        logrus.Error("Error when marshaling struct : " + err.Error())
        return http.StatusInternalServerError, nil, err
    }
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

func DroppedInAppNotif(url string, target interface{}, request req.SendNotifDroppedInAppNotif, access *access.Access) (int, interface{}, error) {
    b, err := jsoniter.Marshal(request)
    if err != nil {
        logrus.Error("Error when marshaling struct : " + err.Error())
        return http.StatusInternalServerError, nil, err
    }
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
    if err != nil {
        logrus.Error(err)
        return http.StatusInternalServerError, req, err
    }
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", access.Token)
    req.Header.Set("AuthID", access.AuthID)
    req.Header.Set("X-REQUEST-NAME", constant.XRequestNameNotif)
    req.Header.Set("X-CHANNEL", access.XChannel)

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
