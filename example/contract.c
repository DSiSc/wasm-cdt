#include <stdio.h>
#include <string.h>
#include "justitia.h"

char *key = "Hello";

//store data to storage
void test_store() {
    char *val = "World!";
    justitia_storage_write(key, val);
}

//store data to storage
char *test_load() {
    return justitia_storage_read(key);
}

//store data to storage
char *test_hash() {
    return justitia_sha256(key);
}

// get contract address
char *self_address() {
    return justitia_self_address();
}

// get caller address
char *caller_address() {
    return justitia_caller_address();
}

//get current block's timestamp
uint64_t current_timestamp() {
    return justitia_timestamp();
}

//get current block's height
uint64_t current_height() {
    return justitia_block_height();
}

//call another contract
char *call(const char *address, int argc, const char *argv[]) {
    return justitia_call_contract(address, argc, argv, 0);
}

//static call another contract
char *static_call(const char *address,int argc, const char *argv[]) {
    return justitia_static_call_contract(address, argc, argv, 0);
}

// contract entry method
char *invoke(int argc, char *argv[]) {
    if (argc <= 0) {
        return NULL;
    }
    // store data
    if (strcmp(argv[0], "store") == 0) {
        test_store();
    }
    // load data from statedb
    if (strcmp(argv[0], "load") == 0) {
        return test_load();
    }
    // compute hash
    if (strcmp(argv[0], "hash") == 0) {
        return test_hash();
    }
    // get contract address
    if (strcmp(argv[0], "self_address") == 0) {
        return self_address();
    }
    // get caller address
    if (strcmp(argv[0], "caller_address") == 0) {
        return caller_address();
    }
    // get current timestamp
    if (strcmp(argv[0], "current_timestamp") == 0) {
        char *str = (char*)malloc(20);
        sprintf(str, "%lld", current_timestamp());
        return str;
    }
    //get current height
    if (strcmp(argv[0], "current_height") == 0) {
        char *str = (char*)malloc(20);
        sprintf(str, "%lld", current_height());
        return str;
    }
    //call another contract
    if (strcmp(argv[0], "call") == 0 && argc > 2) {
        return call(argv[1], argc - 2, &argv[2]);
    }
    //static another contract
    if (strcmp(argv[0], "static_call") == 0 && argc > 2) {
        return static_call(argv[1], argc - 2, &argv[2]);
    }

    char *method = (char *)malloc(50);
    sprintf(method, "call method %s", argv[0]);
    return method;
}