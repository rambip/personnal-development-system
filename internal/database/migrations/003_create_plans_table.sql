CREATE TABLE IF NOT EXISTS plans (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    resources_required TEXT,
    value_id INTEGER NOT NULL,
    FOREIGN KEY (value_id) REFERENCES `values` (id)
);
