package config_loader

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"

	"github.com/ne006/shepherd/supervisor"
)

type ConfigLoaderTestSuite struct {
	suite.Suite
	ConfigPath string
	ExampleApp supervisor.App
	Logger     *zap.SugaredLogger
}

func (suite *ConfigLoaderTestSuite) SetupTest() {
	if logger, err := zap.NewProduction(); err != nil {
		suite.T().Errorf("%v\n", err)
	} else {
		suite.Logger = logger.Sugar()
	}

	suite.ConfigPath = ("./example_config.yml")

	app := supervisor.App{
		Name:   "test_app",
		Config: "./example_config.yml",
	}
	app.SetEnv(map[string]any{
		"APP": "test_app",
	})

	webGroup := supervisor.Group{
		Name: "web",
	}
	webGroup.SetEnv(map[string]any{
		"APP": "test_app",
	})

	bgGroup := supervisor.Group{
		Name:     "bg_jobs",
		Children: []supervisor.SupervisionUnit{},
	}
	bgGroup.SetEnv(map[string]any{
		"APP":     "test_app",
		"LOG_DIR": "/home/user/log/test_app/bg_jobs",
	})

	cronGroup := supervisor.Group{
		Name:     "cron",
		Children: []supervisor.SupervisionUnit{},
	}
	cronGroup.SetEnv(map[string]any{
		"APP":     "test_app",
		"LOG_DIR": "/home/user/log/test_app/bg_jobs",
	})

	webServer := supervisor.Process{
		Name: "web_server",

		Cmd: exec.Cmd{
			Path: "web_srv",
			Args: []string{"-p", "8080", "-b", "0.0.0.0", "-e", "development"},
		},

		Logger: suite.Logger,
	}
	webServer.SetEnv(map[string]any{
		"APP": "test_app",
	})

	bgServer := supervisor.Process{
		Name: "bg_server",
		Cmd: exec.Cmd{
			Path: "bg_srv",
			Args: []string{"-e", "development"},
		},

		Logger: suite.Logger,
	}
	bgServer.SetEnv(map[string]any{
		"APP":     "test_app",
		"LOG_DIR": "/home/user/log/test_app/bg_jobs",
		"API_KEY": "{{env \"API_KEY\"}}",
		"WORKERS": 5,
	})

	cron := supervisor.Process{
		Name: "cron",

		Cmd: exec.Cmd{
			Path: "cron",
			Args: []string{},
		},

		Logger: suite.Logger,
	}
	cron.SetEnv(map[string]any{
		"LOG_DIR": "/home/user/log/test_app/bg_jobs",
		"APP":     "test_app",
	})

	cronGroup.Children = []supervisor.SupervisionUnit{&cron}

	webGroup.Children = []supervisor.SupervisionUnit{&webServer}
	bgGroup.Children = []supervisor.SupervisionUnit{&bgServer, &cronGroup}

	app.Children = []supervisor.SupervisionUnit{&webGroup, &bgGroup}

	suite.ExampleApp = app
}

// All methods that begin with "Test" are run as tests within a
// suite.
func (suite *ConfigLoaderTestSuite) TestLoadConfig() {
	before := os.Getenv("API_KEY")
	os.Setenv("API_KEY", "token")

	if app, err := LoadConfig("./example_config.yml", suite.Logger); err != nil {
		os.Setenv("API_KEY", before)
		suite.T().Errorf("%v\n", err)
	} else {
		os.Setenv("API_KEY", before)
		assert.Equal(suite.T(), &suite.ExampleApp, app)
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigLoaderTestSuite))
}
