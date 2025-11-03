# Contributing to Thaiyyal

First off, thank you for considering contributing to Thaiyyal! It's people like you that make Thaiyyal such a great tool.

## Code of Conduct

This project and everyone participating in it is governed by our Code of Conduct. By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

### Our Pledge

We pledge to make participation in our project a harassment-free experience for everyone, regardless of:
- Age, body size, disability, ethnicity, gender identity and expression
- Level of experience, education, socio-economic status
- Nationality, personal appearance, race, religion
- Sexual identity and orientation

### Our Standards

**Positive behavior includes:**
- Using welcoming and inclusive language
- Being respectful of differing viewpoints and experiences
- Gracefully accepting constructive criticism
- Focusing on what is best for the community
- Showing empathy towards other community members

**Unacceptable behavior includes:**
- Trolling, insulting/derogatory comments, and personal or political attacks
- Public or private harassment
- Publishing others' private information without explicit permission
- Other conduct which could reasonably be considered inappropriate

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. When creating a bug report, include as many details as possible:

**Bug Report Template:**

```markdown
**Describe the bug**
A clear and concise description of what the bug is.

**To Reproduce**
Steps to reproduce the behavior:
1. Go to '...'
2. Click on '....'
3. Scroll down to '....'
4. See error

**Expected behavior**
A clear description of what you expected to happen.

**Workflow JSON**
If applicable, include the workflow JSON that causes the issue.

**Screenshots**
If applicable, add screenshots to help explain your problem.

**Environment:**
 - OS: [e.g., macOS, Linux, Windows]
 - Go Version: [e.g., 1.24.7]
 - Node.js Version: [e.g., 20.x]
 - Browser: [e.g., Chrome 120, Safari 17]

**Additional context**
Add any other context about the problem here.
```

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, include:

**Enhancement Template:**

```markdown
**Is your feature request related to a problem?**
A clear description of what the problem is. Ex. I'm always frustrated when [...]

**Describe the solution you'd like**
A clear and concise description of what you want to happen.

**Describe alternatives you've considered**
A clear description of any alternative solutions or features you've considered.

**Additional context**
Add any other context or screenshots about the feature request here.

**Would you like to implement this feature?**
Let us know if you're interested in implementing this yourself.
```

### Your First Code Contribution

Unsure where to begin? Look for issues labeled:
- `good first issue` - Good for newcomers
- `help wanted` - Extra attention needed
- `documentation` - Improvements or additions to documentation

## Development Workflow

### 1. Set Up Development Environment

```bash
# Fork the repository on GitHub
# Clone your fork
git clone https://github.com/YOUR_USERNAME/thaiyyal.git
cd thaiyyal

# Add upstream remote
git remote add upstream https://github.com/yesoreyeram/thaiyyal.git

# Install dependencies
npm install
cd backend && go mod download
```

### 2. Create a Branch

```bash
# Create a feature branch
git checkout -b feature/my-new-feature

# Or a bugfix branch
git checkout -b fix/issue-123
```

### 3. Make Your Changes

#### Backend (Go) Changes

```bash
cd backend

# Run tests frequently
go test ./...

# Run tests with coverage
go test -cover ./...

# Format code
go fmt ./...

# Run linter (if gol angci-lint is installed)
golangci-lint run

# Build to ensure no compilation errors
go build ./...
```

#### Frontend (TypeScript/React) Changes

```bash
# Run development server
npm run dev

# Run linter
npm run lint

# Build to ensure no errors
npm run build
```

### 4. Write Tests

All code changes should include appropriate tests:

**Backend Testing:**
```go
// Example test structure
func TestMyFeature(t *testing.T) {
    // Arrange
    payload := `{"nodes": [...], "edges": [...]}`
    
    // Act
    engine, err := workflow.NewEngine([]byte(payload))
    require.NoError(t, err)
    result, err := engine.Execute()
    
    // Assert
    require.NoError(t, err)
    assert.Equal(t, expectedValue, result.FinalOutput)
}
```

**Testing Guidelines:**
- Write unit tests for all new functions
- Write integration tests for new features
- Aim for 80%+ code coverage
- Use table-driven tests where appropriate
- Test edge cases and error conditions

### 5. Commit Your Changes

We follow [Conventional Commits](https://www.conventionalcommits.org/) specification:

```bash
# Format: <type>(<scope>): <description>
git commit -m "feat(executor): add custom validation middleware"
git commit -m "fix(engine): resolve topological sort issue"
git commit -m "docs(readme): update installation instructions"
git commit -m "test(executor): add tests for new node types"
```

**Commit Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation only changes
- `style`: Code style changes (formatting, missing semicolons, etc.)
- `refactor`: Code change that neither fixes a bug nor adds a feature
- `perf`: Performance improvement
- `test`: Adding or updating tests
- `chore`: Changes to build process or auxiliary tools

**Commit Message Guidelines:**
- Use the present tense ("add feature" not "added feature")
- Use the imperative mood ("move cursor to..." not "moves cursor to...")
- Limit the first line to 72 characters or less
- Reference issues and pull requests liberally after the first line

### 6. Push and Create Pull Request

```bash
# Push to your fork
git push origin feature/my-new-feature

# Create Pull Request on GitHub
```

**Pull Request Template:**

```markdown
## Description
Brief description of changes made.

## Related Issue
Fixes #123

## Type of Change
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update

## Testing
Describe the tests you ran to verify your changes:
- [ ] Unit tests pass (`go test ./...`)
- [ ] Integration tests pass
- [ ] Manual testing performed
- [ ] Coverage maintained or improved

## Checklist
- [ ] My code follows the project's style guidelines
- [ ] I have performed a self-review of my own code
- [ ] I have commented my code, particularly in hard-to-understand areas
- [ ] I have made corresponding changes to the documentation
- [ ] My changes generate no new warnings
- [ ] I have added tests that prove my fix is effective or that my feature works
- [ ] New and existing unit tests pass locally with my changes
- [ ] Any dependent changes have been merged and published

## Screenshots (if applicable)
Add screenshots to help reviewers understand your changes.
```

## Coding Standards

### Go Code Style

1. **Follow Go Conventions:**
   - Use `gofmt` for formatting
   - Follow [Effective Go](https://go.dev/doc/effective_go)
   - Use [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) guidelines

2. **Package Organization:**
   ```go
   // Good: Clear package purpose
   package executor // Defines node executors
   
   // Bad: Vague package purpose
   package utils
   ```

3. **Error Handling:**
   ```go
   // Good: Wrap errors with context
   if err != nil {
       return fmt.Errorf("failed to parse workflow: %w", err)
   }
   
   // Bad: Lose error context
   if err != nil {
       return err
   }
   ```

4. **Comments:**
   ```go
   // Good: Describe what and why
   // ValidateURL validates a URL for SSRF protection by checking
   // the scheme, domain, and IP address against configured rules.
   func (p *SSRFProtection) ValidateURL(urlStr string) error {
   
   // Bad: State the obvious
   // ValidateURL validates URL
   func (p *SSRFProtection) ValidateURL(urlStr string) error {
   ```

5. **Testing:**
   ```go
   // Use table-driven tests
   func TestOperation(t *testing.T) {
       tests := []struct {
           name     string
           op       string
           inputs   []float64
           expected float64
           wantErr  bool
       }{
           {"add", "add", []float64{1, 2}, 3, false},
           {"divide by zero", "divide", []float64{1, 0}, 0, true},
       }
       
       for _, tt := range tests {
           t.Run(tt.name, func(t *testing.T) {
               // Test logic
           })
       }
   }
   ```

### TypeScript/React Code Style

1. **Use TypeScript:**
   - Define interfaces for all props
   - Avoid `any` type
   - Use strict type checking

2. **Component Structure:**
   ```typescript
   // Good: Typed props and clear structure
   interface NodePaletteProps {
       onNodeSelect: (nodeType: string) => void;
       availableNodes: NodeType[];
   }
   
   export const NodePalette: React.FC<NodePaletteProps> = ({ 
       onNodeSelect, 
       availableNodes 
   }) => {
       // Component logic
   };
   ```

3. **Follow React Best Practices:**
   - Use functional components with hooks
   - Properly handle side effects with `useEffect`
   - Memoize expensive computations with `useMemo`
   - Use `useCallback` for stable function references

### Documentation Standards

1. **Code Documentation:**
   - Document all public APIs
   - Include examples in documentation
   - Explain complex algorithms
   - Document security considerations

2. **README Files:**
   - Clear, concise descriptions
   - Step-by-step instructions
   - Working examples
   - Troubleshooting section

3. **Architecture Documentation:**
   - Include diagrams (use Mermaid)
   - Explain design decisions
   - Document trade-offs
   - Keep updated with code changes

## Security Guidelines

### Security Best Practices

1. **Input Validation:**
   - Validate all user inputs
   - Sanitize data before processing
   - Use type-safe parsing
   - Set resource limits

2. **Dependency Management:**
   - Keep dependencies up to date
   - Review security advisories
   - Use `go mod tidy` regularly
   - Run `npm audit` regularly

3. **Secret Management:**
   - Never commit secrets
   - Use environment variables
   - Use `.gitignore` for sensitive files
   - Review commits before pushing

4. **SSRF Protection:**
   - Validate URLs before making requests
   - Block private IP ranges
   - Block cloud metadata endpoints
   - Use allowlist when possible

### Reporting Security Vulnerabilities

**DO NOT** create public GitHub issues for security vulnerabilities.

Instead:
1. Email security details to the maintainers
2. Include reproduction steps
3. Describe the impact
4. Suggest a fix if possible

We will respond within 48 hours and work with you to address the issue.

## Review Process

### What to Expect

1. **Automated Checks:**
   - CI/CD pipeline runs tests
   - Linters check code style
   - Security scanners check for vulnerabilities
   - Coverage reports generated

2. **Code Review:**
   - At least one maintainer reviews
   - Feedback provided within 7 days
   - Changes may be requested
   - Approval required before merge

3. **Merge:**
   - Squash and merge is preferred
   - Commit message should be clear
   - Branch deleted after merge

### Review Checklist

Reviewers will check:
- [ ] Code follows style guidelines
- [ ] Tests are comprehensive
- [ ] Documentation is updated
- [ ] No security vulnerabilities introduced
- [ ] Performance impact is acceptable
- [ ] Breaking changes are documented
- [ ] Backward compatibility maintained (if applicable)

## Community

### Getting Help

- **Documentation**: Start with our comprehensive docs
- **Discussions**: Use GitHub Discussions for questions
- **Issues**: Search existing issues before creating new ones
- **Chat**: Join our community chat (link TBD)

### Recognition

Contributors are recognized in:
- Release notes
- Contributors list
- Project README

Significant contributions may result in:
- Committer status
- Maintainer status
- Public acknowledgment

## Attribution

This Contributing Guide is adapted from:
- [Contributor Covenant](https://www.contributor-covenant.org/)
- [Open Source Guides](https://opensource.guide/)
- Various open source projects' contribution guidelines

Thank you for contributing to Thaiyyal! 🎉
