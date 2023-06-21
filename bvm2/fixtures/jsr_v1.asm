        ; Version 1
        ; JSR
        MOV     l50 n
        JSR     setVal sret
done:   HLT     ok 0


; setVal
; pass n in n
setVal: MOV     n val
        JMP I   sret 0

.data
n:      0
sret:   0

val:    0
ok:     0
l50:    50
