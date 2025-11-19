# Property Coverage Schema

Este documento descreve o conjunto de tabelas necessárias para suportar o serviço de cobertura de imóveis (`property_coverage`). Ele reflete as decisões mais recentes acordadas com o time de produto/DBA e elimina colunas não utilizadas (`created_at`, `updated_at`) bem como a antiga tabela intermediária `vertical_complex_zip_codes`.

## Regras Gerais

- Todas as tabelas seguem o padrão `InnoDB` e **não** possuem colunas `created_at`/`updated_at`, conforme a diretriz global do projeto.
- O CEP (`zip_code`) é armazenado como `CHAR(8)` e sempre normalizado sem separadores.
- As seeds existentes em `data/*.csv` podem ser importadas via `LOAD DATA` usando exatamente esta estrutura.

## Tabelas

### 1. `vertical_complexes`
| Coluna | Tipo | Notas |
| --- | --- | --- |
| `id` | INT UNSIGNED PK AUTO_INCREMENT | Identificador único. |
| `name` | VARCHAR(255) NOT NULL | Nome comercial do empreendimento. |
| `zip_code` | CHAR(8) NOT NULL | CEP único por linha (um vertical pode ter várias linhas para números distintos). |
| `street` | VARCHAR(255) NOT NULL | Logradouro normalizado. |
| `number` | VARCHAR(30) NOT NULL | Número ou intervalo textual (ex.: `149;189`). |
| `neighborhood`, `city` | VARCHAR(150) NOT NULL | Localização. |
| `state` | CHAR(2) NOT NULL | UF. |
| `reception_phone` | VARCHAR(30) NULL | Opcional. |
| `sector` | TINYINT UNSIGNED NOT NULL | Segmentação interna. |
| `main_registration` | VARCHAR(35) NULL | Matrícula/lote exatamente como recebido (permite múltiplos registros por linha). |
| `type` | INT UNSIGNED NOT NULL | Bitmask de tipos aceitos (consumido pelo serviço). |

Índices relevantes:
- `uk_vertical_zip_number (zip_code, number)` garante busca determinística pelo par CEP+número.
- `idx_vertical_zip (zip_code)` agiliza auditorias e cargas.

### 2. `vertical_complex_towers`
| Coluna | Tipo | Notas |
| --- | --- | --- |
| `id` | INT UNSIGNED PK |
| `vertical_complex_id` | INT UNSIGNED FK -> `vertical_complexes(id)` |
| `tower` | VARCHAR(120) NOT NULL |
| `floors`, `total_units`, `units_per_floor` | SMALLINT UNSIGNED NOT NULL DEFAULT 0 |

Usada para detalhar torres específicas. FK com `ON DELETE CASCADE` para simplificar reimportações.

### 3. `vertical_complex_sizes`
| Coluna | Tipo | Notas |
| --- | --- | --- |
| `id` | INT UNSIGNED PK |
| `vertical_complex_id` | INT UNSIGNED FK |
| `size` | DECIMAL(8,2) NOT NULL | Metragem da unidade. |
| `description` | VARCHAR(255) NULL | Copy amigável (ex.: "Studio Compacto"). |

> **Novo**: tabela adicionada para armazenar os tamanhos das unidades (itens presentes em `data/vertical_complex_sizes.csv`).

### 4. `horizontal_complexes`
Mesma estrutura conceitual de `vertical_complexes`. Índice principal em `zip_code` para inspeções; a cobertura é resolvida via tabela auxiliar de CEPs.

### 5. `horizontal_complex_zip_codes`
| Coluna | Tipo | Notas |
| --- | --- | --- |
| `id` | INT UNSIGNED PK |
| `horizontal_complex_id` | INT UNSIGNED FK |
| `zip_code` | CHAR(8) NOT NULL |

- Índice único `uk_horizontal_zip` garante um CEP por condomínio horizontal.
- Esta é a única tabela auxiliar de CEP; **não existe** mais `vertical_complex_zip_codes` porque o CEP já está em `vertical_complexes`.

### 6. `no_complex_zip_codes`
| Coluna | Tipo | Notas |
| --- | --- | --- |
| `zip_code` | CHAR(8) PK |
| `neighborhood`, `city` | VARCHAR(150) NOT NULL |
| `state` | CHAR(2) NOT NULL |
| `sector` | TINYINT UNSIGNED NOT NULL |
| `type` | INT UNSIGNED NOT NULL |

Tabela usada como fallback quando não há condomínio associado. Nenhuma coluna de auditoria é necessária.

## DDL de Referência

O arquivo `scripts/db_creation.sql` contém as instruções `DROP TABLE IF EXISTS` + `CREATE TABLE` para todas as estruturas acima. Ele já aplica as constraints, índices e elimina os campos proibidos.

## Carga de Dados (Seeds)

Para importar os CSVs disponibilizados:
1. Copie os arquivos para `/var/lib/mysql-files` (ou path equivalente com permissões `LOAD DATA`).
2. Utilize instruções `LOAD DATA INFILE` semelhantes às existentes para outras seeds. Manter o mesmo delimitador `;` e `IGNORE 1 LINES`.
3. A ordem recomendada é:
   1. `vertical_complexes`
   2. `vertical_complex_towers`
   3. `vertical_complex_sizes`
   4. `horizontal_complexes`
   5. `horizontal_complex_zip_codes`
   6. `no_complex_zip_codes`

## Impacto no Código

- O serviço `ResolvePropertyTypes` normaliza o CEP (mantendo somente dígitos) e sanitiza o número (removendo espaços internos e forçando maiúsculas) antes de abrir a transação, evitando consultas inconsistentes.
- O adapter MySQL (`internal/adapter/right/mysql/property_coverage`) agora trata números armazenados como listas delimitadas por `;`, utilizando `FIND_IN_SET` para encontrar correspondências sem depender da antiga `vertical_complex_zip_codes`.
- A busca vertical só ocorre quando um número válido é informado; do contrário, o fluxo segue diretamente para os cenários horizontal/isolado.

Qualquer ajuste futuro deve seguir essas regras de normalização e continuar evitando colunas de auditoria desnecessárias.

    | `main_registration` | VARCHAR(35) NULL | Matrícula/lote exatamente como recebido (permite múltiplos registros por linha). |
