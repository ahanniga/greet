# Greet

## A Desktop NOSTR Client (beta) 

Outstanding issues/features/missing:

- Backend currently using `nostr.Query()`, not `nostr.Subscribe()`
- Zaps, bookmarks, profile editing are missing
- Much refactoring/optimisation to do

## Building

Requires:

- [Go](https://go.dev/learn/) >= 1.18
- [Wails](https://wails.io/) 


Steps:

```bash
git clone https://github.com/ahanniga/greet.git
cd greet
wails build
```

The binary will be in `build/bin`. 