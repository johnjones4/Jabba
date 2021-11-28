CREATE TABLE IF NOT EXISTS events (
  id INT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  event_vendor_type VARCHAR(255) NOT NULL,
  event_vendor_id VARCHAR(255) NOT NULL,
  created TIMESTAMP NOT NULL,
  vendor_info JSON NOT NULL,
  alerts JSON NOT NULL,
  is_normal BOOLEAN NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_vendor ON events (event_vendor_type, event_vendor_id);
