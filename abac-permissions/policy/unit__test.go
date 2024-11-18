package policy

import (
	"testing"
	"time"
)

func TestAdminHasFullAccess(t *testing.T) {
	subject := Subject{ID: "1", Roles: []Role{Admin}}
	resource := Comment{ID: "1", AuthorId: "2", CreatedAt: time.Now()}
	policy := ABACPolicy{}

	if !policy.Evaluate(subject, View, resource) {
		t.Errorf("Admin should have full access to view")
	}
	if !policy.Evaluate(subject, Update, resource) {
		t.Errorf("Admin should have full access to update")
	}
	if !policy.Evaluate(subject, Create, resource) {
		t.Errorf("Admin should have full access to create")
	}
	if !policy.Evaluate(subject, Delete, resource) {
		t.Errorf("Admin should have full access to delete")
	}
}

func TestModeratorPermissions(t *testing.T) {
	subject := Subject{ID: "1", Roles: []Role{Moderator}}
	comment := Comment{ID: "1", AuthorId: "2", CreatedAt: time.Now()}
	todo := Todo{ID: "1", AuthorId: "2", CreatedAt: time.Now(), IsDone: true}
	policy := ABACPolicy{}

	if !policy.Evaluate(subject, View, comment) {
		t.Errorf("Moderator should have access to view comments")
	}
	if !policy.Evaluate(subject, Update, comment) {
		t.Errorf("Moderator should have access to update comments")
	}
	if !policy.Evaluate(subject, Create, comment) {
		t.Errorf("Moderator should have access to create comments")
	}
	if !policy.Evaluate(subject, Delete, comment) {
		t.Errorf("Moderator should have access to delete comments")
	}
	if !policy.Evaluate(subject, Delete, todo) {
		t.Errorf("Moderator should have access to delete completed todos")
	}
	todo.IsDone = false
	if policy.Evaluate(subject, Delete, todo) {
		t.Errorf("Moderator should not have access to delete incomplete todos")
	}
}

func TestUserPermissions(t *testing.T) {
	subject := Subject{ID: "1", Roles: []Role{User}}
	comment := Comment{ID: "1", AuthorId: "1", CreatedAt: time.Now()}
	todo := Todo{ID: "1", AuthorId: "1", CreatedAt: time.Now()}
	policy := ABACPolicy{}

	if !policy.Evaluate(subject, View, comment) {
		t.Errorf("User should have access to view their own comments")
	}
	if !policy.Evaluate(subject, Create, comment) {
		t.Errorf("User should have access to create comments")
	}
	if !policy.Evaluate(subject, Update, comment) {
		t.Errorf("User should have access to update their own comments")
	}
	if !policy.Evaluate(subject, Delete, comment) {
		t.Errorf("User should have access to delete their own comments")
	}
	if !policy.Evaluate(subject, Update, todo) {
		t.Errorf("User should have access to update their own todos")
	}
	if !policy.Evaluate(subject, Delete, todo) {
		t.Errorf("User should have access to delete their own todos")
	}
}

func TestUserBlockedPermissions(t *testing.T) {
	subject := Subject{ID: "1", Roles: []Role{User}, BlockedBy: []string{"2"}}
	comment := Comment{ID: "1", AuthorId: "2", CreatedAt: time.Now()}
	policy := ABACPolicy{}

	if policy.Evaluate(subject, View, comment) {
		t.Errorf("User should not have access to view comments from users who blocked them")
	}
	if policy.Evaluate(subject, Update, comment) {
		t.Errorf("User should not have access to update comments from users who blocked them")
	}
	if policy.Evaluate(subject, Delete, comment) {
		t.Errorf("User should not have access to delete comments from users who blocked them")
	}
}
