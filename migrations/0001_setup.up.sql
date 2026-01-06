CREATE TABLE IF NOT EXISTS descriptions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    description VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS categories (
    id INT AUTO_INCREMENT PRIMARY KEY,
    category VARCHAR(255)
);

CREATE TABLE IF NOT EXISTS entries (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    quantity INT NOT NULL,
    category INT NOT NULL,
    CONSTRAINT fk_entries_category
        FOREIGN KEY (category) REFERENCES categories(id),
    description INT,
    CONSTRAINT fk_entries_description
        FOREIGN KEY (description) REFERENCES descriptions(id)
) ENGINE=InnoDB;

