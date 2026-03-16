# Raylib-Go Binding Analysis

Comprehensive analysis of the Go bindings for raylib v5.6-dev, covering 631 functions across 12 wrapper files.

## File & Function Count

| File | Functions | Purpose |
|------|-----------|---------|
| rcore.go | 153 | Window, input (keyboard/mouse/gamepad), drawing modes, shaders, camera |
| rlgl_cgo.go | 141 | OpenGL state, matrix ops, vertices, textures, framebuffers |
| rtextures.go | 94 | Image/texture loading, manipulation, conversion |
| rmodels.go | 71 | 3D model loading, drawing, mesh, materials, collision |
| rshapes.go | 67 | 2D/3D shape drawing |
| raudio.go | 61 | Audio loading, streaming, wave manipulation, callbacks |
| rtext.go | 22 | Font/text loading, drawing, measurement |
| rcamera.go | 14 | Camera control helpers |
| rgestures.go | 8 | Touch gesture detection |
| videomode.go | 1 | **EQOA addition**: GetVideoModes via GLFW |

## Known Bugs

### UpdateAudioStream Frame Count (raudio.go:491)

C API: `void UpdateAudioStream(AudioStream stream, const void *data, int frameCount)`
Go wrapper passes `len(data)` (element count) instead of frame count. For stereo, element count = frames x 2.

**Status**: Documented, deliberately NOT fixed. Fixing causes silence because raylib's internal buffer sizing depends on the inflated count.

**Workaround**: Use `LoadMusicStreamFromMemory` for streaming audio.

### MatrixRotateY Sign Swap (raymath.go)

Go MatrixRotateY has M2/M8 signs swapped vs C raylib. Causes inverted Y-axis rotation.

**Status**: Not fixed. EQOA constructs matrices manually.

## Missing rlgl Functions

Functions available in C `rlgl.h` but not wrapped in `rlgl_cgo.go`:

| Function | Purpose | EQOA Impact |
|----------|---------|-------------|
| `glDepthFunc` | Depth test function (GL_LESS, GL_LEQUAL) | Worked around via `internal/gfx/gl.go` |
| `glPolygonOffset` | Depth bias for overlapping polygons | Worked around via `internal/gfx/gl.go` |
| `glClear(DEPTH_BIT)` | Clear depth buffer only | Worked around via `internal/gfx/gl.go` |
| `rlLoadFramebuffer` | Create empty FBO | Not needed yet |
| `rlUnloadFramebuffer` | Delete FBO | Not needed yet |
| `rlUpdateTexture` | Update texture region on GPU | Not needed yet |
| `rlLoadTexture` | Load raw texture data to GPU | Not needed yet |
| `rlPolygonMode` | Wireframe/fill per face | EnableWireMode covers our use |

EQOA wraps the missing GL calls in `internal/gfx/gl.go` using direct CGO to system GL headers.

## Audio System

### Thread Safety
- Audio callbacks run on miniaudio's audio thread
- Go callback wrapper uses `unsafe.Slice` for buffer access
- Mixed processor callbacks protected with RWMutex
- `UpdateAudioStream` locks via C mutex internally

### Buffer Management
- Go slices passed as `unsafe.Pointer(&data[0])` to C
- `LoadWaveSamples` returns unsafe slice (caller must not let GC collect backing data)
- Music type handles double-buffering internally — preferred for streaming

### EQOA Audio Architecture
- **BGM**: `LoadMusicStreamFromMemory(".wav", wavData)` — raylib streams internally
- **SFX**: `LoadSoundFromWave(wave)` — pre-decoded PCM cached in memory
- **UpdateMusicStream** called each frame from render loop

## OpenGL ES Compatibility

All 141 wrapped rlgl functions work with ES2/ES3. Activated via build tag:
```bash
go build -tags es2
```

EQOA's 6 direct GL calls (in `internal/gfx/gl.go`) use ES-compatible functions:
- `glDepthFunc(GL_LEQUAL/GL_LESS)` — ES2+
- `glPolygonOffset(-1, -1)` — ES2+
- `glClear(GL_DEPTH_BUFFER_BIT)` — ES2+
- `glEnable/glDisable(GL_POLYGON_OFFSET_FILL)` — ES2+

## Window Management

### Fullscreen Modes
- `ToggleFullscreen()` — exclusive, changes monitor resolution (restores on exit)
- `ToggleBorderlessWindowed()` — borderless, keeps monitor resolution

### Monitor Queries
- `GetMonitorCount/Width/Height/RefreshRate/Name/Position` — full multi-monitor support
- `GetCurrentMonitor()` — which monitor the window is on
- `GetWindowScaleDPI()` — HiDPI awareness
- `GetRenderWidth/Height()` — frame buffer size (may differ from screen on HiDPI)

### EQOA Window Modes
- **Windowed**: `SetWindowSize` at game resolution
- **Borderless**: `SetWindowSize` then `ToggleBorderlessWindowed`
- **Fullscreen**: `SetWindowSize` then `ToggleFullscreen` with focus-loss auto-restore

## Input System

### Gamepad
- SDL2-based controller mappings (SDL_GameControllerDB format)
- `SetGamepadMappings(string)` for custom controllers
- `SetGamepadVibration(gamepad, leftMotor, rightMotor, duration)` — haptic feedback
- 6 axes (2 sticks + 2 triggers), ~17 buttons per controller

### Keyboard
- `IsKeyPressed` (one-frame), `IsKeyDown` (held), `IsKeyPressedRepeat` (OS repeat)
- `GetKeyPressed()` / `GetCharPressed()` — queue-based input

### Mouse
- Button states, position, delta, wheel (X/Y)
- Cursor control (show/hide/lock)
- Custom cursor types

## Build System

### Platform Compilation

| Platform | Backend | C Sources | Link Libraries |
|----------|---------|-----------|----------------|
| Linux | GLFW/X11 | 25 GLFW + 7 raylib | -lm -pthread -ldl -lrt |
| Linux | GLFW/Wayland | 25 GLFW + 7 raylib | + -lwayland-client/cursor/egl |
| Windows | GLFW/Win32 | 20 GLFW + 7 raylib | -lgdi32 -lwinmm -lopengl32 |
| macOS | GLFW/Cocoa | 20 GLFW + 7 raylib | Cocoa, IOKit, CoreVideo, OpenGL frameworks |

### Graphics Backend Tags
```
(default)     → OpenGL 3.3
-tags es2     → OpenGL ES 2.0
-tags es3     → OpenGL ES 3.0
-tags opengl43 → OpenGL 4.3
```

### Cross-Compilation
```bash
# Linux → Windows (MinGW)
CC=x86_64-w64-mingw32-gcc GOOS=windows go build

# Linux → Windows (Zig)
CC="zig cc -target x86_64-windows-gnu" GOOS=windows go build

# EQOA uses Zig for glibc 2.31 compat:
CC="zig cc -target x86_64-linux-gnu.2.31" go build
```

### Bundled Libraries
- **miniaudio** — single-header audio (WAV, OGG, MP3, FLAC)
- **stb_image** — image loading (PNG, JPG, BMP, TGA, GIF, PSD)
- **stb_truetype** — TrueType font rendering
- **stb_image_resize** — image scaling
- **GLFW 3.4** — window/input/context management

## Version

**raylib v5.6-dev** (development build)
- Forked from gen2brain/raylib-go master
- 2 custom commits on top of upstream
