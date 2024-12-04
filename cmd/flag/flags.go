package flag

import (
	"fmt"
	cfg "github.com/pete911/flowlogs/internal/aws"
	"github.com/spf13/cobra"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

var (
	Global    Flags
	logLevels = map[string]slog.Level{"debug": slog.LevelDebug, "info": slog.LevelInfo, "warn": slog.LevelWarn, "error": slog.LevelError}
)

type Flags struct {
	Region   string
	logLevel string
}

func (f Flags) Logger() *slog.Logger {
	if level, ok := logLevels[strings.ToLower(f.logLevel)]; ok {
		opts := &slog.HandlerOptions{Level: level}
		return slog.New(slog.NewJSONHandler(os.Stderr, opts))
	}

	fmt.Printf("invalid log level %s", f.logLevel)
	os.Exit(1)
	return nil
}

func (f Flags) AWSConfig() cfg.Config {
	cfg, err := cfg.NewConfig(f.Region)
	if err != nil {
		fmt.Printf("new aws config: %v\n", err)
		os.Exit(1)
	}
	return cfg
}

func InitPersistentFlags(cmd *cobra.Command, flags *Flags) {
	cmd.PersistentFlags().StringVar(
		&flags.Region,
		"region",
		getStringEnv("REGION", ""),
		"aws region",
	)
	cmd.PersistentFlags().StringVar(
		&flags.logLevel,
		"log-level",
		"info",
		"log level - debug, info, warn, error",
	)
}

func getStringEnv(envName string, defaultValue string) string {
	env, ok := os.LookupEnv(fmt.Sprintf("AWSFL_%s", envName))
	if !ok {
		return defaultValue
	}
	return env
}

func getBoolEnv(envName string, defaultValue bool) bool {
	env, ok := os.LookupEnv(fmt.Sprintf("AWSFL_%s", envName))
	if !ok {
		return defaultValue
	}
	if out, err := strconv.ParseBool(env); err == nil {
		return out
	}
	return defaultValue
}

func getIntEnv(envName string, defaultValue int) int {
	env, ok := os.LookupEnv(fmt.Sprintf("AWSFL_%s", envName))
	if !ok {
		return defaultValue
	}
	if out, err := strconv.Atoi(env); err == nil {
		return out
	}
	return defaultValue
}
