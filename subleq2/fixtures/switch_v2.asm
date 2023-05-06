            ; Version 2
            ; SWITCH
            ; Uses a table to find the case location
            
            ; MOV     l8 cnt
            cnt cnt
            lm8 cnt

loop:       ; Lookup case location using cnt as an index in switchTable
            caseLoc caseLoc
            tableLoc tableLoc
            cnt z
            switchTable z
            z tableLoc
            ; z z   not needed because of z z in jmpi
            [tableLoc] caseLoc

            ; JMP I   caseLoc
jmpi:       z z [caseLoc]

decCnt:     ; DJNZ    cnt loop
            l1 cnt done
            z z loop

done:       ; HLT
            lm1 1000

// TODO: be able to put here or switchTable here not rely on done+3
switchTable: done+3
0-case0
0-case1
0-case2
0-case3
0-case4
0-case5
0-case6
0-case7       

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

tableLoc: 0
caseLoc: 0
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