        ; Version 1
        ; JSR
        ; Set n to 50
        n n
        lm50 n
        ; Store the return location
        sret sret
        lmdone sret

        ; Jump to setVal
        z z setVal

        ; HLT
done:   lm1 1000

; variables
z:      0
lm1:    1
lm50:   -50
lmdone: 0-done     ; Location of done negated
val:    0

; setVal
; pass n in n
sret:   0
setVal: val val
        n z
        z val
        ; Return
        z z [sret]
n:      0
