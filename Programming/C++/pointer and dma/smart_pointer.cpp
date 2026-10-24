#include<iostream>
#include<memory>

using namespace std;

struct  Tree {
    Tree() {
        cout << "construction" << endl;
    }

    ~Tree() {
        cout << "destruction" << endl;
    }

    void Print() {
        cout << "Hello unique pointer" << endl;
    }
};


unique_ptr<int> Get(int n = 0) {
    unique_ptr<int> v(new int(n));
    cout << "Get(): " << v << endl;
    // return std::move(v);
    return v;   // automatically calls move..
}

int32_t main() {

    // no need to call delete manually,
    // when scope ends, delete is automatically called..
    // single owner..
    unique_ptr<int> value(new int(15));

    cout << value << " " << *value << endl;
    // vector<unique_ptr<int>> ptrs;

    value.reset(new int(42));   // reallocates & delete the prev allocation..
    cout << value << " " << *value << endl; 

    // value = 35; // error
    *value = 35;
    cout << value << " " << *value << endl;

    unique_ptr<string> strPtr(new string);
    strPtr->assign("Hello World");

    cout << strPtr << " " << *strPtr << endl;

    // no copy, but you can move
    unique_ptr<int> valueNew(std::move(value)); // transfer ownership to valueNow, value becomes nullptr
    cout << valueNew << " " << *valueNew << endl;
    // cout << value << " " << *value << endl;
    cout << (value == nullptr ? "true" : "false") << endl;

    // what if, the value escapes the scope ??
    unique_ptr<int> vPtr = Get();
    cout << "main(): " << vPtr << endl;

    {
        unique_ptr<Tree> newTree(new Tree());
        // unique_ptr<Tree> newTree = make_unique<Tree>(); // something is wrong here ?
        newTree->Print();
    }


    shared_ptr<int> sValue(new int(19));
    shared_ptr<int> nNewValue = sValue;

    {
        shared_ptr<Tree> sTree;
        {
            cout << "Inner scope: ";
            shared_ptr<Tree> st(new Tree());
            sTree = st;
        }
        cout << "Outer scope: ";
    }

    // weak_ptr;

    return 0;
}