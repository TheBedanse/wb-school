package database

import (
	"context"
	"fmt"

	"L0/internal/models"

	"github.com/jackc/pgx/v5"
)

func NewOrderRepository(db *pgx.Conn) *Database {
	return &Database{Conn: db}
}

func (r *Database) GetOrderByUID(ctx context.Context, orderUID string) (*models.Order, error) {
	order := &models.Order{}
	query := `
		SELECT order_uid, track_number, entry, locale, internal_signature, 
		       customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
		FROM orders 
		WHERE order_uid = $1`

	err := r.Conn.QueryRow(ctx, query, orderUID).Scan(
		&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale,
		&order.InternalSignature, &order.CustomerID, &order.DeliveryService,
		&order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("order not found")
		}
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	delivery, err := r.getDeliveryByOrderUID(ctx, orderUID)
	if err != nil {
		return nil, err
	}
	order.Delivery = *delivery

	payment, err := r.getPaymentByOrderUID(ctx, orderUID)
	if err != nil {
		return nil, err
	}
	order.Payment = *payment

	items, err := r.getItemsByOrderUID(ctx, orderUID)
	if err != nil {
		return nil, err
	}
	order.Items = items

	return order, nil
}

func (r *Database) getDeliveryByOrderUID(ctx context.Context, orderUID string) (*models.Delivery, error) {
	delivery := &models.Delivery{}
	query := `
		SELECT name, phone, zip, city, address, region, email 
		FROM deliveries 
		WHERE order_uid = $1`

	err := r.Conn.QueryRow(ctx, query, orderUID).Scan(
		&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City,
		&delivery.Address, &delivery.Region, &delivery.Email,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get delivery: %w", err)
	}
	return delivery, nil
}

func (r *Database) getPaymentByOrderUID(ctx context.Context, orderUID string) (*models.Payment, error) {
	payment := &models.Payment{}
	query := `
		SELECT transaction, request_id, currency, provider, amount, payment_dt, 
		       bank, delivery_cost, goods_total, custom_fee 
		FROM payments 
		WHERE order_uid = $1`

	err := r.Conn.QueryRow(ctx, query, orderUID).Scan(
		&payment.Transaction, &payment.RequestID, &payment.Currency, &payment.Provider,
		&payment.Amount, &payment.PaymentDt, &payment.Bank, &payment.DeliveryCost,
		&payment.GoodsTotal, &payment.CustomFee,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}
	return payment, nil
}

func (r *Database) getItemsByOrderUID(ctx context.Context, orderUID string) ([]models.Item, error) {
	query := `
        SELECT chrt_id, track_number, price, rid, name, sale, size, 
               total_price, nm_id, brand, status 
        FROM items 
        WHERE order_uid = $1`

	rows, err := r.Conn.Query(ctx, query, orderUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get items: %w", err)
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		err := rows.Scan(
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name,
			&item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		items = append(items, item)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating items: %w", err)
	}

	return items, nil
}

func (r *Database) GetAllOrderUIDs(ctx context.Context) ([]string, error) {
	query := `SELECT order_uid FROM orders ORDER BY date_created DESC`

	rows, err := r.Conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get order UIDs: %w", err)
	}
	defer rows.Close()

	var orderUIDs []string
	for rows.Next() {
		var uid string
		if err := rows.Scan(&uid); err != nil {
			return nil, fmt.Errorf("failed to scan order UID: %w", err)
		}
		orderUIDs = append(orderUIDs, uid)
	}

	return orderUIDs, nil
}
