DROP TABLE IF EXISTS sales;
DROP TABLE IF EXISTS units;
DROP TABLE IF EXISTS taxes;

ALTER TABLE invoices DROP CONSTRAINT IF EXISTS invoices_status_check;
ALTER TABLE invoices ADD CONSTRAINT invoices_status_check CHECK (status IN ('unpaid', 'paid', 'void'));
ALTER TABLE invoices ALTER COLUMN status SET DEFAULT 'unpaid';
ALTER TABLE invoices DROP COLUMN IF EXISTS due_date;

ALTER TABLE purchase_requests DROP CONSTRAINT IF EXISTS purchase_requests_status_check;
ALTER TABLE purchase_requests ADD CONSTRAINT purchase_requests_status_check
    CHECK (status IN ('new', 'in_progress', 'approved', 'rejected', 'completed'));
ALTER TABLE purchase_requests ALTER COLUMN status SET DEFAULT 'new';
DROP INDEX IF EXISTS idx_purchase_requests_number;
ALTER TABLE purchase_requests DROP COLUMN IF EXISTS number;
ALTER TABLE purchase_requests DROP COLUMN IF EXISTS client_name;

ALTER TABLE products DROP COLUMN IF EXISTS stock;
ALTER TABLE products DROP COLUMN IF EXISTS unit;
ALTER TABLE products DROP COLUMN IF EXISTS tax_rate;

ALTER TABLE users DROP CONSTRAINT IF EXISTS users_status_check;
ALTER TABLE users DROP COLUMN IF EXISTS full_name;
ALTER TABLE users DROP COLUMN IF EXISTS email;
ALTER TABLE users DROP COLUMN IF EXISTS status;
