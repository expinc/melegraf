package smoke

import (
	"os"
	"path"
	"testing"

	"github.com/expinc/melegraf/test/util"
)

func TestEmptyConfig(t *testing.T) {
	cfgStr := "{}"
	workDir, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	env, err := util.SetupEnvironment(path.Join(workDir, "testdir"), cfgStr)
	defer util.TearDownEnvironment(env)
	if err != nil {
		t.Error(err)
	}

	// FIXME: health check
}
