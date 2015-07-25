package main

/*
#cgo LDFLAGS: -lsass
#cgo CFLAGS: -Wall
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <sass_context.h>

int C_CompileSassFile(const char *path, char **output) {
  // create the file context and get all related structs
  struct Sass_File_Context* file_ctx = sass_make_file_context(path);
  struct Sass_Context* ctx = sass_file_context_get_context(file_ctx);
  // struct Sass_Options* ctx_opt = sass_context_get_options(ctx);

  int status = sass_compile_file_context(file_ctx);
  if (status == 0) {
    *output = strdup(sass_context_get_output_string(ctx));
  } else {
    *output = strdup(sass_context_get_error_message(ctx));
  }

  // release allocated memory
  sass_delete_file_context(file_ctx);

  return status;
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func CompileSassFile(path string) (output string, err error) {
	var cpath *C.char = C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	var cout *C.char = nil
	var status = C.C_CompileSassFile(cpath, &cout)
	defer C.free(unsafe.Pointer(cout))

	output = C.GoString(cout)
	if status != 0 {
		return "", fmt.Errorf("Can not compile sass: %s", output)
	} else {
		return output, nil
	}
}
