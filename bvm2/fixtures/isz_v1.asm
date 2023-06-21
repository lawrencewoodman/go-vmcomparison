         ;Version 1
         ;PDP-8 ISZ
        MOV     memBase memLoc
        ADD     opAddr memLoc
        ADD  DI l1 memLoc
        AND  DI mask12 memLoc
        JNZ  I  memLoc done
        ADD     l1 pc
        AND     mask12 pc
done:   HLT     ok 0

.data
memBase: 0
opAddr:  6
memLoc:  0
mask12:  4095    ; 0o7777
pc:      9
ok:      0
tmp:     23
l1:      1

