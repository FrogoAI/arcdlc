# Workflow Policy (POL-ENG-002)

| Field | Value |
| :---- | :---- |
| Policy Name | Workflow Policy |
| Unique ID | POL-ENG-002 |
| Status | Draft |
| Owner | CTO / Engineering Management |
| Approver | CTO |
| Created | 2026-04-27 |
| Effective Date | After CTO approval |
| Next Review Date | 2026-10-27 |
| Related Policies | POL-GOV-001 Policy of Policies, POL-TECH-001 Policy of Initiatives, POL-ENG-001 Engineering Principles |

## Purpose

This policy defines the day-to-day product, engineering, QA, and deployment workflow for the team.
It turns the initiative lifecycle, Engineering Principles, trunk-based development, AIC/TSC documentation, plan execution, QA, and Kubernetes deployment into one operational instruction.

The document also clarifies Product Owner responsibilities. The Product Owner owns product value, priorities, acceptance, route readiness, and
continuous clarification for the team.

## Policy Statement

All product and engineering work must flow through a visible lifecycle:

```text
Idea -> Product Discovery -> Initiative -> AIC/TSC -> Refinement -> Plan Tasks
     -> Trunk-Based Development -> QA -> Deployment -> Production Validation -> Done
```

Work must not start unless the team understands the product outcome, acceptance criteria, architecture constraints, test strategy, and deployment path. Work is not done until it is deployed to production, validated by QA, and accepted against the agreed criteria.

## Scope

This policy applies to the following team roles:

- CTO
- DevOps Engineer
- Engineering Team Lead
- Fullstack Engineer
- AQA Engineer
- Project Manager
- Product Owner
- Frontend Engineer
- Backend Engineer

The policy applies to:

- Product discovery and backlog work.
- Initiative creation and architectural definition.
- AIC and Tech Stack Canvas documentation.
- Plan task queue management in `docs/aics/<slug>/plan.md`.
- Trunk-based development and merge request flow.
- Testing, QA, and acceptance.
- Kubernetes, ArgoCD, route migration, DNS, and production deployment.
- Production validation, rollback, and incident follow-up.

## Governing References

This policy operationalizes the following materials:

- `source/Policy of Policies.md` (POL-GOV-001) for policy lifecycle and RACI definitions.
- `source/Policy of Initiatives.md` (POL-TECH-001) for initiative lifecycle and AIC requirements.
- `source/Engineering Principles.md` (POL-ENG-001) for engineering, testing, review, deployment, branch, API, and NATS standards.
- `source/AIC Template.md` for Architecture Inception Canvas structure.
- `source/Tech Stack Canvas.md` for technology decision structure.
- The organization's `<org>/k8s` repo, `proxy.md`, for current Kubernetes frontend migration and proxy deployment rules.
- Trunk-Based Development reference: https://trunkbaseddevelopment.com/

## Definitions

- **AIC**: Architecture Inception Canvas. Mandatory architecture/business framing document for
  initiatives.
- **TSC**: Tech Stack Canvas. Mandatory technology decision document for initiatives.
- **Trunk**: The main integration branch. It is the only permanent development branch.
- **Short-lived branch**: A branch created for a small task and merged back quickly, normally within
  one working day and no later than two working days unless the Engineering Team Lead approves an
  exception.
- **Feature flag**: A runtime switch that lets incomplete or risky behavior be merged without being
  exposed to users.
- **Release candidate**: A tested version that is ready for environment deployment.
- **Route takeover**: Moving a specific frontend route from legacy EC2 behavior to Kubernetes
  frontend ownership.
- **Done**: Deployed to production, QA validated, and accepted against AC.
- **Accountable (A)**: Owns the correct completion of the task.
- **Responsible (R)**: Performs the work.
- **Consulted (C)**: Provides required input before completion.
- **Informed (I)**: Must be kept up to date.

## Core Workflow Principles

The executable unit of work is a task block in `docs/aics/<slug>/plan.md`, worked by `/arcdlc:execute` —
not a ticket in an external tracker. Any tracker is a reporting mirror of that plan, never its source of
truth: scope, acceptance criteria, and status live inside the task block.

1. **Performance by Design**
   - Performance is part of scope, architecture, review, testing, and deployment.
   - Each initiative must define expected sizing numbers or explicitly state why they are unknown.

2. **KISS**
   - Prefer the smallest clear, efficient, modular solution that solves the product and operational problem.
   - Keep code simple to read, understand, scale, and extend; do not make it simplistic.
   - Avoid unnecessary services, queues, flags, abstractions, or dependencies.

3. **DRY**
   - Reuse proven modules and patterns.
   - Do not create abstractions before repeated behavior is stable and understood.

4. **MDCA and DDD**
   - Organize work by business domain, not technical layer.
   - Keep bounded contexts clear. Modules communicate through explicit contracts.

5. **Clean Code**
   - Code must be readable, tested, reviewed, and maintainable.
   - Go code follows `gofumpt`, `golangci-lint v2`, effective Go naming, error wrapping, and
     package ownership rules from POL-ENG-001.

6. **Event-Driven Internal Communication**
   - Internal service communication uses NATS unless explicitly approved otherwise.
   - NATS subjects follow `{project}.{system}.{service}.{sync|async}.{event_description}`.

7. **REST With Backoff for External APIs**
   - External and third-party integrations use REST with retry/backoff behavior and clear failure
     handling.

8. **API Versioning**
   - New APIs start with `/v1/...`.
   - Breaking changes require a new version and an explicit deprecation plan.

9. **Kubernetes-First Deployment**
   - If a component can be deployed in Kubernetes, it must be deployed in Kubernetes.
   - All deployable components must expose health/readiness behavior suitable for the platform.

10. **Documentation Before Implementation**
    - Standard initiatives require AIC and TSC before development starts.
    - Architecture changes require AIC/TSC updates before implementation.

## Workflow Stages

### 1. Product Discovery

Lead: Product Owner.

Trigger: new product idea, customer request, operational pain, incident follow-up, compliance need,
or technical risk.

Required output:

- Problem statement.
- User or stakeholder group.
- Business value.
- Initial scope and non-goals.
- Success criteria.
- Known constraints.
- Open questions.

Rules:

- The Product Owner owns the "why" and priority.
- The Project Manager owns visibility and follow-up.
- The Engineering Team Lead confirms whether the work needs architecture review; the CTO makes the
  final call for strategic or high-risk architecture questions.
- No development task is created until the problem and expected outcome are clear.

### 2. Initiative Creation

Lead: Engineering Team Lead.

Standard initiatives must follow POL-TECH-001 and use AIC plus TSC.

Required AIC inputs:

- Business Case.
- Functional Overview.
- Quality Goals from ISO-25010.
- Organizational Constraints.
- Technical Constraints, including POL-ENG-001 references.
- Business Context.
- Architectural Hypotheses.
- Technical Challenges and Risks.
- Tasks after refinement.

Required TSC inputs:

- Business Feature Description.
- System Name and Description.
- Sizing Numbers.
- Major Quality Attributes.
- Services.
- Architecture Representation.
- Sequence Flow when useful.
- Frontend, backend, and data technology choices.
- APIs and integrations.
- Security and compliance.
- Testing and QA.
- Infrastructure and deployment.
- Monitoring and analytics.
- Development workflow and collaboration.

Rules:

- Missing information must be asked directly. Do not guess.
- AIC/TSC must be updated when decisions change.
- Diagrams must live in the initiative `images/` folder as `.dot` source plus `.png` render where
  Graphviz is used.
- Technical initiatives may use the simplified AIC path, but must still include risks and tasks.

### 3. Refinement and Task Decomposition

Lead: Project Manager.

Trigger: AIC/TSC are reviewed enough to plan implementation.

Definition of Ready for a plan task:

- The task has a clear user/business or technical outcome.
- AC is written and testable.
- Design/architecture link is attached when relevant.
- Dependencies and affected services are known.
- Test approach is known.
- Deployment impact is known.
- The task is small enough for a short-lived branch.
- Owner is assigned.

Task rules:

- Decomposition produces task blocks in `docs/aics/<slug>/plan.md`, in the ArcDLC plan format;
  `/arcdlc:plan` writes them and `/arcdlc:execute` works them one at a time.
- One task per `###` block, sized so a single session can implement, test, and commit it.
- Every block carries `WHAT`, `WHERE`, `WHY`, `Acceptance`, `References`, and `Status`; add `HOW`
  whenever a decision would otherwise be left to the executor's judgment.
- Blocks are ordered by dependency — the queue runs top-to-bottom.
- Blocks link back to the AIC/TSC through `References`.
- Status lives on the `- Status:` line inside the block. A tracker ticket, where one exists, mirrors
  the block for reporting; it never overrides it.

### 4. Trunk-Based Development

Lead: Engineering Team Lead. Trunk discipline, short-lived branches, feature flags, CI gates, branch
naming, and merge request review are defined in `source/tbd.md`.

### 5. Engineering Implementation Rules

Lead: Engineering Team Lead.

General rules from POL-ENG-001:

- Use Clean Code principles.
- Avoid `fmt.Print*`, `log.Print*`, and bare `print*` in application code. Use structured logging.
- Never ignore errors.
- Error strings are lowercase and wrap context with `%w`.
- Prefer early returns and left-aligned happy path.
- Keep functions focused; split functions that grow beyond clear readability.
- Avoid flag arguments that make one function do two different things.
- Use full-word naming and correct acronym casing, such as `userID` and `HTTPClient`.
- No `Get` prefix for simple accessors.
- Use compile-time interface checks where appropriate.
- Each package owns its env config prefix.
- Logs go to stdout only.
- Services shut down gracefully.
- Jobs must be idempotent.

Go/MDCA rules:

- Domain models have storage tags only.
- DTOs have JSON tags only.
- Handler never sees domain storage models.
- Repository never sees DTOs.
- Consumer-side packages define ports/interfaces where they are used.
- Cross-service internal communication uses events, not service-to-service HTTP.

Testing rules:

- Unit tests cover business logic not covered by integration tests.
- Integration tests are required for service and database communication.
- Handler tests use framework test helpers and mocked services.
- Flow tests validate cross-service behavior against a running stack where required.
- QA automation should be added for stable product paths and critical regressions.

### 6. QA and Acceptance

Lead: AQA Engineer.

QA starts before coding.

Rules:

- The AQA Engineer defines test strategy during refinement.
- Each task must have "how to test" or equivalent QA evidence.
- Engineers support QA with test data, feature flag instructions, environment notes, and logs.
- The Product Owner validates product behavior; the AQA Engineer validates quality and regression
  risk.
- A task cannot be moved to `DONE` only because code was merged.

Minimum QA evidence:

- AC pass/fail result.
- Environment tested.
- Build or commit reference.
- Regression scope.
- Known issues or explicit "none known".
- Product acceptance result for user-facing changes.

### 7. Deployment Readiness

Lead: DevOps Engineer.

Before any non-trivial deployment, the team must confirm:

- AIC/TSC or task context is current.
- CI passes.
- MR is approved and merged according to branch policy.
- Database migrations are reviewed and reversible or forward-safe.
- Feature flags are configured.
- Rollback plan exists.
- Monitoring and logs are available.
- QA sign-off exists for the target environment.
- Product go/no-go is explicit for user-facing changes.
- The DevOps Engineer confirms platform readiness for infrastructure-affecting changes.
- The Engineering Team Lead confirms technical readiness.
- The CTO is consulted for high-risk, architecture-impacting, or exception-based production
  releases.

Deployment environments:

- `dev`: early engineering validation.
- `stage`: integrated QA and product review.
- `beta`: controlled production-like exposure.
- `prod`: production users.

Done means:

- Production deployed.
- QA validated production or the agreed production-safe smoke scope.
- Product acceptance completed when user-facing.
- The plan task block is marked `DONE`.
- Monitoring checked after release.

## Deployment Policy

### Kubernetes and ArgoCD

- Deployable infrastructure and services must be deployed through Kubernetes where possible.
- ArgoCD is the expected synchronization mechanism for Kubernetes application state.
- Runtime configuration must be environment-specific and reviewed before sync.
- Health, readiness, and liveness behavior must be available for deployable services.
- Logs must go to stdout and be visible in the standard observability stack.

## Allowed and Prohibited Conduct

### Allowed

- Start discovery before all answers are known.
- Use AIC/TSC to expose gaps and risks.
- Merge small, safe slices frequently.
- Hide incomplete work behind feature flags.
- Deploy incrementally, route by route, through an explicit route-ownership list.
- Roll back by removing route ownership from Kubernetes.
- Ask for CTO/Architect review when architecture risk is high.

### Prohibited

- Starting development without clear AC.
- Starting standard initiative implementation without AIC and TSC.
- Keeping long-lived feature branches.
- Decreasing test coverage in a merge request.
- Deploying a breaking API change without version bump.
- Using service-to-service HTTP for internal communication without explicit approval.
- Claiming the catch-all `/` route before an approved final cutover.
- Changing DNS before legacy fallback, ALB health, TLS, and rollback are validated.
- Marking a plan task `DONE` before production deployment and QA validation.
- Treating Product Owner approval as optional for user-facing production changes.
