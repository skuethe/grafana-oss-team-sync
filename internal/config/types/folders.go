package types

type Folders map[string]Folder

type Folder struct {
	Title       string            `yaml:"title"`
	Description string            `yaml:"description"`
	Permissions FolderPermissions `yaml:"permissions"`
}

type FolderPermissions struct {
	Teams map[string]GrafanaPermission `yaml:"teams"`
}
