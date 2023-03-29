         ; SIVM2 - Version 1
         ; LOOP UNTIL
         MOV     l150 cnt
         MOV     memBase memloc
         ADD     opAddr memloc
loop:    ADD I   memloc lac
         DJNZ    cnt loop
         AND     mask12 lac
done:    HLT     ok 0
memBase: 14
opAddr:  4
lac:     9
ok:      0
val:     23
l150:    150
cnt:     0
memloc:  0
mask12:  4095     ; 0o7777
