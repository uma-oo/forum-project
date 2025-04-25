 # Web Forum Project

## Overview
This project is a web forum that allows users to communicate, associate categories to posts, like and dislike posts and comments,  filter posts and greate post .

## Features

- User registration and login

- Post creation and commenting

- Category association with posts

- Liking and disliking posts and comments

- Filtering posts by categories, created posts, and liked posts

 ## Usage
Clone the repository and then ./runner.sh will start the server for you.

Change into the project directory: cd web-forum

Build the Docker image: docker build -t web-forum .

Run the Docker container: docker run -p 8081:8080 web-forum

Open a web browser and navigate to http://localhost:8080

## Dependencies
This project was built using the following technologies:

- Go 1.17 or higher

- uuid 1.6 or higher

- bcrypt 0.32 or higher

- SQLite 3.36 or higher

- Docker 20.10 or higher

## Authors

- **Hassan El Ouazizi**     [GitHub Profile](https://github.com/helouazizi)
- **Ismail Haji**           [GitHub Profile](https://github.com/hajji-Ismail)
- **Ayoub Nachti**          [GitHub Profile](https://github.com/DarkMethoss)
- **Oumayma EL-FAHSI**      [GitHub Profile](https://github.com/uma-oo)
