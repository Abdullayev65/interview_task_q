package moduls

type User struct {
	ID       int
	Name     string
	Password string
	UserName string `gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
}
type Post struct {
	ID     int
	UserID int
	User   *User
	Title  string
	Text   string
}
type Comment struct {
	ID     int
	UserID int
	User   *User
	Post   *Post
	PostID int
	Text   string
}
type Like struct {
	ID     int
	UserID int `gorm:"index:like_index"`
	User   *User
	// type can be "post" or "comment"
	Type string `gorm:"index:like_index"`
	// id of post or comment
	LikedID int `gorm:"index:like_index"`
}

func NewLike(UserID int, Type string, LikedID int) (*Like, error) {
	if Type == "post" || Type == "comment" {
		return &Like{UserID: UserID, Type: Type, LikedID: LikedID}, nil
	}
	return nil, nil
}
