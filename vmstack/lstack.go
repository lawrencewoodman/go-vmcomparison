/*
 * A limited stack - It will keep pushing and popping but won't error of it
 * reaches the beginning or end of the stack instead it will just keep
 * acting on the same position in the stack.
 *
 * Copyright (C) 2023 Lawrence Woodman <lwoodman@vlifesystems.com>
 *
 * Licensed under an MIT licence.  Please see LICENCE.md for details.
 */
package vmstack

// 8 element limited stack
type LStack struct {
	// Using 8 as on some platforms AND mask may be quicker than
	// condition
	stack [8]uint
	sp    int
}

func NewLStack() *LStack {
	return &LStack{}
}

// TODO: research best to decrement then get or otherway around
func (s *LStack) pop() uint {
	ts := s.stack[s.sp]

	if s.sp > 0 {
		s.sp--
	}
	return ts
}

func (s *LStack) push(v uint) {
	if s.sp < 7 {
		s.sp++
	}
	s.stack[s.sp] = v
	// NOTE: keep sp at TOS so the we can use peek and replace
}

// TODO: check name
func (s *LStack) peek() uint {
	return s.stack[s.sp]
}

// TODO: check name
func (s *LStack) replace(n uint) {
	s.stack[s.sp] = n
}
