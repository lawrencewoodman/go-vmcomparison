            ; Version 1
            ; SWITCH
            ; Using JMP DI
            ; TODO: Reimplement this, the problem is that instructions are now 3 words

            MOV     l8 cnt
loop:       MOV     cnt caseOff
            SHL     l2 caseOff
            ; We use switch-4 so we don't have to dec cnt
            ; JMP DI  switch-4 caseOff
            JMP DI  8 caseOff
decCnt:     DJNZ    cnt loop
            HLT     ok 0

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
caseOff: 0