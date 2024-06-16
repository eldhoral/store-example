// Validation
//
// Usage:
//	 func(h Handler) handler(ctx *app.Context) {
//
//		id := ctx.GetVar("id")
//
//		if ctx.HasError() {
//			return h.MobileBadRequest(ctx, ctx.GetFirstError())
//		}

package app

import (
    "context"
    "encoding/json"
    "errors"
    "fmt"
    "net/http"
    "os"
    "strconv"
    "strings"
    "time"

    "store-api/pkg/data/filedata"
    "store-api/pkg/errs"
    "store-api/pkg/helper/realiphelper"
    "store-api/pkg/pagination"

    "github.com/gorilla/mux"
)

const defaultMemory = 32 << 20

// Context extend core Context
// Used for Authentication, and userIdentity used for whole request life cycle
// Note: every request it will generate new Context to hold data between handlers, or services
type Context struct {
    Writer  http.ResponseWriter
    Request *http.Request

    isGuest bool // Check user is authenticated
    ssoID   string
    ip      string

    errors []*errs.Error // Validation data. Use ctx.HasError() to check all params is valid

    hasBody    bool                   // For POST, PUT
    isFormData bool                   // For JSON or FormData
    data       map[string]interface{} // For JSON
    formData   map[string]string      // For FormData

    // Extend more on your demand. Consider small data for performance factor

    startTime time.Time
    isStaging bool
}

func (ctx Context) GetSsoID() string          { return ctx.ssoID }
func (ctx Context) IsGuest() bool             { return ctx.isGuest }
func (ctx Context) GetIP() string             { return ctx.ip }
func (ctx *Context) Context() context.Context { return ctx.Request.Context() }

func (ctx Context) MethodName() string { return ctx.Request.Method }

func (ctx Context) GetElapsed() time.Duration { return time.Since(ctx.startTime) }
func (ctx Context) GetURI() string            { return ctx.Request.RequestURI }
func (ctx Context) IsStaging() bool           { return ctx.isStaging }

func (ctx Context) GetUploadFile(name string) (*filedata.UploadFile, error) {
    return filedata.NewUploadFileInstance(ctx.Request, name)
}

func (ctx Context) GetMultipleUploadFile(name string) ([]filedata.UploadFile, error) {
    return filedata.NewMultipleUploadFile(ctx.Request, name)
}

// NewPaginator
func (ctx *Context) NewPaginator() *pagination.Paginator {
    return pagination.NewPaginator(ctx.Request)
}

func (ctx Context) IsWebRoute() bool {
    return strings.Contains(ctx.Request.RequestURI, "/api/web/")
}

func (ctx Context) IsMobileRoute() bool {
    return strings.Contains(ctx.Request.RequestURI, "/api/v1/")
}

func (ctx *Context) AppendError(err *errs.Error) { ctx.errors = append(ctx.errors, err) }
func (ctx *Context) HasError() bool {
    return ctx.errors != nil && len(ctx.errors) > 0
}

func (ctx Context) Errors() []string {
    if len(ctx.errors) == 0 {
        return nil
    }
    errorStrings := make([]string, len(ctx.errors))
    for i, err := range ctx.errors {
        errorStrings[i] = err.Err.Error()
    }
    return errorStrings
}

func (ctx Context) GetFirstError() error {
    if len(ctx.errors) == 0 {
        return nil
    }
    return ctx.errors[0].Err
}

func (ctx *Context) GetJsonBody() map[string]interface{} {
    if ctx.hasBody {
        return ctx.data
    }
    return nil
}

func (ctx *Context) GetFormBody() map[string]string {
    if ctx.hasBody {
        return ctx.formData
    }
    return nil
}

func (ctx Context) IsContentTypeJson() bool {
    if strings.Contains(ctx.Request.Header.Get("Content-Type"), "application/json") {
        return true
    } else {
        return false
    }
}

// Validate params, short circuit on first parse data fail
// If one params fail to parse value, no further parse data
//
//
// GetVarInt: 	/uri/{key}
// GetQueryInt: /uri/id={key}
//
// GetValue:  	/uri return value as string
// GetValueInt: /uri JSONData {"key": "value"}
// GetValueInt: /uri FormData key=value
//
// Int, Int64, Float64, Bool, String

// ======================== Get Var GET /uri/{key} =====================================================================

// GetVar: /uri/{key} as string
func (ctx Context) GetVar(key string) string {
    vars := mux.Vars(ctx.Request)
    value := vars[key]
    if value != "" {
        return value
    } else {
        ctx.AppendError(errs.NewMissingError(key))
        return ""
    }
}

func (ctx Context) CheckVar(key string) bool {
    if key != "" {
        return true
    } else {
        return false
    }
}

// GetVarInt64: /uri/{key} as int
func (ctx *Context) GetVarInt(key string) int {
    return ctx.stringToInt(key, ctx.GetVar(key))
}

// GetVarInt64: /uri/{key} as int64
func (ctx *Context) GetVarInt64(key string) int64 {
    return ctx.stringToInt64(key, ctx.GetVar(key))
}

// GetVarInt64: /uri/{key} as float64
func (ctx *Context) GetVarFloat64(key string) float64 {
    return ctx.stringToFloat64(key, ctx.GetVar(key))
}

// ======================== Get Query  /uri/{key}={value} ==============================================================

// GetQuery: GET /uri/{key}={value} as string
func (ctx *Context) GetQuery(key string) string {
    if !ctx.Request.URL.Query().Has(key) {
        ctx.AppendError(errs.NewMissingQueryError(key))
        return ""
    }
    vars := ctx.Request.URL.Query().Get(key)
    return vars
}

// GetQueryInt64: /uri/{key}={value} as int
func (ctx *Context) GetQueryInt(key string) int {
    return ctx.stringToInt(key, ctx.GetQuery(key))
}

// GetQueryInt64: /uri/{key}={value} as int64
func (ctx *Context) GetQueryInt64(key string) int64 {
    return ctx.stringToInt64(key, ctx.GetQuery(key))
}

// GetQueryInt64: /uri/{key}={value} as float64
func (ctx *Context) GetQueryFloat64(key string) float64 {
    return ctx.stringToFloat64(key, ctx.GetQuery(key))
}

// GetQueryInt64: /uri/{key}={value} as bool
func (ctx *Context) GetQueryBool(key string) bool {
    return ctx.stringToBool(key, ctx.GetQuery(key))
}

// HasParam Check if parameter contain specific key
func (ctx *Context) HasParam(key string) bool {
    _, hasParam := ctx.Request.Form[key]
    return hasParam
}

// ======================== Post Param json body {"key" : "value", ... : ... } =========================================

// GetValue: POST /uri {"key" : "value"} value as string
func (ctx *Context) GetValue(key string) string {
    if ctx.hasBody {

        if ctx.isFormData {
            if s, exists := ctx.formData[key]; exists {
                return s
            }
        } else {
            if s, exists := ctx.data[key]; exists {
                val, ok := s.(string)
                if ok {
                    return val
                } else {
                    ctx.AppendError(errs.NewParseError(key))
                }
            }
        }

        ctx.AppendError(errs.NewMissingError(key))
    }
    return ""
}

// GetValueInt: POST /uri {"key" : "value"} value as int
func (ctx *Context) GetValueInt(key string) int {
    if ctx.isFormData {
        return ctx.stringToInt(key, ctx.getFormData(key))
    } else {
        return ctx.anyToInt(key, ctx.getJSONData(key))
    }
}

// GetValueInt64: POST /uri {"key" : "value"} value as int64
func (ctx *Context) GetValueInt64(key string) int64 {
    if ctx.isFormData {
        return ctx.stringToInt64(key, ctx.getFormData(key))
    } else {
        return ctx.anyToInt64(key, ctx.getJSONData(key))
    }
}

// GetValueFloat64: POST /uri {"key" : "value"} value as float64
func (ctx *Context) GetValueFloat64(key string) float64 {
    if ctx.isFormData {
        return ctx.stringToFloat64(key, ctx.getFormData(key))
    } else {
        return ctx.anyToFloat64(key, ctx.getJSONData(key))
    }
}

// GetValueBool: POST /uri {"key" : "value"} value as bool
func (ctx *Context) GetValueBool(key string) bool {
    if ctx.isFormData {
        return ctx.stringToBool(key, ctx.getFormData(key))
    } else {
        return ctx.anyToBool(key, ctx.getJSONData(key))
    }
}

// ======================== Parse JSON or Form data, depend on request content-type ====================================

// ParseJson for convenience, auto parse into Context based on Content-type
func (ctx *Context) ParseJson() {

    if ctx.IsStaging() {
        fmt.Printf("[Debug] Parse JSON. Content-type: %s - %s %s\n",
            ctx.Request.Header.Get("Content-type"), ctx.Request.Method, ctx.Request.URL)
    }

    ctx.isFormData = false
    if ctx.hasBody {
        decoder := json.NewDecoder(ctx.Request.Body)
        err := decoder.Decode(&ctx.data)
        if err != nil {
            errString := err.Error()

            if errString != "EOF" && errString != "unexpected EOF" {
                ctx.AppendError(errs.NewMalformedJSONError(fmt.Errorf("invalid json body request: %v", err)))
            }
        }

    } else {
        ctx.data = make(map[string]interface{})
    }
}

func (ctx *Context) ParseForm() {
    if ctx.IsStaging() {
        fmt.Printf("[Debug] Parse Form. Content-type: %s - %s %s\n",
            ctx.Request.Header.Get("Content-type"), ctx.Request.Method, ctx.Request.URL)
    }

    ctx.isFormData = true
    if err := ctx.Request.ParseForm(); err != nil {
        ctx.AppendError(errs.NewMalformedFormError(err))
    }

    if err := ctx.Request.ParseMultipartForm(defaultMemory); err != nil && !errors.Is(err, http.ErrNotMultipart) {
        ctx.AppendError(errs.NewMalformedFormError(err))
    }

    formData := make(map[string]string)
    for k, v := range ctx.Request.Form {
        val := v[len(v)-1]
        formData[k] = val // pick last
    }
    ctx.formData = formData
}

// ======================== NewContext =================================================================================

func NewContext(rw http.ResponseWriter, r *http.Request, isStaging bool) *Context {

    ctx := &Context{
        Writer:     rw,
        Request:    r,
        hasBody:    r.Method != http.MethodGet,
        isFormData: false,

        ip:        realiphelper.FromRequest(r),
        startTime: time.Now(),
        isStaging: isStaging,
    }

    if strings.Contains(ctx.Request.Header.Get("Content-type"), "application/json") {
        ctx.ParseJson()
    } else {
        ctx.ParseForm()
    }

    return ctx
}

// Private method to param convert =====================================================================================

func (ctx *Context) stringToInt(key, value string) int {
    if !ctx.HasError() { // short circuit, no more parse data
        val, err := strconv.Atoi(value)
        if err != nil {
            ctx.AppendError(errs.NewParseError(key))
        }
        return val
    }
    return 0
}

func (ctx *Context) stringToInt64(key, value string) int64 {
    if !ctx.HasError() { // short circuit, no more parse data
        val, err := strconv.ParseInt(value, 10, 64)
        if err != nil {
            ctx.AppendError(errs.NewParseError(key))
        }
        return val
    }
    return 0
}

func (ctx *Context) stringToFloat64(key, value string) float64 {
    if !ctx.HasError() { // short circuit, no more parse data
        val, err := strconv.ParseFloat(value, 64)
        if err != nil {
            ctx.AppendError(errs.NewParseError(key))
        }
        return val
    }
    return 0.0
}

func (ctx *Context) stringToBool(key, value string) bool {
    if !ctx.HasError() { // short circuit, no more parse data
        val, err := strconv.ParseBool(value)
        if err != nil {
            ctx.AppendError(errs.NewParseError(key))
        }
        return val
    }
    return false
}

func (ctx *Context) anyToInt(key string, value interface{}) int {
    if !ctx.HasError() { // short circuit, no more parse data
        return int(ctx.anyToFloat64(key, value))
    }
    return 0
}

func (ctx *Context) anyToInt64(key string, value interface{}) int64 {
    if !ctx.HasError() { // short circuit, no more parse data
        return int64(ctx.anyToFloat64(key, value))
    }
    return 0
}

func (ctx *Context) anyToFloat64(key string, value interface{}) float64 {
    if !ctx.HasError() { // short circuit, no more parse data
        val, ok := value.(float64)
        if ok {
            return val
        } else {
            ctx.AppendError(errs.NewParseError(key))
        }
    }
    return 0
}

func (ctx *Context) anyToBool(key string, value interface{}) bool {
    if !ctx.HasError() { // short circuit, no more parse data
        val, ok := value.(bool)
        if ok {
            return val
        } else {
            ctx.AppendError(errs.NewParseError(key))
        }
    }
    return false
}

func (ctx *Context) getFormData(key string) string {
    //if ctx.formData == nil {
    //	ctx.AppendError(errs.NewError(
    //		errors.New("ctx.data not init. Call ParseJSON() or ParseForm() first"), "data body"))
    //	return ""
    //}

    if ctx.hasBody {

        if value, exists := ctx.formData[key]; exists {
            return value
        } else {
            ctx.AppendError(errs.NewMissingKeyError(key))
        }

    }
    return ""
}

func (ctx *Context) GetRequestInfo() string {
    if ctx.IsStaging() {
        return fmt.Sprintf("%s http://%s%s - elapse: %s",
            ctx.Request.Method, ctx.Request.Host, ctx.Request.URL.String(), ctx.GetElapsed())
    }
    return fmt.Sprintf("%s http://%s%s - elapse: %s",
        ctx.Request.Method, ctx.Request.Host, ctx.Request.URL.String(), ctx.GetElapsed())
}

func (ctx *Context) getJSONData(key string) interface{} {
    if ctx.hasBody {
        if value, exists := ctx.data[key]; exists {
            return value
        } else {
            ctx.AppendError(errs.NewMissingKeyError(key))
        }

    }
    return nil
}

func (ctx *Context) DebugGetPayload() string {

    // Note: If APP_DEBUG show for production as well
    if os.Getenv("APP_DEBUG") == "True" {

        if ctx.hasBody {
            note := ":warning: NOTE: APP_DEBUG=True"
            if ctx.IsContentTypeJson() {
                return fmt.Sprintf("\nPayload: %v - %s\n", ctx.GetJsonBody(), note)
            } else {
                return fmt.Sprintf("\nPayload: %v - %s\n", ctx.GetFormBody(), note)
            }
        }
    }
    return ""
}
