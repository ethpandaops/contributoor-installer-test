package update

import (
	"fmt"

	"github.com/ethpandaops/contributoor-installer/cmd/cli/options"
	"github.com/ethpandaops/contributoor-installer/internal/service"
	"github.com/ethpandaops/contributoor-installer/internal/tui"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func RegisterCommands(app *cli.App, opts *options.CommandOpts) error {
	app.Commands = append(app.Commands, cli.Command{
		Name:      opts.Name(),
		Aliases:   opts.Aliases(),
		Usage:     "Update Contributoor to the latest version",
		UsageText: "contributoor update [options]",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "version, v",
				Usage: "The contributoor version to update to",
				Value: "latest",
			},
		},
		Action: func(c *cli.Context) error {
			log := opts.Logger()

			configService, err := service.NewConfigService(log, c.GlobalString("config-path"))
			if err != nil {
				return fmt.Errorf("error loading config: %w", err)
			}

			dockerService, err := service.NewDockerService(log, configService)
			if err != nil {
				return fmt.Errorf("error creating docker service: %w", err)
			}

			binaryService := service.NewBinaryService(log, configService)
			githubService := service.NewGitHubService("ethpandaops", "contributoor")

			return updateContributoor(c, log, configService, dockerService, binaryService, githubService)
		},
	})

	return nil
}

func updateContributoor(
	c *cli.Context,
	log *logrus.Logger,
	config service.ConfigManager,
	docker service.DockerService,
	binary service.BinaryService,
	github service.GitHubService,
) error {
	var (
		success        bool
		targetVersion  string
		cfg            = config.Get()
		currentVersion = cfg.Version
	)

	log.WithField("version", currentVersion).Info("Current version")

	defer func() {
		if !success {
			if err := rollbackVersion(log, config, currentVersion); err != nil {
				log.Error(err)
			}
		}
	}()

	// Determine target version.
	targetVersion, err := determineTargetVersion(c, github)
	if err != nil {
		// Flag as success, there's nothing to update on rollback if we fail to determine the target version.
		success = true

		return err
	}

	// Check if update is needed.
	if targetVersion == currentVersion {
		// Flag as success, there's nothing to update.
		success = true

		logUpdateStatus(log, c.IsSet("version"), targetVersion)

		return nil
	}

	// Update config version.
	if uerr := updateConfigVersion(config, targetVersion); uerr != nil {
		return uerr
	}

	// Refresh our config state, given it was updated above.
	cfg = config.Get()

	// Update the service.
	log.WithField("version", cfg.Version).Info("Updating Contributoor")

	success, err = updateService(log, cfg, docker, binary)
	if err != nil {
		return err
	}

	log.Infof(
		"%sContributoor updated successfully to version %s%s",
		tui.TerminalColorGreen,
		cfg.Version,
		tui.TerminalColorReset,
	)

	return nil
}

func updateService(log *logrus.Logger, cfg *service.ContributoorConfig, docker service.DockerService, binary service.BinaryService) (bool, error) {
	switch cfg.RunMethod {
	case service.RunMethodDocker:
		return updateDocker(log, cfg, docker)
	case service.RunMethodBinary:
		return updateBinary(log, cfg, binary)
	default:
		return false, fmt.Errorf("invalid run method: %s", cfg.RunMethod)
	}
}

func updateBinary(log *logrus.Logger, cfg *service.ContributoorConfig, binary service.BinaryService) (bool, error) {
	// Check if service is currently running.
	running, err := binary.IsRunning()
	if err != nil {
		log.Errorf("could not check service status: %v", err)

		return false, err
	}

	// If the service is running, we need to stop it before we can update the binary.
	if running {
		if tui.Confirm("Service is running. In order to update, it must be stopped. Would you like to stop it?") {
			if err := binary.Stop(); err != nil {
				return false, fmt.Errorf("failed to stop service: %w", err)
			}
		} else {
			log.Error("update process was cancelled")

			return false, nil
		}
	}

	if err := binary.Update(); err != nil {
		log.Errorf("could not update service: %v", err)

		return false, err
	}

	// if err := config.Save(); err != nil {
	// 	log.Errorf("could not save updated config: %v", err)
	// 	return false, err
	// }

	// If it was running, start it again for them.
	if running {
		if err := binary.Start(); err != nil {
			return true, fmt.Errorf("failed to start service: %w", err)
		}
	}

	return true, nil
}

func updateDocker(log *logrus.Logger, cfg *service.ContributoorConfig, docker service.DockerService) (bool, error) {
	if err := docker.Update(); err != nil {
		log.Errorf("could not update service: %v", err)

		return false, err
	}

	// Check if service is currently running.
	running, err := docker.IsRunning()
	if err != nil {
		log.Errorf("could not check service status: %v", err)

		return true, err
	}

	// If the service is running, we need to restart it with the new version.
	if running {
		if tui.Confirm("Service is running. Would you like to restart it with the new version?") {
			if err := docker.Stop(); err != nil {
				return true, fmt.Errorf("failed to stop service: %w", err)
			}

			if err := docker.Start(); err != nil {
				return true, fmt.Errorf("failed to start service: %w", err)
			}
		} else {
			log.Info("service will continue running with the previous version until next restart")
		}
	} else {
		if tui.Confirm("Service is not running. Would you like to start it?") {
			if err := docker.Start(); err != nil {
				return true, fmt.Errorf("failed to start service: %w", err)
			}
		}
	}

	return true, nil
}

func determineTargetVersion(c *cli.Context, github service.GitHubService) (string, error) {
	if c.IsSet("version") {
		version := c.String("version")

		exists, err := github.VersionExists(version)
		if err != nil {
			return "", fmt.Errorf("failed to check version: %w", err)
		}

		if !exists {
			return "", fmt.Errorf(
				"%sversion %s not found. Use 'contributoor update' without --version to get the latest version%s",
				tui.TerminalColorRed,
				version,
				tui.TerminalColorReset,
			)
		}

		return version, nil
	}

	version, err := github.GetLatestVersion()
	if err != nil {
		return "", fmt.Errorf("failed to get latest version: %w", err)
	}

	return version, nil
}

func updateConfigVersion(config service.ConfigManager, version string) error {
	if err := config.Update(func(cfg *service.ContributoorConfig) {
		cfg.Version = version
	}); err != nil {
		return fmt.Errorf("failed to update config version: %w", err)
	}

	if err := config.Save(); err != nil {
		return fmt.Errorf("could not save updated config: %w", err)
	}

	return nil
}

func rollbackVersion(log *logrus.Logger, config service.ConfigManager, version string) error {
	if err := config.Update(func(cfg *service.ContributoorConfig) {
		cfg.Version = version
	}); err != nil {
		return fmt.Errorf("failed to roll back version in config: %w", err)
	}

	if err := config.Save(); err != nil {
		return fmt.Errorf("failed to save config after version rollback: %w", err)
	}

	return nil
}

func logUpdateStatus(log *logrus.Logger, isVersionSet bool, version string) {
	if isVersionSet {
		log.Infof(
			"%scontributoor is already running version %s%s",
			tui.TerminalColorGreen,
			version,
			tui.TerminalColorReset,
		)
	} else {
		log.Infof(
			"%scontributoor is up to date at version %s%s",
			tui.TerminalColorGreen,
			version,
			tui.TerminalColorReset,
		)
	}
}
