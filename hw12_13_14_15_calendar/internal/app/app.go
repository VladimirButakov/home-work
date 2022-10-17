package app

import (
	"fmt"
	"github.com/VladimirButakov/home-work/tree/master/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
	"time"
)

const userID string = "main_owner_id"

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Debug(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
}

type Storage interface {
	AddEvent(event storage.Event) error
	UpdateEvent(event storage.Event) error
	RemoveEvent(eventID string) error
	DailyEvents(date time.Time) ([]storage.Event, error)
	WeeklyEvents(date time.Time) ([]storage.Event, error)
	MonthEvents(date time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{logger, storage}
}

func (a *App) CreateEvent(title string, date int64) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("can not create unique id, %w", err)
	}

	event := storage.Event{ID: id.String(), Title: title, EventDate: date, UserID: userID}

	return a.storage.AddEvent(event)
}

func (a *App) UpdateEvent(id, title string, date int64, description string, durationUntil int64, ownerID string, noticeBefore int64) error {
	event := storage.Event{ID: id, Title: title, EventDate: date, Duration: durationUntil, Description: description, UserID: ownerID, NoticeBefore: noticeBefore}

	return a.storage.UpdateEvent(event)
}

func (a *App) RemoveEvent(id string) error {
	return a.storage.RemoveEvent(id)
}

func (a *App) DailyEvents(date time.Time) ([]storage.Event, error) {
	return a.storage.DailyEvents(date)
}

func (a *App) WeeklyEvents(date time.Time) ([]storage.Event, error) {
	return a.storage.WeeklyEvents(date)
}

func (a *App) MonthEvents(date time.Time) ([]storage.Event, error) {
	return a.storage.MonthEvents(date)
}

func (a *App) GetLogger() Logger {
	return a.logger
}
