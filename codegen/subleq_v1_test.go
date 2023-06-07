// Generated test file by main_test.go

package codegen

func initsubleq_v1() ([]uint, []func(*CGVM)) {
	const (
		p_exec = 9
		p_fetch = 0
		p_halt = 20
		p_jmpC = 17
	)
	const (
		m_hltVal = 2
		m_l1000 = 0
		m_memBase = 3
		m_ok = 7
		m_opA = 4
		m_opB = 5
		m_opC = 6
		m_pc = 1
		m_program = 8
		m_sum = 22
	)
	program := []func(v *CGVM){
		func(v *CGVM) { op_LDA(v, calcBaseIndexAddr(v, m_memBase, m_pc)) },
		func(v *CGVM) { op_STA(v, m_opA) },
		func(v *CGVM) { op_INC(v, m_pc) },
		func(v *CGVM) { op_LDA(v, calcBaseIndexAddr(v, m_memBase, m_pc)) },
		func(v *CGVM) { op_STA(v, m_opB) },
		func(v *CGVM) { op_INC(v, m_pc) },
		func(v *CGVM) { op_LDA(v, calcBaseIndexAddr(v, m_memBase, m_pc)) },
		func(v *CGVM) { op_STA(v, m_opC) },
		func(v *CGVM) { op_INC(v, m_pc) },
		func(v *CGVM) { op_LDA(v, calcBaseIndexAddr(v, m_memBase, m_opB)) },
		func(v *CGVM) { op_SUB(v, calcBaseIndexAddr(v, m_memBase, m_opA)) },
		func(v *CGVM) { op_STA(v, calcBaseIndexAddr(v, m_memBase, m_opB)) },
		func(v *CGVM) { op_LDA(v, m_l1000) },
		func(v *CGVM) { op_SUB(v, m_opB) },
		func(v *CGVM) { op_JEQ(v, p_halt) },
		func(v *CGVM) { op_LDA(v, calcBaseIndexAddr(v, m_memBase, m_opB)) },
		func(v *CGVM) { op_JGT(v, p_fetch) },
		func(v *CGVM) { op_LDA(v, m_opC) },
		func(v *CGVM) { op_STA(v, m_pc) },
		func(v *CGVM) { op_JMP(v, p_fetch) },
		func(v *CGVM) { op_LDA(v, calcBaseIndexAddr(v, m_memBase, m_opB)) },
		func(v *CGVM) { op_STA(v, m_hltVal) },
		func(v *CGVM) { op_HLT(v, m_ok) },
	}
	memory := []uint{
		1000,
		0,
		0,
		m_program,
		0,
		0,
		0,
		0,
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
		0,
		4999,
		18446744073709551615,
	}
	return memory, program
}

func init() {
	addTest("subleq_v1.asm", initsubleq_v1, map[uint]uint{22: 5000,})
}
