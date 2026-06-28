// createprocess.c
// Simple CreateProcess example to launch October.exe and wait for it
#include <windows.h>
#include <stdio.h>

int main(void) {
    STARTUPINFOA si;
    PROCESS_INFORMATION pi;
    ZeroMemory(&si, sizeof(si));
    si.cb = sizeof(si);
    ZeroMemory(&pi, sizeof(pi));

    // Adjust path as necessary
    // Prefer a relative path so the example works from a checked-out repo or build directory
    LPCSTR app = ".\\October.exe"; // ensure working directory is the repo root or adjust as needed

    BOOL ok = CreateProcessA(
        app,   // application name
        NULL,  // command line
        NULL, NULL, FALSE, 0, NULL, NULL, &si, &pi);
    if (!ok) {
        printf("CreateProcess failed: %lu\n", GetLastError());
        return 1;
    }

    // Wait for process to exit
    WaitForSingleObject(pi.hProcess, INFINITE);

    CloseHandle(pi.hThread);
    CloseHandle(pi.hProcess);
    return 0;
}
