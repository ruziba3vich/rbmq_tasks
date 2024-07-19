package consumer

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

type ConsumerStorage struct {
	consumerId          primitive.ObjectID
	consumersCollection *mongo.Collection
	tasksCollection     *mongo.Collection
	logger              *log.Logger
	channel             *amqp.Channel
	queue               amqp.Queue
}

func New(cid primitive.ObjectID, c *mongo.Collection, t *mongo.Collection, l *log.Logger, ch *amqp.Channel, q amqp.Queue) *ConsumerStorage {
	return &ConsumerStorage{
		consumerId:          cid,
		consumersCollection: c,
		tasksCollection:     t,
		logger:              l,
		channel:             ch,
		queue:               q,
	}
}

func (c *ConsumerStorage) ListenToProducer(ctx context.Context) error {
	if err := c.channel.QueueBind(
		c.queue.Name,
		"",
		"prodonik",
		false,
		nil,
	); err != nil {
		c.logger.Println(err)
		return err
	}

	messages, err := c.channel.Consume(
		c.queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		c.logger.Println(err)
		return err
	}

	for message := range messages {
		var task models.Task
		if err := json.Unmarshal(message.Body, &task); err == nil {
			var err error
			if task.Method == models.POST {
				err = c.postTask(ctx, &task)
			} else if task.Method == models.UPDATE {
				err = c.updateTask(ctx, &task)
			} else if task.Method == models.DELETE {
				err = c.deleteTask(ctx, &task)
			} else {
				c.logger.Println("no matching methods found!")
			}
			if err != nil {
				c.logger.Println(err)
				message.Nack(false, true)
			} else {
				message.Ack(false)
			}
		}
	}
	return nil
}

func (t *ConsumerStorage) postTask(ctx context.Context, req *models.Task) error {
	pureTask := req.ConvertToPureTask()
	_, err := t.tasksCollection.InsertOne(ctx, pureTask)
	return err
}

func (t *ConsumerStorage) updateTask(ctx context.Context, req *models.Task) error {
	filter := bson.M{"task_id": req.TaskId}

	fields := bson.M{}

	if len(req.TaskContent) > 0 {
		fields["task_content"] = req.TaskContent
	}
	if len(req.TaskStatus) > 0 {
		fields["task_status"] = req.TaskStatus
	}

	if len(fields) == 0 {
		return errors.New("no field given to be updated")
	}
	query := bson.M{
		"$set": fields,
	}
	result, err := t.tasksCollection.UpdateOne(ctx, filter, query)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (t *ConsumerStorage) deleteTask(ctx context.Context, req *models.Task) error {
	filter := bson.M{
		"task_id": req.TaskId,
	}
	result, err := t.tasksCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
