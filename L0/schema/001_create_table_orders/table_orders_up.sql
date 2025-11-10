-- +goose Up
CREATE TABLE orders (
    order_uid 			VARCHAR(100) PRIMARY KEY,
    track_number 		VARCHAR(50) NOT NULL,
    entry 				VARCHAR(10) NOT NULL,
    locale				VARCHAR(10) NOT NULL,
    internal_signature  VARCHAR(200) DEFAULT '',
    customer_id 		VARCHAR(100) NOT NULL,
    delivery_service 	VARCHAR(50) NOT NULL,
    shardkey 			VARCHAR(10) NOT NULL,
    sm_id 				INTEGER NOT NULL,
    date_created		TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    oof_shard 			VARCHAR(10) NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS orders;