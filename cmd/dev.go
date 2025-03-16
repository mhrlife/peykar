package cmd

import (
	"fmt"
	"github.com/mhrlife/peykar/internal"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
)

// devCmd represents the dev command
var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Runs your local work as a peykar server",
	Run: func(cmd *cobra.Command, args []string) {
		logger := slog.Default()

		if !isValidPeykarProject(logger) {
			logger.Error("Not a valid Peykar project")

			return
		}

		moduleDir, err := os.Getwd()
		if err != nil {
			logger.Error("Failed to get current working directory", "err", err)

			return
		}

		tempDir, err := os.MkdirTemp("", "peykar-dev-*")
		if err != nil {
			logger.Error("Failed to create temporary directory", "err", err)

			return
		}

		defer os.RemoveAll(tempDir)

		modulePath, err := buildModuleTemp(moduleDir, tempDir)
		if err != nil {
			logger.Error("Failed to build module", "err", err)

			return
		}

		logger.Info("✅ Module built.", "path", modulePath)

		pluginManager := internal.NewPluginManager(logger)

		if err := pluginManager.LoadPlugin(modulePath); err != nil {
			logger.Error("Failed to load plugin", "err", err)

			return
		}

		fmt.Println(pluginManager.All())
	},
}

func buildModuleTemp(moduleDir, tempDir string) (string, error) {
	outputPath := filepath.Join(tempDir, "module.so")

	pluginDir, _ := os.Getwd()

	cmd := exec.Command(
		"go",
		"build",
		"-buildmode=plugin",
		"-mod=vendor",
		"--trimpath",
		"-gcflags", fmt.Sprintf("-trimpath %s", pluginDir),
		"-asmflags", fmt.Sprintf("-trimpath %s", pluginDir),
		"-o",
		outputPath,
	)

	cmd.Dir = moduleDir
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("build failed: %w", err)
	}

	return outputPath, nil
}

func isValidPeykarProject(logger *slog.Logger) bool {
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		logger.Error("Not a valid Peykar project", "err", err)

		return false
	}

	mainFiles, _ := filepath.Glob("main.go")
	mainGo := len(mainFiles) > 0
	if !mainGo {
		logger.Error("Not a valid Peykar project", "err", "main.go not found")
	}

	return mainGo
}

func init() {
	rootCmd.AddCommand(devCmd)
}
