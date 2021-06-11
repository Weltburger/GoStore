package store

import (
	"GoStore/pkg/models"
	"context"
	"golang.org/x/sync/errgroup"
)

type OrderRepository struct {
	store *Store
}

func (orderRepository *OrderRepository) InsertOrder(ctx context.Context, order *models.Order) (*models.Order, error) {
	if err := orderRepository.store.DB.QueryRowxContext(ctx, `insert into "public"."orders"(product_uuid, price, quantity, email, status) 
VALUES($1, $2, $3, $4, $5) returning uuid`, order.ProductUUID, order.Price, order.Quantity, order.Email,
	order.Status).Scan(&order.UUID); err != nil {
		return nil, err
	}

	return order, nil
}

func (orderRepository *OrderRepository) OrdersByEmail(ctx context.Context, order *models.Order) ([]*models.Order, int64, error) {
	orders := make([]*models.Order, 0)
	total := int64(0)
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		rows, err := orderRepository.store.DB.QueryxContext(ctx, `SELECT uuid, 
       product_uuid, 
       price,
       quantity,
       email,
       status,
       create_at, 
       updated_at, 
       deleted_at 
FROM "public"."orders" ps 
WHERE email = $1`, order.Email)

		if err != nil {
			return err
		}

		defer rows.Close()

		for rows.Next() {
			order := &models.Order{}
			if err := rows.StructScan(order); err != nil {
				return err
			}
			orders = append(orders, order)
		}
		return nil
	})

	eg.Go(func() error {
		if err := orderRepository.store.DB.QueryRowxContext(ctx, `SELECT count(1)
FROM "public"."orders" 
WHERE email = $1`, order.Email).Scan(&total); err != nil {
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (orderRepository *OrderRepository) AllOrders(ctx context.Context, take int64, skip int64) ([]*models.Order, int64, error) {
	orders := make([]*models.Order, 0)
	total := int64(0)
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		rows, err := orderRepository.store.DB.QueryxContext(ctx, `SELECT uuid, 
       product_uuid, 
       price,
       quantity,
       email,
       status,
       create_at, 
       updated_at, 
       deleted_at 
FROM "public"."orders" ps offset $1 limit $2`, skip, take)

		if err != nil {
			return err
		}

		defer rows.Close()

		for rows.Next() {
			order := &models.Order{}
			if err := rows.StructScan(order); err != nil {
				return err
			}
			orders = append(orders, order)
		}
		return nil
	})

	eg.Go(func() error {
		if err := orderRepository.store.DB.QueryRowxContext(ctx, `SELECT count(1)
FROM "public"."orders"`).Scan(&total); err != nil {
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (orderRepository *OrderRepository) UpdateOrderStatus(ctx context.Context, order *models.Order) error {
	if _, err := orderRepository.store.DB.ExecContext(ctx, `update "public"."orders" SET status = $1,
        deleted_at = NULL
		where uuid = $2`, order.Status, order.UUID); err != nil {
		return err
	}

	return nil
}

func (orderRepository *OrderRepository) UpdateOrder(ctx context.Context, order *models.Order) error {
	if _, err := orderRepository.store.DB.ExecContext(ctx, `update "public"."orders" SET product_uuid = $1,
        price = $2,
		quantity = $3, 
		email = $4, 
		status = $5, 
		updated_at = now(),
        deleted_at = NULL 
		where uuid = $6`,
		order.ProductUUID, order.Price, order.Quantity, order.Email, order.Status, order.UUID); err != nil {
		return err
	}

	return nil
}

func (orderRepository *OrderRepository) DeleteOrder(ctx context.Context, order *models.Order) error {
	if _, err := orderRepository.store.DB.ExecContext(ctx, `update "public"."orders" SET deleted_at = now()
		where uuid = $1`, order.UUID); err != nil {
		return err
	}

	return nil
}
