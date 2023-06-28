            ; Version 3
            ; SWITCH
            ; Using a table
            LIT 8

loop:       STORE cnt
            FETCH cnt
            DUP
            LIT switchTable
            ADD
            FETCH
            JSR

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
spacer:  0

switchTable: 1     ; TODO: Allow switchTable here or similar
!case0
!case1
!case2
!case3
!case4
!case5
!case6
!case7

lac:     3
ok:      0
cnt:     0