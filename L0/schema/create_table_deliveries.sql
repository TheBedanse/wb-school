CREATE TABLE deliveries (
    order_uid VARCHAR(100) PRIMARY KEY REFERENCES orders(order_uid),
    name      VARCHAR(100) NOT NULL,
    phone     VARCHAR(20) NOT NULL,
    zip       VARCHAR(20) NOT NULL,
    city      VARCHAR(60) NOT NULL,
    address   VARCHAR(150) NOT NULL,
    region    VARCHAR(50) NOT NULL,
    email     VARCHAR(100) NOT NULL
);