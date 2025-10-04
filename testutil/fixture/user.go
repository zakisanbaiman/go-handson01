package fixture

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/zakisanbaiman/go-handson01/entity"
)

// デフォルト値を基本に、setterで指定された値で上書きする
// 参考: https://engineering.mercari.com/blog/entry/20220411-42fc0ba69c/
func User(setter func(*entity.User)) *entity.User {
	// デフォルト値
	user := &entity.User{
		ID:         entity.UserID(rand.Int()),
		Name:       "test" + strconv.Itoa(rand.Int())[:5],
		Password:   "password",
		Role:       "admin",
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	if setter != nil {
		setter(user)
	}
	return user
}
