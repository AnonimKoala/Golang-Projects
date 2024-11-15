package main

var Roles = map[Role]RolePermissions{
	Admin: {
		Comments: CommentPermissions{View: true, Update: true, Delete: true},
	},
	Moderator: {
		Comments: CommentPermissions{View: true, Update: true, Delete: false},
	},
	User: {
		Comments: CommentPermissions{View: true, Update: false, Delete: false},
	},
}
