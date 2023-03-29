/*
 * A circular stack
 *
 * Copyright (C) 2023 Lawrence Woodman <lwoodman@vlifesystems.com>
 *
 * Licensed under an MIT licence.  Please see LICENCE.md for details.
 */
package vmstack

// 8 element circular stack
type CStack struct {
	// Using 8 as on some platforms AND mask may be quicker than
	// condition
	stack [8]uint
	sp    int
}

func NewCStack() *CStack {
	return &CStack{}
}

func (s *CStack) pop() uint {
	ts := s.stack[s.sp]
	if s.sp == 0 {
		s.sp = 7
	} else {
		s.sp--
	}
	return ts
}

func (s *CStack) push(v uint) {
	if s.sp == 7 {
		s.sp = 0
	} else {
		s.sp++
	}
	s.stack[s.sp] = v
	// NOTE: keep sp at TOS so the we can use peek and replace
}

// TODO: check name
func (s *CStack) peek() uint {
	return s.stack[s.sp]
}

// TODO: check name
func (s *CStack) replace(n uint) {
	s.stack[s.sp] = n
}
