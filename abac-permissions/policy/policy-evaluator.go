package policy

type PolicyEvaluator interface {
	Evaluate(subject Subject, action Action, resource Resource) bool
}

type ABACPolicy struct{}

func (p ABACPolicy) Evaluate(subject Subject, action Action, resource Resource) bool {
	// Admin has full access
	if contains(subject.Roles, Admin) {
		return true
	}

	// Moderator permissions
	if contains(subject.Roles, Moderator) {
		switch action {
		case View, Update, Create:
			return true
		case Delete:
			if todo, ok := resource.(Todo); ok {
				return todo.IsDone
			}

			if _, ok := resource.(Comment); ok {
				return true
			}

			return false
		}
	}

	// User permissions
	if contains(subject.Roles, User) {
		if isBlocked(subject, resource) {
			return false
		}

		switch action {
		case View, Create:
			return true
		case Update, Delete:
			return resource.GetAuthorID() == subject.ID
		}
	}

	return false
}

// contains checks if a given role is present in the list of Roles.
func contains(roles []Role, role Role) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

// isBlocked checks if the subject is blocked by the author of the resource.
func isBlocked(subject Subject, resource Resource) bool {
	for _, blockedByID := range subject.BlockedBy {
		if blockedByID == resource.GetAuthorID() {
			return true
		}
	}
	return false
}
