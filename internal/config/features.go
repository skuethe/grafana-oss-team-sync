package config

type FeatureSchema struct {
	DisableFolders       bool `koanf:"disableFolders"`
	DisableUserSync      bool `koanf:"disableUserSync"`
	AddLocalAdminToTeams bool `koanf:"addLocalAdminToTeams"`
}

var Feature *FeatureSchema

func valdidateFeatureSchema() error {
	return K.Unmarshal(ConfigParamFeatures, &Feature)
}
