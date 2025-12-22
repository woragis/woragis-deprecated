# Logging Security and Compliance

**Version:** 1.0  
**Last Updated:** 2025-12-22  
**Status:** Active

## Overview

This document outlines security best practices and compliance requirements for the Woragis logging aggregation system.

## Security Principles

### 1. No Sensitive Data in Logs

**Never log:**
- Passwords or password hashes
- API keys or tokens
- Credit card numbers
- Social security numbers
- Personal identification numbers (PINs)
- Authentication tokens (except trace IDs)
- Database connection strings with passwords
- Encryption keys

**Safe to log:**
- User IDs (non-sensitive identifiers)
- Request IDs and trace IDs
- Service names
- Error messages (without sensitive data)
- Performance metrics
- HTTP status codes
- Request paths (without query parameters containing sensitive data)

### 2. Access Control

#### Grafana Access
- **Default credentials**: Change immediately in production
- **User accounts**: Create individual accounts for team members
- **Role-based access**: Use Grafana roles (Viewer, Editor, Admin)
- **SSO Integration**: Configure SSO for production (recommended)

#### Loki Access
- **Network isolation**: Loki should not be exposed publicly
- **Internal only**: Access only from within Docker network
- **Authentication**: Consider adding authentication layer (future)

### 3. Encryption

#### Encryption in Transit
- **Grafana**: Use HTTPS in production
- **Loki**: Internal network (Docker network is isolated)
- **Promtail**: Internal network only

#### Encryption at Rest
- **Loki data**: Currently not encrypted (future enhancement)
- **Grafana data**: Currently not encrypted (future enhancement)
- **Recommendation**: Use encrypted volumes in production

### 4. Network Security

#### Current Configuration
- **Loki**: Port 3100 (internal only)
- **Grafana**: Port 3000 (exposed, should be behind reverse proxy in production)
- **Promtail**: Port 9080 (internal only)

#### Production Recommendations
- Use reverse proxy (nginx/traefik) for Grafana
- Enable HTTPS/TLS
- Restrict access by IP (if needed)
- Use VPN or private network for access

## Compliance Requirements

### GDPR Compliance

#### Right to Erasure
- Logs may contain user data (user_id, IP addresses)
- Must be able to delete logs for specific users
- Implement log redaction/deletion procedures

#### Data Minimization
- Only log necessary information
- Avoid logging PII when possible
- Use pseudonymization for user identifiers (future)

#### Data Retention
- Follow retention policies
- Delete logs after retention period
- Document data deletion procedures

### Audit Requirements

#### Access Logging
- Log all access to Grafana
- Log all log queries (future)
- Retain access logs for compliance period

#### Change Tracking
- Track configuration changes
- Track dashboard modifications
- Maintain audit trail

## Best Practices

### 1. Log Sanitization

Implement log sanitization in services:

```python
# Python example
def sanitize_log_data(data):
    """Remove sensitive fields from log data"""
    sensitive_fields = ['password', 'api_key', 'token', 'secret']
    sanitized = data.copy()
    for field in sensitive_fields:
        if field in sanitized:
            sanitized[field] = '[REDACTED]'
    return sanitized
```

```go
// Go example
func sanitizeLogData(data map[string]interface{}) map[string]interface{} {
    sensitiveFields := []string{"password", "api_key", "token", "secret"}
    sanitized := make(map[string]interface{})
    for k, v := range data {
        for _, field := range sensitiveFields {
            if strings.Contains(strings.ToLower(k), field) {
                sanitized[k] = "[REDACTED]"
                continue
            }
        }
        sanitized[k] = v
    }
    return sanitized
}
```

### 2. Environment-Specific Logging

- **Development**: More verbose logging (DEBUG level)
- **Production**: Minimal logging (INFO level and above)
- **Staging**: Similar to production with some DEBUG

### 3. Log Rotation

- Configure Docker log rotation
- Set maximum log file sizes
- Limit log retention in Docker

### 4. Monitoring Access

- Monitor Grafana access logs
- Alert on unusual access patterns
- Review access regularly

## Configuration Security

### Environment Variables

Store sensitive configuration in environment variables:

```bash
# .env file (do not commit to git)
GRAFANA_ADMIN_PASSWORD=secure-password-here
GRAFANA_SECRET_KEY=generate-secure-key
```

### Secrets Management

For production:
- Use secrets management (HashiCorp Vault, AWS Secrets Manager, etc.)
- Rotate secrets regularly
- Never commit secrets to version control

### Docker Secrets

Use Docker secrets for sensitive data:
```yaml
services:
  grafana:
    secrets:
      - grafana_admin_password
secrets:
  grafana_admin_password:
    external: true
```

## Security Checklist

### Initial Setup
- [x] Change default Grafana admin password
- [ ] Configure HTTPS for Grafana (production)
- [ ] Set up reverse proxy (production)
- [ ] Configure firewall rules
- [ ] Enable access logging

### Ongoing Maintenance
- [ ] Review access logs monthly
- [ ] Rotate passwords quarterly
- [ ] Audit user access annually
- [ ] Review and update security policies
- [ ] Test backup and restore procedures

### Compliance
- [ ] Document data retention policies
- [ ] Implement log deletion procedures
- [ ] Configure audit logging
- [ ] Review GDPR compliance
- [ ] Document security procedures

## Incident Response

### Security Incident Procedures

1. **Detect**: Monitor for unusual access patterns
2. **Contain**: Restrict access if breach detected
3. **Investigate**: Review logs and access patterns
4. **Remediate**: Fix security issues
5. **Document**: Record incident and response

### Log Access in Incidents

- Preserve logs during security incidents
- Export relevant logs for investigation
- Maintain chain of custody
- Document all access to logs

## Future Enhancements

1. **Authentication**
   - Add authentication to Loki API
   - Implement OAuth2/SSO for Grafana
   - Add API key authentication

2. **Encryption**
   - Encrypt Loki data at rest
   - Encrypt Grafana data at rest
   - Use encrypted volumes

3. **Access Control**
   - Fine-grained access control
   - Service-level access restrictions
   - Query-level permissions

4. **Audit Logging**
   - Log all Grafana queries
   - Log all Loki API access
   - Maintain audit trail

5. **Data Protection**
   - Automatic PII detection
   - Log redaction capabilities
   - Data deletion procedures

## References

- [Grafana Security Documentation](https://grafana.com/docs/grafana/latest/setup-grafana/configure-security/)
- [Loki Security Best Practices](https://grafana.com/docs/loki/latest/operations/security/)
- [OWASP Logging Cheat Sheet](https://cheatsheetseries.owasp.org/cheatsheets/Logging_Cheat_Sheet.html)
