-- === Налог по продаже ===
-- Ставка фиксируется в строке продажи на момент её создания,
-- чтобы суммы с налогом считались стабильно даже при изменении ставки продукта.
ALTER TABLE sales ADD COLUMN IF NOT EXISTS tax_rate NUMERIC(5, 2) NOT NULL DEFAULT 0 CHECK (tax_rate >= 0);
