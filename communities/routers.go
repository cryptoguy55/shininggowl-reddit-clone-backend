package communities

import (
	"errors"
	"first/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CommunitiesRegister(router *gin.RouterGroup) {
	router.POST("/", CommunityCreate)

}

func CommunitiesAnonymousRegister(router *gin.RouterGroup) {
	router.GET("/", CommunityList)
}

func CommunityCreate(c *gin.Context) {
	communityModelValidator := NewCommunityModelValidator()
	if err := communityModelValidator.Bind(c); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewValidatorError(err))
		return
	}
	//fmt.Println(articleModelValidator.articleModel.Author.UserModel)

	if err := SaveOne(&communityModelValidator.communityModel); err != nil {
		c.JSON(http.StatusUnprocessableEntity, common.NewError("database", err))
		return
	}

	serializer := CommunitySerializer{c, communityModelValidator.communityModel}
	c.JSON(http.StatusCreated, gin.H{"article": serializer.Response()})
}

func CommunityList(c *gin.Context) {
	communityModels, err := getAllCommunities()
	if err != nil {
		c.JSON(http.StatusNotFound, common.NewError("Community", errors.New("Invalid param")))
		return
	}
	serializer := CommunitiesSerializer{c, communityModels}
	c.JSON(http.StatusOK, gin.H{"communities": serializer.Response()})
}
