        ; Version 1
        ; SUBLEQ emulator

        LIT 0
        DROP
        
        ; Fetch operands
fetch:  FETCH pc
        DUP
        ; opA
        FETCH memBase
        ADD
        FETCH
        
        SWAP
        INC            ; increment pc
        DUP

        ; opB
        FETCH memBase
        ADD
        FETCH

        SWAP
        INC            ; increment pc
        DUP

        ; opC
        FETCH memBase
        ADD
        FETCH

        SWAP
        INC            ; increment pc
        STORE pc


        ; Execute
exec:   ROT            ; (opA opB opC -- opB opC opA)
        ROT            ; (opB opC opA -- opC opA opB)
        DUP
        STORE opB
        ; [opB]
        FETCH memBase
        ADD
        FETCH

        ; [opA]
        SWAP
        FETCH memBase
        ADD
        FETCH

        SUB

        FETCH memBase
        FETCH opB
        ADD
        STORE


        ; If opB == 1000 THEN halt
        FETCH opB
        LIT 1000
        SUB
        JZ halt

        ; IF mem[opB] > 0 THEN jump to fetch
        FETCH memBase
        FETCH opB
        ADD
        FETCH
        JGT next

        ; ELSE jump to opC
jmpC:   STORE pc       ; (opC -- )
        JMP fetch

next:   DROP
        JMP fetch

halt:   DROP
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