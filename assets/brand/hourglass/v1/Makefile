SVG := source/hourglass_favicon.svg
OUT := exports

.PHONY: all clean

all: $(OUT)/favicon-16x16.png $(OUT)/favicon-32x32.png $(OUT)/apple-touch-icon.png $(OUT)/favicon.ico

$(OUT):
	@mkdir -p $(OUT)

$(OUT)/favicon-16x16.png: $(SVG) | $(OUT)
	inkscape $(SVG) --export-width=16 --export-height=16 --export-filename=$@

$(OUT)/favicon-32x32.png: $(SVG) | $(OUT)
	inkscape $(SVG) --export-width=32 --export-height=32 --export-filename=$@

$(OUT)/apple-touch-icon.png: $(SVG) | $(OUT)
	inkscape $(SVG) --export-width=180 --export-height=180 --export-filename=$@

$(OUT)/favicon.ico: $(OUT)/favicon-16x16.png $(OUT)/favicon-32x32.png
	convert $(OUT)/favicon-16x16.png $(OUT)/favicon-32x32.png $@

clean:
	@rm -rf $(OUT)
