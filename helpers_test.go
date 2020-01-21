package main

import (
	"errors"
	"testing"
)

func TestBase(t *testing.T) {
	want := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	got := base
	if got != want {
		t.Errorf("want: %q, but got: %q", want, base)
	}
}

type mockDriver struct {
	returnString string
	returnVal    int
	returnErr    error
}

func (m *mockDriver) Get(s string) (string, int, error) {
	return m.returnString, m.returnVal, m.returnErr
}

func (m *mockDriver) Set(short, long string) (int, error) {
	return m.returnVal, m.returnErr
}

func TestDBGet_string(t *testing.T) {
	mDriver := &mockDriver{
		returnString: "http://foo.example.com",
		returnVal:    200,
		returnErr:    nil,
	}
	want := mDriver.returnString
	got, i, err := dbGet(mDriver, want)
	if got != want {
		t.Errorf("want: %q, but got: %q", want, got)
	}
	if i != 200 {
		t.Errorf("want: 200, but got: %d", i)
	}
	if err != nil {
		t.Errorf("want: nil, but got: %v", err)
	}
}

func TestDBGet_err(t *testing.T) {
	mDriver := &mockDriver{
		returnString: "",
		returnVal:    404,
		returnErr:    errors.New("whoops, something wrong happens"),
	}
	want := mDriver.returnString
	got, i, err := dbGet(mDriver, want)
	if got != want {
		t.Errorf("want: %q, but got: %q", want, got)
	}
	if i != 404 {
		t.Errorf("want: 404, but got: %d", i)
	}
	if err != nil {
		if err.Error() != mDriver.returnErr.Error() {
			t.Errorf("want: %q, but got: %q", mDriver.returnErr, err.Error())
		}
	}
}

func TestDBSet_string(t *testing.T) {
	mDriver := &mockDriver{
		returnString: "",
		returnVal:    200,
		returnErr:    nil,
	}
	want := 200
	got, err := dbSet(mDriver, "foo", "bar")
	if got != want {
		t.Errorf("want: %d, but got: %d", want, got)
	}
	if err != nil {
		t.Errorf("want: nil, but got: %v", err)
	}
}

func TestDBSet_err(t *testing.T) {
	mDriver := &mockDriver{
		returnString: "",
		returnVal:    404,
		returnErr:    errors.New("whoops, something wrong happens"),
	}
	want := 404
	got, err := dbSet(mDriver, "foo", "bar")
	if got != want {
		t.Errorf("want: %d, but got: %d", want, got)
	}
	if err != nil {
		if err.Error() != mDriver.returnErr.Error() {
			t.Errorf("want: %q, but got: %q", mDriver.returnErr, err.Error())
		}
	}
}

func TestStrToInt_shortString(t *testing.T) {
	want := 10
	got := strToInt("a")
	if got != want {
		t.Errorf("want: %d, but got: %d", want, got)
	}
}

func TestStrToInt_longString(t *testing.T) {
	want := strToInt("https://azertyuiopqsdfghjklmwxcvbnAZERTYUIOPQSDFGHJKLMWXCVBNazertyuiopqsdfghjklmwxcvbnAZERTYUIOPQSDFGHJKLMWXCVBNazertyuiopqsdfghjklmwxcvbnAZERTYUIOPQSDFGHJKLMWXCVBNazertyuiopqsdfghjklmwxcvbnAZERTYUIOPQSDFGHJKLMWXCVBNazertyuiopqsdfghjklmwxcvbnAZERTYUIOPQSDFGHJKLMWXCVBNazertyuiopqsdfghjklmwxcvbnAZERTYUIOPQSDFGHJKLMWXCVBNazertyuiopqsdfghjklmwxcvbnAZERTYUIOPQSDFGHJKLMWXCVBNazertyuiopqsdfghjklmwxcvbnAZERTYUIOPQSDFGHJKLMWXCVBNazertyuiopqsdfghjklmwxcvbnAZERTYUIOPQSDFGHJKLMWXCVBN")
	got := strToInt("https://nbvcxwmlkjhgfdsqpoiuytrezaNBVCXWMLKJHGFDSQPOIUYTREZAnbvcxwmlkjhgfdsqpoiuytrezaNBVCXWMLKJHGFDSQPOIUYTREZAnbvcxwmlkjhgfdsqpoiuytrezaNBVCXWMLKJHGFDSQPOIUYTREZAnbvcxwmlkjhgfdsqpoiuytrezaNBVCXWMLKJHGFDSQPOIUYTREZAnbvcxwmlkjhgfdsqpoiuytrezaNBVCXWMLKJHGFDSQPOIUYTREZAnbvcxwmlkjhgfdsqpoiuytrezaNBVCXWMLKJHGFDSQPOIUYTREZAnbvcxwmlkjhgfdsqpoiuytrezaNBVCXWMLKJHGFDSQPOIUYTREZAnbvcxwmlkjhgfdsqpoiuytrezaNBVCXWMLKJHGFDSQPOIUYTREZAnbvcxwmlkjhgfdsqpoiuytrezaNBVCXWMLKJHGFDSQPOIUYTREZA")
	if got == want {
		t.Errorf("%d must be different from %d", want, got)
	}
}

func TestStrToInt_rune(t *testing.T) {
	want := 0
	got := strToInt("://*%%*$$^&`")
	if got == want {
		t.Errorf("want %d, but got %d", want, got)
	}
}
