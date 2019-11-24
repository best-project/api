package storage

import (
	"fmt"
	"github.com/best-project/api/internal"
	"github.com/jinzhu/gorm"
	"time"
)

type User struct {
	db *gorm.DB
}

func (u *User) SaveUser(user *internal.User) error {
	return u.db.Save(user).Error
}

func (u *User) UpdateUser(user *internal.User) error {
	user.UpdatedAt = time.Now()
	return u.db.Save(user).Error
}

func (u *User) Exist(user *internal.User) bool {
	u.db.RLock()
	defer u.db.RUnlock()
	users := make([]internal.User, 0)

	u.db.Where(user).Find(&users)
	if len(users) > 0 {
		return true
	}
	return false
}

func (u *User) GetByName(username string) ([]internal.User, error) {
	u.db.RLock()
	defer u.db.RUnlock()
	users := make([]internal.User, 0)

	u.db.Where(&internal.User{Username: username}).Find(&users)
	if len(users) > 1 {
		return nil, fmt.Errorf("found more then one user with name %s: ", username)
	}

	return users, nil
}

func (u *User) GetManyByID(ids []uint) ([]internal.User, error) {
	u.db.Lock()
	defer u.db.Unlock()

	users := make([]internal.User, 0)
	u.db.Where(ids).Find(&users)

	return users, nil
}
