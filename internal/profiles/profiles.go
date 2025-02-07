package profiles

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Profiles struct {
	Profiles    map[string]*Profile
	persistPath string
}

func (p *Profiles) AddProfile(profile *Profile) {
	p.Profiles[profile.Name] = profile
}

func (p *Profiles) RemoveProfile(name string) {
	delete(p.Profiles, name)
}

func (p *Profiles) LoadProfile(name string) *Profile {
	if val, ok := p.Profiles[name]; ok {
		return val
	}

	return nil
}

func LoadFromFile(path string) (*Profiles, error) {
	privateViper := viper.New()
	privateViper.SetConfigType("yaml")
	privateViper.SetConfigFile(path)

	if err := privateViper.ReadInConfig(); err != nil {
		if _, ok := err.(*os.PathError); ok {
			fmt.Printf("No profile file found at %s, creating a new one\n", path)
			return &Profiles{Profiles: map[string]*Profile{}, persistPath: path}, nil
		}
		return &Profiles{Profiles: map[string]*Profile{}, persistPath: path}, err
	}

	savedProfiles := privateViper.Get("profiles").(map[string]interface{})
	profiles := make(map[string]*Profile, len(savedProfiles))
	for k, v := range savedProfiles {
		p := v.(map[string]interface{})

		profile := &Profile{
			Name:              p["name"].(string),
			TestModeSecretKey: p["test_mode_secret_key"].(string),
			BaseURL:           p["base_url"].(string),
			GrpcServerAddress: p["grpc_server_address"].(string),
		}

		profiles[k] = profile
	}

	return &Profiles{Profiles: profiles, persistPath: path}, nil
}

func (p *Profiles) Persist() error {
	privateViper := viper.New()
	privateViper.SetConfigType("yaml")
	privateViper.Set("profiles", p.Profiles)
	return privateViper.WriteConfigAs(p.persistPath)
}
