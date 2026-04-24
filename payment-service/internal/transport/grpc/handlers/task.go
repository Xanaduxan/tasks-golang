package handlers

import (
	"context"
	"errors"
	"time"

	tasks2 "github.com/Xanaduxan/tasks-golang/task-service/internal/service/tasks"
	"github.com/Xanaduxan/tasks-golang/task-service/internal/transport/grpc/mapper"
	taskv2 "github.com/Xanaduxan/tasks-golang/task-service/pkg/pb/task/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TaskHandler struct {
	taskv2.UnimplementedTaskServiceServer
	service *tasks2.Service
}

func NewTaskHandler(service *tasks2.Service) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

func mapError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, tasks2.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, tasks2.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, tasks2.ErrForbidden):
		return status.Error(codes.PermissionDenied, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}

func (h *TaskHandler) CreateTask(
	ctx context.Context,
	req *taskv2.CreateTaskRequest,
) (*taskv2.CreateTaskResponse, error) {

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, mapError(err)
	}

	var groupID *uuid.UUID
	if req.GroupId != "" {
		gid, err := uuid.Parse(req.GroupId)
		if err != nil {
			return nil, mapError(err)
		}
		groupID = &gid
	}

	var deadline *time.Time
	if req.Deadline != "" {
		t, err := time.Parse(time.RFC3339, req.Deadline)
		if err != nil {
			return nil, mapError(err)
		}
		deadline = &t
	}

	taskID, err := h.service.CreateTask(userID, req.Name, deadline, groupID)
	if err != nil {
		return nil, mapError(err)
	}

	return &taskv2.CreateTaskResponse{
		TaskId: taskID.String(),
	}, nil
}

func (h *TaskHandler) GetTask(
	ctx context.Context,
	req *taskv2.GetTaskRequest,
) (*taskv2.GetTaskResponse, error) {

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, mapError(err)
	}

	taskID, err := uuid.Parse(req.TaskId)
	if err != nil {
		return nil, mapError(err)
	}

	task, err := h.service.GetTask(userID, taskID)
	if err != nil {
		return nil, mapError(err)
	}

	return &taskv2.GetTaskResponse{
		Task: mapper.TaskToProto(task),
	}, nil
}
func (h *TaskHandler) ListTasks(
	ctx context.Context,
	req *taskv2.ListTasksRequest,
) (*taskv2.ListTasksResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, mapError(err)
	}

	tasks, err := h.service.ListTasks(userID)
	if err != nil {
		return nil, mapError(err)
	}

	resp := &taskv2.ListTasksResponse{
		Tasks: make([]*taskv2.Task, 0, len(tasks)),
	}

	for _, task := range tasks {
		resp.Tasks = append(resp.Tasks, mapper.TaskToProto(task))
	}

	return resp, nil
}
func (h *TaskHandler) UpdateTask(
	ctx context.Context,
	req *taskv2.UpdateTaskRequest,
) (*taskv2.UpdateTaskResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	taskID, err := uuid.Parse(req.TaskId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid task_id")
	}

	var groupID *uuid.UUID
	if req.GroupId != "" {
		gid, err := uuid.Parse(req.GroupId)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid group_id")
		}
		groupID = &gid
	}

	var deadline *time.Time
	if req.Deadline != "" {
		t, err := time.Parse(time.RFC3339, req.Deadline)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid deadline")
		}
		deadline = &t
	}

	err = h.service.UpdateTask(userID, taskID, req.Name, deadline, groupID)
	if err != nil {
		return nil, mapError(err)
	}

	return &taskv2.UpdateTaskResponse{}, nil
}

func (h *TaskHandler) DeleteTask(
	ctx context.Context,
	req *taskv2.DeleteTaskRequest,
) (*taskv2.DeleteTaskResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id")
	}

	taskID, err := uuid.Parse(req.TaskId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid task_id")
	}

	err = h.service.DeleteTask(userID, taskID)
	if err != nil {
		return nil, mapError(err)
	}

	return &taskv2.DeleteTaskResponse{}, nil
}
func (h *TaskHandler) SearchTasks(
	ctx context.Context,
	req *taskv2.SearchTasksRequest,
) (*taskv2.SearchTasksResponse, error) {

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	tasks, err := h.service.SearchTasks(userID, req.Query)
	if err != nil {
		return nil, err
	}

	resp := &taskv2.SearchTasksResponse{
		Tasks: make([]*taskv2.Task, 0, len(tasks)),
	}

	for _, t := range tasks {
		resp.Tasks = append(resp.Tasks, mapper.TaskToProto(t))
	}

	return resp, nil
}
