# Issue Tracker Configuration

## GitHub Issues

**Repository:** https://github.com/JT2M0L3Y/gopl

**Tracking Work:** Issues and feature requests are managed in the GitHub Issues system for this repository. All engineering work should be tracked as GitHub Issues when possible.

**Workflow:**

1. **Feature Requests & Enhancement Requests:** Created as GitHub Issues with appropriate labels
2. **Bug Reports:** Created as GitHub Issues with bug labels
3. **Engineering Work:** Issues triaged and assigned to appropriate agents
4. **Documentation & Maintenance:** Created as GitHub Issues

**Issue Creation:**
- Use the GitHub Issues interface (web or `gh` CLI)
- Include clear descriptions, reproduction steps (for bugs), and any relevant context
- Apply relevant labels from the triage label vocabulary

**Integration with Engineering Skills:**
- The `to-tickets` skill reads from GitHub Issues to create tracer-bullet tickets
- The `triage` skill processes issues through the triage workflow
- The `to-spec` skill can generate specifications from triaged issues

**PRs as Request Surface:**
GitHub Pull Requests can serve as a request surface for external contributions. However, the default posture is that PRs should not be automatically added to the triage queue. Users who want PRs in the triage queue can enable this by editing this file and setting `pr-as-request-surface: true`.

**Repository Access:**
This repository is publicly accessible on GitHub. All team members can access issues and contributions through the standard GitHub interface.