package profiles

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Profiles struct {
	Profiles    map[string]Profile
	persistPath string
}

func (p *Profiles) AddProfile(profile Profile) {
	p.Profiles[profile.Name] = profile
}

func (p *Profiles) RemoveProfile(name string) {
	delete(p.Profiles, name)
}

func (p *Profiles) GetProfile(name string) Profile {
	return p.Profiles[name]
}

func LoadFromFile(path string) (*Profiles, error) {
	privateViper := viper.New()
	privateViper.SetConfigType("yaml")
	privateViper.SetConfigFile(path)

	if err := privateViper.ReadInConfig(); err != nil {
		if _, ok := err.(*os.PathError); ok {
			fmt.Printf("No profile file found at %s, creating a new one\n", path)
			return &Profiles{Profiles: map[string]Profile{}, persistPath: path}, nil
		}
		return &Profiles{Profiles: map[string]Profile{}, persistPath: path}, err
	}

	profiles := make(map[string]Profile, 0)
	if err := privateViper.UnmarshalKey("profiles", &profiles); err != nil {
		return &Profiles{Profiles: map[string]Profile{}, persistPath: path}, err
	}

	return &Profiles{Profiles: profiles, persistPath: path}, nil
}

func (p *Profiles) Persist() error {
	privateViper := viper.New()
	privateViper.SetConfigType("yaml")
	privateViper.Set("profiles", p.Profiles)
	return privateViper.WriteConfigAs(p.persistPath)
}
