         ; SIVM2 - Version 2
         ; JSR
         LIT     150 lcnt
loop:    LIT     4000 dlen
         JSR     delay dret
         ADD     dtot ltot
         DJNZ    lcnt loop
done:    HLT     ok 0

ltot:    0
lcnt:    0
ok:      0
l1:      1

; delay loop
; passes length of delay in dlen
delay:   LIT     50 dtot
dloop:   ADD     l1 dtot
         DJNZ    dlen dloop
         JMP I   dret 0
dtot:    0
dlen:    0
dret:    0