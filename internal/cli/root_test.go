package cli

import (
	"reflect"
	"testing"
)

func TestPrepareArgsPhrase(t *testing.T) {
	got := PrepareArgs([]string{"--json", "cep 01001000"})
	want := []string{"--json", "__phrase", "cep 01001000"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("PrepareArgs() = %#v, want %#v", got, want)
	}
}

func TestPrepareArgsSubcommand(t *testing.T) {
	got := PrepareArgs([]string{"cep", "01001000"})
	want := []string{"cep", "01001000"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("PrepareArgs() = %#v, want %#v", got, want)
	}
}
