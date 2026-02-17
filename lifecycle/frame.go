package lifecycle

type FrameLifecycleAware interface {
	StartFrame()
	EndFrame()
}
