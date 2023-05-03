            ; Version 2
            ; PDP-8 TAD
            FETCH memBase
	        FETCH opAddr
	        FETCHBI
	        FETCH lac
	        ADD
			AND 8191    ; 13-bit mask
        	STORE lac
	        HLT 1
memBase:    8
opAddr:     3
lac:        9
val:        23