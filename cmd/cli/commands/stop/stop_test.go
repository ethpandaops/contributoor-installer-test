package stop

import (
	"errors"
	"flag"
	"testing"

	"github.com/ethpandaops/contributoor-installer/cmd/cli/options"
	"github.com/ethpandaops/contributoor-installer/internal/sidecar"
	"github.com/ethpandaops/contributoor-installer/internal/sidecar/mock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli"
	"go.uber.org/mock/gomock"
)

func TestStopContributoor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name          string
		runMethod     string
		setupMocks    func(*mock.MockConfigManager, *mock.MockDockerSidecar, *mock.MockBinarySidecar)
		expectedError string
	}{
		{
			name:      "docker - stops running service successfully",
			runMethod: sidecar.RunMethodDocker,
			setupMocks: func(cfg *mock.MockConfigManager, d *mock.MockDockerSidecar, b *mock.MockBinarySidecar) {
				cfg.EXPECT().Get().Return(&sidecar.Config{
					RunMethod: sidecar.RunMethodDocker,
					Version:   "latest",
				}).Times(1)
				d.EXPECT().IsRunning().Return(true, nil)
				d.EXPECT().Stop().Return(nil)
			},
		},
		{
			name:      "docker - service not running",
			runMethod: sidecar.RunMethodDocker,
			setupMocks: func(cfg *mock.MockConfigManager, d *mock.MockDockerSidecar, b *mock.MockBinarySidecar) {
				cfg.EXPECT().Get().Return(&sidecar.Config{
					RunMethod: sidecar.RunMethodDocker,
				}).Times(1)
				d.EXPECT().IsRunning().Return(false, nil)
			},
		},
		{
			name:      "docker - stop fails",
			runMethod: sidecar.RunMethodDocker,
			setupMocks: func(cfg *mock.MockConfigManager, d *mock.MockDockerSidecar, b *mock.MockBinarySidecar) {
				cfg.EXPECT().Get().Return(&sidecar.Config{
					RunMethod: sidecar.RunMethodDocker,
				}).Times(1)
				d.EXPECT().IsRunning().Return(true, nil)
				d.EXPECT().Stop().Return(errors.New("stop failed"))
			},
			expectedError: "stop failed",
		},
		{
			name:      "binary - stops running service successfully",
			runMethod: sidecar.RunMethodBinary,
			setupMocks: func(cfg *mock.MockConfigManager, d *mock.MockDockerSidecar, b *mock.MockBinarySidecar) {
				cfg.EXPECT().Get().Return(&sidecar.Config{
					RunMethod: sidecar.RunMethodBinary,
				}).Times(1)
				b.EXPECT().IsRunning().Return(true, nil)
				b.EXPECT().Stop().Return(nil)
			},
		},
		{
			name:      "binary - service not running",
			runMethod: sidecar.RunMethodBinary,
			setupMocks: func(cfg *mock.MockConfigManager, d *mock.MockDockerSidecar, b *mock.MockBinarySidecar) {
				cfg.EXPECT().Get().Return(&sidecar.Config{
					RunMethod: sidecar.RunMethodBinary,
				}).Times(1)
				b.EXPECT().IsRunning().Return(false, nil)
			},
		},
		{
			name:      "invalid sidecar run method",
			runMethod: "invalid",
			setupMocks: func(cfg *mock.MockConfigManager, d *mock.MockDockerSidecar, b *mock.MockBinarySidecar) {
				cfg.EXPECT().Get().Return(&sidecar.Config{
					RunMethod: "invalid",
				}).Times(1)
			},
			expectedError: "invalid sidecar run method",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockConfig := mock.NewMockConfigManager(ctrl)
			mockDocker := mock.NewMockDockerSidecar(ctrl)
			mockBinary := mock.NewMockBinarySidecar(ctrl)

			tt.setupMocks(mockConfig, mockDocker, mockBinary)

			app := cli.NewApp()
			ctx := cli.NewContext(app, nil, nil)

			err := stopContributoor(ctx, logrus.New(), mockConfig, mockDocker, mockBinary)

			if tt.expectedError != "" {
				assert.ErrorContains(t, err, tt.expectedError)

				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestRegisterCommands(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name          string
		configPath    string
		expectedError string
	}{
		{
			name:       "successfully registers command",
			configPath: "testdata/valid", // "testdata" is an ancillary dir provided by go-test.
		},
		{
			name:          "fails when config service fails",
			configPath:    "/invalid/path/that/doesnt/exist",
			expectedError: "error loading config",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create CLI app, with the config flag.
			app := cli.NewApp()
			app.Flags = []cli.Flag{
				cli.StringFlag{
					Name: "config-path",
				},
			}

			// Ensure we set the config path flag.
			globalSet := flag.NewFlagSet("test", flag.ContinueOnError)
			globalSet.String("config-path", "", "")
			err := globalSet.Set("config-path", tt.configPath)
			require.NoError(t, err)

			// Create the cmd context.
			globalCtx := cli.NewContext(app, globalSet, nil)
			app.Metadata = map[string]interface{}{
				"flagContext": globalCtx,
			}

			// Now test!
			RegisterCommands(
				app,
				options.NewCommandOpts(
					options.WithName("stop"),
					options.WithLogger(logrus.New()),
					options.WithAliases([]string{"s"}),
				),
			)

			if tt.expectedError != "" {
				// Ensure the command registration succeeded.
				assert.NoError(t, err)

				// Assert that the action execution fails as expected.
				cmd := app.Commands[0]
				ctx := cli.NewContext(app, nil, globalCtx)

				// Assert that the action is the func we expect, mainly because the linter is having a fit otherwise.
				action, ok := cmd.Action.(func(*cli.Context) error)
				require.True(t, ok, "expected action to be func(*cli.Context) error")

				// Execute the action and assert the error.
				actionErr := action(ctx)
				assert.Error(t, actionErr)
				assert.ErrorContains(t, actionErr, tt.expectedError)
			} else {
				// Ensure the command registration succeeded.
				assert.NoError(t, err)
				assert.Len(t, app.Commands, 1)

				// Ensure the command is registered as expected by dumping the command.
				cmd := app.Commands[0]
				assert.Equal(t, "stop", cmd.Name)
				assert.Equal(t, []string{"s"}, cmd.Aliases)
				assert.Equal(t, "Stop Contributoor", cmd.Usage)
				assert.Equal(t, "contributoor stop [options]", cmd.UsageText)
				assert.NotNil(t, cmd.Action)
			}
		})
	}
}
