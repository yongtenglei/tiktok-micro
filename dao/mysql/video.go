package mysql

import (
	"gorm.io/gorm"
	"time"
	"yunyandz.com/tiktok/pkg/constant"
	"yunyandz.com/tiktok/pkg/errorx"
	usermodel "yunyandz.com/tiktok/user-part/model"
	videomodel "yunyandz.com/tiktok/video-part/model"
)

// 创建一个新的视频。
func CreateVideo(video *videomodel.Video) (uint64, error) {
	//事务
	err := DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(video).Error; err != nil {
			return err
		}
		var user usermodel.User
		if err := tx.First(&user, video.AuthorID).Error; err != nil {
			return err
		}
		if err := tx.Model(&user).Association("Videos").Append(video); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	return uint64(video.ID), nil
}

// 更新视频的播放地址，通常用于视频上传完成后
func UpdateVideo(id uint64, playurl string) error {
	if err := DB.Model(&videomodel.Video{}).Where("id = ?", id).Update("playurl", playurl).Error; err != nil {
		return err
	}
	return nil
}

// 获取最新的视频条目，按照时间降序排列，按照文档中的要求，这里只返回前30条
func GetNewVideos() ([]*videomodel.Video, error) {
	var videos []*videomodel.Video
	if err := DB.Limit(constant.FeedLimit).Find(&videos).Error; err != nil {
		return nil, err
	}
	return videos, nil
}

// 获取指定时间戳之前创建时间的视频列表。这里依然是最多返回30条。
func GetVideosBeforeTime(time time.Time) ([]*videomodel.Video, error) {
	var videos []*videomodel.Video
	if err := DB.Where("created_at < ?", time).Limit(constant.FeedLimit).Find(&videos).Error; err != nil {
		return nil, err
	}
	return videos, nil
}

// 获取视频的详情
func GetVideo(videoId uint64) (*videomodel.Video, error) {
	var video videomodel.Video
	if err := DB.First(&video, videoId).Error; err != nil {
		return nil, err
	}
	return &video, nil
}

// 获取用户的视频列表
func GetVideosByUser(userId uint64) ([]*videomodel.Video, error) {
	var videos []*videomodel.Video
	if err := DB.Where("author_id = ?", userId).Find(&videos).Error; err != nil {
		return nil, err
	}
	return videos, nil
}

// 获取用户的视频点赞列表
func GetUserLikeVideos(userId uint64) ([]*videomodel.Video, error) {
	var videos []*videomodel.Video
	if err := DB.Raw("SELECT * FROM videos WHERE id IN (SELECT video_id FROM user_likes WHERE user_id = ?)", userId).Scan(&videos).Error; err != nil {
		return nil, err
	}
	return videos, nil
}

// 点赞视频
func LikeVideo(userId uint64, videoId uint64) error {
	err := DB.Transaction(func(tx *gorm.DB) error {
		if row := DB.Exec("INSERT INTO user_likes (user_id, video_id) VALUES (?, ?)", userId, videoId).RowsAffected; row != 1 {
			return errorx.ErrUserAlreadyLikeVideo
		}
		if err := DB.Exec("UPDATE videos SET likecount = likecount + 1 WHERE id = ?", videoId).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// 取消点赞视频
func UnLikeVideo(userId uint64, videoId uint64) error {
	err := DB.Transaction(func(tx *gorm.DB) error {
		if DB.Exec("DELETE FROM user_likes WHERE user_id = ? AND video_id = ?", userId, videoId).RowsAffected == 0 {
			return errorx.ErrUserNotLikeVideo
		}
		if err := DB.Exec("UPDATE videos SET likecount = likecount - 1 WHERE id = ?", videoId).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

// 获取视频的点赞数
func GetVideoLikesCount(id uint64) (int64, error) {
	var count int64
	if err := DB.Model(&videomodel.Video{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

//查询视频点赞
func IsFavorite(userId uint64, videoId uint64) (bool, error) {
	var count int64
	if err := DB.Raw("SELECT COUNT(*) FROM user_likes WHERE user_id = ? AND video_id = ?", userId, videoId).Scan(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
