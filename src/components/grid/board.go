package grid

import rl "github.com/gen2brain/raylib-go/raylib"

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
	for i := range mask.Length {
		maskBit := mask.MaskShape&(1<<i) != 0
		rowBit := r.IsBitSet(startCol + int(i))

		var op bool
		switch mask.MaskType {
		case MaskTypeOR:
			op = maskBit || rowBit
		case MaskTypeXOR:
			op = maskBit != rowBit
		case MaskTypeNOT:
			op = !maskBit
			// case MaskTypeAND:
			// 	op = maskBit && rowBit
		}

		if op {
			r.SetBit(startCol + i)
		} else {
			r.ClearBit(startCol + i)
		}
	}
}

func GenerateRandomRow() Row {
	return Row{bits: uint8(rl.GetRandomValue(0, 255))}
}

type Board struct {
	rows            []Row
	currentRowIndex int

	validRowStates          []uint8
	lastRowSpawnTime        float32
	currentRowSpawnInterval float32
}
