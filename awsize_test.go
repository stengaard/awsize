package main

import (
	"fmt"
	"net"
	"testing"
)

func TestHostsInNet(t *testing.T) {
	template := "10.0.0.0/%d"

	expect := map[int]int{
		30: 2,
		19: 8190,
		18: 16382,
		17: 32766,
		16: 65534,
	}
	for input, expected := range expect {
		inString := fmt.Sprintf(template, input)
		_, n, err := net.ParseCIDR(inString)
		if err != nil {
			t.Errorf("error parsing cidr : %q - %s", inString, err)
			continue
		}
		output := hostsInNet(n)
		if output != expected {
			t.Errorf("%s returned %d hosts - expected %d", inString, output, expected)
		}
	}

}
