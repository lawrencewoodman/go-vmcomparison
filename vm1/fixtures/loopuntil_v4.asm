         ; Version 4
         ; Uses an x index register
         ; LOOP UNTIL
         LDA     l150
         STA     cnt
         LDA     lac
         LDX     opAddr
loop:    ADDIX    memBase
         DSZ     cnt
         JMP     loop
         STA12   lac
done:    HLT     ok
memBase: 9
opAddr:  4
lac:     9
ok:      0
val:     23
cnt:     0
l150:    150

