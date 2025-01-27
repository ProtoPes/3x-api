package main

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
)

// Default constant values
const (
	envLogging            = "LOG_LEVEL"
	templatesFolder       = "../templates"
	paramFile             = templatesFolder + "/parameters.txt"
	clientTemplate        = templatesFolder + "/client.conf"
	serverTemplate        = templatesFolder + "/wg0.conf"
	messageInitiationSize = 148
	messageResponseSize   = 92
	jcmin                 = 3
	jcmax                 = 10
	jmin                  = 50
	jmax                  = 1000
)

// Structs
type file struct {
	name        string
	content     []byte
	permissions os.FileMode
}

// Creating new logger
var (
	LogLevel = &slog.LevelVar{}
	Logger   = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     LogLevel,
	}))
	levels = map[string]slog.Level{
		"DEBUG": slog.LevelDebug,
		"INFO":  slog.LevelInfo,
		"WARN":  slog.LevelWarn,
		"ERROR": slog.LevelError,
	}
)

// Set logger to [Default] and parse environment [LOG_LEVEL] to set appropriate
// level
func configureLogger(logger *slog.Logger, logLevel *slog.LevelVar) {
	slog.SetDefault(logger)
	envLevel, isSet := os.LookupEnv(envLogging)
	if !isSet {
		return
	}
	newLevel, exists := levels[strings.ToUpper(envLevel)]
	var msg string
	if !exists {
		msg = "Invalid"
	} else {
		msg = "Set"
		logLevel.Set(newLevel)
	}
	logger.Warn(msg, envLogging, envLevel)
}

// Try to write files
func writeFiles(files ...file) error {
	for _, j := range files {
		Logger.Debug("Saving", "filename", j.name)
		err := os.WriteFile(j.name, j.content, j.permissions)
		if err != nil {
			return err
		}
		Logger.Debug("Success")
	}
	return nil
}

// Try to read a file, exit with 1 if something wrong
func tryReadFile(name string) string {
	res, err := os.ReadFile(name)
	checkError(err)
	return string(res)
}

// Check if error occured and if so, exit program with non-zero code
func checkError(err error) {
	if err != nil {
		Logger.Error(err.Error())
		panic(err)
	}
}

// If environment variable is set, parse it to int
func getDefaultValue(envKey string, defaultValue int) int {
	key, exists := os.LookupEnv(envKey)
	if !exists {
		return defaultValue
	}
	value, err := strconv.ParseInt(key, 10, 32)
	if err != nil {
		Logger.Warn(envKey+" Couldn't be parsed", "Given", key, "Fallback to", defaultValue)
		return defaultValue
	} else if (strings.Contains(envKey, "MIN") && value < 0) || (strings.Contains(envKey, "MAX") && value > int64(defaultValue)) {
		Logger.Warn(envKey+" is invalid", "Given", value, "Fallback to", defaultValue)
		return defaultValue
	}
	Logger.Info("Set", envKey, value)
	return int(value)
}
