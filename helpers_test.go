package main

import (
	"errors"
	"testing"
)

func TestBase(t *testing.T) {
	want := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	got := base
	if got != want {
		t.Errorf("const base = %q, but want %q", base, want)
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
		t.Errorf("dbGet() returned string = %q, want: %q", got, want)
	}
	if i != 200 {
		t.Errorf("dbGet() returned int = %d, want: 200 ", i)
	}
	if err != nil {
		t.Errorf("dbGet() returned err = %s,want: %s", err, mDriver.returnErr)
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
		t.Errorf("dbGet() returned string = %q, want: %q", got, want)
	}
	if i != 404 {
		t.Errorf("dbGet() returned int = %d, want: 404 ", i)
	}
	if err != nil {
		if err != mDriver.returnErr {
			t.Errorf("dbGet() returned err = %s, want: %s", err, mDriver.returnErr)
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
		t.Errorf("dbSet() returned int = %d, want: %d ", got, want)
	}
	if err != nil {
		t.Errorf("dbSet() returned err = %s, want: %s", err, mDriver.returnErr)
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
		t.Errorf("dbSet() returned int = %d, want: %d ", got, want)
	}
	if err != nil {
		if err.Error() != mDriver.returnErr.Error() {
			t.Errorf("dbSet() returned err = %s, want: %s", err, mDriver.returnErr)
		}
	}
}

func TestStrToInt_shortString(t *testing.T) {
	want := 10
	got := strToInt("a")
	if got != want {
		t.Errorf("strToInt(shrotString) = %d, want: %d", got, want)
	}
}

func TestStrToInt_longString(t *testing.T) {
	want := strToInt("https://azertyuiopqsdfghjklmwxcvbnAZERTYUIOPQSDFGHJKLMWXCVBNazertyuiopqsdfghjklmwxcvbnAZERTYUIOPQSDFGHJKLMWXCVBNazertyuiopqsdfghjklmwxcvbnAZERTYUIOPQSDFGHJKLMWXCVBNazertyuiopqsdfghjklmwxcvbnAZERTYUIOPQSDFGHJKLMWXCVBNazertyuiopqsdfghjklmwxcvbnAZERTYUIOPQSDFGHJKLMWXCVBNazertyuiopqsdfghjklmwxcvbnAZERTYUIOPQSDFGHJKLMWXCVBNazertyuiopqsdfghjklmwxcvbnAZERTYUIOPQSDFGHJKLMWXCVBNazertyuiopqsdfghjklmwxcvbnAZERTYUIOPQSDFGHJKLMWXCVBNazertyuiopqsdfghjklmwxcvbnAZERTYUIOPQSDFGHJKLMWXCVBN")
	got := strToInt("https://nbvcxwmlkjhgfdsqpoiuytrezaNBVCXWMLKJHGFDSQPOIUYTREZAnbvcxwmlkjhgfdsqpoiuytrezaNBVCXWMLKJHGFDSQPOIUYTREZAnbvcxwmlkjhgfdsqpoiuytrezaNBVCXWMLKJHGFDSQPOIUYTREZAnbvcxwmlkjhgfdsqpoiuytrezaNBVCXWMLKJHGFDSQPOIUYTREZAnbvcxwmlkjhgfdsqpoiuytrezaNBVCXWMLKJHGFDSQPOIUYTREZAnbvcxwmlkjhgfdsqpoiuytrezaNBVCXWMLKJHGFDSQPOIUYTREZAnbvcxwmlkjhgfdsqpoiuytrezaNBVCXWMLKJHGFDSQPOIUYTREZAnbvcxwmlkjhgfdsqpoiuytrezaNBVCXWMLKJHGFDSQPOIUYTREZAnbvcxwmlkjhgfdsqpoiuytrezaNBVCXWMLKJHGFDSQPOIUYTREZA")
	if got == want {
		t.Errorf("strToInt(longString) = %d, must be different from %d", got, want)
	}
}

func TestStrToInt_idempotence(t *testing.T) {
	url := "https://www.google.com"
	want := strToInt(url)
	got := strToInt(url)
	if got != want {
		t.Errorf("strToInt(%s) = %d, want: %d", url, got, want)
	}
}

func TestStrToInt_rune(t *testing.T) {
	want := 0
	got := strToInt("://*%%*$$^&`")
	if got == want {
		t.Errorf("want %d, but got %d", want, got)
	}
}
