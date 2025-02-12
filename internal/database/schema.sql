/* create a table called "users" */
/* Expiration date to be added and also */
CREATE TABLE
    IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        userName VARCHAR(50) NOT NULL,
        userEmail VARCHAR(100) NOT NULL,
        userPassword VARCHAR(100) NOT NULL,
        token VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, /*for the user*/   
        token_created_at DATETIME ,
        expiration_date DATETIME
    );

-- table of session ??
/* create posts table*/
CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id INT,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    total_likes INT DEFAULT 0 CHECK (total_likes >= 0),
    total_dislikes INT DEFAULT 0 CHECK (total_dislikes >= 0),
    total_comments INT DEFAULT 0 CHECK (total_comments >= 0),  
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE 
);


/* create category table */
CREATE TABLE
    IF NOT EXISTS post_categories (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        category VARCHAR(255) NOT NULL,
        post_id INT,
        FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE
    );

/* lwts drop the categories table*/
DROP TABLE IF EXISTS categories;

CREATE TABLE
    IF NOT EXISTS categories (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        category TEXT
    );

INSERT INTO
    categories (category)
VALUES
    ('Technology'),
    ('Science'),
    ('Health'),
    ('Lifestyle'),
    ('Education'),
    ('Gaming'),
    ('Business');

CREATE TABLE IF NOT EXISTS post_reaction (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    reaction_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (post_id) REFERENCES posts (id),
    FOREIGN KEY (reaction_id) REFERENCES reaction (id), 
    CONSTRAINT unique_columns UNIQUE (user_id, post_id)
);

--CREATE INDEX IF NOT EXISTS idx_user_post ON likes (user_id, post_id);
/* craete comments table*/
CREATE TABLE
    IF NOT EXISTS comments (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        post_id INT,
        user_id INT,
        total_likes INT DEFAULT 0,
        total_dislikes INT DEFAULT 0,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        content TEXT NOT NULL,
        FOREIGN KEY (post_id) REFERENCES posts (id) ON DELETE CASCADE,
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
    );

DROP TABLE IF EXISTS reactions;

CREATE TABLE IF NOT EXISTS reactions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    reaction INTEGER NOT NULL

);

INSERT INTO 
    reactions (reaction)
VALUES (0),(1),(-1);


CREATE TABLE
    IF NOT EXISTS comment_reactions (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        comment_id INTEGER NOT NULL,
        reaction_id INTEGER NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES users (id),
        FOREIGN KEY (comment_id) REFERENCES comments (id),
        FOREIGN KEY (reaction_id) REFERENCES reaction (id), 
         CONSTRAINT unique_columns UNIQUE (user_id, comment_id)
    );

