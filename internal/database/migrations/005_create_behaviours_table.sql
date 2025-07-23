CREATE TABLE IF NOT EXISTS behaviours (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    description TEXT,
    conflicting_aim_id INTEGER NOT NULL,
    FOREIGN KEY (conflicting_aim_id) REFERENCES aims (id)
);