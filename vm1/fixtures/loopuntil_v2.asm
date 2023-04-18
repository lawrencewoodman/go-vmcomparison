        ; Version 2
        ; Uses y register
        ; LOOP UNTIL
        LDY     l5000
        LDA     l0
loop:   ADD     l1
        ; Decrement Y and Jump if Not Zero
        DYJNZ  loop
        STA    sum
done:   HLT    ok
sum:    9
ok:     0
l0:     0
l1:     1
l5000:  5000
