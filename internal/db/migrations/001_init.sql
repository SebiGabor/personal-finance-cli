CREATE TABLE IF NOT EXISTS transactions (
                                            id INTEGER PRIMARY KEY AUTOINCREMENT,
                                            date TEXT NOT NULL,
                                            description TEXT,
                                            amount REAL NOT NULL,
                                            category TEXT,
                                            account TEXT,
                                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS categories (
                                          id INTEGER PRIMARY KEY AUTOINCREMENT,
                                          name TEXT NOT NULL UNIQUE,
                                          regex_rule TEXT
);

CREATE TABLE IF NOT EXISTS budgets (
                                       id INTEGER PRIMARY KEY AUTOINCREMENT,
                                       category TEXT NOT NULL,
                                       amount REAL NOT NULL,
                                       period TEXT NOT NULL
);

