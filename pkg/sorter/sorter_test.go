package sorter

import (
	"io/ioutil"
	"os"
	"testing"
)

const (
	maxGoRoutines = 1000
	maxMemory     = 2048
)

func logChannel(t *testing.T) chan string {
	messages := make(chan string)
	go func() {
		for msg := range messages {
			t.Log(msg)
		}
	}()
	return messages
}

func testSort(t *testing.T, input string) {
	messages := logChannel(t)
	defer close(messages)

	output, err := ioutil.TempFile("", "output.idx")
	if err != nil {
		t.Errorf("temp file error %s", err)
		return
	}
	defer os.Remove(output.Name())

	tempDir, err := ioutil.TempDir("", "leakdb_")
	if err != nil {
		t.Errorf("Temp error: %s\n", err)
		return
	}
	defer os.RemoveAll(tempDir)

	err = Start(messages, input, output.Name(), maxMemory, maxGoRoutines, tempDir, false)
	if err != nil {
		t.Errorf("Sort error: %s\n", err)
		return
	}

	output.Seek(0, 0)
	sorted, err := CheckSort(messages, output.Name(), false)
	if err != nil {
		t.Errorf("Check sort error: %s\n", err)
		return
	}
	if !sorted {
		t.Error("Failed to correctly sort index")
		return
	}
}

func TestSorterSmallEmail(t *testing.T) {
	testSort(t, "../../test/small-email-unsorted.idx")
}

func TestSorterLargeEmail(t *testing.T) {
	testSort(t, "../../test/large-email-unsorted.idx")
}

func TestSorterSmallUser(t *testing.T) {
	testSort(t, "../../test/small-user-unsorted.idx")
}

func TestSorterLargeUser(t *testing.T) {
	testSort(t, "../../test/large-user-unsorted.idx")
}

func TestSorterSmallDomain(t *testing.T) {
	testSort(t, "../../test/small-domain-unsorted.idx")
}

func TestSorterLargeDomain(t *testing.T) {
	testSort(t, "../../test/large-domain-unsorted.idx")
}