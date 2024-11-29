package service

import (
	"collectYoutubeData/config"
	"collectYoutubeData/repository"
	"collectYoutubeData/service/types"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
	"net/http"
	"sync"
)

type YoutubeService struct {
	APiKey        string
	VideoApiURL   string
	CommentApiURL string
	client        *resty.Client
	CommentGoCnt  int
	repository    *repository.PostgresRepository
}

func NewYoutubeService(cfg *config.Config) *YoutubeService {
	y := &YoutubeService{
		APiKey:        cfg.ApiKey,
		VideoApiURL:   cfg.VideoApiURL,
		CommentApiURL: cfg.CommentApiURL,
		client:        resty.New(),
		CommentGoCnt:  cfg.GoRoutineCnt,
		repository:    repository.NewPostgresRepository(cfg.PgDsn),
	}
	return y
}

func (y *YoutubeService) Run() {
	logId, err := y.repository.InsertBatchLog()

	videoCH := y.GetPopularVideo()
	ResultCH := y.GetVideoComments(videoCH)

	tx, err := y.repository.GetTx()
	if err != nil {
		panic(err)
	}

	for res := range ResultCH {
		videoId, err := y.repository.InsertVideo(tx, res.Video)
		if err != nil {
			log.Fatalf("err : %v", err)
		}
		err = y.repository.InsertComments(tx, res.Comments, videoId)
		if err != nil {
			panic(err)
		}
	}

	err = tx.Commit()
	fmt.Printf("err = %v", err)

	err = y.repository.UpdateBatchLogToSuccess(logId)
	y.repository.Close()
}

func (y *YoutubeService) getApiDataToBiteArr(baseUrl string, queryParams map[string]string) ([]byte, error) {

	resp, err := y.client.R().
		SetQueryParams(queryParams).
		SetHeader("Accept", "application/json").
		Get(baseUrl)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		log.Printf("status code = %d\n", resp.StatusCode())
		return nil, err
	}

	return resp.Body(), nil
}

func (y *YoutubeService) GetPopularVideo() <-chan *types.VideoItem {

	videoCh := make(chan *types.VideoItem)

	queryParams := map[string]string{
		"key":        y.APiKey,
		"part":       "snippet,statistics",
		"chart":      "mostPopular",
		"regionCode": "KR",
		"PageToken":  "",
	}

	go func() {
		for i := 0; i < 10; i++ {
			var youtubeVideoRes types.YoutubeVideoRes
			biteArr, err := y.getApiDataToBiteArr(y.VideoApiURL, queryParams)
			if err != nil {
				panic(err)
			}
			err = json.Unmarshal(biteArr, &youtubeVideoRes)
			queryParams["pageToken"] = youtubeVideoRes.NextPageToken

			items := youtubeVideoRes.Items
			for _, item := range items {
				videoCh <- &item
			}
		}
		close(videoCh)
	}()
	return videoCh
}

func (y *YoutubeService) GetVideoComments(videoCh <-chan *types.VideoItem) <-chan *types.ResultDto {

	dtoChan := make(chan *types.ResultDto)

	var wg sync.WaitGroup

	for i := 0; i < y.CommentGoCnt; i++ {
		wg.Add(1)
		go func(goId int) {
			defer wg.Done()
			fmt.Printf("GORUTINE %d start \n", goId)

			var commentRes types.CommentRes
			for videoItem := range videoCh {
				queryParams := map[string]string{
					"key":     y.APiKey,
					"part":    "snippet",
					"order":   "relevance",
					"videoId": videoItem.Id,
				}
				biteArr, err := y.getApiDataToBiteArr(y.CommentApiURL, queryParams)
				if err != nil {
					fmt.Println(videoItem.Id)
					continue
				}

				err = json.Unmarshal(biteArr, &commentRes)
				if err != nil {
					continue
				}

				video := types.NewVideoDto(videoItem)
				comments := make([]*types.CommentDto, len(commentRes.Items))

				for i, _ := range comments {
					comments[i] = types.NewCommentDto(commentRes.Items[i])
				}

				dto := &types.ResultDto{
					Video:    video,
					Comments: comments,
				}

				log.Printf("videotitle = %s \ncommentCount=%d \n\n", video.Title, len(comments))

				dtoChan <- dto

			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(dtoChan)
	}()

	return dtoChan
}
