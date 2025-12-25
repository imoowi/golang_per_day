package repos

import (
	"golang_per_day_24/internal/components"
	"golang_per_day_24/internal/interfaces"
	"golang_per_day_24/internal/models"
)

var User *UserRepo

type UserRepo struct {
	interfaces.Repo[*models.User]
}

func NewUserRepo() {
	User = &UserRepo{
		Repo: *interfaces.NewRepo[*models.User](components.DB),
	}
}

func init() {
	RegisterRepos(NewUserRepo)
}

// func (r *UserRepo) GetById(uid uint) (user models.User, err error) {
// 	err = r.DB.Where(`id=?`, uid).First(&user).Error
// 	if err != nil {
// 		return
// 	}
// 	return
// }
