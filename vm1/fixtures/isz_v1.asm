         ;Version 1
         ;PDP-8 ISZ
         LDA     memBase
         ADD     opAddr
         STA     memLoc
         INC I   memLoc
         LDA I   memLoc
         AND     mask12
         STA I   memLoc
         JNZ     done
         INC     pc
         LDA     pc
         AND     mask12
         STA     pc
done:    HLT     ok
memBase: 13
opAddr:  6
memLoc:  0
mask12:  4095    ; 0o7777
pc:      9
ok:      0
tmp:     23

