# Cloudflare DDNS
A Dynamic DNS for Cloudflare that allows your zones to be constantly up-to-date with your constantly changing public IP. 

## Description
A simple Go application allows your DNS records to have their A and AAAA records updated with your current dynamic public IP.

You can specify your configuration by means of a `yaml` configuration file.

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