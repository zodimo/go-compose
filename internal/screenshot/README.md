# Screenshot of Recorded Ops

Example from gio tests

```go
import 	(
    "gioui.org/gpu/headless"
)

width := lineWidth
height := 100
cap := image.NewRGBA(image.Rect(0, 0, width, height))
w, _ := headless.NewWindow(width, height)
defer w.Release()

ops := new(op.Ops)
gtx := layout.Context{
    Constraints: layout.Constraints{Max: image.Pt(width, height)},
    Ops:         ops,
}

w.Frame(ops)
w.Screenshot(cap)
b := new(bytes.Buffer)
_ = png.Encode(b, cap)
screenshotName := tc.name + ".png"
_ = os.WriteFile(screenshotName, b.Bytes(), 0o644)
t.Logf("wrote %q", screenshotName)
```

# Create a Recording
```go

import 	(
    "gioui.org/op"
)
//gtx layout.Context
macro := op.Record(gtx.Ops)
// call node.Layout(gtx) or similar
callOp := macro.Stop()
```