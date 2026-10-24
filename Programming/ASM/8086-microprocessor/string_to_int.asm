.model small
.stack 100h
.data
    num1Msg db "A [0-255]: $"
    newLine db 0Dh, 0Ah, "$"
    num1 dw ?

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
    call input_int
    mov num1, ax

    ; exit
    mov ah, 4ch
    int 21h
main endp
input_int proc
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

    ; set result environment
    xor ax, ax
    xor bx, bx
    xor dx, dx

    ; print buffer[2] to buffer[2+cl]
    convert_loop:
        mov bl, [si]    ; bl = buffer[si]
        sub bl, '0'     ; bl = '9' - '0'

        ; ax = ax * 10
        mov dx, ax
        shl ax, 1   ; ax = ax * 2
        ; shl dx, 3   ; dx = dx * 8 = ax * 8
        shl dx, 1
        shl dx, 1
        shl dx, 1
        add ax, dx  ; ax = 2ax + 8ax = 10ax = ax * 10

        ; ax = ax + digit
        add ax, bx  ; ax = ax + bx[bl]

        inc si
        loop convert_loop

    ret
input_int endp
print_msg proc
    mov ah, 9
    int 21h
    ret
print_msg endp
end main