package model

import (
	"gorm.io/gorm"
	usermodel "yunyandz.com/tiktok/user-part/model"
)

type Video struct {
	gorm.Model

	AuthorID uint64 `gorm:"column:author_id"`

	Title       string `gorm:"size:128"`
	Description string `gorm:"size:1024"`
	Playurl     string `gorm:"size:1024"`
	Coverurl    string `gorm:"size:1024"`

	Commentcount uint64
	Likecount    uint64 `gorm:"default:0"`

	Likes    []*usermodel.User `gorm:"many2many:user_likes"`
	Comments []*Comment        `gorm:"many2many:video_comments"`
}
