package fdtd

import (
	"bufio"
	"os"
	"path"
	"testing"
)

func TestParserNormal(t *testing.T) {
	if file, err := os.Open(path.Join(os.Getenv("TestSampleDir"), "normal.out")); err != nil {
		t.Errorf("wrong sample file")
		return
	} else {
		defer file.Close()
		reader := bufio.NewReader(file)
		for {
			if line, err := reader.ReadString('\n'); err != nil {
				return
			} else {
				l := ParseStdOutLine(line)
				t.Logf("%+v", l)
			}
		}
	}
}

func TestParserError(t *testing.T) {
	if file, err := os.Open(path.Join(os.Getenv("TestSampleDir"), "error.out")); err != nil {
	} else {
		reader := bufio.NewReader(file)
		for {
			if line, err := reader.ReadString('\n'); err != nil {
				return
			} else {
				l := ParseStdOutLine(line)
				t.Logf("%+v", l)
			}
		}
	}
}
