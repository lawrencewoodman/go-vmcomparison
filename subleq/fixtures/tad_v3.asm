                ; Version 3
                ; PDP-8 TAD
                ; Assumes lac will only every contain a 13-bit value and
                ; memory will only ever contain a 12-bit value
                ; and uses the AND routine as a subroutine

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
                ; z z     not needed because of later: z z and13

                ; Store the mask13 mask
                andMask andMask
                mmask13 andMask

                ; Store the return location
                and13ret+2 and13ret+2
                lmdone and13ret+2

                ; Jump to and13
                z z and13

                ; HLT
done:           lm1 1000

z:      0
lm1:    -1
lmdone: 0-done

cnt:    0
memBase: 49  ; TODO: be able to put memBase or similar here
opAddr:  3
lac:     9
val:     23

res:    0
m:      0
lhbitval:   8192
lhbitvalc:  8192
l13:    13
lmhbitval: -8192
mmask13: -8191    ; -0o17777

                ;--------------------------
                ; AND lac
                ; andMask contains value
                ; to mask with
                ;--------------------------

andMask:        0
and13ret:       z z 0
and13:          ; Loop 14 times
                l13 cnt
loop:
                ; m := 0
                m m

                ; res << 2
                res z
                z res
                z z

                ; IF andMask >= lhbitval JUMP to mhbit
                andMask lhbitvalc mhbit
                ; Zero lhbitvalc and JUMP to acheck
                lhbitvalc lhbitvalc acheck

mhbit:          ; INC m
                lm1 m
                ; andMask -= lhbitval
                lhbitval andMask
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

cont:           ; andMask << 2
                andMask z
                z andMask
                z z

                ; lac << 2
                lac z
                z lac
                z z

                ; ADD l1 to cnt and JUMP to loop if <= 0
                lm1 cnt loop

                ; Copy result to lac
                lac lac
                res z
                z lac

                ; Return and zero z
                z z and13ret

                ;------------
                ; End of AND
                ;------------
