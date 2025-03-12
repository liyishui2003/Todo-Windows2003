package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"todolist/internal/model"
	"todolist/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskHandler struct {
	taskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

// 增
func (h *TaskHandler) CreateTask(ctx *gin.Context) {
	var req model.CreateTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fmt.Printf("参数解析失败: %v", err)
		fmt.Printf("请求体内容: %+v", req)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	// 从 Context 中获取 user_id
	userID := ctx.MustGet("user_id").(float64)

	// 创建任务
	task := &model.Task{
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		UserID:      userID,
		Importance:  req.Importance,
		Urgency:     req.Urgency,
	}
	if err := h.taskService.CreateTask(task); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "创建任务失败"})
		return
	}

	ctx.JSON(http.StatusOK, task)
}

// 删
func (h *TaskHandler) DeleteTask(ctx *gin.Context) {
	// 从 URL 路径中提取任务 ID
	taskID := ctx.Param("taskID")
	id, err := strconv.Atoi(taskID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// 调用 Service 层删除任务
	err = h.taskService.DeleteTask(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

// 改
func (h *TaskHandler) EditTask(ctx *gin.Context) {
	// 从 URL 路径中提取任务 ID
	taskID := ctx.Param("taskID")
	id, err := strconv.Atoi(taskID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// 解析请求体
	var req model.EditTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	// 调用 Service 层更新任务信息
	updatedTask, err := h.taskService.EditTask(uint(id), &req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// 返回更新后的任务
	ctx.JSON(http.StatusOK, updatedTask)
}

// 查
func (h *TaskHandler) GetTask(ctx *gin.Context) {
	taskID := ctx.Param("taskID")
	id, err := strconv.Atoi(taskID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid task ID"})
		return
	}

	task, err := h.taskService.GetTask(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	ctx.JSON(http.StatusOK, task)
}

// 根据优先级查
func (h *TaskHandler) GetTasksByPriority(ctx *gin.Context) {
	// 从查询参数中获取 importance 和 urgency
	importanceStr := ctx.DefaultQuery("importance", "-1")
	urgencyStr := ctx.DefaultQuery("urgency", "-1")

	// 转换为 int 类型
	importance, err := strconv.Atoi(importanceStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid importance value"})
		return
	}
	urgency, err := strconv.Atoi(urgencyStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid urgency value"})
		return
	}

	// 调用 Service 层获取任务
	tasks, err := h.taskService.GetTasksByPriority(importance, urgency)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 返回任务列表
	ctx.JSON(http.StatusOK, tasks)
}

// 根据userid查
func (c *TaskHandler) GetTasksByUserID(ctx *gin.Context) {
	// 从请求中获取 user_id
	userID := ctx.Param("userID")
	id, err := strconv.Atoi(userID)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userID"})
		return
	}
	// 调用 Service 层的方法
	tasks, err := c.taskService.GetTasksByUserID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get tasks",
		})
		return
	}

	// 返回成功响应
	ctx.JSON(http.StatusOK, gin.H{
		"tasks": tasks,
	})
}
