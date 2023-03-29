package vmstack

import (
	"testing"
)

const (
	HLT     = 0 << 24
	FETCH   = 1 << 24
	STORE   = 2 << 24
	ADD     = 3 << 24
	STORE13 = 9 << 24
	DJNZ    = 11 << 24
	STORE12 = 14 << 24
	LITO    = 15 << 24
	FETCHO  = 16 << 24
	DJNZO   = 17 << 24
	DROP    = 18 << 24
	SWAP    = 19 << 24
	FETCHBI = 20 << 24
	STOREO  = 21 << 24
	DSZO    = 22 << 24
	JMPO    = 23 << 24
	ADDBI   = 24 << 24
	R_PUSH  = 25 << 24
	R_POP   = 26 << 24
	FETCHI  = 27 << 24
)

var routineStack_TADv1 = []uint{
	LITO + 13, // membase
	FETCH,
	LITO + 14, // opAddr
	FETCH,
	ADD,
	FETCH,
	LITO + 15, // lac
	FETCH,
	ADD,
	LITO + 15, // lac - Would DUP or something similar be handy?
	STORE13,
	LITO + 16, // ok
	HLT,
	13, // membase
	4,  // opAddr
	9,  // lac
	0,  // ok
	23, // val
}

var routineStack_TADv2 = []uint{
	FETCHO + 10, // membase
	FETCHO + 11, // opaddr
	ADD,
	FETCH,
	FETCHO + 12, // lac
	ADD,
	LITO + 12, // lac
	STORE13,
	LITO + 13, // ok
	HLT,
	10, // membase
	4,  // opAddr
	9,  // lac
	0,  // ok
	23, // val
}

var routineStack_TADv3 = []uint{
	FETCHO + 9,  // membase
	FETCHO + 10, // opaddr
	FETCHBI,
	FETCHO + 11, // lac
	ADD,
	LITO + 11, // lac
	STORE13,
	LITO + 0, // ok
	HLT,
	9,  // membase
	4,  // opAddr
	9,  // lac
	0,  // ok
	23, // val
}

/*
	; Version 2
	; LOOP UNTIL
	LDA     l150
	STA     cnt
	LDA     lac

loop:    ADD II  memBase,opAddr

	DSZ     cnt
	JMP     loop
	STA12   lac

done:    HLT     ok
memBase: 8
opAddr:  4
lac:     9
ok:      0
val:     23
cnt:     0
l150:    150
*/
var routineStack_loopUntilv1 = []uint{
	FETCHO + 16, // lac
	LITO + 150,  // (sum -- sum cnt)
	// loop:
	SWAP,        // (sum cnt -- cnt sum)
	FETCHO + 14, // memBase
	FETCHO + 15, // opAddr
	FETCHBI,     // (cnt sum base index -- cnt sum val)
	ADD,         // (cnt sum val -- cnt sum)
	SWAP,        // (cnt sum -- sum cnt)
	DJNZO + 2,   // loop
	DROP,        // (sum cnt -- sum)
	LITO + 16,   // lac
	STORE12,
	LITO + 0,
	HLT,
	14, // memBase
	3,  // opAddr
	9,  // lac
	23, // val
}

var routineStack_loopUntilv2 = []uint{
	FETCHO + 15, // lac
	LITO + 150,  // (sum -- sum cnt)
	STOREO + 17, // cnt
	// loop:
	FETCHO + 13, // memBase
	FETCHO + 14, // opAddr
	FETCHBI,     // (sum base index -- sum val)
	ADD,         // (sum val -- sum)
	DSZO + 17,   // cnt
	JMPO + 3,
	LITO + 15, // lac
	STORE12,
	LITO + 0,
	HLT,
	13, // memBase
	3,  // opAddr
	9,  // lac
	23, // val
	0,  // cnt
}

var routineStack_loopUntilv3 = []uint{
	FETCHO + 14, // lac
	LITO + 150,  // (sum -- sum cnt)
	STOREO + 16, // cnt
	// loop:
	FETCHO + 12, // membase
	FETCHO + 13, // opaddr
	ADDBI,
	DSZO + 16, // cnt
	JMPO + 3,
	LITO + 14, // lac
	STORE12,
	LITO + 0,
	HLT,
	12, // memBase
	3,  // opAddr
	9,  // lac
	23, // val
	0,  // cnt
}

var routineStack_loopUntilv4 = []uint{
	LITO + 22, // memBase
	FETCH,
	LITO + 23, // opAddr
	FETCH,
	ADD,
	LITO + 24, // memloc
	STORE,
	LITO + 25, // lac
	FETCH,
	LITO + 150, // (sum -- sum cnt)
	// loop:
	R_PUSH,
	LITO + 24, // memloc
	FETCHI,
	ADD, // (sum val -- sum)
	R_POP,
	LITO + 10, // loop
	DJNZ,
	DROP,      // (sum cnt -- sum)
	LITO + 25, // lac
	STORE12,
	LITO + 0,
	HLT,
	22, // memBase
	4,  // opAddr
	0,  // memloc
	9,  // lac
	23, // val
}

var tests = []struct {
	name    string
	routine []uint
	want    map[uint]uint // [memloc]value
}{
	{"tad_v1", routineStack_TADv1, map[uint]uint{15: 32}},
	{"tad_v2", routineStack_TADv2, map[uint]uint{12: 32}},
	{"tad_v3", routineStack_TADv3, map[uint]uint{11: 32}},
	{"loopuntil_v1", routineStack_loopUntilv1, map[uint]uint{16: 3459}},
	{"loopuntil_v2", routineStack_loopUntilv2, map[uint]uint{15: 3459}},
	{"loopuntil_v3", routineStack_loopUntilv3, map[uint]uint{14: 3459}},
	{"loopuntil_v4", routineStack_loopUntilv4, map[uint]uint{25: 3459}},
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
	}
}
