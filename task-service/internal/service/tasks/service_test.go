package tasks

import (
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/Xanaduxan/tasks-golang/task-service/internal/events"
	authmocks "github.com/Xanaduxan/tasks-golang/task-service/internal/service/auth/mocks"
	groupmembersmocks "github.com/Xanaduxan/tasks-golang/task-service/internal/service/group_members/mocks"
	groupsmocks "github.com/Xanaduxan/tasks-golang/task-service/internal/service/groups/mocks"
	"github.com/Xanaduxan/tasks-golang/task-service/internal/service/tasks/mocks"
	storage2 "github.com/Xanaduxan/tasks-golang/task-service/internal/storage"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newTestService(
	t *testing.T,
) (
	*Service,
	*mocks.TaskInterface,
	*authmocks.UserInterface,
	*groupsmocks.GroupInterface,
	*groupmembersmocks.GroupMemberInterface,
	*mocks.Notifier,
) {
	t.Helper()

	tasksMock := mocks.NewTaskInterface(t)
	usersMock := authmocks.NewUserInterface(t)
	groupsMock := groupsmocks.NewGroupInterface(t)
	groupMembersMock := groupmembersmocks.NewGroupMemberInterface(t)
	notifierMock := mocks.NewNotifier(t)

	svc := NewService(tasksMock, usersMock, groupsMock, groupMembersMock, notifierMock)

	return svc, tasksMock, usersMock, groupsMock, groupMembersMock, notifierMock
}

func TestService_CreateTask_SuccessWithoutGroup(t *testing.T) {
	svc, tasksMock, usersMock, groupsMock, groupMembersMock, _ := newTestService(t)

	userID := uuid.New()
	deadline := time.Now().Add(24 * time.Hour)

	usersMock.
		On("GetByID", userID).
		Return(storage2.User{ID: userID}, nil).
		Once()

	tasksMock.
		On("Create", mock.MatchedBy(func(task storage2.Task) bool {
			return task.UserID == userID &&
				task.Name == "test task" &&
				task.GroupID == nil &&
				task.Status == storage2.StatusCreated &&
				task.Deadline != nil
		})).
		Return(nil).
		Once()

	taskID, err := svc.CreateTask(userID, "test task", &deadline, nil)

	require.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, taskID)

	groupsMock.AssertNotCalled(t, "GetByID", mock.Anything)
	groupMembersMock.AssertNotCalled(t, "IsMember", mock.Anything, mock.Anything)
}

func TestService_CreateTask_EmptyName(t *testing.T) {
	svc, tasksMock, usersMock, groupsMock, groupMembersMock, _ := newTestService(t)

	taskID, err := svc.CreateTask(uuid.New(), "", nil, nil)

	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidInput)
	assert.Equal(t, uuid.Nil, taskID)

	usersMock.AssertNotCalled(t, "GetByID", mock.Anything)
	tasksMock.AssertNotCalled(t, "Create", mock.Anything)
	groupsMock.AssertNotCalled(t, "GetByID", mock.Anything)
	groupMembersMock.AssertNotCalled(t, "IsMember", mock.Anything, mock.Anything)
}

func TestService_CreateTask_UserNotFound(t *testing.T) {
	svc, tasksMock, usersMock, _, _, _ := newTestService(t)

	userID := uuid.New()

	usersMock.
		On("GetByID", userID).
		Return(storage2.User{}, sql.ErrNoRows).
		Once()

	taskID, err := svc.CreateTask(userID, "task", nil, nil)

	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNotFound)
	assert.Equal(t, uuid.Nil, taskID)

	tasksMock.AssertNotCalled(t, "Create", mock.Anything)
}

func TestService_CreateTask_UserNotInGroup(t *testing.T) {
	svc, tasksMock, usersMock, groupsMock, groupMembersMock, _ := newTestService(t)

	userID := uuid.New()
	groupID := uuid.New()
	deadline := time.Now().Add(24 * time.Hour)

	usersMock.
		On("GetByID", userID).
		Return(storage2.User{ID: userID}, nil).
		Once()

	groupsMock.
		On("GetByID", groupID).
		Return(storage2.Group{ID: groupID}, nil).
		Once()

	groupMembersMock.
		On("IsMember", groupID, userID).
		Return(false, nil).
		Once()

	taskID, err := svc.CreateTask(userID, "task", &deadline, &groupID)

	require.Error(t, err)
	assert.ErrorIs(t, err, ErrForbidden)
	assert.Equal(t, uuid.Nil, taskID)

	tasksMock.AssertNotCalled(t, "Create", mock.Anything)
}

func TestService_GetTask_Success(t *testing.T) {
	svc, tasksMock, usersMock, _, _, _ := newTestService(t)

	userID := uuid.New()
	taskID := uuid.New()

	expected := storage2.Task{
		ID:     taskID,
		UserID: userID,
		Name:   "task",
		Status: storage2.StatusCreated,
	}

	usersMock.
		On("GetByID", userID).
		Return(storage2.User{ID: userID}, nil).
		Once()

	tasksMock.
		On("GetByID", taskID).
		Return(expected, nil).
		Once()

	tasksMock.
		On("HasAccess", taskID, userID).
		Return(true, nil).
		Once()

	got, err := svc.GetTask(userID, taskID)

	require.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestService_GetTask_Forbidden(t *testing.T) {
	svc, tasksMock, usersMock, _, _, _ := newTestService(t)

	userID := uuid.New()
	taskID := uuid.New()

	usersMock.
		On("GetByID", userID).
		Return(storage2.User{ID: userID}, nil).
		Once()

	tasksMock.
		On("GetByID", taskID).
		Return(storage2.Task{ID: taskID, UserID: uuid.New()}, nil).
		Once()

	tasksMock.
		On("HasAccess", taskID, userID).
		Return(false, nil).
		Once()

	_, err := svc.GetTask(userID, taskID)

	require.Error(t, err)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestService_GetTaskForWorker_Success(t *testing.T) {
	svc, tasksMock, _, _, _, _ := newTestService(t)

	taskID := uuid.New()
	expected := storage2.Task{
		ID:     taskID,
		Name:   "worker-task",
		Status: storage2.StatusCreated,
	}

	tasksMock.
		On("GetByID", taskID).
		Return(expected, nil).
		Once()

	got, err := svc.GetTaskForWorker(taskID)

	require.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestService_DeleteTask_Success(t *testing.T) {
	svc, tasksMock, usersMock, _, _, _ := newTestService(t)

	userID := uuid.New()
	taskID := uuid.New()

	usersMock.
		On("GetByID", userID).
		Return(storage2.User{ID: userID}, nil).
		Once()

	tasksMock.
		On("GetByID", taskID).
		Return(storage2.Task{
			ID:     taskID,
			UserID: userID,
			Name:   "task",
		}, nil).
		Once()

	tasksMock.
		On("DeleteByID", taskID).
		Return(nil).
		Once()

	err := svc.DeleteTask(userID, taskID)

	require.NoError(t, err)
}

func TestService_DeleteTask_Forbidden(t *testing.T) {
	svc, tasksMock, usersMock, _, _, _ := newTestService(t)

	userID := uuid.New()
	taskID := uuid.New()
	ownerID := uuid.New()

	usersMock.
		On("GetByID", userID).
		Return(storage2.User{ID: userID}, nil).
		Once()

	tasksMock.
		On("GetByID", taskID).
		Return(storage2.Task{
			ID:     taskID,
			UserID: ownerID,
			Name:   "task",
		}, nil).
		Once()

	err := svc.DeleteTask(userID, taskID)

	require.Error(t, err)
	assert.ErrorIs(t, err, ErrForbidden)

	tasksMock.AssertNotCalled(t, "DeleteByID", mock.Anything)
}

func TestService_UpdateTask_Success(t *testing.T) {
	svc, tasksMock, usersMock, groupsMock, groupMembersMock, _ := newTestService(t)

	userID := uuid.New()
	taskID := uuid.New()
	groupID := uuid.New()
	deadline := time.Now().Add(48 * time.Hour)

	usersMock.
		On("GetByID", userID).
		Return(storage2.User{ID: userID}, nil).
		Once()

	groupsMock.
		On("GetByID", groupID).
		Return(storage2.Group{ID: groupID}, nil).
		Once()

	groupMembersMock.
		On("IsMember", groupID, userID).
		Return(true, nil).
		Once()

	existingTask := storage2.Task{
		ID:      taskID,
		UserID:  userID,
		Name:    "old",
		Status:  storage2.StatusCreated,
		GroupID: nil,
	}

	tasksMock.
		On("GetByID", taskID).
		Return(existingTask, nil).
		Once()

	tasksMock.
		On("HasAccess", taskID, userID).
		Return(true, nil).
		Once()

	tasksMock.
		On("Update", mock.MatchedBy(func(task storage2.Task) bool {
			return task.ID == taskID &&
				task.Name == "new name" &&
				task.GroupID != nil &&
				*task.GroupID == groupID &&
				task.Deadline != nil
		})).
		Return(nil).
		Once()

	err := svc.UpdateTask(userID, taskID, "new name", &deadline, &groupID)

	require.NoError(t, err)
}

func TestService_UpdateTask_EmptyName(t *testing.T) {
	svc, tasksMock, usersMock, groupsMock, groupMembersMock, _ := newTestService(t)

	err := svc.UpdateTask(uuid.New(), uuid.New(), "", nil, nil)

	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidInput)

	tasksMock.AssertNotCalled(t, "Update", mock.Anything)
	usersMock.AssertNotCalled(t, "GetByID", mock.Anything)
	groupsMock.AssertNotCalled(t, "GetByID", mock.Anything)
	groupMembersMock.AssertNotCalled(t, "IsMember", mock.Anything, mock.Anything)
}

func TestService_UpdateTaskStatus_SuccessWithNotifier(t *testing.T) {
	svc, tasksMock, _, _, _, notifierMock := newTestService(t)

	taskID := uuid.New()
	userID := uuid.New()
	groupID := uuid.New()

	task := storage2.Task{
		ID:      taskID,
		UserID:  userID,
		GroupID: &groupID,
		Name:    "task",
		Status:  storage2.StatusCreated,
	}

	tasksMock.
		On("GetByID", taskID).
		Return(task, nil).
		Once()

	tasksMock.
		On("UpdateStatus", taskID, storage2.StatusDone).
		Return(nil).
		Once()

	notifierMock.
		On("NotifyTaskStatusUpdated", mock.MatchedBy(func(e events.TaskStatusUpdated) bool {
			return e.TaskID == taskID &&
				e.UserID == userID &&
				e.GroupID != nil &&
				*e.GroupID == groupID &&
				e.Status == storage2.StatusDone
		})).
		Return(nil).
		Once()

	err := svc.UpdateTaskStatus(taskID, storage2.StatusDone)

	require.NoError(t, err)
}

func TestService_UpdateTaskStatus_InvalidStatus(t *testing.T) {
	svc, tasksMock, _, _, _, notifierMock := newTestService(t)

	err := svc.UpdateTaskStatus(uuid.New(), storage2.TaskStatus("bad-status"))

	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidInput)

	tasksMock.AssertNotCalled(t, "GetByID", mock.Anything)
	tasksMock.AssertNotCalled(t, "UpdateStatus", mock.Anything, mock.Anything)
	notifierMock.AssertNotCalled(t, "NotifyTaskStatusUpdated", mock.Anything)
}

func TestService_UpdateTaskStatus_TaskNotFound(t *testing.T) {
	svc, tasksMock, _, _, _, notifierMock := newTestService(t)

	taskID := uuid.New()

	tasksMock.
		On("GetByID", taskID).
		Return(storage2.Task{}, sql.ErrNoRows).
		Once()

	err := svc.UpdateTaskStatus(taskID, storage2.StatusDone)

	require.Error(t, err)
	assert.ErrorIs(t, err, ErrNotFound)

	tasksMock.AssertNotCalled(t, "UpdateStatus", mock.Anything, mock.Anything)
	notifierMock.AssertNotCalled(t, "NotifyTaskStatusUpdated", mock.Anything)
}

func TestService_UpdateTaskStatus_NotifierError(t *testing.T) {
	svc, tasksMock, _, _, _, notifierMock := newTestService(t)

	taskID := uuid.New()
	userID := uuid.New()

	task := storage2.Task{
		ID:     taskID,
		UserID: userID,
		Name:   "task",
		Status: storage2.StatusCreated,
	}

	tasksMock.
		On("GetByID", taskID).
		Return(task, nil).
		Once()

	tasksMock.
		On("UpdateStatus", taskID, storage2.StatusDone).
		Return(nil).
		Once()

	notifierMock.
		On("NotifyTaskStatusUpdated", mock.Anything).
		Return(errors.New("notify failed")).
		Once()

	err := svc.UpdateTaskStatus(taskID, storage2.StatusDone)

	require.Error(t, err)
	assert.EqualError(t, err, "notify failed")
}

func TestService_GetAllNotDone(t *testing.T) {
	svc, tasksMock, _, _, _, _ := newTestService(t)

	expected := []storage2.Task{
		{ID: uuid.New(), Name: "task-1", Status: storage2.StatusCreated},
		{ID: uuid.New(), Name: "task-2", Status: storage2.StatusInProgress},
	}

	tasksMock.
		On("GetAllNotDone").
		Return(expected, nil).
		Once()

	got, err := svc.GetAllNotDone()

	require.NoError(t, err)
	assert.Equal(t, expected, got)
}
