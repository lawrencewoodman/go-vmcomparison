        ; Version 1
        ; ADD12
        FETCH a
        FETCH b
        ADD
        AND 4095   ; 0o7777
        STORE b

        HLT 1      ; ok

a:      4094
b:      6