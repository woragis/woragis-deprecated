# Email Notifications for Test Failures

This document explains how email notifications work for CI/CD test failures.

## Overview

When tests fail in CI/CD, the system automatically sends email notifications using the email-worker service via RabbitMQ.

## How It Works

### Flow

1. **Test Failure Detected**: A workflow (Integration Tests, Performance Tests, etc.) fails
2. **Notification Triggered**: The workflow triggers the `notify-test-failure.yml` workflow
3. **Message Published**: A notification script publishes an email message to RabbitMQ
4. **Email Worker Processes**: The email-worker consumes the message and sends the email
5. **Email Delivered**: Recipients receive the notification email

### Architecture

```
GitHub Actions Workflow (Test Failure)
    ↓
Trigger notify-test-failure.yml
    ↓
send_test_notification.py
    ↓
Publish to RabbitMQ (woragis.notifications exchange)
    ↓
Email Worker consumes message
    ↓
SMTP sends email
    ↓
Recipients receive notification
```

## Configuration

### Required Secrets

Add to GitHub repository secrets:

- `NOTIFICATION_EMAIL` - Email address(es) to receive notifications
  - Can be a single email: `dev-team@woragis.com`
  - Or multiple emails (comma-separated): `dev1@woragis.com,dev2@woragis.com`

### Email Worker Configuration

The email-worker must be configured with:
- `SMTP_HOST` - SMTP server hostname
- `SMTP_PORT` - SMTP server port (default: 587)
- `SMTP_USERNAME` - SMTP username
- `SMTP_PASSWORD` - SMTP password
- `SMTP_FROM` - From email address

### RabbitMQ Configuration

The notification script uses:
- Exchange: `woragis.notifications` (configurable via `EMAIL_EXCHANGE`)
- Routing Key: `emails.send` (configurable via `EMAIL_ROUTING_KEY`)

## Email Format

### Subject
```
⚠️ Test Failure: [Workflow Name]
```

### Content
- Workflow name and status
- Triggered by (actor)
- Commit SHA
- Workflow URL (link to GitHub Actions)
- Time of failure
- Action required steps

### Example Email

```
Subject: ⚠️ Test Failure: Integration Tests

Workflow: Integration Tests
Status: FAILURE
Triggered by: john-doe
Time: 2024-01-15 10:30:45 UTC

Commit: a1b2c3d4
Workflow URL: https://github.com/org/repo/actions/runs/123456789

Action Required:
- Review the workflow logs
- Check for test failures or errors
- Fix issues and re-run tests
```

## Workflow Integration

### Integration Tests Workflow

Notifications sent when:
- Any test suite fails (Server, Email Worker, Translation Worker, WhatsApp Worker)

### Performance Tests Workflow

Notifications sent when:
- Performance tests fail
- Performance regression detected

### Performance Regression Workflow

Notifications sent when:
- Performance regressions exceed thresholds

## Manual Testing

### Test Notification Script Locally

1. Start RabbitMQ:
   ```bash
   docker run -d --name test-rabbitmq \
     -p 5673:5672 \
     -e RABBITMQ_DEFAULT_USER=test \
     -e RABBITMQ_DEFAULT_PASS=test \
     -e RABBITMQ_DEFAULT_VHOST=test \
     rabbitmq:3.13-management-alpine
   ```

2. Start email-worker (in another terminal):
   ```bash
   cd backend/email-worker
   export SMTP_HOST=smtp.example.com
   export SMTP_FROM=noreply@woragis.com
   # ... other SMTP config
   go run cmd/email-worker/main.go
   ```

3. Run notification script:
   ```bash
   python3 .github/scripts/send_test_notification.py \
     --workflow "Test Workflow" \
     --status "failure" \
     --url "https://github.com/org/repo/actions/runs/123" \
     --commit "abc123" \
     --actor "test-user" \
     --email "test@example.com"
   ```

### Test via GitHub Actions

1. Go to Actions → Notify Test Failure
2. Click "Run workflow"
3. Fill in the inputs:
   - Workflow name: "Test Notification"
   - Workflow URL: Any GitHub Actions URL
   - Status: "failure"
   - Commit SHA: Any commit SHA
   - Actor: Your username
4. Run workflow
5. Check email inbox

## Troubleshooting

### No Email Received

**Possible Causes:**
- Email-worker not running
- SMTP configuration incorrect
- RabbitMQ connection issues
- Email in spam folder

**Solutions:**
1. Check email-worker logs
2. Verify SMTP configuration
3. Check RabbitMQ connection
4. Check spam folder

### Notification Not Triggered

**Possible Causes:**
- Workflow not configured to trigger notifications
- Notification workflow disabled
- Secrets not configured

**Solutions:**
1. Check workflow files for notification triggers
2. Verify `NOTIFICATION_EMAIL` secret is set
3. Check workflow permissions

### RabbitMQ Connection Failed

**Possible Causes:**
- RabbitMQ service not available
- Wrong connection URL
- Authentication failed

**Solutions:**
1. Check RabbitMQ service health
2. Verify connection URL
3. Check credentials

## Advanced Configuration

### Multiple Recipients

Set `NOTIFICATION_EMAIL` to comma-separated emails:
```
dev1@woragis.com,dev2@woragis.com,manager@woragis.com
```

**Note**: The current implementation sends to a single email. For multiple recipients, you would need to:
1. Split the comma-separated list
2. Send separate messages for each recipient
3. Or modify email-worker to handle multiple recipients

### Custom Email Template

Edit `.github/scripts/send_test_notification.py` to customize:
- Email subject format
- Email body template
- HTML styling

### Different Notifications for Different Workflows

Modify the notification script to send different emails based on workflow name:

```python
if workflow_name == "Performance Tests":
    subject = "🚨 Performance Test Failure"
elif workflow_name == "Integration Tests":
    subject = "⚠️ Integration Test Failure"
```

## Best Practices

1. **Set Appropriate Recipients**: Only notify relevant team members
2. **Monitor Email Worker**: Ensure email-worker is running and healthy
3. **Check Spam Folder**: Add sender to contacts to avoid spam
4. **Review Notifications**: Don't ignore test failure notifications
5. **Fix Promptly**: Address test failures quickly

## Integration with Email Worker

The notification system integrates seamlessly with the email-worker:

1. **Uses Same Queue**: Messages go to the same RabbitMQ queue as other emails
2. **Same Format**: Uses the same `EmailEnvelope` format
3. **Same Processing**: Email-worker processes notifications like any other email
4. **Reliable Delivery**: Benefits from email-worker's retry logic and DLQ

## Production Setup

### For Production CI/CD

1. **Configure Production RabbitMQ**:
   - Update `RABBITMQ_URL` in workflow
   - Use production RabbitMQ instance
   - Ensure email-worker is running

2. **Set Notification Email**:
   - Add `NOTIFICATION_EMAIL` secret in GitHub
   - Use team distribution list or monitoring email

3. **Ensure Email Worker Running**:
   - Deploy email-worker to production
   - Configure SMTP settings
   - Monitor email-worker health

## Next Steps

1. Configure `NOTIFICATION_EMAIL` secret in GitHub
2. Ensure email-worker is running and configured
3. Test notification flow
4. Monitor email delivery
5. Customize email templates if needed
