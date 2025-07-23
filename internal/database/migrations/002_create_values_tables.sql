-- Create the "values" table (note the use of quotes to avoid conflicts with SQL keywords)
CREATE TABLE IF NOT EXISTS "values" (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Create the join table for parent-child relationships among values.
CREATE TABLE IF NOT EXISTS value_parents (
    value_id INTEGER NOT NULL,
    parent_value_id INTEGER NOT NULL,
    PRIMARY KEY (value_id, parent_value_id),
    FOREIGN KEY (value_id) REFERENCES "values"(id) ON DELETE CASCADE,
    FOREIGN KEY (parent_value_id) REFERENCES "values"(id) ON DELETE CASCADE
);
