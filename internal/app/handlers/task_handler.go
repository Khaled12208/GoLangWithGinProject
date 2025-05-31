package handlers

import (
	"github.com/gin-gonic/gin"
	"golangwithgin/internal/domain"
	"net/http"
	"strconv"
)

type TaskHandler struct {
	taskService domain.TaskService
}

func NewTaskHandler(taskService domain.TaskService) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
	}
}

// @Summary Create a new task
// @Description Submit a new task for processing
// @Tags tasks
// @Accept json
// @Produce json
// @Security Bearer
// @Param task body domain.SwaggerTask true "Task details"
// @Success 202 {object} domain.SwaggerTask
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /tasks [post]
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task domain.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.taskService.SubmitTask(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, task)
}

// @Summary Get task by ID
// @Description Get a task's details and status by its ID
// @Tags tasks
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Task ID"
// @Success 200 {object} domain.SwaggerTask
// @Failure 400 {object} domain.ErrorResponse
// @Failure 401 {object} domain.ErrorResponse
// @Failure 404 {object} domain.ErrorResponse
// @Router /tasks/{id} [get]
func (h *TaskHandler) GetTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	task, err := h.taskService.GetTaskStatus(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// @Summary Get all tasks
// @Description Get a list of all tasks
// @Tags tasks
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {array} domain.SwaggerTask
// @Failure 401 {object} domain.ErrorResponse
// @Failure 500 {object} domain.ErrorResponse
// @Router /tasks [get]
func (h *TaskHandler) GetAllTasks(c *gin.Context) {
	tasks, err := h.taskService.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
} 