package consulkv

import (
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	c, err := NewClient([]string{"10.10.28.2:8500"}, "http")
	if err != nil {
		t.Errorf("Construct ConsulClient failed: %s", err.Error())
	}

	// ------- string, int
	v, err := c.Get("dev/test", strings.ToLower)
	if err != nil {
		t.Errorf("Get string failed: %s", err.Error())
	}
	r := v.Validate(func(s string) bool {
		return strings.HasPrefix(s, "1")
	})
	if !r {
		t.Error("Validate failed.")
	}
	s := v.String()
	if s != "123" {
		t.Errorf("Get string failed: expect '123' but %s", s)
	}
	i := v.MustInt()
	if i != 123 {
		t.Errorf("Convert Int failed: expect 123 but %v", i)
	}

	// ------- json object
	v1, err := c.Get("dev/json_obj", nil)
	if err != nil {
		t.Errorf("Get failed: %s", err.Error())
	}
	m, err := v1.JsonObject()
	if err != nil {
		t.Errorf("Get json object failed: %s", err.Error())
	}
	if m["foo"] != "foo" {
		t.Errorf("Unexpected json object: %v", m)
	}

	m1 := v1.MustJsonObject()
	if m1["foo"] != "foo" {
		t.Errorf("Unexpected json object: %v", m1)
	}

	// ------- json array
	v2, err := c.Get("dev/json_array", nil)
	if err != nil {
		t.Errorf("Get failed: %s", err.Error())
	}
	arr, err := v2.JsonArray()
	if err != nil {
		t.Errorf("Get json array failed: %s", err.Error())
	}
	if arr[0] != "foo" {
		t.Errorf("Unexpected json array: %v", arr)
	}
	arr1 := v2.MustJsonArray()
	if arr1[0] != "foo" {
		t.Errorf("Unexpected json array: %v", arr1)
	}

	// ---- Default value
	i1 := v2.MustInt(5)
	if i1 != 5 {
		t.Error("Test default value failed.")
	}

	// ------ Not exist
	_, err = c.Get("not_exist", nil)
	if err != ErrNotExist {
		t.Error("Test getting not exist key failed")
	}

}
