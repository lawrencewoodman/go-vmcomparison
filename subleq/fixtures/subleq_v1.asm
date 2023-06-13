        ; Version 1
        ; SUBLEQ emulator

        ; Fetch operands
fetch:  ; Modify getA to create: [pc] z
        getA getA
        pc z
        z getA
        z z

        ; [pc] z
getA:   0 z
        ; Modify exec to create: [memA] [memB]   - this just does [memA]
        exec exec
        z exec
        mmemBase exec
        z z

        ; opB and memB
        lm1 pc
        opB opB
        memB memB
        mmemBase memB

        ; Modify getB to create: [pc] z
        getB getB
        pc z
        z getB
        z z

        ; [pc] z
getB:   0 z
        z opB
        z memB
        z z

        ; Modify exec to create: [memA] [memB]   - this just does [memB]
        ; Modify cmemB to create: z [memB] jmpC
        ; Modify halt to create: [memB] z
        exec+1 exec+1
        cmemB+1 cmemB+1
        halt halt
        memB z
        z exec+1
        z cmemB+1
        z halt
        z z

        ; opC
        ; Store opC as a negative number to make mov to PC quicker
        lm1 pc
        opC opC

        ; Modify getC to create: [pc] opC
        getC getC
        pc z
        z getC
        z z

        ; [pc] opC
getC:   0 opC

        ; Execute
        ; [memA] [memB]
exec:   0 0

        ; If opB == 1000 THEN halt
        opB l1000 bge1000
        z z ifJmp

bge1000: l1000 z halt

        ; If mem[opB] <= 0 THEN jump to opC
ifJmp:  l1000 l1000
        lm1000 l1000

        ; z [memB] jmpC
cmemB:  z 0 jmpC

incPC:  lm1 pc
        z z fetch

jmpC:   pc pc
        mmemBase pc
        opC pc
        z z fetch

        ; HLT
        ; [memB] z
halt:   0 z
        z hltVal
        z z
        lm1 1000

lm1:    -1
lm1000: -1000
l1000:  1000
z:      0
; TODO: work out why need 0+program instead of just program
pc:     0+program
hltVal: 0
mmemBase: 0-program
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