<div align="center">
<h1>go-yt-dlp</h1>

<br />

<img alt="go-yt-dlp" width="300" src="https://raw.githubusercontent.com/vaaski/go-yt-dlp/main/.github/yt-dlp-gopher.svg" />

</div>

<br />

> go-yt-dlp is a small wrapper around the excellent [yt-dlp project][yt-dlp].
> It provides a simple interface to download videos from YouTube and a
> [few other sites][supportedsites].

## USAGE

Either double-click the binary or launch it in the command line.
It will ask you to either paste a link or enter a search term.
If you enter a search term, it will search YouTube and download the
first result. If you happen to have a URL in your clipboard,
it will automatically pre-fill it for you.

At the query prompt, you can also press <kbd>Tab</kbd> to enable
searching on YouTube Music only.

All this isn't particularly fancy or the most efficient, but rather serves
as an exercise for me to learn go.

## INSTALLATION

The latest commit will be built on GitHub actions.
Currently, there are only binaries for macOS and Windows,
because those are the platforms I tested on. It should
probably run on other platforms as well.

- [macOS arm64][macOS arm64]
- [macOS amd64][macOS amd64]
- [windows amd64][windows amd64]

At the moment you'll also need to [download yt-dlp yourself][yt-dlp installation]
and either put it in your $PATH or adjacent to the go-yt-dlp binary.

## RUNNING FROM SOURCE

To run the project from source, just install go and
execute the following commands:

- `go mod tidy`
- `go run .`

## BUILD FROM SOURCE

- Install go 1.21 or higher
- [Install goreleaser](https://goreleaser.com/install/#go-install)
- [Install go-winres](https://github.com/tc-hib/go-winres#installation)
- Clone this repo
- Run `go generate`

[yt-dlp]: https://github.com/yt-dlp/yt-dlp
[supportedsites]: https://github.com/yt-dlp/yt-dlp/blob/master/supportedsites.md
[yt-dlp installation]: https://github.com/yt-dlp/yt-dlp#installation
[macOS arm64]: https://nightly.link/vaaski/go-yt-dlp/workflows/build/main/go-yt-dlp%20darwin%20arm64.zip
[macOS amd64]: https://nightly.link/vaaski/go-yt-dlp/workflows/build/main/go-yt-dlp%20darwin%20amd64.zip
[windows amd64]: https://nightly.link/vaaski/go-yt-dlp/workflows/build/main/go-yt-dlp%20windows%20amd64.zip
