package cmd

import "github.com/stretchr/testify/suite"

func newTestConfig(s *suite.Suite) *Config {
	config, err := NewConfig()
	s.NoError(err)

	return config
}
