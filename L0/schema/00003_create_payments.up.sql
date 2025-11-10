-- +goose Up
CREATE TABLE payments (
    order_uid     VARCHAR(100) PRIMARY KEY REFERENCES orders(order_uid),
    transaction   VARCHAR(50) NOT NULL,
    request_id    VARCHAR(100) DEFAULT '',
    currency      VARCHAR(3) NOT NULL,
    provider      VARCHAR(50) NOT NULL,
    amount        INTEGER NOT NULL,
    payment_dt    BIGINT NOT NULL,
    bank          VARCHAR(50) NOT NULL,
    delivery_cost INTEGER NOT NULL,
    goods_total   INTEGER NOT NULL,
    custom_fee    INTEGER DEFAULT 0
);

-- +goose Down
DROP TABLE IF EXISTS payments;