CREATE TABLE IF NOT EXISTS shops
(
    id         UUID PRIMARY KEY,
    name       TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now()

);

CREATE TABLE IF NOT EXISTS payments
(
    id                          UUID PRIMARY KEY,
    status                      TEXT      NOT NULL DEFAULT 'NEW'
        CHECK (status in ('NEW', 'WAITING_FOR_VALIDATION_2', 'FAILED', 'READY_FOR_CLOSURE', 'CLOSED')),
    shop_id                     UUID      NOT NULL REFERENCES shops (id),
    attempts                    INTEGER   NOT NULL DEFAULT 0 CHECK (attempts >= 0),
    amount                      NUMERIC   NOT NULL DEFAULT 0 CHECK (amount >= 0),
    waiting_for_validation_2_at TIMESTAMP NULL,
    created_at                  TIMESTAMP NOT NULL DEFAULT now(),
    updated_at                  TIMESTAMP NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_payments_shop_id ON payments (shop_id);