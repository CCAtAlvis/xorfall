package grid

import (
	"sort"

	"github.com/CCAtAlvis/xorfall/src/configs"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Row struct {
	bits uint8
}

func (r *Row) IsBitSet(col int) bool {
	return r.bits&(1<<col) != 0
}

func (r *Row) SetBit(col int) {
	r.bits |= (1 << col)
}

func (r *Row) ClearBit(col int) {
	r.bits &^= (1 << col)
}

func (r *Row) ApplyMask(mask Mask, startCol int) {
	newBits := r.PreviewMaskResult(mask, startCol)
	r.bits = newBits
}

// PreviewMaskResult returns the row bits as they would be after applying the mask (non-mutating).
func (r *Row) PreviewMaskResult(mask Mask, startCol int) uint8 {
	bits := r.bits
	for i := range mask.Length {
		col := startCol + int(i)
		maskBit := mask.MaskShape&(1<<i) != 0
		rowBit := bits&(1<<col) != 0

		var op bool
		switch mask.MaskType {
		case MaskTypeOR:
			op = maskBit || rowBit
		case MaskTypeXOR:
			op = maskBit != rowBit
		case MaskTypeNOT:
			op = !rowBit
		case MaskTypeAND:
			op = maskBit && rowBit
		case MaskTypeXNOR:
			op = maskBit == rowBit
		}

		if op {
			bits |= 1 << col
		} else {
			bits &^= 1 << col
		}
	}
	return bits
}

// GenerateRow produces a row based on current phase (from configs.GameState().GetPhase()).
// Rows are generated only at game start and when a row is cleared and replaced at the bottom.
func GenerateRow(phase configs.Phase) Row {
	switch phase {
	case configs.Phase1Learning:
		return generatePhase1Row()
	case configs.Phase2SkillBuilding:
		return generatePhase2Row()
	default:
		return generatePhase3Row()
	}
}

// Phase 1: Start from all 0s or all 1s; flip exactly 1 or 2 non-adjacent bits.
func generatePhase1Row() Row {
	base := uint8(0)
	if rl.GetRandomValue(0, 1) == 1 {
		base = 255
	}
	flipCount := 1
	if rl.GetRandomValue(0, 1) == 1 {
		flipCount = 2
	}
	positions := pickNonAdjacentPositions(flipCount)
	flipMask := uint8(0)
	for _, p := range positions {
		flipMask |= 1 << p
	}
	return Row{bits: base ^ flipMask}
}

func pickNonAdjacentPositions(n int) []int {
	if n == 1 {
		return []int{int(rl.GetRandomValue(0, 7))}
	}
	// Pick two positions with |i-j| > 1
	for {
		a := rl.GetRandomValue(0, 7)
		b := rl.GetRandomValue(0, 7)
		if a > b {
			a, b = b, a
		}
		if b-a > 1 {
			return []int{int(a), int(b)}
		}
	}
}

// Phase 2: Start from solved row; flip 2–4 bits; no run of 3+ consecutive flipped bits.
func generatePhase2Row() Row {
	base := uint8(0)
	if rl.GetRandomValue(0, 1) == 1 {
		base = 255
	}
	flipCount := 2 + int(rl.GetRandomValue(0, 2)) // 2, 3, or 4
	positions := pickPositionsNoRun3(flipCount)
	flipMask := uint8(0)
	for _, p := range positions {
		flipMask |= 1 << p
	}
	return Row{bits: base ^ flipMask}
}

func pickPositionsNoRun3(n int) []int {
	positions := make([]int, n)
	for {
		used := make(map[int]bool)
		for i := 0; i < n; i++ {
			for {
				p := rl.GetRandomValue(0, 7)
				if !used[int(p)] {
					used[int(p)] = true
					positions[i] = int(p)
					break
				}
			}
		}
		slices := make([]int, 0, n)
		for p := range used {
			slices = append(slices, p)
		}
		sort.Ints(slices)
		if !hasRunOf3Consecutive(slices) {
			return slices
		}
	}
}

func hasRunOf3Consecutive(sorted []int) bool {
	if len(sorted) < 3 {
		return false
	}
	run := 1
	for i := 1; i < len(sorted); i++ {
		if sorted[i] == sorted[i-1]+1 {
			run++
			if run >= 3 {
				return true
			}
		} else {
			run = 1
		}
	}
	return false
}

// Phase 3: Fully random; reject if already solved (all 0s or all 1s).
func generatePhase3Row() Row {
	for {
		bits := uint8(rl.GetRandomValue(0, 255))
		if bits != 0 && bits != 255 {
			return Row{bits: bits}
		}
	}
}

type Board struct {
	rows            []Row
	currentRowIndex int

	validRowStates   []uint8
	lastRowSpawnTime float32
	rowSpawnInterval float32
}
