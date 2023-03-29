         ; Version 2
         ; PDP-8 ISZ
         LDA     memBase
         ADD     opAddr
         STA     memLoc
         INC12 I memLoc
         JNZ     done
         INC12   pc
done:    HLT     ok
memBase: 7
opAddr:  5
memLoc:  0
pc:      9
ok:      0
tmp:     23
