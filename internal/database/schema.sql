
/* create a table called "users" */
CREATE TABLE IF NOT EXISTS  users (
    id INTEGER PRIMARY KEY AUTOINCREMENT ,
    userName VARCHAR(50) NOT NULL , 
    userEmail VARCHAR(100) NOT NULL ,
    userPassword VARCHAR(255) NOT NULL,
    token VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
);

/* create posts table*/
CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT ,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
    user_id INT ,
    title VARCHAR(255) NOT NULL ,
    content TEXT NOT NULL ,
    total_likes INT DEFAULT 0,
    total_dislikes INT DEFAULT 0,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

/* craete comments table*/
CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT ,
    post_id INT ,
    user_id INT ,
    total_likes INT DEFAULT 0,
    total_dislikes INT DEFAULT 0,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
    content TEXT NOT NULL ,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

/* create category table */
CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT ,
    category VARCHAR(255) NOT NULL,
    post_id INT ,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
);
