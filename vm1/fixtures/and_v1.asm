        ; Version 1
        ; PDP-8 AND
        LDA     memBase
        ADD     opAddr
        STA     memLoc
        LDA I   memLoc
        OR      maskl
        AND     lac
        STA     lac
        HLT     ok

memLoc:  0        
memBase: 17
opAddr:  6
memLoc:  0
maskl:   4096    ; 0o10000
lac:     4503
ok:      0
val:     3003
