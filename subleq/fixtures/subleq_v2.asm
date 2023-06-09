        ; Version 2
        ; SUBLEQ emulator
        ; Store PC as a negative number to increase speed
        ; Use this fact to combine an update to PC with a jump

        ; Fetch operands
fetch:  ; Modify getA to create: [pc] z
        getA getA
        pc getA
        z z

        ; [pc] z
getA:   0 z
        ; Modify exec to create: [memA] [memB]   - this just does [memA]
        exec exec
        z exec
        mmemBase exec
        z z

        ; opB and memB
        l1 pc
        opB opB

        ; Modify getB to create: [pc] z
        getB getB
        pc getB

        ; [pc] z
getB:   0 z
        z opB
        ; Use this to create memB := memBase + opB
        memBase z

        ; Modify exec to create: [memA] [memB]   - this just does [memB]
        ; Modify cmemB to create: z [memB] jmpC
        ; Modify halt to create: [memB] z
        exec+1 exec+1
        cmemB+1 cmemB+1
        halt halt
        z exec+1
        z cmemB+1
        z halt
        z z

        ; opC
        l1 pc
        opC opC

        ; Modify getC to create: [pc] z
        getC getC
        pc getC
        z z

        ; [pc] z
getC:   0 z
        z opC
        ; z z - this is done below

        ; Execute
        ; [memA] [memB]
exec:   0 0

        ; If opB == 1000 THEN halt
        opB l1000 bge1000
        z z ifJmp

bge1000: z z
        l1000 z halt

        ; If mem[opB] <= 0 THEN jump to opC
ifJmp:  l1000 l1000
        lm1000 l1000

        ; z [memB] jmpC
cmemB:  z 0 jmpC

incPC:  l1 pc fetch

jmpC:   pc pc
        memBase pc
        opC pc fetch

        ; HLT
        ; [memB] z
halt:   0 z
        z hltVal
        z z
        lm1 1000

lm1:    -1
l1:     1
lm1000: -1000
l1000:  1000
z:      0
pc:     0-program
hltVal: 0
memBase: program
mmemBase: 0-program
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