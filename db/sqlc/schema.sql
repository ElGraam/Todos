-- schema.sql
CREATE TABLE todos (
    id INT AUTO_INCREMENT PRIMARY KEY,
    body TEXT NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
