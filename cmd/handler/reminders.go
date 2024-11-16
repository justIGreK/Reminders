package handler

import (
	"context"

	"github.com/justIGreK/Reminders/internal/models"
	remindersProto "github.com/justIGreK/Reminders/pkg/go/reminders"
	"google.golang.org/protobuf/types/known/emptypb"
)

type RemindersServiceServer struct {
	remindersProto.UnimplementedRemindersServiceServer
	RmsSRV RemindersService
}

type RemindersService interface {
	AddReminder(ctx context.Context, createRms models.CreateRms) (string, error)
	GetUpcomingReminders(ctx context.Context) ([]models.Reminder, error)
	MarkReminderAsSent(ctx context.Context, userID, rmID string) error
	DeleteReminder(ctx context.Context, userID, rmID string) error
}

const (
	DateTimeFormat string = "2006-01-02 15:04"
	DateFormat     string = "2006-01-02"
	TimeFormat     string = "15:04"
)

func (s *RemindersServiceServer) AddReminder(ctx context.Context, req *remindersProto.AddReminderRequest) (*remindersProto.AddReminderResponse, error) {
	reminderReq := models.CreateRms{
		UserID: req.Request.UserId,
		Action: req.Request.Action,
		Time:   req.Request.Time,
	}
	if req.Request.Date != nil {
		reminderReq.Date = &req.Request.Date.Value
	}
	resp, err := s.RmsSRV.AddReminder(ctx, reminderReq)
	if err != nil {
		return nil, err
	}

	return &remindersProto.AddReminderResponse{
		Response: resp,
	}, nil

}

func (s *RemindersServiceServer) GetUpcomingReminders(ctx context.Context, empty *emptypb.Empty) (*remindersProto.GetUpcomingRemindersResponse, error) {
	rms, err := s.RmsSRV.GetUpcomingReminders(ctx)
	if err != nil {
		return nil, err
	}

	protoRms := s.convertToProtoReminders(rms)
	return &remindersProto.GetUpcomingRemindersResponse{
		Reminders: protoRms,
	}, nil

}

func (s *RemindersServiceServer) convertToProtoReminders(rms []models.Reminder) []*remindersProto.Reminder {
	protoRms := make([]*remindersProto.Reminder, len(rms))
	for i, b := range rms {
		protoRms[i] = &remindersProto.Reminder{
			Id:           b.ID,
			UserId:       b.UserID,
			Action:       b.Action,
			Utctime:      b.Time.UTC().Format(DateTimeFormat),
			Originaltime: b.OriginalTime.Format(DateTimeFormat),
		}
	}
	return protoRms
}
func (s *RemindersServiceServer) MarkReminderAsSent(ctx context.Context, req *remindersProto.DeleteReminderRequest) (*emptypb.Empty, error) {
	err := s.RmsSRV.MarkReminderAsSent(ctx, req.UserId, req.RmId)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *RemindersServiceServer) DeleteReminder(ctx context.Context, req *remindersProto.DeleteReminderRequest) (*emptypb.Empty, error) {
	err := s.RmsSRV.DeleteReminder(ctx, req.UserId, req.RmId)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
