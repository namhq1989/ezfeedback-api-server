package domain

type ProjectRole string

const (
	ProjectRoleUnknown ProjectRole = ""
	ProjectRoleOwner   ProjectRole = "owner"
	ProjectRoleEditor  ProjectRole = "editor"
	ProjectRoleViewer  ProjectRole = "viewer"
)

func (s ProjectRole) String() string {
	return string(s)
}

func (s ProjectRole) IsValid() bool {
	return s != ProjectRoleUnknown
}

func (s ProjectRole) IsOwner() bool {
	return s == ProjectRoleOwner
}

func (s ProjectRole) IsEditor() bool {
	return s == ProjectRoleEditor
}

func (s ProjectRole) IsViewer() bool {
	return s == ProjectRoleViewer
}

func ToProjectRole(s string) ProjectRole {
	switch s {
	case ProjectRoleOwner.String():
		return ProjectRoleOwner
	case ProjectRoleEditor.String():
		return ProjectRoleEditor
	case ProjectRoleViewer.String():
		return ProjectRoleViewer
	default:
		return ProjectRoleUnknown
	}
}
