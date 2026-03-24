package handlers

import (
	"context"
	"errors"
	"time"

	"github.com/Xanaduxan/tasks-golang/internal/service/tasks"
	"github.com/Xanaduxan/tasks-golang/internal/transport/grpc/mapper"
	taskv1 "github.com/Xanaduxan/tasks-golang/pkg/pb/task/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TaskHandler struct {
	taskv1.UnimplementedTaskServiceServer
	service *tasks.Service
}

func NewTaskHandler(service *tasks.Service) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

func mapError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, tasks.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, tasks.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, tasks.ErrForbidden):
		return status.Error(codes.PermissionDenied, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}

func (h *TaskHandler) CreateTask(
	ctx context.Context,
	req *taskv1.CreateTaskRequest,
) (*taskv1.CreateTaskResponse, error) {

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

	return &taskv1.CreateTaskResponse{
		TaskId: taskID.String(),
	}, nil
}

func (h *TaskHandler) GetTask(
	ctx context.Context,
	req *taskv1.GetTaskRequest,
) (*taskv1.GetTaskResponse, error) {

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

	return &taskv1.GetTaskResponse{
		Task: mapper.TaskToProto(task),
	}, nil
}
func (h *TaskHandler) ListTasks(
	ctx context.Context,
	req *taskv1.ListTasksRequest,
) (*taskv1.ListTasksResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, mapError(err)
	}

	tasks, err := h.service.ListTasks(userID)
	if err != nil {
		return nil, mapError(err)
	}

	resp := &taskv1.ListTasksResponse{
		Tasks: make([]*taskv1.Task, 0, len(tasks)),
	}

	for _, task := range tasks {
		resp.Tasks = append(resp.Tasks, mapper.TaskToProto(task))
	}

	return resp, nil
}
func (h *TaskHandler) UpdateTask(
	ctx context.Context,
	req *taskv1.UpdateTaskRequest,
) (*taskv1.UpdateTaskResponse, error) {
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

	return &taskv1.UpdateTaskResponse{}, nil
}

func (h *TaskHandler) DeleteTask(
	ctx context.Context,
	req *taskv1.DeleteTaskRequest,
) (*taskv1.DeleteTaskResponse, error) {
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

	return &taskv1.DeleteTaskResponse{}, nil
}
func (h *TaskHandler) SearchTasks(
	ctx context.Context,
	req *taskv1.SearchTasksRequest,
) (*taskv1.SearchTasksResponse, error) {

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, err
	}

	tasks, err := h.service.SearchTasks(userID, req.Query)
	if err != nil {
		return nil, err
	}

	resp := &taskv1.SearchTasksResponse{
		Tasks: make([]*taskv1.Task, 0, len(tasks)),
	}

	for _, t := range tasks {
		resp.Tasks = append(resp.Tasks, mapper.TaskToProto(t))
	}

	return resp, nil
}
