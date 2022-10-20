package storage

type Event struct {
	ID           string `db:"id" json:"id"`
	Title        string `db:"title" json:"title"`
	EventDate    int64  `db:"date" json:"date"`
	Duration     int64  `db:"durationUntil" json:"durationUntil"`
	Description  string `db:"description" json:"description"`
	UserID       string `db:"userId" json:"userId"`
	NoticeBefore int64  `db:"noticeBefore" json:"noticeBefore"`
}
