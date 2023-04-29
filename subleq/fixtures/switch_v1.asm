            ; Version 1
            ; SWITCH
            ; MOV     l8 cnt
            cnt cnt
            lm8 cnt

loop:       ; MOV     cnt caseOff
            caseOff caseOff
            cnt z
            z caseOff
            z z

            ; SHL     l3 caseOff
            caseOff z
            z caseOff
            z z
            caseOff z
            z caseOff
            z z
            caseOff z
            z caseOff
            z z

            ; MOV     switchBase jmpi+2
            jmpi+2 jmpi+2
            switchBase z
            z jmpi+2
            z z

            ; ADD     caseOff caseLoc
            caseOff z
            z jmpi+2
            z z

            ; self modifying instruction
            ; JMP I   caseLoc 0
jmpi:       z z 0

decCnt:     ; DJNZ    cnt loop
            l1 cnt done
            z z loop

done:       ; HLT
            lm1 1000
       

switch:
case0:      ; ADD     l11 lac
            lm11 lac
            ; JMP     decCnt
            z z decCnt
            0    ; padding to make up to 8 words
            0    ; /
            

case1:      ; ADD     l23 lac
            lm23 lac
            ; JMP     decCnt
            z z decCnt
            0
            0

case2:      ; ADD     l56 lac
            lm56 lac
            ; JMP     decCnt
            z z decCnt
            0
            0

case3:      ; ADD     l79 lac
            lm79 lac
            ; JMP     decCnt
            z z decCnt
            0
            0

case4:      ; ADD     l123 lac
            lm123 lac
            ; JMP     decCnt
            z z decCnt
            0
            0

case5:      ; ADD     l367 lac
            lm367 lac
            ; JMP     decCnt
            z z decCnt
            0
            0

case6:      ; ADD     l592 lac
            lm592 lac
            ; JMP     decCnt
            z z decCnt
            0
            0

case7:      ; ADD     l1001 lac
            lm1001 lac
            ; JMP     decCnt
            z z decCnt
            0
            0

switchBase: switch-8  ; -8 so we don't have to DEC cnt
caseOff: 0
lac:     3
ok:      0
z:       0
cnt:     0
l1:      1
lm1:     -1
lm8:     -8
lm11:    -11
lm23:    -23
lm56:    -56
lm79:    -79
lm123:   -123
lm367:   -367
lm592:   -592
lm1001:  -1001