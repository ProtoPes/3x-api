package main

import (
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

const unusedIPs = "unused_ips.txt"

func generateUnusedIPs() {
    log.Default().Println("Creating ", unusedIPs, " file")
    var IPs [254]string
    j := 2
    for i := range IPs {
        IPs[i] = strconv.Itoa(j)
        j++
    }

    output := strings.Join(IPs[0:254], "\n")
    os.WriteFile(unusedIPs, []byte(output), 0666)
}

func findUnusedIP() string {

    ipsBytes, err1 := os.ReadFile(unusedIPs)
    if err1 != nil {
        log.Fatal(err1)
    }

    index := slices.Index(ipsBytes, 10)
    var nextIP []byte = ipsBytes[:index]
    os.WriteFile(unusedIPs, ipsBytes[index + 1:], 0666)

    return string(nextIP)
}
