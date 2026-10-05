package formengine

// ValidatorFunc validates a control value, returning a non-nil error when the
// value is invalid. Validators must not depend on the control's internal lock:
// they run on a snapshot of the value with no node lock held (design D5).
//
// The engine is validator-agnostic: it stores and joins whatever errors a
// validator returns. The public fform package and its rules subpackage build the
// rule catalog on top of this type; the engine ships no rules of its own.
type ValidatorFunc[T any] func(value T) error
