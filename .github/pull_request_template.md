## Summary

Brief description of what this PR does and why.

## Changes

- [ ] Added new data source
- [ ] Modified existing data source
- [ ] Updated provider configuration
- [ ] Fixed a bug
- [ ] Updated documentation
- [ ] Added tests
- [ ] Refactored code
- [ ] Other: ___________

## Type of Change

- [ ] New feature (non-breaking change that adds functionality)
- [ ] Bug fix (non-breaking change that fixes an issue)
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update
- [ ] Test improvement
- [ ] Performance improvement

## Terraform Provider Specific

- [ ] New resource or data source added
- [ ] Schema changes (breaking/non-breaking)
- [ ] Provider configuration changes
- [ ] Import functionality changes
- [ ] State migration required
- [ ] Terraform version compatibility maintained

## Testing

- [ ] All existing tests pass
- [ ] New tests added for new functionality
- [ ] Manual testing completed with Terraform
- [ ] Tested with multiple Terraform versions (if applicable)
- [ ] Acceptance tests pass (if applicable)
- [ ] No regression in existing functionality

## Documentation

- [ ] README.md updated if needed
- [ ] Examples updated if needed
- [ ] Schema documentation updated
- [ ] CHANGELOG.md updated (if applicable)

## Checklist

- [ ] Code follows project conventions
- [ ] All whimsy names follow constraints (max 6 chars, a-z only, alphabetical)
- [ ] No hardcoded values that should be configurable
- [ ] Error messages are helpful and clear
- [ ] Code is backward compatible (or breaking changes are documented)
- [ ] Tests cover edge cases and error conditions

## Terraform Configuration Example

If this PR adds/changes functionality, provide a minimal Terraform configuration example:

```hcl
# Example usage of the changes in this PR
terraform {
  required_providers {
    whimsy = {
      source = "mioi/whimsy"
    }
  }
}

# Your example here
```

## Additional Notes

Any additional context, screenshots, or information that would be helpful for reviewers.

## Breaking Changes

If this PR contains breaking changes, describe:
- What will break
- Migration path for users
- Deprecation timeline (if applicable)