package models

type User struct {
	CurrentPath string
	IsLoged   bool
	UserName  string
	UserEmail string
	// Profile   string // about imges we can store them in databse as blob dont wory
}

type FormsData struct {
	UserNameInput string 
	UserEmailInput string 
	UserPasswordInput string 
	FormErrors
}

type FormErrors struct {
	FormError string
	InvalidUserName string
	InvalidEmail string
	InvalidPassword string
}

type Data struct {
	User        User
	Posts       []Post
	Categories  []Categorie
	FormsData
}

type Post struct { /// after use your own envpreption
	PostCreator                                              string
	PostCreatedAt                                            string
	PostTitle                                                string
	PostContent                                              string
	TotalLikes, TotalDeslikes, TotalComments, PostId, UserID int
	Categories                                               []Categorie
}
type Categorie struct {
	CatergoryName string
}

type Comment struct {
	CommentCretor    string
	CommentCreatedAt string
	CommentContent   string
	CommentId        int
	UserId           int
}
