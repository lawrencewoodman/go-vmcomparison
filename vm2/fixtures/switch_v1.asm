            ; Version 1
            ; SWITCH
            MOV     l8 cnt
loop:       MOV     cnt caseOff
            ; TODO: should we do a single SHL and put result in another location
            ; TODO: or shift left by A number of times
            SHL     l2 caseOff
            MOV     switchBase caseLoc
            ADD     caseOff caseLoc
            ; TODO: remove need to put 0 at end of JMP
            JMP I   caseLoc 0
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

; TODO: Implement simple maths
;switchBase: switch-4  ; -4 so we don't have to DEC cnt
switchBase: 12  ; -4 so we don't have to DEC cnt
caseOff: 0
caseLoc: 0
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
