# plugin-morphe-morphemap

Kalo plugin that scaffolds MorpheMap (`.map`) files from Morphe type definitions.

## Overview

Given source and target Morphe types, this plugin generates skeleton `.map` files with:
- Aliases pre-configured for the source and target types
- All target fields listed with best-effort matching to source fields
- `TODO` placeholders for fields that couldn't be auto-matched

The generated scaffolds are intended for manual refinement -- the plugin provides the structure, the developer fills in the mapping logic. This is designed for **cross-domain structural mapping** where human judgment is required (e.g., mapping external API types to local domain models).

## Input

- **Local Morphe Registry** (`KA:MO1:YAML1`): Local project Morphe schema files
- **External Morphe Registry** (`KA:MO1:YAML1`, optional): Third-party API type definitions

## Output

- **MorpheMap files** (`KA:MM1:YAML1`): Skeleton `.map` transformation files

## Configuration

```yaml
config:
  "@kalo-build/plugin-morphe-morphemap":
    mappings:
      - name: "IS24HouseBuyToRealEstateListing"
        sourceAlias: "IS24"
        sourcePath: "immoscout24/RealEstateHouseBuy"
        targetAlias: "Local"
        targetPath: "RealEstateListing"
```

## Build

```bash
cd scripts && bash build.sh
```

## License

MIT
