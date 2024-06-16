package data

import (
    "errors"
    "strconv"
    "time"

    "store-api/pkg/data/constant"
)

func NewParamsWrapper() Params {
    return Params{body: make(map[string]interface{})}
}

type Params struct {
    body map[string]interface{}
}

func (p Params) Add(key string, val interface{}) {
    p.body[key] = val
}

func (p Params) GetInt(key string) int {
    if val, ok := p.body[key].(int); ok {
        return val
    }

    if val, ok := p.ConvertInterfaceToInt(key); ok {
        return val
    }

    return 0
}

func (p *Params) GetInt64(key string) int64 {
    if val, ok := p.body[key].(int64); ok {
        return val
    }

    if val, ok := p.ConvertInterfaceToInt64(key); ok {
        return val
    }

    return p.body[key].(int64)
}

func (p Params) GetTime(key string) (time.Time, error) {
    return time.Parse(constant.DefaultDatetimeLayout, p.body[key].(string))
}

func (p Params) GetLocalTime(key string) (time.Time, error) {
    return time.ParseInLocation(constant.DefaultDatetimeLayout, p.body[key].(string), time.Local)
}

func (p Params) GetValue(key string) interface{} {
    return p.body[key]
}

func (p Params) GetMapValue(key string) ([]map[string]interface{}, error) {
    if _, ok := p.body[key].([]map[string]interface{}); !ok {
        return nil, errors.New("failed to convert")
    }

    return p.body[key].([]map[string]interface{}), nil
}

func (p Params) GetString(key string) string {
    if val, ok := p.body[key].(string); ok {
        return val
    } else {
        return ""
    }
}

func (p Params) GetBool(key string) bool {
    return p.body[key].(bool)
}

func (p Params) GetIntOrError(key string) (int, error) {
    val := p.body[key]
    if val != nil {
        return val.(int), nil
    }

    return 0, errors.New("Key not found for " + key)
}

func (p Params) HasKey(key string) bool {
    return p.body[key] != nil
}

func (p *Params) ConvertInterfaceToInt(key string) (int, bool) {
    var i int

    var err error
    switch v := p.body[key].(type) {
    case int64:
        i = int(v)
    case float64:
        i = int(v)
    case string:
        i, err = strconv.Atoi(v)
    default:
        i = 0
    }
    return i, err == nil
}

func (p *Params) ConvertInterfaceToInt64(key string) (int64, bool) {
    var i int64

    var err error
    switch v := p.body[key].(type) {
    case int:
        i = int64(v)
    case float64:
        i = int64(v)
    case string:
        i, err = strconv.ParseInt(v, 10, 64)
    default:
        i = 0
    }
    return i, err == nil
}
