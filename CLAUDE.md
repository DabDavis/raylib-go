# Raylib-Go Fork (EQOA)

Local fork of `gen2brain/raylib-go` for the EQOA client project. Keeps the upstream module path (`github.com/gen2brain/raylib-go/raylib`) so the EQOA client uses a `replace` directive — no import path changes needed.

## Why Fork?

The upstream Go bindings have wrapper bugs and missing features. Forking lets us fix issues and add GLFW access without waiting for upstream PRs or maintaining complex workarounds.

## Client Integration

```
# go-eqoa-client/go.mod
replace github.com/gen2brain/raylib-go/raylib => ../raylib-go/raylib
```

## Changes From Upstream

### Bug Fixes

#### `UpdateAudioStream` frame count (raudio.go:491)
- **Issue**: C API `UpdateAudioStream(stream, data, frameCount)` expects frame count, but Go wrapper passes `len(data)` (element count). For stereo, this is 2x the correct value.
- **Status**: DOCUMENTED but NOT fixed. The upstream behavior happens to work because raylib internally clamps to buffer size. Fixing it (dividing by channels) causes silence — raylib's internal buffer sizing depends on the inflated count.
- **Workaround**: Use `LoadMusicStreamFromMemory` for streaming audio instead of manual `AudioStream` + `UpdateAudioStream`. The `Music` type handles double-buffering correctly.

### New Features

#### `GetVideoModes()` (videomode.go)
- Queries `glfwGetVideoModes` for all display modes supported by the primary monitor
- Returns `[]VideoMode{Width, Height, RefreshRate}`
- Deduplication and filtering done by the caller (EQOA client's graphics settings)
- Uses raylib's bundled GLFW headers (no system GLFW dependency)

## Known Upstream Issues (Not Yet Fixed)

| Issue | File | Description |
|-------|------|-------------|
| `UpdateAudioStream` count | `raudio.go` | Passes element count not frame count (see above) |
| `MatrixRotateY` sign swap | `raymath.go` | M2/M8 signs swapped vs C raylib — EQOA works around this with manual matrix construction |

## File Structure

```
raylib/
  raudio.go          — Audio bindings (UpdateAudioStream documented)
  videomode.go       — NEW: GetVideoModes via GLFW
  rcore.go           — Window/input bindings
  rlgl_cgo.go        — OpenGL state bindings (used by EQOA's gfx package)
  raylib.go          — Core types (AudioStream, Music, Vector3, etc.)
  external/glfw/     — Bundled GLFW source (videomode.go uses this)
```

## Updating From Upstream

1. Add upstream remote: `git remote add upstream https://github.com/gen2brain/raylib-go.git`
2. Fetch: `git fetch upstream`
3. Merge: `git merge upstream/master`
4. Resolve conflicts in our changed files (`raudio.go`, `videomode.go`)
5. Test EQOA client builds: `cd ../go-eqoa-client && go build ./cmd/client/`

## Rules

- Keep module path as `github.com/gen2brain/raylib-go/raylib` — never change it
- Document all changes from upstream in this file
- Test EQOA client after any merge from upstream
- Prefer adding new files over modifying upstream files (easier merges)
