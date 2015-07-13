package main

import "os"

type Config struct {
	Env           map[string]string              `yaml:"env"`
	ProjectName   string                         `yaml:"project_name"`
	Notifications map[string](map[string]string) `yaml:"notifications"`
	BeforeInstall []string                       `yaml:"before_install"`
	Install       []string                       `yaml:"install"`
	BeforeScript  []string                       `yaml:"before_script"`
	Script        []string                       `yaml:"script"`
	AfterScript   []string                       `yaml:"afterScript"`
}

func (conf *Config) applyEnvironmentVariables() {
	env := conf.Env

	if env != nil {
		for key, value := range env {
			os.Setenv(key, value)
		}
	}
}
