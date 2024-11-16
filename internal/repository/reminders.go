package repository

import (
	"context"
	"errors"
	"time"

	"github.com/justIGreK/Reminders/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RemindersRepo struct {
	collection *mongo.Collection
}

func NewRemindersRepository(db *mongo.Client) *RemindersRepo {
	return &RemindersRepo{
		collection: db.Database(dbname).Collection(remindersCollection),
	}
}

func (r *RemindersRepo) AddReminder(ctx context.Context, reminder models.Reminder) (string, error) {
	reminder.IsActive = true
	res, err := r.collection.InsertOne(ctx, reminder)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}
func (r *RemindersRepo) GetUpcomingReminders(ctx context.Context) ([]models.Reminder, error) {
	now := time.Now().UTC()
	filter := bson.M{
		"time":      bson.M{"$lte": now.Add(60 * time.Second)},
		"is_active": true,
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var reminders []models.Reminder
	if err := cursor.All(ctx, &reminders); err != nil {
		return nil, err
	}

	return reminders, nil
}

func (r *RemindersRepo) MarkReminderAsInactive(ctx context.Context, userID, id string) error {
	oid, err := convertToObjectIDs(id)
	if err != nil {
		return err
	}
	filter := bson.M{
		"_id":       oid[0],
		"user_id":   userID,
		"is_active": true,
	}
	update := bson.M{
		"$set": bson.M{
			"is_active": false,
		},
	}
	res, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return errors.New("no changes")
	}
	return nil
}

func (r *RemindersRepo) DeleteReminder(ctx context.Context, userID, rmID string) error {
	oid, err := convertToObjectIDs(rmID)
	if err != nil {
		return err
	}
	filter := bson.M{
		"_id":       oid[0],
		"user_id":   userID,
	}
	res, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return errors.New("no changes")
	}
	return nil
}

func (r *RemindersRepo) GetReminders(ctx context.Context, userID string) ([]models.Reminder, error) {
	filter := bson.M{
		"user_id":   userID,
		"is_active": true,
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reminders []models.Reminder
	err = cursor.All(ctx, &reminders)
	if err != nil {
		return nil, err
	}
	return reminders, nil
}
