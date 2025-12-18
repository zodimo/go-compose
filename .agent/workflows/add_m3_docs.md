---
description: Add reference documentation to Material 3 components
---

This workflow guides you through adding Material 3 reference links to component documentation.

## Prerequisites
- Verify `docs/m3_documentation_status.md` exists and lists the components.

## Steps

1. **Select Component**
   - Open `docs/m3_documentation_status.md`.
   - Find the first unchecked component (e.g., `[ ]`).
   - Note the component name and the M3 Web/Specs URLs.

2. **Navigate to Component Directory**
   - Go to `compose/foundation/material3/{component}`.
   - Example: For `button`, go to `compose/foundation/material3/button`.

3. **Locate or Create Documentation File**
   - Look for `doc.go` in the component directory.
   - If `doc.go` exists, open it.
   - If `doc.go` does not exist, create it with the following content (adjust package name):
     ```go
     /*
     Package {component} contains Material 3 {Component Name} components.
     */
     package {component}
     ```

4. **Add Reference Documentation**
   - Add the reference links to the package comment or the main component struct comment.
   - Format:
     ```go
     // Reference: [Material 3 {Component Name}]({M3 Web URL})
     // Specs: [Material 3 {Component Name} Specs]({M3 Specs URL})
     ```
   - Ensure the URLs match those provided in `docs/m3_documentation_status.md`.

5. **Verify URLs**
   - Double-check that the URLs are valid and reachable.
   - If a URL in the status doc is incorrect, update it with the correct one.

6. **Update Progress**
   - Mark the component as completed in `docs/m3_documentation_status.md`:
     - Change `[ ]` to `[x]`.

7. **Repeat**
   - Repeat the process for the next component in the list.
