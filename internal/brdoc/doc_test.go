package brdoc

import "testing"

func TestCPF(t *testing.T) {
	if !ValidateCPF("529.982.247-25") {
		t.Fatal("expected valid CPF")
	}
	if ValidateCPF("111.111.111-11") {
		t.Fatal("expected repeated CPF to be invalid")
	}

	formatted, err := FormatCPF("52998224725")
	if err != nil {
		t.Fatalf("FormatCPF() error = %v", err)
	}
	if formatted != "529.982.247-25" {
		t.Fatalf("formatted CPF = %q", formatted)
	}

	generated, err := GenerateCPF()
	if err != nil {
		t.Fatalf("GenerateCPF() error = %v", err)
	}
	if !ValidateCPF(generated) {
		t.Fatalf("generated CPF is invalid: %s", generated)
	}
}

func TestCNPJ(t *testing.T) {
	if !ValidateCNPJ("04.252.011/0001-10") {
		t.Fatal("expected valid CNPJ")
	}
	if ValidateCNPJ("11.111.111/1111-11") {
		t.Fatal("expected repeated CNPJ to be invalid")
	}

	formatted, err := FormatCNPJ("04252011000110")
	if err != nil {
		t.Fatalf("FormatCNPJ() error = %v", err)
	}
	if formatted != "04.252.011/0001-10" {
		t.Fatalf("formatted CNPJ = %q", formatted)
	}

	generated, err := GenerateCNPJ()
	if err != nil {
		t.Fatalf("GenerateCNPJ() error = %v", err)
	}
	if !ValidateCNPJ(generated) {
		t.Fatalf("generated CNPJ is invalid: %s", generated)
	}
}
