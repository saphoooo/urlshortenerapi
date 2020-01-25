package main

import (
	"testing"
)

func TestBase(t *testing.T) {
	want := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	got := base
	if got != want {
		t.Errorf("const base = %q, but want %q", base, want)
	}
}

func Test_strToInt_basic(t *testing.T) {
	tests := []struct {
		name string
		args string
		want int
	}{
		{"shortString", "a", 10},
		{"rune", "@#&é(')$*%<>", 2450613165425379093},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := strToInt(tt.args); got != tt.want {
				t.Errorf("strToInt(%q) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}

func Test_strToInt_long(t *testing.T) {
	tests := []struct {
		name string
		args string
		want int
	}{
		{"longString", "https://nbvcxwmlkjhgfdsqpoiuytrezaNBVCXWMLKJHGFDSQPOIUYTREZAnbvcxwmlkjhgfdsqpoiuytrezaNBVCXWMLKJHGFDSQPOIUYTREZAnbvcxwmlkjhgfdsqpoiuytrezaNBVCXWMLKJHGFDSQPOIUYTREZAnbvcxwmlkjhgfdsqpoiuytrezaNBVCXWMLKJHGFDSQPOIUYTREZAnbvcxwmlkjhgfdsqpoiuytrezaNBVCXWMLKJHGFDSQPOIUYTREZAnbvcxwmlkjhgfdsqpoiuytrezaNBVCXWMLKJHGFDSQPOIUYTREZAnbvcxwmlkjhgfdsqpoiuytrezaNBVCXWMLKJHGFDSQPOIUYTREZAnbvcxwmlkjhgfdsqpoiuytrezaNBVCXWMLKJHGFDSQPOIUYTREZAnbvcxwmlkjhgfdsqpoiuytrezaNBVCXWMLKJHGFDSQPOIUYTREZA", strToInt("https://azertyuiopqsdfghjklmwxcvbnAZERTYUIOPQSDFGHJKLMWXCVBNazertyuiopqsdfghjklmwxcvbnAZERTYUIOPQSDFGHJKLMWXCVBNazertyuiopqsdfghjklmwxcvbnAZERTYUIOPQSDFGHJKLMWXCVBNazertyuiopqsdfghjklmwxcvbnAZERTYUIOPQSDFGHJKLMWXCVBNazertyuiopqsdfghjklmwxcvbnAZERTYUIOPQSDFGHJKLMWXCVBNazertyuiopqsdfghjklmwxcvbnAZERTYUIOPQSDFGHJKLMWXCVBNazertyuiopqsdfghjklmwxcvbnAZERTYUIOPQSDFGHJKLMWXCVBNazertyuiopqsdfghjklmwxcvbnAZERTYUIOPQSDFGHJKLMWXCVBNazertyuiopqsdfghjklmwxcvbnAZERTYUIOPQSDFGHJKLMWXCVBN")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := strToInt(tt.args); got == tt.want {
				t.Errorf("strToInt(%q) = %v, want %v", tt.args, got, tt.want)
			}
		})
	}
}
