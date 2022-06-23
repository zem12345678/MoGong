package task_jobs

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type TaskModel struct {
	Task *task.TaskModel
	store *TaskStore
}

type TaskStore struct {
	db     *mongo.Collection

}

type TaskModelFn interface {
	
}


func