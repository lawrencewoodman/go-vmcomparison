        ; Version 1
        ; SUBLEQ emulator

        ; Fetch operands
fetch:  LDA     memBase
        ADD     pc
        STA     memLoc

        ; A
        LDA I   memLoc
        ADD     memBase
        STA     memA
        INC     memLoc

        ; B
        LDA I   memLoc
        STA     opB
        ADD     memBase
        STA     memB
        INC     memLoc

        ; C
        LDA I   memLoc
        STA     opC

        LDA     pc
        ADD     l3
        STA     pc


        ; Execute
exec:   LDA I   memB
        SUB I   memA
        STA I   memB

        ; If opB == 1000 THEN halt
        LDA     l1000
        SUB     opB
        JEQ     halt

        ; IF mem[opB] > 0 THEN jump to fetch
        LDA I   memB
        JGT     fetch

        ; ELSE jump to opC
jmpC:   LDA     opC
        STA     pc
        JMP     fetch


halt:   LDA I   memB
        STA     hltVal
        HLT     ok


l3:      3
l1000:   1000
pc:      0
hltVal:  0
memBase: program
memLoc:  0
memA:    0
memB:    0
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