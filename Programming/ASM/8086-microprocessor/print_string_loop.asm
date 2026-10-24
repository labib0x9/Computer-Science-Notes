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

    ; exit
    mov ah, 4ch
    int 21h
main endp
input_string proc
    ; input string 
    mov ah, 0Ah
    int 21h

    ; print '\n'
    mov dx, offset newLine
    call print_msg

    ; set si for index, cl for loop
    mov si, offset buffer   ; si = buf[0]
    mov cl, [si+1]  ; cl = char count
    add si, 2   ; si = buf[2]

    ; print buffer[2] to buffer[2+cl]
    print_char_loop:
        ; print buffer[si]
        mov dl, [si]
        mov ah, 2
        int 21h

        ; print '\n'
        mov dx, offset newLine
        call print_msg

        inc si
        loop print_char_loop

    ret
input_string endp
print_msg proc
    mov ah, 9
    int 21h
    ret
print_msg endp
end main