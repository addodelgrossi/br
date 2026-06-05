package intent

import "testing"

func TestParse(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		action string
		arg    string
	}{
		{name: "cep", input: "cep 01001000", action: ActionCEP, arg: "01001000"},
		{name: "cnpj lookup", input: "cnpj 00.000.000/0001-91", action: ActionCNPJLookup, arg: "00.000.000/0001-91"},
		{name: "cnpj validar", input: "  cnpj   validar  11.222.333/0001-81  ", action: ActionCNPJValidate, arg: "11.222.333/0001-81"},
		{name: "cpf gerar", input: "cpf gerar", action: ActionCPFGenerate},
		{name: "banco", input: "banco 341", action: ActionBank, arg: "341"},
		{name: "bancos", input: "bancos", action: ActionBanks},
		{name: "ddd", input: "ddd 16", action: ActionDDD, arg: "16"},
		{name: "feriados", input: "feriados 2026", action: ActionHolidays, arg: "2026"},
		{name: "ufs", input: "estados", action: ActionUFs},
		{name: "cidades", input: "municipios SP", action: ActionCities, arg: "SP"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			if got.Action != tt.action {
				t.Fatalf("Action = %q, want %q", got.Action, tt.action)
			}
			if got.Arg(0) != tt.arg {
				t.Fatalf("Arg(0) = %q, want %q", got.Arg(0), tt.arg)
			}
		})
	}
}

func TestParseUnknown(t *testing.T) {
	_, err := Parse("wat 123")
	if err == nil {
		t.Fatal("expected error")
	}
}
