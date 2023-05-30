        ; Version 3
        ; SUBLEQ emulator
        ; Keep even more things on the stack

        ; Fetch operands
fetch:  FETCH pc
        DUP
        ; opA
        FETCH memBase
        ADD
        FETCH
        
        SWAP
        ADD 1          ; increment pc
        DUP

        ; opB
        FETCH memBase
        ADD
        FETCH

        SWAP
        ADD 1          ; increment pc
        DUP

        ; opC
        FETCH memBase
        ADD
        FETCH

        SWAP
        ADD 1          ; increment pc
        STORE pc


        ; Execute
exec:   ROT             ; (opA opB opC -- opB opC opA)
        ROT             ; (opB opC opA -- opC opA opB)
        DUP             ; (opC opA opB -- opC opA opB opB)

        ; [opB]
        FETCH memBase
        ADD
        FETCH           ; (? -- opC opA opB [opB])

        ; [opA]
        ROT
        FETCH memBase
        ADD
        FETCH           ; (? -- opC opB [opB] [opA])

        SUB             ; (opC opB [opB] [opA] -- opC opB [opB]-[opA])
        OVER            ; (? -- opC opB [opB]-[opA] opB)
        OVER            ; (? -- opC opB [opB]-[opA] opB [opB]-[opA])
        SWAP

        FETCH memBase
        ADD
        STORE           ; (? -- opC opB [opB]+[opA])


        ; If opB == 1000 THEN halt
        SWAP
        LIT 1000
        SUB
        JZ halt

        ; IF mem[opB] > 0 THEN jump to fetch
        JGT next

        ; ELSE jump to opC
jmpC:   STORE pc       ; (opC -- )
        JMP fetch

next:   DROP
        JMP fetch

halt:   DROP
        DROP
        FETCH memBase
        FETCH opB
        ADD
        FETCH
        STORE hltVal
        HLT 1


pc:      0
hltVal:  0
memBase: !program
opA:     0
opB:     0
opC:     0


; loopuntil_v1 from subleq/fixtures/
program:
15
13
3
16
14
6
16
13
3
16 
1000
12
0
0
sum: 0
4999
-1 