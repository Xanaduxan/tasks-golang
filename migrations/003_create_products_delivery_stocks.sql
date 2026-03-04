CREATE TABLE IF NOT EXISTS products
(
    id         UUID PRIMARY KEY,
    name       TEXT      NOT NULL,
    price      numeric   NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()

);



CREATE TABLE IF NOT EXISTS deliveries
(
    id         UUID PRIMARY KEY,
    status     TEXT      NOT NULL DEFAULT 'awaiting'
        CHECK (status in ('awaiting', 'on_path', 'processing', 'checked', 'accepted')),
    user_id    UUID      NOT NULL REFERENCES users (id),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS delivery_items
(
    id          UUID PRIMARY KEY,
    delivery_id UUID      NOT NULL REFERENCES deliveries (id) ON DELETE CASCADE,
    product_id  UUID      NOT NULL REFERENCES products (id),
    created_at  TIMESTAMP NOT NULL DEFAULT now(),
    quantity    INTEGER   NOT NULL CHECK (quantity > 0),
    UNIQUE (delivery_id, product_id)
);

CREATE TABLE IF NOT EXISTS stocks
(
    product_id UUID PRIMARY KEY REFERENCES products (id),
    quantity   INTEGER   NOT NULL CHECK (quantity >= 0),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_deliveries_user_id ON deliveries (user_id);
CREATE INDEX IF NOT EXISTS idx_delivery_items_delivery_id ON delivery_items (delivery_id);