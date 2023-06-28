            ; Version 2
            ; PDP-8 ISZ
            FETCH memBase
	        FETCH opAddr
	        ADD               ; (memBase opAddr -- memLoc)
			DUP               ; (memLoc -- memLoc memLoc)
			DUP               ; (memLoc memLoc -- memLoc memLoc memLoc)
			FETCH             ; (memLoc memLoc memLoc -- memLoc memLoc val)
			ADD 1             ; (memLoc memLoc val -- memLoc memLoc val)
			AND 4095          ; 0o7777
			SWAP		      ; (memLoc memLoc val -- memLoc val memLoc)
			STORE             ; (memLoc val memLoc -- memLoc)
			FETCH             ; (memLoc -- val)
			JNZ done          ; (val --)

			FETCH pc
			ADD 1
			AND 4095    ; 0o7777
			STORE pc

done:       HLT 1

.data
spacer:     0
memBase:    1
opAddr:     3
pc:			9
val:        23
