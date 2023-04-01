            ; Version 2
            ; Using a table to jump into switch
            ; SWITCH
            MOV     l8 cnt
loop:       JMPX I  caseJumpBase cnt
decCnt:     DJNZ    cnt loop
            HLT     ok 0

; TODO: Implement some sort of here indicator
caseJumpBase:  8   ; here
17 ; case0
21 ; case1
25 ; case2
29 ; case3
33 ; case4
37 ; case5
41 ; case6
45 ; case7


switch:
case0:      ADD     l11 lac
            JMP     decCnt 0

case1:      ADD     l23 lac
            JMP     decCnt 0

case2:      ADD     l56 lac
            JMP     decCnt 0

case3:      ADD     l79 lac
            JMP     decCnt 0

case4:      ADD     l123 lac
            JMP     decCnt 0

case5:      ADD     l367 lac
            JMP     decCnt 0

case6:      ADD     l592 lac
            JMP     decCnt 0

case7:      ADD     l1001 lac
            JMP     decCnt 0


lac:     3
ok:      0
cnt:     0
l2:      2
l8:      8
l11:    11
l23:    23
l56:    56
l79:    79
l123:   123
l367:   367
l592:   592
l1001:  1001
