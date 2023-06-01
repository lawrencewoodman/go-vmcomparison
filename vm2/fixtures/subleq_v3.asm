        ; Version 3
        ; SUBLEQ emulator
        ; Base PC off program

        ; Fetch operands
fetch:  MOV     memBase memA
        ADD I   pc memA
        ADD     l1 pc
        MOV I   pc opB
        ADD     l1 pc
        MOV I   pc opC
        ADD     l1 pc

        ; Execute
exec:   MOV     memBase memB
        ADD     opB memB
        SUB II  memA memB

        ; If opB == 1000 THEN halt
        SNE     opB l1000
        JMP     halt

        ; IF mem[opB] > 0 THEN fetch next instruction
        JGT I   memB fetch

        ; ELSE mem[opB] <= 0 so set pc to opC
jmpC:   MOV     memBase pc
        ADD     opC pc
        JMP     fetch

halt:   MOV I   memB hltVal
        HLT     ok

l1:      1
l1000:   1000
pc:      program
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