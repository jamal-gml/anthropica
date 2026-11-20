# Dockerfile patterns

All examples: BuildKit (`# syntax=docker/dockerfile:1`), pinned digests (replace `@sha256:...`), non-root, exec form.

## Node.js
```dockerfile
# syntax=docker/dockerfile:1
FROM node:22-slim@sha256:... AS deps
WORKDIR /app
COPY package.json package-lock.json ./
RUN --mount=type=cache,target=/root/.npm npm ci --omit=dev

FROM node:22-slim@sha256:... AS build
WORKDIR /app
COPY package.json package-lock.json ./
RUN --mount=type=cache,target=/root/.npm npm ci
COPY . .
RUN npm run build

FROM gcr.io/distroless/nodejs22-debian12:nonroot@sha256:...
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY --from=build /app/dist ./dist
USER 65532:65532
EXPOSE 8080
CMD ["dist/server.js"]
```
- `npm ci`, not `npm install`. Do not run via `npm start` (npm does not forward signals).

## Python
```dockerfile
FROM python:3.12-slim@sha256:... AS build
ENV PIP_NO_CACHE_DIR=1 PIP_DISABLE_PIP_VERSION_CHECK=1
WORKDIR /app
RUN python -m venv /venv
COPY requirements.txt .
RUN --mount=type=cache,target=/root/.cache/pip /venv/bin/pip install -r requirements.txt
COPY . .

FROM python:3.12-slim@sha256:...
ENV PYTHONDONTWRITEBYTECODE=1 PYTHONUNBUFFERED=1 PATH=/venv/bin:$PATH
RUN useradd --uid 10001 --no-create-home app
WORKDIR /app
COPY --from=build /venv /venv
COPY --from=build /app /app
USER 10001:10001
CMD ["gunicorn", "-b", "0.0.0.0:8080", "--graceful-timeout", "20", "app:app"]
```
- Wheels with C extensions: build in a full image, copy the venv to slim with the same Python and libc. Avoid alpine for Python (musl, slow builds).

## Go (static)
```dockerfile
FROM golang:1.23@sha256:... AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/app ./cmd/app

FROM gcr.io/distroless/static-debian12:nonroot@sha256:...
COPY --from=build /out/app /app
USER 65532:65532
ENTRYPOINT ["/app"]
```

## JVM
```dockerfile
FROM eclipse-temurin:21-jdk@sha256:... AS build
WORKDIR /src
COPY . .
RUN ./gradlew --no-daemon bootJar && java -Djarmode=layertools -jar build/libs/app.jar extract --destination /out

FROM eclipse-temurin:21-jre@sha256:...
RUN useradd --uid 10001 app
WORKDIR /app
COPY --from=build /out/dependencies/ ./
COPY --from=build /out/spring-boot-loader/ ./
COPY --from=build /out/application/ ./
USER 10001:10001
ENTRYPOINT ["java", "-XX:MaxRAMPercentage=75", "org.springframework.boot.loader.launch.JarLauncher"]
```
- Layered jars keep dependency layers cached across code changes.

## .dockerignore baseline
```
.git
.github
**/node_modules
**/__pycache__
**/*.pyc
dist
build
target
.env*
*.pem
*.key
coverage
Dockerfile*
docker-compose*
```

## Build secrets and private registries
```bash
docker buildx build --secret id=npmrc,src=$HOME/.npmrc .
# Dockerfile: RUN --mount=type=secret,id=npmrc,target=/root/.npmrc npm ci
docker buildx build --ssh default .     # RUN --mount=type=ssh git clone ...
```

## Reproducibility and provenance
- `docker buildx build --provenance=mode=max --sbom=true` attaches attestations; signing and SBOM policy belong to `secure-supply-chain`.
- `SOURCE_DATE_EPOCH` and pinned package versions for reproducible layers.
- Cache in CI: `--cache-from type=registry,ref=<reg>/<app>:buildcache --cache-to type=registry,ref=...,mode=max` (or `type=gha`).
- Multi-arch: `--platform linux/amd64,linux/arm64`; `exec format error` at runtime = arch mismatch.

## Compose (local/dev)
- `depends_on: { db: { condition: service_healthy } }` with a `healthcheck` on the dependency; plain `depends_on` only orders start.
- Named volumes for data; never bind-mount secrets from the repo. `.env` is for local defaults only and stays git-ignored.
- `HEALTHCHECK --interval=30s --timeout=3s --start-period=20s --retries=3 CMD ["/app","healthcheck"]` (distroless has no curl; ship a healthcheck subcommand or static binary).

## Image review checklist
- [ ] Final stage has no compilers, package managers or shells unless needed
- [ ] `docker history --no-trunc` shows no secrets, tokens or `.env`
- [ ] User is numeric non-root; app writes only to declared volumes/tmp
- [ ] Base pinned by digest; automated bumps configured
- [ ] Size reasonable (static Go < 30 MB, slim runtimes < 200 MB typical)
- [ ] Trivy: 0 fixable CRITICAL
- [ ] Labels: `org.opencontainers.image.source`, `.revision`, `.version`
