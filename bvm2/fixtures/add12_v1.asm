        ; Version 1
        ; ADD12
        ADD     a b
        AND     mask12 b
done:   HLT     ok 0

.data
ok:     0
a:      4094
b:      6
mask12: 4095    ; 0o7777