---
name: 🐛 Bug Report
about: Create a report to help us improve
title: '[BUG] '
labels: bug
assignees: ''
---

## Bug Description
A clear and concise description of what the bug is.

## To Reproduce
Steps to reproduce the behavior:
1. Start mockapi with config '...'
2. Send request to '....'
3. See error

## Expected Behavior
A clear and concise description of what you expected to happen.

## Actual Behavior
What actually happened.

## Environment
- **OS**: [e.g. macOS 14.0, Windows 11, Ubuntu 22.04]
- **MockAPI Version**: [e.g. v1.0.0]
- **Go Version**: [e.g. 1.21.0]
- **Installation Method**: [go install / Docker / Binary]

## Configuration
```json
{
  "port": 8088,
  "routes": [...]
}
```

## Request/Response
**Request**:
```bash
curl -X GET http://localhost:8088/mock/users/1
```

**Response**:
```json
{
  "error": "..."
}
```

## Logs
```
[paste logs here]
```

## Screenshots
If applicable, add screenshots to help explain your problem.

## Additional Context
Add any other context about the problem here.
