package server

import "bytes"

//App is a server application abstraction,
//can be HTTP, Consumer (kafka/pub-sub), Redis, etc
type App interface {
	Run() error
}

type LogMessage struct {
	CompanyID  int64
	EmployeeID int64
	Email      string
	Log        string
}

const (
	DefaultResponseType     = iota // For mobile and web
	MobileStatusOKType             // For mobile, always return 200
	MobileSetStatusCodeType        // For mobile, in catch() sometimes, we see $this->setStatusCode(400),
	MobileNotAllowedType           // For MobileBadRequest and Method NotAllowMethod response
	StreamResponseType             // For Excel

	// Extend more custom response your need
)

// MobileResponse always http.StatusOK but Status field 200,403,400.
type MobileResponse struct {
	Message string      `json:"message"`
	Status  int         `json:"status"`
	Data    interface{} `json:"data"`
	Version string      `json:"version"`
}

// WebResponse different format with Mobile Response
type WebResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Yii2HTTPException for yii\web\HttpException
// Backward compatibility with Yii2, Throw new BadRequestHttpException, MethodNotAllowedHttpException, without handle it
// This struct add this to Data field of Response
type NotAllowedMethod struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  int    `json:"status"`
	//Type    string `json:"type"`  // "yii\\web\\BadRequestHttpException". No need this on golang
}

// Reponse to save log to ULMS, used for web and mobile. Extend more custom
type Response struct {
	Status       int
	Message      string
	Data         interface{}
	Version      string
	Log          *LogMessage
	ResponseType int
}

func (r *Response) IsOk() bool {
	return r.Status >= 200 && r.Status < 300
}

func (r *Response) GetDataAsString() string {
	if val, ok := r.Data.(string); ok {
		return val
	}
	return ""
}

// GetStatus return http.StatusOK, http.StatusBadRequest, ...
func (r Response) GetStatus() int {
	return r.Status
}

func (r Response) GetMobileException() interface{} {
	return r.Data
}

func (r Response) GetStream() bytes.Buffer {
	return r.Data.(bytes.Buffer)
}
