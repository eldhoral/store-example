package service

import (
    "database/sql"
    "errors"
    "net/http"
    "time"

    modelCart "store-api/internal/store/domain/cart"
    modelTransaction "store-api/internal/store/domain/transaction"
    presenterCart "store-api/internal/store/presenter/cart"
    presenterMember "store-api/internal/store/presenter/member"
    presenterProduct "store-api/internal/store/presenter/product"
    presenterTransaction "store-api/internal/store/presenter/transaction"
    "store-api/internal/store/repository"
    "store-api/pkg/security"

    "github.com/jinzhu/copier"
    "github.com/spf13/cast"
    "golang.org/x/crypto/bcrypt"
)

// NewService creates new user service
func NewService(repo repository.StoreRepository) StoreService {
    return &service{
        repo: repo,
    }
}

type service struct {
    repo repository.StoreRepository
}

func (s service) ListProduct(request presenterProduct.ProductRequest) (result []presenterProduct.ProductResponse, httpStatus int, err error) {
    findAllProduct, err := s.repo.ListProduct(request.Category)
    if err == sql.ErrNoRows {
        httpStatus = http.StatusNotFound
        err = errors.New("Product not found")
        return
    }
    if err != nil {
        httpStatus = http.StatusInternalServerError
        return
    }
    copier.Copy(&result, &findAllProduct)
    return
}

func (s service) AddToCart(request presenterCart.CartRequest) (httpStatus int, err error) {
    var (
        cart = modelCart.Cart{}
    )

    copier.Copy(&cart, &request)
    err = s.repo.CreateCart(cart)
    if err != nil {
        httpStatus = http.StatusInternalServerError
        return
    }

    return
}

func (s service) ViewCart(request presenterCart.CartViewRequest) (result []presenterCart.CartResponse, httpStatus int, err error) {
    findAllCart, err := s.repo.GetCart(request.MemberID)
    if err == sql.ErrNoRows {
        httpStatus = http.StatusNotFound
        err = errors.New("Cart not found")
        return
    }
    if err != nil {
        httpStatus = http.StatusInternalServerError
        return
    }
    copier.Copy(&result, &findAllCart)
    return
}

func (s service) DeleteProductInCart(request presenterCart.CartProductDeleteRequest) (httpStatus int, err error) {
    err = s.repo.DeleteProductInCart(request.MemberID, request.ProductID)
    if err == sql.ErrNoRows {
        httpStatus = http.StatusNotFound
        err = errors.New("Cart not found")
        return
    }
    if err != nil {
        httpStatus = http.StatusInternalServerError
        return
    }
    return
}

func (s service) CreateTransaction(request presenterTransaction.TransactionRequest) (httpStatus int, err error) {
    var (
        transaction = modelTransaction.Transactions{}
    )

    getProduct, err := s.repo.GetProduct(request.ProductID)
    if err == sql.ErrNoRows {
        httpStatus = http.StatusNotFound
        err = errors.New("Product not found")
        return
    }
    if err != nil {
        httpStatus = http.StatusInternalServerError
        return
    }

    defer func() {
        if err != nil {
            transaction.Status = "failed"
            httpStatus = http.StatusInternalServerError
            err = s.repo.InsertFailedTransaction(transaction)
        }
    }()

    totalTransactionAmount := getProduct.Price * float64(request.Quantity)
    deductedStockProduct := getProduct.Stock - request.Quantity
    if deductedStockProduct < 0 {
        httpStatus = http.StatusOK
        err = errors.New("Quantity not enough")
        return
    }
    copier.Copy(&transaction, &request)

    transaction.Amount = totalTransactionAmount
    transaction.AmountFee = 0
    transaction.Status = "success"

    err = s.repo.CreateTransaction(transaction, deductedStockProduct)
    if err != nil {
        httpStatus = http.StatusInternalServerError
        return
    }

    return
}

func (s service) Login(request presenterMember.LoginRequest) (result presenterMember.LoginResponse, httpStatus int, err error) {
    memberData, err := s.repo.GetMemberByUsername(request.Username)
    if err == sql.ErrNoRows {
        httpStatus = http.StatusNotFound
        err = errors.New("Member not found")
        return
    }
    if err != nil {
        httpStatus = http.StatusInternalServerError
        //err = errors.New("Error in store service")
        return
    }

    if !compareBcrypt(memberData.Credential, request.Password+memberData.Salt) {
        httpStatus = http.StatusUnauthorized
        err = errors.New("Unauthorized")
        return
    }

    var (
        crypt, _ = security.New("123")
        sess     = security.Session{
            UserId:   cast.ToString(memberData.ID),
            Username: memberData.Username,
            Name:     "-",
            Role:     "-",
            Iat:      0,
            Expired:  time.Now().Add(time.Hour).Unix(),
        }
    )
    token, err := sess.Encrypt(crypt)
    if err != nil {
        return
    }

    result.Token = token
    result.IdUser = memberData.ID

    return

}

func compareBcrypt(hashedString, plainString string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hashedString), []byte(plainString))
    if err != nil {
        return false
    }

    return true
}
