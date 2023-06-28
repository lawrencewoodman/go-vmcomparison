        ; Version 1
        ; ADD12
        FETCH a
        FETCH b
        ADD
        AND 4095   ; 0o7777
        STORE b

        HLT 1      ; ok

.data
spacer: 0          ; Used so that we can reference a as an operand
a:      4094
b:      6