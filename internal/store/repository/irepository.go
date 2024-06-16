package repository

import (
    modelCart "store-api/internal/store/domain/cart"
    modelMember "store-api/internal/store/domain/member"
    modelProduct "store-api/internal/store/domain/product"
    modelTransaction "store-api/internal/store/domain/transaction"
)

type StoreRepository interface {
    ListProduct(category string) (result []modelProduct.Product, err error)
    GetProduct(productId int) (result modelProduct.Product, err error)
    CreateCart(model modelCart.Cart) (err error)
    GetCart(memberId int) (result []modelCart.Cart, err error)
    DeleteProductInCart(memberId, productId int) (err error)
    CreateTransaction(model modelTransaction.Transactions, deductedStockProduct int) (err error)
    GetMemberByUsername(username string) (result modelMember.Member, err error)
    InsertFailedTransaction(model modelTransaction.Transactions) (err error)
}
