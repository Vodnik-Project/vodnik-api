package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/Vodnik-Project/vodnik-api/db/sqlc"
	log "github.com/Vodnik-Project/vodnik-api/logger"
	"github.com/Vodnik-Project/vodnik-api/util"
	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
)

type createTaskRequest struct {
	Title     string `json:"title"`
	Info      string `json:"info"`
	Tag       string `json:"tag"`
	Beggining string `json:"beggining"`
	Deadline  string `json:"deadline"`
	Color     string `json:"color"`
}

type TaskDataResponse struct {
	TaskID    string `json:"task_id"`
	ProjectID string `json:"project_id"`
	Title     string `json:"title"`
	Info      string `json:"info"`
	Tag       string `json:"tag"`
	CreatedBy string `json:"created_by"`
	CreatedAt string `json:"created_at"`
	Beggining string `json:"beggining"`
	Deadline  string `json:"deadline"`
	Color     string `json:"color"`
}

func (s Server) CreateTask(c echo.Context) error {
	var task createTaskRequest
	err := c.Bind(&task)
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("can't parse input data")
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"message": "invalid input",
			"traceid": traceid,
		})
	}
	err = util.CheckEmpty(task, []string{"Title"})
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg(err.Error())
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"message": err.Error(),
			"traceid": traceid,
		})
	}
	userUUID := c.Get("userUUID").(uuid.UUID)
	var beggining time.Time
	if task.Beggining != "" {
		beggining, err = time.Parse(time.RFC3339, task.Beggining)
		if err != nil {
			traceid := util.RandomString(8)
			log.Logger.Err(err).Str("traceid", traceid).Msg("can't parse input data")
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{
				"message": "invalid beggining time format. time format must be RFC3339.",
				"traceid": traceid,
			})
		}
	}
	var deadline time.Time
	if task.Deadline != "" {
		deadline, err = time.Parse(time.RFC3339, task.Deadline)
		if err != nil {
			traceid := util.RandomString(8)
			log.Logger.Err(err).Str("traceid", traceid).Msg("can't parse input data")
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{
				"message": "invalid deadline time format. time format must be RFC3339.",
				"traceid": traceid,
			})
		}
	}
	err = s.store.CreateTaskTx(c, sqlc.CreateTaskParams{
		ProjectID: c.Get("projectUUID").(uuid.UUID),
		Title:     task.Title,
		Info:      sql.NullString{String: task.Info, Valid: true},
		Tag:       sql.NullString{String: task.Tag, Valid: true},
		CreatedBy: userUUID,
		Beggining: sql.NullTime{Time: beggining, Valid: true},
		Deadline:  sql.NullTime{Time: deadline, Valid: true},
		Color:     sql.NullString{String: task.Color, Valid: true},
	})
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "an error occurred while processing your request",
			"traceid": traceid,
		})
	}
	createdTask := c.Get("task").(sqlc.Task)
	responseData := TaskDataResponse{
		TaskID:    createdTask.TaskID.String(),
		ProjectID: createdTask.ProjectID.String(),
		Title:     createdTask.Title,
		Info:      createdTask.Info.String,
		Tag:       createdTask.Tag.String,
		CreatedBy: userUUID.String(),
		CreatedAt: createdTask.CreatedAt.Time.Format(time.RFC3339),
		Beggining: createdTask.Beggining.Time.Format(time.RFC3339),
		Deadline:  createdTask.Deadline.Time.Format(time.RFC3339),
		Color:     createdTask.Color.String,
	}
	log.Logger.Info().Msgf("task created: %+v", responseData)
	return c.JSON(http.StatusOK, echo.Map{
		"message": "task created successfully",
		"task":    responseData,
	})
}

func (s Server) GetTaskData(c echo.Context) error {
	ctx := c.Request().Context()
	taskID := c.Param("taskid")
	taskUUID, err := uuid.FromString(taskID)
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("invalid taskid")
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"message": "invalid taskID",
			"traceid": traceid,
		})
	}
	task, err := s.store.GetTaskData(ctx, taskUUID)
	if err != nil {
		traceid := util.RandomString(8)
		if err == sql.ErrNoRows {
			log.Logger.Err(err).Str("traceid", traceid).Msg("")
			return c.JSON(http.StatusNotFound, echo.Map{
				"message": "no task found",
				"traceid": traceid,
			})
		}
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "an error occurred while processing your request",
			"traceid": traceid,
		})
	}
	responseData := TaskDataResponse{
		TaskID:    task.TaskID.String(),
		ProjectID: task.ProjectID.String(),
		Title:     task.Title,
		Info:      task.Info.String,
		Tag:       task.Tag.String,
		CreatedBy: task.CreatedBy.String(),
		CreatedAt: task.CreatedAt.Time.Format(time.RFC3339),
		Beggining: task.Beggining.Time.Format(time.RFC3339),
		Deadline:  task.Deadline.Time.Format(time.RFC3339),
		Color:     task.Color.String,
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "task found",
		"task":    responseData,
	})
}

func (s Server) GetTasksByProjectID(c echo.Context) error {
	return nil
}

func (s Server) UpdateTask(c echo.Context) error {
	ctx := c.Request().Context()
	var updateTaskData createTaskRequest
	err := c.Bind(&updateTaskData)
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("can't parse input data")
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"message": "invalid input",
			"traceid": traceid,
		})
	}
	if updateTaskData == (createTaskRequest{}) {
		traceid := util.RandomString(8)
		log.Logger.Err(errors.New("input data is empty")).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"message": "input data is empty",
			"traceid": traceid,
		})
	}
	var beggining time.Time
	if updateTaskData.Beggining != "" {
		beggining, err = time.Parse(time.RFC3339, updateTaskData.Beggining)
		if err != nil {
			traceid := util.RandomString(8)
			log.Logger.Err(err).Str("traceid", traceid).Msg("can't parse input data")
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{
				"message": "invalid beggining time format. time format must be RFC3339.",
				"traceid": traceid,
			})
		}
	}
	var deadline time.Time
	if updateTaskData.Deadline != "" {
		deadline, err = time.Parse(time.RFC3339, updateTaskData.Deadline)
		if err != nil {
			traceid := util.RandomString(8)
			log.Logger.Err(err).Str("traceid", traceid).Msg("can't parse input data")
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{
				"message": "invalid deadline time format. time format must be RFC3339.",
				"traceid": traceid,
			})
		}
	}
	taskID := c.Param("taskid")
	taskUUID, err := uuid.FromString(taskID)
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("invalid taskid")
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"message": "invalid taskID",
			"traceid": traceid,
		})
	}
	updatedTask, err := s.store.UpdateTask(ctx, sqlc.UpdateTaskParams{
		Title:     updateTaskData.Title,
		Info:      updateTaskData.Info,
		Tag:       updateTaskData.Tag,
		Beggining: sql.NullTime{Time: beggining, Valid: true},
		Deadline:  sql.NullTime{Time: deadline, Valid: true},
		Color:     updateTaskData.Color,
		TaskID:    taskUUID,
	})
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "an error occurred while processing your request",
			"traceid": traceid,
		})
	}
	responseData := TaskDataResponse{
		TaskID:    updatedTask.TaskID.String(),
		ProjectID: updatedTask.ProjectID.String(),
		Title:     updatedTask.Title,
		Info:      updatedTask.Info.String,
		Tag:       updatedTask.Tag.String,
		CreatedBy: updatedTask.CreatedBy.String(),
		CreatedAt: updatedTask.CreatedAt.Time.Format(time.RFC3339),
		Beggining: updatedTask.Beggining.Time.Format(time.RFC3339),
		Deadline:  updatedTask.Deadline.Time.Format(time.RFC3339),
		Color:     updatedTask.Color.String,
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "task updated successfully",
		"task":    responseData,
	})
}

func (s Server) DeleteTask(c echo.Context) error {
	return nil
}

func (s Server) GetUsersInTask(c echo.Context) error {
	return nil
}

func (s Server) AddUserToTask(c echo.Context) error {
	return nil
}

func (s Server) DeleteUserFromTask(c echo.Context) error {
	return nil
}
