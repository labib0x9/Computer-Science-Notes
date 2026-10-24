#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include<stdbool.h>
#include<unistd.h>

int main() {

	/* code - 01 */

	char buffer[3];
	gets(buffer);

	printf("%s\n", buffer);

	/* code - 02 */
	char buf[1024];
	fgets(&buf, sizeof(buf), stdin);
	buffer[strcspn(buf, "\n")] = '\0';

	char newStr[3];
	strcpy(newStr, buf);

	printf("%s\n", newStr);

	return 0;
}