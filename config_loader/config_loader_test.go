package config_loader

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/ne006/shepherd/supervisor"
)

type ConfigLoaderTestSuite struct {
	suite.Suite
	ConfigPath string
	ExampleApp supervisor.App
}

func (suite *ConfigLoaderTestSuite) SetupTest() {
	suite.ConfigPath = ("./example_config.yml")

	app := supervisor.App{
		Name:   "test_app",
		Config: "./example_config.yml",
	}

	webGroup := supervisor.Group{
		Name: "web",
	}

	bgGroup := supervisor.Group{
		Name:     "bg_jobs",
		Children: []supervisor.SupervisionUnit{},
	}

	webServer := supervisor.Process{
		Name: "web_server",

		Cmd: exec.Cmd{
			Path: "web_srv",
			Args: []string{"-p", "8080", "-b", "0.0.0.0", "-e", "development"},
		},
	}

	bgServer := supervisor.Process{
		Name: "bg_server",
		Cmd: exec.Cmd{
			Path: "bg_srv",
			Args: []string{"-e", "development"},
		},
	}

	webGroup.Children = []supervisor.SupervisionUnit{webServer}
	bgGroup.Children = []supervisor.SupervisionUnit{bgServer}

	app.Children = []supervisor.SupervisionUnit{webGroup, bgGroup}

	suite.ExampleApp = app
}

// All methods that begin with "Test" are run as tests within a
// suite.
func (suite *ConfigLoaderTestSuite) TestLoadConfig() {
	if app, err := LoadConfig("./example_config.yml"); err != nil {
		suite.T().Errorf("%v\n", err)
	} else {
		assert.Equal(suite.T(), &suite.ExampleApp, app)
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigLoaderTestSuite))
}
