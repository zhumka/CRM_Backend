-- Пользователи системы (роли: admin, user)
CREATE TABLE IF NOT EXISTS users (
    id            SERIAL PRIMARY KEY,
    username      VARCHAR(50)  NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role          VARCHAR(20)  NOT NULL DEFAULT 'user' CHECK (role IN ('admin', 'user')),
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT now()
);

-- Категории продукции
CREATE TABLE IF NOT EXISTS categories (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    description TEXT         NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
);

-- Поставщики
CREATE TABLE IF NOT EXISTS suppliers (
    id           SERIAL PRIMARY KEY,
    name         VARCHAR(100) NOT NULL,
    contact_name VARCHAR(100) NOT NULL DEFAULT '',
    phone        VARCHAR(20)  NOT NULL DEFAULT '',
    email        VARCHAR(100) NOT NULL DEFAULT '',
    address      TEXT         NOT NULL DEFAULT '',
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT now()
);

-- Продукты (товары)
CREATE TABLE IF NOT EXISTS products (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100)   NOT NULL,
    description TEXT           NOT NULL DEFAULT '',
    price       NUMERIC(10, 2) NOT NULL DEFAULT 0 CHECK (price >= 0),
    category_id INTEGER        REFERENCES categories(id) ON DELETE SET NULL,
    supplier_id INTEGER        REFERENCES suppliers(id) ON DELETE SET NULL,
    created_at  TIMESTAMPTZ    NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ    NOT NULL DEFAULT now()
);

-- Заявки на закупку продукта
CREATE TABLE IF NOT EXISTS purchase_requests (
    id          SERIAL PRIMARY KEY,
    user_id     INTEGER     NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_id  INTEGER     NOT NULL REFERENCES products(id) ON DELETE RESTRICT,
    quantity    INTEGER     NOT NULL DEFAULT 1 CHECK (quantity >= 1),
    status      VARCHAR(20) NOT NULL DEFAULT 'new'
        CHECK (status IN ('new', 'in_progress', 'approved', 'rejected', 'completed')),
    comment     TEXT        NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Счета-фактуры
CREATE TABLE IF NOT EXISTS invoices (
    id                  SERIAL PRIMARY KEY,
    number              VARCHAR(50)    NOT NULL UNIQUE,
    purchase_request_id INTEGER        REFERENCES purchase_requests(id) ON DELETE SET NULL,
    amount              NUMERIC(12, 2) NOT NULL DEFAULT 0 CHECK (amount >= 0),
    status              VARCHAR(20)    NOT NULL DEFAULT 'unpaid'
        CHECK (status IN ('unpaid', 'paid', 'void')),
    issued_date         DATE           NOT NULL DEFAULT CURRENT_DATE,
    created_at          TIMESTAMPTZ    NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ    NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_products_category ON products(category_id);
CREATE INDEX IF NOT EXISTS idx_products_supplier ON products(supplier_id);
CREATE INDEX IF NOT EXISTS idx_purchase_requests_user ON purchase_requests(user_id);
CREATE INDEX IF NOT EXISTS idx_purchase_requests_status ON purchase_requests(status);
CREATE INDEX IF NOT EXISTS idx_invoices_request ON invoices(purchase_request_id);
