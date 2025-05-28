package config

type TeamsSchema struct {
	List []string `koanf:"teams"`
}

var Teams *TeamsSchema = &TeamsSchema{}

func validateTeamsSchema() error {
	return K.Unmarshal("", &Teams)
}
