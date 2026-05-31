-- Отчёты (сущность «Отчёт» из диаграммы классов ТЗ)
CREATE TABLE IF NOT EXISTS reports (
    id           SERIAL PRIMARY KEY,
    category_id  INTEGER      REFERENCES categories(id) ON DELETE SET NULL,
    title        VARCHAR(150) NOT NULL,
    content      TEXT         NOT NULL DEFAULT '',
    created_date DATE         NOT NULL DEFAULT CURRENT_DATE,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_reports_category ON reports(category_id);
CREATE INDEX IF NOT EXISTS idx_reports_created_date ON reports(created_date);
