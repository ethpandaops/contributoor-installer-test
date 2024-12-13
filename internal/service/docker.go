package service

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ethpandaops/contributoor-installer/internal/tui"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -package mock -destination mock/docker.mock.go github.com/ethpandaops/contributoor-installer/internal/service DockerService

type DockerService interface {
	ServiceRunner
}

// DockerService is a basic service for interacting with the docker container.
type dockerService struct {
	logger        *logrus.Logger
	composePath   string
	configPath    string
	configService ConfigManager
}

// NewDockerService creates a new DockerService.
func NewDockerService(logger *logrus.Logger, configService ConfigManager) (DockerService, error) {
	composePath, err := findComposeFile()
	if err != nil {
		return nil, fmt.Errorf("failed to find docker-compose.yml: %w", err)
	}

	if err := validateComposePath(composePath); err != nil {
		return nil, fmt.Errorf("invalid docker-compose file: %w", err)
	}

	return &dockerService{
		logger:        logger,
		composePath:   filepath.Clean(composePath),
		configPath:    configService.GetConfigPath(),
		configService: configService,
	}, nil
}

// Start starts the docker container using docker-compose.
func (s *dockerService) Start() error {
	//nolint:gosec // validateComposePath() and filepath.Clean() in-use.
	cmd := exec.Command("docker", "compose", "-f", s.composePath, "up", "-d", "--pull", "always")
	cmd.Env = s.getComposeEnv()

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to start containers: %w\nOutput: %s", err, string(output))
	}

	s.logger.Info("Service started successfully")

	return nil
}

// Stop stops and removes the docker container using docker-compose.
func (s *dockerService) Stop() error {
	// Stop and remove containers, volumes, and networks
	//nolint:gosec // validateComposePath() and filepath.Clean() in-use.
	cmd := exec.Command("docker", "compose", "-f", s.composePath, "down",
		"--remove-orphans",
		"-v",
		"--rmi", "local",
		"--timeout", "30")
	cmd.Env = s.getComposeEnv()

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to stop containers: %w\nOutput: %s", err, string(output))
	}

	s.logger.Info("Service stopped and cleaned up successfully")

	return nil
}

// IsRunning checks if the docker container is running.
func (s *dockerService) IsRunning() (bool, error) {
	//nolint:gosec // validateComposePath() and filepath.Clean() in-use.
	cmd := exec.Command("docker", "compose", "-f", s.composePath, "ps", "--format", "{{.State}}")
	cmd.Env = s.getComposeEnv()

	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to check container status: %w", err)
	}

	states := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, state := range states {
		if strings.Contains(strings.ToLower(state), "running") {
			return true, nil
		}
	}

	return false, nil
}

// Update pulls the latest image and restarts the container.
func (s *dockerService) Update() error {
	cfg := s.configService.Get()

	//nolint:gosec // validateComposePath() and filepath.Clean() in-use.
	cmd := exec.Command("docker", "pull", fmt.Sprintf("ethpandaops/contributoor:%s", cfg.Version))
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to pull image: %w\nOutput: %s", err, string(output))
	}

	s.logger.WithField("version", cfg.Version).Infof(
		"%sImage updated successfully%s",
		tui.TerminalColorGreen,
		tui.TerminalColorReset,
	)

	return nil
}

func (s *dockerService) getComposeEnv() []string {
	cfg := s.configService.Get()

	return append(os.Environ(),
		fmt.Sprintf("CONTRIBUTOOR_CONFIG_PATH=%s", s.configService.GetConfigDir()),
		fmt.Sprintf("CONTRIBUTOOR_VERSION=%s", cfg.Version),
	)
}

func findComposeFile() (string, error) {
	// Get binary directory
	ex, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("could not get executable path: %w", err)
	}

	binDir := filepath.Dir(ex)

	// Check release mode (next to binary)
	composePath := filepath.Join(binDir, "docker-compose.yml")
	if _, e := os.Stat(composePath); e == nil {
		return composePath, nil
	}

	// Check dev mode paths
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("could not get working directory: %w", err)
	}

	// Try current directory
	if _, err := os.Stat(filepath.Join(cwd, "docker-compose.yml")); err == nil {
		return filepath.Join(cwd, "docker-compose.yml"), nil
	}

	// Try repo root
	if _, err := os.Stat(filepath.Join(cwd, "..", "..", "docker-compose.yml")); err == nil {
		return filepath.Join(cwd, "..", "..", "docker-compose.yml"), nil
	}

	return "", fmt.Errorf("docker-compose.yml not found")
}

func validateComposePath(path string) error {
	// Check if path exists and is a regular file
	fi, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("invalid compose file path: %w", err)
	}

	if fi.IsDir() {
		return fmt.Errorf("compose path is a directory, not a file")
	}

	// Ensure absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Verify file extension
	if !strings.HasSuffix(strings.ToLower(absPath), ".yml") &&
		!strings.HasSuffix(strings.ToLower(absPath), ".yaml") {
		return fmt.Errorf("compose file must have .yml or .yaml extension")
	}

	return nil
}
