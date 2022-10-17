package storage

type Event struct {
	ID           string `db:"id" json:"id"`
	Title        string `db:"title" json:"title"`
	EventDate    int64  `db:"date" json:"date"`
	Duration     int64  `db:"duration_until" json:"duration_until"`
	Description  string `db:"description" json:"description"`
	UserID       string `db:"user_id" json:"user_id"`
	NoticeBefore int64  `db:"notice_before" json:"notice_before"`
}
