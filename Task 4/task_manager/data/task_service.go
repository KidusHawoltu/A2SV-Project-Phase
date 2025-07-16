package data

import (
	"A2SV_ProjectPhase/Task4/TaskManager/models"
	"fmt"
	"sync"
)

type TaskList struct {
	Tasks  map[int]*models.Task
	nextId int
	// Used RWMutex so that readers doesn't lock eachother
	mu sync.RWMutex
}

type TaskManager interface {
	GetTasks() []*models.Task
	GetTaskById(id int) (*models.Task, error)
	UpdateTask(id int, updatedTask models.Task) (*models.Task, error)
	DeleteTask(id int) error
	AddTask(task models.Task) *models.Task
}

func NewTaskManager() TaskManager {
	return &TaskList{
		Tasks:  make(map[int]*models.Task),
		nextId: 1,
		// no need to initialize mu
	}
}

func (tl *TaskList) GetTasks() []*models.Task {
	tl.mu.RLock()
	defer tl.mu.RUnlock()

	tasks := make([]*models.Task, 0, len(tl.Tasks))
	for _, task := range tl.Tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

func (tl *TaskList) GetTaskById(id int) (*models.Task, error) {
	tl.mu.RLock()
	defer tl.mu.RUnlock()

	task, exists := tl.Tasks[id]
	if !exists {
		return nil, fmt.Errorf("there is no task with id %v", id)
	}
	return task, nil
}

func (tl *TaskList) UpdateTask(id int, updatedTask models.Task) (*models.Task, error) {
	tl.mu.Lock()
	defer tl.mu.Unlock()

	task, exists := tl.Tasks[id]
	if !exists {
		return nil, fmt.Errorf("there is no task with id %v", id)
	}
	task.Title = updatedTask.Title
	task.Description = updatedTask.Description
	task.DueDate = updatedTask.DueDate
	task.Status = updatedTask.Status
	return task, nil
}

func (tl *TaskList) DeleteTask(id int) error {
	tl.mu.Lock()
	defer tl.mu.Unlock()

	if _, exists := tl.Tasks[id]; !exists {
		return fmt.Errorf("there is no task with id %v", id)
	}
	delete(tl.Tasks, id)
	return nil
}

func (tl *TaskList) AddTask(task models.Task) *models.Task {
	tl.mu.Lock()
	defer tl.mu.Unlock()

	newTask := task
	newTask.Id = tl.nextId
	tl.Tasks[tl.nextId] = &newTask
	tl.nextId++
	return &newTask
}
