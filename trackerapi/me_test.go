package trackerapi

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func TestTrackerFile(t *testing.T) {
	// Should pass
	currentUser.Login("jasondemps1@gmail.com", "PAssw0rd12345")

	err := parse(makeRequest())

	if err != nil {
		t.Error("Could not parse request")
	}

	err = ioutil.WriteFile(FileLocation, []byte(currentUser.APIToken), 0644)

	if err != nil {
		t.Errorf("Error writing file to: %s\n", FileLocation)
	}

	file, ok := ioutil.ReadFile(FileLocation)

	if ok != nil {
		t.Errorf("Tracker API file was not written to %s\n", FileLocation)
	} else {
		t.Logf("File output: %s\n", file)
		fmt.Printf("File output: %s\n", file)
	}
}

func TestSetCredentials(t *testing.T) {
	// Should fail, credentials aren't set
	setCredentials()
	err := parse(makeRequest())

	if err != nil {
		t.Error("Parsed invalid request")
	}

	// Should pass
	currentUser.Login("jasondemps1@gmail.com", "PAssw0rd12345")

	err = parse(makeRequest())

	if err != nil {
		t.Error("Could not parse request")
	}
}
