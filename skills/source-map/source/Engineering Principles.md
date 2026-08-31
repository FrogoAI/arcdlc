# Engineering Principles (POL-ENG-001)

This document records one organization's engineering practice, kept as a worked example of a
delivery standard. It is not a universal mandate: where a rule names a specific stack — event-driven
internal communication, Swagger, Kubernetes — it applies **only** to projects built on that stack. On
any other stack, those rules are conditional and should be read as "the equivalent rule for your
transport / API-doc / deployment target", or ignored outright. The stack-neutral rules (KISS, DRY,
Clean Code, non-decreasing test coverage, API versioning discipline, review duties, branch, commit
and MR conventions) are the portable part.

---

## Definitions

* **Acceptance Criteria (AC)**: A set of predefined requirements, constraints, and edge cases that must be met for a task to be considered complete.
* **Incident**: An unplanned interruption to an IT service or a reduction in the quality of an IT service.
* **System Down-time**: Any period during which a system is unavailable or fails to perform its primary functions, impacting users or business operations.
* **Post-mortem**: A formal document analyzing an incident, its impact, the actions taken to resolve it, the root cause, and the follow-up actions required to prevent recurrence.
* **Root Cause**: The fundamental issue which, if addressed, will prevent the incident from recurring.
* **KISS (Keep It Simple and Smart)**: See `source/KISS.md` for the full definition and its rules.
* **DRY (Don't Repeat Yourself)**: A principle aimed at reducing repetition of software patterns.
* **Conventional Commits**: The commit message convention this policy mandates. See
  `source/Conventional Commits.md` for the full specification.
* **Message Broker**: Any queue, log, or pub/sub middleware carrying asynchronous messages between
  services (Kafka, NATS, RabbitMQ, SQS/SNS, Pub/Sub, Redis Streams, …).
* **Subject**: The broker's addressable channel for a message — the term a concrete broker may call a
  topic, subject, queue, stream, or exchange/routing key.

---

## Procedures

### Architectural & Design Principles

* **Communication**: All internal service communication must use an event-driven approach. REST communication with a backoff mechanism must be used for all third-party integrations.
* **Documentation**: All external communication and APIs must be documented in Swagger.
* **Core Principles**: All design must prioritize **Performance by Design**. Solutions must adhere to **KISS** to avoid unnecessary complexity and **DRY** to prevent code repetition.
* **Architecture Changes**: Any changes to architecture require discussion and documentation via an Architecture Inception Canvas or Architecture Communication Canvas before implementation.

### Development & Testing Procedures

* **Code Quality**: Developers must follow **Clean Code** principles.
* **Test Coverage**: Test coverage must be improved over time and is not permitted to decrease in any merge request. Reviewers must check test coverage.
* **Test Types**: Integration tests are required for communication with other services and databases. Unit tests are required for functionality not covered by integration tests.
* **Task Completion**: A task is "Done" when it is implemented, every one of its acceptance criteria is demonstrably met, the change is validated (tests and lint green), and the work is committed. This matches the `DONE` status of the ArcDLC plan format, which `/arcdlc:execute` enforces; nothing beyond the task's own scope gates it.
* **Release Gate**: Deployment to production and QA verification form a separate release gate that follows task completion. A release is not shippable until its changes are deployed and QA-verified, but that gate is tracked on the release — never on the individual task's status.

### Code Review & Deployment Procedures

* **Reviewer Responsibility**: Reviewers must provide feedback within two business days. The review must focus on logic, algorithms, adherence to Engineering Principles, and fulfillment of the Acceptance Criteria.
* **Developer Responsibility**: The developer who authors the code is responsible for finding a reviewer, ensuring all pipeline checks pass, merging approved code, and shepherding the change through the release gate above.
* **Deployment**: All infrastructure must be deployed in a unified way. If a component is deployable in Kubernetes, it must be deployed in Kubernetes.

### Branch Naming Convention

* **Epic Branches**: All branches related to an epic must follow this naming structure:
  * **Format**: `epic/{TASK}_{TITLE}`
  * **{TASK}**: The full Task ID (e.g., `DAF-100`).
  * **{TITLE}**: A short, descriptive title, with words separated by a dash (`-`).
  * **Example**: `epic/DAF-123-user-authentication-flow`

### Commit Message Convention

* **Mandatory Convention**: Every commit message must follow **Conventional Commits 1.0.0** — the
  specification in `source/Conventional Commits.md` is normative. Reviewers must reject an MR whose
  commits do not conform.
* **Format**: `<type>[optional scope][!]: <description>`, followed by an optional body and optional
  footers, separated by blank lines.
* **Type**: `feat` for a new capability, `fix` for a bug fix; `docs`, `test`, `refactor`, `perf`,
  `build`, `ci`, `chore`, `style` for everything else. The type is lowercase.
* **Scope**: A scope in parentheses names the affected area — the initiative slug, service, or module
  (e.g. `feat(scoring): …`). Use one consistent vocabulary per repository.
* **Description**: Imperative mood, lowercase, no trailing period, and it must state the change, not
  the task ID.
* **Breaking Changes**: A backward-incompatible change must carry `!` after the type/scope **and** a
  `BREAKING CHANGE: <what breaks and how to migrate>` footer.
* **Task Reference**: Every commit must reference its work item in a footer — `Refs: {TASK}` with the
  full Task ID (e.g. `Refs: DAF-123`) — so a commit is always traceable to its acceptance criteria.
* **Granularity**: One task, one commit. When a change spans several types, the dominant type wins;
  do not merge unrelated tasks into a single commit.
* **AI-Assisted Work**: A commit produced with AI assistance must be marked as such in a footer, per
  the project's disclosure rule.

### Merge Request (MR) Submission

* **Decomposition**: An epic may be broken down into multiple subtasks. Developers are permitted to create separate branches and MRs for subtasks.
* **Final Merge**: The final solution for a single epic or bug must be consolidated into **one final Merge Request**. All subtask branches must be merged into the primary epic branch before final review and merge.

### API Development & Versioning

* **Mandatory Versioning**: Each API endpoint must be deployed with a specific version in the URL path (e.g., `/v1/users`). The initial version of any new service must be `v1`.
* **Backward Incompatibility**: Each change that is not backward-compatible must be deployed to a new version number higher than the previous one (e.g., `v2`).
* **Deprecation**: When a new API version is released, partners must be informed how long the previous version will be supported. By default, a previous version must exist for at least one month.
* **Legacy Endpoints**: All legacy endpoints without a version must have a corresponding versioned endpoint created (starting from `v1`). A migration plan with a due date must be defined.
* **Exceptions**: Adding a "non-versioned" API must be discussed individually and requires explicit written approval from the CTO and the BE Lead.

### Asynchronous Messaging and Subject Naming

These rules are broker-agnostic: they hold for any message broker (Kafka, NATS, RabbitMQ, SQS/SNS,
Pub/Sub, …). Read "subject" as whatever the chosen broker calls its addressable channel — topic,
subject, queue, stream, or exchange with routing key — and map the naming pattern onto that broker's
own constraints (its allowed separator and character set) without changing the segment order.

* **Naming Convention**: A subject name must follow the pattern:
  `{project}.{system}.{service}.{kind of event (sync or async)}.{what is the message about}`.
  * Example: `[company].scoring.api-server.async.scoring_event`.
  * On a broker that forbids `.` as a separator, substitute its idiomatic separator (e.g. `-`) and
    keep the five segments in the same order.
* **Responsibility**: The service that writes/produces the message owns the subject name and the
  message's data structure — including its schema and compatibility guarantees.
* **Versioning**: To create a new version of an existing subject, append a version number (e.g. `v2`,
  `v3`) to the end of the subject name. The initial version carries no `v1` suffix.
* **Renaming Subjects**: Rename by addition, never in place: publish to the new subject alongside the
  old one, deploy new versions of all consumers, and remove the old producer and subject only once
  100% of traffic has moved to the new name.
* **Request/Reply**: A subject used in a request/reply pattern must have exactly one consumer. Every
  request/reply handler must send a reply back to the caller, even when the outcome is an error.
* **Delivery Semantics**: Assume at-least-once delivery. Consumers must be idempotent, and every
  message must carry enough identity (message ID or business key) for a consumer to detect a replay.
