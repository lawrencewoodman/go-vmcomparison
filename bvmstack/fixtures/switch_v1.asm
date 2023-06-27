            ; Version 1
            ; SWITCH
            LIT 8

loop:       STORE cnt
            FETCH cnt
            SHL
            SHL
            FETCH switchBase
            ADD
            JMP
decCnt:     FETCH cnt
            DJNZ loop
            HLT 1               ; ok


switch:
case0:      FETCH lac
            ADD 11
            STORE lac
            JMP decCnt

case1:      FETCH lac
            ADD 23
            STORE lac
            JMP decCnt

case2:      FETCH lac
            ADD 56
            STORE lac
            JMP decCnt

case3:      FETCH lac
            ADD 79
            STORE lac
            JMP decCnt

case4:      FETCH lac
            ADD 123
            STORE lac
            JMP decCnt

case5:      FETCH lac
            ADD 367
            STORE lac
            JMP decCnt

case6:      FETCH lac
            ADD 592
            STORE lac
            JMP decCnt

case7:      FETCH lac
            ADD 1001
            STORE lac
            JMP decCnt

; TODO: Implement simple maths
;switchBase: switch-4  ; -4 so we don't have to DEC cnt
switchBase: 7
lac:     3
ok:      0
cnt:     0