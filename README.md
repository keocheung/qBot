# qBot
A custom robot for qBittorrent

## Install
### Docker Compose
```yaml
version: '3'
services:
  qbot:
    image: ghcr.io/keocheung/qbot
    container_name: qbot
    volumes:
      - ./config.yaml:/config/config.yaml # YAML config file location. Support YAML & JSON
    environment:
      - CONFIG_PATH=/config/config.yaml # Optional. Default is ./config.yaml
    network_mode: bridge
    restart: unless-stopped
```

## Usage
### Config File Example
qBot uses [Expr](https://github.com/antonmedv/expr) as condition expression parser. Please refer to [Expr](https://github.com/antonmedv/expr) for more documents.
```yaml
web_url: http://127.0.0.1:8080
api_key: abc123
get_torrents:
  limit: 20
  sort: added_on
  reverse: true
rules:
  # Set the max ratio to 2.0 for torrents that match the condition.
  - condition: "(Category == 'Movie' || any(Tags, {# == 'Public' || # == 'Private'})) && MaxRatio == -1"
    action:
      max_ratio: 2.0
```
More examples are in [example](./example)