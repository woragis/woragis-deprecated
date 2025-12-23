# Contributing Guide

**Last Updated:** 2025-12-22  
**Purpose:** Guidelines for contributing to the Woragis backend project

---

## Welcome!

Thank you for your interest in contributing to Woragis! This guide will help you get started.

---

## Code of Conduct

- Be respectful and inclusive
- Welcome newcomers and help them learn
- Focus on constructive feedback
- Respect different viewpoints and experiences

---

## How to Contribute

### Reporting Bugs

1. **Check existing issues** - Search GitHub issues to see if the bug is already reported
2. **Create a new issue** if it doesn't exist
3. **Include:**
   - Clear description of the bug
   - Steps to reproduce
   - Expected vs actual behavior
   - Environment details (OS, Docker version, etc.)
   - Relevant logs or error messages

### Suggesting Features

1. **Check existing issues** - See if the feature is already requested
2. **Create a feature request issue**
3. **Include:**
   - Clear description of the feature
   - Use case and benefits
   - Possible implementation approach (if you have ideas)

### Contributing Code

1. **Fork the repository**
2. **Create a feature branch** (see Git Workflow below)
3. **Make your changes**
4. **Write/update tests**
5. **Update documentation**
6. **Submit a pull request**

---

## Git Workflow

### Branch Naming

Use descriptive branch names:
- `feature/add-user-authentication`
- `fix/rabbitmq-connection-error`
- `docs/update-api-documentation`
- `refactor/improve-error-handling`

### Commit Messages

Follow conventional commits format:

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**
```
feat(ai-service): add support for Anthropic Claude

Add Anthropic API integration to ai-service with fallback
to OpenAI if Anthropic is unavailable.

Closes #123
```

```
fix(resume-worker): resolve syntax error in main.py

Fixed unclosed try block that was causing SyntaxError
on worker startup.

Fixes #456
```

### Pull Request Process

1. **Create a branch** from `main` (or `develop` if using Git Flow)
2. **Make your changes** following coding standards
3. **Write tests** for new features or bug fixes
4. **Update documentation** if needed
5. **Ensure all tests pass:**
   ```bash
   make test-all
   # or
   docker-compose up -d
   go test ./...
   pytest
   ```
6. **Create pull request:**
   - Clear title and description
   - Reference related issues
   - Request review from maintainers
   - Respond to feedback promptly

---

## Development Setup

See [Development Setup Guide](./setup-guide.md) for detailed instructions.

**Quick Start:**
```bash
# Clone repository
git clone <repository-url>
cd woragis/backend

# Create .env file
cp .env.example .env
# Edit .env with your configuration

# Start services
docker-compose up -d

# Verify setup
curl http://localhost:8080/healthz
```

---

## Coding Standards

### Go Code

- Follow [Effective Go](https://go.dev/doc/effective_go) guidelines
- Use `gofmt` for formatting
- Run `golint` or `golangci-lint` for linting
- Write unit tests for new functions
- Document exported functions and types

**Example:**
```go
// ProcessUser processes a user request and returns the result.
// It validates the input and handles errors appropriately.
func ProcessUser(ctx context.Context, userID string) (*User, error) {
    // Implementation
}
```

### Python Code

- Follow [PEP 8](https://pep8.org/) style guide
- Use `black` for formatting
- Use `pylint` or `flake8` for linting
- Write docstrings for functions and classes
- Type hints where appropriate

**Example:**
```python
def process_user(user_id: str) -> User:
    """
    Process a user request and return the user object.
    
    Args:
        user_id: The unique identifier for the user
        
    Returns:
        User object with user data
        
    Raises:
        ValueError: If user_id is invalid
        NotFoundError: If user is not found
    """
    # Implementation
```

### JavaScript/TypeScript Code

- Follow [Airbnb JavaScript Style Guide](https://github.com/airbnb/javascript)
- Use `prettier` for formatting
- Use `eslint` for linting
- Write JSDoc comments for functions
- Use TypeScript where possible

---

## Testing

### Writing Tests

- Write tests for all new features
- Aim for >80% code coverage
- Test both success and error cases
- Use descriptive test names

**Go Example:**
```go
func TestProcessUser(t *testing.T) {
    tests := []struct {
        name    string
        userID  string
        want    *User
        wantErr bool
    }{
        {
            name:    "valid user",
            userID:  "123",
            want:    &User{ID: "123"},
            wantErr: false,
        },
        {
            name:    "invalid user",
            userID:  "",
            want:    nil,
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ProcessUser(tt.userID)
            if (err != nil) != tt.wantErr {
                t.Errorf("ProcessUser() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            // Assertions
        })
    }
}
```

**Python Example:**
```python
def test_process_user():
    """Test processing a valid user."""
    user = process_user("123")
    assert user.id == "123"
    assert user.name is not None

def test_process_user_invalid():
    """Test processing an invalid user."""
    with pytest.raises(ValueError):
        process_user("")
```

### Running Tests

```bash
# All tests
make test-all

# Go tests
cd server && go test ./...

# Python tests
cd ai-service && pytest

# Integration tests
docker-compose up -d
go test -tags=integration ./...
```

---

## Documentation

### When to Update Documentation

- Adding new features
- Changing API endpoints
- Modifying configuration
- Adding new services or workers
- Changing deployment procedures

### Documentation Standards

- Use clear, concise language
- Include code examples
- Update related documentation
- Add diagrams where helpful
- Keep documentation up to date

### Documentation Locations

- **API Docs:** `docs/api/`
- **Architecture:** `docs/architecture/`
- **Development:** `docs/development/`
- **Deployment:** `docs/deployment/`
- **Runbooks:** `docs/runbooks/`
- **Service READMEs:** `{service}/README.md`

---

## Code Review Process

### For Contributors

1. **Address feedback** promptly and professionally
2. **Make requested changes** or explain why not
3. **Keep PRs focused** - one feature or fix per PR
4. **Respond to comments** and questions
5. **Update PR** as needed based on feedback

### For Reviewers

1. **Be constructive** and respectful
2. **Explain reasoning** for requested changes
3. **Approve promptly** when satisfied
4. **Test changes** if possible
5. **Thank contributors** for their work

---

## Adding New Features

### Before Starting

1. **Discuss in issue** or team chat first
2. **Check existing code** for similar patterns
3. **Plan the implementation**
4. **Consider backward compatibility**

### Implementation Steps

1. **Create feature branch**
2. **Write tests first** (TDD approach)
3. **Implement feature**
4. **Update documentation**
5. **Run all tests**
6. **Create pull request**

### Adding a New Service

See [Adding a Service](./adding-service.md) guide.

### Adding a New Worker

See [Adding a Worker](./adding-worker.md) guide.

---

## Security

### Security Best Practices

- **Never commit secrets** (API keys, passwords, etc.)
- **Use environment variables** for sensitive data
- **Validate all inputs**
- **Use parameterized queries** (prevent SQL injection)
- **Sanitize user input**
- **Keep dependencies updated**
- **Report security issues** privately

### Reporting Security Issues

**DO NOT** create a public issue for security vulnerabilities.

Instead:
1. Email security team or maintainers privately
2. Include details about the vulnerability
3. Wait for response before disclosing publicly

---

## Getting Help

### Resources

- **Documentation:** `docs/` directory
- **Issues:** GitHub Issues
- **Discussions:** GitHub Discussions (if enabled)
- **Team Chat:** (if available)

### Asking Questions

When asking for help:
1. **Search existing issues** first
2. **Provide context** (what you're trying to do)
3. **Include error messages** and logs
4. **Show what you've tried**
5. **Be patient** - maintainers are volunteers

---

## Recognition

Contributors will be:
- Listed in CONTRIBUTORS.md (if exists)
- Credited in release notes
- Appreciated by the team! 🎉

---

## Questions?

If you have questions about contributing:
1. Check the documentation
2. Search existing issues
3. Create a new issue with the `question` label
4. Reach out to maintainers

---

## Thank You!

Your contributions make Woragis better for everyone. Thank you for taking the time to contribute! 🙏

---

**Last Updated:** 2025-12-22
