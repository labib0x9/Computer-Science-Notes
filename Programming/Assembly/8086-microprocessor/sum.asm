; summation and subtraction of two number 
.model small
.stack 100h
.data
    in1 db "A [0-9]: $"
    in2 db "B [0-9]: $"
    A db ?
    B db ?

    rem db ?
    qut db ?

    ou1 db "A + B = $"
    ou2 db "A - B = $"
    
    nl db 0Dh, 0Ah, "$"

.code
main proc
    mov ax, @data
    mov ds, ax

    ; first int
    mov dx, offset in1
    mov ah, 9
    int 21h

    call input_int
    mov A, al

    ; nl
    mov dx, offset nl
    mov ah, 9
    int 21h

    ; second int
    mov dx, offset in2
    mov ah, 9
    int 21h

    call input_int
    mov B, al

    ; nl
    mov dx, offset nl
    mov ah, 9
    int 21h

    ; sum
    mov dx, offset ou1
    mov ah, 9
    int 21h

    mov al, A
    mov bl, B
    add al, bl
    mov dl, al
    call print_int

    ; nl
    mov dx, offset nl
    mov ah, 9
    int 21h

    ; sub
    mov dx, offset ou2
    mov ah, 9
    int 21h

    mov al, A
    mov bl, B
    sub al, bl
    mov dl, al
    call print_int

    mov ah, 4Ch
    int 21h
main endp


; ------------------------------------------------
; input single digit 0–9
; ------------------------------------------------
input_int proc
    mov ah, 1
    int 21h
    sub al, '0'
    ret
input_int endp


; ------------------------------------------------
; print integer in DL (0–255)
; supports: 1 or 2 digits (00–99)
; ------------------------------------------------
print_int proc
    cmp dl, 10
    jl one_digit       ; if dl < 10 → jump

two_digit:
    mov al, dl
    mov bl, 10
    xor ah, ah
    div bl             ; AL = quotient, AH = remainder

    mov dl, al
    add dl, '0'
    mov ah, 2
    int 21h            ; print quotient

    mov dl, ah
    add dl, '0'
    mov ah, 2
    int 21h            ; print remainder
    ret

one_digit:
    add dl, '0'
    mov ah, 2
    int 21h
    ret
print_int endp

end main
