                ; Version 2
                ; PDP-8 TAD
                ; Assumes lac will only every contain a 13-bit value and
                ; memory will only ever contain a 12-bit value

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

                ;--------------------------
                ; AND mask13 lac
                ; Routine taken from add12
                ;--------------------------
                ; Loop 14 times
                l13 cnt
loop:
                ; m := 0
                m m

                ; res << 2
                res z
                z res
                z z

                ; IF mask13 >= lhbitval JUMP to mhbit
                mask13 lhbitvalc mhbit
                ; Zero lhbitvalc and JUMP to acheck
                lhbitvalc lhbitvalc acheck

mhbit:          ; INC m
                lm1 m
                ; mask13 -= lhbitval
                lhbitval mask13
                ; COPY lhbitval to lhbitvalc (using subtraction of lmhbitval)
                lhbitvalc lhbitvalc

acheck:         lmhbitval lhbitvalc
                ; IF lac >= lhbitval JUMP to ahbit
                lac lhbitvalc ahbit
                ; COPY lhbitval to lhbitvalc (using subtraction of lmhbitval)
                lhbitvalc lhbitvalc
                lmhbitval lhbitvalc
                z z cont

ahbit:          ; lac -= lhbitval
                lhbitval lac
                ; COPY lhbitval to lhbitvalc (using subtraction of lmhbitval)
                lhbitvalc lhbitvalc
                lmhbitval lhbitvalc
                ; IF m <= 0 JUMP to cont
                z m cont

                ; High bits match
                ; res++
                lm1 res

cont:           ; mask13 << 2
                mask13 z
                z mask13
                z z

                ; lac << 2
                lac z
                z lac
                z z

                ; ADD l1 to cnt and JUMP to loop if <= 0
                lm1 cnt loop

                ;------------
                ; End of AND
                ;------------

                lac lac
                res z
                z lac
                z z

                ; HLT
done:           lm1 1000

z:      0
m:      0
cnt:    0
res:    0
lhbitval:   8192
lhbitvalc:  8192
lm1:    -1
l13:    13
lmhbitval: -8192
mask13: 8191    ; 0o17777
memBase: 136
opAddr:  3
lac:     9
val:     23