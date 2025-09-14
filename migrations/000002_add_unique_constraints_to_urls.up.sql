DROP INDEX IF EXISTS idx_urls_code;
DROP INDEX IF EXISTS idx_urls_url;

CREATE UNIQUE INDEX idx_urls_code ON urls(code);
CREATE UNIQUE INDEX idx_urls_url ON urls(url);
