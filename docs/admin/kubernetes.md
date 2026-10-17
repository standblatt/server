# Kubernetes

The repository contains a Helm chart in `helm/`. It deploys the app as a Deployment with Service, Ingress (Traefik, TLS via cert-manager), optionally a PodDisruptionBudget and a PodMonitor.

## Values

| Value | Default | Description |
|---|---|---|
| `image.repository`, `image.tag` | empty | the image, e.g. tag `latest` (releases) or `bleeding_edge` (master) |
| `image.pullPolicy` | `Always` | keep `Always` when the tag moves, otherwise nodes keep their cached image |
| `imagePullSecrets` | `[]` | e.g. `[{name: regcred}]` for the private registry |
| `envFromSecret` | empty | existing Secret whose keys become environment variables (`DATABASE_URL`, `STANDBLATT_ADMIN_TOKEN`) |
| `env` | `STANDBLATT_TRUST_PROXY: "true"`, `STANDBLATT_LOG_LEVEL: info` | plain environment variables |
| `web.replicas` | `2` | replicas; with more than one a PodDisruptionBudget is created |
| `metrics.podMonitor.enabled` | `false` | needs the Prometheus operator CRDs |
| `ingress.hosts` | `[]` | host names |
| `ingress.tls.clusterIssuer` | `letsencrypt` | cert-manager ClusterIssuer |
| `priorityClassName` | empty | e.g. `production` or `staging` |

The pods run as non-root with a read-only root filesystem and all capabilities dropped.

## Example

```sh
kubectl -n standblatt create secret generic standblatt-env \
  --from-literal=DATABASE_URL='postgres://standblatt:secret@pg-rw.postgres:5432/standblatt?sslmode=require' \
  --from-literal=STANDBLATT_ADMIN_TOKEN="$(openssl rand -base64 32)"

helm upgrade --install standblatt-app helm -n standblatt \
  --set image.repository=ghcr.io/standblatt/server \
  --set image.tag=v0.3.0 \
  --set envFromSecret=standblatt-env \
  --set 'ingress.hosts[0]=standblatt.example.org'
```

Do not put a connection pooler in transaction mode between the app and PostgreSQL: every replica keeps one `LISTEN` connection for realtime events.

## The hosted installation

app.standblatt.ch runs on this chart through Flux: the GitOps repo of the cluster holds one HelmRelease per environment with the values inline; the chart comes from `oci://ghcr.io/standblatt/charts/standblatt` (pushed by CI on every tag); image and chart are public, so Flux and Keel need no credentials. The images move by Keel, not by Flux: staging follows `bleeding_edge`, production `latest`, and Keel rolls out a new digest within a minute on staging and five minutes on production. See [CI/CD](../developer/ci-cd.md).
