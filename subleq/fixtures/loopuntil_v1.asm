        ; Version 1
        ; LOOP UNTIL

        ; MOV lm4999 to cnt (Loop 5000 times)
        lm4999 z
        z cnt
        z z
loop:
        ; ADD l1 to sum  (Take -1 from sum)
        lm1 sum

        ; ADD l1 to cnt and JUMP to loop if <= 0
        lm1 cnt loop

        ; HLT
done:   lm1 1000
        
z:      0
cnt:    0
sum:    0
lm4999: -4999
lm1:    -1
