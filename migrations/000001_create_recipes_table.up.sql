CREATE TABLE IF NOT EXISTS recipes
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    uuid        CHAR(36)                          NOT NULL UNIQUE,
    name        VARCHAR(255)                      NOT NULL,
    description TEXT,
    author      VARCHAR(100),
    dough       JSON      DEFAULT (JSON_OBJECT()) NOT NULL,
    topping     JSON      DEFAULT (JSON_OBJECT()) NOT NULL,
    steps       JSON      DEFAULT (JSON_ARRAY())  NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_uuid (uuid)
);
