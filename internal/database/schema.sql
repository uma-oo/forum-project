
/* create a table called "users" */
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT ,
    userName VARCHAR(50) NOT NULL ,
    userEmail VARCHAR(100) NOT NULL ,
    userPassword VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
);

/* create posts table*/
CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT ,
    user_id INT ,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP ,
    title VARCHAR(255) NOT NULL ,
    content TEXT NOT NULL ,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
-- edit the above query to include the following fields:
