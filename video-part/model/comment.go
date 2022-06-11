package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model

	UserID  uint64 `gorm:"index"`
	VideoID uint64 `gorm:"index"`
	Content string `gorm:"size:1024"`
}

// 获取视频的评论列表
func (v *VideoModel) GetVideoComments(id uint64) ([]*Comment, error) {
	var comments []*Comment
	if err := v.db.Where("video_id = ?", id).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}

func (v *VideoModel) GetVideoCommentsCount(id uint64) (int64, error) {
	var count int64
	if err := v.db.Model(Comment{}).Where("video_id = ?", id).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// 使用Create创建一条评论
func (v *VideoModel) CreateAComment(videoId uint64, userId uint64, content string) (*Comment, error) {
	comment := Comment{UserID: userId, VideoID: videoId, Content: content}
	err := v.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&comment).Error
		if err != nil {
			return err
		}
		if err := tx.Model(&Video{}).Where("id = ?", videoId).
			Update("commentcount", gorm.Expr("commentcount + ?", 1)).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &comment, err
}

// 查询一条评论
func (v *VideoModel) FindAComment(commentId uint64) (*Comment, error) {
	var comment Comment
	if err := v.db.First(&comment, commentId).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

// 删除一个评论
func (v *VideoModel) DeleteAComment(commentId uint64) error {
	err := v.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&Comment{}, commentId).Error; err != nil {
			return err
		}
		if err := tx.Model(&Video{}).Where("id = ?", commentId).
			Update("commentcount", gorm.Expr("commentcount - ?", 1)).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// gorm.Model包含deletedat字段，可以软删除
// 调用Delete后，记录仍存在但无法查询到
// deletedat -> “删除”时间
