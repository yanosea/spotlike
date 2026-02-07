<div align="right">

![golangci-lint](https://github.com/yanosea/spotlike/actions/workflows/golangci-lint.yml/badge.svg)
![release](https://github.com/yanosea/spotlike/actions/workflows/release.yml/badge.svg)

</div>

<div align="center">

# ⚪ spotlike

![Language:Go](https://img.shields.io/static/v1?label=Language&message=Go&color=blue&style=flat-square)
![License:MIT](https://img.shields.io/static/v1?label=License&message=MIT&color=blue&style=flat-square)
[![Latest Release](https://img.shields.io/github/v/release/yanosea/spotlike?style=flat-square)](https://github.com/yanosea/spotlike/releases/latest)
<br/>
[Coverage Report](https://yanosea.github.io/spotlike/coverage.html)
<br/>
![demo](docs/demo.gif "demo")

</div>

## ℹ️ About

`spotlike` is the CLI tool to LIKE contents in Spotify.
This tool uses [Spotify Web API](https://developer.spotify.com/documentation/web-api) with [Go wrapper library](https://github.com/zmb3/spotify).

## 💻 Usage

```
Available Commands:
  auth,       au,   a  🔑 Authenticate your Spotify client.
  get,        ge,   g  📚 Get the information of the content on Spotify by ID.
  like,       li,   l  🤍 Like content on Spotify by ID.
  unlike,     un,   u  💔 Unlike content on Spotify by ID.
  search,     se,   s  🔍 Search for the ID of content in Spotify.
  completion, comp, c  🔧 Generate the autocompletion script for the specified shell.
  version,    ver,  v  🔖 Show the version of spotlike.
  help                 🤝 Help for spotlike.

Flags:
  -h, --help     🤝 help for spotlike
  -v, --version  🔖 version for spotlike
```

### 🔍 search

Search for the ID of content in Spotify.

```
Flags:
  -A, --artist  🎤 search for artists
  -a, --album   💿 search for albums
  -t, --track   🎵 search for tracks
  -m, --max     🔢 maximum number of search results (default 10)
  -f, --format  📝 format of the output (default "table", e.g: "plain")
  -h, --help    🤝 help for search

Arguments:
  keywords  🔡 search content by keywords (multiple keywords are separated by a space)
```

### 🤍 like

Like content on Spotify by ID.

#### 🤍🎵 like track

```
Flags:
  -A, --artist  🆔 an ID of the artist to like all albums released by the artist
  -a, --album   🆔 an ID of the album to like all tracks in the album
  --no-confirm  🚫 do not confirm before liking the track
  -f, --format  📝 format of the output (default "table", e.g: "plain")
  -h, --help    🤝 help for track

Arguments:
  ID  🆔 ID of the tracks (e.g: "20q73dOrP7ceLGAJQVtuTq 5A7nqzXUt5IZIOA7oNBv6M")
```

#### 🤍💿 like album

```
Flags:
  -A, --artist  🆔 an ID of the artist to like all albums released by the artist
  --no-confirm  🚫 do not confirm before liking the album
  -f, --format  📝 format of the output (default "table", e.g: "plain")
  -h, --help    🤝 help for album

Arguments:
  ID  🆔 ID of the albums (e.g: "1dGzXXa8MeTCdi0oBbvB1J 6Xy481vVb9vPK4qbCuT9u1")
```

#### 🤍🎤 like artist

```
Flags:
  --no-confirm  🚫 do not confirm before liking the artist
  -f, --format  📝 format of the output (default "table", e.g: "plain")
  -h, --help    🤝 help for artist

Arguments:
  ID  🆔 ID of the artists (e.g: "00DuPiLri3mNomvvM3nZvU 3B9O5mYYw89fFXkwKh7jCS")
```

### 💔 unlike

Unlike content on Spotify by ID.
Subcommands and flags are the same as the `like` command.

### 📚 get

Get the information of the content on Spotify by ID.

#### 📚💿 get albums

```
Flags:
  -f, --format  📝 format of the output (default "table", e.g: "plain")
  -h, --help    🤝 help for albums

Argument:
  ID  🆔 ID of the albums (e.g: "1dGzXXa8MeTCdi0oBbvB1J")
```

#### 📚🎵 get tracks

```
Flags:
  -f, --format  📝 format of the output (default "table", e.g: "plain")
  -h, --help    🤝 help for tracks

Argument:
  ID  🆔 ID of the artist or album (e.g: "00DuPiLri3mNomvvM3nZvU")
```

## 📝 Preparation

1. Login [Spotify Developer](https://developer.spotify.com).
2. Open [Dashboard](https://developer.spotify.com/dashboard).
3. Press `Create App` button and input below.
    1. `App name` (e.g. `spotlike`)
    2. `App description` (e.g. `spotlike is a CLI tool to LIKE the contents in Spotify.`)
    3. `Redirect URIs` (recommended: `http://localhost:8080/callback`)
    4. `Which API/SDKs are you planning to use` (check `Web API`)
4. Open `Basic Information` from created app in Dashboard.
5. Get `Client ID` and `Client secret`.
6. Set environments below.
    1. `SPOTIFY_ID`
    2. `SPOTIFY_SECRET`
    3. `SPOTIFY_REDIRECT_URI`
7. Now, you're ready for authenticate in `spotlike`!

## 🌍 Environments

### 🆔 Spotify client ID

```sh
export SPOTIFY_ID=your_client_id
```

### 🔑 Spotify client secret

```sh
export SPOTIFY_SECRET=your_client_secret
```

### 🔗 Spotify redirect URI

Default : `http://localhost:8080/callback`

```sh
export SPOTIFY_REDIRECT_URI=http://localhost:8080/callback
```

### 🔄 Spotify refresh token

This is automatically obtained after running `spotlike auth`.

```sh
export SPOTIFY_REFRESH_TOKEN=your_refresh_token
```

## 🔧 Installation

### 🐭 Using go

```sh
go install github.com/yanosea/spotlike/app/presentation/cli/spotlike@latest
```

### 🍺 Using homebrew

```sh
brew tap yanosea/tap
brew install yanosea/tap/spotlike
```

### 📦 Download from release

Go to the [Releases](https://github.com/yanosea/spotlike/releases) and download the latest binary for your platform.

## ✨ Update

### 🐭 Using go

Reinstall `spotlike`!

```sh
go install github.com/yanosea/spotlike/app/presentation/cli/spotlike@latest
```

### 🍺 Using homebrew

```sh
brew update
brew upgrade spotlike
```

### 📦 Download from release

Download the latest binary from the [Releases](https://github.com/yanosea/spotlike/releases) page and replace the old binary in your `$PATH`.

## 🧹 Uninstallation

### 🐭 Using go

```sh
rm $GOPATH/bin/spotlike
# maybe you have to execute with sudo
rm -fr $GOPATH/pkg/mod/github.com/yanosea/spotlike*
```

### 🍺 Using homebrew

```sh
brew uninstall spotlike
brew untap yanosea/tap/spotlike
```

### 📦 Download from release

Remove the binary you downloaded and placed in your `$PATH`.

## 📃 License

[🔓MIT](./LICENSE)

## 🖊️ Author

[🏹 yanosea](https://github.com/yanosea)

## 🔥 Motivation

- Spotify's smartphone app or web app does not have the way below.
    - LIKE all tracks in one album.
    - LIKE all albums from one artist.
- I wanted to LIKE them at once, so I created it!!

## 🤝 Contributing

Feel free to point me in the right direction🙏
