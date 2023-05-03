         ;Version 2
         ;PDP-8 ISZ
         INC II  memBase,opAddr
         LDA II  memBase,opAddr
         AND     mask12
         STA II  memBase,opAddr
         JNZ     done
         INC     pc
         LDA     pc
         AND     mask12
         STA     pc
done:    HLT     ok
memBase: 10
opAddr:  5
mask12:  4095    ; 0o7777
pc:      9
ok:      0
val:     23

