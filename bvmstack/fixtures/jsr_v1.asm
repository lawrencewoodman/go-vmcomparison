        ; Version 1
        ; JSR
		LIT  50     ; value to set
		JSR setVal
		HLT  1      ; ok

setVal:	STORE val
		RET
val:    0