---
description: Implement a new Material 3 component in GoCompose
---

# Component Implementation Workflow

This workflow guides you through implementing a new UI component in `go-compose-dev`.

## 1. Planning & Analysis
1. Check `docs/ROADMAP.md` to confirm the component's priority and requirements.
2. Check `docs/component_inventory.md` (if available) to see existing gap analysis.
3. Search the codebase for similar components or prior art (e.g., in `gio-mw` or other packages).
4. Create an `implementation_plan.md` artifact detailing:
   - Target package (usually `compose/foundation/material3/[component]`).
   - Planned API (Composables, Modifiers, Colors).
   - Demo verification plan.

## 2. Implementation
1. **Create Component Package**:
   - Create `compose/foundation/material3/[component]/[component].go`.
   - Define the main Composable function (e.g., `func MyComponent(...)`).
   - Implement logic using `gioui.org` primitives or `gio-mw` helpers.
2. **Define Styles & Defaults**:
   - Create `types.go` or `[component]_options.go` for options patterns.
   - Define standard modifications via `Modifier`.

## 3. Demo Creation
1. **Create Demo Directory**:
   - `mkdir -p cmd/demo/[component]`
2. **Implement UI**:
   - Create `cmd/demo/[component]/ui.go` showcasing variants (e.g., different sizes, states).
3. **Entry Point**:
   - Create `cmd/demo/[component]/main.go` to run the demo standalone.

## 4. Verification
1. **Manual Verification**:
   - Run `go run cmd/demo/[component]/main.go` and ensure it renders correctly.
2. **Screenshot Test**:
   - Create `cmd/demo/[component]/screenshot_test.go`.
   - Use `screenshot.TakeScreenshot` to capture the `UI()` state.
   - Run `go test ./cmd/demo/[component]/...` to generate the screenshot.
   - **Visual Verification**: This is for visual inspection by you (the agent), not a unit test per se. Use the resulting screenshot to verify that visual expectations are met.

## 5. Documentation & Finalization
1. **Update Roadmap**:
   - Mark the item as completed in `docs/ROADMAP.md`.
2. **Update Inventory**:
   - Update `docs/component_inventory.md` status.
3. **Walkthrough**:
   - Create/Update `walkthrough.md` with the generated screenshot.
4. **Cleanup**:
   - Run `go mod tidy` to ensure dependencies are clean.
