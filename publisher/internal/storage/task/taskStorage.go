package task

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/ruziba3vich/task-publisher/internal/models"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskStorage struct {
	queue           amqp.Queue
	channel         *amqp.Channel
	tasksCollection *mongo.Collection
	logger          *log.Logger
}

func New(queue amqp.Queue, channel *amqp.Channel, tasksCollection *mongo.Collection, logger *log.Logger) *TaskStorage {
	return &TaskStorage{
		queue:           queue,
		channel:         channel,
		tasksCollection: tasksCollection,
		logger:          logger,
	}
}

func (t *TaskStorage) SendTask(req *models.Task) error {
	req.TaskId = primitive.NewObjectID()
	req.Method = models.POST
	if err := t.channel.ExchangeDeclare(
		"prodonik",
		"fanout",
		false,
		false,
		false,
		false,
		nil,
	); err != nil {
		t.logger.Println(err)
		return err
	}

	message, err := json.Marshal(req)
	if err != nil {
		t.logger.Println(err)
		return err
	}
	return t.channel.Publish(
		"prodonik",
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        message,
		},
	)
}

func (t *TaskStorage) GetAllTasks(ctx context.Context) (*models.GetAllTasksResponse, error) {
	filter := bson.M{}
	cursor, err := t.tasksCollection.Find(ctx, filter)

	if err != nil {
		t.logger.Println(err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var response models.GetAllTasksResponse
	for cursor.Next(ctx) {
		var task models.Task
		if err := cursor.Decode(&task); err != nil {
			t.logger.Println(err)
			return nil, err
		}
		response.Tasks = append(response.Tasks, &task)
	}

	if err := cursor.Err(); err != nil {
		t.logger.Println(err)
		return nil, err
	}
	return &response, nil
}

func (t *TaskStorage) GetTaskById(ctx context.Context, req *models.ByIdRequest) (*models.Task, error) {
	filter := bson.M{
		"task_id": req.TaskId,
	}

	var response models.Task

	err := t.tasksCollection.FindOne(ctx, filter).Decode(&response)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &response, nil
}

func (t *TaskStorage) UpdateTask(req *models.Task) error {
	req.Method = models.UPDATE
	return t.SendTask(req)
}

func (t *TaskStorage) DeleteTaskById(ctx context.Context, req *models.ByIdRequest) error {
	resp, err := t.GetTaskById(ctx, req)
	if err != nil {
		return err
	}
	resp.Method = models.DELETE
	return t.SendTask(resp)
}
