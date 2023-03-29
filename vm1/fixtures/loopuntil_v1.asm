         ; Version 1
         ; LOOP UNTIL
         LDA     l150
         STA     cnt
         LDA     memBase
         ADD     opAddr
         STA     memloc
         LDA     lac
loop:    ADD I   memloc
         DSZ     cnt
         JMP     loop
         STA12   lac
done:    HLT     ok
memBase: 11
opAddr:  4
lac:     9
ok:      0
val:     23
cnt:     0
l150:    150
memloc:  0

