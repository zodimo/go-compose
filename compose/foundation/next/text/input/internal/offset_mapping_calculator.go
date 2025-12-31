package internal

import (
	"github.com/zodimo/go-compose/compose/ui/next/text"
)

// OffsetMappingCalculator builds bidirectional offset mappings after edit operations.
//
// This class tracks offset transformations between an original text and a transformed
// text that has had one or more edit operations applied. It supports mapping offsets
// in both directions:
//   - Forward: original → transformed
//   - Reverse: transformed → original
//
// When mappings are ambiguous (e.g., a point mapping to a deleted range), the result
// is returned as a TextRange indicating all valid mappings.
//
// This is a port of androidx.compose.foundation.text.input.internal.OffsetMappingCalculator.
type OffsetMappingCalculator struct {
	ops     []offsetOp
	opsSize int
}

// offsetOp represents a single edit operation for offset mapping.
type offsetOp struct {
	offset  int // Position where the edit occurs
	srcLen  int // Length in source text
	destLen int // Length in destination text
}

// NewOffsetMappingCalculator creates a new OffsetMappingCalculator.
func NewOffsetMappingCalculator() *OffsetMappingCalculator {
	return &OffsetMappingCalculator{
		ops: make([]offsetOp, 0, 8),
	}
}

// RecordEditOperation records a replacement from [sourceStart, sourceEnd) with newLength chars.
//
// Parameters:
//   - sourceStart: Start of the range being replaced (in source text)
//   - sourceEnd: End of the range being replaced (in source text)
//   - newLength: Length of the replacement text
func (o *OffsetMappingCalculator) RecordEditOperation(sourceStart, sourceEnd, newLength int) {
	srcLen := sourceEnd - sourceStart
	if srcLen == 0 && newLength == 0 {
		return // No-op
	}

	o.ops = append(o.ops, offsetOp{
		offset:  sourceStart,
		srcLen:  srcLen,
		destLen: newLength,
	})
	o.opsSize++
}

// MapFromSource maps an offset in the original string to the corresponding offset
// or range of offsets in the transformed string.
func (o *OffsetMappingCalculator) MapFromSource(offset int) text.TextRange {
	return o.mapOffset(offset, true)
}

// MapFromDest maps an offset in the transformed string to the corresponding offset
// or range of offsets in the original string.
func (o *OffsetMappingCalculator) MapFromDest(offset int) text.TextRange {
	return o.mapOffset(offset, false)
}

// mapOffset performs the actual offset mapping.
func (o *OffsetMappingCalculator) mapOffset(offset int, fromSource bool) text.TextRange {
	if o.opsSize == 0 {
		return text.NewTextRange(offset, offset)
	}

	currentOffset := offset
	resultStart := offset
	resultEnd := offset

	// Process operations in order for forward mapping, reverse for backward
	if fromSource {
		for i := 0; i < o.opsSize; i++ {
			op := o.ops[i]
			newStart, newEnd := o.mapStep(currentOffset, resultStart, resultEnd, op.offset, op.srcLen, op.destLen, fromSource)
			resultStart, resultEnd = newStart, newEnd
			currentOffset = resultStart // Track for subsequent operations
		}
	} else {
		for i := o.opsSize - 1; i >= 0; i-- {
			op := o.ops[i]
			newStart, newEnd := o.mapStep(currentOffset, resultStart, resultEnd, op.offset, op.destLen, op.srcLen, fromSource)
			resultStart, resultEnd = newStart, newEnd
			currentOffset = resultStart
		}
	}

	return text.NewTextRange(resultStart, resultEnd)
}

// mapStep applies one operation to the offset mapping.
func (o *OffsetMappingCalculator) mapStep(
	offset int,
	rangeStart int,
	rangeEnd int,
	opOffset int,
	untransformedLen int,
	transformedLen int,
	fromSource bool,
) (int, int) {
	opEnd := opOffset + untransformedLen

	if offset < opOffset {
		// Before the operation - no change
		return rangeStart, rangeEnd
	}

	if offset > opEnd {
		// After the operation - shift by delta
		delta := transformedLen - untransformedLen
		return rangeStart + delta, rangeEnd + delta
	}

	// Within the operation range
	if offset == opOffset {
		// At the start of the operation
		if untransformedLen == 0 {
			// Insertion: point maps to start of inserted range
			return opOffset, opOffset
		}
		return opOffset, opOffset
	}

	if offset == opEnd {
		// At the end of the operation
		return opOffset + transformedLen, opOffset + transformedLen
	}

	// Within the operation (ambiguous mapping)
	// Return the full transformed range
	return opOffset, opOffset + transformedLen
}

// Reset clears all recorded operations.
func (o *OffsetMappingCalculator) Reset() {
	o.ops = o.ops[:0]
	o.opsSize = 0
}
