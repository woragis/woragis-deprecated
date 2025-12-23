# Code Review Process

**Last Updated:** 2025-12-22  
**Purpose:** Guidelines for conducting effective code reviews

---

## Overview

Code reviews are essential for maintaining code quality, sharing knowledge, and catching bugs before they reach production.

---

## Review Principles

1. **Be Respectful** - Code review is about the code, not the person
2. **Be Constructive** - Provide actionable feedback
3. **Be Timely** - Review PRs within 24-48 hours
4. **Be Thorough** - Check functionality, style, and security
5. **Be Open** - Accept feedback gracefully

---

## Review Checklist

### Code Quality

- [ ] Code follows project coding standards
- [ ] Code is readable and well-organized
- [ ] Functions/classes are appropriately sized
- [ ] Naming is clear and consistent
- [ ] Comments explain "why", not "what"
- [ ] No commented-out code
- [ ] No debug code or console.logs

### Functionality

- [ ] Code works as intended
- [ ] Edge cases are handled
- [ ] Error handling is appropriate
- [ ] No obvious bugs or logic errors
- [ ] Performance considerations addressed

### Testing

- [ ] Tests are included for new features
- [ ] Tests cover edge cases
- [ ] Tests are meaningful (not just for coverage)
- [ ] Integration tests added if needed
- [ ] All tests pass

### Security

- [ ] No hardcoded secrets
- [ ] Input validation present
- [ ] SQL injection prevention (parameterized queries)
- [ ] XSS prevention (output escaping)
- [ ] Authentication/authorization checks
- [ ] No security vulnerabilities

### Documentation

- [ ] Public APIs are documented
- [ ] Complex logic has comments
- [ ] README updated if needed
- [ ] CHANGELOG updated (if applicable)
- [ ] PR description is clear

### Dependencies

- [ ] New dependencies are necessary
- [ ] Dependencies are up to date
- [ ] No known vulnerabilities
- [ ] License is acceptable

---

## Review Process

### 1. Author Responsibilities

**Before Submitting PR:**
- [ ] Self-review your code
- [ ] Run tests locally
- [ ] Run linters/formatters
- [ ] Update documentation
- [ ] Write clear PR description
- [ ] Link related issues

**PR Description Template:**
```markdown
## Description
[Clear description of changes]

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation

## Testing
[How you tested the changes]

## Checklist
- [ ] Code follows standards
- [ ] Tests added/updated
- [ ] Documentation updated
```

### 2. Reviewer Responsibilities

**Review Steps:**
1. **Understand the Context**
   - Read the PR description
   - Check linked issues
   - Understand the problem being solved

2. **Review the Code**
   - Check functionality
   - Review code quality
   - Look for security issues
   - Verify tests

3. **Provide Feedback**
   - Be specific and actionable
   - Explain reasoning
   - Suggest improvements
   - Approve or request changes

4. **Follow Up**
   - Respond to author's questions
   - Re-review after changes
   - Approve when satisfied

---

## Review Comments

### Good Review Comments

**Specific and Actionable:**
```
❌ "This is wrong"
✅ "This function should handle the case when user is nil. Consider adding a nil check at the beginning."
```

**Explains Reasoning:**
```
✅ "Using a map here would be O(1) lookup instead of O(n) for the slice. Since we're checking membership frequently, a map would be more efficient."
```

**Suggests Alternatives:**
```
✅ "Consider using context.WithTimeout here to prevent the request from hanging indefinitely if the external service is slow."
```

### Types of Comments

1. **Must Fix** - Blocking issues that must be addressed
2. **Should Fix** - Important improvements, but not blocking
3. **Nice to Have** - Optional improvements
4. **Questions** - Clarifications needed

---

## Approval Criteria

### Approve When:
- Code is correct and follows standards
- Tests are adequate
- Documentation is updated
- No security issues
- Performance is acceptable

### Request Changes When:
- Critical bugs found
- Security vulnerabilities
- Missing tests for new features
- Code doesn't follow standards
- Documentation is incomplete

### Comment Only When:
- Minor suggestions
- Questions for clarification
- Optional improvements

---

## Review Best Practices

### For Authors

1. **Keep PRs Small**
   - Focus on one feature/fix
   - Easier to review
   - Faster to merge

2. **Respond to Feedback**
   - Address all comments
   - Ask questions if unclear
   - Explain your reasoning

3. **Don't Take It Personally**
   - Feedback is about code, not you
   - Learn from reviews
   - Improve over time

### For Reviewers

1. **Review Promptly**
   - Don't let PRs sit for days
   - Set aside time for reviews
   - Communicate if delayed

2. **Be Constructive**
   - Explain why, not just what
   - Suggest improvements
   - Acknowledge good work

3. **Focus on Important Issues**
   - Don't nitpick style (use linters)
   - Focus on logic and architecture
   - Prioritize security and bugs

---

## Common Issues to Look For

### Security
- Hardcoded secrets
- SQL injection vulnerabilities
- XSS vulnerabilities
- Missing authentication checks
- Insecure random number generation

### Performance
- N+1 query problems
- Missing database indexes
- Inefficient algorithms
- Memory leaks
- Unbounded loops

### Code Quality
- Code duplication
- Overly complex functions
- Poor error handling
- Missing input validation
- Inconsistent patterns

---

## Review Tools

### GitHub/GitLab Features
- **Inline Comments** - Comment on specific lines
- **Suggestion Mode** - Suggest code changes
- **Review Status** - Track review progress
- **Review Requests** - Assign reviewers

### Code Review Tools
- **GitHub Pull Requests** - Built-in review
- **GitLab Merge Requests** - Built-in review
- **Reviewable** - Advanced review features
- **Phabricator** - Enterprise review tool

---

## Review Metrics

**Track:**
- Average review time
- Number of review rounds
- Time to merge
- Review coverage (all PRs reviewed)

**Goals:**
- Reviews within 24 hours
- 1-2 review rounds average
- All PRs reviewed before merge

---

## Escalation

**If Disagreement:**
1. Discuss in PR comments
2. Request additional reviewers
3. Schedule a meeting if needed
4. Document decision in PR

**If Blocked:**
- Communicate clearly
- Provide alternatives
- Escalate to tech lead if needed

---

## Related Documentation

- [Contributing Guide](./contributing.md)
- [Coding Standards](./coding-standards.md)
- [Git Workflow](./git-workflow.md) (when created)
- [PR Template](../.github/PULL_REQUEST_TEMPLATE.md)

---

**Next Steps:**
1. Set up branch protection rules
2. Require PR reviews
3. Train team on review process
4. Monitor review metrics

---

**Last Updated:** 2025-12-22
