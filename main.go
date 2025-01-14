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
    template = "templates/values.template"
    valuesFile = "values.txt"
)


func generateConfigValues() {
    const messageInitiationSize int = 148
    const messageResponseSize int = 92
    config := make(map[string]string)
    config["$AWG_SERVER_PORT"] = strconv.Itoa(randIntBound(30000, 50000))
    config["$JUNK_PACKET_COUNT"] = strconv.Itoa(randIntBound(3, 10))
    config["$JUNK_PACKET_MIN_SIZE"] = strconv.Itoa(50)
    config["$JUNK_PACKET_MAX_SIZE"] = strconv.Itoa(1000)
    s1 := randIntBound(15, 150)
    s2 := randIntBound(15, 150)
    for (s1 + messageInitiationSize == s2 + messageResponseSize) {
        s2 = randIntBound(15, 150)
    }
    config["$INIT_PACKET_JUNK_SIZE"] = strconv.Itoa(s1)
    config["$RESPONSE_PACKET_JUNK_SIZE"] = strconv.Itoa(s2)
    var headersValue [4]string
    var max int = math.MaxInt32
    for i := range headersValue  {
        headersValue[i] = strconv.Itoa(randIntBound(5, max))
    }
    config["$INIT_PACKET_MAGIC_HEADER"] = headersValue[0];
    config["$RESPONSE_PACKET_MAGIC_HEADER"] = headersValue[1];
    config["$UNDERLOAD_PACKET_MAGIC_HEADER"] = headersValue[2];
    config["$TRANSPORT_PACKET_MAGIC_HEADER"] = headersValue[3];
    config["$AWG_SUBNET_IP"] = "10.8.1.0"
    config["$WIREGUARD_SUBNET_CIDR"] = "24"
    writeConfigFile(config)
}

func writeConfigFile(config map[string]string) {

    values, err1 := os.ReadFile(template)
    if err1 != nil {
        log.Fatal(err1)
    }

    replaceStrings(string(values), valuesFile, config)


}

func replaceStrings(input string, outputFile string, config map[string]string) {
    for i, j := range config {
        input = strings.ReplaceAll(input, i, j)
        os.WriteFile(outputFile, []byte(input), 0666)
    }
}

func showUsageError() {
    fmt.Println("Program accepts one of these arguments:")
    fmt.Println("-c or --config : generate config files from template")
    fmt.Println("-n or --name: generate ranfom name")
    fmt.Println("-i or --ip: find unused ip for client")
    fmt.Println("-g or --gen-ip: generate unused ip adresses file")
}

func main() {
    cliArgs := os.Args[1:]
    if len(cliArgs) != 1 {
        showUsageError()
        os.Exit(1)
    }

    switch cliArgs[0] {
        case "-c", "--config": generateConfigValues()
        case "-n", "--name": fmt.Println(GetRandomName())
        case "-i", "--ip": fmt.Println(findUnusedIP())
        case "-g", "--gen-ip": generateUnusedIPs()
    }

}
