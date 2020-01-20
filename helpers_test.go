package main

import (
	"testing"
)

func TestBase(t *testing.T) {
	want := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if base != want {
		t.Errorf("base = \"%s\", want = \"%s\"", base, want)
	}
}

type mockDriver struct{}

var mDriver *mockDriver

func (m *mockDriver) Get(s string) (string, int, error) {
	return s, 200, nil
}

func (m *mockDriver) Set(short, long string) (int, error) {
	return 200, nil
}

func TestDBGet(t *testing.T) {
	s, i, err := dbGet(mDriver, "this is a test")
	if s != "this is a test" {
		t.Errorf("want this is a test, but got %s", s)
	}
	if i != 200 {
		t.Errorf("want 200, but got %d", i)
	}
	if err != nil {
		t.Errorf("want nil, but got %v", err.Error())
	}
}

func TestDBSet(t *testing.T) {
	i, err := dbSet(mDriver, "foo", "bar")
	if i != 200 {
		t.Errorf("want 200, but got %d", i)
	}
	if err != nil {
		t.Errorf("want nil, but got %v", err.Error())
	}
}

func TestStrToInt(t *testing.T) {
	result := strToInt("a")
	if result != 10 {
		t.Errorf("want 10, but got %d", result)
	}
}
