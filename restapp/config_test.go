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
	t.Log("env.App.RsServer", env.App.RsServer)
	t.Log("env.App.Template", env.App.Template)

	fsm, err := env.GetFsManager()
	if err != nil {
		t.Fatal(err)
	} else {
		for id, fs := range fsm.Map() {
			t.Logf("env.fsManager[%s] = %s", id, fs.Name())
		}
	}
}
