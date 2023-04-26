package vmstack2

import (
	"fmt"
	"testing"
)

const (
	HLT     = 0 << 24
	FETCH   = 1 << 24
	STORE   = 2 << 24
	ADD     = 3 << 24
	AND     = 5 << 24
	STORE13 = 9 << 24
	DJNZ    = 11 << 24
	STORE12 = 14 << 24
	LIT     = 15 << 24
	DROP    = 18 << 24
	SWAP    = 19 << 24
	FETCHBI = 20 << 24
	ADDBI   = 24 << 24
	FETCHI  = 27 << 24
	JSR     = 28 << 24
	RET     = 29 << 24
)

var routineStack_TADv1 = []uint{
	FETCH + 8, // membase
	FETCH + 9, // opAddr
	ADD,
	FETCH,
	FETCH + 10, // lac
	ADD,
	STORE13 + 10, // lac
	HLT + 1,      // ok
	8,            // membase
	3,            // opAddr
	9,            // lac
	23,           // val
}

var routineStack_TADv2 = []uint{
	FETCH + 7, // membase
	FETCH + 8, // opaddr
	FETCHBI,
	FETCH + 9, // lac
	ADD,
	STORE13 + 9, // lac
	HLT + 1,     // ok
	7,           // membase
	3,           // opAddr
	9,           // lac
	23,          // val
}

var routineStack_TADv3 = []uint{
	FETCH + 9,  // membase
	FETCH + 10, // opAddr
	ADD,
	FETCH,
	FETCH + 11, // lac
	ADD,
	AND + 8191, // 13-bit mask
	STORE + 11, // lac
	HLT + 1,    // ok
	9,          // membase
	3,          // opAddr
	9,          // lac
	23,         // val
}

var routineStack_JSRv1 = []uint{
	LIT + 50, // value to set
	JSR + 3,  // setVal
	HLT + 1,  // ok
	// setVal:
	STORE + 5, // val
	RET,
	0, // val
}

var routineStack_loopUntilv1 = []uint{
	LIT + 0,    // (-- sum)
	LIT + 5000, // (sum -- sum cnt)
	// loop:
	SWAP,      // (sum cnt -- cnt sum)
	ADD + 1,   //
	SWAP,      // (cnt sum -- sum cnt)
	DJNZ + 2,  // loop
	DROP,      // (sum cnt -- sum)
	STORE + 9, // sum
	HLT + 1,   // ok
	0,         // sum
}

var routineStack_loopUntilv2 = []uint{
	LIT + 0,    // (-- sum)
	STORE + 8,  // sum
	LIT + 5000, // (-- cnt)
	// loop:
	FETCH + 8, // sum
	ADD + 1,   //
	STORE + 8, // sum
	DJNZ + 3,  // loop (cnt -- cnt)
	HLT + 1,   // ok
	0,         // sum
}

var tests = []struct {
	name    string
	routine []uint
	want    map[uint]uint // [memloc]value
}{
	{"tad_v1", routineStack_TADv1, map[uint]uint{10: 32}},
	{"tad_v2", routineStack_TADv2, map[uint]uint{9: 32}},
	{"tad_v3", routineStack_TADv3, map[uint]uint{11: 32}},
	{"jsr_v1", routineStack_JSRv1, map[uint]uint{5: 50}},
	{"loopuntil_v1", routineStack_loopUntilv1, map[uint]uint{9: 5000}},
	{"loopuntil_v2", routineStack_loopUntilv2, map[uint]uint{8: 5000}},
}

func TestRun(t *testing.T) {
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := New()
			v.LoadRoutine(test.routine)
			_, err := v.Run()
			if err != nil {
				t.Errorf("Run() err: %v", err)
			}
			for memLoc, wantValue := range test.want {
				if v.mem[memLoc] != wantValue {
					t.Errorf("mem[%d] got: %d, want: %d", memLoc, v.mem[memLoc], wantValue)
				}
			}
		})
	}
}

func BenchmarkRun(b *testing.B) {
	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			b.StopTimer()

			for n := 0; n < b.N; n++ {
				v := New()
				v.LoadRoutine(test.routine)

				b.StartTimer()
				_, err := v.Run()
				b.StopTimer()

				if err != nil {
					b.Errorf("Run() err: %v", err)
				}
				for memLoc, wantValue := range test.want {
					if v.mem[memLoc] != wantValue {
						b.Errorf("mem[%d] got: %d, want: %d", memLoc, v.mem[memLoc], wantValue)
					}
				}
			}
		})
		fmt.Printf("Routine: %s size: %d\n", test.name, len(test.routine))
	}
}
