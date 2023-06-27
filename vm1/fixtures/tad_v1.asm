         ; Version 1
         ; PDP-8 TAD
         LDA     memBase
         ADD     opAddr
         STA     memLoc
         LDA I   memLoc
         ADD     lac
         AND     mask13
         STA     lac
done:    HLT     ok
memBase: 16
opAddr:  6
memLoc:  0
mask13:  8191    ; 0o17777
lac:     9
ok:      0
val:     23
