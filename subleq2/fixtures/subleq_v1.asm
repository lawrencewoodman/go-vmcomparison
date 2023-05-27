        ; Version 1
        ; SUBLEQ emulator

        ; Fetch operands
fetch:  memLoc memLoc
        opA opA
        opB opB
        opC opC
        mmemBase memLoc
        pc z
        z memLoc
        z z
        [memLoc] z
        z opA
        z z
        lm1 memLoc
        [memLoc] z
        z opB
        z z
        lm1 memLoc
        [memLoc] z
        z opC
        z z

exec:   ; Execute
        mmemBase opA
        memB memB
        mmemBase memB
        opB z
        z memB
        z z
        [opA] [memB]


        ; If opB == 1000 THEN halt
        opB l1000 bge1000
        z z ifJmp

bge1000: l1000 z halt

        ; If mem[opB] <= 0 THEN jump to opC
ifJmp:  l1000 l1000
        lm1000 l1000

        z [memB] jmpC

incPC:  lm3 pc
        z z fetch

jmpC:   pc pc
        opC z
        z pc
        z z fetch

        ; HLT
halt:   hltVal hltVal
        [memB] z
        z hltVal
        z z
        lm1 1000

lm1:    -1
lm3:    -3
lm1000: -1000
l1000:  1000
z:      0
pc:     0
hltVal: 0
mmemBase: 0-program
memBase: program
memLoc:  0
opA:     0
opB:     0
opC:     0
memB:    0

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