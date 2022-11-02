package storage

// Event Событие
type Event struct {
	// ID - уникальный идентификатор события (можно воспользоваться UUID)
	ID string `db:"id" json:"id"`
	// Заголовок - короткий текст
	Title string `db:"title" json:"title"`
	// Дата и время события
	EventDate int64 `db:"date" json:"date"`
	// Длительность события (или дата и время окончания)
	Duration int64 `db:"duration_until" json:"duration_until"`
	// Описание события - длинный текст, опционально
	Description string `db:"description" json:"description"`
	// ID пользователя, владельца события
	UserID string `db:"owner_id" json:"owner_id"`
	// За сколько времени высылать уведомление, опционально.
	NoticeBefore int64 `db:"notice_before" json:"notice_before"`
}
