package restapp

import "testing"

func Test_NewConfigByFile(t *testing.T) {
	cfg, err := NewConfigByFile("./config_test.yaml")
	if err != nil {
		t.Fatal(err)
		return
	}
	env, err := cfg.GetEnvConfig("dev_lxd")
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log("env.fs", env.Fs)
	t.Log("env.App.HServer", env.App.HServer)
	t.Log("env.App.Template", env.App.Template)

}
