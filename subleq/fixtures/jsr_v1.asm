        ; Version 1
        ; JSR
        ; Set n to 50
        n n
        lm50 n
        ; Store the return location
        ; sret+2 sret+2
        25 25
        ; lmdone sret+2
        lmdone 25

        ; Jump to setVal
        z z setVal

        ; HLT
done:   lm1 1000

; variables
z:      0
lm1:    1
lm50:   -50
lmdone: -15     ; Location of done negated
val:    0

; setVal
; pass n in n
sret:   z z 0
setVal: val val
        n z
        z val
        ; Return via sret
        z z sret
n:      0
