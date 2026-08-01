# Diagram Conventions

**Purpose**: The single home for how initiative diagrams are produced, named, and colored. Applies to every notation in this library — C4, UML, BPMN, flowcharts, and the diagrams embedded in an arc42 document.

Each notation reference (`C4.md`, `UML.md`, `BPMN.md`, `Flowchart.md`, `Arc42.md`) points here instead of repeating these rules. Notation-specific content — symbols, templates, selection guidance — stays in those files.

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

**Applies to C4 diagrams only** (`C4.md`: context, container, component, code, dynamic, deployment, landscape). Color encodes the *role* of the element.

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

**Applies to UML (`UML.md`), BPMN (`BPMN.md`), and flowchart (`Flowchart.md`) diagrams.** Color encodes the *kind of step*.

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
- One deliberate crossover: UML use-case and deployment diagrams color external actors with the C4 person and external-system hexes (`#08427B`, `#999999`) so actors read the same across notations — see the DOT templates in `UML.md`.
- ArchiMate diagrams use a third, layer-based palette that stays with its notation in `TOGAF.md` (`## ArchiMate Diagram Conventions`).
