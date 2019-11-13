package storage

import (
	"fmt"
	"github.com/best-project/api/internal"
	"github.com/jinzhu/gorm"
)

type User struct {
	db *gorm.DB
}

func (u *User) SaveUser(user *internal.User) {
	u.db.Lock()
	defer u.db.Unlock()

	u.db.Save(user)
}

func (u *User) Exist(user *internal.User) bool {
	u.db.RLock()
	defer u.db.RUnlock()
	users := make([]internal.User, 0)

	u.db.Where(user).Find(users)
	if len(users) > 0 {
		return true
	}
	return false
}

func (u *User) GetByName(username string) (*internal.User, error) {
	u.db.RLock()
	defer u.db.RUnlock()
	users := make([]internal.User, 0)

	u.db.Where(&internal.User{Username: username}).Find(users)
	if len(users) == 0 {
		return nil, fmt.Errorf("not found user with name %s: ", username)
	}

	return &users[0], nil
}
