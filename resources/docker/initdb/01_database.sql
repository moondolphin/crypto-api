CREATE TABLE IF NOT EXISTS coins (
  id BIGINT AUTO_INCREMENT PRIMARY KEY,
  symbol VARCHAR(16) NOT NULL UNIQUE,
  enabled BOOLEAN NOT NULL DEFAULT TRUE,
  coingecko_id VARCHAR(64),
  binance_symbol VARCHAR(32),
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO coins (symbol, enabled, coingecko_id, binance_symbol)
VALUES
  ('BTC', TRUE, 'bitcoin', 'BTCUSDT'),
  ('ETH', TRUE, 'ethereum', 'ETHUSDT')
ON DUPLICATE KEY UPDATE
  enabled = VALUES(enabled),
  coingecko_id = VALUES(coingecko_id),
  binance_symbol = VALUES(binance_symbol);
