package store

import (
	"context"
	"database/sql"
	"errors"
	"golang.org/x/sync/errgroup"

	//"database/sql"
	//"errors"
	"GoStore/pkg/models"
	//"golang.org/x/sync/errgroup"
)

type ProductRepository struct {
	store *Store
}

func (productRepository *ProductRepository) InsertProduct(ctx context.Context, product *models.Product) (*models.Product, error) {
	if err := productRepository.store.DB.QueryRowxContext(ctx, `insert into "public"."products"(user_uuid, name, description, price, quantity) 
VALUES($1, $2, $3, $4, $5) returning uuid`, product.User.UUID, product.Name, product.Description, product.Price,
product.Quantity).Scan(&product.UUID); err != nil {
		return nil, err
	}

	return product, nil
}

func (productRepository *ProductRepository) DeleteProduct(ctx context.Context, product *models.Product) error {
	if _, err := productRepository.store.DB.ExecContext(ctx, `update "public"."products" SET deleted_at = now(),
    price = 0, quantity = 0 where uuid = $1`, product.UUID); err != nil {
		return err
	}

	return nil
}

func (productRepository *ProductRepository) UpdateProduct(ctx context.Context, product *models.Product) error {
	if _, err := productRepository.store.DB.ExecContext(ctx, `update "public"."products" SET name = $1, 
    	description = $2,
        price = $3,
    	quantity = $4,
        updated_at = now(),
        deleted_at = NULL 
		where uuid = $5`,
		product.Name, product.Description, product.Price, product.Quantity, product.UUID); err != nil {
		return err
	}

	return nil
}

func (productRepository *ProductRepository) ProductByUUID(ctx context.Context, product *models.Product) (*models.Product, error) {
	if err := productRepository.store.DB.QueryRowxContext(ctx, `SELECT uuid, 
       name, 
       description, 
       price, 
       quantity, 
       create_at, 
       updated_at, 
       deleted_at
	FROM "public"."products" ps 
	WHERE ps.uuid = $1`, product.UUID).StructScan(product); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("there is no such product")
		}
		return nil, err
	}

	return product, nil
}

func (productRepository *ProductRepository) AllProducts(ctx context.Context, take int64, skip int64) ([]*models.Product, int64, error) {
	products := make([]*models.Product, 0)
	total := int64(0)
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		rows, err := productRepository.store.DB.QueryxContext(ctx, `SELECT uuid, 
       name, 
       description,
       price,
       quantity,
       create_at, 
       updated_at, 
       deleted_at 
FROM "public"."products" ps offset $1 limit $2`, skip, take)

		if err != nil {
			return err
		}

		defer rows.Close()

		for rows.Next() {
			product := &models.Product{}
			if err := rows.StructScan(product); err != nil {
				return err
			}
			products = append(products, product)
		}
		return nil
	})

	eg.Go(func() error {
		if err := productRepository.store.DB.QueryRowxContext(ctx, `SELECT count(1)
FROM "public"."products"`).Scan(&total); err != nil {
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, 0, err
	}

	return products, total, nil
}
