package lazy

// GridCells determines how cells should be arranged in a grid's cross-axis.
type GridCells interface {
	// calculateCrossAxisCellCount calculates the number of columns (for vertical grid)
	// or rows (for horizontal grid) based on the available space.
	calculateCrossAxisCellCount(availableSpacePx int, spacingPx int) int
	// calculateCellSize returns the size in pixels for each cell given the count and available space
	calculateCellSize(availableSpacePx int, count int, spacingPx int) int
}

// gridCellsFixed implements a fixed number of columns/rows.
type gridCellsFixed struct {
	count int
}

// Fixed creates a GridCells that specifies a fixed number of columns (for LazyVerticalGrid)
// or rows (for LazyHorizontalGrid).
func Fixed(count int) GridCells {
	if count < 1 {
		count = 1
	}
	return &gridCellsFixed{count: count}
}

func (g *gridCellsFixed) calculateCrossAxisCellCount(availableSpacePx int, spacingPx int) int {
	return g.count
}

func (g *gridCellsFixed) calculateCellSize(availableSpacePx int, count int, spacingPx int) int {
	if count <= 0 {
		return availableSpacePx
	}
	// Total spacing between cells
	totalSpacing := spacingPx * (count - 1)
	usableSpace := availableSpacePx - totalSpacing
	if usableSpace < 0 {
		usableSpace = 0
	}
	return usableSpace / count
}

// gridCellsAdaptive implements adaptive sizing based on minimum cell size.
type gridCellsAdaptive struct {
	minSize int
}

// Adaptive creates a GridCells that calculates the number of columns/rows dynamically
// so that each cell is at least minSize dp wide/tall.
// The actual size will be larger to use all available space.
func Adaptive(minSize int) GridCells {
	if minSize <= 0 {
		minSize = 64 // Default minimum
	}
	return &gridCellsAdaptive{minSize: minSize}
}

func (g *gridCellsAdaptive) calculateCrossAxisCellCount(availableSpacePx int, spacingPx int) int {
	if availableSpacePx <= 0 {
		return 1
	}
	// Calculate how many cells of minSize can fit
	minSizePx := g.minSize // Already in dp, assuming 1:1 for now
	if minSizePx <= 0 {
		minSizePx = 64
	}
	// Account for spacing: n cells need (n-1) spacing gaps
	// Formula: n * minSize + (n-1) * spacing <= available
	// n * (minSize + spacing) <= available + spacing
	// n <= (available + spacing) / (minSize + spacing)
	count := (availableSpacePx + spacingPx) / (minSizePx + spacingPx)
	if count < 1 {
		count = 1
	}
	return count
}

func (g *gridCellsAdaptive) calculateCellSize(availableSpacePx int, count int, spacingPx int) int {
	if count <= 0 {
		return availableSpacePx
	}
	// Total spacing between cells
	totalSpacing := spacingPx * (count - 1)
	usableSpace := availableSpacePx - totalSpacing
	if usableSpace < 0 {
		usableSpace = 0
	}
	return usableSpace / count
}
