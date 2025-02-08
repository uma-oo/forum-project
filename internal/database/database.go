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
	"forum/pkg/logger"

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

func Fetch_Database(r *http.Request, query string, userid int, liked bool) (*models.Data, error) {
	var finalQuery string
	if userid > 0 && !liked {
		finalQuery = fmt.Sprintf("%s WHERE users.id = %d ORDER  BY posts.created_at DESC;", query, userid)
	} else if userid > 0 && liked {
		finalQuery = fmt.Sprintf("%s WHERE  post_reaction.user_id = %d AND  post_reaction.reaction = 1", query, userid)
	} else {
		finalQuery = fmt.Sprintf("%s ORDER BY posts.created_at DESC", query)
	}

	stm, err := Database.Prepare(finalQuery)
	if err != nil {
		logger.LogWithDetails(err)
		return nil, err
	}
	rows, err := stm.Query()
	if err != nil {
		logger.LogWithDetails(err)
		return nil, err
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
		stm, err := Database.Prepare("SELECT userEmail FROM users WHERE userName = ? ")
		if err != nil {
			logger.LogWithDetails(err)
			return nil, err
		}
		stm.QueryRow(userName).Scan(&Email)
	}
	data.User.UserName = userName
	data.User.UserEmail = Email

	for rows.Next() {
		post := &models.Post{}
		err := rows.Scan(
			&post.PostId, &post.PostTitle, &post.PostContent, &post.TotalLikes, &post.TotalDeslikes, &post.PostCreatedAt, &post.PostCreator, &post.UserID,
		)
		if err != nil {
			logger.LogWithDetails(err)
			return nil, err
		}
		// Fetch categories for the post
		query := "SELECT category FROM categories WHERE post_id = ?"
		stm, err := Database.Prepare(query)
		if err != nil {
			logger.LogWithDetails(err)
			return nil, err
		}
		rows2, err := stm.Query(post.PostId)
		if err != nil {
			logger.LogWithDetails(err)
			return nil, err
		}
		defer rows2.Close()
		for rows2.Next() {
			categ := &models.Categorie{}
			err := rows2.Scan(&categ.CatergoryName)
			if err != nil {
				logger.LogWithDetails(err)
				return nil, err
			}
			post.Categories = append(post.Categories, *categ)
		}
		data.Posts = append(data.Posts, *post)
	}

	if err := rows.Err(); err != nil {
		logger.LogWithDetails(err)
		return nil, err
	}
	/// lets fetch cetegories
	query2 := `SELECT category FROM stoke_categories`
	stm, err = Database.Prepare(query2)
	if err != nil {
		logger.LogWithDetails(err)
		return nil, err
	}

	rows, err = stm.Query()
	if err != nil {
		logger.LogWithDetails(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		category := &models.Categorie{}
		err := rows.Scan(&category.CatergoryName)
		if err != nil {
			logger.LogWithDetails(err)
			return nil, err
		}
		data.Categories = append(data.Categories, *category)
	}
	if err := rows.Err(); err != nil {
		logger.LogWithDetails(err)
		return nil, err
	}

	return data, nil
}
