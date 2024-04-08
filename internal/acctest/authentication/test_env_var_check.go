package authentication

import (
	"errors"
	"os"
	"testing"
)

func TestEnvVarSlice(vars []string, fileName string, t *testing.T) map[string]string {
	var errorSlice []error
	var envMap = make(map[string]string)
	for _, v := range vars {
		if os.Getenv(v) == "" {
			errorSlice = append(errorSlice, errors.New(v+" is not set for "+fileName))
		} else {
			envMap[v] = os.Getenv(v)
		}
	}
	if len(errorSlice) > 0 {
		for _, err := range errorSlice {
			t.Error(err)
		}
		t.FailNow()
	}
	return envMap
}
