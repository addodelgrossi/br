package cli

import (
	"io"
	"strings"

	"github.com/addodelgrossi/br/internal/brasilapi"
	"github.com/addodelgrossi/br/internal/intent"
	"github.com/addodelgrossi/br/internal/output"
	"github.com/spf13/cobra"
)

type rootFlags struct {
	json    bool
	csv     bool
	quiet   bool
	noColor bool
	apiBase string
}

var knownCommands = map[string]bool{
	"cep":      true,
	"cnpj":     true,
	"cpf":      true,
	"banco":    true,
	"bancos":   true,
	"ddd":      true,
	"feriados": true,
	"feriado":  true,
	"ufs":      true,
	"uf":       true,
	"estados":  true,
	"cidades":  true,
	"cidade":   true,
}

func NewRootCommand(out io.Writer) *cobra.Command {
	flags := &rootFlags{apiBase: brasilapi.DefaultBaseURL}

	root := &cobra.Command{
		Use:           "br",
		Short:         "CLI brasileira para devs",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}
			return runPhrase(cmd, out, flags, strings.Join(args, " "))
		},
	}

	root.SetOut(out)
	root.PersistentFlags().BoolVar(&flags.json, "json", false, "imprime saida em JSON")
	root.PersistentFlags().BoolVar(&flags.csv, "csv", false, "imprime listas em CSV")
	root.PersistentFlags().BoolVarP(&flags.quiet, "quiet", "q", false, "imprime somente o valor principal")
	root.PersistentFlags().BoolVar(&flags.noColor, "no-color", false, "desativa cores")
	root.PersistentFlags().StringVar(&flags.apiBase, "api-base", brasilapi.DefaultBaseURL, "URL base da BrasilAPI")
	_ = root.PersistentFlags().MarkHidden("api-base")

	phrase := &cobra.Command{
		Use:    "__phrase [texto]",
		Hidden: true,
		Args:   cobra.ArbitraryArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runPhrase(cmd, out, flags, strings.Join(args, " "))
		},
	}
	root.AddCommand(phrase)
	root.AddCommand(cepCommand(out, flags))
	root.AddCommand(cnpjCommand(out, flags))
	root.AddCommand(cpfCommand(out, flags))
	root.AddCommand(bankCommand(out, flags))
	root.AddCommand(banksCommand(out, flags))
	root.AddCommand(dddCommand(out, flags))
	root.AddCommand(holidaysCommand(out, flags))
	root.AddCommand(ufsCommand(out, flags))
	root.AddCommand(citiesCommand(out, flags))

	return root
}

func PrepareArgs(args []string) []string {
	index := firstNonFlagIndex(args)
	if index < 0 {
		return args
	}

	candidate := args[index]
	key := intent.NormalizeWord(candidate)
	if strings.Contains(candidate, " ") || !knownCommands[key] {
		prepared := make([]string, 0, len(args)+1)
		prepared = append(prepared, args[:index]...)
		prepared = append(prepared, "__phrase")
		prepared = append(prepared, args[index:]...)
		return prepared
	}

	return args
}

func firstNonFlagIndex(args []string) int {
	skipNext := false
	for i, arg := range args {
		if skipNext {
			skipNext = false
			continue
		}
		if arg == "--" {
			if i+1 < len(args) {
				return i + 1
			}
			return -1
		}
		if strings.HasPrefix(arg, "--api-base") && !strings.Contains(arg, "=") {
			skipNext = true
			continue
		}
		if strings.HasPrefix(arg, "-") {
			continue
		}
		return i
	}
	return -1
}

func runPhrase(cmd *cobra.Command, out io.Writer, flags *rootFlags, phrase string) error {
	parsed, err := intent.Parse(phrase)
	if err != nil {
		return err
	}
	return newApp(out, flags).RunIntent(cmd.Context(), parsed)
}

func newApp(out io.Writer, flags *rootFlags) App {
	client := brasilapi.NewClient(nil)
	client.BaseURL = flags.apiBase
	return App{
		Client: client,
		Output: output.Options{
			JSON:    flags.json,
			CSV:     flags.csv,
			Quiet:   flags.quiet,
			NoColor: flags.noColor,
		},
		Writer: out,
	}
}

func runIntent(out io.Writer, flags *rootFlags, in intent.Intent) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		return newApp(out, flags).RunIntent(cmd.Context(), in)
	}
}

func cepCommand(out io.Writer, flags *rootFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "cep <cep>",
		Short: "Consulta um CEP",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return newApp(out, flags).RunIntent(cmd.Context(), intent.Intent{Action: intent.ActionCEP, Args: args})
		},
	}
}

func cnpjCommand(out io.Writer, flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cnpj <cnpj>",
		Short: "Consulta, gera, valida ou formata CNPJ",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return newApp(out, flags).RunIntent(cmd.Context(), intent.Intent{Action: intent.ActionCNPJLookup, Args: args})
		},
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "gerar",
		Short: "Gera um CNPJ valido",
		Args:  cobra.NoArgs,
		RunE:  runIntent(out, flags, intent.Intent{Action: intent.ActionCNPJGenerate}),
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "validar <cnpj>",
		Short: "Valida um CNPJ",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return newApp(out, flags).RunIntent(cmd.Context(), intent.Intent{Action: intent.ActionCNPJValidate, Args: args})
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "formatar <cnpj>",
		Short: "Formata um CNPJ",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return newApp(out, flags).RunIntent(cmd.Context(), intent.Intent{Action: intent.ActionCNPJFormat, Args: args})
		},
	})
	return cmd
}

func cpfCommand(out io.Writer, flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cpf",
		Short: "Gera, valida ou formata CPF",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "gerar",
		Short: "Gera um CPF valido",
		Args:  cobra.NoArgs,
		RunE:  runIntent(out, flags, intent.Intent{Action: intent.ActionCPFGenerate}),
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "validar <cpf>",
		Short: "Valida um CPF",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return newApp(out, flags).RunIntent(cmd.Context(), intent.Intent{Action: intent.ActionCPFValidate, Args: args})
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "formatar <cpf>",
		Short: "Formata um CPF",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return newApp(out, flags).RunIntent(cmd.Context(), intent.Intent{Action: intent.ActionCPFFormat, Args: args})
		},
	})
	return cmd
}

func bankCommand(out io.Writer, flags *rootFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "banco <codigo>",
		Short: "Consulta um banco pelo codigo",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return newApp(out, flags).RunIntent(cmd.Context(), intent.Intent{Action: intent.ActionBank, Args: args})
		},
	}
}

func banksCommand(out io.Writer, flags *rootFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "bancos",
		Short: "Lista bancos brasileiros",
		Args:  cobra.NoArgs,
		RunE:  runIntent(out, flags, intent.Intent{Action: intent.ActionBanks}),
	}
}

func dddCommand(out io.Writer, flags *rootFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "ddd <ddd>",
		Short: "Consulta cidades de um DDD",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return newApp(out, flags).RunIntent(cmd.Context(), intent.Intent{Action: intent.ActionDDD, Args: args})
		},
	}
}

func holidaysCommand(out io.Writer, flags *rootFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "feriados <ano>",
		Short: "Lista feriados nacionais",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return newApp(out, flags).RunIntent(cmd.Context(), intent.Intent{Action: intent.ActionHolidays, Args: args})
		},
	}
}

func ufsCommand(out io.Writer, flags *rootFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "ufs",
		Short: "Lista estados brasileiros",
		Args:  cobra.NoArgs,
		RunE:  runIntent(out, flags, intent.Intent{Action: intent.ActionUFs}),
	}
}

func citiesCommand(out io.Writer, flags *rootFlags) *cobra.Command {
	return &cobra.Command{
		Use:   "cidades <uf>",
		Short: "Lista municipios de uma UF",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return newApp(out, flags).RunIntent(cmd.Context(), intent.Intent{Action: intent.ActionCities, Args: args})
		},
	}
}
