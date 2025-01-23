package main

import (
	"fmt"
	"log"
	"log/slog"
	"math"
	"os"
	"strconv"
	"strings"
)

// default values
const (
	templatesFolder       = "../templates"
	paramFile             = templatesFolder + "/parameters.txt"
	clientTemplate        = templatesFolder + "/client.conf"
	serverTemplate        = templatesFolder + "/wg0.conf"
	valuesFile            = "scripts/values"
	messageInitiationSize = 148
	messageResponseSize   = 92
	jcmin                 = 3
	jcmax                 = 10
	jmin                  = 50
	jmax                  = 1000
)

func tryReadFile(name string) string {
	/* Try to read a file, exit with 1 if something wrong */
	res, err := os.ReadFile(name)
	if err != nil {
		log.Fatal(err)
	}
	return string(res)
}

func getDefaultValue(envKey string, defaultValue int) int {
	/* If environment variable is set, parse it to int */
	key, exists := os.LookupEnv(envKey)
	if !exists {
		return defaultValue
	}
	value, err := strconv.ParseInt(key, 10, 32)
	if err != nil {
		slog.Info("Cannot parse env value, fallback to default ", envKey, defaultValue)
		return defaultValue
	} else if (strings.Contains(envKey, "MIN") && value < 0) || (strings.Contains(envKey, "MAX") && value > int64(defaultValue)) {
		slog.Info("Value is not legal, fallback to default: ", "given", value, envKey, defaultValue)
		return defaultValue
	}
	return int(value)
}

func generateConfigValues() {
	// Defaults
	junkPacketCountMin := getDefaultValue("JUNK_PACKET_COUNT_MIN", jcmin)
	junkPacketCountMax := getDefaultValue("JUNK_PACKET_COUNT_MAX", jcmax)
	junkPacketMinSize := getDefaultValue("JUNK_PACKET_MIN_SIZE", jmin)
	junkPacketMaxSize := getDefaultValue("JUNK_PACKET_MAX_SIZE", jmax)
	os.Setenv("Jc", strconv.Itoa(randIntBound(junkPacketCountMin, junkPacketCountMax)))
	os.Setenv("Jmin", strconv.Itoa(junkPacketMinSize))
	os.Setenv("Jmax", strconv.Itoa(junkPacketMaxSize))
	s1 := randIntBound(15, 150)
	s2 := randIntBound(15, 150)
	for s1+messageInitiationSize == s2+messageResponseSize {
		s2 = randIntBound(15, 150)
	}

	os.Setenv("S1", strconv.Itoa(s1))
	os.Setenv("S2", strconv.Itoa(s2))
	for i := range 4 {
		os.Setenv("H"+strconv.Itoa(i+1), strconv.Itoa(randIntBound(5, math.MaxInt32)))
	}

	param := os.ExpandEnv(tryReadFile(paramFile))
	os.WriteFile("parameters", []byte(param), 0o600)
	os.Setenv("Parameters", param)
	os.WriteFile("wg0.conf", []byte(os.ExpandEnv(tryReadFile(serverTemplate))), 0o600)
	os.WriteFile("client.conf", []byte(os.ExpandEnv(tryReadFile(clientTemplate))), 0o600)

}

func showUsageMessage(message string) {
	fmt.Println(message)
	fmt.Println("Program accepts one of these arguments:")
	fmt.Println("-c or --config : generate config files from template")
	fmt.Println("-n or --name: generate ranfom name")
	fmt.Println("-i or --ip: find unused ip for client")
	fmt.Println("-g or --gen-ip: generate unused ip adresses file")
	fmt.Println("-h or --help: show usage")
}

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	slog.Debug("Start of program")
	cliArgs := os.Args[1:]
	if len(cliArgs) != 1 {
		log.Fatal("Provide exactly one argument! Pass -h for help")
	}

	switch cliArgs[0] {
	case "init":
		generateConfigValues()
	case "-n", "--name":
		fmt.Println(GetRandomName())
	case "-i", "--ip":
		fmt.Println(findUnusedIP())
	case "-g", "--gen-ip":
		generateUnusedIPs()
	case "-h", "--help":
		showUsageMessage("AWG configs generator. Usage:")
	default:
		log.Fatal("Unrecognised flag! Pass -h for help")
	}
}
