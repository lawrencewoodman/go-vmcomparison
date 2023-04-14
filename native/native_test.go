package native

import (
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
		9,  // lac
		0,  // cnt
		23, // value
	}
	action := func(v *Native) {
		for v.mem[1] = 150; v.mem[1] != 0; v.mem[1] = mask12(v.mem[1] - 1) {
			v.mem[0] = mask12(v.mem[0] + v.mem[2])
		}
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
	{"tad", initTad, map[uint]uint{0: 32}, 0},
	{"isz", initIsz, map[uint]uint{0: 24}, 0},
	{"jsr", initJsr, map[uint]uint{0: 50}, 0},
	{"loopuntil", initLoopUntil, map[uint]uint{0: 3459}, 0},
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
