package cmd

import "testing"

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
