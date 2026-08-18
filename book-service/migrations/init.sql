CREATE TABLE books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description VARCHAR(200) NOT NULL,
    author VARCHAR(100) NOT NULL /* масштабирование до REFERENCES authors(id) */
);