package pagination

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPaginator(t *testing.T) {
	url := "test.com"

	t.Run("Success return Paginator struct with default settings", func(t *testing.T) {
		req, _ := http.NewRequest("GET", url, nil)
		paginator := NewPaginator(req)

		assert.Equal(t, url, paginator.r.URL.Path)
		assert.Equal(t, 0, paginator.offset)
		assert.Equal(t, 1, paginator.page)
		assert.Equal(t, 10, paginator.perPage)
		assert.Equal(t, 10, paginator.limit)
	})

	t.Run("Success return Paginator with specific per page", func(t *testing.T) {
		url := "test.com?per_page=500"
		req, _ := http.NewRequest("GET", url, nil)
		paginator := NewPaginator(req)

		assert.Equal(t, 500, paginator.limit)
		assert.Equal(t, 500, paginator.perPage)
	})

	t.Run("Read per-page query param if per_page is empty", func(t *testing.T) {
		url := "test.com?per-page=500"
		req, _ := http.NewRequest("GET", url, nil)
		paginator := NewPaginator(req)

		assert.Equal(t, 500, paginator.limit)
		assert.Equal(t, 500, paginator.perPage)
	})

	t.Run("Success return Paginator with page=1 and offset=1", func(t *testing.T) {
		url := "test.com?page=1"
		req, _ := http.NewRequest("GET", url, nil)
		paginator := NewPaginator(req)

		assert.Equal(t, 10, paginator.limit)
		assert.Equal(t, 10, paginator.perPage)
		assert.Equal(t, 1, paginator.page)
		assert.Equal(t, 0, paginator.offset)
	})

	t.Run("Success return Paginator with page=99 and offset 980", func(t *testing.T) {
		url := "test.com?page=99"
		req, _ := http.NewRequest("GET", url, nil)
		paginator := NewPaginator(req)

		assert.Equal(t, 10, paginator.limit)
		assert.Equal(t, 10, paginator.perPage)
		assert.Equal(t, 99, paginator.page)
		assert.Equal(t, 980, paginator.offset)
	})
}

func TestGetParams(t *testing.T) {
	url := "test.com"
	mySSOID := "asdf"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("X-SSO-User-ID", mySSOID)
	paginator := NewPaginator(req)
	myLimit := 999
	myPage := 2
	myOffset := 888

	t.Run("Success return Params with specified fields", func(t *testing.T) {
		paginator.limit = myLimit
		paginator.page = myPage
		paginator.offset = myOffset

		params := paginator.GetParams()
		assert.Equal(t, myLimit, params.GetInt("limit"))
		assert.Equal(t, myPage, params.GetInt("page"))
		assert.Equal(t, myOffset, params.GetInt("offset"))
		assert.Equal(t, mySSOID, params.GetString("sso_id"))
	})
}

func TestGetData(t *testing.T) {
	url := "test.com"
	mySSOID := "asdf"

	req, _ := http.NewRequest("GET", url, nil)
	req.RequestURI = url
	req.Header.Add("X-SSO-User-ID", mySSOID)

	t.Run("Success GetData", func(t *testing.T) {
		paginator := NewPaginator(req)
		items := []int{1, 2}
		data := paginator.GetData(items)

		nextURI := "http://" + url + "?page=2"
		prevURI := "http://" + url + "?page=1"

		assert.Equal(t, items, data.Data)
		assert.Equal(t, nextURI, data.NextPage)
		assert.Equal(t, prevURI, data.PreviousPage)
	})
}

func TestGetWebEmployeeListPagination(t *testing.T) {
	url := "test.com"
	items := []int{1, 2, 3}
	links := map[string]interface{}{
		"self":  "",
		"first": "",
		"last":  "",
		"prev":  "",
		"next":  "",
	}
	responseData := map[string]interface{}{
		"data":   items,
		"_links": links,
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.RequestURI = url

	t.Run("Success GetWebEmployeeListPagination", func(t *testing.T) {
		paginator := NewPaginator(req)
		paginator.GetWebEmployeeListPagination(responseData)
	})

	t.Run("Success GetWebEmployeeListPagination with pagination query params", func(t *testing.T) {
		links := map[string]interface{}{
			"self":  "test.com?page=2&limit=5&per-page=5",
			"first": "test.com?page=1&limit=5&per-page=5",
			"last":  "test.com?page=3&limit=5&per-page=5",
			"prev":  "test.com?page=1&limit=5&per-page=5",
			"next":  "test.com?page=3&limit=5&per-page=5",
		}
		responseData := map[string]interface{}{
			"data":   items,
			"_links": links,
		}
		req, _ := http.NewRequest("GET", url, nil)
		req.RequestURI = url + "?"
		paginator := NewPaginator(req)
		paginator.GetWebEmployeeListPagination(responseData)
	})
}

func TestGetDataResponse(t *testing.T) {
	url := "test.com"
	mySSOID := "asdf"

	req, _ := http.NewRequest("GET", url, nil)
	req.RequestURI = url
	req.Header.Add("X-SSO-User-ID", mySSOID)

	t.Run("Success GetDataResponse", func(t *testing.T) {
		paginator := NewPaginator(req)
		items := []int{1, 2}
		response := map[string]interface{}{
			"items": items,
			"total": 2,
		}
		data := paginator.GetDataResponse(response)

		nextURI := "http://" + url + "?page=1"

		assert.Equal(t, items, data.Data)
		assert.Equal(t, nextURI, data.NextPageUrl)
		assert.Equal(t, nextURI, data.PreviousPageUrl)
	})

	t.Run("Success GetDataResponse page param order", func(t *testing.T) {
		url := "test.com"
		form := "?per-page=10&page=1"
		finalURI := url + form
		req, _ := http.NewRequest("GET", finalURI, nil)
		req.RequestURI = finalURI

		paginator := NewPaginator(req)
		items := []int{1, 2}
		response := map[string]interface{}{
			"items": items,
			"total": 2,
		}
		data := paginator.GetDataResponse(response)
		nextURI := "http://" + url + form

		assert.Equal(t, items, data.Data)
		assert.Equal(t, nextURI, data.NextPageUrl)
		assert.Equal(t, nextURI, data.PreviousPageUrl)
	})
}

func TestGetDataResponseAssigneList(t *testing.T) {
	url := "test.com"
	mySSOID := "asdf"

	req, _ := http.NewRequest("GET", url, nil)
	req.RequestURI = url
	req.Header.Add("X-SSO-User-ID", mySSOID)

	t.Run("Success GetDataResponseAssigneList", func(t *testing.T) {
		paginator := NewPaginator(req)
		items := []int{1, 2}
		response := map[string]interface{}{
			"items": items,
			"total": 2,
		}
		data := paginator.GetDataResponseAssigneList(response)

		nextURI := "http://" + url + "?page=1"

		assert.Equal(t, items, data.Data)
		assert.Equal(t, nextURI, data.Pagination.NextPageUrl)
		assert.Equal(t, nextURI, data.Pagination.PreviousPageUrl)
	})

	t.Run("Success GetDataResponseAssigneList", func(t *testing.T) {
		form := "?per-page=10&page=1"
		finalURI := url + form
		req, _ := http.NewRequest("GET", finalURI, nil)
		req.RequestURI = finalURI

		paginator := NewPaginator(req)
		items := []int{1, 2}
		response := map[string]interface{}{
			"items": items,
			"total": 2,
		}
		data := paginator.GetDataResponseAssigneList(response)

		nextURI := "http://" + url + form

		assert.Equal(t, items, data.Data)
		assert.Equal(t, nextURI, data.Pagination.NextPageUrl)
		assert.Equal(t, nextURI, data.Pagination.PreviousPageUrl)
	})
}

func TestGetWebLinkPagination(t *testing.T) {
	url := "test.com"
	mySSOID := "asdf"

	req, _ := http.NewRequest("GET", url, nil)
	req.RequestURI = url
	req.Header.Add("X-SSO-User-ID", mySSOID)

	t.Run("Success GetWebLinkPagination with only 1 page", func(t *testing.T) {
		paginator := NewPaginator(req)
		items := []int{1, 2}
		response := map[string]interface{}{
			"items": items,
			"total": 2,
		}
		data := paginator.GetWebLinkPagination(response)

		finalURI := url + "?page=1&per-page=10"

		assert.Equal(t, finalURI, data.Self)
		assert.Equal(t, "", data.Next)
		assert.Equal(t, "", data.Prev)
		assert.Equal(t, finalURI, data.First)
		assert.Equal(t, finalURI, data.Last)
	})

	t.Run("Success GetWebLinkPagination with more than 1 page and current page = 1", func(t *testing.T) {
		form := "?page=1&per-page=1"
		finalUrl := url + form
		req, _ := http.NewRequest("GET", finalUrl, nil)
		req.RequestURI = finalUrl

		paginator := NewPaginator(req)
		items := []int{1, 2, 3}
		response := map[string]interface{}{
			"items": items,
			"total": 3,
		}
		data := paginator.GetWebLinkPagination(response)

		selfURI := url + form
		nextURI := url + "?page=2&per-page=1"
		lastURI := url + "?page=3&per-page=1"

		assert.Equal(t, selfURI, data.Self)
		assert.Equal(t, nextURI, data.Next)
		assert.Equal(t, "", data.Prev)
		assert.Equal(t, selfURI, data.First)
		assert.Equal(t, lastURI, data.Last)
	})

	t.Run("Success GetWebLinkPagination with more than 1 page and current page = 2", func(t *testing.T) {
		form := "?page=2&per-page=1"
		finalUrl := url + form
		req, _ := http.NewRequest("GET", finalUrl, nil)
		req.RequestURI = finalUrl

		paginator := NewPaginator(req)
		items := []int{1, 2, 3}
		response := map[string]interface{}{
			"items": items,
			"total": 3,
		}
		data := paginator.GetWebLinkPagination(response)

		selfURI := url + form
		nextURI := url + "?page=3&per-page=1"
		prevURI := url + "?page=1&per-page=1"

		assert.Equal(t, selfURI, data.Self)
		assert.Equal(t, nextURI, data.Next)
		assert.Equal(t, prevURI, data.Prev)
		assert.Equal(t, prevURI, data.First)
		assert.Equal(t, nextURI, data.Last)
	})

	t.Run("Success GetWebLinkPagination with page param order", func(t *testing.T) {
		form := "?per-page=1&page=2"
		finalUrl := url + form
		req, _ := http.NewRequest("GET", finalUrl, nil)
		req.RequestURI = finalUrl
		paginator := NewPaginator(req)
		items := []int{1, 2, 3}
		response := map[string]interface{}{
			"items": items,
			"total": 3,
		}
		data := paginator.GetWebLinkPagination(response)

		self := url + "?per-page=1&page=2"
		nextURI := url + "?per-page=1&page=3"
		prevURI := url + "?per-page=1&page=1"

		assert.Equal(t, self, data.Self)
		assert.Equal(t, nextURI, data.Next)
		assert.Equal(t, nextURI, data.Last)
		assert.Equal(t, prevURI, data.First)
		assert.Equal(t, prevURI, data.Prev)
	})
}

func TestGetWebTimesheetSummaryPagination(t *testing.T) {
	url := "test.com"
	form := "?page=2&per-page=1"
	finalUrl := url + form

	req, _ := http.NewRequest("GET", finalUrl, nil)
	req.RequestURI = finalUrl

	t.Run("Success GetWebTimesheetSummaryPagination", func(t *testing.T) {
		paginator := NewPaginator(req)
		items := []int{1, 2, 3}
		response := map[string]interface{}{
			"data":  items,
			"total": 3,
		}
		data := paginator.GetWebTimesheetSummaryPagination(response)

		self := url + "?page=2&per-page=1"
		nextURI := url + "?page=3&per-page=1"
		prevURI := url + "?page=1&per-page=1"

		assert.Equal(t, self, data.Self)
		assert.Equal(t, nextURI, data.Next)
		assert.Equal(t, nextURI, data.Last)
		assert.Equal(t, prevURI, data.First)
		assert.Equal(t, prevURI, data.Prev)
	})

	t.Run("Success GetWebTimesheetSummaryPagination with more than 1 page and current page = 2", func(t *testing.T) {
		form := "?page=2&per-page=1"
		finalUrl := url + form
		req, _ := http.NewRequest("GET", finalUrl, nil)
		req.RequestURI = finalUrl
		paginator := NewPaginator(req)
		items := []int{1, 2, 3}
		response := map[string]interface{}{
			"data":  items,
			"total": 3,
		}
		data := paginator.GetWebTimesheetSummaryPagination(response)

		self := url + "?page=2&per-page=1"
		nextURI := url + "?page=3&per-page=1"
		prevURI := url + "?page=1&per-page=1"

		assert.Equal(t, self, data.Self)
		assert.Equal(t, nextURI, data.Next)
		assert.Equal(t, nextURI, data.Last)
		assert.Equal(t, prevURI, data.First)
		assert.Equal(t, prevURI, data.Prev)
	})

	t.Run("Success GetWebTimesheetSummaryPagination with page param order", func(t *testing.T) {
		form := "?per-page=1&page=2"
		finalUrl := url + form
		req, _ := http.NewRequest("GET", finalUrl, nil)
		req.RequestURI = finalUrl
		paginator := NewPaginator(req)
		items := []int{1, 2, 3}
		response := map[string]interface{}{
			"data":  items,
			"total": 3,
		}
		data := paginator.GetWebTimesheetSummaryPagination(response)

		self := url + "?per-page=1&page=2"
		nextURI := url + "?per-page=1&page=3"
		prevURI := url + "?per-page=1&page=1"

		assert.Equal(t, self, data.Self)
		assert.Equal(t, nextURI, data.Next)
		assert.Equal(t, nextURI, data.Last)
		assert.Equal(t, prevURI, data.First)
		assert.Equal(t, prevURI, data.Prev)
	})
}
