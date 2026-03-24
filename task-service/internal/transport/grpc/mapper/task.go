package mapper

import (
	"time"

	"github.com/Xanaduxan/tasks-golang/task-service/internal/storage"
	taskv1 "github.com/Xanaduxan/tasks-golang/task-service/pkg/pb/task/v1"
)

func TaskToProto(task storage.Task) *taskv1.Task {
	var deadline string
	if task.Deadline != nil {
		deadline = task.Deadline.Format(time.RFC3339)
	}

	var groupID string
	if task.GroupID != nil {
		groupID = task.GroupID.String()
	}

	return &taskv1.Task{
		Id:       task.ID.String(),
		Name:     task.Name,
		Status:   string(task.Status),
		UserId:   task.UserID.String(),
		GroupId:  groupID,
		Deadline: deadline,
	}
}
