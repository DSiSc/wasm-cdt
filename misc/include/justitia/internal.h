#include <stdint.h>
char *justitia_internal_storage_read(const char *key);
void justitia_internal_storage_write(const char *key, const char *val);
char *justitia_internal_sha256(const char *valptr);
char *justitia_internal_call_contract(const char *contract_address, const char *input, uint64_t value);
char *justitia_internal_static_call_contract(const char *contract_address, const char *input, uint64_t value);