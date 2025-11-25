FROM quay.io/prometheus/prometheus:v3.2.1

COPY --chown=nobody:nobody prometheus.yml /etc/prometheus/prometheus.yml
COPY --chown=nobody:nobody scrape-configs /etc/prometheus/scrape-configs
