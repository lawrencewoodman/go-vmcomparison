         ; Version 2
         ; LOOP UNTIL
         LDA     l150
         STA     cnt
         LDA     lac
loop:    ADD II  memBase,opAddr
         DSZ     cnt
         JMP     loop
         STA12   lac
done:    HLT     ok
memBase: 8
opAddr:  4
lac:     9
ok:      0
val:     23
cnt:     0
l150:    150

