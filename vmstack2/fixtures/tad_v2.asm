            ; Version 2
            ; PDP-8 TAD
            FETCH memBase
	        FETCH opAddr
	        FETCHBI
	        FETCH lac
	        ADD
        	STORE13 lac
	        HLT 1
memBase:    7
opAddr:     3
lac:        9
val:        23