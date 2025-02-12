package models

import "fmt"

type User struct {
	CurrentPath string
	IsLoged     bool
	UserName    string
	UserEmail   string
	UserId      string 
	// Profile   string // about imges we can store them in databse as blob dont wory
}

type FormsData struct {
	UserNameInput       string
	UserEmailInput      string
	UserPasswordInput   string
	PostGategoriesInput []string
	PostTitleInput      string
	PostContentInput    string
	FormErrors
}

type FormErrors struct {
	FormError             string
	InvalidUserName       string
	InvalidEmail          string
	InvalidPassword       string
	InvalidPostTitle      string
	InvalidPostContent    string
	InvalidPostCategories string
}

type Data struct {
	User       User
	Posts      []Post
	Categories []Categorie
	FormsData
}

type Categorie struct {
	CatergoryName string
}

type Comment struct {
	CommentId                 int
	UserId                    int
	PostId                    int
	CommentCreator            string
	CommentCreatedAt          string
	CommentContent            string
	TotalLikes, TotalDeslikes int
}

type Post struct { /// after use your own envpreption
	PostCreator                                              string
	PostCreatedAt                                            string
	PostTitle                                                string
	PostContent                                              string
	TotalLikes, TotalDeslikes, TotalComments, PostId, UserID int
	Categories                                               []Categorie
	Comments                                                 []Comment
}

// Functions helping us to debug

func (c *Comment) String() string {
	return fmt.Sprintf(`commentId %v , UserId %v, PostId %v CommentCreator %v CommentCreatedAt %v 
	CommentContent %v TotalLikes %v TotalDeslikes %v`, c.CommentId, c.UserId, c.PostId, c.CommentCreator, c.CommentCreatedAt, c.CommentContent, c.TotalLikes, c.TotalDeslikes)
}

func (p *Post) String() string {
	return fmt.Sprintf(`PostId %v , UserId %v, Post Creator %v
	 PostCreatedAt %v PostTitle %v 
	PostContent %v TotalLikes %v TotalDeslikes %v Total Comments %v
	Categories %v Comments %v`, p.PostId, p.UserID, p.PostCreator, p.PostCreatedAt, p.PostTitle, p.PostContent, p.TotalLikes, p.TotalDeslikes, p.TotalComments, p.Categories, p.Comments)
}

func (u *User) String() string {
	return fmt.Sprintf(`Path %v, isLogged %v, username %v, email %v`, u.CurrentPath, u.IsLoged, u.UserName, u.UserEmail)
}

type PageError struct {
	StatusCode   int
	ErrorMessage string
}

var (
	BadRequest          = PageError{400, "Bad Request"}
	PageNotFound        = PageError{404, "Page Not Found"}
	MethodNotAllowed    = PageError{405, "Method Not Allowed"}
	Unauthorized        = PageError{403, "Access Forbidden"}
	InternalServerError = PageError{500, "Internal server error"}
)
