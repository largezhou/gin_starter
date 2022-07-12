package user

import (
	"context"

	"github.com/google/uuid"
	"github.com/largezhou/gin_starter/app/model"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	model.Model
	Uuid     string  `json:"uuid"`
	Username string  `json:"username"`
	Password string  `json:"-"`
	Email    *string `json:"email"` // 可为 null 的值，用指针配合 v.V 和 v.P 辅助方法 比用 sql.NullString 方便多了
}

func Create(ctx context.Context, user *User) (*User, error) {
	user.Uuid = uuid.NewString()

	b, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(b)

	res := model.FromCtx(ctx).Create(&user)
	return user, res.Error
}

func GetUserList(ctx context.Context) ([]*User, error) {
	var userList []*User
	res := model.FromCtx(ctx).Find(&userList)
	return userList, res.Error
}
