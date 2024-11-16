package client

import (
	"context"

	tz "github.com/justIGreK/Reminders-Timezone/pkg/go/timezone"
	"github.com/justIGreK/Reminders/internal/models"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TimezoneClient struct {
	client tz.TimezoneServiceClient
}

func NewTimezoneClient(serviceAddress string) (*TimezoneClient, error) {
	conn, err := grpc.NewClient(serviceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &TimezoneClient{
		client: tz.NewTimezoneServiceClient(conn),
	}, nil
}

func (tc *TimezoneClient) GetTimezone(ctx context.Context, userID string) (*models.Timezone, error) {
	req := &tz.GetTimezoneRequest{UserId: userID}
	res, err := tc.client.GetTimezone(ctx, req)
	if err != nil {
		return nil, err
	}
	return &models.Timezone{
		UserID: res.Timezone.UserId,
		Latitude: float64(res.Timezone.Latitude),
		Longitude: float64(res.Timezone.Longitude),
		DiffHour: int(res.Timezone.Diffhout),
	}, nil
}