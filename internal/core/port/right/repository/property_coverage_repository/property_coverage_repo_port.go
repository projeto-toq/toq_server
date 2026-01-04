package propertycoveragerepository

import (
	"context"
	"database/sql"

	globalmodel "github.com/projeto-toq/toq_server/internal/core/model/global_model"
	propertycoveragemodel "github.com/projeto-toq/toq_server/internal/core/model/property_coverage_model"
)

// RepositoryInterface expõe o contrato de persistência para consultas e gestão de cobertura de imóveis.
//
// Regras de contrato (alinhadas ao guia Seções 2.1, 7.3 e 8.5):
//   - Todos os métodos retornam erros “puros” (sql.ErrNoRows quando não há resultado / 0 rows afetadas).
//   - Não há abertura/commit de transação: services fornecem tx (pode ser nil para leituras).
//   - Nenhum mapeamento para HTTP/domínio é feito aqui; apenas acesso a dados e conversão entity→domínio.
//   - InstrumentedAdapter deve ser usado em todas as operações de banco para métricas/tracing.
type RepositoryInterface interface {
	// GetVerticalCoverage retorna cobertura vertical por CEP e número; sql.ErrNoRows se não houver match.
	GetVerticalCoverage(ctx context.Context, tx *sql.Tx, zipCode, number string) (propertycoveragemodel.CoverageInterface, error)
	// GetHorizontalCoverage retorna cobertura horizontal por CEP; sql.ErrNoRows se não houver match.
	GetHorizontalCoverage(ctx context.Context, tx *sql.Tx, zipCode string) (propertycoveragemodel.CoverageInterface, error)
	// GetNoComplexCoverage retorna cobertura standalone para o CEP; sql.ErrNoRows se inexistente.
	GetNoComplexCoverage(ctx context.Context, tx *sql.Tx, zipCode string) (propertycoveragemodel.CoverageInterface, error)
	// GetVerticalComplexByZipNumber retorna complexo vertical (com torres/tamanhos) por CEP+numero; sql.ErrNoRows se não achar.
	GetVerticalComplexByZipNumber(ctx context.Context, tx *sql.Tx, zipCode, number string) (propertycoveragemodel.ManagedComplexInterface, error)
	// GetHorizontalComplexByZip retorna complexo horizontal por CEP (principal ou associado); sql.ErrNoRows se não achar.
	GetHorizontalComplexByZip(ctx context.Context, tx *sql.Tx, zipCode string) (propertycoveragemodel.ManagedComplexInterface, error)

	// ListManagedComplexes lista complexos (vertical/horizontal/standalone) com filtros; slice vazio quando não houver registros.
	ListManagedComplexes(ctx context.Context, tx *sql.Tx, params ListManagedComplexesParams) ([]propertycoveragemodel.ManagedComplexInterface, error)
	// GetManagedComplex retorna complexo pelo id/kind; sql.ErrNoRows se não existir.
	GetManagedComplex(ctx context.Context, tx *sql.Tx, id int64, kind propertycoveragemodel.CoverageKind) (propertycoveragemodel.ManagedComplexInterface, error)
	// CreateManagedComplex insere um complexo conforme kind; retorna ID criado.
	CreateManagedComplex(ctx context.Context, tx *sql.Tx, entity propertycoveragemodel.ManagedComplexInterface) (int64, error)
	// UpdateManagedComplex atualiza um complexo; sql.ErrNoRows se não encontrado.
	UpdateManagedComplex(ctx context.Context, tx *sql.Tx, entity propertycoveragemodel.ManagedComplexInterface) (int64, error)
	// DeleteManagedComplex remove complexo pelo id/kind; sql.ErrNoRows se não encontrado.
	DeleteManagedComplex(ctx context.Context, tx *sql.Tx, id int64, kind propertycoveragemodel.CoverageKind) (int64, error)

	// CreateVerticalComplexTower cria torre; retorna ID.
	CreateVerticalComplexTower(ctx context.Context, tx *sql.Tx, tower propertycoveragemodel.VerticalComplexTowerInterface) (int64, error)
	// UpdateVerticalComplexTower atualiza torre; sql.ErrNoRows se não encontrada.
	UpdateVerticalComplexTower(ctx context.Context, tx *sql.Tx, tower propertycoveragemodel.VerticalComplexTowerInterface) (int64, error)
	// DeleteVerticalComplexTower remove torre; sql.ErrNoRows se não encontrada.
	DeleteVerticalComplexTower(ctx context.Context, tx *sql.Tx, id int64) (int64, error)
	// GetVerticalComplexTower busca torre por id; sql.ErrNoRows se não encontrada.
	GetVerticalComplexTower(ctx context.Context, tx *sql.Tx, id int64) (propertycoveragemodel.VerticalComplexTowerInterface, error)
	// ListVerticalComplexTowers lista torres; slice vazio quando não houver registros.
	ListVerticalComplexTowers(ctx context.Context, tx *sql.Tx, params ListVerticalComplexTowersParams) ([]propertycoveragemodel.VerticalComplexTowerInterface, error)

	// CreateVerticalComplexSize cria tamanho; retorna ID.
	CreateVerticalComplexSize(ctx context.Context, tx *sql.Tx, size propertycoveragemodel.VerticalComplexSizeInterface) (int64, error)
	// UpdateVerticalComplexSize atualiza tamanho; sql.ErrNoRows se não encontrado.
	UpdateVerticalComplexSize(ctx context.Context, tx *sql.Tx, size propertycoveragemodel.VerticalComplexSizeInterface) (int64, error)
	// DeleteVerticalComplexSize remove tamanho; sql.ErrNoRows se não encontrado.
	DeleteVerticalComplexSize(ctx context.Context, tx *sql.Tx, id int64) (int64, error)
	// GetVerticalComplexSize busca tamanho por id; sql.ErrNoRows se não encontrado.
	GetVerticalComplexSize(ctx context.Context, tx *sql.Tx, id int64) (propertycoveragemodel.VerticalComplexSizeInterface, error)
	// ListVerticalComplexSizes lista tamanhos; slice vazio quando não houver registros.
	ListVerticalComplexSizes(ctx context.Context, tx *sql.Tx, params ListVerticalComplexSizesParams) ([]propertycoveragemodel.VerticalComplexSizeInterface, error)

	// CreateHorizontalComplexZipCode cria CEP associado a complexo horizontal; retorna ID.
	CreateHorizontalComplexZipCode(ctx context.Context, tx *sql.Tx, zip propertycoveragemodel.HorizontalComplexZipCodeInterface) (int64, error)
	// UpdateHorizontalComplexZipCode atualiza CEP; sql.ErrNoRows se não encontrado.
	UpdateHorizontalComplexZipCode(ctx context.Context, tx *sql.Tx, zip propertycoveragemodel.HorizontalComplexZipCodeInterface) (int64, error)
	// DeleteHorizontalComplexZipCode remove CEP; sql.ErrNoRows se não encontrado.
	DeleteHorizontalComplexZipCode(ctx context.Context, tx *sql.Tx, id int64) (int64, error)
	// GetHorizontalComplexZipCode busca CEP por id; sql.ErrNoRows se não encontrado.
	GetHorizontalComplexZipCode(ctx context.Context, tx *sql.Tx, id int64) (propertycoveragemodel.HorizontalComplexZipCodeInterface, error)
	// ListHorizontalComplexZipCodes lista CEPs; slice vazio quando não houver registros.
	ListHorizontalComplexZipCodes(ctx context.Context, tx *sql.Tx, params ListHorizontalComplexZipCodesParams) ([]propertycoveragemodel.HorizontalComplexZipCodeInterface, error)
}

// ListManagedComplexesParams configura filtros/paginação para listagem admin de complexos.
// Campos string são filtrados com LIKE (%value%) quando não vazios; Kind/PropertyType/Sector aplicam igualdade.
type ListManagedComplexesParams struct {
	Name         string
	ZipCode      string
	Number       string
	City         string
	State        string
	Sector       *propertycoveragemodel.Sector
	PropertyType *globalmodel.PropertyType
	Kind         *propertycoveragemodel.CoverageKind
	Limit        int
	Offset       int
}

// ListVerticalComplexTowersParams filtra torres por complexo, nome e paginação.
type ListVerticalComplexTowersParams struct {
	VerticalComplexID int64
	Tower             string
	Limit             int
	Offset            int
}

// ListVerticalComplexSizesParams filtra tamanhos por complexo e paginação.
type ListVerticalComplexSizesParams struct {
	VerticalComplexID int64
	Limit             int
	Offset            int
}

// ListHorizontalComplexZipCodesParams filtra CEPs associados a complexo horizontal.
type ListHorizontalComplexZipCodesParams struct {
	HorizontalComplexID int64
	ZipCode             string
	Limit               int
	Offset              int
}
