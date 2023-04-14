        ; Version 2
        ; JSR
        LIT     50 n
        JSR     setVal sret
done:   HLT     ok 0

val:    0
ok:     0

; setVal
; pass n in n
setVal: MOV     n val
        JMP I   sret 0
n:      0
sret:   0