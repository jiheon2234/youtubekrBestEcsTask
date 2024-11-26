package types

import (
	"time"
)

// VideoDto 구조체 정의
type VideoDto struct {
	YID           string
	PublishedAt   time.Time
	ChannelId     string
	Title         string
	Description   string
	ThumbnailURL  string
	ViewCount     int
	LikeCount     int
	FavoriteCount int
	CommentCount  int
	Tags          []string
}

func NewVideoDto(videoItem *VideoItem) *VideoDto {
	return &VideoDto{
		YID:           videoItem.Id,
		PublishedAt:   videoItem.Snippet.PublishedAt,
		ChannelId:     videoItem.Snippet.ChannelId,
		Title:         videoItem.Snippet.Title,
		Description:   videoItem.Snippet.Description,
		ThumbnailURL:  videoItem.Snippet.Thumbnails.Default.Url,
		ViewCount:     videoItem.Statistics.ViewCount,
		LikeCount:     videoItem.Statistics.LikeCount,
		FavoriteCount: videoItem.Statistics.FavoriteCount,
		CommentCount:  videoItem.Statistics.CommentCount,
		Tags:          getEmptySliceIfNil(videoItem.Snippet.Tags),
	}
}

func getEmptySliceIfNil[T any](sl []T) []T {
	if sl == nil {
		return []T{}
	}
	return sl
}

type CommentDto struct {
	//VideoId         int
	TextDisplay     string
	YID             string
	AuthorName      string
	AuthorImage     string
	LikeCount       int
	PublishedAt     time.Time
	UpdatedAt       time.Time
	TotalReplyCount int
	SavedAt         time.Time
}

func NewCommentDto(commentItem *CommentItem) *CommentDto {
	// can get VideoId when video saved to DB
	return &CommentDto{
		TextDisplay:     commentItem.Snippet.TopLevelComment.Snippet.TextDisplay,
		YID:             commentItem.Id,
		AuthorName:      commentItem.Snippet.TopLevelComment.Snippet.AuthorDisplayName,
		AuthorImage:     commentItem.Snippet.TopLevelComment.Snippet.AuthorProfileImageUrl,
		LikeCount:       commentItem.Snippet.TopLevelComment.Snippet.LikeCount,
		PublishedAt:     commentItem.Snippet.TopLevelComment.Snippet.PublishedAt,
		UpdatedAt:       commentItem.Snippet.TopLevelComment.Snippet.UpdatedAt,
		TotalReplyCount: commentItem.Snippet.TotalReplyCount,
	}
}

type ResultDto struct {
	Video    *VideoDto
	Comments []*CommentDto
}
