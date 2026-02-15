package lifecycle

type FrameLifecycleAwarePersistentState interface {
	StartFrame()
	EndFrame()
}
