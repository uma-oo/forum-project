package models

type User struct {
	CurrentPath string
	IsLoged   bool
	UserName  string
	UserEmail string
	// Profile   string // about imges we can store them in databse as blob dont wory
}

type Data struct {
	User        User
	Posts       []Post
	Categories  []Categorie
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
