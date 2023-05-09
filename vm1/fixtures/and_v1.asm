        ; Version 1
        ; PDP-8 AND
        LDA II  memBase,opAddr
        OR      maskl
        AND     lac
        STA     lac
        HLT     ok
memBase: 5
opAddr:  6
memLoc:  0
maskl:   4096    ; 0o10000
lac:     4503
ok:      0
val:     3003
