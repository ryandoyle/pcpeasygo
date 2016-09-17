package pmapi

import (
	"testing"
	"reflect"
)

func TestPmapiContext_PmGetContextHostname(t *testing.T) {
	c, _ := PmNewContext(PmContextHost, "localhost")
	hostname, _ := c.PmGetContextHostname()

	assertEquals(t, hostname, "ryandesktop")
}

func TestPmNewContext_withAnInvalidHostHasANilContext(t *testing.T) {
	c, _ := PmNewContext(PmContextHost, "not-a-host")

	assertNil(t, c)
}

func TestPmNewContext_withAnInvalidHostHasAnError(t *testing.T) {
	_, err := PmNewContext(PmContextHost, "not-a-host")

	assertNotNil(t, err)
}

func TestPmNewContext_hasAValidContext(t *testing.T) {
	c, _ := PmNewContext(PmContextHost, "localhost")

	assertNotNil(t, c)
}

func TestPmNewContext_hasANilError(t *testing.T) {
	_, err := PmNewContext(PmContextHost, "localhost")

	assertNil(t, err)
}

func TestPmNewContext_supportsALocalContext(t *testing.T) {
	c, _ := PmNewContext(PmContextLocal, "")

	assertNotNil(t, c)
}

func assertEquals(t *testing.T, a interface{}, b interface{}) {
	if(a != b) {
		t.Errorf("expected %v, got %v", b, a)
	}
}

func assertNil(t *testing.T, a interface{}) {
	if ! isNil(a) {
		t.Errorf("expected nil but got %v", a)
	}
}

func assertNotNil(t *testing.T, a interface{}) {
	if isNil(a) {
		t.Errorf("expected not nil but got %v", a)
	}
}

func isNil(object interface{}) bool {
	if object == nil {
		return true
	}
	return reflect.ValueOf(object).IsNil()
}
