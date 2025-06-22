// SPDX-FileCopyrightText: 2025 Sebastian Küthe and (other) contributors to project grafana-oss-team-sync <https://github.com/skuethe/grafana-oss-team-sync>
// SPDX-License-Identifier: GPL-3.0-or-later

package configtypes

type Folders map[string]Folder

type Folder struct {
	Title       string            `yaml:"title"`
	Description string            `yaml:"description"`
	Permissions FolderPermissions `yaml:"permissions"`
}

type FolderPermissions struct {
	Teams map[string]GrafanaPermission `yaml:"teams"`
}
