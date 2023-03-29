         ; Version 3
         ; PDP-8 ISZ
         INC12 II  memBase,opAddr
         LDA   II  memBase,opAddr
         JNZ       done
         INC12     pc
done:    HLT       ok
memBase: 5
opAddr:  4
pc:      9
ok:      0
val:     23

