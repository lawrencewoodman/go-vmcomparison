        ; Version 1
        ; The subleq version of this test operates on 31 bits instead of 32 because
        ; all words are signed
        ; ADD12

        ; ADD a to b
        a z
        z b
        z z

        ; AND mask12 b
        ; hbitval := 2^30
        ; res:= 0
	; for x := 0; x < 31; x++ {
	;   m := 0
        ;   res += res
	;   if mask12 >= hbitval {
	;     m++
	;     mask12 -= hbitval
	;   }
	;   if b >= hbitval {
        ;     b -= hbitval
	;     if m == 1 {
        ;       res++
        ;     }
	;   }
	;   mask12 += mask12
	;   b += b
	; }

        ; Loop 31 times
        l30 cnt
loop:
        ; m := 0
        m m

        ; res << 2
        res z
        z res
        z z

        ; IF mask12 >= lhbitval JUMP to mhbit
        mask12 lhbitvalc mhbit
        ; Zero lhbitvalc and JUMP to bcheck
        lhbitvalc lhbitvalc bcheck

mhbit:  ; INC m
        lm1 m
        ; mask12 -= lhbitval
        lhbitval mask12
        ; COPY lhbitval to lhbitvalc (using subtraction of lmhbitval)
        lhbitvalc lhbitvalc
bcheck: lmhbitval lhbitvalc
        ; IF B >= lhbitval JUMP to bhbit
        b lhbitvalc bhbit
        ; COPY lhbitval to lhbitvalc (using subtraction of lmhbitval)
        lhbitvalc lhbitvalc
        lmhbitval lhbitvalc
        z z cont

bhbit:  ; b -= lhbitval
        lhbitval b
        ; COPY lhbitval to lhbitvalc (using subtraction of lmhbitval)
        lhbitvalc lhbitvalc
        lmhbitval lhbitvalc
        ; IF m <= 0 JUMP to cont
        z m cont

        ; High bits match
        ; res++
        lm1 res

cont:   ; mask12 << 2
        mask12 z
        z mask12
        z z

        ; b << 2
        b z
        z b
        z z

        ; ADD l1 to cnt and JUMP to loop if <= 0
        lm1 cnt loop

        ; HLT
done:   lm1 1000

z:      0
a:      4094
b:      6
m:      0
cnt:    0
res:    0
lhbitval:   1073741824
lhbitvalc:  1073741824
lm1:    -1
l30:    30
lmhbitval: -1073741824
mask12: 4095    ; 0o7777