         ; Version 4
         ; Uses an x and y index register
         ; LOOP UNTIL
         LDY     l150
         LDA     lac
         LDX     opAddr
loop:    ADDIX   memBase
         ; Decrement Y and Jump if Not Zero
         DYJNZ   loop
         STA12   lac
done:    HLT     ok
memBase: 8
opAddr:  3
lac:     9
ok:      0
val:     23
l150:    150
