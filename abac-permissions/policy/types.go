package policy

import "time"

type Role string

const (
	Admin     Role = "admin"
	Moderator Role = "moderator"
	User      Role = "user"
)

type Action string

const (
	View   Action = "view"
	Update Action = "update"
	Create Action = "create"
	Delete Action = "delete"
)

type Subject struct {
	ID        string
	Roles     []Role
	BlockedBy []string
}

type Resource interface {
	GetID() string
	GetAuthorID() string
	GetCreatedAt() time.Time
}

type Comment struct {
	ID        string
	Body      string
	AuthorId  string
	CreatedAt time.Time
}

func (c Comment) GetID() string {
	return c.ID
}

func (c Comment) GetAuthorID() string {
	return c.AuthorId
}

func (c Comment) GetCreatedAt() time.Time {
	return c.CreatedAt
}

type Todo struct {
	ID        string
	Body      string
	IsDone    bool
	AuthorId  string
	CreatedAt time.Time
}

func (t Todo) GetID() string {
	return t.ID
}

func (t Todo) GetAuthorID() string {
	return t.AuthorId
}

func (t Todo) GetCreatedAt() time.Time {
	return t.CreatedAt
}
