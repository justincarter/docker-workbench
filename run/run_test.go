package run

import (
	"os"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	err := Run("echo", "test", "pass")
	if err != nil {
		t.Fail()
	}
}

func TestOutput(t *testing.T) {
	out, err := Output("echo", "test", "pass")
	if len(out) == 0 || err != nil {
		t.Fail()
	}

}

func TestVBoxManagePath_Path(t *testing.T) {
	testVBoxManagePath(t)
}

func TestVBoxManagePath_InstallPath(t *testing.T) {
	os.Setenv("VBOX_INSTALL_PATH", "/dummy/install/path")
	testVBoxManagePath(t)
}

func TestVBoxManagePath_MSIInstallPath(t *testing.T) {
	os.Setenv("VBOX_MSI_INSTALL_PATH", "/dummy/msi/install/path")
	testVBoxManagePath(t)
}

func testVBoxManagePath(t *testing.T) {
	path := VBoxManagePath()
	s := strings.Split(path, string(os.PathSeparator))
	if s[len(s)-1] != "VBoxManage" {
		t.Fail()
	}
}
