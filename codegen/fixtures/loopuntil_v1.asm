        ; Version 1
        ; LOOP UNTIL
        LDA     l5000
        STA     cnt
        LDA     l0
loop:   ADD     l1
        DSZ     cnt
        JMP     loop
        STA     sum
done:   HLT     ok

.data
sum:    0
ok:     0
cnt:    0
l0:     0
l1:     1
l5000:  5000
