package api

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/Vodnik-Project/vodnik-api/db/sqlc"
	log "github.com/Vodnik-Project/vodnik-api/logger"
	"github.com/Vodnik-Project/vodnik-api/types"
	"github.com/Vodnik-Project/vodnik-api/util"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx"
	"github.com/labstack/echo/v4"
	"gopkg.in/validator.v2"
)

func (s Server) CreateTask(c echo.Context) error {
	var task types.CreateTaskParams
	err := c.Bind(&task)
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("can't parse input data")
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"message": "invalid input",
			"traceid": traceid,
		})
	}
	if err = validator.Validate(task); err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("invalid input data")
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"message": "invalid input data",
			"error":   err.Error(),
			"traceid": traceid,
		})
	}
	userUUID := c.Get("userUUID").(uuid.UUID)
	beggining, _ := time.Parse(time.RFC3339, task.Beggining)
	deadline, _ := time.Parse(time.RFC3339, task.Deadline)
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
	responseData := types.TaskData{
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
	taskUUID := c.Get("taskUUID").(uuid.UUID)
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
	responseData := types.TaskData{
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
	ctx := c.Request().Context()
	var taskRequest types.GetTaskParams
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
		Column2:        taskRequest.Limit,
		Offset:         taskRequest.Page * taskRequest.Page,
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
	var responseData []types.TaskData
	for _, j := range projects {
		responseData = append(responseData, types.TaskData{
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
	var updateTaskData types.CreateTaskParams
	err := c.Bind(&updateTaskData)
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("can't parse input data")
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"message": "invalid input",
			"traceid": traceid,
		})
	}
	if updateTaskData == (types.CreateTaskParams{}) {
		traceid := util.RandomString(8)
		log.Logger.Err(errors.New("input data is empty")).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"message": "input data is empty",
			"traceid": traceid,
		})
	}
	beggining, _ := time.Parse(time.RFC3339, updateTaskData.Beggining)
	deadline, _ := time.Parse(time.RFC3339, updateTaskData.Deadline)
	taskUUID := c.Get("taskUUID").(uuid.UUID)
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
	responseData := types.TaskData{
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
	taskUUID := c.Get("taskUUID").(uuid.UUID)
	err := s.store.DeleteTask(ctx, taskUUID)
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
	ctx := c.Request().Context()
	taskUUID := c.Get("taskUUID").(uuid.UUID)
	users, err := s.store.GetUsersByTaskID(ctx, taskUUID)
	if err != nil {
		traceid := util.RandomString(8)
		if err == sql.ErrNoRows {
			log.Logger.Err(err).Str("traceid", traceid).Msg("")
			return c.JSON(http.StatusNotFound, echo.Map{
				"message": "no users found",
				"traceid": traceid,
			})
		}
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "an error occurred while processing your request",
			"traceid": traceid,
		})
	}
	var responseData []types.UsersInTaskData
	for _, i := range users {
		responseData = append(responseData, types.UsersInTaskData{
			UserID:   i.UserID.String(),
			Username: i.Username,
			Bio:      i.Bio.String,
			AddedAt:  i.AddedAt.Time.Format(time.RFC3339),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "users found",
		"users":   responseData,
	})
}

func (s Server) AddUserToTask(c echo.Context) error {
	ctx := c.Request().Context()
	userToAdd := c.Param("userid")
	userToAddUUID, err := uuid.FromString(userToAdd)
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"message": "invalid userid",
			"traceid": traceid,
		})
	}
	taskUUID := c.Get("taskUUID").(uuid.UUID)
	_, err = s.store.AddUserToTask(ctx, sqlc.AddUserToTaskParams{
		UserID: userToAddUUID,
		TaskID: taskUUID,
	})
	if err != nil {
		traceid := util.RandomString(8)
		if err.(pgx.PgError).Code == "23505" {
			log.Logger.Err(err).Str("traceid", traceid).Msg("")
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "user already exist in task",
				"traceid": traceid,
			})
		}
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "an error occurred while processing your request",
			"traceid": traceid,
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "user added to task successfully",
	})
}

func (s Server) DeleteUserFromTask(c echo.Context) error {
	ctx := c.Request().Context()
	userToAdd := c.Param("userid")
	userToDeleteUUID, err := uuid.FromString(userToAdd)
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusUnprocessableEntity, echo.Map{
			"message": "invalid userid",
			"traceid": traceid,
		})
	}
	taskUUID := c.Get("taskUUID").(uuid.UUID)
	err = s.store.DeleteUserFromTask(ctx, sqlc.DeleteUserFromTaskParams{
		UserID: userToDeleteUUID,
		TaskID: taskUUID,
	})
	if err != nil {
		traceid := util.RandomString(8)
		log.Logger.Err(err).Str("traceid", traceid).Msg("")
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "an error occurred while processing your request",
			"traceid": traceid,
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"message": "user deleted from task successfully",
	})
}
