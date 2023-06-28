            ; Version 1
            ; PDP-8 TAD
            FETCH memBase
	        FETCH opAddr
	        ADD
	        FETCH
	        FETCH lac
	        ADD
			AND 8191    ; 13-bit mask
        	STORE lac
	        HLT 1

.data
spacer:     0
memBase:    1
opAddr:     3
lac:        9
val:        23