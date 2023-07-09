// Generated test file by main_test.go

package codegen

func initloopuntil_v1() ([]int64, []func(*CGVM) error) {
	const (
		p_done = 7
		p_loop = 3
	)
	const (
		m_cnt = 2
		m_l0 = 3
		m_l1 = 4
		m_l5000 = 5
		m_ok = 1
		m_sum = 0
	)
	program := []func(v *CGVM) error {
		func(v *CGVM) error { return op_LDA(v, m_l5000) },
		func(v *CGVM) error { return op_STA(v, m_cnt) },
		func(v *CGVM) error { return op_LDA(v, m_l0) },
		func(v *CGVM) error { return op_ADD(v, m_l1) },
		func(v *CGVM) error { return op_DSZ(v, m_cnt) },
		func(v *CGVM) error { return op_JMP(v, p_loop) },
		func(v *CGVM) error { return op_STA(v, m_sum) },
		func(v *CGVM) error { return op_HLT(v, m_ok) },
	}
	memory := []int64{
		0,
		0,
		0,
		0,
		1,
		5000,
	}
	return memory, program
}

func init() {
	addTest("loopuntil_v1.asm", initloopuntil_v1, map[int64]int64{0: 5000,})
}
