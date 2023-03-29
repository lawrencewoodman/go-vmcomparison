         ; Version 2
         ; PDP-8 TAD
         LDA     memBase
         ADD     opAddr
         STA     memLoc
         LDA I   memLoc
         ADD     lac
         STA13   lac
done:    HLT     ok
memBase: 7
opAddr:  5
memLoc:  0
lac:     9
ok:      0
val:     23
