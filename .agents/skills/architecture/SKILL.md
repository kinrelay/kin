# Architecture Skill

Use this skill when shaping or reviewing domain/application architecture.

## Objective

Preserve Kin's domain-first design using DDD, Clean Architecture, Hexagonal Architecture, CQRS, and Domain Events without turning the codebase into ceremony-heavy architecture for its own sake.

## Discovery checklist

Before implementation, identify:

- Ubiquitous language used by the product and issue.
- Actor(s) and the interaction being enabled.
- Domain responsibility: which domain owns the rule?
- Entities and Value Objects involved.
- Aggregate / consistency boundary, if one is actually needed.
- Business invariants and invalid states.
- Domain Events representing meaningful facts.
- Commands that express write intent.
- Queries that satisfy read intent.
- Ports required by the application layer.
- Which concerns are adapters/infrastructure and therefore must stay outside the domain.

Do not invent an Aggregate, Domain Service, or Domain Event unless it protects or communicates meaningful domain behavior.

## Dependency rule

Allowed direction:

`delivery/adapters -> application -> domain`

The domain must not import:

- SQL/database packages
- HTTP/router packages
- queue/event-broker clients
- logging/telemetry vendors
- external-provider SDKs
- OpenAI/LLM-specific types
- persistence DTOs
- API DTOs

Application code may depend on domain types and port interfaces, but not concrete adapters.

## Domain modeling

### Entities
Use an Entity when identity and lifecycle matter.

### Value Objects
Use Value Objects for domain concepts defined by value and invariants. Construct them through validated constructors where invalid values are possible.

### Aggregates
Use Aggregates only to enforce a true transactional consistency boundary. Do not create large aggregates merely to make navigation convenient.

### Domain Services
Use a Domain Service when a domain rule does not naturally belong to a single Entity/Value Object and remains free of infrastructure concerns.

### Domain Events
Use past-tense facts such as `FriendshipCreated`, `ActivityIngested`, or `PrivacyPolicyChanged` when downstream behavior benefits from explicit domain communication.

Domain Events do not imply Event Sourcing or a distributed broker.

## Application layer

Application code orchestrates interactions. It may:

- load aggregates through ports
- invoke domain behavior
- persist changed state through ports
- publish domain events through ports
- coordinate multiple domains explicitly

Application code should not contain core business invariants that belong to the domain.

## CQRS

### Write model

- Commands express intent.
- Aggregates protect invariants.
- Persistence follows domain decisions.
- UI/query convenience must not shape write aggregates.

### Read model

- Queries never change domain state.
- Read models may be denormalized.
- Query handlers may use dedicated read ports and projection-specific SQL later.
- A query does not need to reconstruct a write-side aggregate when a read DTO is sufficient.

## Ports and Adapters

Define ports in terms of required capability, not vendor mechanics.

Good:

```go
type ContextGenerator interface {
    Generate(ctx context.Context, input ContextInput) (GeneratedContext, error)
}
```

Avoid vendor-shaped ports such as interfaces that expose model IDs, prompt payloads, SDK request types, or SQL-specific concepts unless those concepts are genuinely part of the domain/application contract.

Adapters translate between external representations and inner contracts.

## Model separation

Keep these distinct:

- Domain model: business meaning and invariants.
- Persistence model: storage representation.
- API / delivery DTO: external transport representation.
- Provider DTO: third-party SDK/API representation.

Mapping code is intentional architectural isolation, not duplication to eliminate reflexively.

## Modular-monolith rules

- Prefer domain/capability-oriented modules over global folders such as `models/`, `services/`, and `repositories/` containing unrelated domains.
- A module owns its language and behavior.
- Avoid arbitrary cross-module imports.
- Cross-module interactions should be explicit through application orchestration, exported contracts, or events.
- Do not split a module into a microservice solely because it has a clear boundary.

## Architecture review questions

Before completion, ask:

1. What business rule changed, and where is it enforced?
2. Could the domain/application code run without Postgres, HTTP, queues, or an LLM provider?
3. Did database/API/provider concerns leak inward?
4. Did read-side needs distort the write model?
5. Is any cross-domain coupling implicit or circular?
6. Did we introduce abstractions that have no current use case?
7. Can a future adapter be replaced without changing domain behavior?
