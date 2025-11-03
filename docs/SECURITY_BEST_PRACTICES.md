# Security Best Practices

## Input Validation

1. **Validate all inputs** before processing
2. **Use strict limits** for user-facing APIs
3. **Sanitize** all outputs before logging
4. **Type-check** all data

## SSRF Protection

1. **Enable SSRF protection** (enabled by default)
2. **Use allowlists** instead of blocklists when possible
3. **Validate URLs** before making requests
4. **Block private IPs** by default

## Resource Limits

1. **Set appropriate timeouts** for your use case
2. **Monitor limit violations** to detect abuse
3. **Use stricter limits** for untrusted workflows
4. **Test with limits** enabled

## Error Handling

1. **Don't expose internal details** in error messages
2. **Log security events** for monitoring
3. **Fail securely** on errors
4. **Sanitize errors** before returning to users

## Deployment

1. **Keep dependencies updated**
2. **Run security scans** regularly
3. **Use environment variables** for secrets
4. **Enable HTTPS** in production
5. **Implement rate limiting**
6. **Monitor for anomalies**

---

**Last Updated:** 2025-11-03
**Version:** 1.0
