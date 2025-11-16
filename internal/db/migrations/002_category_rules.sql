CREATE TABLE IF NOT EXISTS category_rules (
                                              id INTEGER PRIMARY KEY AUTOINCREMENT,
                                              pattern TEXT NOT NULL,
                                              category TEXT NOT NULL
);
