# Kepler Local Development Setup

This is a minimal Docker Compose setup for running Kepler from the local repository with Prometheus and Grafana on ARM Ampere Altra (Fedora).

## What's Included

- **Kepler**: Built from the current repository (local development)
- **Prometheus**: Configured to scrape Kepler metrics
- **Grafana**: Pre-configured with Prometheus datasource and custom dashboard support

## Prerequisites

- Docker and Docker Compose installed
- Running on ARM Ampere Altra with Fedora

## Quick Start

1. **Add your custom Grafana dashboards**:
   Place your dashboard JSON files in `./grafana/dashboards/kepler/` directory:
   ```bash
   mkdir -p grafana/dashboards/kepler
   # Copy your dashboard JSON files here
   cp /path/to/your/dashboard.json grafana/dashboards/kepler/
   ```

2. **Start the services**:
   ```bash
   docker compose up -d
   ```

3. **Access the services**:
   - Kepler metrics: http://localhost:28282/metrics
   - Prometheus: http://localhost:39090
   - Grafana: http://localhost:33000 (username: admin, password: admin)

## Directory Structure

```
local-dev/
├── compose.yaml              # Main docker compose configuration
├── kepler/
│   └── etc/kepler/
│       └── config.yaml       # Kepler configuration
├── prometheus/
│   ├── Dockerfile
│   ├── prometheus.yml        # Prometheus main config
│   └── scrape-configs/
│       └── kepler.yaml       # Kepler scrape configuration
└── grafana/
    ├── Dockerfile
    ├── datasource.yml        # Prometheus datasource config
    ├── dashboards.yml        # Dashboard provisioning config
    └── dashboards/           # Place your custom dashboards here
        └── kepler/           # Dashboard folder (auto-created in Grafana)
            └── *.json        # Your dashboard JSON files go here
```

## Customization

### Adding More Prometheus Scrape Targets

Create additional YAML files in `prometheus/scrape-configs/` directory. For example:

```yaml
# prometheus/scrape-configs/custom.yaml
scrape_configs:
  - job_name: my-custom-service
    static_configs:
      - targets: [my-service:9090]
```

### Modifying Kepler Configuration

Edit `kepler/etc/kepler/config.yaml` to adjust Kepler settings such as log level, metrics level, etc.

### Custom Grafana Dashboards

1. Place your dashboard JSON files in `grafana/dashboards/kepler/`
2. Restart Grafana: `docker compose restart grafana`
3. The dashboards will be automatically provisioned in the "kepler" folder

You can also create subdirectories within `grafana/dashboards/` and they will appear as separate folders in Grafana due to the `foldersFromFilesStructure: true` setting.

## Stopping the Services

```bash
docker compose down
```

To also remove the Prometheus data volume:
```bash
docker compose down -v
```

## Troubleshooting

### Kepler not collecting metrics
- Ensure the container is running with `privileged: true`
- Check that `/proc` and `/sys` are mounted correctly
- Check Kepler logs: `docker compose logs kepler`

### Grafana dashboards not appearing
- Ensure your dashboard JSON files are in `grafana/dashboards/kepler/`
- Check Grafana logs: `docker compose logs grafana`
- Verify file permissions are readable: `chmod -R a+rX grafana/dashboards/`

### Prometheus not scraping Kepler
- Verify the services are on the same network
- Check Prometheus targets: http://localhost:39090/targets
- Check Prometheus logs: `docker compose logs prometheus`
