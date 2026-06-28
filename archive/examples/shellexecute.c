// shellexecute.c
// Launch October.exe via ShellExecuteEx and wait for process handle
#include <windows.h>
#include <stdio.h>

int main(void) {
    SHELLEXECUTEINFOA sei;
    ZeroMemory(&sei, sizeof(sei));
    sei.cbSize = sizeof(sei);
    sei.fMask = SEE_MASK_NOCLOSEPROCESS;
    // Prefer relative path for portability
    sei.lpFile = ".\\October.exe"; // adjust working dir as needed
    sei.nShow = SW_SHOWNORMAL;

    if (!ShellExecuteExA(&sei)) {
        printf("ShellExecuteEx failed: %lu\n", GetLastError());
        return 1;
    }

    if (sei.hProcess) {
        WaitForSingleObject(sei.hProcess, INFINITE);
        CloseHandle(sei.hProcess);
    }
    return 0;
}
