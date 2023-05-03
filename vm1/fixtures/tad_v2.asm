         ; Version 2
         ; PDP-8 TAD
         LDA  II memBase,opAddr
         ADD     lac
         AND     mask13
         STA     lac
done:    HLT     ok
memBase: 5
opAddr:  5
mask13:  8191    ; 0o17777
lac:     9
ok:      0
val:     23
