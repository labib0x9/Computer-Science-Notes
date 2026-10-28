#include<stdio.h>
#include<stdlib.h>
#include<string.h>
#include<stdbool.h>
#include<math.h>

// binary tree of n nodes and d depths
/*
			1
		   / \
		  2   5
		 / \   /\
		3   4   6 7
*/
typedef struct node {
	int val;
	struct node* left;
	struct node* right;
} node;

node* new_node(int val) {
	node* nnode = (node*) malloc(sizeof(node));
	if (nnode == NULL) {
		return NULL;
	}
	nnode->val = val;
	nnode->left = NULL;
	nnode->right = NULL;
	return nnode;
}

int N, depth;
int NowAt = 0;
void create_tree(node* root, int d) {
	if (d == depth) {
		return;
	}

	node* nnode;
	nnode = new_node(++NowAt);
	root->left = nnode;
	create_tree(root->left, d + 1);

	nnode = new_node(++NowAt);
	root->right = nnode;
	create_tree(root->right, d + 1);
}

void visit_tree(node* root, int d) {
	if (root == NULL) {
		return;
	}
	for (int i = 0; i < d; i++) {
		printf("\t");
	}
	printf("%d\n", root->val);

	if (d == depth) {
		return;
	}

	visit_tree(root->left, d + 1);
	visit_tree(root->right, d + 1);
}

void free_tree(node* cur) {
	if (cur == NULL) {
		return;
	}
	free_tree(cur->left);
	free_tree(cur->right);
	free(cur);
}

int main() {

	depth = 2;	// depth -> 0-based indexing..

	node* root = new_node(++NowAt);

	create_tree(root, 0);
	visit_tree(root, 0);

	free_tree(root);

	// printf("%d\n", root->left->left->val);

	return 0;
}