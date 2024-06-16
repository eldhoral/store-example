package service

import (
    presenterCart "store-api/internal/store/presenter/cart"
    presenterMember "store-api/internal/store/presenter/member"
    presenterProduct "store-api/internal/store/presenter/product"
    presenterTransaction "store-api/internal/store/presenter/transaction"
)

type StoreService interface {
    ListProduct(request presenterProduct.ProductRequest) (result []presenterProduct.ProductResponse, httpStatus int, err error)
    AddToCart(request presenterCart.CartRequest) (httpStatus int, err error)
    ViewCart(request presenterCart.CartViewRequest) (result []presenterCart.CartResponse, httpStatus int, err error)
    DeleteProductInCart(request presenterCart.CartProductDeleteRequest) (httpStatus int, err error)
    CreateTransaction(request presenterTransaction.TransactionRequest) (httpStatus int, err error)
    Login(request presenterMember.LoginRequest) (result presenterMember.LoginResponse, httpStatus int, err error)
}
