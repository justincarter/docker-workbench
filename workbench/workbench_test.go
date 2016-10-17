package workbench

import "testing"

func TestValidProxyIP_Valid(t *testing.T) {
	ips := []string{
		"192.168.0.1",
		"10.0.0.1",
		"172.10.10.10",
	}
	for _, k := range ips {
		valid := validProxyIP(k)
		if !valid {
			t.Fail()
		}
	}
}

func TestValidProxyIP_Invalid(t *testing.T) {
	ips := []string{
		"invalid string",
		"::1",
		"192.168.99.1",
	}
	for _, k := range ips {
		valid := validProxyIP(k)
		if valid {
			t.Fail()
		}
	}
}
