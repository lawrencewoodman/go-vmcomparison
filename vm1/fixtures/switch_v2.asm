          ; Version 2
          ; SWITCH
          ; Using a table
          LDA     l8
          STA     cnt
loop:     LDA II  switchTable,cnt
          STA     caseLoc
          JMP I   caseLoc
decCnt:   DSZ     cnt
          JMP     loop
          HLT     ok

switchTable: 8    ; TODO: Allow switchTable or similar here
case0
case1
case2
case3
case4
case5
case6
case7

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

caseLoc: 0
lac:     3
ok:      0
cnt:     0
l8:      8
l11:    11
l23:    23
l56:    56
l79:    79
l123:   123
l367:   367
l592:   592
l1001:  1001
