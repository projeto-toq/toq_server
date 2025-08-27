# 🧱 Inicialização e Injeção de Dependências em Arquitetura Hexagonal com Go

Este sistema em Go segue os princípios da **arquitetura hexagonal (Ports and Adapters)**, promovendo separação clara entre lógica de negócio e interfaces externas. A inicialização do sistema é composta por duas etapas principais: **setup dos servidores e infraestrutura básica**, e **configuração dos adapters e injeção de dependências**.

---

## 🔧 Inicialização dos Componentes Básicos

Antes da composição dos adapters, são inicializados os serviços fundamentais que sustentam o funcionamento da aplicação:

- **Logger**: sistema de log estruturado para rastreamento e diagnóstico  
- **Banco de Dados (DB)**: conexão e pool de acesso ao MySQL  
- **Telemetria**: coleta de métricas e rastreamento distribuído  
- **Cache**: sistema de cache (ex: Redis) para otimização de performance  

---

## 🔌 Configuração dos Adapters

A arquitetura hexagonal define dois tipos de adapters: **Left (entrada)** e **Right (saída)**, que se conectam aos respectivos ports da aplicação.

### ◀️ Left Adapters (Entrada)

Responsáveis por receber requisições externas e traduzi-las para chamadas à aplicação:

| Adapter      | Função        | Conecta ao                  |
|--------------|---------------|-----------------------------|
| `HTTP (Gin)` | Servidor web  | `port/left` → Handlers      |

### ▶️ Right Adapters (Saída)

Responsáveis por interagir com sistemas externos e fornecer dados ou serviços à aplicação:

| Adapter   | Função                    | Conecta ao                          |
|-----------|---------------------------|-------------------------------------|
| `MySQL`   | Persistência de dados     | `port/right` → Repositories         |
| `AWS S3`  | Armazenamento de arquivos | `port/right` → Storage              |
| `CEP`     | Consulta de endereço      | `port/right` → cep                  |
| `SMS`     | Envio de mensagens        | `port/right` → SMS                  |
| `CNPJ`    | Consulta de empresas      | `port/right` → cnpj                 |
| `CPF`     | Consulta de pessoas       | `port/right` → cpf                  |
| `Email`   | Envio de e-mails          | `port/right` → email                |
| `FCM`     | Push notifications        | `port/right` → fcm                  |

---

## 🧬 Injeção de Dependências

A composição dos componentes segue o fluxo de dependências definido pela arquitetura:

- **Services**: são instanciados recebendo os adapters de saída (repositories, storages, etc.) via seus respectivos ports. Cada service encapsula a lógica de negócio.  
- **Handlers**: são criados recebendo os services como dependência, atuando como ponte entre o mundo externo (HTTP) e a lógica da aplicação.  

- user service

Esse padrão garante **baixo acoplamento**, **alta testabilidade** e **flexibilidade para substituição de adapters** sem impactar a lógica central.