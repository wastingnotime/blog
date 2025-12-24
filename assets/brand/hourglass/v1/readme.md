# hourglass favicon — v1

This directory contains **version 1** of the Wasting No Time hourglass favicon.

## what is v1?

Version 1 establishes the **baseline visual identity** of the Wasting No Time brand.

It is defined by:
- a **solid black square background**
- a **filled purple hourglass** (`#6e4fda`)
- a **square silhouette** designed to remain legible at very small sizes (16×16, 32×32)

This version prioritizes:
- clarity over decoration
- silhouette over detail
- recognizability in browser tabs and bookmarks

v1 is intentionally conservative.  
Any future changes that alter color, shape, framing, or symbolism must be introduced as a new version (`v2`, `v3`, …).

## source of truth

- `source/hourglass_favicon.svg` is the canonical source.
- All exported assets are generated from this SVG.

## exports

The following files are generated and shipped:

- `favicon.ico`
- `favicon-16x16.png`
- `favicon-32x32.png`
- `apple-touch-icon.png`

## regenerating assets

To regenerate all exported assets:

```bash
make clean all
```

This requires:

* Inkscape
* ImageMagick (`convert`)

## usage

This version is promoted to the site by copying its exports into:

```bash
assets/site/current/
```

The site build process then copies those files into the generated `public/` directory.

## immutability

Once published, **v1 must not be modified**.

Any visual change — no matter how small — must be introduced as a new version.

