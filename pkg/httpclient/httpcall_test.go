package httpclient

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Pokemon struct {
	Name   string `json:"name"`
	Weight int    `json:"weight"`
	Height int    `json:"height"`
}

type GetUsersResponse struct {
	Message string `json:"message"`
}

type Response struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func mockGetForUsersFilteredByBranchAndorg(w http.ResponseWriter, r *http.Request) {
	sc := http.StatusOK
	m := make(map[string]string)
	m["message"] = "hello"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(sc)
	json.NewEncoder(w).Encode(m)
}

func mockPostEncodedForm(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	m := make(map[string]interface{})

	m["status"] = r.Form.Get("status")
	m["message"] = r.Form.Get("message")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(m)
}

func mockProcessFormData(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var data Response
	err := decoder.Decode(&data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	m := make(map[string]interface{})

	m["status"] = data.Status
	m["message"] = data.Message

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(m)

}

func mockProcessMultipartFormData(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20) // maxMemory 32MB
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	status := r.Form.Get("status")
	message := r.Form.Get("message")

	m := make(map[string]interface{})

	m["status"] = status
	m["message"] = message

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(m)

}

func TestHTTPCall(t *testing.T) {
	httpCall := New()
	client := httpCall.CreateClient()

	var resp *Pokemon
	code, err := client.Get("https://pokeapi.co/api/v2/pokemon/1", nil, &resp)
	if err != nil {
		t.Error("http call error", err.Error())
	}
	assert.Equal(t, http.StatusOK, code, "code equals to 200")
	assert.NotNil(t, resp, "object not nil")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch strings.TrimSpace(r.URL.Path) {
		case "/GetForUsersFilteredByBranchAndOrganizationAndName":
			mockGetForUsersFilteredByBranchAndorg(w, r)
		case "/PostEncodedForm":
			mockPostEncodedForm(w, r)
		case "/PostJSON":
			mockProcessFormData(w, r)
		case "/PutJSON":
			mockProcessFormData(w, r)
		case "/postForm":
			mockProcessFormData(w, r)
		case "/postMultiPart":
			mockProcessMultipartFormData(w, r)
		default:
			http.NotFoundHandler().ServeHTTP(w, r)
		}
	}))

	defer server.Close()

	t.Run("GetForUsersFilteredByBranchAndOrganizationAndName", func(t *testing.T) {

		var resp *GetUsersResponse

		status, err := client.GetForUsersFilteredByBranchAndOrganizationAndName(
			server.URL+"/GetForUsersFilteredByBranchAndOrganizationAndName", nil, &resp,
		)

		assert.Equal(t, status, http.StatusOK)
		assert.NotNil(t, resp)
		assert.NoError(t, err)

	})

	t.Run("PostEncodedForm", func(t *testing.T) {

		var resp *Response

		data := url.Values{}
		data.Set("status", "200")
		data.Set("message", "success")

		headers := map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		}

		status, err := client.PostEncodedForm(
			server.URL+"/PostEncodedForm", strings.NewReader(data.Encode()), headers, &resp,
		)

		assert.Equal(t, status, http.StatusOK)
		assert.NoError(t, err)

	})

	t.Run("PostJSON", func(t *testing.T) {

		var resp *Response

		data := map[string]string{"status": "200", "message": "PostJSon"}

		headers := map[string]string{
			"X-SSO-ID": "12345-6789",
		}

		status, err := client.Post(
			server.URL+"/PostJSON", data, headers, &resp,
		)

		assert.Equal(t, status, http.StatusOK)
		assert.NoError(t, err)
		assert.Equal(t, data["status"], resp.Status)
		assert.Equal(t, data["message"], resp.Message)

	})

	t.Run("PutJSON", func(t *testing.T) {

		var resp *Response

		data := map[string]string{"status": "200", "message": "PostJSon"}

		headers := map[string]string{
			"Content-Type": "application/json",
		}

		status, err := client.Put(
			server.URL+"/PutJSON", data, headers, &resp,
		)

		assert.Equal(t, status, http.StatusOK)
		assert.NoError(t, err)
		assert.Equal(t, data["status"], resp.Status)
		assert.Equal(t, data["message"], resp.Message)

		// failed to unmarshal the result to dest
		status, err = client.Put(
			server.URL+"/PutJSON", map[string]interface{}{"status": true}, headers, resp,
		)

		assert.Error(t, err)
		assert.Equal(t, status, http.StatusInternalServerError)

	})

	t.Run("Post Form", func(t *testing.T) {

		var resp *Response

		headers := map[string]string{
			"X-SSO-ID": "12345-6789",
		}

		data := map[string]string{"status": "200", "message": "testPostForm"}

		status, err := client.PostForm(
			server.URL+"/postForm", data, headers, &resp,
		)

		assert.Equal(t, status, http.StatusOK)
		assert.NoError(t, err)
		assert.Equal(t, resp.Status, data["status"])
		assert.Equal(t, resp.Message, data["message"])

		status, err = client.PostForm(
			server.URL+"/postForm", map[string]interface{}{"status": true}, headers, &resp,
		)

		assert.Error(t, err)
		assert.Equal(t, status, http.StatusInternalServerError)

	})

	t.Run("Post Multipart", func(t *testing.T) {

		var resp *Response

		headers := map[string]string{
			"X-SSO-ID": "12345-6789",
		}

		data := map[string]string{
			"status":  "200",
			"message": "PostMultipart",
		}
		status, err := client.PostMultipart(
			server.URL+"/postMultiPart", data, headers, &resp, nil,
		)

		assert.Equal(t, status, http.StatusOK)
		assert.NoError(t, err)
		assert.Equal(t, data["status"], resp.Status)
		assert.Equal(t, data["message"], resp.Message)

	})

	t.Run("Post Multipart with file", func(t *testing.T) {

		var resp *Response

		headers := map[string]string{
			"X-SSO-ID": "12345-6789",
		}

		data := map[string]string{
			"status":  "200",
			"message": "PostMultipart",
		}
		status, err := client.PostMultipart(
			server.URL+"/postMultiPart", data, headers, &resp, []byte("hello\ngo\n"),
		)

		assert.Equal(t, status, http.StatusOK)
		assert.NoError(t, err)
		assert.Equal(t, data["status"], resp.Status)
		assert.Equal(t, data["message"], resp.Message)

	})

	t.Run("Post Multiparts", func(t *testing.T) {

		// create temporary file
		f, err := os.CreateTemp("", "tmpfile-")
		if err != nil {
			t.Error(err)
		}

		// close and remove the temporary file
		defer f.Close()
		defer os.Remove(f.Name())

		fileData := []byte("halo from go")
		if _, err := f.Write(fileData); err != nil {
			t.Error(err)
		}

		file := make([]string, 0)
		file = append(file, f.Name())

		var resp *Response
		data := map[string]string{
			"status":  "200",
			"message": "PostMultipart",
		}
		headers := map[string]string{
			"X-SSO-ID": "12345-6789",
		}
		status, err := client.PostMultiparts(
			server.URL+"/postMultiPart", data, headers, &resp, file,
		)

		assert.Equal(t, status, http.StatusOK)
		assert.NoError(t, err)
		assert.Equal(t, data["status"], resp.Status)
		assert.Equal(t, data["message"], resp.Message)

	})

	t.Run("Post JSON", func(t *testing.T) {

		var resp *Response

		data := map[string]string{"status": "200", "message": "testPostForm"}

		status, err := client.PostJSON(
			server.URL+"/postForm", "test", data, resp,
		)

		assert.Equal(t, status, 0)
		assert.NoError(t, err)

	})
}
