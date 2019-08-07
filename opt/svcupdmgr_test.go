package opt

import "testing"

func TestSvcUpdMgr_RemoteInfo(t *testing.T) {
	mgr := &SvcUpdMgr{}
	info, err := mgr.RemoteInfo()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", info)
}

func TestSvcUpdMgr_RemoteRestart(t *testing.T) {
	mgr := &SvcUpdMgr{}
	err := mgr.Restart("om")
	if err != nil {
		t.Fatal(err)
	}
}
