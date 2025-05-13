package config

type FolderSchema struct {
	Title       string `koanf:"title"`
	Description string `koanf:"description"`
	Permissions struct {
		Teams map[string]folderPermissions `koanf:"teams"`
	} `koanf:"permissions"`
}

type folderPermissions int64

const (
	PermissionViewer folderPermissions = 1 << iota
	PermissionEditor
	PermissionAdmin
)
