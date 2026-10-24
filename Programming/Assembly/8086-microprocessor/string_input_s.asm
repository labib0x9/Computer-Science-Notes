.model small
.stack 100h
.data
    num1Msg db "A [0-255]: $"
    newLine db 0Dh, 0Ah, "$"

    buffer db 5, 0, 5 dup('$')
.code
main proc
    ; load ds
    mov ax, @data
    mov ds, ax

    ; print msg1
    mov dx, offset num1Msg
    call print_msg

    ; input string
    mov dx, offset buffer
    call input_string

    ; print '\n'
    mov dx, offset newLine
    call print_msg

    ; print string
    mov dx, offset buffer + 2
    call print_msg

    ; exit
    mov ah, 4ch
    int 21h
main endp
input_string proc
    mov ah, 0Ah
    int 21h
    ret
input_string endp
print_msg proc
    mov ah, 9
    int 21h
    ret
print_msg endp
end main