          ; Version 1
          ; SWITCH
          LDA     l8
          STA     cnt
loop:     LDA     cnt
          STA     caseOff
          SHL     caseOff
          SHL     caseOff
          JMP II  switchBase,caseOff
decCnt:   DSZ     cnt
          JMP     loop
          HLT     ok

switch:
case0:    LDA     lac
          ADD     l11
          STA     lac
          JMP     decCnt

case1:    LDA     lac
          ADD     l23
          STA     lac
          JMP     decCnt

case2:    LDA     lac
          ADD     l56
          STA     lac
          JMP     decCnt

case3:    LDA     lac
          ADD     l79
          STA     lac
          JMP     decCnt

case4:    LDA     lac
          ADD     l123
          STA     lac
          JMP     decCnt

case5:    LDA     lac
          ADD     l367
          STA     lac
          JMP     decCnt

case6:    LDA     lac
          ADD     l592
          STA     lac
          JMP     decCnt

case7:    LDA     lac
          ADD     l1001
          STA     lac
          JMP     decCnt

; TODO: Implement simple maths
;switchBase: switch-4  ; -4 so we don't have to DEC cnt
switchBase: 6  ; -4 so we don't have to DEC cnt
caseOff: 0
lac:     3
ok:      0
cnt:     0
l3:      3
l8:      8
l11:    11
l23:    23
l56:    56
l79:    79
l123:   123
l367:   367
l592:   592
l1001:  1001
