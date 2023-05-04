# Greet

## A Desktop NOSTR Client (beta) 

Outstanding issues/features/missing:

- Backend currently using `nostr.Query()`, not `nostr.Subscribe()`
- Zaps, bookmarks and DMs are missing
- Still refactoring/optimisation to do

## Building

### Requires:

- [Go](https://go.dev/learn/) >= 1.18
- [Wails](https://wails.io/) 
- Run `wails doctor` to endure any OS-specific dependencies are installed

Note that Wails for Linux has a hard dependency on libwebkit2gtk-4.0. If the app does not run then check and install if necessary: `apt search libwebkit2gtk-4.0`

Steps:

```bash
git clone https://github.com/ahanniga/greet.git
cd greet
wails build
```

The binary will be in `build/bin`. 