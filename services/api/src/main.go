package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

const (
	valuesFile                = "scripts/values"
	messageInitiationSize int = 148
	messageResponseSize   int = 92
)

var defaults = map[string]int{
	"JCMIN": 0,
	"JCMAX": 10,
	"JMIN":  10,
	"JMAX":  10,
}

const (
	jcmax = 10
)

// func validateDefaults(values

func getDefaultValue(envKey string, defaultValue int) int {
	value, err := strconv.ParseInt(os.Getenv(envKey), 10, 32)
	if err != nil {
		log.Printf("Cannot parse '%s', fallback to default %d\n", envKey, defaultValue)
		return defaultValue
	} else if (strings.Contains(envKey, "MIN") && value < 0) || (strings.Contains(envKey, "MAX") && value > int64(defaultValue)) {
		log.Printf("'%d'is not legal, fallback to default %d\n", value, defaultValue)
		return defaultValue
	}
	return int(value)
}

func generateConfigValues() {
	// Defaults
	junkPacketCountMin := getDefaultValue("JUNK_PACKET_COUNT_MIN", 3)
	junkPacketCountMax := getDefaultValue("JUNK_PACKET_COUNT_MAX", 10)
	junkPacketMinSize := getDefaultValue("JUNK_PACKET_MIN_SIZE", 50)
	junkPacketMaxSize := getDefaultValue("JUNK_PACKET_MAX_SIZE", 1000)
	config := make(map[string]int)
	config["$JUNK_PACKET_COUNT"] = randIntBound(junkPacketCountMin, junkPacketCountMax)
	config["$JUNK_PACKET_MIN_SIZE"] = junkPacketMinSize
	config["$JUNK_PACKET_MAX_SIZE"] = junkPacketMaxSize
	s1 := randIntBound(15, 150)
	s2 := randIntBound(15, 150)
	for s1+messageInitiationSize == s2+messageResponseSize {
		s2 = randIntBound(15, 150)
	}

	config["$INIT_PACKET_JUNK_SIZE"] = s1
	config["$RESPONSE_PACKET_JUNK_SIZE"] = s2
	var headersValue [4]int
	for i := range headersValue {
		headersValue[i] = randIntBound(5, math.MaxInt32)
	}

	config["$INIT_PACKET_MAGIC_HEADER"] = headersValue[0]
	config["$RESPONSE_PACKET_MAGIC_HEADER"] = headersValue[1]
	config["$UNDERLOAD_PACKET_MAGIC_HEADER"] = headersValue[2]
	config["$TRANSPORT_PACKET_MAGIC_HEADER"] = headersValue[3]
	writeConfigFile(config)
}

func writeConfigFile(config map[string]int) {
	values, err1 := os.ReadFile(valuesFile)
	if err1 != nil {
		log.Fatal(err1)
	}
	replaceStrings(string(values), valuesFile, config)
}

func replaceStrings(input string, outputFile string, config map[string]int) {
	for i, j := range config {
		input = strings.ReplaceAll(input, i, strconv.Itoa(j))
		os.WriteFile(outputFile, []byte(input), 0o600)
	}
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
	cliArgs := os.Args[1:]
	if len(cliArgs) != 1 {
		log.Fatal("Provide exactly one argument! Pass -h for help")
	}

	switch cliArgs[0] {
	case "-c", "--config":
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
