Atue como um desenvolvedor GO Senior e faça toda a interação em português.

---

## 1) Objetivo do Pedido
- Tipo: Somente análise e apresentaçao do plano para aprovação (sem implementação).

## 2) Requisição
Continuando com a criação do toq_server, precisamos

- Adicionar as rotas de CRUD para as entidades complex e suas relacionadas: complex_towers, complex_sizes, complex_zipcodes.
  - a tabela complex contem os dados dos empreendimentos de todos os tipos cobertos pela toq. Através do zipcode e number podemos localizar um empreendimento e pelo type saber quais os propertyTypes são possíveis neste empreendimento.
  - a tabela complex_towers - contem os dados das torres dentro de cada empreendimento. Cada registro está associada a um empreendimento específico através de uma chave estrangeira.
  - a tabela complex_sizes - contem os tamanhos disponíveis para cada empreendimento. Cada registro está associado a um empreendimento específico através de uma chave estrangeira.
  - a tabela complex_zipcodes - contem os zipcodes associados a cada empreendimento, para o caso de condomínios horizontais onde diversos zipcodes podem estar presentes para cada empreendimento. Cada registro está associado a um empreendimento específico através de uma chave estrangeira.
- O seed da tabela complex_sizes ja existe, mas o campo description está vazio. Atualize o seed para preencher este campo com descrições apropriadas para cada size. Seja criativo para gerar descrições realistas.
- o acesso aos endpoints de CRUD deve ser permitido apenas para perfis de administradores. Portanto, altere permissions.csv e role_permissions.csv conforme necessário.
- Deverá existir, adicionalmente, uma rota GET HTTP para listar os tamanhos disponíveis para um determinado empreendimento baseado no zipCode e number fornecido no body.
- A tabela complexes, pode localizar o empreendimento baseado no zipCode e number.
- a resposta deve ser uma lista de objetos JSON com os campos id (int), size (float) e description (string). NotFound se não for localizado o empreendimento. Se o empreendimento existir mas a busca em `complex_sizes` retornar `not_found`, trate respondendo `description` como string vazia.
- Todos os perfis devem ter acesso ao endpoint, portanto altere permissions.csv e role_permissions.csv conforme necessário.
- Os endpoints devem ser documentados no Swagger com comentários no código.

- Documentação de referência: `docs/toq_server_go_guide.md`


## 3) Requisitos
  - Observabilidade alinhada (logs/traces/metrics existentes)
  - Performance/concorrência (se aplicável)
  - Sem back compatibility/downtime (ambiente de desenvolvimento)

## 4) Artefatos a atualizar (se aplicável)
- Código (handlers/services/adapters)
- Swagger (anotações no código)
- Documentação (README/docs/*)
- Observabilidade (métricas/dashboards) — apenas quando estritamente pertinente
- Não crie/altere arquivos de testes

## 5) Arquitetura e Fluxo (resumo)
- Siga o guia: `docs/toq_server_go_guide.md` (Seções 1–4 e 7–11).
- Em alto nível: `Handlers` → `Services` → `Repositories`; DI por factories; converters nos repositórios; transações via serviço padrão.

## 6) Erros e Observabilidade (resumo)
- Siga o guia: `docs/toq_server_go_guide.md` (Seções 5, 6, 8, 9 e 10).
- Pontos-chave: spans só fora de handlers HTTP; `SetSpanError` em falhas de infra; handlers usam `http_errors.SendHTTPErrorObj`; adapters retornam erros puros; logs com severidade correta.

## 7) Dados e Transações
- Tabelas/entidades afetadas: <listar>
- Necessidade de transação: <sim/não; justificar>
- Regras de concorrência e idempotência: <se aplicável>

## 8) Interfaces/Contratos
- Endpoints HTTP envolvidos: <listar caminhos e métodos>
- Payloads de request/response: <resumo>
- Erros de domínio esperados e mapeamento HTTP: <ex.: ErrX → 409>

## 9) Entregáveis Esperados do Agente
- Análise detalhada com checklist dos requisitos e plano por etapas (com ordem de execução).
- Quality gates rápidos no final: Build, Lint/Typecheck e mapeamento Requisito → Status.
- Não execute git status, git diff nem go test.

## 10) Restrições e Assunções
- Ambiente: desenvolvimento (sem back compatibility, sem janela de manutenção, sem migração).
- Sem uso de mocks nem soluções temporárias em entregas finais.

## 11) Anexos e Referências
- Arquivos relevantes: `docs/toq_server_go_guide.md`.

---

### Regras Operacionais do Pedido (a serem cumpridas pelo agente)

1) Antes de propor/alterar código:
- Revisar arquivos relevantes (código e configs); evitar suposições sem checagem.
- Produzir um checklist de requisitos explícitos e implícitos, mantendo-o visível.

2) Durante a análise/implementação:
- Seguir o guia `docs/toq_server_go_guide.md` (arquitetura, erros/observabilidade, transações).
- Manter adapters com erros “puros”; sem HTTP/semântica de domínio nessa camada.
- A documentação Swagger/docs deve ser criada por comentários, em inglês em DTO/Handler e execução de make swagger. Sem alterações manuais no swagger.yaml/json.

3) Após mudanças:
- Relatar “requirements coverage” (Requisito → Done/Deferred + motivo).

4) Em casos de dúvidas: consultar o requerente.

5) Caso a tarefa seja grande/demorada, dividir em fases menores e entregáveis curtos.

6) Caso haja remoção de arquivos limpe o conteúdo e informa a lista para remoção manual.

7) Comentários internos em português; docstrings de funções em inglês; Swagger por anotações no código.