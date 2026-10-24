.model small
.stack 100h
.data
    num1Msg db "A [0-255]: $"
    newLine db 0Dh, 0Ah, "$"

    num1 db ?, "$"
.code
main proc
    ; load ds
    mov ax, @data
    mov ds, ax

    ; print msg1
    mov dx, offset num1Msg
    call print_msg

    ; input one char
    mov ah, 1
    int 21h
    mov num1, al

    ; print '\n'
    mov dx, offset newLine
    call print_msg

    ; print num1
    mov dx, offset num1
    call print_msg

    ; exit
    mov ah, 4ch
    int 21h
main endp
print_msg proc
    mov ah, 9
    int 21h
    ret
print_msg endp
end main