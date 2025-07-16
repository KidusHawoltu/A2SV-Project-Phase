package data

import (
	"A2SV_ProjectPhase/Task4/TaskManager/models"
	"fmt"
)

type TaskList struct {
	Tasks  map[int]*models.Task
	nextId int
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
	}
}

func (taskList *TaskList) GetTasks() []*models.Task {
	tasks := make([]*models.Task, 0, len(taskList.Tasks))
	for _, task := range taskList.Tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

func (taskList *TaskList) GetTaskById(id int) (*models.Task, error) {
	task, exists := taskList.Tasks[id]
	if !exists {
		return nil, fmt.Errorf("there is no task with id %v", id)
	}
	return task, nil
}

func (taskList *TaskList) UpdateTask(id int, updatedTask models.Task) (*models.Task, error) {
	task, exists := taskList.Tasks[id]
	if !exists {
		return nil, fmt.Errorf("there is no task with id %v", id)
	}
	task.Title = updatedTask.Title
	task.Description = updatedTask.Description
	task.DueDate = updatedTask.DueDate
	task.Status = updatedTask.Status
	return task, nil
}

func (taskList *TaskList) DeleteTask(id int) error {
	if _, exists := taskList.Tasks[id]; !exists {
		return fmt.Errorf("there is no task with id %v", id)
	}
	delete(taskList.Tasks, id)
	return nil
}

func (taskList *TaskList) AddTask(task models.Task) *models.Task {
	newTask := task
	newTask.Id = taskList.nextId
	taskList.Tasks[taskList.nextId] = &newTask
	taskList.nextId++
	return &newTask
}
