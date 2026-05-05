package handler

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"idp-backend/internal/domain"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ProjectHandler struct {
	service domain.ProjectService
}

func NewProjectHandler(service domain.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}

type CreateProjectRequest struct {
	Name   string `json:"name" binding:"required"`
	GitURL string `json:"git_url" binding:"required"`
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var req CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project, err := h.service.CreateProject(req.Name, req.GitURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, project)
}

func (h *ProjectHandler) GetProjects(c *gin.Context) {
	projects, err := h.service.GetProjects()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projects)
}

func (h *ProjectHandler) GetProjectByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	project, err := h.service.GetProjectByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}

	c.JSON(http.StatusOK, project)
}

func (h *ProjectHandler) GetProjectLogs(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	logs, err := h.service.GetProjectLogs(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Split logs into lines and format as expected by frontend
	lines := strings.Split(logs, "\n")
	var logEntries []gin.H
	for i, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		logEntries = append(logEntries, gin.H{
			"id":        fmt.Sprintf("%d", i),
			"timestamp": time.Now().Format(time.RFC3339), // Fallback timestamp
			"content":   line,
			"level":     "info",
		})
	}

	c.JSON(http.StatusOK, logEntries)
}
