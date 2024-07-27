# Pixxy

A really sweet image manipulation tool for creating aesthetic images on the command line.

---

you can apply any arbitrary set of colors to an image, in this case I used a theme file for the Kitty terminal.
`pixxy` will automatically parse this file for any valid hex colors and apply them to the image. You can use
any file (or string) that contains hex colors. This can be useful for creating themed wallpapers, you could pick
any Gruvbox config file and the image will match those colors.

```sh
pixxy glitch \
    --gif \
    --verbose \
    --input ~/Pictures/cowgirl-thumbnail.jpg \
    --output output.gif \
    --seed sweet \
    --threshold 1.0 \
    --palette-file ~/.config/kitty/themes/ayanami-cold.conf \
```

output:
Input | Output  
:-------------------------:|:-------------------------:
![image of cowgirl](./assets/cowgirl-thumbnail.jpg)|![image of cowgirl glitched as a gif](./assets/cowgirl-glitch.gif)

another example using the terminal text-editor, Helix, Gruvbox theme file:
![hatsune miku remixed with Gruvbox](./assets/screenshot.png)

## Dither

examples

```sh
# simple dither
pix dither           \
  --input input.png  \
  --output out.png   \
  --dither floyd

# use dithering and ordered dither matrices (ie "dots" and manga like toning)
pix dither              \
  --input input.png     \
  --output out.png      \
  --dither floyd        \
  --ordered vertical5x3

# layer multiple dither filters
pix dither              \
  --input input.png     \
  --output out.png      \
  -d floyd -d atkinsons \
  -m vertical5x3        \
  --bayer

# random ordered matrix of 0.3 strength and using a color palette from a random file containing hex colors
pix dither \
--input input.png \
--output out.png \
--palette-file ~/.config/kitty/themes/pretty-pink.conf \
-t 0.3 \
-m $(/pix dither --list-maps | shuf | head -n1) \
--scale --scale-factor 5
```

# Wallpaper-finder

find wallpaper sized images!

```sh
wallpaper-finder -e jpg -e jpeg -e png ~/Pictures ~/images | \
    fzf --preview='kitty icat --clear --transfer-mode=memory --stdin=no --place=${FZF_PREVIEW_COLUMNS}x${FZF_PREVIEW_LINES}@0x0 {}'
```

```sh
Usage:
  wallpaper-finder [OPTIONS]

Application Options:
  -e, --extension= an array of extensions to search for ie (-e png -e jpg)
  -d, --directory= an array of paths to search, can specify more than one
  -r, --ratio=     a ratio to search, in the format <widht>x<height> (16x9)
  -t, --tolerance= percentage of tolerance for the ratio [5%] (default: 5)
  -c, --color      print paths with color
  -f, --follow     follow symlinks
  -v, --verbose    print debugging information and verbose output

Help Options:
  -h, --help       Show this help message

```

# Acknowledgements

[darkliquid/glitch](https://github.com/darkliquid/glitch) MIT license

- the base glitch library I built off of

[koyachi/go-atkinson](https://github.com/koyachi/go-atkinson) MIT licensed

- atkinsons dither function
