        ; Version 2
        ; SUBLEQ emulator
        ; Uses JGT instead of Skip instructions

        ; Fetch operands
fetch:  MOV     memBase memLoc
        ADD     pc memLoc
        MOV I   memLoc opA
        ADD     l1 memLoc
        MOV I   memLoc opB
        ADD     l1 memLoc
        MOV I   memLoc opC

        ; Execute
exec:   MOV     memBase memA
        ADD     opA memA
        MOV     memBase memB
        ADD     opB memB
        SUB II  memA memB

        ; If opB == 1000 THEN halt
        SNE     opB l1000
        JMP     halt

        ; If mem[opB] <= 0 THEN jump to opC
        JGT I   memB incPC

jmpC:   MOV     opC pc
        JMP     fetch

incPC:  ADD     l3 pc
        JMP     fetch



halt:   MOV I   memB hltVal
        HLT     ok

.data
l1:      1
l3:      3
l1000:   1000
pc:      0
hltVal:  0
memBase: program
memLoc:  0
memA:    0
memB:    0
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