# API Microservice Payment

Microserviço responsável pelo processamento de pagamentos, integrando-se com o Mercado Pago para geração de QR Codes e gerenciamento do ciclo de vida de transações financeiras.

---

## 📦 Subindo o projeto localmente com Docker

### Pré-requisitos

- [Docker](https://www.docker.com/) e [Docker Compose](https://docs.docker.com/compose/) instalados
- Arquivo `.env` configurado na raiz do projeto (veja `.env.example` como referência)

### Passos

1. **Clone o repositório:**

   ```bash
   git clone <url-do-repositorio>
   cd api-microservice-payment
   ```

2. **Configure as variáveis de ambiente:**

   Crie um arquivo `.env` na raiz do projeto com as variáveis necessárias (credenciais do banco, Mercado Pago, AWS/SQS, etc.).

3. **Suba os containers:**

   ```bash
   docker-compose up --build
   ```

   Isso irá iniciar:
   - **PostgreSQL** na porta `5432`
   - **ElasticMQ** (emulador SQS local) nas portas `9324` (API) e `9325` (UI)
   - **API** na porta `8080`

4. **Acesse a API:**

   - Swagger UI: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

5. **Parar os containers:**

   ```bash
   docker-compose down
   ```

---

## 🛠️ Comandos do Makefile

| Comando                | Descrição                                                                 |
|------------------------|---------------------------------------------------------------------------|
| `make swagger`         | Gera/atualiza a documentação Swagger em `internal/common/infra/api/swagger` |
| `make test`            | Executa todos os testes unitários do projeto                              |
| `make test-coverage`   | Executa testes com cobertura e exibe o percentual total                   |
| `make test-coverage-html` | Gera relatório HTML de cobertura (`coverage.html`)                     |

---

## 🚀 Tecnologias Utilizadas

| Tecnologia       | Motivo da Escolha                                                                                                                                       |
|------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------|
| **Go (Golang)**  | Linguagem de alta performance, compilada, com excelente suporte a concorrência. Ideal para microserviços que exigem baixa latência e alto throughput.   |
| **PostgreSQL**   | Banco de dados relacional robusto, open-source, com suporte a transações ACID. Garante consistência e integridade dos dados de pagamento.               |
| **Mercado Pago API** | Gateway de pagamentos amplamente utilizado no Brasil, oferecendo integração via QR Code dinâmico para pagamentos Pix, facilitando a experiência do usuário. |
| **Amazon SQS (ElasticMQ local)** | Serviço de filas gerenciado que desacopla a comunicação entre microserviços, garantindo resiliência e processamento assíncrono de eventos (criação de pagamento, rollback, etc.). |
| **Gin**          | Framework HTTP minimalista e performático para Go, facilitando a criação de APIs REST com middlewares e roteamento eficiente.                           |
| **GORM**         | ORM para Go que simplifica operações com banco de dados, oferecendo migrations automáticas e abstração de queries.                                       |
| **Swagger (swag)** | Geração automática de documentação OpenAPI a partir de anotações no código, facilitando testes e integração por consumidores da API.                   |
| **Godog (Cucumber)** | Framework BDD para Go, permitindo testes de aceitação escritos em linguagem Gherkin, aproximando especificações de negócio do código.                 |

---

## 🏛️ Arquitetura Hexagonal

Este projeto segue os princípios da **Arquitetura Hexagonal** (também conhecida como *Ports and Adapters*), proposta por Alistair Cockburn.

### Estrutura de Diretórios

```
internal/
├── adapters/
│   ├── driven/          # Adapters de saída (repositórios, gateways externos, serviços de fila)
│   │   ├── mercado_pago/
│   │   ├── repositories/
│   │   └── kitchen_order_service/
│   └── driver/          # Adapters de entrada (handlers HTTP, consumers de fila)
│       ├── api/
│       └── queue/
├── core/
│   ├── domain/
│   │   ├── entities/    # Entidades de domínio
│   │   ├── exceptions/  # Exceções de negócio
│   │   ├── ports/       # Interfaces (portas) que definem contratos
│   │   └── value_objects/
│   ├── dto/             # Objetos de transferência de dados
│   ├── factory/         # Factories para injeção de dependências
│   └── use_cases/       # Casos de uso (regras de negócio)
└── common/
    ├── config/          # Configurações e constantes
    ├── infra/           # Infraestrutura (servidor HTTP, conexão DB, filas)
    ├── mocks/           # Mocks para testes
    └── pkg/             # Utilitários compartilhados
```

### Por que Arquitetura Hexagonal?

| Benefício                        | Descrição                                                                                                                                 |
|----------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------|
| **Desacoplamento**               | O núcleo de negócio (`core`) não depende de frameworks, bancos de dados ou serviços externos. Toda comunicação ocorre via interfaces (ports). |
| **Testabilidade**                | Facilita a criação de testes unitários e de integração, pois dependências externas podem ser substituídas por mocks/stubs.                |
| **Flexibilidade tecnológica**    | Permite trocar implementações (ex.: mudar de PostgreSQL para MySQL, ou SQS para RabbitMQ) sem alterar as regras de negócio.               |
| **Manutenibilidade**             | Código organizado em camadas bem definidas facilita a evolução e onboarding de novos desenvolvedores.                                     |
| **Alinhamento com DDD**          | Combina naturalmente com Domain-Driven Design, mantendo o foco no domínio e nas regras de negócio.                                        |
| **Diversidade de portas**          | O projeto lida com diversos modelo de entra e saída de dados: HTTP, Banco de dados, Filas e Payment Providers e futuramente terá suporte a GRPC.                                        |


### Fluxo de uma Requisição

```
[Cliente HTTP] 
      │
      ▼
┌─────────────────┐
│  Driver Adapter │  (handlers/api)
│  (Porta de Entrada)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│    Use Case     │  (regras de negócio)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│  Driven Adapter │  (repository, gateway, queue)
│  (Porta de Saída)
└────────┬────────┘
         │
         ▼
   [PostgreSQL / Mercado Pago / SQS]
```

---
