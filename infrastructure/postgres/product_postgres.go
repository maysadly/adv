package postgres

import (
    "context"
    "FoodStore-AdvProg2/domain"
    "github.com/google/uuid"
)

type ProductPostgresRepo struct{}

func NewProductPostgresRepo() *ProductPostgresRepo {
    return &ProductPostgresRepo{}
}

func (r *ProductPostgresRepo) Save(product domain.Product) error {
    query := `INSERT INTO products (id, name, price, stock) VALUES ($1, $2, $3, $4)`
    _, err := DB.Exec(context.Background(), query,
        uuid.New().String(), product.Name, product.Price, product.Stock)
    return err
}

func (r *ProductPostgresRepo) FindByID(id string) (domain.Product, error) {
    query := `SELECT id, name, price, stock FROM products WHERE id = $1`
    row := DB.QueryRow(context.Background(), query, id)

    var p domain.Product
    err := row.Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
    return p, err
}

func (r *ProductPostgresRepo) Update(id string, product domain.Product) error {
    query := `UPDATE products SET name=$1, price=$2, stock=$3 WHERE id=$4`
    _, err := DB.Exec(context.Background(), query, product.Name, product.Price, product.Stock, id)
    return err
}

func (r *ProductPostgresRepo) Delete(id string) error {
    query := `DELETE FROM products WHERE id=$1`
    _, err := DB.Exec(context.Background(), query, id)
    return err
}

func (r *ProductPostgresRepo) FindAll() ([]domain.Product, error) {
    query := `SELECT id, name, price, stock FROM products`
    rows, err := DB.Query(context.Background(), query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var products []domain.Product
    for rows.Next() {
        var p domain.Product
        err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock)
        if err != nil {
            return nil, err
        }
        products = append(products, p)
    }

    return products, nil
}