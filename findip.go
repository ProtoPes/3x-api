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
    var unusedIP string

    IPsBytes, err1 := os.ReadFile(unusedIPs)
    if err1 != nil {
        log.Fatal(err1)
        os.Exit(1)
    }

    IPSstring := strings.Split(string(IPsBytes), "\n")

    unusedIP = IPSstring[0]
    IPSstring = slices.Delete(IPSstring, 0, 1)

    os.WriteFile(unusedIPs, []byte(strings.Join(IPSstring, "\n")), 0666)

    return unusedIP
}
