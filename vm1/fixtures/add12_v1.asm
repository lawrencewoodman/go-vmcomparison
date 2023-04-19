        ; Version 1
        ; ADD12
        LDA     a
        ADD     b
        AND     mask12
        STA     b
done:   HLT     ok
ok:     0
a:      4094
b:      6
mask12: 4095    ; 0o7777