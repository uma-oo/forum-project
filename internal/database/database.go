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

type Database struct {
	*sql.DB
}

// Open forum database
func NewDatabase() (*Database, error) {
	dbPath := os.Getenv("DB_PATH")
	// fmt.Printf("dbPath: %v\n", dbPath)
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		logger.LogWithDetails(err)
		return nil, err
	}
	return &Database{db}, nil
}

func Create_database() {
	db, err := NewDatabase()
	if err != nil {
		log.Fatal(err)
	}

	// lets open the schema file to execute the sql commands inside it
	shema_path := os.Getenv("SCHEMA_PATH")
	schema, err := os.Open(shema_path)
	if err != nil {
		logger.LogWithDetails(err)
		log.Fatal(err)
	}

	defer schema.Close()
	// now lets read the schema file using the bufio package
	scanner := bufio.NewScanner(schema)
	var sql_command string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "--") || strings.HasPrefix(line, "/*") || line == "" {
			continue
		}
		sql_command += line + " "
		// lets execute the sql command
		if strings.HasSuffix(sql_command, "; ") {
			_, err = db.Exec(sql_command)
			if err != nil {
				logger.LogWithDetails(err)
				log.Fatal(err)
			}
			// free up the sql command
			sql_command = ""
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("data base creatd succesfully")
	Triggers()
}

func Fetch_Database(r *http.Request, query string, userid int, liked bool) (*models.Data, error) {
	var finalQuery string
	if userid > 0 && !liked {
		finalQuery = fmt.Sprintf("%s WHERE users.id = %d ORDER BY posts.created_at DESC;", query, userid)
	} else if userid > 0 && liked {
		finalQuery = fmt.Sprintf("%s WHERE  post_reaction.user_id = %d AND  post_reaction.reaction = 1", query, userid)
	} else { // all posts
		finalQuery = fmt.Sprintf("%s ORDER BY posts.created_at DESC", query)
	}

	db, err := NewDatabase()
	if err != nil {
		logger.LogWithDetails(err)
	}
	data := &models.Data{}
	stm, err := db.Prepare(finalQuery)
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
		stm, err := db.Prepare("SELECT userEmail FROM users WHERE userName = ? ")
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
			&post.PostId, &post.PostTitle, &post.PostContent, &post.TotalLikes, &post.TotalDeslikes, &post.TotalComments, &post.PostCreatedAt, &post.PostCreator, &post.UserID,
		)
		if err != nil {
			logger.LogWithDetails(err)
			return nil, err
		}
		// Fetch categories for the post
		query := "SELECT category FROM post_categories WHERE post_id = ?"
		stm, err := db.Prepare(query)
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
	query2 := `SELECT category FROM categories`
	stm, err = db.Prepare(query2)
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

func Triggers() error {
	db, err := NewDatabase()
	if err != nil {
		logger.LogWithDetails(err)
		return err
	}
	trigger_total_comments := `CREATE TRIGGER IF NOT EXISTS increment_total_comments 
			AFTER INSERT ON comments 
			FOR EACH ROW 
			BEGIN 
				UPDATE posts
				SET total_comments=total_comments+1 
				WHERE posts.id = NEW.post_id;
			END;`
	trigger_total_likes_comments_insert := `CREATE TRIGGER IF NOT EXISTS increment_or_decrement_total_likes_comments_insert
	   AFTER INSERT ON comment_reactions
		FOR EACH ROW
		BEGIN

		UPDATE comments
		SET total_likes = total_likes + 1
		WHERE comments.id = NEW.comment_id
		AND NEW.reaction_id = 1 ;
		END;`

	trigger_total_likes_comments_update := `CREATE TRIGGER IF NOT EXISTS increment_or_decrement_total_likes_comments_update
		AFTER UPDATE ON comment_reactions
		FOR EACH ROW
		BEGIN

		UPDATE comments
		SET total_likes = total_likes + 1
		WHERE comments.id = NEW.comment_id
		AND OLD.reaction_id=0
		AND NEW.reaction_id = 1 ;


		UPDATE comments
		SET 
        total_likes = total_likes + 1,
        total_dislikes = CASE 
            WHEN total_dislikes > 0 THEN total_dislikes - 1
            ELSE 0
        END
		WHERE comments.id = NEW.comment_id AND total_dislikes -1 >= 0
		AND OLD.reaction_id=-1
		AND NEW.reaction_id = 1 ;

		UPDATE comments
		SET total_likes = CASE 
        WHEN total_likes > 0 THEN total_likes - 1
        ELSE 0
        END
		WHERE comments.id = NEW.comment_id
		AND OLD.reaction_id=1
		AND NEW.reaction_id = 0;
		END;`

	trigger_total_dislikes_comments_insert := `CREATE TRIGGER IF NOT EXISTS increment_or_decrement_total_dislikes_comments_insert
	   AFTER INSERT ON comment_reactions
		FOR EACH ROW
		BEGIN

		UPDATE comments
		SET total_dislikes = total_dislikes + 1
		WHERE comments.id = NEW.comment_id
		AND NEW.reaction_id = -1 ;
		END;`
	trigger_total_dislikes_comments_update := `CREATE TRIGGER IF NOT EXISTS increment_or_decrement_total_dislikes_comments_update
		AFTER UPDATE ON comment_reactions
		FOR EACH ROW
		BEGIN

		UPDATE comments
		SET total_dislikes = total_dislikes + 1
		WHERE comments.id = NEW.comment_id
		AND OLD.reaction_id=0 
		AND NEW.reaction_id = -1 ;

		UPDATE comments
		SET 
        total_dislikes = total_dislikes + 1,
        total_likes = CASE 
            WHEN total_likes > 0 THEN total_likes - 1
            ELSE 0
        END
		WHERE comments.id = NEW.comment_id
		AND OLD.reaction_id=1
		AND NEW.reaction_id = -1 ;

		UPDATE comments
		SET total_dislikes = CASE 
        WHEN total_dislikes > 0 THEN total_dislikes - 1
        ELSE 0
        END
		WHERE comments.id = NEW.comment_id
		AND OLD.reaction_id=-1
		AND NEW.reaction_id = 0;
		END;`
		trigger_total_likes_posts_insert := `CREATE TRIGGER IF NOT EXISTS increment_or_decrement_total_likes_posts_insert
		AFTER INSERT ON post_reaction
		 FOR EACH ROW
		 BEGIN
 
		 UPDATE posts
		 SET total_likes = total_likes + 1
		 WHERE posts.id = NEW.post_id 
		 AND NEW.reaction_id = 1 ;
		 END;`
		 trigger_total_likes_posts_update := `CREATE TRIGGER IF NOT EXISTS increment_or_decrement_total_likes_posts_update
		 AFTER UPDATE ON post_reaction
		 FOR EACH ROW
		 BEGIN
 
		 UPDATE posts
		 SET total_likes = total_likes + 1
		 WHERE posts.id = NEW.post_id
		 AND OLD.reaction_id=0
		 AND NEW.reaction_id = 1 ;
 
 
		 UPDATE posts
		 SET 
		 total_likes = total_likes + 1,
		 total_dislikes = CASE 
			 WHEN total_dislikes > 0 THEN total_dislikes - 1
			 ELSE 0
		 END
		 WHERE posts.id = NEW.post_id AND total_dislikes -1 >= 0
		 AND OLD.reaction_id=-1
		 AND NEW.reaction_id = 1 ;
 
		 UPDATE posts
		 SET total_likes = CASE 
		 WHEN total_likes > 0 THEN total_likes - 1
		 ELSE 0
		 END
		 WHERE posts.id = NEW.post_id
		 AND OLD.reaction_id=1
		 AND NEW.reaction_id = 0;
		 END;`
		 trigger_total_dislikes_posts_insert := `CREATE TRIGGER IF NOT EXISTS increment_or_decrement_total_dislikes_posts_insert
		 AFTER INSERT ON post_reaction
		  FOR EACH ROW
		  BEGIN
  
		  UPDATE posts
		  SET total_dislikes = total_dislikes + 1
		  WHERE posts.id = NEW.post_id
		  AND NEW.reaction_id = -1 ;
		  END;`
		  trigger_total_dislikes_posts_update := `CREATE TRIGGER IF NOT EXISTS increment_or_decrement_total_dislikes_posts_update
		  AFTER UPDATE ON post_reaction
		  FOR EACH ROW
		  BEGIN
  
		  UPDATE posts
		  SET total_dislikes = total_dislikes + 1
		  WHERE posts.id = NEW.post_id
		  AND OLD.reaction_id=0 
		  AND NEW.reaction_id = -1 ;
  
		  UPDATE posts
		  SET 
		  total_dislikes = total_dislikes + 1,
		  total_likes = CASE 
			  WHEN total_likes > 0 THEN total_likes - 1
			  ELSE 0
		  END
		  WHERE posts.id = NEW.post_id
		  AND OLD.reaction_id=1
		  AND NEW.reaction_id = -1 ;
  
		  UPDATE posts
		  SET total_dislikes = CASE 
		  WHEN total_dislikes > 0 THEN total_dislikes - 1
		  ELSE 0
		  END
		  WHERE posts.id = NEW.post_id
		  AND OLD.reaction_id=-1
		  AND NEW.reaction_id = 0;
		  END;`
	triggers := []string{
		trigger_total_comments, trigger_total_likes_comments_insert,
		trigger_total_likes_comments_update, trigger_total_dislikes_comments_insert,
		trigger_total_dislikes_comments_update,trigger_total_likes_posts_insert, trigger_total_likes_posts_update, trigger_total_dislikes_posts_insert,trigger_total_dislikes_posts_update,
	}
	for _, query := range triggers {
		statement, err := db.Prepare(query)
		if err != nil {
			logger.LogWithDetails(err)
			return err
		}
		statement.Exec()
	}
	return nil
}
