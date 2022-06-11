package model

import (
	"gorm.io/gorm"
	"yunyandz.com/tiktok/proto/pb"
)

type User struct {
	gorm.Model

	Username string `gorm:"size:32;unique_index"`
	Password string `gorm:"size:256"`

	FollowCount   int64 `gorm:"type:int"`
	FollowerCount int64 `gorm:"type:int"`

	Videos    []*Video `gorm:"many2many:user_videos"`
	Followers []*User  `gorm:"many2many:user_follows"`
	Likes     []*Video `gorm:"many2many:user_likes"`
}

func (u *User) ToPB(isFollow bool) pb.User {
	user := pb.User{
		Id:            int64(u.ID),
		Name:          u.Username,
		FollowCount:   u.FollowCount,
		FollowerCount: u.FollowerCount,
		IsFollow:      isFollow,
	}

	return user
}
