// createprocess_capture.c
// CreateProcess example that redirects stdout/stderr of October.exe and captures output
#include <windows.h>
#include <stdio.h>

int main(void) {
    SECURITY_ATTRIBUTES sa;
    sa.nLength = sizeof(sa);
    sa.lpSecurityDescriptor = NULL;
    sa.bInheritHandle = TRUE;

    HANDLE readPipe = NULL, writePipe = NULL;
    if (!CreatePipe(&readPipe, &writePipe, &sa, 0)) {
        printf("CreatePipe failed: %lu\n", GetLastError());
        return 1;
    }
    // Ensure the read handle is not inherited
    SetHandleInformation(readPipe, HANDLE_FLAG_INHERIT, 0);

    STARTUPINFOA si;
    PROCESS_INFORMATION pi;
    ZeroMemory(&si, sizeof(si));
    si.cb = sizeof(si);
    si.dwFlags = STARTF_USESTDHANDLES;
    si.hStdOutput = writePipe;
    si.hStdError = writePipe;
    si.hStdInput = NULL;

    BOOL ok = CreateProcessA(
        ".\\October.exe",
        NULL, NULL, NULL, TRUE, 0, NULL, NULL, &si, &pi);

    CloseHandle(writePipe);
    if (!ok) {
        CloseHandle(readPipe);
        printf("CreateProcess failed: %lu\n", GetLastError());
        return 1;
    }

    CHAR buffer[4096]; DWORD read;
    while (ReadFile(readPipe, buffer, sizeof(buffer)-1, &read, NULL) && read) {
        buffer[read] = '\0';
        printf("%s", buffer);
    }

    WaitForSingleObject(pi.hProcess, INFINITE);
    CloseHandle(pi.hThread);
    CloseHandle(pi.hProcess);
    CloseHandle(readPipe);
    return 0;
}
