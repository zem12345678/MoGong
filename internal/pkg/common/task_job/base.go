package task_jobs

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskStore struct {
	db *mongo.Collection
}

type TaskModelFn interface {
}

func Save() {

}

func CallBack() {

}
