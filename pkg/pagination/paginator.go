package pagination

import (
    "math"
    "net/http"
    "net/url"
    "reflect"
    "strconv"
    "strings"

    "store-api/pkg/data"

    "github.com/mitchellh/mapstructure"
)

const DefaultLimit = 10

type Paginator struct {
    r       *http.Request
    limit   int
    offset  int
    page    int
    perPage int
}

type Data struct {
    Data         interface{} `json:"data"`
    Path         string      `json:"path"`
    NextPage     string      `json:"next_page"`
    PreviousPage string      `json:"previous_page"`
    Total        int         `json:"total"`
    PerPage      int         `json:"per_page"`
    CurrentPage  int         `json:"current_page"`
}

type DataPagination struct {
    Data            interface{} `json:"data"`
    Path            string      `json:"path"`
    FirstPageUrl    string      `json:"first_page_url"`
    NextPageUrl     string      `json:"next_page_url,omitempty"`
    PreviousPageUrl string      `json:"previous_page_url,omitempty"`
    LastPageUrl     string      `json:"last_page_url"`
    From            int         `json:"from"`
    To              int         `json:"to"`
    LastPage        int         `json:"last_page"`
    Total           int         `json:"total"`
    PerPage         int         `json:"per_page"`
    CurrentPage     int         `json:"current_page"`
}

type Response struct {
    Items interface{}
    Total int
}

type WebAPIResponse struct {
    Items        interface{}
    Total        int
    FilterMember bool
}

type WebGetEmployeeListResponse struct {
    Data  interface{} `json:"data"`
    Links Links       `json:"_links"`
}

type WebTimesheetSummaryResponse struct {
    Data  interface{}
    Total int
    Links Links
}

type QueryPagination struct {
    Page    string `json:"page"`
    Limit   string `json:"limit"`
    PerPage string `json:"per-page"`
}

type Links struct {
    Self  string `json:"self"`
    First string `json:"first,omitempty"`
    Last  string `json:"last,omitempty"`
    Prev  string `json:"prev,omitempty"`
    Next  string `json:"next,omitempty"`
}

type DataPaginationAssigneeList struct {
    Path            string `json:"path"`
    FirstPageUrl    string `json:"first_page_url"`
    NextPageUrl     string `json:"next_page_url,omitempty"`
    PreviousPageUrl string `json:"previous_page_url,omitempty"`
    LastPageUrl     string `json:"last_page_url"`
    From            int    `json:"from"`
    To              int    `json:"to"`
    LastPage        int    `json:"last_page"`
    Total           int    `json:"total"`
    PerPage         int    `json:"per_page"`
    CurrentPage     int    `json:"current_page"`
}

type PaginationForAssigneeList struct {
    Data       interface{}                `json:"data"`
    Pagination DataPaginationAssigneeList `json:"pagination"`
}

func NewPaginator(r *http.Request) *Paginator {
    var offset = 0
    page, _ := strconv.Atoi(r.FormValue("page"))
    perPage, errPerpage := strconv.Atoi(r.FormValue("per_page"))

    if errPerpage != nil {
        perPage, _ = strconv.Atoi(r.FormValue("per-page"))
    }

    limit := perPage

    if page == 0 {
        page = 1
    }

    if limit == 0 {
        limit, _ = strconv.Atoi(r.FormValue("limit"))
    }

    if limit == 0 || limit > 500 {
        limit = DefaultLimit
    }

    if perPage == 0 && limit > 0 {
        perPage = limit
    }

    if page > 1 {
        offset = (page * limit) - limit
    }

    return &Paginator{r: r, page: page, limit: limit, offset: offset, perPage: perPage}
}

func (p Paginator) GetParams() data.Params {
    params := data.NewParamsWrapper()
    params.Add("offset", p.offset)
    params.Add("limit", p.limit)
    params.Add("page", p.page)
    params.Add("sso_id", p.r.Header.Get("X-SSO-User-ID")) // todo: Should remove sso_id, use ctx.GetSsoID() isntead
    return params
}

func (p Paginator) GetData(items interface{}) Data {
    r := p.r
    page := p.page
    requestURI := strings.Replace(r.RequestURI, "page="+r.FormValue("page"), "page=1", -1)

    nextPage := page + 1
    prevPage := page - 1

    if prevPage < 1 {
        prevPage = 1
    }

    if !strings.Contains(requestURI, "page") {
        if strings.Contains(requestURI, "?") {
            requestURI = requestURI + "&page=1"
        } else {
            requestURI = requestURI + "?page=1"
        }
    }

    nextURI := strings.Replace(requestURI, "page="+strconv.Itoa(page), "page="+strconv.Itoa(nextPage), -1)
    prevURI := strings.Replace(requestURI, "page="+strconv.Itoa(page), "page="+strconv.Itoa(prevPage), -1)

    scheme := "https://"
    if r.TLS == nil {
        scheme = "http://"
    }

    var total = 0
    switch reflect.TypeOf(items).Kind() {
    case reflect.Slice:
        s := reflect.ValueOf(items)
        total = s.Len()
    }

    return Data{
        Data:         items,
        Path:         scheme + r.Host + r.URL.Path,
        NextPage:     scheme + r.Host + nextURI,
        PreviousPage: scheme + r.Host + prevURI,
        Total:        total,
        PerPage:      p.limit,
        CurrentPage:  page,
    }
}

func (p Paginator) GetDataResponse(items interface{}) DataPagination {
    data := Response{}
    mapstructure.Decode(items, &data)

    r := p.r
    var firstURI, nextURI, prevURI, lastURI string
    requestURI := r.RequestURI
    page := p.page

    nextPage := page + 1
    prevPage := page - 1
    lastPage := int(math.Ceil(float64(data.Total) / float64(p.perPage)))

    if prevPage < 1 {
        prevPage = 1
    }

    if nextPage > lastPage {
        nextPage = lastPage
    }

    if !strings.Contains(requestURI, "page") {
        if strings.Contains(requestURI, "?") {
            requestURI = requestURI + "&page=1"
        } else {
            requestURI = requestURI + "?page=1"
        }
    }

    if strings.Contains(requestURI, "?page") {
        firstURI = strings.Replace(requestURI, "?page="+strconv.Itoa(page), "?page=1", -1)
        nextURI = strings.Replace(requestURI, "?page="+strconv.Itoa(page), "?page="+strconv.Itoa(nextPage), -1)
        prevURI = strings.Replace(requestURI, "?page="+strconv.Itoa(page), "?page="+strconv.Itoa(prevPage), -1)
        lastURI = strings.Replace(requestURI, "?page="+strconv.Itoa(page), "?page="+strconv.Itoa(lastPage), -1)
    } else if strings.Contains(requestURI, "&page") {
        firstURI = strings.Replace(requestURI, "&page="+strconv.Itoa(page), "&page=1", -1)
        nextURI = strings.Replace(requestURI, "&page="+strconv.Itoa(page), "&page="+strconv.Itoa(nextPage), -1)
        prevURI = strings.Replace(requestURI, "&page="+strconv.Itoa(page), "&page="+strconv.Itoa(prevPage), -1)
        lastURI = strings.Replace(requestURI, "&page="+strconv.Itoa(page), "&page="+strconv.Itoa(lastPage), -1)
    }

    scheme := "https://"
    if r.TLS == nil {
        scheme = "http://"
    }

    return DataPagination{
        Data:            data.Items,
        Path:            scheme + r.Host + r.URL.Path,
        FirstPageUrl:    scheme + r.Host + firstURI,
        NextPageUrl:     scheme + r.Host + nextURI,
        PreviousPageUrl: scheme + r.Host + prevURI,
        LastPageUrl:     scheme + r.Host + lastURI,
        From:            p.offset + 1,
        To:              p.limit * page,
        LastPage:        lastPage,
        Total:           data.Total,
        PerPage:         p.perPage,
        CurrentPage:     page,
    }
}

func (p Paginator) GetDataResponseAssigneList(items interface{}) PaginationForAssigneeList {
    data := Response{}
    mapstructure.Decode(items, &data)

    r := p.r
    var firstURI, nextURI, prevURI, lastURI string
    requestURI := r.RequestURI
    page := p.page

    nextPage := page + 1
    prevPage := page - 1
    lastPage := int(math.Ceil(float64(data.Total) / float64(p.perPage)))

    if prevPage < 1 {
        prevPage = 1
    }

    if nextPage > lastPage {
        nextPage = lastPage
    }

    if !strings.Contains(requestURI, "page") {
        if strings.Contains(requestURI, "?") {
            requestURI = requestURI + "&page=1"
        } else {
            requestURI = requestURI + "?page=1"
        }
    }

    if strings.Contains(requestURI, "?page") {
        firstURI = strings.Replace(requestURI, "?page="+strconv.Itoa(page), "?page=1", -1)
        nextURI = strings.Replace(requestURI, "?page="+strconv.Itoa(page), "?page="+strconv.Itoa(nextPage), -1)
        prevURI = strings.Replace(requestURI, "?page="+strconv.Itoa(page), "?page="+strconv.Itoa(prevPage), -1)
        lastURI = strings.Replace(requestURI, "?page="+strconv.Itoa(page), "?page="+strconv.Itoa(lastPage), -1)
    } else if strings.Contains(requestURI, "&page") {
        firstURI = strings.Replace(requestURI, "&page="+strconv.Itoa(page), "&page=1", -1)
        nextURI = strings.Replace(requestURI, "&page="+strconv.Itoa(page), "&page="+strconv.Itoa(nextPage), -1)
        prevURI = strings.Replace(requestURI, "&page="+strconv.Itoa(page), "&page="+strconv.Itoa(prevPage), -1)
        lastURI = strings.Replace(requestURI, "&page="+strconv.Itoa(page), "&page="+strconv.Itoa(lastPage), -1)
    }

    scheme := "https://"
    if r.TLS == nil {
        scheme = "http://"
    }

    dataPagination := DataPaginationAssigneeList{
        Path:            scheme + r.Host + r.URL.Path,
        FirstPageUrl:    scheme + r.Host + firstURI,
        NextPageUrl:     scheme + r.Host + nextURI,
        PreviousPageUrl: scheme + r.Host + prevURI,
        LastPageUrl:     scheme + r.Host + lastURI,
        From:            p.offset + 1,
        To:              p.limit * page,
        LastPage:        lastPage,
        Total:           data.Total,
        PerPage:         p.perPage,
        CurrentPage:     page,
    }

    return PaginationForAssigneeList{
        Data:       data.Items,
        Pagination: dataPagination,
    }
}

func (p Paginator) GetWebLinkPagination(items interface{}) Links {
    data := WebAPIResponse{}
    mapstructure.Decode(items, &data)

    r := p.r
    var currentURI, firstURI, nextURI, prevURI, lastURI string
    requestURI := r.RequestURI
    page := p.page
    currentPage := page

    nextPage := page + 1
    prevPage := page - 1
    lastPage := int(math.Ceil(float64(data.Total) / float64(p.perPage)))

    if prevPage < 1 {
        prevPage = 1
    }

    if nextPage > lastPage {
        nextPage = lastPage
    }

    if page >= lastPage {
        currentPage = lastPage
        prevPage = lastPage - 1
    }

    if page < 1 {
        currentPage = 1
        prevPage = 1
        if lastPage > currentPage {
            nextPage = currentPage + 1
        }

    }

    if !strings.Contains(requestURI, "page") {
        if strings.Contains(requestURI, "?") {
            requestURI = requestURI + "&page=1"
        } else {
            requestURI = requestURI + "?page=1"
        }
    }

    if !strings.Contains(requestURI, "&page") && !strings.Contains(requestURI, "?page") {
        if strings.Contains(requestURI, "?") {
            requestURI = requestURI + "&page=" + strconv.Itoa(page)
        } else {
            requestURI = requestURI + "?page=" + strconv.Itoa(page)
        }
    }

    if strings.Contains(requestURI, "?page") {
        firstURI = strings.Replace(requestURI, "?page="+strconv.Itoa(page), "?page=1", -1)
        nextURI = strings.Replace(requestURI, "?page="+strconv.Itoa(page), "?page="+strconv.Itoa(nextPage), -1)
        prevURI = strings.Replace(requestURI, "?page="+strconv.Itoa(page), "?page="+strconv.Itoa(prevPage), -1)
        lastURI = strings.Replace(requestURI, "?page="+strconv.Itoa(page), "?page="+strconv.Itoa(lastPage), -1)
        currentURI = strings.Replace(requestURI, "?page="+strconv.Itoa(page), "?page="+strconv.Itoa(currentPage), -1)

    } else if strings.Contains(requestURI, "&page") {
        firstURI = strings.Replace(requestURI, "&page="+strconv.Itoa(page), "&page=1", -1)
        nextURI = strings.Replace(requestURI, "&page="+strconv.Itoa(page), "&page="+strconv.Itoa(nextPage), -1)
        prevURI = strings.Replace(requestURI, "&page="+strconv.Itoa(page), "&page="+strconv.Itoa(prevPage), -1)
        lastURI = strings.Replace(requestURI, "&page="+strconv.Itoa(page), "&page="+strconv.Itoa(lastPage), -1)
        currentURI = strings.Replace(requestURI, "&page="+strconv.Itoa(page), "&page="+strconv.Itoa(currentPage), -1)
    }

    next := nextURI
    prev := prevURI
    last := lastURI
    first := firstURI

    if !strings.Contains(next, "per-page") {
        next = next + "&per-page=" + strconv.Itoa(p.limit)
    }

    if !strings.Contains(prev, "per-page") {
        prev = prev + "&per-page=" + strconv.Itoa(p.limit)
    }

    if !strings.Contains(last, "per-page") {
        last = last + "&per-page=" + strconv.Itoa(p.limit)
    }

    if !strings.Contains(first, "per-page") {
        first = first + "&per-page=" + strconv.Itoa(p.limit)
    }

    if !strings.Contains(currentURI, "per-page") {
        currentURI = currentURI + "&per-page=" + strconv.Itoa(p.limit)
    }

    if currentPage >= lastPage {
        next = ""
    }

    if currentPage <= 1 {
        prev = ""
    }

    itemsIsNil := false

    if data.Items == nil {
        itemsIsNil = true
    }

    if !itemsIsNil {
        switch reflect.TypeOf(data.Items).Kind() {
        case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
            itemsIsNil = reflect.ValueOf(data.Items).IsNil()
        }
    }

    if itemsIsNil {
        last = ""
        first = ""
    }

    return Links{
        Self:  currentURI,
        First: first,
        Last:  last,
        Prev:  prev,
        Next:  next,
    }
}

func (p Paginator) GetWebTimesheetSummaryPagination(items interface{}) Links {
    data := WebTimesheetSummaryResponse{}
    mapstructure.Decode(items, &data)

    r := p.r
    var currentURI, firstURI, nextURI, prevURI, lastURI string
    requestURI := r.RequestURI
    page := p.page
    currentPage := page

    nextPage := page + 1
    prevPage := page - 1
    lastPage := int(math.Ceil(float64(data.Total) / float64(p.perPage)))

    if prevPage < 1 {
        prevPage = 1
    }

    if nextPage > lastPage {
        nextPage = lastPage
    }

    if page >= lastPage {
        currentPage = lastPage
        prevPage = lastPage - 1
    }

    if page < 1 {
        currentPage = 1
        prevPage = 1
        if lastPage > currentPage {
            nextPage = currentPage + 1
        }

    }

    if !strings.Contains(requestURI, "page") {
        if strings.Contains(requestURI, "?") {
            requestURI = requestURI + "&page=1"
        } else {
            requestURI = requestURI + "?page=1"
        }
    }

    if !strings.Contains(requestURI, "&page") && !strings.Contains(requestURI, "?page") {
        if strings.Contains(requestURI, "?") {
            requestURI = requestURI + "&page=" + strconv.Itoa(page)
        } else {
            requestURI = requestURI + "?page=" + strconv.Itoa(page)
        }
    }

    if strings.Contains(requestURI, "?page") {
        firstURI = strings.Replace(requestURI, "?page="+strconv.Itoa(page), "?page=1", -1)
        nextURI = strings.Replace(requestURI, "?page="+strconv.Itoa(page), "?page="+strconv.Itoa(nextPage), -1)
        prevURI = strings.Replace(requestURI, "?page="+strconv.Itoa(page), "?page="+strconv.Itoa(prevPage), -1)
        lastURI = strings.Replace(requestURI, "?page="+strconv.Itoa(page), "?page="+strconv.Itoa(lastPage), -1)
        currentURI = strings.Replace(requestURI, "?page="+strconv.Itoa(page), "?page="+strconv.Itoa(currentPage), -1)

    } else if strings.Contains(requestURI, "&page") {
        firstURI = strings.Replace(requestURI, "&page="+strconv.Itoa(page), "&page=1", -1)
        nextURI = strings.Replace(requestURI, "&page="+strconv.Itoa(page), "&page="+strconv.Itoa(nextPage), -1)
        prevURI = strings.Replace(requestURI, "&page="+strconv.Itoa(page), "&page="+strconv.Itoa(prevPage), -1)
        lastURI = strings.Replace(requestURI, "&page="+strconv.Itoa(page), "&page="+strconv.Itoa(lastPage), -1)
        currentURI = strings.Replace(requestURI, "&page="+strconv.Itoa(page), "&page="+strconv.Itoa(currentPage), -1)
    }

    next := nextURI
    prev := prevURI
    last := lastURI
    first := firstURI

    if !strings.Contains(next, "per-page") {
        next = next + "&per-page=" + strconv.Itoa(p.limit)
    }

    if !strings.Contains(prev, "per-page") {
        prev = prev + "&per-page=" + strconv.Itoa(p.limit)
    }

    if !strings.Contains(last, "per-page") {
        last = last + "&per-page=" + strconv.Itoa(p.limit)
    }

    if !strings.Contains(first, "per-page") {
        first = first + "&per-page=" + strconv.Itoa(p.limit)
    }

    if !strings.Contains(currentURI, "per-page") {
        currentURI = currentURI + "&per-page=" + strconv.Itoa(p.limit)
    }

    if currentPage >= lastPage {
        next = ""
    }

    if currentPage <= 1 {
        prev = ""
    }

    itemsIsNil := false

    if data.Data == nil {
        itemsIsNil = true
    }

    if !itemsIsNil {
        switch reflect.TypeOf(data.Data).Kind() {
        case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
            itemsIsNil = reflect.ValueOf(data.Data).IsNil()
        }
    }

    if itemsIsNil {
        last = ""
        first = ""
    }

    return Links{
        Self:  currentURI,
        First: first,
        Last:  last,
        Prev:  prev,
        Next:  next,
    }
}

func (p Paginator) GetWebEmployeeListPagination(responseData interface{}) WebGetEmployeeListResponse {
    /******
    * format of inputData = {"self": "xxx", "first": "xxxx", "last": "xxxx", "prev": "xxxx", "next": "xxxx"}
    ******/
    dataLinks := Links{}
    data := responseData.(map[string]interface{})
    links := data["_links"].(map[string]interface{})
    items := data["data"]
    mapstructure.Decode(links, &dataLinks)

    r := p.r
    var requestURI, self, first, last, prev, next string
    requestURI = r.RequestURI
    self = requestURI + dataLinks.Self

    if strings.Contains(requestURI, "?") {
        requestURI = requestURI[:strings.Index(requestURI, "?")]
    }

    if strings.Contains(dataLinks.Self, "?") {
        self = requestURI + generateQueryParams(dataLinks.Self)
    }
    if strings.Contains(dataLinks.First, "?") {
        first = requestURI + generateQueryParams(dataLinks.First)
    }
    if strings.Contains(dataLinks.Last, "?") {
        last = requestURI + generateQueryParams(dataLinks.Last)
    }
    if strings.Contains(dataLinks.Prev, "?") {
        prev = requestURI + generateQueryParams(dataLinks.Prev)
    }

    if strings.Contains(dataLinks.Next, "?") {
        next = requestURI + generateQueryParams(dataLinks.Next)
    }

    resultLinks := Links{
        Self:  self,
        First: first,
        Last:  last,
        Prev:  prev,
        Next:  next,
    }

    return WebGetEmployeeListResponse{
        Data:  items,
        Links: resultLinks,
    }
}

// must output like this format: ?page=xx&limit=xx&per-page=xx
func generateQueryParams(urlPath string) string {
    requestURI, err := url.Parse(urlPath)
    if err != nil {
        return ""
    }

    q, err := url.ParseQuery(requestURI.RawQuery)
    if err != nil {
        return ""
    }

    paginationStruct := QueryPagination{
        Page:    q.Get("page"),
        Limit:   q.Get("limit"),
        PerPage: q.Get("per-page"),
    }

    result := "?page=" + paginationStruct.Page
    result += "&limit=" + paginationStruct.Limit
    result += "&per-page=" + paginationStruct.PerPage

    return result
}
