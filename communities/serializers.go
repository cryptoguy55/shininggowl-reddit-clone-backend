package communities

import (
	"first/users"

	"github.com/gin-gonic/gin"
)

type CommunitySerializer struct {
	C *gin.Context
	users.CommunityModel
}
type CommunityResponse struct {
	ID          uint   `json:"ID"`
	Name        string `json:"Name"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

type CommunitiesSerializer struct {
	C           *gin.Context
	Communities []users.CommunityModel
}

func (s *CommunitiesSerializer) Response() []CommunityResponse {
	response := []CommunityResponse{}
	for _, community := range s.Communities {
		serializer := CommunitySerializer{s.C, community}
		response = append(response, serializer.Response())
	}
	return response
}
func (s *CommunitySerializer) Response() CommunityResponse {
	response := CommunityResponse{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		Active:      true,
	}
	return response
}
