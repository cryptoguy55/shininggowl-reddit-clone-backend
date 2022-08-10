package communities

import (
	"first/common"
	"first/users"
)

func getAllCommunities() ([]users.CommunityModel, error) {
	db := common.GetDB()
	var models []users.CommunityModel
	err := db.Find(&models).Error
	return models, err
}
func SaveOne(data interface{}) error {
	db := common.GetDB()
	err := db.Save(data).Error
	return err
}
