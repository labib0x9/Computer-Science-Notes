// https://codeforces.com/contest/1902/problem/E

#include<stdio.h>
#include<stdlib.h>

struct Node {
	char exit;
	int freq;
	struct Node* child[26];
};

struct Node* NewNode() {
	struct Node* node = malloc(sizeof(struct Node));
	node->exit = 0;
	node->freq = 0;
	for (int i = 0; i < 26; i++) {
		node->child[i] = NULL;
	}
	return node;
}

struct Node* Trie() {
	return NewNode();
};

void insert(struct Node* root, const char* s, int n) {
	struct Node* cur = root;
	cur->freq++;
	for (int i = 0; i < n; i++) {
		int ch = s[i] - 'a';
		if (!cur->child[ch]) {
			cur->child[ch] = NewNode();
		}
		cur = cur->child[ch];
		cur->freq++;
	}
	cur->exit = 1;
}

long long LCP(struct Node* root, const char* s, int n) {
	long long ans = 0;
	struct Node* cur = root;
	for (int i = 0; i < n; i++) {
		int ch = s[i] - 'a';
		if (!cur->child[ch]) {
			break;
		}
		cur = cur->child[ch];
		ans += cur->freq;
	}
	return ans;
}

void freeTrie(struct Node* node) {
	if (!node) return;
	for (int i = 0; i < 26; i++) {
		if (node->child[i]) {
			freeTrie(node->child[i]);
		}
	}
	free(node);
}

const int N = 1e6 + 5;

int main() {

	char buffer[N];
	struct Node* root = NewNode();

	int n;
	scanf("%d", &n); getchar();

	char *arr[n];

	for (int i = 0; i < n; i++) {
		fgets(buffer, N, stdin);
		buffer[strcspn(buffer, "\n")] = '\0';

		// printf("%s\n", buffer);
		arr[i] = malloc(strlen(buffer) + 1);
		strcpy(arr[i], buffer);
	}

	long long ans = 0;
	for (int i = 0; i < n; i++) {
		ans += strlen(arr[i]);
		insert(root, arr[i], strlen(arr[i]));
	}

	ans = 2ll * ans * n;
	for (int i = 0; i < n; i++) {
		int sLen = strlen(arr[i]);
		for (int j = 0; j < sLen / 2; j++) {
			char temp = arr[i][j];
			arr[i][j] = arr[i][sLen - j - 1];
			arr[i][sLen - j - 1] = temp;
		}
		// printf("%s\n", arr[i]);
		ans -= 2ll * LCP(root, arr[i], sLen);
		// printf("%d", LCP(root, arr[i], sLen));
	}

	printf("%lld\n", ans);


	freeTrie(root);
	for (int i = 0; i < n; i++) {
		free(arr[i]);
	}

	return 0;
}