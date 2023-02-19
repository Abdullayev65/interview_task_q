package endpoint

import (
	"fmt"
	"github.com/Abdullayev65/interview_task_q/internal/app/moduls"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
)

type Service interface {
	SignUp(sign Sign) (*moduls.User, error)
	LogIn(sign Sign) (string, error)
	AddPost(post *moduls.Post) error
	PostsOfUser(userId int) []moduls.Post
	Like(like *moduls.Like) error
	UserIdFromToken(tokenStr string) (int, error)
	AddComment(comment *moduls.Comment) error
	CommentsByPostId(postId int) ([]moduls.Comment, error)
}

type Endpoint struct {
	s Service
}

func New(s Service) *Endpoint {
	return &Endpoint{
		s: s,
	}
}

//	func (e *Endpoint) Status(c echo.Context) error {
//		d := e.s.DaysLeft()
//
//		s := fmt.Sprintf("Days left: %d", d)
//
//		err := c.String(http.StatusOK, s)
//		if err != nil {
//			return err
//		}
//
//		return nil
//	}
func (e *Endpoint) SignUp(c echo.Context) error {
	var sign Sign
	err := c.Bind(&sign)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return err
	}
	user, err := e.s.SignUp(sign)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return nil
	}
	c.JSON(http.StatusOK, &user)
	return nil
}
func (e *Endpoint) LogIn(c echo.Context) error {
	var sign Sign
	err := c.Bind(&sign)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(":57")
		c.JSON(http.StatusBadRequest, err.Error())
		return nil
	}
	token, err := e.s.LogIn(sign)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return nil
	}
	c.JSON(http.StatusOK, &token)
	return nil
}
func (e *Endpoint) AddPost(c echo.Context) error {
	userId := c.Get("userId").(int)
	var postInput PostInput
	err := c.Bind(&postInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return err
	}
	post := postInput.toModule(userId)
	err = e.s.AddPost(&post)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return nil
	}
	c.JSON(http.StatusOK, newPostOutput(&post))
	return nil
}
func (e *Endpoint) Like(c echo.Context) error {
	// todo
	likeId, err := strconv.Atoi(c.QueryParam("id"))
	if err != nil {
		return err
	}
	typeStr := c.QueryParam("type")
	userId := c.Get("userId").(int)
	like, err := moduls.NewLike(userId, typeStr, likeId)
	if err != nil {
		return err
	}
	err = e.s.Like(like)
	if err != nil {
		return err
	}
	c.JSON(http.StatusOK, "liked 👍")
	return nil
}
func (e *Endpoint) PostsOfUser(c echo.Context) error {
	userId, err := strconv.Atoi(c.QueryParam("userId"))
	if err != nil {
		return err
	}
	postsOfUser := e.s.PostsOfUser(userId)
	c.JSON(http.StatusOK, &postsOfUser)
	return nil
}
func (e *Endpoint) SetUserIdForMW(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get("Authorization")
		if header == "" {
			c.String(http.StatusUnauthorized, "Authorization header is empty")
			return nil
		}
		fields := strings.Fields(header)
		if len(fields) != 2 || fields[0] != "Bearer" {
			c.String(http.StatusUnauthorized, "Authorization header is invalid")
			return nil
		}
		userId, err := e.s.UserIdFromToken(fields[1])
		if err != nil {
			c.String(http.StatusUnauthorized, err.Error())
			return nil
		}
		c.Set("userId", userId)

		return next(c)
	}
}

func (e *Endpoint) AddComment(c echo.Context) error {
	userId := c.Get("userId").(int)
	var commentInput CommentInput
	c.Bind(&commentInput)
	comment := commentInput.toModule(userId)
	err := e.s.AddComment(&comment)
	if err != nil {
		return err
	}
	commentOutput := newCommentOutPut(comment)
	c.JSON(http.StatusOK, &commentOutput)
	return nil
}

func (e *Endpoint) CommentsByPostId(c echo.Context) error {
	postId, err := strconv.Atoi(c.QueryParam("postId"))
	if err != nil {
		return err
	}
	comments, err := e.s.CommentsByPostId(postId)
	if err != nil {
		return err
	}
	var commentOutputSlice []CommentOutput
	for _, comment := range comments {
		commentOutputSlice = append(commentOutputSlice, newCommentOutPut(comment))
	}
	c.JSON(http.StatusOK, &commentOutputSlice)
	return nil
}
