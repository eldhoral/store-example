package handler

import (
    "net/http"

    "store-api/internal/base/app"
    presenterCart "store-api/internal/store/presenter/cart"
    presenterMember "store-api/internal/store/presenter/member"
    presenterProduct "store-api/internal/store/presenter/product"
    presenterTransaction "store-api/internal/store/presenter/transaction"
    "store-api/pkg/data/constant"
    "store-api/pkg/server"

    jsoniter "github.com/json-iterator/go"
)

func (h HTTPHandler) ListProduct(ctx *app.Context) *server.Response {
    ctx.ParseJson()
    isJson := ctx.IsContentTypeJson()
    if !isJson {
        return h.AsWebResponse(ctx, http.StatusBadRequest, "invalid content type", constant.EmptyArray)
    }

    jsonBody := ctx.GetJsonBody()
    if jsonBody == nil {
        return h.AsWebResponse(ctx, http.StatusInternalServerError, "Json Body is required", constant.EmptyArray)
    }

    convertToJsonString, err := jsoniter.Marshal(jsonBody)
    if err != nil {
        return h.AsWebResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
    }

    productReq := presenterProduct.ProductRequest{}
    jsoniter.Unmarshal(convertToJsonString, &productReq)

    result, httpStatus, err := h.StoreService.ListProduct(productReq)
    if err != nil {
        return h.AsWebResponse(ctx, httpStatus, err.Error(), nil)
    }

    return h.AsMobileJson(ctx, httpStatus, "List Product Success", result)
}

func (h HTTPHandler) AddToCart(ctx *app.Context) *server.Response {
    ctx.ParseJson()
    isJson := ctx.IsContentTypeJson()
    if !isJson {
        return h.AsWebResponse(ctx, http.StatusBadRequest, "invalid content type", constant.EmptyArray)
    }

    jsonBody := ctx.GetJsonBody()
    if jsonBody == nil {
        return h.AsWebResponse(ctx, http.StatusInternalServerError, "Json Body is required", constant.EmptyArray)
    }

    convertToJsonString, err := jsoniter.Marshal(jsonBody)
    if err != nil {
        return h.AsWebResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
    }

    cartReq := presenterCart.CartRequest{}
    jsoniter.Unmarshal(convertToJsonString, &cartReq)

    httpStatus, err := h.StoreService.AddToCart(cartReq)
    if err != nil {
        return h.AsWebResponse(ctx, httpStatus, err.Error(), nil)
    }

    return h.AsMobileJson(ctx, httpStatus, "Add To Cart Success", nil)
}

func (h HTTPHandler) ViewCart(ctx *app.Context) *server.Response {
    ctx.ParseJson()
    isJson := ctx.IsContentTypeJson()
    if !isJson {
        return h.AsWebResponse(ctx, http.StatusBadRequest, "invalid content type", constant.EmptyArray)
    }

    jsonBody := ctx.GetJsonBody()
    if jsonBody == nil {
        return h.AsWebResponse(ctx, http.StatusInternalServerError, "Json Body is required", constant.EmptyArray)
    }

    convertToJsonString, err := jsoniter.Marshal(jsonBody)
    if err != nil {
        return h.AsWebResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
    }

    cartReq := presenterCart.CartViewRequest{}
    jsoniter.Unmarshal(convertToJsonString, &cartReq)

    result, httpStatus, err := h.StoreService.ViewCart(cartReq)
    if err != nil {
        return h.AsWebResponse(ctx, httpStatus, err.Error(), nil)
    }

    return h.AsMobileJson(ctx, httpStatus, "View Cart Success", result)
}

func (h HTTPHandler) DeleteProductInCart(ctx *app.Context) *server.Response {
    ctx.ParseJson()
    isJson := ctx.IsContentTypeJson()
    if !isJson {
        return h.AsWebResponse(ctx, http.StatusBadRequest, "invalid content type", constant.EmptyArray)
    }

    jsonBody := ctx.GetJsonBody()
    if jsonBody == nil {
        return h.AsWebResponse(ctx, http.StatusInternalServerError, "Json Body is required", constant.EmptyArray)
    }

    convertToJsonString, err := jsoniter.Marshal(jsonBody)
    if err != nil {
        return h.AsWebResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
    }

    cartReq := presenterCart.CartProductDeleteRequest{}
    jsoniter.Unmarshal(convertToJsonString, &cartReq)

    httpStatus, err := h.StoreService.DeleteProductInCart(cartReq)
    if err != nil {
        return h.AsWebResponse(ctx, httpStatus, err.Error(), nil)
    }

    return h.AsMobileJson(ctx, httpStatus, "Delete Cart Success", nil)
}

func (h HTTPHandler) CreateTransaction(ctx *app.Context) *server.Response {
    ctx.ParseJson()
    isJson := ctx.IsContentTypeJson()
    if !isJson {
        return h.AsWebResponse(ctx, http.StatusBadRequest, "invalid content type", constant.EmptyArray)
    }

    jsonBody := ctx.GetJsonBody()
    if jsonBody == nil {
        return h.AsWebResponse(ctx, http.StatusInternalServerError, "Json Body is required", constant.EmptyArray)
    }

    convertToJsonString, err := jsoniter.Marshal(jsonBody)
    if err != nil {
        return h.AsWebResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
    }

    transactionReq := presenterTransaction.TransactionRequest{}
    jsoniter.Unmarshal(convertToJsonString, &transactionReq)

    httpStatus, err := h.StoreService.CreateTransaction(transactionReq)
    if err != nil {
        return h.AsWebResponse(ctx, httpStatus, err.Error(), nil)
    }

    return h.AsMobileJson(ctx, httpStatus, "Transaction Success", nil)
}

func (h HTTPHandler) Login(ctx *app.Context) *server.Response {
    ctx.ParseJson()
    isJson := ctx.IsContentTypeJson()
    if !isJson {
        return h.AsWebResponse(ctx, http.StatusBadRequest, "invalid content type", constant.EmptyArray)
    }

    jsonBody := ctx.GetJsonBody()
    if jsonBody == nil {
        return h.AsWebResponse(ctx, http.StatusInternalServerError, "Json Body is required", constant.EmptyArray)
    }

    convertToJsonString, err := jsoniter.Marshal(jsonBody)
    if err != nil {
        return h.AsWebResponse(ctx, http.StatusInternalServerError, err.Error(), nil)
    }

    memberReq := presenterMember.LoginRequest{}
    jsoniter.Unmarshal(convertToJsonString, &memberReq)

    result, httpStatus, err := h.StoreService.Login(memberReq)
    if err != nil {
        return h.AsWebResponse(ctx, httpStatus, err.Error(), nil)
    }

    return h.AsMobileJson(ctx, httpStatus, "Login Success", result)
}
