package greetings

import (
	"regexp"
	"strings"
	"testing"
)

func TestHelloName(t *testing.T) {
	name := "Chris"
	want := regexp.MustCompile(`\b` + name + `\b`)

	msg, err := Hello(name)

	if !want.MatchString(msg) || err != nil {
		t.Fatalf(`Hello("Chris") = %q, %v, want match for %#q, nil`, msg, err, want)
	}
}

func TestHelloEmpty(t *testing.T) {
	msg, err := Hello("")

	if msg != "" || err == nil {
		t.Fatalf(`Hello("") = %q, %v, want "", error`, msg, err)
	}
}

func TestHellosNames(t *testing.T) {
	names := []string{
		"Chris",
		"John",
	}

	wants := make(map[string]*regexp.Regexp)

	for _, name := range names {
		wants[name] = regexp.MustCompile(`\b` + name + `\b`)
	}

	msgs, err := Hellos(names)

	for name, want := range wants {
		if !want.MatchString(msgs[name]) || err != nil {
			t.Fatalf(`Hellos([]string{%v}) = %q, %v, want match for %#q, nil`, strings.Join(names, ", "), msgs[name], err, want)
		}
	}
}
