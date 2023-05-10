            ; Version 2
            ; SWITCH
            ; Using table

            MOV     l8 cnt
loop:       MOV     switchTable tableLoc
            ADD     cnt tableLoc
            MOV I   tableLoc caseLoc
            JMP I   caseLoc 0
decCnt:     DJNZ    cnt loop
            HLT     ok 0

switchTable: 14    ; TODO: Allow switchTable or similar here
case0
case1
case2
case3
case4
case5
case6
case7


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
tableLoc: 0
caseLoc: 0