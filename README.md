# ilse
TUI grep tool respect for IntelliJ

## Requirements
- [ripgrep](https://github.com/BurntSushi/ripgrep)
  - for fast grep
- [bat](https://github.com/sharkdp/bat)
  - for beautiful preview

## Features
- support HeadMatch(FirstMatch), WordMatch, Regex, FuzzySearch
- preview surrounding code
- auto resize
- You can open the hit in the editor

## How to Work
![ilse](https://user-images.githubusercontent.com/31027514/107879359-b7992800-6f1b-11eb-9408-bea84deedafa.gif)

## How to Use
```command
<Ctrl-W> switch to WordMatch
<Ctrl-E> switch to HeadMatch
<Ctrl-R> switch to Regex
<Ctrl-G> switch to ripgrep
<Ctrl-F> switch to fuzzy search
<Ctrl-T> Toggle case sensitive
<Ctrl-D> specify search target directory visually
<Ctrl-N> clear search target directory
<Ctrl-B> clear your input
```

## Flag
```bash
ilse - ilse is TUI grep tool like IntelliJ

  Flags:
       --version              Displays the program version string.
    -h --help                 Displays help with available flag, subcommand, and positional value parameters.
    -m --max-search-results   Max number of search results (default: 100)
    -f --filter               select filter ('rg', 'fuzzy') (default: rg)
    -fm --filter-mode          select filter mode ('head', 'word', 'regex') (default: head)
    -t --preview-theme        select bat theme for preview (default: OneHalfDark)
    -c --case                 case sensitive
```

## Caution
I intendedly ignore if only one letter. Because, it takes a lot of time, but it's of little value.
