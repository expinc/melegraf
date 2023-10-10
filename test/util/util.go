package util

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

type TestEnvironment struct {
	tmpDir  string
	process *os.Process
}

func SetupEnvironment(tmpDir string, cfgStr string) (TestEnvironment, error) {
	env := TestEnvironment{tmpDir: tmpDir}
	err := os.RemoveAll(tmpDir)
	if err != nil {
		return env, err
	}
	err = os.MkdirAll(tmpDir, os.ModeDir)
	if err != nil {
		return env, err
	}

	cfgPath := path.Join(tmpDir, "melegraf.json")
	err = ioutil.WriteFile(cfgPath, []byte(cfgStr), os.ModePerm)
	if err != nil {
		return env, err
	}
	env.process, err = LaunchMelegraf(cfgPath)
	if err != nil {
		return env, err
	}

	return env, nil
}

func LaunchMelegraf(cfgPath string) (*os.Process, error) {
	gopath, err := exec.LookPath("go")
	if err != nil {
		return nil, err
	}
	return os.StartProcess(gopath, []string{"go", "run", "main.go", "--config", cfgPath}, &os.ProcAttr{})
}

func TearDownEnvironment(env TestEnvironment) error {
	err := env.process.Kill()
	if err != nil {
		return err
	}

	return nil
}
