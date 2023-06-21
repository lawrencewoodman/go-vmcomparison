        ; Version 1
        ; LOOP UNTIL
        MOV     l5000 cnt
loop:   ADD     l1 sum
        DJNZ    cnt loop
done:   HLT     ok 0

.data
sum:    0
ok:     0
l1:     1
l5000:  5000
cnt:    0