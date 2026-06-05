package intent

import (
	"errors"
	"fmt"
	"strings"
)

const (
	ActionCEP          = "cep"
	ActionCNPJLookup   = "cnpj_lookup"
	ActionCNPJGenerate = "cnpj_generate"
	ActionCNPJValidate = "cnpj_validate"
	ActionCNPJFormat   = "cnpj_format"
	ActionCPFGenerate  = "cpf_generate"
	ActionCPFValidate  = "cpf_validate"
	ActionCPFFormat    = "cpf_format"
	ActionBank         = "bank"
	ActionBanks        = "banks"
	ActionDDD          = "ddd"
	ActionHolidays     = "holidays"
	ActionUFs          = "ufs"
	ActionCities       = "cities"
)

type Intent struct {
	Action string
	Args   []string
	Raw    string
}

func (i Intent) Arg(index int) string {
	if index < 0 || index >= len(i.Args) {
		return ""
	}
	return i.Args[index]
}

func Parse(input string) (Intent, error) {
	raw := strings.TrimSpace(input)
	if raw == "" {
		return Intent{}, errors.New("informe uma acao, ex: br \"cep 01001000\"")
	}

	tokens := strings.Fields(raw)
	key := NormalizeWord(tokens[0])

	switch key {
	case "cep":
		return singleArg(raw, tokens, ActionCEP, "informe um CEP, ex: br \"cep 01001000\"")
	case "cnpj":
		return parseCNPJ(raw, tokens)
	case "cpf":
		return parseCPF(raw, tokens)
	case "banco":
		return singleArg(raw, tokens, ActionBank, "informe o codigo do banco, ex: br \"banco 341\"")
	case "bancos":
		return Intent{Action: ActionBanks, Raw: raw}, nil
	case "ddd":
		return singleArg(raw, tokens, ActionDDD, "informe um DDD, ex: br \"ddd 16\"")
	case "feriado", "feriados":
		return singleArg(raw, tokens, ActionHolidays, "informe o ano, ex: br \"feriados 2026\"")
	case "uf", "ufs", "estado", "estados":
		return Intent{Action: ActionUFs, Raw: raw}, nil
	case "cidade", "cidades", "municipio", "municipios":
		return singleArg(raw, tokens, ActionCities, "informe uma UF, ex: br \"cidades SP\"")
	default:
		return Intent{}, fmt.Errorf("acao desconhecida: %q", tokens[0])
	}
}

func parseCNPJ(raw string, tokens []string) (Intent, error) {
	if len(tokens) < 2 {
		return Intent{}, errors.New("informe um CNPJ ou uma acao: gerar, validar ou formatar")
	}

	op := NormalizeWord(tokens[1])
	switch op {
	case "gerar", "gera", "generate":
		return Intent{Action: ActionCNPJGenerate, Raw: raw}, nil
	case "validar", "valida", "validate", "check":
		if len(tokens) < 3 {
			return Intent{}, errors.New("informe o CNPJ para validar")
		}
		return Intent{Action: ActionCNPJValidate, Args: []string{tokens[2]}, Raw: raw}, nil
	case "formatar", "formata", "format":
		if len(tokens) < 3 {
			return Intent{}, errors.New("informe o CNPJ para formatar")
		}
		return Intent{Action: ActionCNPJFormat, Args: []string{tokens[2]}, Raw: raw}, nil
	default:
		return Intent{Action: ActionCNPJLookup, Args: []string{tokens[1]}, Raw: raw}, nil
	}
}

func parseCPF(raw string, tokens []string) (Intent, error) {
	if len(tokens) < 2 {
		return Intent{}, errors.New("informe uma acao para CPF: gerar, validar ou formatar")
	}

	op := NormalizeWord(tokens[1])
	switch op {
	case "gerar", "gera", "generate":
		return Intent{Action: ActionCPFGenerate, Raw: raw}, nil
	case "validar", "valida", "validate", "check":
		if len(tokens) < 3 {
			return Intent{}, errors.New("informe o CPF para validar")
		}
		return Intent{Action: ActionCPFValidate, Args: []string{tokens[2]}, Raw: raw}, nil
	case "formatar", "formata", "format":
		if len(tokens) < 3 {
			return Intent{}, errors.New("informe o CPF para formatar")
		}
		return Intent{Action: ActionCPFFormat, Args: []string{tokens[2]}, Raw: raw}, nil
	default:
		return Intent{}, fmt.Errorf("acao de CPF desconhecida: %q", tokens[1])
	}
}

func singleArg(raw string, tokens []string, action string, message string) (Intent, error) {
	if len(tokens) < 2 {
		return Intent{}, errors.New(message)
	}
	return Intent{Action: action, Args: []string{tokens[1]}, Raw: raw}, nil
}

func NormalizeWord(word string) string {
	word = strings.ToLower(strings.Trim(word, " \t\r\n:;,"))
	replacements := map[rune]rune{
		'á': 'a', 'à': 'a', 'â': 'a', 'ã': 'a',
		'é': 'e', 'ê': 'e',
		'í': 'i',
		'ó': 'o', 'ô': 'o', 'õ': 'o',
		'ú': 'u',
		'ç': 'c',
	}

	var b strings.Builder
	for _, r := range word {
		if replacement, ok := replacements[r]; ok {
			b.WriteRune(replacement)
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}
