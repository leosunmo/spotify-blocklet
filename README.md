# spotify-blocklet

![Spotify-blocklet](/assets/screenshot.png?raw=true "Spotify-blocket")

Displays current Artist - Song. No short text form implemented yet.
Yellow if paused.

Mouse events:  
Left click - Play/Pause  
Right click - Next  
Middle click - Previous

These are easy to modify in `main.go`.

Spotify-blocklet listens to dbus events so should be ran as a daemon in i3blocks (`interval=persist`).

## Example config
i3blocks.conf
```
# Spotify player
[spotify-blocklet]
command=/usr/share/i3blocks/spotify-blocklet
interval=persist
```

## TODO
* [ ] Mouse scroll events (Volume)
* [ ] Short text (scrolling text possible?)
* [ ] Easily configured styling from i3blocks config
