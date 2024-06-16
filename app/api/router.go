package api

import (
    "fmt"
    "net/http"

    "store-api/internal/base/handler"
)

func (h *HttpServe) setupRouter() {
    // StrictSlash will treat /projects/ to be same as /projects
    h.v1 = h.router.PathPrefix("/api/v1/").Subrouter()

    h.Route("POST", "/product/list", h.store.ListProduct)
    h.Route("POST", "/cart/add", h.store.AddToCart)
    h.Route("POST", "/cart/view", h.store.ViewCart)
    h.Route("POST", "/cart/delete", h.store.DeleteProductInCart)
    h.Route("POST", "/transaction/create", h.store.CreateTransaction)
    h.Route("POST", "/login", h.store.Login)

    // assign method not allowed handler
    h.v1.MethodNotAllowedHandler = h.base.MethodNotAllowedHandler()
}

func (h *HttpServe) Route(method string, path string, f handler.HandlerFn) {
    if method != http.MethodGet &&
            method != http.MethodPost &&
            method != http.MethodDelete &&
            method != http.MethodPut {
        panic(fmt.Sprintf(":%s method not allow", method))
    }

    h.v1.HandleFunc(path, h.base.RunAction(f)).Methods(method)
}
