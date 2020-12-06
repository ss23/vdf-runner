# vdf-runner
Tool for running Steam/SteamWorks `.vdf` install files manually

## Notes
This is very much an alpha tool, and will likely require manual modification to be used. As an example, if you want to use it on a real `vdf` file, you'll need to modify and set values in `main.go`, such that variables like the installation directory are correctly replaced.

Additionally, it does not yet run the processes required, however these are trivial to run manually, as opposed to manually adding registry keys which can be tedious.

## Why?
For some reason, Steam refuses to run `vdf` install files from a network drive on my computer, so I am using this to run them manually. For some games, there are few registry keys, so manually is fine, but for some, writing this tool is faster than adding 50 keys by hand.

Another reason is to help with older game compatibility. For example, [Spore](https://en.wikipedia.org/wiki/Spore_%282008_video_game%29) requires modifications to the `vdf` before it can be used, namely moving all registry keys to `HKLM/Software/Wow6432Node` (instead of the configured `HKLM/Software` location), and this can be done trivially with this tool.

## Roadmap
The roadmap is I will accept PRs from people if they want, but I have no plans to work on this tool more than is required for me to run games as I need them.
