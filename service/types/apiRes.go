package types

import "time"

type YoutubeVideoRes struct {
	Kind          string      `json:"kind"`
	Etag          string      `json:"etag"`
	Items         []VideoItem `json:"items"`
	NextPageToken string      `json:"nextPageToken"`
	PageInfo      struct {
		TotalResults   int `json:"totalResults"`
		ResultsPerPage int `json:"resultsPerPage"`
	} `json:"pageInfo"`
}

type VideoItem struct {
	Id      string `json:"id"`
	Snippet struct {
		PublishedAt time.Time `json:"publishedAt"`
		ChannelId   string    `json:"channelId"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Thumbnails  struct {
			Default struct {
				Url string `json:"url"`
			} `json:"default"`
		} `json:"thumbnails"`
		Tags []string `json:"tags"`
	} `json:"snippet"`
	Statistics struct {
		ViewCount     int `json:"viewCount,string"`
		LikeCount     int `json:"likeCount,string"`
		FavoriteCount int `json:"favoriteCount,string"`
		CommentCount  int `json:"commentCount,string"`
	} `json:"statistics"`
}

type CommentRes struct {
	Kind  string         `json:"kind"`
	Items []*CommentItem `json:"items"`
}

type CommentItem struct {
	//Kind    string `json:"kind"`
	//Etag    string `json:"etag"`
	Id      string `json:"id"`
	Snippet struct {
		//ChannelId       string `json:"channelId"`
		VideoId         string `json:"videoId"`
		TopLevelComment struct {
			//Kind    string `json:"kind"`
			//Etag    string `json:"etag"`
			Id      string `json:"id"`
			Snippet struct {
				//ChannelId   string `json:"channelId"`
				//VideoId     string `json:"videoId"`
				TextDisplay string `json:"TextDisplay"`
				//TextOriginal          string `json:"textOriginal"`
				AuthorDisplayName     string `json:"authorDisplayName"`
				AuthorProfileImageUrl string `json:"authorProfileImageUrl"`
				//AuthorChannelUrl      string `json:"authorChannelUrl"`
				//AuthorChannelId       struct {
				//	Value string `json:"value"`
				//} `json:"authorChannelId"`
				//CanRate      bool      `json:"canRate"`
				ViewerRating string    `json:"viewerRating"`
				LikeCount    int       `json:"likeCount"`
				PublishedAt  time.Time `json:"publishedAt"`
				UpdatedAt    time.Time `json:"updatedAt"`
			} `json:"snippet"`
		} `json:"topLevelComment"`
		//CanReply        bool `json:"canReply"`
		TotalReplyCount int `json:"totalReplyCount"`
		//IsPublic        bool `json:"isPublic"`
	} `json:"snippet"`
}

//type ResultDto struct {
//	VideoItem    *VideoItem
//	CommentItems []*CommentItem
//}
