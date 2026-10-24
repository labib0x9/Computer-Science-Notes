.model small
.stack 100h
.data
    num1Msg db "A [0-255]: $"
    newLine db 0Dh, 0Ah, "$"

    num1 db ?
.code
main proc
    ; load ds
    mov ax, @data
    mov ds, ax

    ; print msg1
    mov dx, offset num1Msg
    mov ah, 9
    int 21h

    ; input one char
    mov ah, 1
    int 21h
    mov num1, al

    ; print '\n'
    mov dx, offset newLine
    mov ah, 9
    int 21h

    ; print char
    mov dl, num1
    mov ah, 2
    int 21h

    ; exit
    mov ah, 4ch
    int 21h
main endp
end main