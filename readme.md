# ⚠️ Work in Progress ⚠️

Form handler written in Go.

## Features
- [x] Global configuration
- [x] Receive form submissions
- [x] Email form submissions
- [x] Email templating
- [x] Handle multiple forms
- [x] Global keyword blocklist for message field
- [x] Form configuration
	- [x] Designate fields, e.g. "name", "email", "message"
	- [x] Additional keyword blocklist
- [x] Honeypot field
- [x] Cloudflare Turnstile validation
- [x] Akismet validation
- [ ] Submission logging
	- [ ] Multiple levels, such as "spam", "email failed", "success", "all"
- [ ] Mailgun integration

## Development features
- [ ] End-to-end submission testing (send POST request to fohago, receive and verify email)
- [ ] Unit tests

Inspired by:
- [Mailbear](https://github.com/DenBeke/mailbear)
- [MailyGo](https://git.jlel.se/jlelse/MailyGo)