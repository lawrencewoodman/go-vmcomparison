            ; Version 1
            ; PDP-8 TAD
            FETCH memBase
	        FETCH opAddr
	        ADD
	        FETCH
	        FETCH lac
	        ADD
        	STORE13 lac
	        HLT 1
memBase:    8
opAddr:     3
lac:        9
val:        23