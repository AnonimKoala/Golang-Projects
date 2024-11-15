package main

import "time"

type Comment struct {
	id        string
	body      string
	authorId  Account
	createdAt time.Time
}

type Role string

const (
	Admin     Role = "admin"
	Moderator Role = "moderator"
	User      Role = "user"
)

type Account struct {
	id        string
	roles     []Role
	blockedBy []Account
}

type RolePermissions struct {
	Comments CommentPermissions
}

type CommentPermissions struct {
	View   bool
	Update bool
	Delete bool
}
