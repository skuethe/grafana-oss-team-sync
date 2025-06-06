package config

import (
	"errors"
	"fmt"
)

type GrafanaPermissions int

const (
	PermissionViewer GrafanaPermissions = 1 << iota
	PermissionEditor
	PermissionAdmin
)

func ValidateGrafanaPermission(in GrafanaPermissions) error {
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
