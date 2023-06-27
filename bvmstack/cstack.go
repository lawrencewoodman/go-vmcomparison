/*
 * A circular stack using Big numbers
 *
 * Copyright (C) 2023 Lawrence Woodman <lwoodman@vlifesystems.com>
 *
 * Licensed under an MIT licence.  Please see LICENCE.md for details.
 */
package bvmstack

import "math/big"

// 8 element circular stack
type CStack struct {
	// Using 8 as on some platforms AND mask may be quicker than
	// condition
	stack [8]*big.Int
	sp    int
}

func NewCStack() *CStack {
	return &CStack{}
}

func (s *CStack) pop() *big.Int {
	ts := s.stack[s.sp]
	if s.sp == 0 {
		s.sp = 7
	} else {
		s.sp--
	}
	return ts
}

func (s *CStack) push(n *big.Int) {
	if s.sp == 7 {
		s.sp = 0
	} else {
		s.sp++
	}
	s.stack[s.sp] = n
	// NOTE: keep sp at TOS so the we can use peek and replace
}

// TODO: check name
func (s *CStack) peek() *big.Int {
	return s.stack[s.sp]
}

// TODO: check name
func (s *CStack) replace(n *big.Int) {
	s.stack[s.sp] = n
}
