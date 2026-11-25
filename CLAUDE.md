# CLAUDE.md

Single-file Go CLI (`main.go`) that converts JPEG/PNG/GIF to WebP, with optional fit-within-box resizing.

## Commands

```sh
go build -o make-webp .
go vet ./...
go run . -width 800 'path/to/*.png'   # needs cwebp installed
```

There are no tests.

## Architecture

- `main` parses flags into package-level vars, expands args via `getPossiblePaths` (`filepath.Abs` + `filepath.Glob`), then calls `processImage` per file.
- `processImage`: decode with `kovidgoyal/imaging` → optional `fitImage` (Lanczos, no upscaling) → encode by shelling out to `cwebp` via `nickalie/go-webpbin`.
- `outputFilename` builds the output name from the input basename plus the `-safe`, `-lowercase-filenames`, and `-with-dimensions` flags.

## Conventions and gotchas

- Flags are package-level globals read directly by helpers; keep it that way unless the tool grows enough to justify a config struct.
- `cwebp` is an external binary. `-cwebp-path` is the *directory* containing it (default `/opt/homebrew/bin`), passed to `webpbin.SetVendorPath`.
- `kovidgoyal/imaging` is a maintained fork of `disintegration/imaging`; don't switch back.
- Input decoders are registered via blank imports (`image/gif`, `image/jpeg`, `image/png`). Adding a format means adding its decoder import.
- Per-file errors are logged and skipped. Failing to open the output file currently panics.
- Dependency bumps are their own commits with `build:` prefixed messages.
