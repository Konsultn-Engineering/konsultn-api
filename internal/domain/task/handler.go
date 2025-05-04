package task

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	repo *Repository[Task] // use the interface for flexibility
}

func NewHandler(repo *Repository[Task]) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) CreateTask(ctx *gin.Context) {
	var task Task
	if err := ctx.ShouldBindJSON(&task); err != nil {
		// If there is an error (invalid JSON or missing fields), return an error response
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.repo.Save(task)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.repo.Preload(&task, []string{"Assignee"}, "id", task.ID)

	ctx.JSON(http.StatusOK, gin.H{
		"message": "task created successfully",
		"task":    task,
	})

}

func (h *Handler) GetTaskById(ctx *gin.Context) {
	taskId := ctx.Param("id")

	task, err := h.repo.FindById(taskId)
	h.repo.Preload(&task, []string{"Assignee"}, "id", task.ID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, task)
}
