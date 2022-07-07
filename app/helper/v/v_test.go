package v

import (
	"testing"
)

func Test(t *testing.T) {
	a := P("string")
	if V(a) != "string" {
		t.Error("error")
	}

	a = nil
	if V(a) != "" {
		t.Error("error")
	}

	b := P(struct{ key string }{key: "val"})

	if V(b).key != "val" {
		t.Error("error")
	}

	b = nil
	if V(b).key != "" {
		t.Error("error")
	}

	var ct map[string]string
	c := P(ct)
	if c == nil || *c != nil {
		t.Error("error")
	}
}
