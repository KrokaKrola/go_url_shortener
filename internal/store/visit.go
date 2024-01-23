package store

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Visit struct {
	ID        int64     `json:"id"`
	LinkID    int64     `json:"link_id"`
	UserAgent string    `json:"user_agent"`
	IPAddress string    `json:"ip_address"`
	Referer   string    `json:"referer"`
	CreatedAt time.Time `json:"created_at"`
}

func NewVisit(linkId int, ua string, ip string, referer string) *Visit {
	return &Visit{
		LinkID:    int64(linkId),
		UserAgent: ua,
		IPAddress: ip,
		Referer:   referer,
	}
}

func (v *Visit) CreateOne(ctx context.Context, connection *pgx.Conn) error {
	_, err := connection.Exec(
		ctx,
		"INSERT INTO visits (link_id, user_agent, ip_address, referer) VALUES ($1, $2, $3, $4)",
		v.LinkID,
		v.UserAgent,
		v.IPAddress,
		v.Referer,
	)

	if err != nil {
		return err
	}

	return nil
}

func (v *Visit) GetVisitsCount(ctx context.Context, connection *pgx.Conn) (int64, error) {
	var count int64

	err := connection.QueryRow(
		ctx,
		"SELECT COUNT(*) FROM visits WHERE link_id = $1",
		v.LinkID,
	).Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}
