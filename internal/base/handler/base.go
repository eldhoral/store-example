package handler

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "os"

    "store-api/pkg/metric"

    "gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

    "store-api/internal/base/app"
    storeService "store-api/internal/store/service"
    "store-api/pkg/errs"
    "store-api/pkg/httpclient"
    "store-api/pkg/server"

    "github.com/jmoiron/sqlx"
    "github.com/sirupsen/logrus"
)

type HandlerFn func(*app.Context) *server.Response

type BaseHTTPHandler struct {
    Handlers         interface{}
    DB               *sqlx.DB
    HTTPClient       httpclient.Client
    Params           map[string]string
    StatsdMonitoring metric.StatsdMonitoring
    StoreService     storeService.StoreService
}

func NewBaseHTTPHandler(
        db *sqlx.DB,
        httpClient httpclient.Client,
        params map[string]string,
        statsdMonitoring metric.StatsdMonitoring,
        storeSvc storeService.StoreService,

) *BaseHTTPHandler {

    return &BaseHTTPHandler{DB: db, HTTPClient: httpClient, Params: params,
        StatsdMonitoring: statsdMonitoring,
        StoreService:     storeSvc,
    }
}

// AsJson to response custom message: 200, 201 with message (Mobile use 500 error)
func (h BaseHTTPHandler) AsJson(ctx *app.Context, status int, message string, data interface{}) *server.Response {
    serverResponse := server.Response{
        Status:       status,
        Message:      message,
        Data:         data,
        Version:      os.Getenv("APP_VERSION"),
        ResponseType: server.DefaultResponseType,
    }

    return &serverResponse
}

// AsMobileStatusOK always return httpStatus: 200, but Status field: 500,400,200...
func (h BaseHTTPHandler) AsMobileStatusOK(ctx *app.Context, status int, message string, data interface{}) *server.Response {
    serverResponse := server.Response{
        Status:       status,
        Message:      message,
        Data:         data,
        Version:      os.Getenv("APP_VERSION"),
        ResponseType: server.DefaultResponseType,
    }

    return &serverResponse
}

// AsMobileJsonSetStatusCode for special case, in catch() it call       'status' => $this->setStatusCode(400),
func (h BaseHTTPHandler) AsMobileJsonSetStatusCode(ctx *app.Context, status int, message string, data interface{}) *server.Response {
    return &server.Response{
        Status:       status,
        Message:      message,
        Data:         data,
        Version:      os.Getenv("APP_VERSION"),
        ResponseType: server.MobileSetStatusCodeType,
    }
}

// ThrowExceptionJson for some exception not handle in Yii2 framework
func (h BaseHTTPHandler) ThrowExceptionJson(ctx *app.Context, status, code int, name, message string) *server.Response {
    return &server.Response{
        Status:       status,
        Message:      "",
        Data:         server.NotAllowedMethod{Name: name, Message: message, Code: code, Status: status},
        Version:      "",
        Log:          nil,
        ResponseType: server.MobileNotAllowedType,
    }
}

// AsJsonWithLog for custom log
func (h BaseHTTPHandler) AsJsonWithLog(status int, message string, data interface{}, log *server.LogMessage) *server.Response {
    return &server.Response{Status: status, Message: message, Data: data, Version: os.Getenv("APP_VERSION"), Log: log}
}

func (h BaseHTTPHandler) IsStaging() bool {
    return h.Params["APP_ENV"] == "development"
}

func (h BaseHTTPHandler) IsProd() bool {
    return h.Params["APP_ENV"] == "production"
}

func (h BaseHTTPHandler) GetParam(key string) string {
    return h.Params[key]
}

// RunAction entry point to handle route.
func (h BaseHTTPHandler) RunAction(fn HandlerFn) http.HandlerFunc {
    return h.CapturePanic(h.Execute(fn))
}

// SendPanicFlock used only for CapturePanic() to send some clue
func (h BaseHTTPHandler) SendPanic(r *http.Request, errMsg string, err interface{}) {
    logrus.Errorln(errMsg) // Need to notify

    errStack, file := errs.StackAndFile(3)
    errInfo := fmt.Sprintf("\n:red_circle paylater-customer-api service \n* MUST FIX :boom: :boom: :boom: "+
            "Panic Error: %v*\nRequest: %s %s\nAccess Token: %s\nFile: %s", err, r.Method, r.RequestURI, r.Header.Get("Authorization"), file)
    msg := fmt.Sprintf("%s\n\nStack trace: \n%s...", errInfo, errStack)

    fmt.Println("\nPANIC:", msg)
}

// Execute SpecificHandler.Method(ctx *app.Context)
func (f BaseHTTPHandler) Execute(handler HandlerFn) http.HandlerFunc {
    return func(rw http.ResponseWriter, r *http.Request) {

        // 1. Authentication
        ctx, _ := f.Authentication(rw, r)
        //if err != nil {
        //    span, _ := tracer.StartSpanFromContext(r.Context(), "Unauthorized", tracer.ResourceName(r.RequestURI))
        //    defer span.Finish(tracer.WithError(fmt.Errorf("accessToken: %s. Err: %v", r.Header.Get("accessToken"), err)))
        //
        //    WriteJSON(rw, http.StatusUnauthorized, server.MobileResponse{
        //        Status:  http.StatusUnauthorized,
        //        Message: err.Error(),
        //        Data:    check, // Show more err for debug instead of empty []
        //        Version: os.Getenv("APP_VERSION"),
        //    })
        //    return
        //}

        // 2. Capture handler error to avoid infinite loop SendFlock
        defer func() {
            if err0 := recover(); err0 != nil {
                f.SendPanic(r, "CaptureLastPanic NEED TO FIX NOW", err0)

                WriteJSON(rw, http.StatusInternalServerError,
                    "Request is halted unexpectedly, please contact the administrator.")
            }
        }()

        // 3. Process route action, and return *server.Response
        resp := handler(ctx)
        httpStatus := resp.GetStatus()

        // Start Timer
        //start := carbon.Now().ToDateTimeString()

        //header, err := ctx.GetAccess()
        //if err != nil {
        //    logrus.Error("Error getting header for sending to graylog : " + err.Error())
        //}

        //stringRequestHeader, err := jsoniter.MarshalToString(header)
        //if err != nil {
        //    logrus.Error("Error when marshaling request header to string for sending to graylog : " + err.Error())
        //}
        //
        //stringResponse, err := jsoniter.MarshalToString(resp)
        //if err != nil {
        //    logrus.Error("Error when marshaling response to string for sending to graylog : " + err.Error())
        //}
        //
        //logrus.WithFields(logrus.Fields{
        //    "clientip":        ctx.GetIP(),
        //    "timestamp":       start,
        //    "method":          ctx.Request.Method,
        //    "path":            ctx.GetURI(),
        //    "statuscode":      httpStatus,
        //    "latency":         ctx.GetElapsed().Milliseconds(),
        //    "request_header":  stringRequestHeader,
        //    "request":         ctx.GetRequestInfo(),
        //    "response_header": rw.Header(),
        //    "response":        stringResponse,
        //}).Info()

        // Except mobile with server.MobileStatusOKType always response httpStatus 200
        if httpStatus >= http.StatusInternalServerError && resp.ResponseType != server.MobileStatusOKType {
            // Send more clue for other internal app can debug.
            WriteJSON(rw, httpStatus, server.NotAllowedMethod{
                Name:    "Server Error",
                Message: resp.Message,
                Code:    0,
                Status:  httpStatus,
            })
            return
        }

        if httpStatus >= 300 {
            span, _ := tracer.StartSpanFromContext(ctx.Context(), "response", tracer.ResourceName(ctx.Request.RequestURI))
            defer span.Finish(tracer.WithError(fmt.Errorf("%v %v", resp.Message, resp.Data)))
        }

        if f.IsStaging() {
            fmt.Printf("INFO: %s   - code: %d\n\n", ctx.GetRequestInfo(), httpStatus)
        }

        // 5. Response JSON data for web route, and mobile route
        if ctx.IsWebRoute() {
            if resp.ResponseType == server.MobileNotAllowedType {
                WriteJSON(rw, httpStatus, resp.GetMobileException())
            } else if resp.ResponseType == server.StreamResponseType {
                WriteStream(rw, resp.Message, resp.GetStream())
            } else {
                WriteJSON(rw, httpStatus, server.WebResponse{
                    Status:  httpStatus,
                    Message: resp.Message,
                    Data:    resp.Data})
            }
            return
        } else {
            if resp.ResponseType == server.MobileNotAllowedType {
                // CodeResponseRequest: happen due to Throw Exception without handle. Or in catch().
                // Often has Code field
                WriteJSON(rw, httpStatus, resp.GetMobileException())
                return
            } else if resp.ResponseType == server.MobileSetStatusCodeType {
                WriteJSON(rw, httpStatus, server.MobileResponse{
                    Status:  httpStatus,
                    Message: resp.Message,
                    Data:    resp.Data,
                    Version: resp.Version})
                return

            } else if resp.ResponseType == server.StreamResponseType {
                WriteStream(rw, resp.Message, resp.GetStream())
            } else {
                // MobileResponse status always http.StatusOK but Status field 200,403,400,500,...
                WriteJSON(rw, http.StatusOK, server.MobileResponse{
                    Status:  httpStatus,
                    Message: resp.Message,
                    Data:    resp.Data,
                    Version: resp.Version})
                return
            }

            // Can extend more type of response here
        }
    }
}

func (h BaseHTTPHandler) Authentication(rw http.ResponseWriter, r *http.Request) (*app.Context, error) {
    return app.NewContext(rw, r, h.IsStaging()), nil
}

// CapturePanic Last layer to capture panic which might halt the whole application.
func (h BaseHTTPHandler) CapturePanic(next http.HandlerFunc) http.HandlerFunc {
    return func(rw http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                h.SendPanic(r, "CaptureLastPanic NEED TO FIX NOW", err)

                WriteJSON(rw, http.StatusInternalServerError,
                    "Request is halted unexpectedly, please contact the administrator.")
            }
        }()
        next(rw, r)
    }
}

func WriteJSON(rw http.ResponseWriter, httpCode int, data interface{}) {
    rw.Header().Set("Content-Type", "application/json")
    if httpCode != 0 {
        rw.WriteHeader(httpCode)
    }
    if data != nil {
        _ = json.NewEncoder(rw).Encode(data)
    }
}
func WriteStream(rw http.ResponseWriter, fileName string, b bytes.Buffer) {
    rw.Header().Set("File-Name", fileName)
    rw.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")
    rw.Header().Set("Content-Type", "application/octet-stream")
    rw.Header().Set("Content-Length", fmt.Sprintf("%d", b.Len()))
    rw.Header().Set("Content-Transfer-Encoding", "binary")
    rw.Header().Set("Expires", "0")

    rw.WriteHeader(http.StatusOK)

    _, err := rw.Write(b.Bytes())
    if err != nil {
        logrus.Errorln(fmt.Errorf("Stream %s %v", fileName, err))
        return
    }
}

// common missing id response, can be used by any endpoint route
func (h BaseHTTPHandler) NotFoundHandler(w http.ResponseWriter, r *http.Request) {
    WriteJSON(w, http.StatusBadRequest, server.NotAllowedMethod{
        Name:    "Bad Request",
        Message: "Missing required parameters: id",
        Code:    0,
        Status:  http.StatusBadRequest,
    })
}

func (h BaseHTTPHandler) NotFoundHandlerNIK(w http.ResponseWriter, r *http.Request) {
    WriteJSON(w, http.StatusBadRequest, server.NotAllowedMethod{
        Name:    "Bad Request",
        Message: "Missing required parameters: nik",
        Code:    0,
        Status:  http.StatusBadRequest,
    })
}

// common method not allowed response (405)
func (h BaseHTTPHandler) MethodNotAllowedHandler() http.Handler {
    return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
        WriteJSON(rw, http.StatusMethodNotAllowed, server.NotAllowedMethod{
            Name: "Method Not Allowed",
            // hardcode the allowed method list for now
            Message: fmt.Sprintf("Method Not Allowed. This URL can only handle the following request methods: %s", "GET"),
            Code:    0,
            Status:  http.StatusMethodNotAllowed,
        })
    })
}

func (h BaseHTTPHandler) NotFoundRoute(ctx *app.Context) *server.Response {
    return h.AsJson(ctx, http.StatusNotFound, "page not found", nil)
}

func (h BaseHTTPHandler) MissingIDParameter(ctx *app.Context) *server.Response {
    return h.AsJson(ctx, http.StatusBadRequest, "page not found", nil)
}
