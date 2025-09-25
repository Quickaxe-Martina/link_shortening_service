CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    code VARCHAR(10) NOT NULL,
    url TEXT NOT NULL
);

CREATE INDEX idx_urls_code ON urls(code);

CREATE INDEX idx_urls_url ON urls(url); 