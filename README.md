# make-webp

A small CLI that converts JPEG, PNG, and GIF images to WebP, optionally resizing them first. Encoding is done by Google's `cwebp` binary.

## Requirements

- Go 1.27+
- `cwebp` (`brew install webp`). By default the tool looks for it in `/opt/homebrew/bin`; override with `-cwebp-path`.

## Install

```sh
go install github.com/jimmysawczuk/make-webp@latest
```

## Usage

```sh
make-webp [flags] <file-or-glob>...
```

Arguments are file paths or glob patterns (Go `filepath.Glob` syntax; `**` is not supported). Quote globs if you want the tool, rather than your shell, to expand them.

### Flags

| Flag | Default | Description |
| --- | --- | --- |
| `-q` | `90` | WebP quality (0–100). |
| `-width` | `0` | Max width in pixels. `0` means unconstrained. |
| `-height` | `0` | Max height in pixels. `0` means unconstrained. |
| `-out` | input's directory | Directory to write `.webp` files to. Must already exist. |
| `-safe` | `true` | Replace whitespace with `-` and strip characters other than letters, digits, `_`, and `-` from the output filename. |
| `-with-dimensions` | `true` | Append the requested size (e.g. `-800x`, `-x600`, `-800x600`) to the output filename when resizing. |
| `-lowercase-filenames` | `false` | Lowercase the output filename. |
| `-cwebp-path` | `/opt/homebrew/bin` | Directory containing the `cwebp` binary. |

Boolean flags that default to `true` are disabled with `=false`, e.g. `-safe=false`.

### Examples

Convert every PNG in a directory at the default quality:

```sh
make-webp 'images/*.png'
```

Resize to fit within 1200px wide, quality 80, into `dist/`:

```sh
make-webp -width 1200 -q 80 -out dist 'photos/*.jpg'
# photos/My Photo.jpg -> dist/My-Photo-1200x.webp
```

Fit within a 800×600 box without the dimension suffix:

```sh
make-webp -width 800 -height 600 -with-dimensions=false hero.png
# hero.png -> hero.webp
```

## Behavior notes

- Resizing preserves aspect ratio (fit within the box, Lanczos filter) and never upscales.
- The filename suffix reflects the *requested* bounds, not the actual output dimensions.
- Existing output files are overwritten.
- A file that fails to decode or encode is logged and skipped; the rest continue.
