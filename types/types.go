package types

// Auth
type LoginParams struct {
	Email    string `json:"email" validate:"nonzero,regexp=^[a-zA-Z0-9]+(?:.[a-zA-Z0-9]+)*@[a-zA-Z0-9]+(?:.[a-zA-Z0-9]+)*$"`
	Password string `json:"password" validate:"nonzero"`
}

type RefreshTokenParams struct {
	RefreshToken string `json:"refresh_token"`
	UserID       string `'json:"userid"`
}

// User
type UserData struct {
	UserID   string `json:"userid"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
	JoinedAt string `json:"joinedAt"`
}

type CreateUserParams struct {
	Username string `json:"username" validate:"nonzero"`
	Email    string `json:"email" validate:"nonzero,regexp=^[a-zA-Z0-9]+(?:.[a-zA-Z0-9]+)*@[a-zA-Z0-9]+(?:.[a-zA-Z0-9]+)*$"`
	Password string `json:"password" validate:"nonzero"`
	Bio      string `json:"bio"`
}

type UpdateUserParams struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

// Project
type ProjectData struct {
	ProjectID string `json:"project_id"`
	Title     string `json:"title"`
	Info      string `json:"info"`
	OwnerID   string `json:"owner_id"`
	CreatedAt string `json:"created_at"`
}

type CreateProjectParams struct {
	Title string `json:"title" validate:"nonzero"`
	Info  string `json:"info"`
}

type UpdateProjectParams struct {
	Title   string `json:"title"`
	Info    string `json:"info"`
	OwnerID string `json:"owner_id"`
}

type UsersInProjectData struct {
	UserID         string `json:"userid"`
	Username       string `json:"username"`
	Bio            string `json:"bio"`
	AddedToProject string `json:"added_to_project"`
	Admin          bool   `json:"admin"`
}

// Task
type TaskData struct {
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

type CreateTaskParams struct {
	Title     string `json:"title" validate:"nonzero"`
	Info      string `json:"info"`
	Tag       string `json:"tag"`
	Beggining string `json:"beggining" validate:"regexp=(^$|^([0-9]+)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])[Tt]([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]|60)(\\.[0-9]+)?(([Zz])|([\\+|\\-]([01][0-9]|2[0-3]):[0-5][0-9]))$)"`
	Deadline  string `json:"deadline" validate:"regexp=(^$|^([0-9]+)-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])[Tt]([01][0-9]|2[0-3]):([0-5][0-9]):([0-5][0-9]|60)(\\.[0-9]+)?(([Zz])|([\\+|\\-]([01][0-9]|2[0-3]):[0-5][0-9]))$)"`
	Color     string `json:"color"`
}

type GetTaskParams struct {
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

type UsersInTaskData struct {
	UserID   string `json:"userid"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
	AddedAt  string `json:"added_at"`
}
