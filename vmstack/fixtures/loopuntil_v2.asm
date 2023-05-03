        ; Version 2
        ; LOOP UNTIL
		LIT 0      ; (-- sum)
		STORE sum  ; (--)
		LIT 5000   ; (-- cnt)

loop:   FETCH sum
		ADD 1
		STORE sum
		DJNZ loop
		HLT 1      ; ok

sum:    0