CREATE TABLE IF NOT EXISTS address 
(
  address_id           INTEGER PRIMARY KEY AUTOINCREMENT,   
  address_address      TEXT    NOT NULL,
  address_number       TEXT    NOT NULL,
  address_neighborhood TEXT    NOT NULL,
  address_city         TEXT    NOT NULL,
  address_state        TEXT    NOT NULL,
  address_cep          TEXT    NOT NULL,
  address_created_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  user_id              INTEGER NOT NULL UNIQUE,
  FOREIGN KEY (user_id) REFERENCES users(user_id)
);