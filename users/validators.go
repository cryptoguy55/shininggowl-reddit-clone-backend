package users

import (
	"first/common"
	"mime/multipart"
	"os"

	"github.com/gin-gonic/gin"
)

// *ModelValidator containing two parts:
// - Validator: write the form/json checking rule according to the doc https://github.com/go-playground/validator
// - DataModel: fill with data from Validator after invoking common.Bind(c, self)
// Then, you can just call model.save() after the data is ready in DataModel.
type UserModelValidator struct {
	User struct {
		Username string                `form:"username" binding:"required"`
		Email    string                `form:"email" binding:"required" `
		Password string                `form:"password" binding:"required" `
		Bio      string                `form:"bio"`
		Image    *multipart.FileHeader `form:"image" binding:"required"`
	} `json:"user"`
	userModel UserModel `json:"-"`
}

// There are some difference when you create or update a model, you need to fill the DataModel before
// update so that you can use your origin data to cheat the validator.
// BTW, you can put your general binding logic here such as setting password.
func (self *UserModelValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, self)
	if err != nil {
		return err
	}
	self.userModel.Username = self.User.Username
	self.userModel.Email = self.User.Email
	self.userModel.Bio = self.User.Bio
	self.userModel.Image = self.User.Image.Filename
	if self.User.Password != os.Getenv("NBRandomPassword") {
		self.userModel.setPassword(self.User.Password)
	}

	// if self.User.Image != nil {
	// }
	return nil
}

// You can put the default value of a Validator here
func NewUserModelValidator() UserModelValidator {
	userModelValidator := UserModelValidator{}
	//userModelValidator.User.Email ="w@g.cn"
	return userModelValidator
}

func NewUserModelValidatorFillWith(userModel UserModel) UserModelValidator {
	userModelValidator := NewUserModelValidator()
	userModelValidator.User.Username = userModel.Username
	userModelValidator.User.Email = userModel.Email
	userModelValidator.User.Bio = userModel.Bio
	userModelValidator.User.Password = os.Getenv("NBRandomPassword")

	// if userModel.Image != nil {
	// 	userModelValidator.User.Image = *userModel.Image
	// }
	return userModelValidator
}

type LoginValidator struct {
	User struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required" validate:"min=8, max=255"`
	} `json:"user"`
	userModel UserModel `json:"-"`
}
type ResendCommand struct {
	// We only need the email to initialize an email sendout
	Email string `json:"email" binding:"required"`
}
type Token struct {
	// We only need the email to initialize an email sendout
	Token string `json:"token" binding:"required"`
}

type PasswordResetCommand struct {
	// We only need the email to initialize an email sendout
	Password string `json:"password" binding:"required"`
}

func (self *LoginValidator) Bind(c *gin.Context) error {
	err := common.Bind(c, self)
	if err != nil {
		return err
	}

	self.userModel.Email = self.User.Email
	return nil
}

// You can put the default value of a Validator here
func NewLoginValidator() LoginValidator {
	loginValidator := LoginValidator{}
	return loginValidator
}
