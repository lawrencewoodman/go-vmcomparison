        ; Version 1
        ; SUBLEQ emulator

        ; Fetch operands
fetch:  LDA II memBase,pc
        STA     opA
        INC     pc
        LDA II memBase,pc
        STA     opB
        INC     pc
        LDA II memBase,pc
        STA     opC
        INC     pc


        ; Execute
exec:   LDA II  memBase,opB
        SUB II  memBase,opA
        STA II  memBase,opB

        ; If opB == 1000 THEN halt
        LDA     l1000
        SUB     opB
        JEQ     halt

        ; IF mem[opB] > 0 THEN jump to fetch
        LDA II  memBase,opB
        JGT     fetch

        ; ELSE jump to opC
jmpC:   LDA     opC
        STA     pc
        JMP     fetch


halt:   LDA II  memBase,opB
        STA     hltVal
        HLT     ok


l1000:   1000
pc:      0
hltVal:  0
memBase: program
opA:     0
opB:     0
opC:     0
ok:      0


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