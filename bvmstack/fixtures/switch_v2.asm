            ; Version 2
            ; SWITCH
            LIT 8

loop:       STORE cnt
            FETCH cnt
            SHL
            SHL
            FETCH switchBase
            ADD
            JSR

            FETCH cnt
            DJNZ loop
            HLT 1               ; ok


switch:
case0:      FETCH lac
            ADD 11
            STORE lac
            RET

case1:      FETCH lac
            ADD 23
            STORE lac
            RET

case2:      FETCH lac
            ADD 56
            STORE lac
            RET

case3:      FETCH lac
            ADD 79
            STORE lac
            RET

case4:      FETCH lac
            ADD 123
            STORE lac
            RET

case5:      FETCH lac
            ADD 367
            STORE lac
            RET

case6:      FETCH lac
            ADD 592
            STORE lac
            RET

case7:      FETCH lac
            ADD 1001
            STORE lac
            RET

.data
spacer: 0
; TODO: Implement simple maths
;switchBase: switch-4  ; -4 so we don't have to DEC cnt
switchBase: 7
lac:     3
ok:      0
cnt:     0