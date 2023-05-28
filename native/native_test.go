package native

import (
	"math"
	"testing"
)

func initIsz() ([]uint, func(v *Native)) {
	mem := []uint{
		23, // opAddr
	}
	action := func(v *Native) {
		opAddr := 0
		v.mem[opAddr] = mask12(v.mem[opAddr] + 1)
		if v.mem[opAddr] == 0 {
			v.pc = mask12(v.pc + 1)
		}
	}
	return mem, action
}

func initLoopUntil() ([]uint, func(v *Native)) {
	mem := []uint{
		0, // sum
		0, // cnt
	}
	action := func(v *Native) {
		for v.mem[1] = 5000; v.mem[1] != 0; v.mem[1]-- {
			v.mem[0] += 1
		}
	}
	return mem, action
}

func initAnd() ([]uint, func(v *Native)) {
	mem := []uint{
		4503, // lac
		3003, // value
	}
	action := func(v *Native) {
		opAddr := 1
		v.mem[0] &= v.mem[opAddr] | 0o10000
	}
	return mem, action
}

func initTad() ([]uint, func(v *Native)) {
	mem := []uint{
		9,  // lac
		23, // value
	}
	action := func(v *Native) {
		opAddr := 1
		v.mem[0] = mask13(v.mem[0] + v.mem[opAddr])
	}
	return mem, action
}

func initSubleq() ([]uint, func(v *Native)) {
	mem := []uint{
		// loopuntil_v1 from subleq/fixtures/
		// program
		15,
		13,
		3,
		16,
		14,
		6,
		16,
		13,
		3,
		16,
		1000,
		12,
		0,
		0,
		0, // sum
		4999,
		math.MaxUint32, // -1
	}
	action := func(v *Native) {
		var hltVal uint
		pc := uint(0)
		for {
			operandA := v.mem[pc]
			operandB := v.mem[pc+1]
			operandC := v.mem[pc+2]
			v.mem[operandB] = mask32(v.mem[operandB] - v.mem[operandA])
			// If hlt location
			if operandB == 1000 {
				hltVal = v.mem[operandB]
				break
			}
			if v.mem[operandB] == 0 || v.mem[operandB] > math.MaxInt32 {
				pc = operandC
			} else {
				pc += 3
			}
		}
		if hltVal != 1 {
			panic("htlVal != 1")
		}
	}
	return mem, action
}

func initAdd12() ([]uint, func(v *Native)) {
	mem := []uint{
		4094, // a
		6,    // b
	}
	action := func(v *Native) {
		v.mem[1] = mask12(v.mem[0] + v.mem[1])
	}
	return mem, action
}

// Used by initJsr
// The go:noinline directive is there so that we actually test a subroutine rather than
// just having the routine inlined in the code
//
//go:noinline
func setVal(v *Native, n uint) {
	v.mem[0] = n
}

func initJsr() ([]uint, func(v *Native)) {
	mem := []uint{
		0, // val
	}
	action := func(v *Native) {
		setVal(v, 50)
	}
	return mem, action
}

func initSwitch() ([]uint, func(v *Native)) {
	mem := []uint{
		3, // lac
	}
	action := func(v *Native) {
		for i := 8; i != 0; i-- {
			switch i - 1 {
			case 0:
				v.mem[0] += 11
			case 1:
				v.mem[0] += 23
			case 2:
				v.mem[0] += 56
			case 3:
				v.mem[0] += 79
			case 4:
				v.mem[0] += 123
			case 5:
				v.mem[0] += 367
			case 6:
				v.mem[0] += 592
			case 7:
				v.mem[0] += 1001
			}
		}
	}
	return mem, action
}

var tests = []struct {
	name   string
	init   func() ([]uint, func(v *Native))
	want   map[uint]uint // [memloc]value
	wantPC uint
}{
	{"add12", initAdd12, map[uint]uint{1: 4}, 0},
	{"and", initAnd, map[uint]uint{0: 4499}, 0},
	{"tad", initTad, map[uint]uint{0: 32}, 0},
	{"isz", initIsz, map[uint]uint{0: 24}, 0},
	{"jsr", initJsr, map[uint]uint{0: 50}, 0},
	{"loopuntil", initLoopUntil, map[uint]uint{0: 5000}, 0},
	{"subleq", initSubleq, map[uint]uint{14: 5000}, 0},
	{"switch", initSwitch, map[uint]uint{0: 2255}, 0},
}

func TestNative(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(b *testing.T) {
			mem, action := test.init()
			v := New()
			v.LoadMem(mem)
			action(v)
			for memLoc, wantValue := range test.want {
				if v.mem[memLoc] != wantValue {
					t.Errorf("mem[%d] got: %d, want: %d", memLoc, v.mem[memLoc], wantValue)
				}
			}
			if v.pc != test.wantPC {
				t.Errorf("PC got: %d, want: %d", v.pc, test.wantPC)
			}
		})
	}
}

func BenchmarkNative(b *testing.B) {

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			b.StopTimer()
			mem, action := test.init()
			for n := 0; n < b.N; n++ {
				v := New()
				v.LoadMem(mem)
				b.StartTimer()
				action(v)
				b.StopTimer()
				for memLoc, wantValue := range test.want {
					if v.mem[memLoc] != wantValue {
						b.Errorf("mem[%d] got: %d, want: %d", memLoc, v.mem[memLoc], wantValue)
					}
				}
				if v.pc != test.wantPC {
					b.Errorf("PC got: %d, want: %d", v.pc, test.wantPC)
				}
			}
		})
	}
}
