package mysql

import (
	"errors"
	"gorm.io/gorm"
	"yunyandz.com/tiktok/user-part/model"
)

func GetUser(id uint64) (*model.User, error) {
	var user model.User
	if err := DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByName(username string) (*model.User, error) {
	var user model.User
	if err := DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateUser(user *model.User) (id uint64, err error) {
	err = DB.Model(&model.User{}).Save(user).Error

	return uint64(user.ID), err
}

// 获取用户的关注列表
func GetFollowList(userId uint64) ([]*model.User, error) {
	var users []*model.User
	if err := DB.Where("id in (select follower_id from user_follows where user_id = ?)", userId).Find(&users).Error; err != nil {
		return nil, err
	}
	// Todo: redis缓存
	return users, nil
}

// 获取用户的粉丝列表
func GetFollowerList(userId uint64) ([]*model.User, error) {
	var users []*model.User
	if err := DB.Where("id in (select user_id from user_follows where follower_id = ?)", userId).Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func IsFollow(userId uint64, followerId uint64) (bool, error) {
	var count int64
	if err := DB.Raw("SELECT COUNT(*) FROM user_follows WHERE user_id = ? AND follower_id = ?", userId, followerId).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// 关注一个用户
func CreateFollow(userId uint64, followId uint64) error {
	// 使用事务保证一致性
	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("insert into user_follows (user_id, follower_id) values (?, ?)", userId, followId).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.User{}).Where("id = ?", userId).Update("follow_count", gorm.Expr("follow_count + 1")).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.User{}).Where("id = ?", followId).Update("follower_count", gorm.Expr("follower_count + 1")).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	// Todo: Redis缓存
	return nil
}

// 取消关注一个用户
func DeleteFollow(userId uint64, followId uint64) error {
	// 使用事务保证一致性
	err := DB.Transaction(func(tx *gorm.DB) error {
		if rows := tx.Exec("delete from user_follows where user_id = ? and follower_id = ?", userId, followId).RowsAffected; rows == 0 {
			return errors.New("关系不存在")
		}
		if err := tx.Model(&model.User{}).Where("id = ?", userId).Update("follow_count", gorm.Expr("follow_count - 1")).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.User{}).Where("id = ?", followId).Update("follower_count", gorm.Expr("follower_count - 1")).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	// Todo: Redis缓存
	return nil
}
