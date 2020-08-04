package fdtd

import (
	"os"
	"path"
	"strconv"
	"testing"
)

func TestListfile(t *testing.T) {
	temp := os.TempDir()
	tp := path.Join(temp, "testListfile")
	os.MkdirAll(tp, 0777)
	for i := 1; i <= 10; i++ {
		os.Create(path.Join(tp, strconv.Itoa(i)+".fsp"))
	}
	l, _ := ListAllFile(tp, ".fsp")
	if len(l) != 10 {
		t.Errorf("unexpected num")
	}
	l, _ = ListAllFile(tp, ".fs")
	if len(l) != 0 {
		t.Errorf("unexpected num")
	}

}
