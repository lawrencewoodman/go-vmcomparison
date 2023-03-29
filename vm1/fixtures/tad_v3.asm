         ; Version 3
         ; PDP-8 TAD
         LDA II  memBase,opAddr
         ADD     lac
         STA13   lac
done:    HLT     ok
memBase: 4
opAddr:  4
lac:     9
ok:      0
val:     23
