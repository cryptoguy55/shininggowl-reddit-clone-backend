package users

import (
	"errors"
	"first/common"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UsersRegister(router *gin.RouterGroup) {
	router.POST("/", UsersRegistration)
	router.POST("/login", UsersLogin)
	router.POST("/resetPassword", ResetLink)
	router.POST("/verifyEmail", VerifyAccount)
	router.POST("/verifyPassword", VerifyPassword)
	// router.GET("/password-reset", PasswordReset)
}

func UserRegister(router *gin.RouterGroup) {
	router.GET("/", UserRetrieve)
	router.PUT("/", UserUpdate)
	router.GET("/all", UserAll)
}

func ProfileRegister(router *gin.RouterGroup) {
	router.GET("/:username", ProfileRetrieve)
	router.POST("/:username/follow", ProfileFollow)
	router.DELETE("/:username/follow", ProfileUnfollow)
}

func ProfileRetrieve(c *gin.Context) {
	username := c.Param("username")
	userModel, err := FindOneUser(&UserModel{Username: username})
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("profile", errors.New("Invalid username")))
		return
	}
	profileSerializer := ProfileSerializer{c, userModel}
	c.JSON(http.StatusOK, gin.H{"profile": profileSerializer.Response()})
}

func ProfileFollow(c *gin.Context) {
	username := c.Param("username")
	userModel, err := FindOneUser(&UserModel{Username: username})
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("profile", errors.New("Invalid username")))
		return
	}
	myUserModel := c.MustGet("my_user_model").(UserModel)
	err = myUserModel.following(userModel)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	serializer := ProfileSerializer{c, userModel}
	c.JSON(http.StatusOK, gin.H{"profile": serializer.Response()})
}

func ProfileUnfollow(c *gin.Context) {
	username := c.Param("username")
	userModel, err := FindOneUser(&UserModel{Username: username})
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("profile", errors.New("Invalid username")))
		return
	}
	myUserModel := c.MustGet("my_user_model").(UserModel)

	err = myUserModel.unFollowing(userModel)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	serializer := ProfileSerializer{c, userModel}
	c.JSON(http.StatusOK, gin.H{"profile": serializer.Response()})
}

func UsersRegistration(c *gin.Context) {
	userModelValidator := NewUserModelValidator()
	if err := userModelValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	if file, err := c.FormFile("image"); err == nil {
		if errs := c.SaveUploadedFile(file, "./public/"+file.Filename); errs != nil {

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"message": "Unable to save the file",
			})
			return
		}
	}

	// single file
	if err := SaveOne(&userModelValidator.userModel); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	resetToken, _ := common.GenerateNonAuthToken(userModelValidator.userModel.Email)

	// Define email body
	link := "http://localhost:3000/verifyEmail?verify_token=" + resetToken
	body := "Here is your reset <a href='" + link + "'>link</a>"
	html := "<strong>" + body + "</strong>"

	// Initialize email sendout
	err := common.SendMail("Verify Account", body, userModelValidator.userModel.Email, html, userModelValidator.userModel.Username)
	if err != true {
		c.JSON(500, gin.H{"message": "An issue occured sending you an email"})
	}
	c.Set("my_user_model", userModelValidator.userModel)
	// serializer := UserSerializer{c}
	c.JSON(http.StatusCreated, gin.H{"user": "success"})

}

func UsersLogin(c *gin.Context) {
	loginValidator := NewLoginValidator()
	if err := loginValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("login", errors.New("Not empty email or password")))
		return
	}
	userModel, err := FindOneUser(&UserModel{Email: loginValidator.userModel.Email})

	if err != nil {
		c.JSON(http.StatusForbidden, common.NewError("login", errors.New("Not Registered email or invalid password")))
		return
	}
	if userModel.Active == false {
		c.JSON(http.StatusForbidden, common.NewError("login", errors.New("Please Email Verify.")))
		return
	}
	if userModel.checkPassword(loginValidator.User.Password) != nil {
		c.JSON(http.StatusForbidden, common.NewError("login", errors.New("Not Registered email or invalid password")))
		return
	}
	UpdateContextUserModel(c, userModel.ID)
	serializer := UserSerializer{c}
	c.JSON(http.StatusOK, gin.H{"user": serializer.Response()})
}

func UserRetrieve(c *gin.Context) {
	serializer := UserSerializer{c}
	c.JSON(http.StatusOK, gin.H{"user": serializer.Response()})
}

func UserAll(c *gin.Context) {
	userModels, err := getAllUsers()
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("Community", errors.New("Invalid param")))
		return
	}
	serializer := UsersSerializer{c, userModels}
	c.JSON(http.StatusOK, gin.H{"users": serializer.Response()})
}

func UserUpdate(c *gin.Context) {
	myUserModel := c.MustGet("my_user_model").(UserModel)
	userModelValidator := NewUserModelValidatorFillWith(myUserModel)
	if err := userModelValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		return
	}

	userModelValidator.userModel.ID = myUserModel.ID
	if err := myUserModel.Update(userModelValidator.userModel); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}
	UpdateContextUserModel(c, myUserModel.ID)
	serializer := UserSerializer{c}
	c.JSON(http.StatusOK, gin.H{"user": serializer.Response()})
}
func ResetLink(c *gin.Context) {

	var data ResendCommand
	// Ensure the user provides all values from the request.body
	if (c.BindJSON(&data)) != nil {
		// Return 400 status if they don't provide the email
		c.JSON(400, gin.H{"message": "Provided all fields"})
		c.Abort()
		return
	}

	// Fetch the account from the database based on the email
	// provided

	result, err := FindOneUser(&UserModel{Email: data.Email})
	// Return 404 status if an account was not found
	if result.Email == "" {
		c.JSON(404, gin.H{"message": "User account was not found"})
		c.Abort()
		return
	}

	// Return 500 status if something went wrong while fetching
	// account
	if err != nil {
		c.JSON(500, gin.H{"message": "Something wrong happened, try again later"})
		c.Abort()
		return
	}

	resetToken, _ := common.GenerateNonAuthToken(result.Email)

	// The link to be clicked in order to perform a password reset
	link := "http://localhost:3000/forgot-password?reset_token=" + resetToken
	// Define the body of the email
	body := "Here is your reset <a href='" + link + "'>link</a>"
	html := "<strong>" + body + "</strong>"

	// Initialize email sendout
	email := common.SendMail("Reset Password", body, result.Email, html, result.Username)

	// If email was sent, return 200 status code
	if email == true {
		c.JSON(200, gin.H{"message": "Please Check mail"})
		c.Abort()
		return
		// Return 500 status when something wrong happened
	} else {
		c.JSON(500, gin.H{"message": "An issue occured sending you an email"})
		c.Abort()
		return
	}
}

func VerifyPassword(c *gin.Context) {
	var data PasswordResetCommand

	if c.BindJSON(&data) != nil {
		c.JSON(406, gin.H{"message": "Provide relevant fields"})
		c.Abort()
		return
	}

	resetToken, _ := c.GetQuery("reset_token")
	fmt.Println(resetToken)
	// Decode the token
	userID, err := common.DecodeNonAuthToken(resetToken)
	if err != nil {
		// Return response when we get an error while fetching user
		c.JSON(500, gin.H{"message": "Invaild token"})
		c.Abort()
		return
	}
	result, err1 := FindOneUser(&UserModel{Email: userID})

	if err1 != nil {
		// Return response when we get an error while fetching user
		c.JSON(500, gin.H{"message": "Something wrong happened, try again later"})
		c.Abort()
		return
	}

	if result.Email == "" {
		c.JSON(404, gin.H{"message": "User account was not found"})
		c.Abort()
		return
	}

	// Update user account

	// _err := result.Update(UserModel{Password: true})

	// if _err != nil {
	// 	// Return response if we are not able to update user password
	// 	c.JSON(500, gin.H{"message": "Something happened while verifying you account, try again"})
	// 	c.Abort()
	// 	return
	// }

	c.JSON(201, gin.H{"message": "Password has been updated successfully, Please log in"})
}
func VerifyAccount(c *gin.Context) {
	// Get token from link query
	var data Token
	// Ensure the user provides all values from the request.body
	if (c.BindJSON(&data)) != nil {
		// Return 400 status if they don't provide the email
		c.JSON(400, gin.H{"message": "Token is needed"})
		c.Abort()
		return
	}
	verifyToken := data.Token

	// Decode verify token

	userID, err := common.DecodeNonAuthToken(verifyToken)
	if err != nil {
		fmt.Println(err)
		// Return response when we get an error while fetching user
		c.JSON(500, gin.H{"message": "Invaild token"})
		c.Abort()
		return
	}

	// Fetch user based on details from decoded token
	result, err1 := FindOneUser(&UserModel{Email: userID})

	if err1 != nil {
		// Return response when we get an error while fetching user
		c.JSON(500, gin.H{"message": "Something wrong happened, try again later"})
		c.Abort()
		return
	}

	if result.Email == "" {
		c.JSON(404, gin.H{"message": "User account was not found"})
		c.Abort()
		return
	}

	// Update user account

	_err := result.Update(UserModel{Active: true})

	if _err != nil {
		// Return response if we are not able to update user password
		c.JSON(500, gin.H{"message": "Something happened while verifying you account, try again"})
		c.Abort()
		return
	}

	c.JSON(201, gin.H{"message": "Account verified, log in"})
}
func RefreshToken(c *gin.Context) {
	// Get refresh token from header
	refreshToken := c.Request.Header["Refreshtoken"]

	// Check if refresh token was provided
	if refreshToken == nil {
		c.JSON(403, gin.H{"message": "No refresh token provided"})
		c.Abort()
		return
	}

	// Decode token to get data
	// id, err := common.DecodeNonAuthToken(refreshToken[0])

	// if err != nil {
	// 	c.JSON(500, gin.H{"message": "Problem refreshing your session"})
	// 	c.Abort()
	// 	return
	// }

	// Create new token
	// accessToken, _refreshToken, _err := common.GenRefreshToken(id)

	// if _err != nil {
	// 	c.JSON(500, gin.H{"message": "Problem creating new session"})
	// 	c.Abort()
	// 	return
	// }

	// c.JSON(200, gin.H{"message": "Log in success", "token": accessToken, "refresh_token": _refreshToken})
}
