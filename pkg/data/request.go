package data

import (
    "context"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "reflect"
    "strconv"

    "store-api/pkg/data/filedata"

    "github.com/gorilla/mux"
    "github.com/sirupsen/logrus"
)

// Deprecated: NewRequestWrapper. Use ctx *app.Context instead
func NewRequestWrapper(r *http.Request) *Request {

    hasBody := false

    if r.Method != http.MethodGet {
        hasBody = true
        body := make(map[string]interface{})

        raw, err := ioutil.ReadAll(r.Body)
        if err != nil {
            logrus.Warning(fmt.Errorf("Deprecated: NewRequestWrapper: Parse Json Fail: %v", err))
        }
        err = json.Unmarshal(raw, &body)
        if err != nil {
            logrus.Warning(fmt.Errorf("Deprecated: NewRequestWrapper: Parse Json Fail: %v", err))
        }

        return &Request{r: r, hasBody: hasBody, body: body}
    }

    return &Request{r: r, hasBody: hasBody}
}

type Request struct {
    r       *http.Request
    hasBody bool // For POST, PUT,
    body    map[string]interface{}

    hasError  bool
    errors    []string
    errorCode int
}

func (req Request) Context() context.Context {
    return context.Background()
}

func (req Request) GetUploadFile(name string) (*filedata.UploadFile, error) {
    return filedata.NewUploadFileInstance(req.r, name)
}

func (req Request) GetMultipleUploadFile(name string) ([]filedata.UploadFile, error) {
    return filedata.NewMultipleUploadFile(req.r, name)
}

func (req Request) GetVarInt(key string) (int, error) {
    vars := mux.Vars(req.r)
    id, err := strconv.Atoi(vars[key])

    if err != nil {
        return 0, err
    }
    return id, nil
}

func (req Request) GetVar(key string) string {
    vars := mux.Vars(req.r)
    return vars[key]
}

func (req Request) GetInt(key string) int {
    if req.hasBody { // For POST method
        if req.body[key] != nil && req.isIntValue(req.body[key]) {
            return req.body[key].(int)
        }
        if req.isFloatValue(req.body[key]) {
            return int(req.body[key].(float64))
        }
        // Handle string int eg: "123"
        val, err := strconv.Atoi(req.body[key].(string))
        if err != nil {
            val = 0
        }
        return val
    }

    val, err := strconv.Atoi(req.r.FormValue(key))
    if err != nil {
        val = 0
    }
    return val
}

func (req Request) GetFloat(key string) float64 {
    if req.hasBody {
        if req.isFloatValue(req.body[key]) {
            return req.body[key].(float64)
        }
    }

    val, err := strconv.ParseFloat(req.r.FormValue(key), 64)
    if err != nil {
        val = 0.0
    }
    return val
}

func (req Request) GetValue(key string) string {
    if req.hasBody { // For POST method
        if req.body[key] != nil && req.isStringValue(req.body[key]) {
            return req.body[key].(string)
        }

        return ""
    }
    return req.r.FormValue(key)
}

func (req Request) GetBool(key string) bool {
    if req.hasBody { // For POST method
        if req.body[key] != nil && req.isBoolValue(req.body[key]) {
            return req.body[key].(bool)
        }

        return false
    }
    return false
}

func (req Request) GetBody() map[string]interface{} {
    if req.hasBody {
        return req.body
    }
    return nil
}

func (req Request) GetSsoID() string {
    return req.r.Header.Get("X-SSO-User-ID")
}

func (req Request) isStringValue(v interface{}) bool {
    switch reflect.TypeOf(v).Kind() {
    case reflect.String:
        return true
    }
    return false
}

func (req Request) isIntValue(v interface{}) bool {
    switch reflect.TypeOf(v).Kind() {
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        return true
    }
    return false
}

func (req Request) isFloatValue(v interface{}) bool {
    switch reflect.TypeOf(v).Kind() {
    case reflect.Float32, reflect.Float64:
        return true
    }
    return false
}

func (req Request) isBoolValue(v interface{}) bool {
    switch reflect.TypeOf(v).Kind() {
    case reflect.Bool:
        return true
    }
    return false
}

func (req Request) GetQueryVar(key string) string {
    vars := req.r.URL.Query().Get(key)
    return vars
}

func (req Request) GetQueryVarInt(key string) int {
    val, err := strconv.Atoi(req.r.URL.Query().Get(key))
    if err != nil {
        val = 0
    }
    return val
}
