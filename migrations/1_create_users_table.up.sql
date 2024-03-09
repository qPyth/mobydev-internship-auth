CREATE TABLE IF NOT EXISTS users (
                                     id          INTEGER PRIMARY KEY AUTOINCREMENT,
                                     name        TEXT,
                                     email       TEXT NOT NULL UNIQUE,
                                     hash_pass   TEXT NOT NULL,
                                     phone_number TEXT,
                                     b_day       DATETIME,
                                     created_at  DATETIME NOT NULL,
                                     updated_at  DATETIME NOT NULL
);
