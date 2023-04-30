                ; Version 4
                ; PDP-8 TAD
                ; Assumes lac will only every contain a 13-bit value and
                ; memory will only ever contain a 12-bit value
                ; and checks for overflow rather than using AND

                ; MOV memBase addit (memLoc)
                addit addit
                memBase z
                z addit
                z z

                ; ADD opAddr addit (memLoc)
                opAddr z
                z addit
                z z

                ; ADD I memLoc lac
                ; memLoc is stored in first location of self-modifying instruction
addit:          0 z
                z lac
                z z

                ; IF lac >= 8192 JUMP to overflow
                lac l8192 overflow
                ; ELSE JUMP to done
                z z done

overflow:       ; Remove overflow amount
                l8192c lac

done:           ; Restore l8192
                l8192 l8192
                lm8192 l8192

                ; HLT
                lm1 1000

z:      0
l8192:  8192
l8192c: 8192    ; Used because l8192 gets temporarily corrupted
lm1:    -1
lm8192: -8192

memBase: 53     ; TODO: be able to put memBase or similar here
opAddr:  3
lac:     9
val:     23