        ; Version 1
        ; PDP-8 AND
        MOV     memBase memLoc
        ADD     opAddr memLoc
        MOV I   memLoc tmp
        OR      maskl tmp
        AND     tmp lac
        HLT     ok 0
memBase: 18
opAddr:  7
memLoc:  0
maskl:   4096    ; 0o10000
lac:     4503
tmp:     0
ok:      0
val:     3003
