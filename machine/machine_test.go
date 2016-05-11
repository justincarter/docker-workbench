package machine

import (
	"reflect"
	"testing"
)

func TestParseEnvOutput(t *testing.T) {

	input := `export DOCKER_TLS_VERIFY="1"
export DOCKER_HOST="tcp://192.168.99.100:2376"
export DOCKER_CERT_PATH="d:\docker\machines\workbench"
export DOCKER_MACHINE_NAME="workbench"
# Run this command to configure your shell:
# eval $("C:\Program Files\Docker Toolbox\docker-machine.exe" env workbench)
`
	result := parseEnvOutput([]byte(input))

	expected := map[string]string{
		"DOCKER_TLS_VERIFY":   "1",
		"DOCKER_HOST":         "tcp://192.168.99.100:2376",
		"DOCKER_CERT_PATH":    "d:\\docker\\machines\\workbench",
		"DOCKER_MACHINE_NAME": "workbench",
	}

	if !reflect.DeepEqual(expected, result) {
		t.Fail()
	}

}
