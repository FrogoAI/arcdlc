# Entity-Component-System (ECS) — Definition & Usage Guide

**Source**: Scott Bilas, "A Data-Driven Game Object System" (GDC 2002); Adam Martin, "Entity Systems are the Future of MMOG Development" (2007), which named the pattern and split systems out of components; Mick West, "Evolve Your Hierarchy" (Game Developer, 2007), on composition over inheritance for game objects.
**Adaptation**: Go-idiomatic — plain-struct components, an indexed registry for entities, stateless systems. Aligned with MDCA (`mdca.md`, Appendix B, which names the ECS-shaped `internal/<subsystem>/` package as the Go **client**'s domain-module unit), `Go Client.md`, `ddd.md`, and `solid.md`.

---

## What is ECS?

ECS is an architectural pattern that separates data from behavior:

- **Entity** = unique identifier (uint64). Has no data, no logic. Just an ID that links components.
- **Component** = pure data struct. No methods, no callbacks, no logic. Only fields.
- **System** = stateless function that reads/writes components. All logic lives here.

Key principle: **composition over inheritance**. An entity's capabilities are defined by which components it has, not by a class hierarchy. Add a `Draggable` component → entity becomes draggable. Remove it → it stops being draggable. No class changes needed.

## Rules

These five blocks, together with the anti-pattern table below, are the conformance criteria of this
document — the only part of it an audit may file a gap against. Each block carries a stable identifier
(`ECS-C1`, `ECS-C2` for components; `ECS-S1`, `ECS-S2` for systems; `ECS-E1` for entities); cite it in a gap
block rather than the block's title. Identifiers are never reused or renumbered; new rules are appended with
the next free number inside their letter group.

### `ECS-C1` — Components MUST:
- Be plain structs with exported fields
- Contain only data (primitives, pointers, slices)
- Have single responsibility (Transform = position only, Sprite = image only)
- Be reusable across different entity types

### `ECS-C2` — Components MUST NOT:
- Contain methods with logic
- Hold function callbacks
- Reference other components directly
- Import system packages

### `ECS-S1` — Systems MUST:
- Be stateless — all state lives in components
- Have single responsibility (one system = one behavior)
- Operate on components by querying the registry
- Be independent of each other (no system calls another system)

### `ECS-S2` — Systems MUST NOT:
- Store game state in struct fields
- Know about specific entity types
- Call other systems directly
- Access raw OS or engine APIs (go through the scene/context accessor)

### `ECS-E1` — Entities MUST:
- Be just an ID (uint64) in the registry
- Be composed of components, never subclassed
- Belong to a group (archetype) that defines their component set

## Anti-Patterns

Each row is auditable in any ECS codebase; the `Rule` column names the identifier a gap block should cite.

| Don't | Do | Rule |
|-------|-----|------|
| Logic in components | Logic in systems | `ECS-C1`, `ECS-C2` |
| State in system struct fields | State in components, or behind the scene/context accessor | `ECS-S1`, `ECS-S2` |
| Coordinate/scaling math written inline (`x * scale + offset`) at each call site | One shared coordinate helper that every system calls | `ECS-S1` |
| A system reading input straight from the engine or OS API | Read it through the scene/context accessor interface | `ECS-S2` |
| A system assigning to its own state field | State lives in a component; the state system owns transitions | `ECS-S1`, `ECS-S2` |
| System calls another system | Systems communicate through components | `ECS-S1` |
| One massive `Update()` function | Multiple focused systems | `ECS-S1` |

## Reference Implementation (example only — not a conformance requirement)

Everything below this heading describes **one specific Ebiten/Go game codebase** — its component set, system
order, scene accessor, helper functions, and package names. It is illustrative only: **an audit MUST NOT file
a gap against a project for lacking these types, packages, groups, or function names.** The conformance
criteria of this document are `ECS-C1`–`ECS-E1` and the anti-pattern table above; the sections below merely
show one way to satisfy them.

### Registry as Indexed ECS
```
Registry[group string, id uint64, value any]
```
- **group** = archetype name ("cursor", "tray_piece", "button")
- **id** = entity identifier
- **value** = `*components.Entity` with component pointers

### Entity Container
```go
type Entity struct {
    Transform   *Transform    // position, size
    Sprite      *Sprite       // image reference
    Renderable  *Renderable   // rotation, alpha, z-order
    Button      *Button       // label, colors
    Clickable   *Clickable    // RTree spatial
    Scrollable  *Scrollable   // scroll offset
    TextBlock   *TextBlock    // multi-line text
    Avatar      *Avatar       // character portrait
    Cursor      *Cursor       // virtual cursor state
    State       *State        // game state machine
    Progress    *ProgressBar  // loading progress
    PuzzleGrid  *PuzzleGrid   // puzzle workspace
    PuzzlePiece *PuzzlePiece  // tray piece link
    Draggable   *Draggable    // drag-and-drop flag
    ImageLayer  *ImageLayerComp // TMX image layer
    TileLayer   *TileLayerComp  // TMX tile layer
}
```
Components are pointers — nil means entity doesn't have that component.

### System Execution Order
Systems register in `BaseScene.Systems` and run in order every frame:
```
1. StateSystem    — state machine, loading, transitions
2. CursorSystem   — virtual cursor from mouse/touch
3. InputSystem    — RTree hit testing, fires button clicks
4. ScrollSystem   — mouse wheel / swipe scrolling
5. DragDropSystem — piece pickup, rotation, placement
6. CameraSystem   — zoom keys / pinch gesture
7. RenderSystem   — ALL drawing
```
Each system has `Update(ctx)` (logic) and `Draw(ctx, screen)` (rendering).

### SceneAccessor Interface
Systems access scene data through an interface, never a concrete type:
```go
type SceneAccessor interface {
    GetRegistry() *Registry
    GetCamera() *Camera
    GetInputSystem() *InputSystem
    GetTileMap() *tilemap.Map
    GetCases() []*CaseConfig
    GetSelectedCase() int
    SetSelectedCase(int)
    GetCursorPos() (int, int)
    SetCursorPos(x, y int)
    GetCasesScroll() int
    SetCasesScroll(int)
    // ... etc
}
```

### Coords Helper
All map-to-screen coordinate conversion goes through `Coords`:
```go
co := CoordsFromScene(scene)
screenX := co.MapToScreenX(mapX)     // mapX * scale + offsetX
screenY := co.MapToScreenY(mapY)     // mapY * scale + offsetY
mapX := co.ScreenToMapX(screenX)     // (screenX - offsetX) / scale
pixels := co.MapToScreenSize(mapUnits) // mapUnits * scale
sx, sy, sw, sh := co.MapRect(x, y, w, h) // full rectangle
```
**Never** write `obj.X * scaleX + offsetX` directly.

### Assemblage Functions
Factories that create entity groups for a state:
```go
func CreateAppLayoutEntities(reg, cases, selectedCase) {
    CleanStateGroups(reg)  // remove old entities
    reg.Add("case_list", 0, &Entity{Scrollable: &Scrollable{...}})
    reg.Add("avatar", 0, &Entity{Avatar: &Avatar{...}})
    // ...
}
```
Called by StateSystem on state transitions.

### Adding a New Feature

1. **New data needed?** → New component in `components/`
2. **New behavior needed?** → New system in `fsystems/`
3. **New entity type?** → New group constant in `components/groups.go`
4. **Entities created/destroyed?** → Assemblage function in `entities/`
5. **Scene data access?** → New method on `SceneAccessor` + implement on GameScene
6. **Register system** in `GameScene.Init()` at correct execution position

### Anti-Patterns, in That Codebase's Own Symbols

The project-neutral rules above, spelled out with this codebase's names. Do not audit another project
against these symbols.

| Don't | Do | Rule |
|-------|-----|------|
| `obj.X * scaleX + offsetX` | `co.MapToScreenX(obj.X)` | `ECS-S1` |
| `ebiten.CursorPosition()` | `scene.GetCursorPos()` | `ECS-S2` |
| `s.state = StateXxx` | `s.setECSState(c.StateXxx)` | `ECS-S1`, `ECS-S2` |
