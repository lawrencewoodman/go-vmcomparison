// Generated test file by main_test.go

package codegen

func inittad_v1() ([]uint, []func(*CGVM)) {
	const (
		p_done = 4
	)
	const (
		m_lac = 3
		m_mask13 = 2
		m_memBase = 0
		m_ok = 4
		m_opAddr = 1
		m_val = 5
	)
	program := []func(v *CGVM){
		func(v *CGVM) { op_LDA(v, calcBaseIndexAddr(v, m_memBase, m_opAddr)) },
		func(v *CGVM) { op_ADD(v, m_lac) },
		func(v *CGVM) { op_AND(v, m_mask13) },
		func(v *CGVM) { op_STA(v, m_lac) },
		func(v *CGVM) { op_HLT(v, m_ok) },
	}
	memory := []uint{
		0,
		5,
		8191,
		9,
		0,
		23,
	}
	return memory, program
}

func init() {
	addTest("tad_v1.asm", inittad_v1, map[uint]uint{3: 32,})
}
