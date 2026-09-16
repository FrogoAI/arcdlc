# Diagram Conventions

**Purpose**: The single home for how initiative diagrams are produced, named, and colored. Applies to every notation in this library: C4, UML, BPMN, flowcharts, and the diagrams embedded in an arc42 document.

This file holds the shared rules plus the DOT templates for UML, BPMN, and flowcharts. C4 keeps its own templates in `C4.md` and ArchiMate keeps its own in `TOGAF.md`, because both are also `/arcdlc:aic` output formats.

---

## File Convention

All diagrams live in the `images/` folder inside the initiative directory, in **dual format**: the DOT source is committed next to the compiled image, so a diagram is always editable rather than redrawn.

1. Write DOT source: `images/<name>.dot`
2. Compile: `dot -Tpng images/<name>.dot -o images/<name>.png`
3. Embed in docs: `![Title](images/<name>.png)`

**Naming**: lowercase, underscores, descriptive, prefixed by notation.

| Notation | DOT source | Embed form | Examples |
|----------|------------|------------|----------|
| C4 | `images/c4_<level>_<name>.dot` | `![C4 <Level>: <Name>](…)` | `c4_context`, `c4_container`, `c4_component_list_manager`, `c4_dynamic_upload`, `c4_deployment_prod`, `c4_landscape` |
| UML | `images/<type>_<name>.dot` | `![Title](…)` | `component_lists`, `sequence_upload`, `state_list_sync`, `deployment_k8s`, `usecase_lists`, `class_domain`, `activity_cbf_rebuild` |
| BPMN | `images/bpmn_<process_name>.dot` | `![BPMN: <Process Name>](…)` | `bpmn_order_intake`, `bpmn_list_approval` |
| Flowchart | `images/flow_<name>.dot` | `![Flow: <Name>](…)` | `flow_request_validation`, `flow_cbf_rebuild`, `flow_list_entry_routing`, `flow_migration_procedure`, `flow_error_handling` |
| ArchiMate (TOGAF) | `images/archimate_<layer>.dot` | `![ArchiMate: <Layer>](…)` | `archimate_business`, `archimate_application` |

- C4 `<level>`: `context`, `container`, `component`, `code`, `dynamic`, `deployment`, `landscape`.
- UML `<type>`: the diagram type — `component`, `sequence`, `state`, `deployment`, `usecase`, `class`, `activity`.

**Tooling**:

- **Primary**: Graphviz DOT, for every notation in this library.
- **Fallback**: Markdown ASCII for simple sequence flows — inline in the document, no image file needed.

---

## Color Palettes

Two palettes are in use and they are **not interchangeable**. Choose by the **notation of the diagram, not by the document it appears in**: a C4 diagram keeps the C4 palette even inside a flowchart-heavy document, and a UML, BPMN, or flowchart diagram keeps the shared pastel palette even inside a C4-heavy document. Never mix the two inside one diagram.

### C4 diagrams — blue element family

**Applies to C4 diagrams only** (context, container, component, code, dynamic, deployment, landscape). Color encodes the *role* of the element.

| Element | Hex | Usage |
|---------|-----|-------|
| Person | `#08427B` | Dark blue, white text |
| Your System / Container | `#1168BD` / `#438DD5` | Blue family, white text |
| New Container | `#2694AB` | Teal, white text — highlights new additions |
| Component | `#85BBF0` | Light blue, black text |
| External System | `#999999` | Grey, white text |
| Deprecated | `#FFB5B5` | Red-tinted, black text |
| Infrastructure Node | `#C9E7B7` | Green (matches ArchiMate Technology layer) |
| System Boundary | `#1168BD` border | Dashed cluster outline |

### UML, BPMN, and flowchart diagrams — shared pastel palette

**Applies to UML UML, BPMN, and flowchart diagrams.** Color encodes the *kind of step*.

| Element | Color | Hex | Usage |
|---------|-------|-----|-------|
| Start / Success | Pale green | `#C9E7B7` | Entry points, successful outcomes |
| Process step | Pale cyan | `#B5FFFF` | Normal processing steps |
| Decision | Pale yellow | `#FFFFB5` | Decision diamonds, gateways |
| Error / End (failure) | Pale red | `#FFB5B5` | Error states, failure exits |
| Data store | Pale cyan | `#B5FFFF` | Database, file storage |
| Annotation | White | `#FFFFFF` | Comments, notes |

### Reading the overlap

- Two hexes appear in both palettes with different meanings: `#C9E7B7` (C4: infrastructure node — shared: start/success) and `#FFB5B5` (C4: deprecated — shared: error/failure). Interpret them in the palette of the diagram's notation.
- One deliberate crossover: UML use-case diagrams color external actors with the C4 person and external-system hexes (`#08427B`, `#999999`) so actors read the same across notations — see the use case template below.
- ArchiMate diagrams use a third, layer-based palette that stays with its notation in `TOGAF.md` (`## ArchiMate Diagram Conventions`).

---

## Choosing a notation

Draw a diagram only when it shows something the reader cannot get from the text. Then pick by what you
need to show:

| You need to show | Draw |
|---|---|
| High-level services or modules | C4 container or component (`C4.md`) |
| Infrastructure and deployment | C4 deployment, or a UML deployment diagram |
| Domain types and their relationships | class diagram |
| How components interact in one scenario | sequence diagram (ASCII inline is fine) |
| A workflow with decisions or parallel paths | activity diagram, or BPMN when roles matter |
| An entity's lifecycle | state machine |
| What actors can do | use case diagram |
| A business process with roles, events, and handoffs | BPMN |
| A procedure, an algorithm, or error handling | flowchart |
| Enterprise layers | ArchiMate (`TOGAF.md`) |

BPMN over an activity diagram when the process crosses roles or waits on external events; activity
diagram when it is one system's internal logic.

---

## UML templates

Component, activity, state machine, use case, and class. Sequence diagrams are usually clearer as inline ASCII than as DOT.

### Component Diagram

```dot
digraph Component {
    graph [label="Component Diagram: <System>" labelloc=t fontsize=16 fontname="Arial" rankdir=LR]
    node [shape=box style="filled,rounded" fontname="Arial" fontsize=10]
    edge [fontname="Arial" fontsize=9]

    api [label="<<component>>\npolicy-api" fillcolor="#B5FFFF"]
    mgr [label="<<component>>\nlist-manager" fillcolor="#B5FFFF"]
    est [label="<<component>>\nestimator" fillcolor="#B5FFFF"]

    mongo [label="<<component>>\nMongoDB" fillcolor="#C9E7B7" shape=cylinder]
    nats [label="<<component>>\nNATS" fillcolor="#C9E7B7" shape=hexagon]

    api -> mongo [label="reads/writes"]
    api -> nats [label="publishes"]
    nats -> mgr [label="subscribes"]
    mgr -> mongo [label="reads"]
    mgr -> nats [label="bloom.sync"]
    nats -> est [label="subscribes"]
}
```

### Activity Diagram

```dot
digraph Activity {
    graph [label="Activity: CBF Rebuild" labelloc=t fontsize=16 fontname="Arial" rankdir=TB]
    node [shape=box style="filled,rounded" fontname="Arial" fontsize=10 fillcolor="#B5FFFF"]
    edge [fontname="Arial" fontsize=9]

    start [label="" shape=circle fillcolor=black width=0.3]
    receive [label="Receive list.updated\nNATS event"]
    read [label="Read list entries\nfrom MongoDB"]
    compute [label="Compute CBF\n(17 hash functions)"]
    serialize [label="Serialize CBF\nto binary"]
    publish [label="Publish bloom.sync\nvia NATS"]
    end_node [label="" shape=doublecircle fillcolor=black width=0.3]

    decision [label="Full rebuild\nor granular?" shape=diamond fillcolor="#FFFFB5"]
    increment [label="Increment/Decrement\nsingle item in CBF"]

    start -> receive
    receive -> decision
    decision -> read [label="full rebuild"]
    decision -> increment [label="granular (Phase 2)"]
    read -> compute
    compute -> serialize
    increment -> serialize
    serialize -> publish
    publish -> end_node
}
```

### State Machine Diagram

```dot
digraph StateMachine {
    graph [label="State Machine: List Sync" labelloc=t fontsize=16 fontname="Arial" rankdir=LR]
    node [shape=box style="filled,rounded" fontname="Arial" fontsize=10 fillcolor="#B5FFFF"]
    edge [fontname="Arial" fontsize=9]

    start [label="" shape=circle fillcolor=black width=0.3]
    idle [label="Idle\n(CBF current)"]
    rebuilding [label="Rebuilding\n(computing CBF)"]
    syncing [label="Syncing\n(broadcasting CBF)"]
    stale [label="Stale\n(sync failed)"]

    start -> idle
    idle -> rebuilding [label="list.updated\nevent received"]
    rebuilding -> syncing [label="CBF computed"]
    syncing -> idle [label="bloom.sync\ndelivered"]
    syncing -> stale [label="NATS publish\nfailed"]
    stale -> rebuilding [label="retry timer\nor manual trigger"]
}
```

### Use Case Diagram

```dot
digraph UseCase {
    graph [label="Use Cases: Lists System" labelloc=t fontsize=16 fontname="Arial" rankdir=LR]
    node [fontname="Arial" fontsize=10]
    edge [fontname="Arial" fontsize=9]

    // Actors
    partner [label="Partner" shape=box fillcolor="#08427B" fontcolor=white style=filled]
    policy [label="Policy Engine" shape=box fillcolor="#999999" fontcolor=white style=filled]
    scoring [label="Scoring Engine" shape=box fillcolor="#999999" fontcolor=white style=filled]

    // Use cases
    subgraph cluster_system {
        label="Lists System" style=dashed
        node [shape=ellipse fillcolor="#B5FFFF" style=filled]

        uc1 [label="Manage Lists\n(CRUD)"]
        uc2 [label="Upload CSV"]
        uc3 [label="View Partner/\nPrivate Lists"]
        uc4 [label="Auto-add to\nPartner List"]
        uc5 [label="Check List\nMembership"]
    }

    partner -> uc1
    partner -> uc2
    partner -> uc3
    policy -> uc4
    scoring -> uc5
}
```

### Class Diagram (Go Adaptation)

```dot
digraph ClassDiagram {
    graph [label="Domain Model: List" labelloc=t fontsize=16 fontname="Arial" rankdir=TB]
    node [shape=record style=filled fillcolor="#B5FFFF" fontname="Courier" fontsize=9]
    edge [fontname="Arial" fontsize=9]

    List [label="{\<\<struct\>\>\nList|+ ID : string\n+ OrgSlug : string\n+ Kind : ListKind\n+ Name : string\n+ Summary : string\n+ CreatedAt : time.Time|+ Validate() error}"]

    Entry [label="{\<\<struct\>\>\nEntry|+ ID : string\n+ ListID : string\n+ Type : string\n+ Value : string\n+ CreatedAt : time.Time|}"]

    Repository [label="{\<\<interface\>\>\nRepository|+ Save(List) error\n+ Get(id string) (List, error)\n+ ListEntries(listID, cursor, limit) ([]Entry, error)\n+ AddEntry(Entry) error\n+ RemoveEntry(listID, value) error}" fillcolor="#FFFFB5"]

    List -> Entry [label="1..*" arrowhead=diamond]
    Repository -> List [label="manages" style=dashed]
    Repository -> Entry [label="manages" style=dashed]
}
```

---

## BPMN templates

Exclusive gateway (XOR, exactly one path), parallel gateway (AND, all paths), inclusive gateway (OR, one or more), event-based gateway (whichever event fires first). Swimlanes carry the roles.

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

## Flowchart templates

ISO 5807 shapes: oval terminator, rectangle process, diamond decision, parallelogram input/output, cylinder data store.

### Simple Linear Flow

```dot
digraph Flowchart {
    graph [label="Flow: <Process Name>" labelloc=t fontsize=16 fontname="Arial" rankdir=TB]
    node [shape=box style=filled fillcolor="#B5FFFF" fontname="Arial" fontsize=10]
    edge [fontname="Arial" fontsize=9]

    start [label="Start" shape=oval fillcolor="#C9E7B7"]
    step1 [label="Step 1:\nDo something"]
    step2 [label="Step 2:\nDo next thing"]
    step3 [label="Step 3:\nFinalize"]
    end_node [label="End" shape=oval fillcolor="#FFB5B5"]

    start -> step1 -> step2 -> step3 -> end_node
}
```

### Decision Flow

```dot
digraph Flowchart {
    graph [label="Flow: Request Validation" labelloc=t fontsize=16 fontname="Arial" rankdir=TB]
    node [shape=box style=filled fillcolor="#B5FFFF" fontname="Arial" fontsize=10]
    edge [fontname="Arial" fontsize=9]

    start [label="Receive\nRequest" shape=oval fillcolor="#C9E7B7"]
    validate [label="Validate\nInput" shape=diamond fillcolor="#FFFFB5"]
    auth [label="Check\nAuthorization" shape=diamond fillcolor="#FFFFB5"]
    process [label="Process\nRequest"]
    ok [label="Return 200" shape=oval fillcolor="#C9E7B7"]
    err400 [label="Return 400\nBad Request" shape=oval fillcolor="#FFB5B5"]
    err401 [label="Return 401\nUnauthorized" shape=oval fillcolor="#FFB5B5"]

    start -> validate
    validate -> auth [label="Valid"]
    validate -> err400 [label="Invalid"]
    auth -> process [label="Authorized"]
    auth -> err401 [label="Denied"]
    process -> ok
}
```

### Loop Flow

```dot
digraph Flowchart {
    graph [label="Flow: Process Batch Items" labelloc=t fontsize=16 fontname="Arial" rankdir=TB]
    node [shape=box style=filled fillcolor="#B5FFFF" fontname="Arial" fontsize=10]
    edge [fontname="Arial" fontsize=9]

    start [label="Start" shape=oval fillcolor="#C9E7B7"]
    init [label="Load batch\nfrom MongoDB"]
    check [label="More items?" shape=diamond fillcolor="#FFFFB5"]
    process [label="Process\ncurrent item"]
    advance [label="Move to\nnext item"]
    done [label="Return\nresults" shape=oval fillcolor="#C9E7B7"]

    start -> init -> check
    check -> process [label="Yes"]
    check -> done [label="No"]
    process -> advance -> check
}
```

### Multi-Way Decision

```dot
digraph Flowchart {
    graph [label="Flow: List Entry Type Routing" labelloc=t fontsize=16 fontname="Arial" rankdir=TB]
    node [shape=box style=filled fillcolor="#B5FFFF" fontname="Arial" fontsize=10]
    edge [fontname="Arial" fontsize=9]

    start [label="Receive\nlist entry" shape=oval fillcolor="#C9E7B7"]
    type_check [label="Entry\ntype?" shape=diamond fillcolor="#FFFFB5"]
    email [label="Normalize\nemail\n(lowercase, trim)"]
    ip [label="Validate\nIP format\n(v4 or v6)"]
    card [label="Hash\ncard token\n(SHA-256)"]
    merge [label="Add to CBF"]
    end_node [label="Done" shape=oval fillcolor="#C9E7B7"]

    start -> type_check
    type_check -> email [label="email"]
    type_check -> ip [label="ip"]
    type_check -> card [label="card"]
    email -> merge
    ip -> merge
    card -> merge
    merge -> end_node
}
```

### Error Handling Flow (Go Pattern)

```dot
digraph Flowchart {
    graph [label="Flow: Save List Entry (Go Error Handling)" labelloc=t fontsize=16 fontname="Arial" rankdir=TB]
    node [shape=box style=filled fillcolor="#B5FFFF" fontname="Arial" fontsize=10]
    edge [fontname="Arial" fontsize=9]

    start [label="SaveEntry\n(entry Entry)" shape=oval fillcolor="#C9E7B7"]
    validate [label="Validate\nentry fields"]
    val_err [label="err != nil?" shape=diamond fillcolor="#FFFFB5"]
    insert [label="repo.Insert\n(ctx, entry)"]
    ins_err [label="err != nil?" shape=diamond fillcolor="#FFFFB5"]
    publish [label="nats.Publish\n(list.updated)"]
    pub_err [label="err != nil?" shape=diamond fillcolor="#FFFFB5"]
    ok [label="return nil" shape=oval fillcolor="#C9E7B7"]
    ret_err [label="return\nfmt.Errorf(...)" shape=oval fillcolor="#FFB5B5"]

    start -> validate -> val_err
    val_err -> insert [label="nil"]
    val_err -> ret_err [label="err"]
    insert -> ins_err
    ins_err -> publish [label="nil"]
    ins_err -> ret_err [label="err"]
    publish -> pub_err
    pub_err -> ok [label="nil"]
    pub_err -> ret_err [label="err"]
}
```

### Parallel / Fork-Join Flow

```dot
digraph Flowchart {
    graph [label="Flow: CBF Sync Broadcast" labelloc=t fontsize=16 fontname="Arial" rankdir=TB]
    node [shape=box style=filled fillcolor="#B5FFFF" fontname="Arial" fontsize=10]
    edge [fontname="Arial" fontsize=9]

    start [label="CBF\ncomputed" shape=oval fillcolor="#C9E7B7"]
    serialize [label="Serialize CBF\nto binary"]
    fork [label="" shape=point width=0.1]
    est1 [label="Deliver to\nestimator-1"]
    est2 [label="Deliver to\nestimator-2"]
    est3 [label="Deliver to\nestimator-3"]
    join [label="" shape=point width=0.1]
    done [label="All instances\nupdated" shape=oval fillcolor="#C9E7B7"]

    start -> serialize -> fork
    fork -> est1
    fork -> est2
    fork -> est3
    est1 -> join
    est2 -> join
    est3 -> join
    join -> done
}
```
