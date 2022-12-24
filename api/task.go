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
	"gopkg.in/validator.v2"
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

type getTaskRequest struct {
	Title          string `json:"title"`
	Info           string `json:"info"`
	Tag            string `json:"tag"`
	CreatedBy      string `json:"created_by"`
	CreatedAtFrom  string `json:"created_at_from" validate:"regexp=(^$|^([0-9]+)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])[Tt]([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]|60)(\\.[0-9]+)?(([Zz])|([\\+|\\-]([01][0-9]|2[0-3]):[0-5][0-9]))$)"`
	CreatedAtUntil string `json:"created_at_until" validate:"regexp=(^$|^([0-9]+)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])[Tt]([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]|60)(\\.[0-9]+)?(([Zz])|([\\+|\\-]([01][0-9]|2[0-3]):[0-5][0-9]))$)"`
	BegginingFrom  string `json:"beggining_from" validate:"regexp=(^$|^([0-9]+)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])[Tt]([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]|60)(\\.[0-9]+)?(([Zz])|([\\+|\\-]([01][0-9]|2[0-3]):[0-5][0-9]))$)"`
	BegginingUntil string `json:"beggining_until" validate:"regexp=(^$|^([0-9]+)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])[Tt]([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]|60)(\\.[0-9]+)?(([Zz])|([\\+|\\-]([01][0-9]|2[0-3]):[0-5][0-9]))$)"`
	DeadlineFrom   string `json:"deadline_from" validate:"regexp=(^$|^([0-9]+)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])[Tt]([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]|60)(\\.[0-9]+)?(([Zz])|([\\+|\\-]([01][0-9]|2[0-3]):[0-5][0-9]))$)"`
	DeadlineUntil  string `json:"deadline_until" validate:"regexp=(^$|^([0-9]+)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])[Tt]([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]|60)(\\.[0-9]+)?(([Zz])|([\\+|\\-]([01][0-9]|2[0-3]):[0-5][0-9]))$)"`
	Sortdirection  string `json:"sort_direction" validate:"regexp=(^$|asc|desc)"`
	SortBy         string `json:"sort_by" validate:"regexp=(^$|beggining|deadline|created_at)"`
	Limit          int32  `json:"limit"`
	Page           int32  `json:"page"`
}

func (s Server) GetTasksByProjectID(c echo.Context) error {
	ctx := c.Request().Context()
	var taskRequest getTaskRequest
	err := c.Bind(&taskRequest)
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("can't parse input data")
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"message": "invalid input",
			"traceid": traceid,
		})
	}
	var createdBy uuid.UUID
	if taskRequest.CreatedBy != "" {
		createdBy, err = uuid.FromString(taskRequest.CreatedBy)
		if err != nil {
			traceid := util.RandomString(8)
			log.Logger.Err(err).Str("traceid", traceid).Msg("invalid taskid")
			return c.JSON(http.StatusUnprocessableEntity, echo.Map{
				"message": "invalid created_by id",
				"traceid": traceid,
			})
		}
	}
	if err = validator.Validate(taskRequest); err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("invalid input data")
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"message": "invalid input data",
			"error":   err.Error(),
			"traceid": traceid,
		})
	}
	createdAtFrom, _ := time.Parse(time.RFC3339, taskRequest.CreatedAtFrom)
	createdAtUntil, _ := time.Parse(time.RFC3339, taskRequest.CreatedAtUntil)
	begginingFrom, _ := time.Parse(time.RFC3339, taskRequest.BegginingFrom)
	begginingUntil, _ := time.Parse(time.RFC3339, taskRequest.BegginingUntil)
	deadlineFrom, _ := time.Parse(time.RFC3339, taskRequest.DeadlineFrom)
	deadlineUntil, _ := time.Parse(time.RFC3339, taskRequest.DeadlineUntil)
	projects, err := s.store.GetTasksByProjectID(ctx, sqlc.GetTasksByProjectIDParams{
		ProjectID:      c.Get("projectUUID").(uuid.UUID),
		Title:          taskRequest.Title,
		Info:           taskRequest.Info,
		Tag:            taskRequest.Tag,
		CreatedBy:      createdBy,
		CreatedAtFrom:  createdAtFrom,
		CreatedAtUntil: createdAtUntil,
		BegginingFrom:  begginingFrom,
		BegginingUntil: begginingUntil,
		DeadlineFrom:   deadlineFrom,
		DeadlineUntil:  deadlineUntil,
		Sortdirection:  taskRequest.Sortdirection,
		Sortby:         taskRequest.SortBy,
		Limit:          taskRequest.Limit,
		Offset:         taskRequest.Page * taskRequest.Limit,
	})
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
	var responseData []TaskDataResponse
	for _, j := range projects {
		responseData = append(responseData, TaskDataResponse{
			TaskID:    j.TaskID.String(),
			ProjectID: j.ProjectID.String(),
			Title:     j.Title,
			Info:      j.Info.String,
			Tag:       j.Tag.String,
			CreatedBy: j.CreatedBy.String(),
			CreatedAt: j.CreatedAt.Time.Format(time.RFC3339),
			Beggining: j.Beggining.Time.Format(time.RFC3339),
			Deadline:  j.Deadline.Time.Format(time.RFC3339),
			Color:     j.Color.String,
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "tasks found",
		"count":   len(responseData),
		"tasks":   responseData,
	})
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
	err = s.store.DeleteTask(ctx, taskUUID)
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "an error occurred while processing your request",
			"traceid": traceid,
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "task deleted successfully",
	})
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
