package httpclient

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"reflect"
	"time"

	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
)

//New creates client Factory
func New() ClientFactory {
	return *&clientFactory{}
}

//ClientFactory creates specific client implementation
type ClientFactory interface {
	CreateClient() Client
}

type clientFactory struct{}

func (c clientFactory) CreateClient() Client {
	return client{}
}

//Client abstracts third party request client
type Client interface {
	GetForUsersFilteredByBranchAndOrganizationAndName(string, map[string]string, interface{}) (int, error)
	Get(string, map[string]string, interface{}) (int, error)
	PostJSON(string, string, map[string]string, interface{}) (int, error)
	Post(string, interface{}, map[string]string, interface{}) (int, error)
	PostForm(string, interface{}, map[string]string, interface{}) (int, error)
	PostEncodedForm(string, interface{}, map[string]string, interface{}) (int, error)
	Put(string, interface{}, map[string]string, interface{}) (int, error)
	PostMultipart(string, interface{}, map[string]string, interface{}, interface{}) (int, error)
	PostMultiparts(string, interface{}, map[string]string, interface{}, []string) (int, error)
}

type client struct {
}

func (g client) GetForUsersFilteredByBranchAndOrganizationAndName(url string, headers map[string]string, dest interface{}) (int, error) {
	getRequest := gorequest.New()
	getRequest.Timeout(1 * time.Minute)
	agent := getRequest.Get(url)
	for k, v := range headers {
		agent = agent.Set(k, v)
	}
	resp, bytes, errs := agent.EndBytes()
	if len(errs) > 0 {
		return http.StatusInternalServerError, errs[0]
	}
	if err := json.Unmarshal(bytes, &dest); err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "error while reading response content")
	}
	//fmt.Println(string(bytes))
	return resp.StatusCode, nil
}

func (g client) Get(url string, headers map[string]string, dest interface{}) (int, error) {
	getRequest := gorequest.New()
	agent := getRequest.Get(url)
	for k, v := range headers {
		agent = agent.Set(k, v)
	}
	resp, bytes, errs := agent.EndBytes()
	if len(errs) > 0 {
		return http.StatusInternalServerError, errs[0]
	}
	if err := json.Unmarshal(bytes, &dest); err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "error while reading response content")
	}
	//fmt.Println(string(bytes))
	return resp.StatusCode, nil
}

func (g client) PostJSON(string, string, map[string]string, interface{}) (int, error) {
	return 0, nil
}

// url encoded sample https://golang.cafe/blog/how-to-make-http-url-form-encoded-request-golang.html
func (g client) PostEncodedForm(url string, data interface{}, headers map[string]string, dest interface{}) (int, error) {
	postRequest := gorequest.New()
	agent := postRequest.Put(url)
	for k, v := range headers {
		agent = agent.Set(k, v)
	}
	resp, bytes, errs := agent.Type("urlencoded").Send(data).EndBytes()
	if len(errs) > 0 {
		return http.StatusInternalServerError, errs[0]
	}
	if err := json.Unmarshal(bytes, &dest); err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "error while reading response content")
	}
	return resp.StatusCode, nil
}

// Post JSON
func (g client) Post(url string, data interface{}, headers map[string]string, dest interface{}) (int, error) {
	postRequest := gorequest.New()
	agent := postRequest.Post(url)
	for k, v := range headers {
		agent = agent.Set(k, v)
	}
	resp, bytes, errs := agent.Send(data).EndBytes()
	if len(errs) > 0 {
		return http.StatusInternalServerError, errs[0]
	}
	if err := json.Unmarshal(bytes, &dest); err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "error while reading response content")
	}
	return resp.StatusCode, nil
}

// Put JSON
func (g client) Put(url string, data interface{}, headers map[string]string, dest interface{}) (int, error) {
	postRequest := gorequest.New()
	agent := postRequest.Put(url)
	for k, v := range headers {
		agent = agent.Set(k, v)
	}

	resp, bytes, errs := agent.Send(data).EndBytes()
	if len(errs) > 0 {
		return http.StatusInternalServerError, errs[0]
	}
	if err := json.Unmarshal(bytes, &dest); err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "error while reading response content")
	}
	return resp.StatusCode, nil
}

func (g client) PostForm(url string, data interface{}, headers map[string]string, dest interface{}) (int, error) {
	postRequest := gorequest.New()
	agent := postRequest.Post(url)
	for k, v := range headers {
		agent = agent.Set(k, v)
	}
	resp, bytes, errs := agent.Type("form-data").Send(data).EndBytes()

	if len(errs) > 0 {
		return http.StatusInternalServerError, errs[0]
	}

	if err := json.Unmarshal(bytes, &dest); err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "error while reading response content")
	}

	return resp.StatusCode, nil
}

func (g client) PostMultipart(url string, data interface{}, headers map[string]string, dest interface{}, file interface{}) (int, error) {
	postRequest := gorequest.New()
	agent := postRequest.Post(url)
	for k, v := range headers {
		agent = agent.Set(k, v)
	}
	if file != nil && reflect.ValueOf(file).IsNil() == false {
		agent = agent.SendFile(file, "", "upload")
	}

	resp, bytes, errs := agent.Type("multipart").Send(data).EndBytes()

	if len(errs) > 0 {
		return http.StatusInternalServerError, errs[0]
	}

	if err := json.Unmarshal(bytes, &dest); err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "error while reading response content")
	}

	return resp.StatusCode, nil
}

func (g client) PostMultiparts(url string, data interface{}, headers map[string]string, dest interface{}, file []string) (int, error) {
	postRequest := gorequest.New()
	agent := postRequest.Post(url)
	for k, v := range headers {
		agent = agent.Set(k, v)
	}

	for _, v := range file {
		f, _ := filepath.Abs(v)
		bytesOfFile, _ := ioutil.ReadFile(f)
		agent = agent.SendFile(bytesOfFile, v, "file_name[]")
	}

	resp, bytes, errs := agent.Type("multipart").Send(data).EndBytes()

	if len(errs) > 0 {
		return http.StatusInternalServerError, errs[0]
	}

	if err := json.Unmarshal(bytes, &dest); err != nil {
		return http.StatusInternalServerError, errors.Wrap(err, "error while reading response content")
	}

	return resp.StatusCode, nil
}
