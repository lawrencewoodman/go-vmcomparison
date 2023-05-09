            ; Version 1
            ; PDP-8 AND
            FETCH memBase
	        FETCH opAddr
	        ADD
	        FETCH
			OR 4096     ; 0o10000
	        FETCH lac
	        AND
        	STORE lac
	        HLT 1
memBase:    9
opAddr:     3
lac:        4503
val:        3003