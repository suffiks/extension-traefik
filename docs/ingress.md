---
category: Basic
group: Ingress
title: Overview
weight: 1
---

# Ingress

Ingresses are a way to expose `Application`-services to the internet.

## Example

```yaml
ingresses:
	- host: my-app.example.com
	- host: other.example.com
		paths:
			- /subpath
```
