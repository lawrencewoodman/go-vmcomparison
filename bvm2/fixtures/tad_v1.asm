        ; Version 1
        ; PDP-8 TAD
        MOV     memBase memLoc
        ADD     opAddr memLoc
        ADD I   memLoc lac
        AND     mask13 lac
        HLT     ok 0

.data
memBase: 0
opAddr:  6
memLoc:  0
mask13:  8191    ; 0o17777
lac:     9
ok:      0
val:     23
