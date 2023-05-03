            ; Version 3
            ; PDP-8 ISZ
            FETCH memBase
	        FETCH opAddr
	        ADD               ; (memBase opAddr -- valAddr)
			DUP	              ; (valAddr -- valAddr valAddr)
			STORE memLoc      ; (valAddr valAddr -- valAddr)
			FETCH             ; (valAddr -- val)

			ADD 1             ; (val -- val)
			AND 4095          ; 0o7777
			DUP               ; (val -- val val)
			JNZ done          ; (val val -- val)

			FETCH pc
			ADD 1
			AND 4095    ; 0o7777
			STORE pc

done:       FETCH memLoc      ; (val -- val valAddr)
			STORE
			HLT 1
memBase:    17
opAddr:     3
pc:			9
val:        23
memLoc:		0
