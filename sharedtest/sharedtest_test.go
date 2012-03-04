// Package sharedtest implements shared UpdateStore interface tests
//
// The USTest* functions tests should be run by implementations
// of the UpdateStore interface to ensure compatability between them.
package sharedtest

import (
	"partisci/store"
	"partisci/version"
	"testing"
	"time"
)

func USTestAppSummary(s store.UpdateStore, t *testing.T) {
	v := version.Version{App: "app1", Ver: "ver", Host: "a"}
	v.Prepare()
	s.Update(v)

	if _, ok := s.App("non-existant"); ok {
		t.Error("got ok for non-existant AppId")
	}
	if as, ok := s.App("app1"); ok {
		if as.HostCount != 1 {
			t.Error("expected HostCount: 1, actual: ", as.HostCount)
		}
	} else {
		t.Error("missing expected AppId")
	}

	v2 := version.Version{App: "app1", Ver: "ver", Host: "b"}
	v2.Prepare()
	s.Update(v2)
	if as, ok := s.App("app1"); ok {
		if as.HostCount != 2 {
			t.Error("expected HostCount: 2, actual: ", as.HostCount)
		}
	}
}

func USTestHostSummary(s store.UpdateStore, t *testing.T) {
	v := version.Version{App: "app1", Ver: "ver", Host: "a"}
	v.Prepare()
	s.Update(v)

	if _, ok := s.Host("non-existant"); ok {
		t.Error("got ok for non-existant Host")
	}
	if as, ok := s.Host("a"); ok {
		if as.AppCount != 1 {
			t.Error("expected AppCount: 1, actual: ", as.AppCount)
		}
	} else {
		t.Error("missing expected Host")
	}

	v2 := version.Version{App: "app2", Ver: "ver", Host: "a"}
	v2.Prepare()
	s.Update(v2)
	if as, ok := s.Host("a"); ok {
		if as.AppCount != 2 {
			t.Error("expected AppCount: 2, actual: ", as.AppCount)
		}
	}
}

// test Clear() & Update() interactions
func USTestClearUpdate(s store.UpdateStore, t *testing.T) {
	if len(s.Versions("", "", "")) > 0 {
		t.Error("Versions should be empty")
	}
	v := *new(version.Version)
	v.Prepare()
	s.Update(v)
	if len(s.Versions("", "", "")) != 1 {
		t.Error(v.ExactUpdate, v.LastUpdate)
		t.Error("Versions should have one entry")
	}
	s.Clear()
	if len(s.Versions("", "", "")) > 0 {
		t.Error("Versions should be empty")
	}
	s.Update(v)
	if len(s.Versions("", "", "")) > 0 {
		t.Error("updates older than threshold should be discarded")
	}
}

func USTestTrim(s store.UpdateStore, t *testing.T) {
	// setup one version in the future and a few more
	v1a := version.Version{App: "app1", Ver: "ver", Host: "a"}
	v1a.Prepare()
	v1a.ExactUpdate = v1a.ExactUpdate.Add(time.Duration(10 * time.Second))
	s.Update(v1a)

	v1b := version.Version{App: "app1", Ver: "ver", Host: "b"}
	v1b.Prepare()
	s.Update(v1b)

	v2 := version.Version{App: "app2", Ver: "ver", Host: "a"}
	v2.Prepare()
	s.Update(v2)

	// sanity check
	if l := len(s.Versions("", "", "")); l != 3 {
		t.Fatal("before: version count - expected: 2, actual: ", l)
	}
	if l := len(s.Hosts()); l != 2 {
		t.Fatal("before: host count - expected: 2, actual: ", l)
	}
	if l := len(s.Apps()); l != 2 {
		t.Fatal("before: app count - expected: 2, actual: ", l)
	}

	// trim every version before 1 second in the future of one version
	count := s.Trim(v2.ExactUpdate.Add(time.Duration(1 * time.Second)))
	if count != 2 {
		t.Fatal("after: trim should have removed 2 versions")
	}
	if l := len(s.Versions("", "", "")); l != 1 {
		t.Fatal("after: version count - expected: 1, actual: ", l)
	}
	if l := len(s.Hosts()); l != 1 {
		t.Fatal("after: host count - expected: 1, actual: ", l)
	}
	if l := len(s.Apps()); l != 1 {
		t.Fatal("after: app count - expected: 1, actual: ", l)
	}

	// trim every version
	count = s.Trim(v2.ExactUpdate.Add(time.Duration(20 * time.Second)))
	if count != 1 {
		t.Fatal("after all: trim should have removed the last one version")
	}
	if l := len(s.Versions("", "", "")); l != 0 {
		t.Fatal("after all: version count - expected: 0, actual: ", l)
	}
	if l := len(s.Hosts()); l != 0 {
		t.Fatal(s.Hosts())
		t.Fatal("after all: host count - expected: 0, actual: ", l)
	}
	if l := len(s.Apps()); l != 0 {
		t.Fatal(s.Apps())
		t.Fatal("after all: app count - expected: 0, actual: ", l)
	}
}
