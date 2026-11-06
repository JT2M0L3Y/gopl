# Domain Documentation

## Layout: Single-context

### Single-context Layout

The domain documentation follows a single-context layout:

- **Root Domain Docs:** `GLOSSARY.md` at the repository root
- **ADR Directory:** `docs/adr/` at the repository root (one ADR per file)
- **Consumer Rules:** Rules for consuming domain documentation are defined in `docs/agents/domain.md`

### Domain Docs Structure

1. **GLOSSARY.md** - Defines key terms and concepts used throughout the codebase
2. **ADR Directory (docs/adr/)** - Architecture Decision Records describing design decisions
3. **Consumer Rules** - Guidelines for how other components should reference and use domain documentation

### Consumer Rules

- **GLOSSARY.md**: Should be referenced by components when discussing key concepts
- **ADRs**: Should be linked from relevant code files when discussing architectural decisions
- **Cross-referencing**: Components should link to GLOSSARY entries and relevant ADRs when discussing related concepts

### Maintenance

- New terms and concepts should be added to GLOSSARY.md
- Design decisions should be documented in ADRs with appropriate tags
- The domain documentation should be reviewed during major architectural changes