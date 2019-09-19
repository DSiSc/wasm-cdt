#include <string.h>
#include "internal.h"
#include "justitia.h"


char *clone(const char *valptr) {
    unsigned long val_len = strlen(valptr);
    char *val_cloned =(char*)malloc(val_len + 1);
    memcpy(val_cloned, valptr, val_len + 1);
    return val_cloned;
}

char *concat_str(const uint32_t argc, const char *argv[]) {
    unsigned long total_len = 0;
    for(int i=0;i<argc;i++)
    {
        total_len += strlen(argv[i]) + 1;
    }
    char *val =  (char *)malloc(total_len);
    unsigned long start = 0;
    for(int i=0;i<argc;i++)
    {
        unsigned long argv_len = strlen(argv[i]);
        memcpy(val+start*sizeof(char), argv[i], argv_len);
        start += argv_len;
        val[start] = ',';
        start++;
    }
    return val;
}

void justitia_storage_write(const char *keyptr, const char *valptr){
    char *keyptr_tmp = clone(keyptr);
    char *valptr_tmp = clone(valptr);
    justitia_internal_storage_write(keyptr_tmp, valptr_tmp);
    free(keyptr_tmp);
    free(valptr_tmp);
}

char *justitia_storage_read(const char *keyptr){
    char *keyptr_tmp = clone(keyptr);
    char *val = justitia_internal_storage_read(keyptr_tmp);
    free(keyptr_tmp);
    return val;
}

char *justitia_call_contract(const char *contract_address, const uint32_t argc, const char *argv[], uint64_t value) {
    char *addr = clone(contract_address);
    char *input = concat_str(argc, argv);
    char *ret = justitia_internal_call_contract(addr, input, value);
    free(addr);
    free(input);
    return ret;
}
char *justitia_static_call_contract(const char *contract_address, const uint32_t argc, const char *argv[], uint64_t value) {
    char *addr = clone(contract_address);
    char *input = concat_str(argc, argv);
    char *ret = justitia_internal_static_call_contract(addr, input, value);
    free(addr);
    free(input);
    return ret;
}

char *justitia_sha256(const char *valptr){
    char *valptr_tmp = clone(valptr);
    char *hash = justitia_internal_sha256(valptr_tmp);
    free(valptr_tmp);
    return hash;
}
