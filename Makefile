VER ?= v1

HG_DIR   := assets/brand/hourglass/$(VER)
HG_MAKE  := $(HG_DIR)/Makefile
EXPORTS  := $(HG_DIR)/exports
CURRENT  := assets/site/current

REQUIRED := favicon.ico favicon-16x16.png favicon-32x32.png apple-touch-icon.png

.PHONY: favicon-build favicon-clean favicon-promote favicon-all

favicon-build:
	@test -f "$(HG_MAKE)" || (echo "Missing Makefile: $(HG_MAKE)" && exit 1)
	@$(MAKE) -C "$(HG_DIR)" all

favicon-clean:
	@test -f "$(HG_MAKE)" || (echo "Missing Makefile: $(HG_MAKE)" && exit 1)
	@$(MAKE) -C "$(HG_DIR)" clean

favicon-promote: favicon-build
	@test -d "$(EXPORTS)" || (echo "Missing exports: $(EXPORTS)" && exit 1)
	@for f in $(REQUIRED); do \
		test -f "$(EXPORTS)/$$f" || (echo "Missing required export: $(EXPORTS)/$$f" && exit 1); \
	done
	@mkdir -p "$(CURRENT)"
	@for f in $(REQUIRED); do \
		cp -f "$(EXPORTS)/$$f" "$(CURRENT)/"; \
	done
	@echo "✔ Promoted hourglass $(VER) -> $(CURRENT)"

favicon-all: favicon-promote

.PHONY: favicon-promote-commit
favicon-promote-commit: favicon-promote
	@git add "$(CURRENT)"
	@git commit -m "chore: assets - set current favicon to hourglass $(VER)" || true
