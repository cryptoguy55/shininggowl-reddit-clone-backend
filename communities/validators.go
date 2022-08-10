package communities

import (
	"first/common"
	"first/users"
	"github.com/gin-gonic/gin"
)

func NewCommunityModelValidator() CommunityModelValidator {
	return CommunityModelValidator{}
}

type CommunityModelValidator struct {
	Community struct {
		Name        string `form:"title" json:"name"`
		Description string `form:"description" json:"description"`
	} `json:"community"`
	communityModel users.CommunityModel `json:"-"`
}

func (s *CommunityModelValidator) Bind(c *gin.Context) error {

	err := common.Bind(c, s)
	if err != nil {
		return err
	}
	s.communityModel.Name = s.Community.Name
	s.communityModel.Description = s.Community.Description
	return nil
}
