package regexp

// #ifndef USE_LIBSQLITE3
// #include <sqlite3-binding.h>
// #else
// #include <sqlite3.h>
// #endif
//
// // Extension function defined in regexp.c.
// extern int sqlite3_regexp_init(sqlite3*, char**, const sqlite3_api_routines*);
//
// // Use constructor to register extension function with sqlite.
// void __attribute__((constructor)) init(void) {
//   sqlite3_auto_extension((void*) sqlite3_regexp_init);
// }
import "C"
