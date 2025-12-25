#include <stdio.h>  
extern void GoPrint(int);  
int main() {  
    printf("C start\n");
    GoPrint(666);
    GoPrint(888);
    printf("C end\n");
    return 0;
}