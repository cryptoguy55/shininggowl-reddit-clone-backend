package users

import (
	"github.com/gin-gonic/gin"

	"first/common"
)

type ProfileSerializer struct {
	C *gin.Context
	UserModel
}

// Declare your response schema here
type ProfileResponse struct {
	ID        uint   `json:"-"`
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Image     string `json:"image"`
	Following bool   `json:"following"`
}

// Put your response logic including wrap the userModel here.
func (self *ProfileSerializer) Response() ProfileResponse {
	myUserModel := self.C.MustGet("my_user_model").(UserModel)
	profile := ProfileResponse{
		ID:        self.ID,
		Username:  self.Username,
		Bio:       self.Bio,
		Image:     self.Image,
		Following: myUserModel.isFollowing(self.UserModel),
	}
	return profile
}

type UserSerializer struct {
	c *gin.Context
}

type UserResponse struct {
	Username     string `json:"username"`
	Id           uint   `json:"id"`
	Email        string `json:"email"`
	Bio          string `json:"bio"`
	Image        string `json:"image"`
	Token        string `json:"token"`
	RefreshToken string `json:"refreshtoken"`
}

func (self *UserSerializer) Response() UserResponse {
	myUserModel := self.c.MustGet("my_user_model").(UserModel)
	user := UserResponse{
		Username:     myUserModel.Username,
		Email:        myUserModel.Email,
		Bio:          myUserModel.Bio,
		Image:        myUserModel.Image,
		Token:        common.GenToken(myUserModel.ID),
		RefreshToken: common.GenRefreshToken(myUserModel.ID),
		Id:           myUserModel.ID,
	}
	return user
}

type CommonUserSerializer struct {
	c *gin.Context
	UserModel
}
type UsersSerializer struct {
	C     *gin.Context
	Users []UserModel
}
type AllResponse struct {
	Username string `json:"label"`
	Id       uint   `json:"value"`
	Image    string `json:"image"`
}

func (s *CommonUserSerializer) Response() AllResponse {
	user := AllResponse{
		Username: s.Username,
		Id:       s.ID,
		Image:    s.Image,
	}
	return user
}
func (s *UsersSerializer) Response() []AllResponse {
	response := []AllResponse{}
	for _, user := range s.Users {
		serializer := CommonUserSerializer{s.C, user}
		response = append(response, serializer.Response())
	}
	return response
}
