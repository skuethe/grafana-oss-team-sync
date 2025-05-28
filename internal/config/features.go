package config

type FeatureSchema struct {
	DisableFolders  bool `koanf:"disableFolders"`
	DisableUserSync bool `koanf:"disableUserSync"`
}

var Feature *FeatureSchema

func valdidateFeatureSchema() error {
	return K.Unmarshal("features", &Feature)
}
