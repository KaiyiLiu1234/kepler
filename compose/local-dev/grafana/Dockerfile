FROM quay.io/ceph/grafana:10.4.2

COPY --chown=grafana:grafana datasource.yml /etc/grafana/provisioning/datasources/
COPY --chown=grafana:grafana dashboards.yml /etc/grafana/provisioning/dashboards/
COPY --chown=grafana:grafana dashboards /var/lib/grafana/dashboards
