package fdtd

import (
	"io/ioutil"
	"strings"
)

func ListAllFile(dir string, suffix string) ([]string, error) {
	var res []string
	if files, err := ioutil.ReadDir(dir); err != nil {
		return res, nil
	} else {
		for _, f := range files {
			if strings.HasSuffix(f.Name(), suffix) {
				res = append(res, f.Name())
			}
		}
		return res, nil
	}
}
