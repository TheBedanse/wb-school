package database

import (
	"context"
	"fmt"
	"log"

	"L0/internal/models"

	"github.com/jackc/pgx/v5"
)

func (r *Database) SaveOrder(ctx context.Context, order *models.Order) error {
	tx, err := r.Conn.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	if err := r.saveOrderMain(ctx, tx, order); err != nil {
		return err
	}

	if err := r.saveDelivery(ctx, tx, order); err != nil {
		return err
	}

	if err := r.savePayment(ctx, tx, order); err != nil {
		return err
	}

	if err := r.saveItems(ctx, tx, order); err != nil {
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("Order saved to database(SaveOrder): %s", order.OrderUID)
	return nil
}

func (r *Database) saveOrderMain(ctx context.Context, tx pgx.Tx, order *models.Order) error {
	query := `
		INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, 
		                  customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (order_uid) DO NOTHING`

	_, err := tx.Exec(ctx, query,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale,
		order.InternalSignature, order.CustomerID, order.DeliveryService,
		order.Shardkey, order.SmID, order.DateCreated, order.OofShard,
	)
	return err
}

func (r *Database) saveDelivery(ctx context.Context, tx pgx.Tx, order *models.Order) error {
	query := `
		INSERT INTO deliveries (order_uid, name, phone, zip, city, address, region, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (order_uid) DO UPDATE SET
			name = EXCLUDED.name, phone = EXCLUDED.phone, zip = EXCLUDED.zip,
			city = EXCLUDED.city, address = EXCLUDED.address, region = EXCLUDED.region,
			email = EXCLUDED.email`

	_, err := tx.Exec(ctx, query,
		order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Address, order.Delivery.Region, order.Delivery.Email,
	)
	return err
}

func (r *Database) savePayment(ctx context.Context, tx pgx.Tx, order *models.Order) error {
	query := `
		INSERT INTO payments (order_uid, transaction, request_id, currency, provider, 
		                     amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (order_uid) DO UPDATE SET
			transaction = EXCLUDED.transaction, request_id = EXCLUDED.request_id,
			currency = EXCLUDED.currency, provider = EXCLUDED.provider,
			amount = EXCLUDED.amount, payment_dt = EXCLUDED.payment_dt,
			bank = EXCLUDED.bank, delivery_cost = EXCLUDED.delivery_cost,
			goods_total = EXCLUDED.goods_total, custom_fee = EXCLUDED.custom_fee`

	_, err := tx.Exec(ctx, query,
		order.OrderUID, order.Payment.Transaction, order.Payment.RequestID,
		order.Payment.Currency, order.Payment.Provider, order.Payment.Amount,
		order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost,
		order.Payment.GoodsTotal, order.Payment.CustomFee,
	)
	return err
}

func (r *Database) saveItems(ctx context.Context, tx pgx.Tx, order *models.Order) error {
	_, err := tx.Exec(ctx, "DELETE FROM items WHERE order_uid = $1", order.OrderUID)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		query := `
            INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, 
                              sale, size, total_price, nm_id, brand, status)
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`

		_, err := tx.Exec(ctx, query,
			order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.Rid,
			item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID,
			item.Brand, item.Status,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
