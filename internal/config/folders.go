package config

type FolderSchema struct {
	Title       string `koanf:"title"`
	Description string `koanf:"description"`
	Permissions struct {
		Teams map[string]GrafanaPermissions `koanf:"teams"`
	} `koanf:"permissions"`
}
