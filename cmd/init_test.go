package cmd

import (
	"fmt"
	"os"
	"testing"
)

func mockGetRegion() string {
	return "foo"
}

func TestNewConfig(t *testing.T) {

	expected_backup := "us-east-1"
	expected_main := "us-west-2"

	actual := newConfig("us-west-2", "us-east-1")

	if actual.MainRegion != expected_main {
		t.Errorf("main region actual %s, expected %s", actual.MainRegion, expected_main)
	}

	if actual.BackupRegion != expected_backup {
		t.Errorf("backup region actual %s, expected %s", actual.BackupRegion, expected_backup)
	}
}

func TestGenConfig(t *testing.T) {
	br := mockGetRegion
	mr := mockGetRegion
	expected := Config{
		"foo",
		"foo",
	}
	actual := genConfig(mr, br)
	if actual != expected {
		t.Errorf("expected %v got %v", expected, actual)
	}
}

func TestWriteConfig(t *testing.T) {
	filename := "/tmp/config.json"
	conf := Config{
		"foo",
		"bar",
	}
	writeConfig(conf, filename)
	dat, err := os.ReadFile(filename)
	if err != nil {
		t.Errorf("File not created")
	}
	if string(dat) != "fwf" {
		fmt.Print(string(dat))
	}
	os.Remove(filename)
}

func TestReadConfig(t *testing.T) {
	filename := "/tmp/config.json"
	mconf := Config{
		"foo",
		"bar",
	}
	writeConfig(mconf, filename)
	conf, err := readConfig(filename)
	if err != nil {
		t.Errorf("got error: %s", err)
	}
	if conf.MainRegion != "foo" {
		t.Errorf("%s, %s", conf.MainRegion, conf.BackupRegion)
	}
	os.Remove(filename)
}
