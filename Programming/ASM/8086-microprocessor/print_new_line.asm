.model small
.stack 100h
.data
    num1 db "A [0-255]: $"
    num2 db "B [0-255]: $"
    newLine db 0Dh, 0Ah, "$"
.code
main proc
    ; load ds
    mov ax, @data
    mov ds, ax

    ; print msg1
    mov dx, offset num1
    mov ah, 9
    int 21h

    ; print '\n'
    mov dx, offset newLine
    mov ah, 9
    int 21h

    ; print msg2
    mov dx, offset num2
    mov ah, 9
    int 21h

    ; exit
    mov ah, 4ch
    int 21h
main endp
end main