package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/justIGreK/Reminders/internal/models"
)

type RemindersRepository interface {
	AddReminder(ctx context.Context, reminder models.Reminder) (string, error)
	GetReminders(ctx context.Context, userID string) ([]models.Reminder, error)
	GetUpcomingReminders(ctx context.Context) ([]models.Reminder, error)
	MarkReminderAsInactive(ctx context.Context, userID, rmID string) error
	DeleteReminder(ctx context.Context, userID, rmID string) error
}

type RemindersService struct {
	RemindersRepo RemindersRepository
	Tz            TzService
}

func NewRemindersService(rsRepo RemindersRepository, tz TzService) *RemindersService {
	return &RemindersService{RemindersRepo: rsRepo,
		Tz: tz}
}

type TzService interface {
	GetTimezone(ctx context.Context, userID string) (*models.Timezone, error)
}

const (
	DateTimeFormat string = "2006-01-02 15:04"
	DateFormat     string = "2006-01-02"
	TimeFormat     string = "15:04"
)

func (s *RemindersService) AddReminder(ctx context.Context, createRms models.CreateRms) (string, error) {
	tz, err := s.Tz.GetTimezone(ctx, createRms.UserID)
	if err != nil {
		return "", err
	}
	rmTime, err := time.Parse(TimeFormat, createRms.Time)
	if err != nil {
		log.Println(err)
		return "", err
	}
	now := time.Now().UTC()
	rmDate := now
	if createRms.Date != nil {
		rmDate, err = time.Parse(DateFormat, *createRms.Date)
		if err != nil {
			log.Println(err)
			return "", err
		}
	}
	originalTime := time.Date(rmDate.Year(), rmDate.Month(), rmDate.Day(), rmTime.Hour(), rmTime.Minute(), rmTime.Second(), 0, time.Now().UTC().Location())
	utcTime := originalTime.Add(-time.Duration(tz.DiffHour) * time.Hour)
	if utcTime.Before(now) {
		utcTime = utcTime.Add(24 * time.Hour)
		originalTime = originalTime.Add(24 * time.Hour)
	}
	reminder := models.Reminder{
		UserID:       createRms.UserID,
		Action:       createRms.Action,
		Time:         utcTime,
		OriginalTime: originalTime,
	}
	_, err = s.RemindersRepo.AddReminder(ctx, reminder)
	if err != nil {
		return "", err
	}
	response := fmt.Sprintf("Reminder *%s* was set on %s", reminder.Action, reminder.OriginalTime.Format("2006-01-02 15:04"))
	return response, nil
}

func (s *RemindersService) GetUpcomingReminders(ctx context.Context) ([]models.Reminder, error) {
	rms, err := s.RemindersRepo.GetUpcomingReminders(ctx)
	if err != nil{
		log.Println(err)
		return nil, err
	}
	return rms, nil
}

func (s *RemindersService) MarkReminderAsSent(ctx context.Context, userID, rmID string) error {
	err := s.RemindersRepo.MarkReminderAsInactive(ctx, userID, rmID)
	if err != nil{
		log.Println(err)
		return err
	}
	return nil
}
func (s *RemindersService) DeleteReminder(ctx context.Context, userID, rmID string) error {
	err := s.RemindersRepo.DeleteReminder(ctx, userID, rmID)
	if err != nil{
		log.Println(err)
		return err
	}
	return nil
}
