package brasilapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const DefaultBaseURL = "https://brasilapi.com.br/api"

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}
	return &Client{
		BaseURL:    DefaultBaseURL,
		HTTPClient: httpClient,
	}
}

type CEP struct {
	CEP          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type CNPJ struct {
	CNPJ                         string  `json:"cnpj"`
	RazaoSocial                  string  `json:"razao_social"`
	NomeFantasia                 string  `json:"nome_fantasia"`
	DescricaoSituacaoCadastral   string  `json:"descricao_situacao_cadastral"`
	DataSituacaoCadastral        string  `json:"data_situacao_cadastral"`
	CNAEFiscalDescricao          string  `json:"cnae_fiscal_descricao"`
	Municipio                    string  `json:"municipio"`
	UF                           string  `json:"uf"`
	CapitalSocial                float64 `json:"capital_social"`
	DescricaoIdentificadorMatriz string  `json:"descricao_identificador_matriz_filial"`
}

type Bank struct {
	ISPB     string `json:"ispb"`
	Name     string `json:"name"`
	Code     *int   `json:"code"`
	FullName string `json:"fullName"`
}

func (b Bank) CodeString() string {
	if b.Code == nil {
		return ""
	}
	return strconv.Itoa(*b.Code)
}

type DDD struct {
	State  string   `json:"state"`
	Cities []string `json:"cities"`
}

type Holiday struct {
	Date string `json:"date"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type Region struct {
	ID    int    `json:"id"`
	Sigla string `json:"sigla"`
	Nome  string `json:"nome"`
}

type UF struct {
	ID     int    `json:"id"`
	Sigla  string `json:"sigla"`
	Nome   string `json:"nome"`
	Region Region `json:"regiao"`
}

type Municipio struct {
	Nome       string         `json:"nome"`
	CodigoIBGE FlexibleString `json:"codigo_ibge"`
}

type FlexibleString string

func (f *FlexibleString) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*f = ""
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*f = FlexibleString(s)
		return nil
	}

	var n json.Number
	if err := json.Unmarshal(data, &n); err == nil {
		*f = FlexibleString(n.String())
		return nil
	}

	return fmt.Errorf("valor flexivel invalido: %s", data)
}

func (f FlexibleString) String() string {
	return string(f)
}

func (c *Client) CEP(cep string) (CEP, error) {
	var out CEP
	err := c.get("cep/v1/"+onlyDigits(cep), &out)
	return out, err
}

func (c *Client) CNPJ(cnpj string) (CNPJ, error) {
	var out CNPJ
	err := c.get("cnpj/v1/"+onlyDigits(cnpj), &out)
	return out, err
}

func (c *Client) Banks() ([]Bank, error) {
	var out []Bank
	err := c.get("banks/v1", &out)
	return out, err
}

func (c *Client) BankByCode(code string) (Bank, error) {
	banks, err := c.Banks()
	if err != nil {
		return Bank{}, err
	}

	code = onlyDigits(code)
	for _, bank := range banks {
		if bank.CodeString() == code || bank.ISPB == code {
			return bank, nil
		}
	}
	return Bank{}, fmt.Errorf("banco %q nao encontrado", code)
}

func (c *Client) DDD(ddd string) (DDD, error) {
	var out DDD
	err := c.get("ddd/v1/"+onlyDigits(ddd), &out)
	return out, err
}

func (c *Client) Holidays(year string) ([]Holiday, error) {
	var out []Holiday
	err := c.get("feriados/v1/"+onlyDigits(year), &out)
	return out, err
}

func (c *Client) UFs() ([]UF, error) {
	var out []UF
	err := c.get("ibge/uf/v1", &out)
	return out, err
}

func (c *Client) Municipalities(uf string) ([]Municipio, error) {
	var out []Municipio
	err := c.get("ibge/municipios/v1/"+strings.ToUpper(uf), &out)
	return out, err
}

func (c *Client) get(path string, target any) error {
	baseURL := c.BaseURL
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}

	url := strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(path, "/")
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "br/0.1")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("BrasilAPI retornou HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

func onlyDigits(input string) string {
	var b strings.Builder
	for _, r := range input {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
