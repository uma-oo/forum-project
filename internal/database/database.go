// internal/database/database.go
// internal/database/db.go
package database

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"forum/internal/models"

	_ "github.com/mattn/go-sqlite3"
)

var Database *sql.DB

func Create_database() {
	var err error
	Database, err = sql.Open("sqlite3", "./internal/database/forum.db")
	if err != nil {
		log.Fatal(err)
	}

	// lets open the schema file to execute the sql commands inside it
	schema, err := os.Open("./internal/database/schema.sql")
	if err != nil {
		log.Fatal(err)
	}
	defer schema.Close()

	// now lets read the schema file using the bufio package
	scanner := bufio.NewScanner(schema)
	var sql_command string
	lineIndex := 0
	for scanner.Scan() {

		lineIndex++
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "--") || strings.HasPrefix(line, "/*") || line == "" {
			continue
		}
		sql_command += line + " "
		// lets execute the sql command
		if strings.HasSuffix(sql_command, "; ") {
			_, err = Database.Exec(sql_command)
			if err != nil {
				log.Fatal(err, " line hna :", lineIndex)
			}
			// free up the sql command
			sql_command = ""
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("data base creatd succesfully")
}

func Fetch_Database(r *http.Request, query string, userid int) *models.Data {
	rows, err := Database.Query(query, userid)
	if err != nil {
		fmt.Println("Error executing query:", err)
		log.Fatal("Error executing query:", err)
	}
	defer rows.Close()
	// lets iterate over rows and store them in our models
	data := &models.Data{}

	// lets check if the user have a token

	if t, err := r.Cookie("token"); err == nil {
		if t.Value != "" {
			data.User.IsLoged = true
		}
	}
	// lets extract his username
	userName := r.FormValue("userName")
	Email := r.FormValue("userEmail")
	if Email == "" {
		// fmt.Println("email empty")
		Database.QueryRow("SELECT userEmail FROM users WHERE userName = $1 ", userName).Scan(&Email)
	}

	data.User.UserName = userName
	data.User.UserEmail = Email
	// fmt.Println(data.Userr.UserName, data.Userr.UserEmail, data.Userr.IsLoged, " after in data base ")

	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(
			&post.PostId, &post.PostTitle, &post.PostContent, &post.TotalLikes, &post.TotalDeslikes, &post.PostCreatedAt, &post.PostCreator, &post.UserID,
		)
		if err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		// Fetch categories for the post
		query := "SELECT category FROM categories WHERE post_id = ?"
		rows2, err := Database.Query(query, post.PostId)
		if err != nil {
			log.Fatal(err)
		}
		defer rows2.Close()
		for rows2.Next() {
			categ := &models.Categorie{}
			err := rows2.Scan(&categ.CatergoryName)
			if err != nil {
				log.Fatal(err)
			}
			post.Categories = append(post.Categories, *categ)
		}
		data.Posts = append(data.Posts, *post)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	/// lets fetch cetegories
	query2 := `SELECT category FROM stoke_categories`
	rows, err = Database.Query(query2)
	if err != nil {
		fmt.Println("Error executing query:", err)
		log.Fatal("Error executing query:", err)
	}
	defer rows.Close()
	for rows.Next() {
		category := &models.Categorie{}
		err := rows.Scan(&category.CatergoryName)
		if err != nil {
			log.Fatalf("Failed to scan row: %v", err)
		}
		data.Categories = append(data.Categories, *category)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	return data
}
