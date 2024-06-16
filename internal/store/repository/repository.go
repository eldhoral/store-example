package repository

import (
    "fmt"
    "time"

    "github.com/jmoiron/sqlx"

    modelCart "store-api/internal/store/domain/cart"
    modelMember "store-api/internal/store/domain/member"
    modelProduct "store-api/internal/store/domain/product"
    modelTransaction "store-api/internal/store/domain/transaction"
)

// NewStoreRepository creates new repository
func NewStoreRepository(db *sqlx.DB) StoreRepository {
    return &repo{db: db}
}

type repo struct {
    db *sqlx.DB
}

func (r repo) ListProduct(category string) (result []modelProduct.Product, err error) {
    query := fmt.Sprintf("SELECT id, name, category, price, stock FROM %s", modelProduct.TableName)
    if category != "" {
        query += fmt.Sprintf(" WHERE category = '%s'", category)
    }

    err = r.db.Select(&result, query)
    return
}

func (r repo) GetProduct(productId int) (result modelProduct.Product, err error) {
    query := fmt.Sprintf("SELECT id, name, category, price, stock FROM %s", modelProduct.TableName)
    query += fmt.Sprintf(" WHERE id = %d", productId)

    err = r.db.Get(&result, query)
    return
}

func (r repo) CreateCart(model modelCart.Cart) (err error) {
    arg := map[string]interface{}{
        "member_id":  model.MemberID,
        "product_id": model.ProductID,
        "quantity":   model.Quantity,
        "is_active":  true,
    }

    query := fmt.Sprintf(`INSERT INTO %s SET member_id = :member_id, product_id = :product_id, 
quantity = :quantity, is_active = :is_active`, modelCart.TableName)

    _, err = r.db.NamedExec(query, arg)
    if err != nil {
        return err
    }

    return
}

func (r repo) GetCart(memberId int) (result []modelCart.Cart, err error) {
    query := fmt.Sprintf("SELECT id, member_id, product_id, quantity, is_active FROM %s", modelCart.TableName)
    query += fmt.Sprintf(" WHERE member_id = %d", memberId)

    err = r.db.Select(&result, query)
    return
}

func (r repo) DeleteProductInCart(memberId, productId int) (err error) {
    query := fmt.Sprintf("UPDATE %s SET is_active = false", modelCart.TableName)
    query += fmt.Sprintf(" WHERE member_id = %d AND product_id = %d", memberId, productId)

    _, err = r.db.Exec(query)
    if err != nil {
        return
    }
    return
}

func (r repo) CreateTransaction(model modelTransaction.Transactions, deductedStockProduct int) (err error) {
    arg := map[string]interface{}{
        "member_id":      model.MemberID,
        "product_id":     model.ProductID,
        "trx_code":       model.TrxCode,
        "channel_id":     model.ChannelID,
        "channel_ref_no": model.ChannelRefNo,
        "channel_time":   model.ChannelTime,
        "channel_date":   model.ChannelDate,
        "amount":         model.Amount,
        "amount_fee":     model.AmountFee,
        "status":         model.Status,
        "quantity":       model.Quantity,
        "created_date":   time.Time{},
        "updated_date":   time.Time{},
    }

    tx, err := r.db.Beginx()
    defer func() {
        if err == nil {
            err = tx.Commit()
        } else {
            err = tx.Rollback()
        }
    }()

    // Delete product in cart
    query := fmt.Sprintf("UPDATE %s SET is_active = false", modelCart.TableName)
    query += fmt.Sprintf(" WHERE member_id = %d AND product_id = %d", model.MemberID, model.ProductID)
    _, err = tx.Exec(query)
    if err != nil {
        return
    }

    // Deduct Stock in Product
    query = fmt.Sprintf("UPDATE %s SET stock = %d where id = %d", modelProduct.TableName, deductedStockProduct, model.ProductID)
    _, err = tx.Exec(query)
    if err != nil {
        fmt.Println(err)
        return
    }

    // Create Transaction
    query = fmt.Sprintf(`INSERT INTO %s SET member_id = :member_id, product_id = :product_id, 
    trx_code = :trx_code, channel_id = :channel_id, channel_ref_no = :channel_ref_no, channel_time = :channel_time, 
    channel_date = :channel_date, amount = :amount, amount_fee = :amount_fee, status = :status,
    quantity = :quantity, created_date = :created_date, updated_date = :updated_date`, modelTransaction.TableName)
    _, err = tx.NamedExec(query, arg)
    if err != nil {
        return err
    }

    return
}

func (r repo) GetMemberByUsername(username string) (result modelMember.Member, err error) {
    query := fmt.Sprintf("SELECT id, channel_id, username, credential, salt, created_date FROM %s WHERE username = '%s'", modelMember.TableName, username)

    err = r.db.Get(&result, query)
    if err != nil {
        return
    }
    return
}

func (r repo) InsertFailedTransaction(model modelTransaction.Transactions) (err error) {
    arg := map[string]interface{}{
        "member_id":      model.MemberID,
        "product_id":     model.ProductID,
        "trx_code":       model.TrxCode,
        "channel_id":     model.ChannelID,
        "channel_ref_no": model.ChannelRefNo,
        "channel_time":   model.ChannelTime,
        "channel_date":   model.ChannelDate,
        "amount":         model.Amount,
        "amount_fee":     model.AmountFee,
        "status":         model.Status,
        "quantity":       model.Quantity,
        "created_date":   time.Time{},
        "updated_date":   time.Time{},
    }

    query := fmt.Sprintf(`INSERT INTO %s SET member_id = :member_id, product_id = :product_id, 
    trx_code = :trx_code, channel_id = :channel_id, channel_ref_no = :channel_ref_no, channel_time = :channel_time, 
    channel_date = :channel_date, amount = :amount, amount_fee = :amount_fee, status = :status,
    quantity = :quantity, created_date = :created_date, updated_date = :updated_date`, modelTransaction.TableName)

    _, err = r.db.NamedExec(query, arg)
    if err != nil {
        return err
    }

    return
}
