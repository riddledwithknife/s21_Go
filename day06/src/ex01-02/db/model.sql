CREATE TABLE articles
(
    article_id   SERIAL PRIMARY KEY,
    title        VARCHAR(255),
    article_date DATE,
    article_text TEXT
)