package profiles

type Profile struct {
	Name              string `yaml:"name"`
	TestModeSecretKey string `yaml:"test_mode_secret_key"`
}
