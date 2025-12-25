package repos

import (
	"context"
	"gindemo2/config"
	"gindemo2/models"
)

var (
	User *UserRepo
	Ctx  context.Context = context.Background()
)

type UserRepo struct {
}

func init() {
	User = &UserRepo{}
}

func (s *UserRepo) GetById(uid uint) (user models.User, err error) {
	err = config.DB.Where(`id=?`, uid).First(&user).Error
	if err != nil {
		return
	}
	return
}
