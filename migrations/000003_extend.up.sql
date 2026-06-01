-- === Расширение пользователей ===
ALTER TABLE users ADD COLUMN IF NOT EXISTS full_name VARCHAR(150) NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN IF NOT EXISTS email     VARCHAR(150) NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN IF NOT EXISTS status    VARCHAR(20)  NOT NULL DEFAULT 'active';
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_status_check;
ALTER TABLE users ADD CONSTRAINT users_status_check CHECK (status IN ('active', 'blocked'));

-- === Расширение продуктов (склад, единица, налоговая ставка) ===
ALTER TABLE products ADD COLUMN IF NOT EXISTS stock    INTEGER       NOT NULL DEFAULT 0;
ALTER TABLE products ADD COLUMN IF NOT EXISTS unit     VARCHAR(20)   NOT NULL DEFAULT 'шт';
ALTER TABLE products ADD COLUMN IF NOT EXISTS tax_rate NUMERIC(5, 2) NOT NULL DEFAULT 0;
-- Наполняем разумными значениями уже существующие строки.
UPDATE products SET stock = ((id * 17) % 100), tax_rate = 12 WHERE stock = 0;

-- === Заявки на закупку: номер, клиент, выравнивание статусов ===
ALTER TABLE purchase_requests ADD COLUMN IF NOT EXISTS number      VARCHAR(50);
ALTER TABLE purchase_requests ADD COLUMN IF NOT EXISTS client_name VARCHAR(150) NOT NULL DEFAULT '';
-- number заполняется приложением (CTE при вставке) и миграцией для старых строк.
-- NOT NULL не ставим: внутри INSERT...RETURNING значение проставляется следующим
-- шагом того же запроса, поэтому колонка должна допускать временный NULL.
UPDATE purchase_requests SET number = 'REQ-' || LPAD(id::text, 4, '0') WHERE number IS NULL;
CREATE UNIQUE INDEX IF NOT EXISTS idx_purchase_requests_number ON purchase_requests(number);

ALTER TABLE purchase_requests DROP CONSTRAINT IF EXISTS purchase_requests_status_check;
UPDATE purchase_requests SET status = CASE status
    WHEN 'new'         THEN 'pending'
    WHEN 'in_progress' THEN 'checking'
    ELSE status
END;
ALTER TABLE purchase_requests ADD CONSTRAINT purchase_requests_status_check
    CHECK (status IN ('draft', 'pending', 'checking', 'approved', 'ordered', 'completed', 'rejected'));
ALTER TABLE purchase_requests ALTER COLUMN status SET DEFAULT 'pending';

-- === Счета-фактуры: срок оплаты, выравнивание статусов ===
ALTER TABLE invoices ADD COLUMN IF NOT EXISTS due_date DATE;
UPDATE invoices SET due_date = issued_date + INTERVAL '14 days' WHERE due_date IS NULL;

ALTER TABLE invoices DROP CONSTRAINT IF EXISTS invoices_status_check;
UPDATE invoices SET status = CASE status
    WHEN 'unpaid' THEN 'issued'
    WHEN 'void'   THEN 'draft'
    ELSE status
END;
ALTER TABLE invoices ADD CONSTRAINT invoices_status_check
    CHECK (status IN ('draft', 'issued', 'paid', 'overdue'));
ALTER TABLE invoices ALTER COLUMN status SET DEFAULT 'issued';

-- === Справочник налоговых ставок ===
CREATE TABLE IF NOT EXISTS taxes (
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(100)  NOT NULL,
    rate       NUMERIC(5, 2) NOT NULL DEFAULT 0 CHECK (rate >= 0),
    active     BOOLEAN       NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ   NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ   NOT NULL DEFAULT now()
);

-- === Справочник единиц измерения ===
CREATE TABLE IF NOT EXISTS units (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    short_name  VARCHAR(20)  NOT NULL DEFAULT '',
    description  TEXT        NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT now()
);

-- === Продажи (выполненные продажи и установки) ===
CREATE TABLE IF NOT EXISTS sales (
    id                  SERIAL PRIMARY KEY,
    invoice_id          INTEGER        REFERENCES invoices(id) ON DELETE SET NULL,
    product_name        VARCHAR(150)   NOT NULL DEFAULT '',
    quantity            INTEGER        NOT NULL DEFAULT 1 CHECK (quantity >= 0),
    amount              NUMERIC(12, 2) NOT NULL DEFAULT 0 CHECK (amount >= 0),
    sold_at             DATE           NOT NULL DEFAULT CURRENT_DATE,
    installation_status VARCHAR(20)    NOT NULL DEFAULT 'not_required'
        CHECK (installation_status IN ('not_required', 'scheduled', 'completed')),
    created_at          TIMESTAMPTZ    NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ    NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_sales_invoice ON sales(invoice_id);
