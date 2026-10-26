#include<stdio.h>
#include<stdlib.h>

// you can also use custom values...

typedef enum DaysOfWeek {
    FRIDAY,     // 0
    SATURSAY,   // 1
    SUNDAY,     // 2
    MONDAY,     // 3
} days_of_week_t;

enum {
    PROTOCOL_TCP,
    PROTOCOL_UDP,
};

int32_t main() {

    days_of_week_t day = FRIDAY;

    return 0;
}