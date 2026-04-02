---
name: analyzing-system-design
description: Use before any implementation phase to analyse atomicity, idempotency, consistency, concurrency, resilience, architectural patterns, and database fit. Mandatory for all Standard and Complex tasks. Produces system-design-analysis.md that the Coder must read before writing any code.
---

# System Design Advisor

Analyses every phase of an implementation plan through the lens of distributed systems
correctness, data integrity, and production reliability. Never writes code. Produces a
structured analysis, presents trade-offs and concrete proposals — including optional
Mermaid diagrams — and waits for the developer to approve or reject each proposal before
the Coder starts.

> This skill is **mandatory** for all Standard and Complex tasks. It runs **after**
> Researcher and **before** Coder. Its output (`system-design-analysis.md`) must be read
> by the Coder at the start of every phase.

---

## When to Use

| Condition | Action |
|---|---|
| Standard or Complex task (new feature, new domain, cross-service change) | Mandatory |
| Simple task (single-file config/typo fix) | Skip |
| Resuming a plan that already has `system-design-analysis.md` | Skip unless the plan changed significantly |

---

## Steps

1. Read `~/ai-plans/{repo-name}/{slug}/implementation-plan*.md`
2. Read `~/ai-plans/{repo-name}/{slug}/research.md` for codebase context
3. Run all seven lenses for each phase (see Analysis Protocol below)
4. For each finding that involves data flow or state transitions, offer a Mermaid diagram
   (see Diagram Protocol below)
5. Write `~/ai-plans/{repo-name}/{slug}/system-design-analysis.md`
6. Present the summary to the developer and wait for approval
7. Record decisions in `system-design-analysis.md` under `## Approved Decisions`

---

## Analysis Protocol

For each phase in the implementation plan, run all seven lenses.
If a lens has nothing to flag, write `No concerns identified.` — never skip a section.

---

### Lens 1 — Atomicity

**Question**: Which operations must succeed or fail together (all-or-nothing)?

Check:
- Identify every write path (DB inserts, updates, external calls, event publications)
- Determine whether writes that must be atomic share a single transaction boundary
- Flag dual-write risks: writing to two stores without a coordination mechanism
  (e.g., write to DB + publish event without an Outbox Pattern)
- Flag external calls (HTTP, gRPC, message brokers) inside DB transactions — external
  calls are slow and unreliable; they must never hold a DB transaction open

**Output format**:
```
ATOMICITY RISK — {operation name}
Scope: {files/functions involved}
Risk: {what inconsistency occurs on partial failure}
Proposal: {DB transaction / Outbox Pattern / saga / compensating transaction}
```

---

### Lens 2 — Idempotency

**Question**: Can this operation be safely retried without unintended side effects?

Check:
- Every consumer handler (queue, event, webhook, gRPC retry): is there a deduplication guard?
- Every state-changing external call: is there an idempotency key or check-before-execute guard?
- Retry logic: does retrying produce a different final state than a single execution?

**Output format**:
```
IDEMPOTENCY GAP — {operation name}
Scope: {files/functions involved}
Risk: {what happens on duplicate or retry execution}
Proposal: {idempotency key / deduplication table / check-before-execute / UNIQUE constraint}
```

---

### Lens 3 — Consistency

**Question**: Are all data stores in a consistent state after every possible execution path,
including partial failures?

Check:
- Identify all stores involved in a single logical operation (e.g., DB + cache + queue)
- When do they diverge? Is that divergence acceptable (eventual consistency) or not?
- Does retry logic produce a consistent final state or can it leave data in a partial state?
- Are domain invariants validated and enforced before any write reaches the store?

**Output format**:
```
CONSISTENCY CONCERN — {area}
Scope: {stores/services involved}
Risk: {divergence scenario}
Consistency model required: {strong / eventual — and why}
Proposal: {reconciliation strategy / strict guard / accepted divergence with rationale}
```

---

### Lens 4 — Concurrency

**Question**: What happens when two goroutines or two service instances execute this
operation simultaneously?

For each identified race condition or deadlock risk, evaluate the concurrency controls
below and recommend the most appropriate fit for the project's scale and complexity:

| Pattern | When to use |
|---|---|
| DB transaction + `SELECT FOR UPDATE` | Single-instance serialisation, low contention, correctness critical |
| Optimistic lock (`version` / `updated_at` column) | High read, low write contention |
| DB-level `UNIQUE` constraint | Last line of defence against duplicate inserts |
| Distributed lock (Redis `SET NX PX`) + fencing token | Cross-instance serialisation — multiple replicas competing for the same resource |
| Channel semaphore (Go) | In-process goroutine fan-out control |
| `sync.Mutex` / `sync.RWMutex` | Shared in-memory state within a single process |

> **Fencing token reminder**: a distributed lock alone is insufficient when the lock can
> expire while a process is still running. Always pair Redis locks with a monotonically
> increasing fencing token persisted in the DB. The DB rejects writes with a stale token.
> `SET lock NX PX ttl` + store `fence_token = INCR fence_counter` → DB `WHERE fence_token < new_token`.

Deadlock checklist:
- Are locks always acquired in the same order across all code paths?
- Can a long-running operation hold a lock while waiting for an external call?
- Is there a lock TTL or timeout to prevent indefinite blocking?

**Output format**:
```
CONCURRENCY RISK — {operation name}
Scope: {functions/tables involved}
Race condition / deadlock: {describe the scenario}
Proposed solution: {pattern from table above}
Fencing token required: {yes / no — reason}
Trade-off: {performance cost vs correctness guarantee}
Recommended: {which solution fits this project's scale and complexity}
```

---

### Lens 5 — Resilience, Fault Tolerance, and Scalability

Apply the checklist below. Flag only what is relevant to the phase.

**Resilience**
- Are retries bounded (max attempts + backoff) or potentially infinite?
- Is there a dead-letter mechanism for unprocessable messages or failed jobs?
- Does the system self-heal on transient failures, or does it require manual intervention?

**Fault Tolerance**
- What is the blast radius if this component fails completely?
- Are failures surfaced via logs, metrics, or traces — or silent?
- Is there a circuit breaker for calls to unreliable external dependencies?

**Scalability**
- Are there DB queries on high-cardinality columns without an index?
- Is any in-memory state shared across instances (breaks horizontal scaling)?
- Are consumers designed for single-instance or competing-consumers (scale-out)?
- Is there a hot-key or bottleneck that becomes a single point of contention at load?

**Complexity vs Robustness Trade-off**
Rate the necessity of each mechanism for this project's current stage:
- `ESSENTIAL`: system correctness breaks without it
- `RECOMMENDED`: improves resilience, low implementation cost
- `DEFER`: valid concern, but out of scope for current stage

**Output format**:
```
SYSTEM DESIGN NOTE — {concern}
Category: Resilience / Fault Tolerance / Scalability
Necessity: ESSENTIAL / RECOMMENDED / DEFER
Scope: {component affected}
Finding: {what the issue is}
Proposal: {concrete solution}
```

---

### Lens 6 — Architectural Patterns (EDA / CQRS / SAGA)

**Question**: Does the current design introduce problems that established distributed
systems patterns would solve?

Do not recommend patterns for their own sake. Only flag a pattern when there is a
concrete problem it solves that the current design cannot solve cleanly.

---

#### Event-Driven Architecture (EDA)

Applicable when:
- Services are tightly coupled via synchronous calls where decoupling would improve resilience
- A state change in one service must trigger reactions in multiple others
- The producing service should not know about its consumers
- Temporal decoupling is required

Costs: eventual consistency, harder debugging, event schema evolution, ordering guarantees.

```
EDA OPPORTUNITY — {integration point}
Current design: {synchronous call / tight coupling description}
Problem it causes: {fragility / blocking / coupling issue}
Proposal: {replace with event / add event alongside existing call}
Trade-off: {consistency model change / infra requirement / debugging overhead}
Recommendation: ADOPT / DEFER
```

---

#### CQRS

Applicable when:
- Read and write access patterns are significantly different
- Read performance requires a denormalised projection that conflicts with the write model
- The domain has a high read-to-write ratio requiring independent read optimisation

Costs: eventual consistency between write and read models, synchronisation complexity.

```
CQRS OPPORTUNITY — {domain / aggregate}
Current design: {describe the read/write coupling}
Problem it causes: {slow queries / model tension / scaling bottleneck}
Proposal: {separate read model / projection / materialised view}
Trade-off: {synchronisation overhead / consistency lag}
Recommendation: ADOPT / DEFER
```

---

#### SAGA

Applicable when:
- A business operation spans multiple services or databases and must be atomic at the
  business level, but a distributed transaction is not feasible
- Partial failure must trigger compensating actions

| Style | When to use |
|---|---|
| **Choreography** | Few services, simple flow. Lower infra cost, harder to trace. |
| **Orchestration** | Many services, complex flow, explicit failure handling. Easier to trace. Higher complexity. |

```
SAGA OPPORTUNITY — {workflow name}
Current design: {how the multi-step operation is handled today}
Problem it causes: {partial failure / no rollback path}
Steps in the saga: {list each step and its compensating action}
Recommended style: {choreography / orchestration — with rationale}
Trade-off: {complexity cost vs consistency guarantee}
Recommendation: ADOPT / DEFER
```

---

### Lens 7 — CAP Theorem and Database Selection

**Question**: Is the database choice the right fit for this component's consistency,
availability, and partition tolerance requirements?

Only flag when the current choice creates a concrete problem the access pattern reveals.

#### CAP Theorem

| Trade-off | Behaviour | Typical databases |
|---|---|---|
| **CP** | Rejects requests rather than return stale data | PostgreSQL (sync replication), etcd, ZooKeeper |
| **AP** | Returns potentially stale data rather than error | Cassandra, DynamoDB (eventual), NATS |
| **CA** | No partitions assumed — single-node only | SQLite, single-node PostgreSQL |

```
CAP TRADE-OFF — {component / store}
Current choice: {database used}
Required guarantee: {CP / AP — and why}
Risk of current choice: {over/under-engineered consistency}
Proposal: {keep / adjust replication config / consider alternative}
```

#### Database Type Selection

| Category | Best for |
|---|---|
| Relational (PostgreSQL) | Transactional data, financial records, normalised domain models |
| Document (MongoDB, DynamoDB) | Flexible schema, variable-structure content, catalogues |
| Key-Value (Redis) | Caches, sessions, distributed locks, rate limiting, counters |
| Wide-Column (Cassandra) | High-volume append workloads, time-series, activity feeds |
| Time-Series (TimescaleDB) | Metrics, monitoring, financial tick data |
| Search (Elasticsearch) | Full-text search, log aggregation |

```
DATABASE SELECTION — {component / store}
Current choice: {database used}
Access pattern: {point reads / range queries / aggregations / full-text / graph}
Transaction requirement: {ACID / eventual}
Mismatch identified: {yes / no}
Proposal: {keep / consider {alternative} for {specific access pattern}}
Recommendation: KEEP / CONSIDER ALTERNATIVE / DEFER
```

---

## Diagram Protocol

After completing the lenses, offer Mermaid diagrams for any finding that involves:
- State transitions (transaction lifecycle, saga steps, retry flows)
- Data flow across services or stores (write paths, event chains)
- Concurrency scenarios (lock acquisition, fencing token flow)
- Architectural patterns (Outbox, SAGA choreography/orchestration)

### When to generate a diagram

Always ask the developer first — do not generate unrequested diagrams:

```
Diagram available for: {finding title}
Type: {state diagram / sequence diagram / flowchart}
Useful for: {what it clarifies}
Generate? [Y/N]
```

If the developer confirms, generate the diagram inline using a fenced Mermaid block.

### Diagram types and use cases

**State diagram** — transaction / saga lifecycle:
```mermaid
stateDiagram-v2
    [*] --> Pending
    Pending --> Submitted : send tx
    Submitted --> Confirming : included in block
    Confirming --> Confirmed : N confirmations
    Confirming --> Failed : reorg detected
    Submitted --> Failed : reverted
```

**Sequence diagram** — Outbox Pattern write path:
```mermaid
sequenceDiagram
    participant Handler
    participant DB
    participant OutboxPoller
    participant NATS

    Handler->>DB: BEGIN TX<br/>INSERT cashback<br/>INSERT outbox_event
    DB-->>Handler: COMMIT
    OutboxPoller->>DB: SELECT pending events
    OutboxPoller->>NATS: Publish event
    OutboxPoller->>DB: UPDATE status = published
```

**Flowchart** — idempotency guard:
```mermaid
flowchart TD
    A[Receive event] --> B{event_id in processed_events?}
    B -- Yes --> C[ACK + skip]
    B -- No --> D[Execute business logic]
    D --> E[INSERT processed_events]
    E --> F[ACK message]
    D -- Error --> G[NAK message]
```

**Sequence diagram** — distributed lock + fencing token:
```mermaid
sequenceDiagram
    participant ServiceA
    participant Redis
    participant DB

    ServiceA->>Redis: SET lock NX PX 5000
    Redis-->>ServiceA: OK (token=7)
    ServiceA->>DB: UPDATE nonces SET value=N WHERE fence_token < 7
    DB-->>ServiceA: 1 row updated
    ServiceA->>Redis: DEL lock
    Note over ServiceA,DB: If Redis lock expires mid-flight,<br/>next holder gets token=8.<br/>DB rejects writes with token ≤ 8.
```

---

## Output Contract

Write to `~/ai-plans/{repo-name}/{slug}/system-design-analysis.md`.

### File Format

```markdown
# System Design Analysis — {slug}
**Date**: {YYYY-MM-DD}
**Phases analysed**: {list of phase names}

---

## Phase {N} — {Phase Name}

### Lens 1 — Atomicity
{findings or "No concerns identified."}

### Lens 2 — Idempotency
{findings or "No concerns identified."}

### Lens 3 — Consistency
{findings or "No concerns identified."}

### Lens 4 — Concurrency
{findings or "No concerns identified."}

### Lens 5 — Resilience, Fault Tolerance, and Scalability
{findings or "No concerns identified."}

### Lens 6 — Architectural Patterns (EDA / CQRS / SAGA)
{findings or "No EDA, CQRS, or SAGA opportunities identified for this phase."}

### Lens 7 — CAP Theorem and Database Selection
{findings or "Current database choices are appropriate for the identified access patterns."}

---

## Proposals Requiring Developer Approval

| # | Lens | Title | Risk if ignored | Complexity | Decision |
|---|---|---|---|---|---|
| 1 | ATOMICITY | {title} | {consequence} | low/med/high | [ ] Approve [ ] Reject [ ] Defer |
| 2 | CONCURRENCY | {title} | {consequence} | low/med/high | [ ] Approve [ ] Reject [ ] Defer |

---

## Diagram Requests

| # | Finding | Type | Status |
|---|---|---|---|
| 1 | {finding title} | state diagram | [ ] Requested [ ] Generated [ ] Skipped |

---

## Approved Decisions

{Filled in after developer responds.}
```

---

## Presentation Protocol

After writing the file, present a concise summary to the developer:

```
System Design Analysis complete — {N} phase(s) analysed.

{N} proposals require your decision:

1. [ATOMICITY] {title} — {one-line risk}
2. [IDEMPOTENCY] {title} — {one-line risk}
3. [CONCURRENCY] {title} — {one-line risk} (fencing token: yes/no)
4. [EDA] {title} — {opportunity or risk}
5. [DATABASE] {title} — {trade-off}

Diagrams available: {list titles} — generate any? [Y/N per item]

Full analysis: ~/ai-plans/{repo-name}/{slug}/system-design-analysis.md

Approve all / reject individual items / request changes before Coder starts.
```

**Wait for explicit developer approval. The Coder must not start until the developer responds.**

---

## Rules

- Never generate code — describe patterns and reference files only
- Never mark a proposal as approved on behalf of the developer
- Every lens must be answered for every phase — never skipped
- Never generate Mermaid diagrams without the developer confirming
- Apply text-sanitizer rules before writing the output file

---

## Permissions

- ✅ Read any file in the repository
- ✅ Read `~/ai-plans/{repo-name}/{slug}/implementation-plan*.md`
- ✅ Read `~/ai-plans/{repo-name}/{slug}/research.md`
- ✅ Write `~/ai-plans/{repo-name}/{slug}/system-design-analysis.md`
- ❌ Write production code
- ❌ Run terminal commands
- ❌ Mark proposals as approved without developer response
