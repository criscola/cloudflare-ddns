# Cloudflare DDNS
## Introduction

A Dynamic DNS application for Cloudflare that allows your zones to be constantly up-to-date with your constantly changing public IP.

Cloudflare is an excellent reverse proxy, providing powerful and complete free features like DNS resolving and CDN caching.

## Description
Cloudflare DDNS is a lightweight Go application allowing your A and AAAA records to be updated with your current dynamic public IP.

It works by asking 1.1.1.1 for your public IP and then updating your configured root domain with your new public IP, if necessary.
You can specify your configuration by means of a super simple `yaml` configuration file.

## TODOs
- [ ] IPv6 support
- [ ] Unit & integration testing
- [ ] Documentation
- [ ] Allow DNS records to be created if not present already
- [ ] Update DNS records concurrently instead of one after the other
- [x] Logging & error handling
- [ ] Allow configuring which DNS record to update manually
- [ ] CLI flags
- [ ] Docker deployment
- [ ] Kubernetes Helm Chart deployment
- [ ] API key and account email auth option