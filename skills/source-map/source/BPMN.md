# BPMN Reference

**Source**: Business Process Model and Notation (BPMN) 2.0 Specification (OMG)
**Purpose**: Offline instruction for modeling business processes in initiative documentation. Adapted for Graphviz DOT rendering.

---

## When to Use BPMN

Use BPMN when you need to document **business processes** — workflows that involve multiple actors, decisions, events, and handoffs. BPMN is the standard notation that both technical and business stakeholders can read.

**Use in initiatives when**:
- A business process spans multiple departments or systems (order fulfillment, incident response, list management)
- You need to show who does what and in what order
- The process has complex branching, parallel paths, or error handling
- Stakeholders outside engineering need to understand the flow
- Arc42 Section 6 (Runtime View) or TOGAF Section 2 (Business Architecture) requires process documentation

**Don't use BPMN when**:
- The flow is purely technical (service-to-service calls) — use UML Sequence Diagram instead
- The flow is trivial (linear, no branching) — use ASCII inline
- You're modeling data, not process — use UML Class or ER diagrams

---

## Core Elements

### Flow Objects (the nodes)

| Element | Symbol | Description | DOT Shape |
|---------|--------|-------------|-----------|
| **Task** | Rounded rectangle | A unit of work performed by a participant | `shape=box style="filled,rounded"` |
| **Sub-process** | Rounded rectangle with `+` marker | Collapsed group of tasks | `shape=box style="filled,rounded"` + label with `[+]` |
| **Event** (start) | Thin circle ○ | Triggers the process | `shape=circle width=0.4 style=filled fillcolor=white` |
| **Event** (intermediate) | Double circle ◎ | Something happens during the process | `shape=doublecircle width=0.4` |
| **Event** (end) | Thick circle ● | Process terminates | `shape=circle width=0.4 style=filled fillcolor=black penwidth=3` |
| **Gateway** (exclusive) | Diamond with X | XOR — exactly one path taken | `shape=diamond` with `X` label |
| **Gateway** (parallel) | Diamond with + | AND — all paths taken simultaneously | `shape=diamond` with `+` label |
| **Gateway** (inclusive) | Diamond with O | OR — one or more paths taken | `shape=diamond` with `O` label |
| **Gateway** (event-based) | Diamond with pentagon | Wait for one of several events | `shape=diamond` with `⬠` label |

In DOT, represent a task's type with a stereotype in its label: `[label="<<Service Task>>\nCompute CBF"]`.
The templates below use `<<User Task>>` (human), `<<Service Task>>` (automated), `<<Send Task>>` and
`<<Receive Task>>` (message out / in).

### Connecting Objects (the edges)

| Element | Style | Description |
|---------|-------|-------------|
| **Sequence Flow** | Solid arrow → | Order of activities within a pool |
| **Message Flow** | Dashed arrow → | Communication between pools (different organizations/systems) |
| **Association** | Dotted line | Links artifacts (data objects, annotations) to flow objects |

### Swimlanes (the containers)

| Element | Description | DOT Equivalent |
|---------|-------------|---------------|
| **Pool** | Represents a participant (organization, system, major actor) | `subgraph cluster_<name>` |
| **Lane** | Subdivision within a pool (role, department, service) | Nested `subgraph cluster_<name>` or visual grouping |

### Artifacts (supplementary info)

| Element | Description |
|---------|-------------|
| **Data Object** | Information consumed or produced (document, message) |
| **Data Store** | Persistent storage (database, file system) |
| **Annotation** | Free-text comment attached to an element |
| **Group** | Visual grouping of elements for documentation purposes |

---

## Gateway Patterns

### Exclusive Gateway (XOR) — Choose ONE path

```
        [condition A]
       /              \
  --> <X> -----------> [Task A]
       \              /
        [condition B]
         \           /
          [Task B] --
```

Use when: exactly one condition is true. Default flow for "else" case.

### Parallel Gateway (AND) — ALL paths simultaneously

```
       +--> [Task A] --+
       |                |
  --> <+>              <+> -->
       |                |
       +--> [Task B] --+
```

Use when: multiple activities must happen concurrently. The closing gateway waits for ALL paths to complete.

### Inclusive Gateway (OR) — ONE or MORE paths

```
       +--> [Task A] --+
       |                |
  --> <O>              <O> -->
       |                |
       +--> [Task B] --+
```

Use when: multiple conditions can be true simultaneously. The closing gateway waits for all taken paths.

### Event-Based Gateway — Wait for FIRST event

```
       +--> (Message received) --> [Handle message]
       |
  --> <⬠>
       |
       +--> (Timer expired) --> [Handle timeout]
```

Use when: the process waits for one of several possible events. First event wins.

---

## DOT Templates

### Process with Exclusive Gateway

```dot
digraph BPMN_XOR {
    graph [label="Process: CBF Sync Strategy" labelloc=t fontsize=16 fontname="Arial" rankdir=TB]
    node [shape=box style="filled,rounded" fontname="Arial" fontsize=10 fillcolor="#B5FFFF"]
    edge [fontname="Arial" fontsize=9]

    start [label="" shape=circle fillcolor="#C9E7B7" width=0.4]
    receive [label="<<Receive Task>>\nReceive\nlist event"]
    gw1 [label="X" shape=diamond fillcolor="#FFFFB5" width=0.5]
    full_rebuild [label="<<Service Task>>\nFull CBF rebuild\nfrom MongoDB"]
    granular [label="<<Service Task>>\nIncrement/Decrement\nCBF item"]
    publish [label="<<Send Task>>\nPublish\nbloom.sync"]
    end_node [label="" shape=circle fillcolor="#FFB5B5" width=0.4 penwidth=3]

    start -> receive
    receive -> gw1
    gw1 -> full_rebuild [label="list created/\ndeleted/bulk update"]
    gw1 -> granular [label="single item\nadd/remove (Phase 2)"]
    full_rebuild -> publish
    granular -> publish
    publish -> end_node
}
```

### Process with Parallel Gateway and Swimlanes

```dot
digraph BPMN_Parallel {
    graph [label="Process: List Upload End-to-End" labelloc=t fontsize=16 fontname="Arial" rankdir=TB]
    node [shape=box style="filled,rounded" fontname="Arial" fontsize=10 fillcolor="#B5FFFF"]
    edge [fontname="Arial" fontsize=9]

    // Pool: Backoffice
    subgraph cluster_backoffice {
        label="Backoffice" style="filled,rounded" fillcolor="#FFFFB520"
        upload [label="<<User Task>>\nPartner uploads CSV"]
        save_mongo [label="<<Service Task>>\nSave to MongoDB"]
        publish_event [label="<<Send Task>>\nPublish list.updated"]
    }

    // Pool: Estimator
    subgraph cluster_estimator {
        label="Estimator" style="filled,rounded" fillcolor="#B5FFFF20"
        receive_event [label="<<Receive Task>>\nlist-manager\nreceives event"]
        compute_cbf [label="<<Service Task>>\nCompute CBF"]
        broadcast [label="<<Send Task>>\nBroadcast bloom.sync"]
    }

    // Pool: Scoring
    subgraph cluster_scoring {
        label="Scoring" style="filled,rounded" fillcolor="#C9E7B720"
        receive_cbf [label="<<Receive Task>>\nestimator receives\nbloom.sync"]
        hotswap [label="<<Service Task>>\nHot-swap CBF\nin memory"]
        ready [label="<<Service Task>>\nReady for scoring\n(~850ns lookups)"]
    }

    start [label="" shape=circle fillcolor="#C9E7B7" width=0.4]
    end_node [label="" shape=circle fillcolor="#FFB5B5" width=0.4 penwidth=3]

    start -> upload
    upload -> save_mongo
    save_mongo -> publish_event
    publish_event -> receive_event [style=dashed label="NATS\nmessage flow"]
    receive_event -> compute_cbf
    compute_cbf -> broadcast
    broadcast -> receive_cbf [style=dashed label="NATS\nmessage flow"]
    receive_cbf -> hotswap
    hotswap -> ready
    ready -> end_node
}
```

---

## BPMN vs UML Activity Diagram

| Aspect | BPMN | UML Activity |
|--------|------|-------------|
| **Audience** | Business + technical stakeholders | Primarily technical |
| **Swimlanes** | Pools (organizations) + Lanes (roles) | Partitions |
| **Events** | Rich event taxonomy (message, timer, signal, error, etc.) | Start/end nodes only |
| **Gateways** | XOR, AND, OR, Event-based (explicit symbols) | Decision + fork/join (generic) |
| **Message flow** | Dashed arrows between pools (cross-org communication) | Not native |
| **Standardization** | ISO 19510 standard, widely adopted | UML spec, part of larger standard |
| **Best for** | Cross-functional business processes | Technical algorithms and workflows |

**Rule of thumb**: If the process involves **multiple organizations or departments** communicating via messages, use BPMN. If it's a **single-service algorithm with decisions and loops**, use UML Activity.

---

## Modeling Best Practices

1. **One pool per participant** (system, organization, major actor). Don't mix responsibilities in a single pool.
2. **Start simple** -- model the happy path first, then add error paths and edge cases.
3. **Name tasks with verb + noun** -- "Upload CSV", "Compute CBF", "Send notification". Not "CSV" or "Processing".
4. **Every gateway must have a matching closing gateway** for parallel and inclusive gateways. Exclusive gateways may converge implicitly.
5. **Label all sequence flows from gateways** with the condition. The default flow (else) should be marked.
6. **Use message flows (dashed) between pools**, sequence flows (solid) within a pool. Never mix them.
7. **Keep it on one page** -- if the diagram is too complex, decompose into sub-processes.
8. **Avoid crossing lines** -- rearrange layout if edges cross.
9. **Use intermediate events for waiting** -- a timer event is clearer than a "wait" task.
10. **Don't over-detail** -- BPMN shows the flow, not the implementation. Leave technical details to UML or code.

---

## Diagram Conventions

BPMN diagrams follow the shared DOT → PNG workflow, the `images/bpmn_<process_name>` naming rule, and the shared pastel palette — see `Diagram Conventions.md`.
