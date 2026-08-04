package render

// CallOp is an opaque handle to a recorded sequence of draw calls.
// It is produced by Backend.Record() and consumed by Backend.Apply().
//
// Macro discard semantics (EC2): a CallOp that is never applied to a backend
// produces nothing. This is guaranteed by design — the CallOp is just a
// handle; no side effects occur until Apply is called. Dropping a CallOp
// (e.g. the pointer phase's discarded recording) is safe and free.
type CallOp struct {
	backendID string // identifies the backend that created this
	seqID     uint64 // sequence index within the backend's recording store
	payload   any    // backend-specific opaque recorded data
}

// NewCallOp creates a CallOp with the given backend ID, sequence ID, and payload.
// Backend implementations use this to construct CallOp values from Record().
func NewCallOp(backendID string, seqID uint64, payload any) CallOp {
	return CallOp{backendID: backendID, seqID: seqID, payload: payload}
}

// CallOpPayload returns the backend-specific payload stored in a CallOp.
// Backend implementations use this in Apply() to extract recorded data.
func CallOpPayload(co CallOp) any {
	return co.payload
}

// CallOpBackendID returns the backend identifier of a CallOp.
// Backend implementations use this in Apply() to verify origin.
func CallOpBackendID(co CallOp) string {
	return co.backendID
}
