package service

import (
	"errors"
	"github.com/Abdullayev65/interview_task_q/internal/app/endpoint"
	"github.com/Abdullayev65/interview_task_q/internal/app/moduls"
	"github.com/Abdullayev65/interview_task_q/internal/app/utill"
	"gorm.io/gorm"
)

type Service struct {
	DB       *gorm.DB
	TokenJWT *utill.TokenJWT
}

func New(DB *gorm.DB, TokenJWT *utill.TokenJWT) *Service {
	return &Service{DB: DB, TokenJWT: TokenJWT}
}

func (s *Service) SignUp(sign endpoint.Sign) (*moduls.User, error) {
	exists := s.existsBYUserName(sign.UserName)
	if exists {
		return nil, errors.New("username exists")
	}
	user := moduls.User{UserName: sign.UserName, Password: sign.Password, Name: sign.Name}
	s.DB.Create(&user)
	return &user, nil
}

func (s *Service) LogIn(sign endpoint.Sign) (string, error) {
	exists := s.existsBYUserName(sign.UserName)
	if !exists {
		return "", errors.New("username or password wrong")
	}
	var userId int
	s.DB.Model(&moduls.User{}).
		Select("id").
		Where("user_name = ?", sign.UserName).
		Find(&userId)
	token, err := s.TokenJWT.GenerateToken(userId)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *Service) AddPost(post *moduls.Post) error {
	err := s.DB.Create(post).Error
	if err != nil {
		return err
	}
	return nil
}
func (s *Service) PostsOfUser(userId int) []moduls.Post {
	var posts []moduls.Post
	s.DB.Model(&moduls.Post{}).
		Where("user_id = ?", userId).
		Find(&posts)
	return posts
}
func (s *Service) Like(like *moduls.Like) error {
	var exists bool
	s.DB.Model(&moduls.Like{}).
		Select("count(*) > 0").
		Where("user_id = ? AND liked_id = ? AND type = ?", like.UserID, like.LikedID, like.Type).
		Find(&exists)
	if exists {
		return errors.New("already liked")
	}
	{ // chack type and likeId are correct
		var module interface{}
		switch like.Type {
		case "post":
			module = &moduls.Post{}
		case "comment":
			module = &moduls.Comment{}
		default:
			return errors.New("type should be post or comment")
		}
		exists := s.exists(module, " id = ?", like.LikedID)
		if !exists {
			return errors.New(like.Type + " not exists by likedId")
		}
	}
	s.DB.Create(like)
	return nil
}

func (s *Service) UserIdFromToken(tokenStr string) (int, error) {
	return s.TokenJWT.ParseToken(tokenStr)
}

func (s *Service) AddComment(comment *moduls.Comment) error {
	exists := s.existsById(moduls.Post{}, comment.PostID)
	if !exists {
		return errors.New("post does not exist by id")
	}
	s.DB.Create(comment)
	return nil
}

func (s *Service) CommentsByPostId(postId int) ([]moduls.Comment, error) {
	exists := s.existsById(moduls.Post{}, postId)
	if !exists {
		return nil, errors.New("post does not exist by id")
	}
	var comments []moduls.Comment
	s.DB.Where("post_id = ?", postId).
		Find(&comments)
	return comments, nil
}

func (s *Service) existsBYUserName(userName string) bool {
	var exists bool
	s.DB.Model(&moduls.User{}).
		Select("count(*) > 0").
		Where("user_name = ?", userName).
		Find(&exists)
	return exists
}
func (s *Service) exists(module interface{}, whereQuery string, args ...interface{}) bool {
	var exists bool
	s.DB.Model(module).
		Select("count(*) > 0").
		Where(whereQuery, args...).
		Find(&exists)
	return exists
}
func (s *Service) existsById(module interface{}, id interface{}) bool {
	return s.exists(module, "id = ?", id)
}
