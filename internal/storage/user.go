package storage

import (
	"fmt"
	"github.com/best-project/api/internal"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type User struct {
	db *gorm.DB
}

func (u *User) SaveUser(user *internal.User) error {
	return u.db.Save(user).Error
}

func (u *User) Exist(user *internal.User) bool {
	u.db.RLock()
	defer u.db.RUnlock()
	users := make([]internal.User, 0)

	u.db.Where(&internal.User{Email: user.Email}).Find(&users)
	if len(users) > 0 {
		return true
	}
	return false
}

func (u *User) GetByID(id uint) (*internal.User, error) {
	u.db.RLock()
	defer u.db.RUnlock()
	user := &internal.User{}

	if err := u.db.First(user, id).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (u *User) GetManyByID(ids []uint) ([]internal.User, error) {
	u.db.Lock()
	defer u.db.Unlock()

	users := make([]internal.User, 0)
	u.db.Where(ids).Find(&users)

	return users, nil
}

func (u *User) GetByMail(email string) ([]internal.User, error) {
	u.db.RLock()
	defer u.db.RUnlock()
	users := make([]internal.User, 0)

	u.db.Where(&internal.User{Email: email}).Find(&users)
	if len(users) > 1 {
		return nil, fmt.Errorf("found more then one user with mail %s: ", email)
	}

	return users, nil
}

func (u *User) GetAll() ([]internal.User, error) {
	u.db.RLock()
	defer u.db.RUnlock()
	users := make([]internal.User, 0)

	if err := u.db.Find(&users).Error; err != nil {
		return nil, errors.Wrapf(err, "while getting courses")
	}

	return users, nil
}
