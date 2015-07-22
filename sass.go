package main

/*
#cgo LDFLAGS: -lsass
#include <stdlib.h>
#include <stdio.h>
#include <sass_context.h>

void C_CompileSassFile(const char *path) {
  // create the file context and get all related structs
  struct Sass_File_Context* file_ctx = sass_make_file_context(path);
  struct Sass_Context* ctx = sass_file_context_get_context(file_ctx);
  struct Sass_Options* ctx_opt = sass_context_get_options(ctx);

  // context is set up, call the compile step now
  int status = sass_compile_file_context(file_ctx);

  // print the result or the error to the stdout
  if (status == 0) puts(sass_context_get_output_string(ctx));
  else puts(sass_context_get_error_message(ctx));

  // release allocated memory
  sass_delete_file_context(file_ctx);
}
*/
import "C"
import "unsafe"

func CompileSassFile(path string) {
	var cpath *C.char = C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	C.C_CompileSassFile(cpath)
}
