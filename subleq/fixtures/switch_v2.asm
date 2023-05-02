            ; Version 2
            ; SWITCH
            ; Uses a table

            ; MOV     l8 cnt
            cnt cnt
            lm8 cnt

loop:       ; Modify jmpi to point to correct case statement for cnt
            jmpi+2 jmpi+2
            cnt z
            z getCase
            z z
getCase:    switchTable-1 z
            z jmpi+2
            z z

            ; Restore switchTable address
            cnt getCase

            ; self modifying instruction
            ; JMP to case location
jmpi:       z z 0

decCnt:     ; DJNZ    cnt loop
            l1 cnt done
            z z loop

done:       ; HLT
            lm1 1000
       
switchTable:
case0
case1
case2
case3
case4
case5
case6
case7

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

lac:     3
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