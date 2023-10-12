<div align="center">

[![YT-DLP](https://raw.githubusercontent.com/vaaski/go-yt-dlp/.github/yt-dlp-gopher.svg)](#readme)

</div>

go-yt-dlp is a small wrapper around the excellent [yt-dlp project][yt-dlp].
It provides a simple interface to download videos from YouTube and a
[few other sites][othersites].

## SYNOPSIS

Either double-click the binary or launch it in the command line.
It will ask you to either paste a link or enter a search term.
If you enter a search term, it will search YouTube for the term and
download the first result. If you happen to have a URL in your clipboard,
it will automatically paste it for you.

This isn't very fancy, but rather an exercise for me to learn go.

## INSTALLATION

The latest commit will be built on GitHub actions.
Currently, there are only binaries for macOS and Windows,
because those are the platforms I tested on.

- [macOS arm64](https://nightly.link/vaaski/go-yt-dlp/workflows/build/main/go-yt-dlp%20darwin%20arm64.zip)
- [macOS amd64](https://nightly.link/vaaski/go-yt-dlp/workflows/build/main/go-yt-dlp%20darwin%20amd64.zip)
- [windows amd64](https://nightly.link/vaaski/go-yt-dlp/workflows/build/main/go-yt-dlp%20windows%20amd64.zip)

At the moment you'll also need to [download yt-dlp yourself][yt-dlp-download] and either put it
in your $PATH or adjacent to the go-yt-dlp binary.

## BUILDING

To build the project yourself, just install go and
execute the following commands:

- `go mod tidy`
- `go run .`

[yt-dlp]: https://github.com/yt-dlp/yt-dlp
[othersites]: https://github.com/yt-dlp/yt-dlp/blob/master/supportedsites.md
[yt-dlp-download]: https://github.com/yt-dlp/yt-dlp#installation