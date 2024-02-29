# Justopen

Just opens the **** location.

## Config
Config location: `~/.config/justopen/config.yaml`

Example content:
```yaml
# This is false by default, which means it'll convert the path to lowercase
casesensitive: false
# This list gets evaluated in order
filetypes:
  # Locations starting with http:// or https:// open with firefox
  - { regex: '^https?://', exec: firefox -- %f }
  # Starting with ftp:// -> open with gftp
  - { prefix: 'ftp://', exec: gftp -- %f }
  # Open all kinds of extensions with nvim
  - { regex: '\.(ac|c|cc|cpp|css|cxx|diff|go|h|hpp|in|java|js|json|jsx|mk|pl|py|rb|rs|xml|zig)$', exec: nvim -- %f }
  # If `exectty` exists, `exectty` gets run when in a terminal, and `exec` otherwise.
  # If `exectty` is missing, it uses `exec` in both cases.
  # Open playlists with 'mpv --playlist=...' if inside a terminal,
  # or if not in a terminal open the st terminal with the same command
  - { regex: '\.(m3u|pls)$', exec: st -e mpv --playlist=%f, exectty: mpv --playlist=%f }
  # Same but for non-playlist audio/video files
  - { regex: '\.(aac|flac|m4a|mp3|mpeg3|ogg|wav)$', exec: st -e mpv -- %f, exectty: mpv -- %f }
  - { regex: '\.(avi|flv|m4v|mkv|mov|mp4|mpeg|mpg|ogm|ogv|ts|vob|wmv)$', exec: mpv -- %f }
  # Open certain image extensions with gimp
  - { regex: '\.(xcf|tif|tiff|xbm|xpm)$', exec: gimp -- %f }
  # Open other ones with qiv
  - { regex: '\.(bmp|gif|jpeg|jpg|png|svg)$', exec: qiv --readonly --transparency --autorotate --scale_down -- %f }
  # You can also use suffix for single extensions
  - { suffix: '.pdf', exec: zathura -- %f }
  - { suffix: '.dia', exec: dia -- %f }
  # HTML and various office files
  - { regex: '\.(htm|html)$', exec: firefox -- %f }
  - { regex: '\.(odb|xlsx|xls|xlt|ods|ots|sxc|sdc|odg|otg|sxd|sda|vsd|vss|vst|odp|otp|sxi|pptx|ppt|pps|pot|sdi|odf|sxm|smf|ott|odt|sxw|rtf|sdw|wbk|doc|docx|dot|wri)$', exec: loffice %f }
  # Directories get the `inode/directory` mimetype
  - { mime: '^inode/directory$', exec: xfe -- %f }
  # fallback
  - { mime: '^text/.*', exectty: nano -- %f, exec: leafpad -- %f }
  # This would catch everything
  # - { regex: '.', exec: leafpad }
```

