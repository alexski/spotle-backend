package model

import (
	"database/sql"
	"time"
)

type Artist struct {
	ID               int       `json:"id"`
	Name             string    `json:"artist_name"`
	MonthlyListeners int32     `json:"monthly_listeners"`
	LastChecked      time.Time `json:"last_checked"`
}

func (r *Artist) GetArtist(db *sql.DB) error {
	return db.QueryRow("SELECT artist_name, monthly_listeners, last_checked FROM artist WHERE id=$1", r.ID).Scan(&r.Name, &r.MonthlyListeners, &r.LastChecked)
}

func (r *Artist) UpdateArtist(db *sql.DB) error {
	_, err := db.Exec("UPDATE artist SET artist_name=$1, monthly_listeners=$2, last_checked=$3 WHERE id=$4", r.Name, r.MonthlyListeners, time.Now().UTC(), r.ID)
	return err
}

func (r *Artist) DeleteArtist(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM artist WHERE id=$1", r.ID)
	return err
}

func (r *Artist) CreateArtist(db *sql.DB) error {
	err := db.QueryRow("INSERT INTO artist (artist_name, monthly_listeners, last_checked) values ($1, $2, $3) RETURNING id", r.Name, r.MonthlyListeners, time.Now().UTC()).Scan(&r.ID)
	return err
}

func GetArtists(db *sql.DB, start int, count int) ([]Artist, error) {
	rows, err := db.Query(
		"SELECT id, artist_name, monthly_listeners, last_checked FROM artist LIMIT $1 OFFSET $2",
		count, start)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	artists := []Artist{}
	for rows.Next() {
		var r Artist
		if err := rows.Scan(&r.ID, &r.Name, &r.MonthlyListeners, &r.LastChecked); err != nil {
			return nil, err
		}
		artists = append(artists, r)
	}

	return artists, nil
}
