.MODEL SMALL
.STACK 100H
.DATA
    TEXT DB "Hello, world!$"
.CODE
MAIN PROC
    ; Set Data Segement
    MOV AX, @DATA
    MOV DS, AX

    ; Copy TEXT into DX, as ah = 09h reads from DS:DX
    MOV DX, OFFSET TEXT
    MOV AH, 9
    INT 21H

    ; tell program to exit, ah = 4ch meaning terminate the program.
    MOV AH, 4CH
    INT 21H
MAIN ENDP
END MAIN