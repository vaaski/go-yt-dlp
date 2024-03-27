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

- [Download the latest release here](https://github.com/vaaski/go-yt-dlp/releases/latest)

If you have Go 1.17+ installed, you can also use go install:

```sh
go install github.com/vaaski/go-yt-dlp@latest
```

### Dependencies

In order to function, go-yt-dlp needs [`yt-dlp`][yt-dlp], [`ffmpeg`][ffmpeg] and [`ffprobe`][ffmpeg] to be installed on your system.

Automatic installation is natively supported on Windows.

On other platforms it'll use [homebrew][brew] to install.
If you don't have [homebrew][brew] on MacOS I strongly recommend you install it.

If you wish to install them manually or have already installed them beforehand, here's how go-yt-dlp will look for them:

- Check `$PATH` or `%PATH%`
- Check for `yt-dlp`, `ffmpeg` and `ffprobe` adjacent to go-yt-dlp
- Check in `.go-yt-dlp/bin` in your home directory
- If none of the above is found, it will install them automatically

## UPDATING

To update go-yt-dlp, simply run `go-yt-dlp -U` to replace the binary with the latest release.

This will also run `yt-dlp -U` for you, updating the yt-dlp binary.

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
[ffmpeg]: https://ffmpeg.org
[brew]: https://brew.sh
[supportedsites]: https://github.com/yt-dlp/yt-dlp/blob/master/supportedsites.md
[yt-dlp installation]: https://github.com/yt-dlp/yt-dlp#installation
