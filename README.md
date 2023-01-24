Binds Alt+Tab in order to focus back and forth the last two visited workspaces.

```
~/.config/i3/config

exec --no-startup-id i3-alt-tab
bindsym Mod1+Tab exec --no-startup-id killall -SIGUSR1 i3-alt-tab
```
