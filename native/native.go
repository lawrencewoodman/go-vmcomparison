/*
 * Routines implemented natively for benchmark comparison
 */

package native

// TODO: Make this configurable
const memSize = 32000

type Native struct {
	mem [memSize]uint // Memory
	pc  uint          // Program Counter
}

func New() *Native {
	return &Native{}
}

func (v *Native) LoadMem(mem []uint) {
	copy(v.mem[:], mem)
}

// Returns the lower 12-bits
func mask12(w uint) uint {
	return w & 0o7777
}

// Returns the lower 13-bits
func mask13(w uint) uint {
	return w & 0o17777
}

func mask32(n uint) uint {
	return n & 0xFFFFFFFF
}
