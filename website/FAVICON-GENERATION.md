# VibeSQL Favicon Generation

## Source Icon

`vibesql-icon.svg` - Simple database cylinder in VibeSQL Blue (#336791)

## Generate Favicons

### Option 1: Online Tool (Fastest)

1. Go to https://realfavicongenerator.net/
2. Upload `vibesql-icon.svg`
3. Download package
4. Extract to `website/` folder

### Option 2: ImageMagick (Command Line)

```bash
# Install ImageMagick first
# Then run:

cd website

# Generate PNGs
magick vibesql-icon.svg -resize 16x16 favicon-16x16.png
magick vibesql-icon.svg -resize 32x32 favicon-32x32.png
magick vibesql-icon.svg -resize 96x96 favicon-96x96.png
magick vibesql-icon.svg -resize 180x180 apple-touch-icon.png
magick vibesql-icon.svg -resize 192x192 android-chrome-192x192.png
magick vibesql-icon.svg -resize 512x512 android-chrome-512x512.png

# Generate ICO (multi-size)
magick favicon-16x16.png favicon-32x32.png favicon.ico
```

### Option 3: Use PayEz Quality Standards

Copy structure from `E:\Repos\api.payez.net\public\`:
- favicon-16x16.png
- favicon-32x32.png
- apple-touch-icon.png
- android-chrome-192x192.png
- android-chrome-512x512.png
- site.webmanifest

## Files Needed for Website

```
website/
├── vibesql-icon.svg          # Source (done)
├── vibesql-logo.svg          # Logo wordmark (done)
├── favicon.ico               # Multi-size ICO
├── favicon-16x16.png         # 16×16 PNG
├── favicon-32x32.png         # 32×32 PNG
├── favicon-96x96.png         # 96×96 PNG (optional)
├── apple-touch-icon.png      # 180×180 PNG
├── android-chrome-192x192.png # 192×192 PNG
├── android-chrome-512x512.png # 512×512 PNG
└── site.webmanifest          # PWA manifest
```

## HTML Head (add to index.html)

```html
<link rel="icon" type="image/svg+xml" href="/vibesql-icon.svg">
<link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png">
<link rel="icon" type="image/png" sizes="16x16" href="/favicon-16x16.png">
<link rel="apple-touch-icon" sizes="180x180" href="/apple-touch-icon.png">
<link rel="manifest" href="/site.webmanifest">
```

## Current Status

✅ `vibesql-icon.svg` - Minimal database cylinder icon
✅ `vibesql-logo.svg` - Wordmark with icon
⏳ PNG favicons - Need to generate
⏳ site.webmanifest - Need to create
