.model small
.stack 100h
.data
    num1Msg db "A [0-255]: $"
    newLine db 0Dh, 0Ah, "$"

    buffer db 5 ; max length of input
            db ? ; count of typed char
            db 5 dup(?) ; stored buffer
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
    mov ah, 0Ah
    int 21h

    ; add '$' to end
    mov si, offset buffer ; si = offset of buffer[0]
    mov cl, [si + 1]      ; cl = offset of buffer[1]
    add si, 2             ; si = offset of buffer[2]
    add si, cx            ; si = offset of buffer[2 + count]; cx = counter register
    mov byte ptr [si], '$'

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
print_msg proc
    mov ah, 9
    int 21h
    ret
print_msg endp
end main