package profiles

type Profile struct {
	Name              string `yaml:"name"`
	TestModeSecretKey string `yaml:"test_mode_secret_key"`
	BaseURL           string `yaml:"base_url"`
	GrpcServerAddress string `yaml:"grpc_server_address"`
}
