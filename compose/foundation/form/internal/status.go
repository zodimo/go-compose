package formengine

// Status is the validation/lifecycle state of a form node.
type Status string

const (
	// StatusValid indicates the node (and its subtree) has no validation errors.
	StatusValid Status = "VALID"
	// StatusInvalid indicates the node or one of its descendants has validation errors.
	StatusInvalid Status = "INVALID"
	// StatusPending is reserved for future asynchronous validation. It is never
	// emitted by any node in this change: async validators are follow-up work and
	// no node currently performs asynchronous validation. It exists so that
	// aggregation logic and callers can handle it when it is introduced.
	StatusPending Status = "PENDING"
	// StatusDisabled indicates the node is effectively disabled (own or inherited).
	StatusDisabled Status = "DISABLED"
)
