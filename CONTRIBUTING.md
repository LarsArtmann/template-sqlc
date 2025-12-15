# Contributing to template-sqlc

Thank you for your interest in contributing to the template-sqlc project! This template is built from real-world usage across 21+ projects, and we welcome contributions that make it even better.

## 🎯 How You Can Contribute

### 🐛 Bug Reports

- **Configuration errors**: Invalid sqlc settings, syntax errors
- **Documentation issues**: Unclear explanations, broken links
- **Template bugs**: Configurations that don't work as expected

### ✨ Feature Requests

- **New project patterns**: Additional use cases or project types
- **Better defaults**: Improved configuration recommendations
- **Missing sqlc features**: New sqlc options that aren't covered

### 📚 Documentation Improvements

- **Better explanations**: Clearer reasoning for configuration choices
- **More examples**: Additional project type configurations
- **Troubleshooting**: Solutions for common problems

### 🔧 Code Contributions

- **Configuration updates**: Keep up with latest sqlc features
- **New patterns**: Real-world configurations from your projects
- **Performance improvements**: Better defaults for performance

## 🚀 Getting Started

### Prerequisites

- [sqlc v1.29.0+](https://docs.sqlc.dev/en/stable/overview/install.html)
- Basic understanding of SQL and Go
- Familiarity with your target database (SQLite, PostgreSQL, MySQL)

### Development Setup

```bash
# 1. Fork the repository on GitHub
# 2. Clone your fork
git clone https://github.com/YOUR_USERNAME/template-sqlc.git
cd template-sqlc

# 3. Test the configuration
sqlc compile

# 4. Make your changes
# Edit sqlc.yaml, README.md, or other files

# 5. Validate your changes
sqlc compile
sqlc vet  # If you have a database configured
```

## 📋 Contribution Guidelines

### Configuration Changes

1. **Test thoroughly**: Ensure `sqlc compile` passes
2. **Add comments**: Explain WHY each setting is chosen
3. **Consider trade-offs**: Document performance/flexibility implications
4. **Keep universal**: Avoid domain-specific hardcoded values

### Documentation Updates

1. **Use clear language**: Explain concepts for beginners
2. **Provide examples**: Show practical usage patterns
3. **Update table of contents**: If adding new sections
4. **Test links**: Ensure all links work correctly

### Adding New Project Patterns

```yaml
# === YOUR_PROJECT_TYPE ===
# Description of when to use this configuration
# sql:
#   - name: "your_config"
#     engine: "postgresql"  # or sqlite, mysql
#     # ... rest of configuration
#     # Include comments explaining choices
```

## 📝 Pull Request Process

### 1. Create a Feature Branch

```bash
git checkout -b feature/amazing-feature
# OR
git checkout -b fix/configuration-bug
# OR
git checkout -b docs/better-examples
```

### 2. Make Your Changes

- Follow existing code style and commenting patterns
- Test your changes with `sqlc compile`
- Update documentation as needed

### 3. Commit Your Changes

Use clear, descriptive commit messages:

```bash
git commit -m "feat: add microservices configuration pattern

- Add service-specific database configurations
- Include examples for user and order services
- Document when to use this pattern vs monolithic setup"
```

### 4. Push and Create PR

```bash
git push origin feature/amazing-feature
```

Then create a Pull Request on GitHub with:

- **Clear title**: Summarize what you changed
- **Detailed description**: Explain why the change is needed
- **Test results**: Show `sqlc compile` output if relevant
- **Breaking changes**: Note any incompatibilities

## ✅ Quality Standards

### Configuration Requirements

- [ ] `sqlc compile` passes without errors
- [ ] All settings have explanatory comments
- [ ] Trade-offs and alternatives documented
- [ ] Universal compatibility (avoid project-specific hardcoding)

### Documentation Requirements

- [ ] Clear, beginner-friendly language
- [ ] Practical examples included
- [ ] Links tested and working
- [ ] Table of contents updated (if applicable)

### Code Style

- Use consistent YAML formatting (2 spaces)
- Add comprehensive comments explaining choices
- Group related settings with clear section headers
- Use emoji/symbols for visual organization (✅❌🚫🛡️🚀)

## 🤔 Questions or Need Help?

### Getting Help

- **GitHub Issues**: Ask questions about configuration
- **Discussions**: Brainstorm new features or patterns
- **SQLc Documentation**: https://docs.sqlc.dev/

### Common Questions

**Q: Should I include generated code in the template?**
A: No, only include the `sqlc.yaml` configuration. Generated code depends on actual SQL schema.

**Q: How do I test my configuration changes?**
A: Run `sqlc compile` to validate syntax. Use `sqlc generate` with actual SQL files to test code generation.

**Q: Can I add configurations for my specific domain (e.g., e-commerce)?**  
A: Yes! Add them as commented examples in the "Alternative Configurations" section. Keep them generic enough for others to adapt.

**Q: Should I update the README?**
A: Yes, if you add significant new features or patterns, update the relevant sections in README.md.

## 🏆 Recognition

Contributors will be recognized in:

- GitHub contributor list
- README.md acknowledgments section
- Release notes for significant contributions

## 📜 Code of Conduct

This project follows the standard open source code of conduct:

- Be respectful and inclusive
- Focus on constructive feedback
- Help newcomers learn and contribute
- Prioritize the community and project goals

## 🎉 Thank You!

Every contribution makes this template better for the entire sqlc community. Whether it's fixing a typo, adding a new pattern, or improving documentation - thank you for helping make database development easier for everyone!

---

_This project is inspired by real-world usage across 21+ different projects. Your contributions help ensure it continues to serve diverse use cases effectively._
