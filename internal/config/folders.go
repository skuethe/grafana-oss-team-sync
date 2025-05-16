package config

import (
	"errors"
	"fmt"
)

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

func ValidateFolderPermission(in folderPermissions) error {
	switch in {
	case PermissionViewer:
		return nil
	case PermissionEditor:
		return nil
	case PermissionAdmin:
		return nil
	}

	return errors.New("invalid permission defined: " + fmt.Sprint(in))
}
