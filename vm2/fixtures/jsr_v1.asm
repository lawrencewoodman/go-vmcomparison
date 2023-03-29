         ; SIVM2 - Version 1
         ; JSR
         MOV     l150 lcnt
loop:    MOV     l4000 dlen
         JSR     delay dret
         ADD     dtot ltot
         DJNZ    lcnt loop
done:    HLT     ok 0

ltot:    0
lcnt:    0
ok:      0
l1:      1
l50:     50
l150:    150
l4000:   4000

; delay loop
; passes length of delay in dlen
delay:   MOV     l50 dtot
dloop:   ADD     l1 dtot
         DJNZ    dlen dloop
         JMP I   dret 0
dtot:    0
dlen:    0
dret:    0