/* create a table called "users" */
CREATE TABLE
    IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        userName VARCHAR(50) NOT NULL,
        userEmail VARCHAR(100) NOT NULL,
        userPassword VARCHAR(255) NOT NULL,
        token VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

-- table of session ??
/* create posts table*/
CREATE TABLE
    IF NOT EXISTS posts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        user_id INT,
        title VARCHAR(255) NOT NULL,
        content TEXT NOT NULL,
        total_likes INT DEFAULT 0,
        total_dislikes INT DEFAULT 0,
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
    );

/* craete comments table*/
CREATE TABLE
    IF NOT EXISTS comments (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        post_id INT,
        user_id INT,
        total_likes INT DEFAULT 0,
        total_dislikes INT DEFAULT 0,
        TotalComments INT DEFAULT 0,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        content TEXT NOT NULL,
        FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
    );

/* create category table */
CREATE TABLE
    IF NOT EXISTS categories (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        category VARCHAR(255) NOT NULL,
        post_id INT,
        FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE
    );

/* lwts drop the categories table*/
DROP TABLE IF EXISTS stoke_categories;

CREATE TABLE
    IF NOT EXISTS stoke_categories (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        category TEXT,
        category_id INT
    );

INSERT INTO
    stoke_categories (category, category_id)
VALUES
    ('Technology', 1),
    ('Science', 2),
    ('Health', 3),
    ('Lifestyle', 4),
    ('Education', 5),
    ('Gaming', 6),
    ('Business', 7);

CREATE TABLE
    IF NOT EXISTS likes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        post_id INTEGER NOT NULL,
        reaction INTEGER NOT NULL CHECK (reaction IN (-1, 1)) DEFAULT -1,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users (id),
        FOREIGN KEY (post_id) REFERENCES posts (id)
    );

CREATE TABLE
    IF NOT EXISTS dislikes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        post_id INTEGER NOT NULL,
        reaction INTEGER NOT NULL CHECK (reaction IN (-1, 1)),
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users (id),
        FOREIGN KEY (post_id) REFERENCES posts (id)
    );