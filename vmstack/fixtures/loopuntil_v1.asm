        ; Version 1
        ; LOOP UNTIL
		LIT 0      ; (-- sum)
		LIT 5000   ; (sum -- sum cnt)

loop:   SWAP       ; (sum cnt -- cnt sum)
		ADD 1
		SWAP       ; (cnt sum -- sum cnt)
		DJNZ loop
		DROP       ; (sum cnt -- sum)
		STORE sum
		HLT 1      ; ok

sum:    0