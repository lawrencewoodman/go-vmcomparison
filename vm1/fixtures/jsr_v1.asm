         ; Version 1
         ; JSR
         LDY     l150
loop:    LDA     l4000
         JSR     delay
         LDA     dtot
         ADD     ltot
         STA     ltot
         ; Decrement Y and Jump if Not Zero
         DYJNZ   loop
done:    HLT     ok
ltot:    0
ok:      0
l50:     50
l150:    150
l4000:   4000

; delay loop
; passes length of delay in AC
; TODO: Allow TAY and RET to assemble without 0 operand
delay:   STY     dystore
         TAY     0
         LDA     l50
         STA     dtot
dloop:   INC     dtot
         DYJNZ   dloop
         LDY     dystore
         RET     0
dystore: 0
dtot:    0


