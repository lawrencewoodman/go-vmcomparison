         ; Version 1
         ; JSR
         LDA     l50
loop:    JSR     setVal
done:    HLT     ok

val:     0
ok:      0
l50:     50

; setVal
; pass n in AC
setVal:  STA     val
         RET     0

