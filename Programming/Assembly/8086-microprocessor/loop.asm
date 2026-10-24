.model small
.stack 100h
.data
    num1Msg db "A [0-255]: $"
    newLine db 0Dh, 0Ah, "$"
    ; num1 db ?

    ; buffer db 5, 0, 5 dup('$')
.code
main proc
    ; load ds
    mov ax, @data
    mov ds, ax

    ; loop range si = 0 to si = cl
    mov cl, 3   ; count
    mov si, 0   ; source index

    print_loop:
        ; print msg1
        mov dx, offset num1Msg
        call print_msg
        
        ; print '\n'
        mov dx, offset newLine
        call print_msg

        inc si  ; si = si + 1
        loop print_loop ; next iteration

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