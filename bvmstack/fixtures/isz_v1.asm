            ; Version 1
            ; PDP-8 ISZ
            FETCH memBase
	        FETCH opAddr
	        ADD               ; (memBase opAddr -- memLoc)
			DUP               ; (memLoc -- memLoc memLoc)
			FETCH             ; (memLoc memLoc -- memLoc val)
			ADD 1             ; (memLoc val -- memLoc val)
			AND 4095          ; 0o7777
			DUP               ; (memLoc val -- memLoc val val)
			JNZ done          ; (memLoc val val -- memLoc val)
			SWAP              ; (memLoc val -- val memLoc)
			STORE             ; (val memLoc --)

			FETCH pc
			ADD 1
			AND 4095    ; 0o7777
			STORE pc

done:       SWAP
			STORE
			HLT 1

.data
spacer:     0
memBase:    1
opAddr:     3
pc:			9
val:        23
