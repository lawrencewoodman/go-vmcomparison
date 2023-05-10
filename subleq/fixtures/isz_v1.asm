                ; Version 2
                ; PDP-8 ISZ
                ; Assumes memory will only every contain a 12-bit value
                ; memory will only ever contain a 12-bit value

                ; MOV memBase+opAddr TO getit+3, setit, setit+1, setit+7 (memLoc)
                getit+3 getit+3
                setit setit
                setit+1 setit+1
                setit+7 setit+7
                memBase z
                opAddr z
                z getit+3
                z setit
                z setit+1
                z setit+7
                z z

                ; MOV I memLoc t
                ; memLoc is stored in the fourth location of self-modifying instruction
getit:          t t
                0 z
                z t
                z z

                ; Increment t
                lm1 t

                ; IF t >= 4096 JUMP to toverflow
                t l4096 toverflow
                ; ELSE JUMP to setit
                z z setit

toverflow:      ; Remove overflow amount
                l4096c t

                ; MOV DI t memLoc
                ; memLoc is stored in the first, second and eighth location
                ; of self-modifying instruction
setit:          0 0
                t z
                z 0
                z z

                ; Restore l4096
                l4096 l4096
                lm4096 l4096

                ; IF t <= 0 JUMP to tle
                z t tle
                ; ELSE JUMP to done
                z z done

tle:            ; IF z <= t JUMP to teq (equivalent to IF t == 0)
                t z teq
                ; ELSE
                z z done

                ; T == 0, so increment PC
teq:            lm1 pc

                ; IF pc >= 4096 JUMP to pcoverflow
                pc l4096 pcoverflow
                ; ELSE JUMP to done
                z z pcresl4096

pcoverflow:     ; Remove overflow amount
                l4096c pc

pcresl4096:     ; Restore l4096
                l4096 l4096
                lm4096 l4096

done:           ; HLT
                lm1 1000

z:      0
t:      0
l4096:  4096
l4096c: 4096    ; Used because l4096 temporarily gets corrupted
lm1:    -1
lm4096: -4096
memBase: 114
opAddr:  3
pc:      9
val:     23