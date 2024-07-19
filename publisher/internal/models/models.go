package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type (
	Status string

	Publisher struct {
		Id       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
		FullName string             `json:"full_name" bson:"full_name"`
	}

	Method string

	Task struct {
		TaskId      primitive.ObjectID `json:"task_id" bson:"task_id"`
		TaskBy      primitive.ObjectID `json:"task_by" bson:"task_by"`
		TaskFor     primitive.ObjectID `json:"task_for" bson:"task_for"`
		TaskContent string             `json:"task_content" bson:"task_content"`
		TaskStatus  Status             `json:"task_status" bson:"task_status"`
		Method      Method             `json:"method" bson:"method"`
	}

	GetAllTasksResponse struct {
		Tasks []*Task
	}

	ByIdRequest struct {
		TaskId primitive.ObjectID `json:"task_id" bson:"task_id"`
	}

	PureTask struct {
		TaskId      primitive.ObjectID `json:"task_id" bson:"task_id"`
		TaskBy      primitive.ObjectID `json:"task_by" bson:"task_by"`
		TaskFor     primitive.ObjectID `json:"task_for" bson:"task_for"`
		TaskContent string             `json:"task_content" bson:"task_content"`
		TaskStatus  Status             `json:"task_status" bson:"task_status"`
	}
)

func (t *Task) ConvertToPureTask() *PureTask {
	return &PureTask{
		TaskId:      t.TaskId,
		TaskBy:      t.TaskBy,
		TaskFor:     t.TaskFor,
		TaskContent: t.TaskContent,
		TaskStatus:  t.TaskStatus,
	}
}

const (
	POST   Method = "POST"
	UPDATE Method = "UPDATE"
	DELETE Method = "DELETE"
	READ   Method = "READ"
)
