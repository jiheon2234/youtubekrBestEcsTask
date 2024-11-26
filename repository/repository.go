package repository

import (
	"collectYoutubeData/service/types"
	"context"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"log"
	"time"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(dsn string) *PostgresRepository {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer func() {
		cancel()
		if ctx.Err() == context.DeadlineExceeded {
			log.Fatalf("Connection attempt to PostgreSQL timed out after 5 seconds.\n")
		}
	}()

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to open database: %v", err))
	}
	if err := db.Ping(); err != nil {
		db.Close()
		log.Fatal(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	fmt.Println("Successfully connected to PostgreSQL database!")
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) GetTx() (*sql.Tx, error) {
	tx, err := r.db.Begin()
	if err != nil {
		log.Fatalf("cannot get TX: %w", err)
	}
	return tx, nil
}

func (r *PostgresRepository) InsertBatchLog() (logID int64, err error) {
	query := `
        INSERT INTO batch_log (start_time)
        VALUES ($1)
        RETURNING id
    `

	// 현재 시간 삽입 후 자동 생성된 ID 반환
	err = r.db.QueryRow(query, time.Now()).Scan(&logID)
	if err != nil {
		return 0, fmt.Errorf("failed to insert batch log: %w", err)
	}

	return logID, nil
}

func (r *PostgresRepository) UpdateBatchLogToSuccess(logId int64) (err error) {
	query := `
        UPDATE batch_log
        SET end_time = $1, success = $2
        WHERE id = $3
    `

	_, err = r.db.Exec(query, time.Now(), true, logId)
	if err != nil {
		return fmt.Errorf("failed to update batch log: %w", err)
	}

	return nil
}

func (r *PostgresRepository) InsertVideo(tx *sql.Tx, video *types.VideoDto) (int64, error) {
	query := `
        INSERT INTO videos 
        (y_id, published_at, channel_id, title, description, thumbnail_url, 
        view_count, like_count, favorite_count, comment_count, tags, saved_at) 
        VALUES 
        ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
        RETURNING id
    `

	var lastInsertID int64
	err := tx.QueryRow(query,
		video.YID,
		video.PublishedAt,
		video.ChannelId,
		video.Title,
		video.Description,
		video.ThumbnailURL,
		video.ViewCount,
		video.LikeCount,
		video.FavoriteCount,
		video.CommentCount,
		pq.Array(video.Tags),
		time.Now(),
	).Scan(&lastInsertID)

	if err != nil {
		return 0, fmt.Errorf("failed to insert video: %w", err)
	}

	return lastInsertID, nil
}

func (r *PostgresRepository) InsertComments(tx *sql.Tx, comments []*types.CommentDto, videoId int64) error {
	// 기본 INSERT 쿼리 생성
	query := `
        INSERT INTO comments 
        (video_id, y_id, author_name, author_image, like_count, published_at, updated_at, total_reply_count, saved_at, text_display)
        VALUES 
    `

	// 쿼리와 파라미터 동적 생성
	values := []interface{}{}
	for i, comment := range comments {
		// 각 행에 대해 인덱스 기반 변수 추가
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d),",
			i*10+1, i*10+2, i*10+3, i*10+4, i*10+5, i*10+6, i*10+7, i*10+8, i*10+9, i*10+10)

		// 각 값들을 values 슬라이스에 추가
		values = append(values,
			videoId,
			comment.YID,
			comment.AuthorName,
			comment.AuthorImage,
			comment.LikeCount,
			comment.PublishedAt,
			comment.UpdatedAt,
			comment.TotalReplyCount,
			time.Now(),          // saved_at
			comment.TextDisplay, // 새 필드 추가
		)
	}

	// 마지막 쉼표 제거
	query = query[:len(query)-1]

	// 트랜잭션 내에서 쿼리 실행
	_, err := tx.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("failed to insert comments: %w", err)
	}

	return nil
}

func (r *PostgresRepository) Close() {
	r.db.Close()
}
