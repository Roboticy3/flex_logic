#ifndef GHDL_EMBEDDED_H
#define GHDL_EMBEDDED_H

#ifdef __cplusplus
extern "C" {
#endif

// Function to elaborate the default design
extern int libghdl_main___elabb(void);



// You can also declare global objects if needed
extern void *libghdl_main_E;  // This is a global object

#ifdef __cplusplus
}
#endif

#endif // GHDL_EMBEDDED_H

