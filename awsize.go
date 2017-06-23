// Cmd awsize estimates the size of the AWS behemoth. It does this by analysing
// the number of available IPv4 address AWS claims to operate
package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"sort"
)

type Prefix struct {
	IPPrefix string `json:"ip_prefix"`
	Region   string `json:"region"`
	Service  string `json:"service"`
}

type AWSAddresses struct {
	SyncToken  string   `json:"syncToken"`
	CreateDate string   `json:"createDate"`
	Prefixes   []Prefix `json:"prefixes"`
}

const source = "https://ip-ranges.amazonaws.com/ip-ranges.json"

func main() {
	resp, err := http.Get(source)
	die(err)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		die(fmt.Errorf("JSON source returned %d - expected 200", resp.StatusCode))
	}

	var (
		dec   = json.NewDecoder(resp.Body)
		addrs = AWSAddresses{}
	)
	die(dec.Decode(&addrs))
	var (
		regionCounts  = map[string]int{}
		serviceCounts = map[string]int{}
		regionNames   = []string{}
		serviceNames  = []string{}
	)
	for _, prefix := range addrs.Prefixes {
		_, n, err := net.ParseCIDR(prefix.IPPrefix)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Bad CIDR string: %s - %s\n", prefix.IPPrefix, err)
			continue
		}

		region := prefix.Service + "/" + prefix.Region

		if _, ok := regionCounts[region]; !ok {
			regionNames = append(regionNames, region)
		}
		regionCounts[region] += hostsInNet(n)

		if _, ok := serviceCounts[prefix.Service]; !ok {
			serviceNames = append(serviceNames, prefix.Service)
		}
		serviceCounts[prefix.Service] += hostsInNet(n)

	}

	sort.Sort(sort.StringSlice(regionNames))
	sort.Sort(sort.StringSlice(serviceNames))

	fmt.Println("Service/Region Count")
	total := 0
	for _, s := range regionNames {
		fmt.Printf("%s %d\n", s, regionCounts[s])
		total += regionCounts[s]
	}
	fmt.Printf("Total: %d\n\n", total)

	fmt.Println("Service Count")
	total = 0
	for _, s := range serviceNames {
		fmt.Printf("%s %d\n", s, serviceCounts[s])
		total += serviceCounts[s]
	}
	fmt.Printf("Total: %d\n", total)

}

func hostsInNet(n *net.IPNet) int {
	ones, bits := n.Mask.Size()

	return (1 << uint(bits-ones)) - 2
}

func die(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %q\n", err)
	}
}
