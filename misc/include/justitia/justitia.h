#pragma once
#ifdef __cplusplus
extern "C" {
#endif
#include <stdlib.h>
#include <stdint.h>

//write content to storage
void justitia_storage_write(const char *keyptr, const char *valptr);
//read content from storage, return content pointer if success, NULL if failed
char *justitia_storage_read(const char *keyptr);
//get current block's timestamp
uint64_t justitia_timestamp();
//get current block's height
uint64_t justitia_block_height();
//get current contract's address
char *justitia_self_address();
//get caller's address
char *justitia_caller_address();
//call another contract(use current caller as the contract caller)
char *justitia_call_contract(const char *contract_address, const uint32_t argc, const char *argv[], uint64_t value);
//call another contract(use current contract as the contract caller)
char *justitia_static_call_contract(const char *contract_address, const uint32_t argc, const char *argv[], uint64_t value);
//compute the content's sha256 hash
char *justitia_sha256(const char *valptr);
#ifdef __cplusplus
}
#endif