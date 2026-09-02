package utils

import (
	"errors"
	"time"

	"github.com/AlifRJ/Go_Auth/app/model"
	"github.com/AlifRJ/Go_Auth/app/payload"
)

func GetUserByID(users []model.User, id uint) (model.User, error) {
	for _, user := range users {
		if user.ID == id {
			return user, nil
		}
	}
	return model.User{}, errors.New("user not found")
}

func UpdateUserByID(users []model.User, id uint, body payload.UpdateUserPayload, newPassword string) (error) {
	for i := range users {
		if users[i].ID == id {
			if body.Password != nil {
				users[i].Password = newPassword
			}
			if body.Name != nil {
				users[i].Name = *body.Name
			}
			if body.Username != nil {
				users[i].Username = *body.Username
			}
			if body.Email != nil {
				users[i].Email = *body.Email
			}
			now := time.Now()
			users[i].Updated_at = &now
			return nil
		}
	}
	return errors.New("user not found")
}

func DeleteUserByID(users []model.User, id uint) (error) {
	for i := range users {
		if users[i].ID == id && users[i].Deleted_at == nil{
			now := time.Now()
			users[i].Deleted_at = &now
			return nil
		}
	}
	return errors.New("user not found")
}