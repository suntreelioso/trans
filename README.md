# Trans is a simple file transfer tool

```
trans - a simple file transfer tool

Usage:
  trans [flags]
  trans [command]

Available Commands:
  help        Help about any command
  server      server mode

Flags:
  -a, --addr string   server address (default "127.0.0.1:8080")
  -g, --get strings   get one or more files, e.g. -g 1.txt,2.mp3
  -G, --get-all       get all files
  -h, --help          help for trans
  -l, --list          list all files (default true)
  -p, --path string   path to share, default: current working directory
  -v, --version       version for trans

Use "trans [command] --help" for more information about a command.
```

### Server

```bash
$ trans server -p ~/Downloads
```

### Client

List files

```bash
$ trans -l
files on remote server:
  1  614.4K  data.zip
  2    2.6M  music.mp3
  3    233B  one.txt
```

Download File

```bash
$ ./trans -g one.txt,music.mp3 -p /tmp
2025/06/18 13:34:56 downloading one.txt
2025/06/18 13:34:56 downloaded file to /tmp/one.txt
2025/06/18 13:34:56 downloading music.mp3
2025/06/18 13:34:56 downloaded file to /tmp/music.mp3
```
