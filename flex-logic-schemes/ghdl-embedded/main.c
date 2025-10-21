#include <stdio.h>
#include "ghdl_embedded.h"

int main() {
    printf("Calling GHDL elaboration...\n");

    int ret = libghdl_main___elabb();
    printf("GHDL elaboration returned: %d\n", ret);

    // If you need to access global objects:
    if (libghdl_main_E) {
        printf("libghdl_main_E is non-null\n");
    }

    return 0;
}

