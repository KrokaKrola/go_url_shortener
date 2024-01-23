package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/krokakrola/url_shortener/internal/utils"
)

type Link struct {
	ID        int64     `json:"id"`
	URL       string    `json:"url"`
	Hash      string    `json:"hash"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateLinkRequest struct {
	URL string `json:"url" validate:"required,url,max=1024"`
}

type CreateLinkResponse struct {
	Hash string `json:"hash"`
}

func NewCreateLinkRequest(url string) *CreateLinkRequest {
	return &CreateLinkRequest{
		URL: url,
	}
}

func NewCreateLinkResponse(hash string) map[string]interface{} {
	return map[string]interface{}{
		"hash": hash,
	}
}

func NewLink(url string, hash string) *Link {
	return &Link{
		URL:  url,
		Hash: hash,
	}
}

func (l *Link) CreateOne(ctx context.Context, conn *pgx.Conn) (string, error) {
	randomHash := utils.GenerateRandomString(8)

	_, err := conn.Exec(
		ctx,
		"INSERT INTO links (url, hash) VALUES ($1, $2)",
		l.URL,
		randomHash,
	)

	if err != nil {
		return "", err
	}

	return randomHash, nil
}

func (l *Link) FindOneByHash(ctx context.Context, conn *pgx.Conn, hash string) (*Link, error) {
	row := conn.QueryRow(ctx, "SELECT id, url, hash, created_at FROM links WHERE hash = $1", hash)

	var link Link

	err := row.Scan(&link.ID, &link.URL, &link.Hash, &link.CreatedAt)

	if err != nil {
		return nil, err
	}

	return &link, nil
}
