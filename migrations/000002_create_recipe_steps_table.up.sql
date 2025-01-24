CREATE TABLE IF NOT EXISTS recipe_steps
(
    id          INT AUTO_INCREMENT PRIMARY KEY,
    recipe_id   INT                                 NOT NULL,
    step_number INT                                 NOT NULL,
    description TEXT                                NOT NULL,
    FOREIGN KEY (recipe_id) REFERENCES recipes (id)
);
