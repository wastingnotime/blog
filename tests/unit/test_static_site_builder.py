from __future__ import annotations

from pathlib import Path
from types import SimpleNamespace

from src.app.domain.models.site_config import AnalyticsConfig, SiteConfig
from src.app.infrastructure.builders import static_site_builder
from src.app.infrastructure.builders.static_site_builder import StaticSiteBuilder


def test_static_site_builder_reloads_the_site_module_before_building(
    tmp_path: Path,
    monkeypatch,
) -> None:
    output_dir = tmp_path / "dist"
    assets_dir = tmp_path / "assets"
    assets_dir.mkdir()
    for filename in (
        "favicon.ico",
        "favicon-16x16.png",
        "favicon-32x32.png",
        "apple-touch-icon.png",
        "social-preview.png",
    ):
        (assets_dir / filename).write_bytes(b"asset")

    calls: list[object] = []
    fake_module = SimpleNamespace(
        build_static_site=lambda config, catalog: {"index.html": "<!doctype html>"}
    )

    def fake_import_module(name: str):
        calls.append(("import", name))
        return fake_module

    def fake_reload(module):
        calls.append(("reload", module))
        return module

    monkeypatch.setattr(
        static_site_builder,
        "importlib",
        SimpleNamespace(
            import_module=fake_import_module,
            reload=fake_reload,
        ),
    )

    builder = StaticSiteBuilder(
        output_dir=output_dir,
        identity_assets_dir=assets_dir,
    )
    written_paths = builder.build(
        SiteConfig(
            title="Example",
            description="Static site",
            base_url="https://example.com/",
            analytics=AnalyticsConfig(
                provider="disabled",
                domain="example.com",
                script_url="https://example.com/script.js",
            ),
        ),
        object(),
    )

    assert calls == [
        ("import", "src.app.application.use_cases.build_site"),
        ("reload", fake_module),
    ]
    assert output_dir.joinpath("index.html").read_text(encoding="utf-8") == "<!doctype html>"
    assert len(written_paths) == 6
