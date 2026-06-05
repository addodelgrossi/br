package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/addodelgrossi/br/internal/brasilapi"
	"github.com/addodelgrossi/br/internal/brdoc"
	"github.com/addodelgrossi/br/internal/intent"
	"github.com/addodelgrossi/br/internal/output"
)

type App struct {
	Client *brasilapi.Client
	Output output.Options
	Writer outputWriter
}

type outputWriter interface {
	Write(p []byte) (n int, err error)
}

func (a App) RunIntent(ctx context.Context, in intent.Intent) error {
	_ = ctx

	switch in.Action {
	case intent.ActionCEP:
		return a.printCEP(in.Arg(0))
	case intent.ActionCNPJLookup:
		return a.printCNPJ(in.Arg(0))
	case intent.ActionCNPJGenerate:
		return a.printGeneratedCNPJ()
	case intent.ActionCNPJValidate:
		return a.printCNPJValidation(in.Arg(0))
	case intent.ActionCNPJFormat:
		return a.printFormattedCNPJ(in.Arg(0))
	case intent.ActionCPFGenerate:
		return a.printGeneratedCPF()
	case intent.ActionCPFValidate:
		return a.printCPFValidation(in.Arg(0))
	case intent.ActionCPFFormat:
		return a.printFormattedCPF(in.Arg(0))
	case intent.ActionBank:
		return a.printBank(in.Arg(0))
	case intent.ActionBanks:
		return a.printBanks()
	case intent.ActionDDD:
		return a.printDDD(in.Arg(0))
	case intent.ActionHolidays:
		return a.printHolidays(in.Arg(0))
	case intent.ActionUFs:
		return a.printUFs()
	case intent.ActionCities:
		return a.printCities(in.Arg(0))
	default:
		return fmt.Errorf("acao sem handler: %s", in.Action)
	}
}

func (a App) printCEP(value string) error {
	cep, err := a.Client.CEP(value)
	if err != nil {
		return err
	}
	fields := []output.Field{
		{Key: "cep", Label: "CEP", Value: cep.CEP},
		{Key: "street", Label: "Logradouro", Value: cep.Street},
		{Key: "neighborhood", Label: "Bairro", Value: cep.Neighborhood},
		{Key: "city", Label: "Cidade", Value: cep.City},
		{Key: "state", Label: "UF", Value: cep.State},
		{Key: "service", Label: "Fonte", Value: cep.Service},
	}
	return output.PrintObject(a.Writer, a.Output, fields, cep, cep.CEP)
}

func (a App) printCNPJ(value string) error {
	cnpj, err := a.Client.CNPJ(value)
	if err != nil {
		return err
	}
	fields := []output.Field{
		{Key: "cnpj", Label: "CNPJ", Value: cnpj.CNPJ},
		{Key: "razao_social", Label: "Razao social", Value: cnpj.RazaoSocial},
		{Key: "nome_fantasia", Label: "Nome fantasia", Value: cnpj.NomeFantasia},
		{Key: "situacao", Label: "Situacao", Value: cnpj.DescricaoSituacaoCadastral},
		{Key: "municipio", Label: "Municipio", Value: cnpj.Municipio},
		{Key: "uf", Label: "UF", Value: cnpj.UF},
		{Key: "cnae", Label: "CNAE", Value: cnpj.CNAEFiscalDescricao},
	}
	return output.PrintObject(a.Writer, a.Output, fields, cnpj, cnpj.RazaoSocial)
}

func (a App) printGeneratedCPF() error {
	digits, err := brdoc.GenerateCPF()
	if err != nil {
		return err
	}
	formatted, err := brdoc.FormatCPF(digits)
	if err != nil {
		return err
	}
	raw := map[string]string{"cpf": formatted, "digits": digits}
	fields := []output.Field{
		{Key: "cpf", Label: "CPF", Value: formatted},
		{Key: "digits", Label: "Digitos", Value: digits},
	}
	return output.PrintObject(a.Writer, a.Output, fields, raw, formatted)
}

func (a App) printCPFValidation(value string) error {
	valid := brdoc.ValidateCPF(value)
	digits := brdoc.OnlyDigits(value)
	fields := []output.Field{
		{Key: "cpf", Label: "CPF", Value: value},
		{Key: "digits", Label: "Digitos", Value: digits},
		{Key: "valid", Label: "Valido", Value: boolString(valid)},
	}
	raw := map[string]any{"cpf": value, "digits": digits, "valid": valid}
	return output.PrintObject(a.Writer, a.Output, fields, raw, boolString(valid))
}

func (a App) printFormattedCPF(value string) error {
	formatted, err := brdoc.FormatCPF(value)
	if err != nil {
		return err
	}
	raw := map[string]string{"cpf": formatted, "digits": brdoc.OnlyDigits(value)}
	fields := []output.Field{{Key: "cpf", Label: "CPF", Value: formatted}}
	return output.PrintObject(a.Writer, a.Output, fields, raw, formatted)
}

func (a App) printGeneratedCNPJ() error {
	digits, err := brdoc.GenerateCNPJ()
	if err != nil {
		return err
	}
	formatted, err := brdoc.FormatCNPJ(digits)
	if err != nil {
		return err
	}
	raw := map[string]string{"cnpj": formatted, "digits": digits}
	fields := []output.Field{
		{Key: "cnpj", Label: "CNPJ", Value: formatted},
		{Key: "digits", Label: "Digitos", Value: digits},
	}
	return output.PrintObject(a.Writer, a.Output, fields, raw, formatted)
}

func (a App) printCNPJValidation(value string) error {
	valid := brdoc.ValidateCNPJ(value)
	digits := brdoc.OnlyDigits(value)
	fields := []output.Field{
		{Key: "cnpj", Label: "CNPJ", Value: value},
		{Key: "digits", Label: "Digitos", Value: digits},
		{Key: "valid", Label: "Valido", Value: boolString(valid)},
	}
	raw := map[string]any{"cnpj": value, "digits": digits, "valid": valid}
	return output.PrintObject(a.Writer, a.Output, fields, raw, boolString(valid))
}

func (a App) printFormattedCNPJ(value string) error {
	formatted, err := brdoc.FormatCNPJ(value)
	if err != nil {
		return err
	}
	raw := map[string]string{"cnpj": formatted, "digits": brdoc.OnlyDigits(value)}
	fields := []output.Field{{Key: "cnpj", Label: "CNPJ", Value: formatted}}
	return output.PrintObject(a.Writer, a.Output, fields, raw, formatted)
}

func (a App) printBank(value string) error {
	bank, err := a.Client.BankByCode(value)
	if err != nil {
		return err
	}
	fields := []output.Field{
		{Key: "code", Label: "Codigo", Value: bank.CodeString()},
		{Key: "name", Label: "Nome", Value: bank.Name},
		{Key: "full_name", Label: "Nome completo", Value: bank.FullName},
		{Key: "ispb", Label: "ISPB", Value: bank.ISPB},
	}
	return output.PrintObject(a.Writer, a.Output, fields, bank, bank.Name)
}

func (a App) printBanks() error {
	banks, err := a.Client.Banks()
	if err != nil {
		return err
	}
	rows := make([]map[string]string, 0, len(banks))
	for _, bank := range banks {
		rows = append(rows, map[string]string{
			"code":      bank.CodeString(),
			"name":      bank.Name,
			"full_name": bank.FullName,
			"ispb":      bank.ISPB,
		})
	}
	columns := []output.Column{
		{Key: "code", Label: "Codigo"},
		{Key: "name", Label: "Nome"},
		{Key: "ispb", Label: "ISPB"},
	}
	return output.PrintList(a.Writer, a.Output, columns, rows, banks, "name")
}

func (a App) printDDD(value string) error {
	ddd, err := a.Client.DDD(value)
	if err != nil {
		return err
	}
	fields := []output.Field{
		{Key: "ddd", Label: "DDD", Value: brdoc.OnlyDigits(value)},
		{Key: "state", Label: "UF", Value: ddd.State},
		{Key: "cities", Label: "Cidades", Value: strings.Join(ddd.Cities, ", ")},
	}
	return output.PrintObject(a.Writer, a.Output, fields, ddd, ddd.State)
}

func (a App) printHolidays(year string) error {
	holidays, err := a.Client.Holidays(year)
	if err != nil {
		return err
	}
	rows := make([]map[string]string, 0, len(holidays))
	for _, holiday := range holidays {
		rows = append(rows, map[string]string{
			"date": holiday.Date,
			"name": holiday.Name,
			"type": holiday.Type,
		})
	}
	columns := []output.Column{
		{Key: "date", Label: "Data"},
		{Key: "name", Label: "Nome"},
		{Key: "type", Label: "Tipo"},
	}
	return output.PrintList(a.Writer, a.Output, columns, rows, holidays, "date")
}

func (a App) printUFs() error {
	ufs, err := a.Client.UFs()
	if err != nil {
		return err
	}
	rows := make([]map[string]string, 0, len(ufs))
	for _, uf := range ufs {
		rows = append(rows, map[string]string{
			"id":     fmt.Sprintf("%d", uf.ID),
			"sigla":  uf.Sigla,
			"nome":   uf.Nome,
			"regiao": uf.Region.Nome,
		})
	}
	columns := []output.Column{
		{Key: "sigla", Label: "UF"},
		{Key: "nome", Label: "Nome"},
		{Key: "regiao", Label: "Regiao"},
	}
	return output.PrintList(a.Writer, a.Output, columns, rows, ufs, "sigla")
}

func (a App) printCities(uf string) error {
	cities, err := a.Client.Municipalities(uf)
	if err != nil {
		return err
	}
	rows := make([]map[string]string, 0, len(cities))
	for _, city := range cities {
		rows = append(rows, map[string]string{
			"codigo_ibge": city.CodigoIBGE.String(),
			"nome":        city.Nome,
		})
	}
	columns := []output.Column{
		{Key: "codigo_ibge", Label: "IBGE"},
		{Key: "nome", Label: "Nome"},
	}
	return output.PrintList(a.Writer, a.Output, columns, rows, cities, "nome")
}

func boolString(value bool) string {
	if value {
		return "true"
	}
	return "false"
}
