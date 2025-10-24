# Distribution Guide

This guide covers packaging and distributing applications that use extractous-go.

## Understanding Dynamic Linking

extractous-go uses CGO with dynamically linked native libraries. At runtime, your application needs:

1. The Go executable
2. Platform-specific shared libraries:
   - Linux: `libextractous_ffi.so` + dependencies
   - macOS: `libextractous_ffi.dylib` + dependencies
   - Windows: `extractous_ffi.dll` + dependencies

## Recommended: Executable Directory Placement

The simplest and most reliable approach is placing libraries in the same directory as your executable.

### Linux

```bash
# Copy libraries
cp -r native/linux_amd64/lib/*.so ./dist/

# Set RPATH if needed (already configured by default)
patchelf --set-rpath '$ORIGIN' ./dist/myapp
```

Libraries are found automatically via RPATH.

### macOS

```
# Copy libraries
cp -r native/darwin_arm64/lib/*.dylib ./dist/

# Update install names for .app bundles
install_name_tool -add_rpath @executable_path ./dist/myapp
```

For .app bundles, place libraries in `Contents/Frameworks/`.

### Windows

```
:: Copy all DLLs to exe directory
copy native\windows_amd64\lib\*.dll .\dist\
```

Windows searches the executable directory first, so no additional configuration is needed.

## Framework-Specific Integration

### Wails Desktop Applications

Wails v3 uses a platform-specific build system. Modify `build/<platform>/Taskfile.yml` to copy libraries during build.

**Linux** (`build/linux/Taskfile.yml`):

```
tasks:
  copy-native-libs:
    summary: Copy extractous-go native libraries
    vars:
      LIB_PATH: '{{.ROOT_DIR}}/native/linux_{{.ARCH}}/lib'
      TARGET_DIR: '{{.ROOT_DIR}}/build/bin'
    cmds:
      - cp {{.LIB_PATH}}/*.so {{.TARGET_DIR}}/
    preconditions:
      - sh: test -d {{.LIB_PATH}}
        msg: "Run: go run github.com/rahulpoonia29/extractous-go/cmd/install@latest"

  build:
    deps:
      - task: common:build-frontend
      - task: copy-native-libs
    cmds:
      - go build -o {{.ROOT_DIR}}/build/bin/{{.APP_NAME}} {{.ROOT_DIR}}
```

**macOS** (`build/darwin/Taskfile.yml`):

```
tasks:
  copy-native-libs:
    summary: Copy extractous-go native libraries
    vars:
      LIB_PATH: '{{.ROOT_DIR}}/native/darwin_{{.ARCH}}/lib'
      TARGET_DIR: '{{.ROOT_DIR}}/build/bin/{{.APP_NAME}}.app/Contents/Frameworks'
    cmds:
      - mkdir -p {{.TARGET_DIR}}
      - cp {{.LIB_PATH}}/*.dylib {{.TARGET_DIR}}/
    preconditions:
      - sh: test -d {{.LIB_PATH}}
        msg: "Run: go run github.com/rahulpoonia29/extractous-go/cmd/install@latest"

  build:
    deps:
      - task: common:build-frontend
      - task: copy-native-libs
    # ... rest of build
```

**Windows** (`build/windows/Taskfile.yml`):

```
tasks:
  copy-native-libs:
    summary: Copy extractous-go native libraries
    vars:
      LIB_PATH: '{{.ROOT_DIR}}\native\windows_{{.ARCH}}\lib'
      TARGET_DIR: '{{.ROOT_DIR}}\build\bin'
    cmds:
      - cmd: if not exist "{{.TARGET_DIR}}" mkdir "{{.TARGET_DIR}}"
      - cmd: copy "{{.LIB_PATH}}\*.dll" "{{.TARGET_DIR}}\"
    preconditions:
      - sh: test -d {{.LIB_PATH}}
        msg: "Run: go run github.com/rahulpoonia29/extractous-go/cmd/install@latest"

  build:
    deps:
      - task: common:build-frontend
      - task: copy-native-libs
    # ... rest of build
```

### Electron Applications

For Electron apps with a Go backend:

```
// main.js
const { app } = require('electron');
const path = require('path');

// Set library path before spawning Go process
const libPath = path.join(app.getAppPath(), 'resources', 'native');
process.env.EXTRACTOUS_LIB_PATH = libPath;

// Now spawn your Go backend
```

Copy libraries during packaging:

```
{
  "build": {
    "extraResources": [
      {
        "from": "native/${os}_${arch}/lib",
        "to": "native",
        "filter": ["**/*"]
      }
    ]
  }
}
```

## Alternative: Environment Variable

If bundling libraries separately:

```
# Linux/macOS
export EXTRACTOUS_LIB_PATH=/path/to/libs
./myapp

# Windows
set EXTRACTOUS_LIB_PATH=C:\path\to\libs
myapp.exe
```

This is useful for development or when libraries are installed system-wide.

## Cross-Platform Builds

Build for multiple platforms:

```
# Install all platform libraries
go run github.com/rahulpoonia29/extractous-go/cmd/install@latest --all

# Build for each platform
GOOS=linux GOARCH=amd64 go build -o dist/myapp-linux-amd64
GOOS=darwin GOARCH=arm64 go build -o dist/myapp-darwin-arm64
GOOS=windows GOARCH=amd64 go build -o dist/myapp-windows-amd64.exe

# Copy corresponding libraries
cp native/linux_amd64/lib/*.so dist/
cp native/darwin_arm64/lib/*.dylib dist/
cp native/windows_amd64/lib/*.dll dist/
```

Note: Cross-compilation with CGO requires platform-specific toolchains.

## Installers

### Linux

For .deb packages, install libraries to `/usr/lib` or bundle with executable:

```
/opt/myapp/
├── myapp
└── lib/
    ├── libextractous_ffi.so
    └── ... (dependencies)
```

Wrapper script:

```
#!/bin/bash
DIR="$(cd "$(dirname "${BASH_SOURCE}")" && pwd)"
export LD_LIBRARY_PATH="${DIR}/lib:${LD_LIBRARY_PATH}"
exec "${DIR}/myapp" "$@"
```

### macOS

For .dmg or .pkg, bundle in `.app/Contents/Frameworks/`:

```
MyApp.app/
└── Contents/
    ├── MacOS/
    │   └── myapp
    └── Frameworks/
        ├── libextractous_ffi.dylib
        └── ... (dependencies)
```

### Windows

For .msi or NSIS installers, place DLLs alongside executable:

```
C:\Program Files\MyApp\
├── myapp.exe
├── extractous_ffi.dll
└── ... (dependencies)
```

## Troubleshooting

### "Library not found" errors

**Linux**: Verify RPATH or set `LD_LIBRARY_PATH`:
```
ldd myapp  # Check dependencies
readelf -d myapp | grep RPATH  # Verify RPATH
```

**macOS**: Verify install names:
```
otool -L myapp  # Check linked libraries
```

**Windows**: Ensure DLLs are in executable directory or PATH:
```
where extractous_ffi.dll
```

### Code Signing (macOS)

After copying dylibs, re-sign your application:

```
codesign --force --deep --sign "Developer ID" MyApp.app
```

### Permission Issues

Ensure libraries have execute permissions:

```
chmod +x native/*/lib/*
```

## Security Considerations

- Verify library integrity (checksums will be added to installer)
- Use signed libraries in production
- Don't bundle libraries in publicly writable directories
- Consider static builds for sensitive deployments (future feature)

## Testing Distribution

Test your packaged application on a clean system without development tools to ensure all dependencies are included.

```
# Create clean container for testing
docker run -it --rm -v $(pwd)/dist:/app ubuntu:latest
cd /app
./myapp  # Should work without errors
```