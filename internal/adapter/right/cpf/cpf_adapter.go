package cpfadapter

import (
	"fmt"
	"net/http"
	"time"

	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
)

type CPFAdapter struct {
	Client   *http.Client
	Token    string
	URLBase  string
	Status   bool   `json:"status"`
	Return   string `json:"return"`
	Consumed int    `json:"consumed"`
	Result   struct {
		NumeroDeCpf            string `json:"numero_de_cpf"`
		NomeDaPf               string `json:"nome_da_pf"`
		DataNascimento         string `json:"data_nascimento"`
		SituacaoCadastral      string `json:"situacao_cadastral"`
		DataInscricao          string `json:"data_inscricao"`
		DigitoVerificador      string `json:"digito_verificador"`
		ComprovanteEmitido     string `json:"comprovante_emitido"`
		ComprovanteEmitidoData string `json:"comprovante_emitido_data"`
	} `json:"result"`
}

func NewCPFAdapter(env *globalmodel.Environment) (*CPFAdapter, error) {
	if env.CPF.Token == "" {
		return nil, fmt.Errorf("CPF token is required")
	}
	if env.CPF.URLBase == "" {
		return nil, fmt.Errorf("CPF URL base is required")
	}

	client := &http.Client{
		Timeout: 30 * time.Second, // Reduced from 600s for better UX
	}
	return &CPFAdapter{
		Client:  client,
		Token:   env.CPF.Token,
		URLBase: env.CPF.URLBase,
	}, nil
}
