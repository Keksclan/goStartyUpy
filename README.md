# goStartyUpy

> **Zero-Dependency Go Library für produktionsreife Startup-Banner mit Build-Metadaten, Runtime-Informationen und strukturierten Health-Checks.**

goStartyUpy ist ein wiederverwendbares Go-Modul, das beim Start eines beliebigen Go-Services einen klar strukturierten, deterministischen Startup-Banner rendert. Es kombiniert Build-Metadaten (Version, Commit, Branch …), Runtime-Informationen (Go-Version, OS/Arch, PID …) und optionale Dependency-Checks (SQL, TCP, HTTP, Redis …) in einer einzigen, sofort lesbaren Konsolenausgabe.

**Keine externen Abhängigkeiten.** Alles basiert ausschließlich auf der Go-Standardbibliothek (`stdlib`). Das Modul fügt deiner `go.sum` exakt **null** Einträge hinzu.

---

## Tags / Topics

`golang` · `go-library` · `banner` · `startup-banner` · `cli` · `microservices` · `health-check` · `startup-checks` · `zero-dependencies` · `devops` · `ascii-art` · `build-metadata` · `spring-boot-style` · `production-ready` · `deterministic`

---

## Feature-Übersicht

| Feature | Beschreibung |
|---------|-------------|
| **6 Banner-Styles** | `spring` (Standard), `classic`, `box`, `mini`, `block`, oder eigenes ASCII-Art via `Banner`-Feld |
| **Build-Metadaten** | Version, BuildTime, Commit, Branch, Dirty – injiziert zur Compile-Zeit via `-ldflags` |
| **Runtime-Info** | Go-Version, OS/Arch, PID – automatisch erfasst zur Laufzeit |
| **4 Built-in Checks** | `SQLPingCheck`, `TCPDialCheck`, `HTTPGetCheck`, `RedisPingCheck` – alle ohne externe Deps |
| **Custom Checks** | `checks.New()`, `checks.Bool()`, `checks.NewGroup()` – oder eigenes `Check`-Interface |
| **Parallel & Sequential** | `Runner` unterstützt beide Modi mit konfigurierbarem Per-Check-Timeout |
| **Environment-Erkennung** | Automatisch via `GO_STARTYUPY_ENV` Umgebungsvariable, wenn nicht explizit gesetzt |
| **ANSI-Farben** | Optional via `Color: true` – standardmäßig reiner Text ohne Escape-Sequenzen |
| **ASCII-Only-Modus** | `ASCIIOnly: true` ersetzt Unicode-Box-Zeichen durch plain ASCII (`+`, `-`, `\|`) |
| **Banner-Breite** | `BannerWidth` beschneidet jede Zeile auf eine Maximalbreite |
| **Deterministisch** | Stabile Ausgabereihenfolge, kein Zufall, keine Seiteneffekte |
| **Panic-Safe** | Alle Fehler werden abgefangen und als strukturierte `Result`-Objekte zurückgegeben |

---

## Installation

```bash
go get github.com/keksclan/goStartyUpy
```

**Voraussetzung:** Go 1.24 oder neuer.

Das Modul hat **keine transitiven Abhängigkeiten**. Nach `go get` enthält deine `go.sum` nur den einen Eintrag für `goStartyUpy` selbst.

### Import-Pfade

```go
import (
    "github.com/keksclan/goStartyUpy/banner"   // Banner-Rendering, Options, BuildInfo
    "github.com/keksclan/goStartyUpy/checks"    // Health-Checks, Runner, Check-Interface
    "github.com/keksclan/goStartyUpy/version"   // Modul-Version (z.B. "0.1.0")
)
```

- **`banner`** — Hauptpaket. Enthält `Options`, `BuildInfo`, `Render()`, `RenderWithChecks()`, alle Banner-Funktionen und die Build-Metadaten-Variablen.
- **`checks`** — Health-Check-System. Enthält `Check`-Interface, `Runner`, alle Built-in-Checks und Hilfskonstruktoren.
- **`version`** — Exponiert die Modul-Version als `ModuleVersion`-Konstante.

---

## Quickstart

### Minimaler Banner (ohne Checks)

```go
package main

import (
    "fmt"

    "github.com/keksclan/goStartyUpy/banner"
)

func main() {
    opts := banner.Options{
        ServiceName: "my-service",
    }
    info := banner.CurrentBuildInfo()
    fmt.Print(banner.Render(opts, info))
}
```

Das erzeugt einen Spring-Boot-ähnlichen ASCII-Art-Wordmark mit allen erkannten Build-/Runtime-Metadaten.

### Vollständiges Beispiel (Banner + Checks)

```go
package main

import (
    "context"
    "fmt"
    "os"

    "github.com/keksclan/goStartyUpy/banner"
    "github.com/keksclan/goStartyUpy/checks"
)

func main() {
    opts := banner.Options{
        ServiceName: "order-service",
        Environment: "production",
        Extra: map[string]string{
            "HTTP":  ":8080",
            "gRPC":  ":9090",
        },
    }
    info := banner.CurrentBuildInfo()

    runner := checks.DefaultRunner() // 2s Timeout, parallel
    results := runner.Run(context.Background(),
        checks.New("env-DATABASE_URL", func(ctx context.Context) error {
            if os.Getenv("DATABASE_URL") == "" {
                return fmt.Errorf("DATABASE_URL is not set")
            }
            return nil
        }),
        checks.TCPDialCheck{Address: "localhost:5432", Label: "postgres-tcp"},
        checks.HTTPGetCheck{URL: "http://localhost:8080/healthz", Label: "self-http"},
    )

    fmt.Print(banner.RenderWithChecks(opts, info, results))
}
```

### Kompilieren mit Build-Metadaten

```bash
make build PKG=./cmd/myservice BIN=bin/myservice
```

Oder direkt:

```bash
go build -ldflags "$(./scripts/ldflags.sh)" ./cmd/myservice/
```

---

## Build-Metadaten (`-ldflags`)

Das `banner`-Paket stellt fünf Link-Time-Variablen bereit, die zur Compile-Zeit via `-ldflags` injiziert werden. Diese Werte erscheinen automatisch im gerenderten Banner.

| Variable | Typ | Beschreibung | Default |
|----------|-----|-------------|---------|
| `banner.Version` | `string` | Semantische Version oder `git describe`-Ausgabe (z.B. `v1.2.3`, `v1.2.3-4-gabcdef1`) | `"dev"` |
| `banner.BuildTime` | `string` | UTC-Zeitstempel des Builds im RFC 3339 Format (z.B. `2026-03-05T18:00:00+01:00`) | `"unknown"` |
| `banner.Commit` | `string` | Kurzer Git-Commit-Hash (z.B. `abcdef1`) | `"unknown"` |
| `banner.Branch` | `string` | Git-Branch, auf dem kompiliert wurde (z.B. `master`, `feature/foo`) | `"unknown"` |
| `banner.Dirty` | `string` | `"true"` wenn der Working Tree beim Build uncommitted Changes hatte, sonst `"false"` | `"false"` |

**Wie funktioniert das?**

Go's Linker erlaubt es, Variablenwerte zur Compile-Zeit zu überschreiben. Die `-X`-Flags setzen die Paketvariablen direkt im kompilierten Binary, ohne den Quellcode zu ändern. `CurrentBuildInfo()` liest diese Werte dann zur Laufzeit aus.

### Makefile (empfohlen)

Das mitgelieferte `Makefile` sammelt Git-Metadaten automatisch über `scripts/ldflags.sh`:

```bash
make build-example   # Kompiliert das Example-Binary mit allen Metadaten
make run-example     # Kompiliert und startet das Example
make test            # Führt alle Unit-Tests aus (go test ./...)
make lint            # go vet + gofmt-Check auf allen Paketen
make clean           # Entfernt Build-Artefakte (bin/)
```

Für deinen eigenen Service:

```bash
make build PKG=./cmd/myservice BIN=bin/myservice
```

### Manueller Build

Wenn du das Makefile nicht verwenden möchtest, kannst du die ldflags direkt setzen:

```bash
VERSION=$(git describe --tags --always --dirty)
COMMIT=$(git rev-parse --short HEAD)
BRANCH=$(git rev-parse --abbrev-ref HEAD)
BUILD_TIME=$(date -Iseconds)
DIRTY=$(git diff --quiet && echo "false" || echo "true")

go build -ldflags "\
  -X 'github.com/keksclan/goStartyUpy/banner.Version=${VERSION}' \
  -X 'github.com/keksclan/goStartyUpy/banner.BuildTime=${BUILD_TIME}' \
  -X 'github.com/keksclan/goStartyUpy/banner.Commit=${COMMIT}' \
  -X 'github.com/keksclan/goStartyUpy/banner.Branch=${BRANCH}' \
  -X 'github.com/keksclan/goStartyUpy/banner.Dirty=${DIRTY}'" \
  ./cmd/myservice/
```

### `scripts/ldflags.sh` Helper

Das POSIX-sh-kompatible Script gibt den vollständigen ldflags-String aus – ideal für CI/CD-Pipelines oder andere Build-Systeme:

```bash
# Standard-Nutzung:
LDFLAGS="$(./scripts/ldflags.sh)" go build -ldflags "$LDFLAGS" ./cmd/myservice

# Eigenen Modul-Pfad überschreiben (falls dein Import-Pfad abweicht):
MODULE=github.com/my/repo ./scripts/ldflags.sh
```

---

## Banner-Styles

Die Library unterstützt **6 verschiedene Banner-Styles**, gesteuert über `Options.BannerStyle`. Wenn `Options.Banner` leer ist, bestimmt der Style den automatisch generierten Banner. Wenn `Options.Banner` gesetzt ist, wird der Wert **direkt verwendet** (Raw-Modus) und `BannerStyle` wird ignoriert.

**Style-Übersicht:**

| Style | `BannerStyle`-Wert | Höhe | Font-Technik | Direkte Funktion |
|-------|-------------------|------|-------------|-----------------|
| Spring (Standard) | `"spring"` oder `""` | 5 Zeilen | Unterstriche / Pipes / Slashes | `SpringLikeBanner(name, asciiOnly)` |
| Classic | `"classic"` | 5 Zeilen | Slashes / Backslashes / Underscores | `ClassicLikeBanner(name, asciiOnly)` |
| Box | `"box"` | 3 Zeilen | Unicode-Box-Zeichen (oder ASCII) | `BoxBanner(name, asciiOnly)` |
| Mini | `"mini"` | 3 Zeilen | Kompakte ASCII-Glyphen | `MiniBanner(name, asciiOnly)` |
| Block | `"block"` | 5 Zeilen | Dicke `#`-Zeichen | `BlockBanner(name, asciiOnly)` |
| Custom (Raw) | — | beliebig | Eigenes ASCII-Art | — |

**Zeichenunterstützung (alle Built-in-Fonts):**

Jeder eingebaute Font unterstützt die gleichen Zeichen: **A–Z**, **0–9**, **`-`**, **`_`** und **Leerzeichen**. Der `ServiceName` wird automatisch in Großbuchstaben umgewandelt. Nicht unterstützte Zeichen werden durch ein **`?`-Fallback-Glyph** ersetzt. Leerzeichen werden zu `-` normalisiert.

### Spring-Style (Standard)

`BannerStyle: "spring"` (oder leer, da `"spring"` der Default ist) erzeugt einen großen ASCII-Art-Wordmark, inspiriert vom Spring-Boot-Startup-Banner. Der Font verwendet Unterstriche (`_`), Pipes (`|`), und Slashes (`/`, `\`).

Unterhalb des Wordmarks steht die Tagline `:: goStartyUpy ::` mit optionalem Environment-Suffix (nur wenn via `GO_STARTYUPY_ENV` erkannt, siehe [Environment-Erkennung](#environment-erkennung-go_startyupy_env)).

```go
opts := banner.Options{
    ServiceName: "my-svc",
    // BannerStyle ist standardmäßig "spring"
}
```

**Beispiel-Ausgabe:**

```
 __  __  __   __         ____   __     __  ____
|  \/  | \ \ / /        / ___| \ \   / / / ___|
| |\/| |  \ V /  _____  \___ \  \ \ / / | |
| |  | |   | |  |_____|  ___) |  \ V /  | |___
|_|  |_|   |_|          |____/    \_/    \____|

 :: goStartyUpy ::
```

**Direkter Aufruf** (ohne `Options`/`Render()`):

```go
art := banner.SpringLikeBanner("my-svc", false)
fmt.Println(art)
```

### Classic-Style

`BannerStyle: "classic"` erzeugt einen Banner mit einer Slash/Backslash/Underscore-Schriftart, die an traditionelle Java-Framework-Startup-Banner erinnert. Unterhalb des Wordmarks werden **zwei konfigurierbare Taglines** gedruckt:

| Tagline | Default-Wert | Beispiel |
|---------|-------------|---------|
| `Tagline1` | `"<ServiceName> <Version>"` | `"my-service v1.2.3"` |
| `Tagline2` | `"Build: <BuildTime>  Commit: <Commit> [Branch] [Dirty]"` | `"Build: 2026-03-05  Commit: abcdef1  Branch: master"` |

Beide Taglines können über `Options.Tagline1` und `Options.Tagline2` überschrieben werden:

```go
opts := banner.Options{
    ServiceName: "my-svc",
    BannerStyle: "classic",
    Tagline1:    "My Service v2.0.0",
    Tagline2:    "Powered by goStartyUpy",
}
```

**`ShowDetails`-Option:** Steuert, ob der Key/Value-Info-Block (Service, Version, Go-Version etc.) im Classic-Modus angezeigt wird. Es handelt sich um einen `*bool`-Pointer. Standardmäßig werden Details angezeigt (`nil` = true). Explizit auf `false` setzen, um sie auszublenden:

```go
hide := false
opts := banner.Options{
    ServiceName: "my-svc",
    BannerStyle: "classic",
    ShowDetails: &hide,   // Details-Block wird nicht gedruckt
}
```

**Direkter Aufruf:**

```go
art := banner.ClassicLikeBanner("my-svc", false)
fmt.Println(art)
```

### Box-Style

`BannerStyle: "box"` erzeugt den klassischen Box-Banner mit Unicode-Box-Zeichen (`┌`, `─`, `┐`, `│`, `└`, `┘`). Der Servicename wird zentriert in der Box angezeigt.

```go
opts := banner.Options{
    ServiceName: "my-service",
    BannerStyle: "box",
}
```

**Ausgabe:**

```
┌───────────────────────────┐
│        MY-SERVICE         │
└───────────────────────────┘
```

**ASCII-Only-Modus:** Setze `Options.ASCIIOnly = true`, um Unicode-Box-Zeichen durch plain ASCII zu ersetzen (`+`, `-`, `|`):

```
+---------------------------+
|        MY-SERVICE         |
+---------------------------+
```

**Direkter Aufruf:**

```go
art := banner.BoxBanner("my-service", true) // true = ASCII-only
fmt.Println(art)
```

### Mini-Style

`BannerStyle: "mini"` erzeugt einen **kompakten 3-zeiligen** ASCII-Art-Wordmark. Ideal für schmale Terminals oder Logs, in denen wenig vertikaler Platz zur Verfügung steht.

```go
opts := banner.Options{
    ServiceName: "go",
    BannerStyle: "mini",
}
```

**Beispiel-Ausgabe (`"GO"`):**

```
 __  _
| _ | |
|__||_|
```

**Direkter Aufruf:**

```go
art := banner.MiniBanner("my-svc", false)
fmt.Println(art)
```

### Block-Style

`BannerStyle: "block"` erzeugt einen **dicken 5-zeiligen** ASCII-Art-Wordmark, bei dem jeder Buchstabe aus `#`-Zeichen aufgebaut ist. Gut sichtbar auch in lauten Log-Ausgaben.

```go
opts := banner.Options{
    ServiceName: "go",
    BannerStyle: "block",
}
```

**Beispiel-Ausgabe (`"GO"`):**

```
  ####   ###
 #      #   #
 #  ##  #   #
 #   #  #   #
  ####   ###
```

**Direkter Aufruf:**

```go
art := banner.BlockBanner("my-svc", false)
fmt.Println(art)
```

### Custom Banner (Raw)

Wenn du dein eigenes ASCII-Art verwenden möchtest, setze einfach `Options.Banner`. Der Wert wird **direkt verwendet**, ohne jede Verarbeitung. `BannerStyle` wird in diesem Fall ignoriert.

```go
opts := banner.Options{
    ServiceName: "my-service",
    Banner: `
   ╔═══════════════════════════════════╗
   ║     ★  MY AWESOME SERVICE  ★     ║
   ╚═══════════════════════════════════╝`,
}
```

**Tipp:** Du kannst Tools wie [patorjk.com/software/taag](http://patorjk.com/software/taag/) nutzen, um eigene ASCII-Art-Fonts zu generieren und sie als `Banner`-String einzufügen.

### Banner-Breite (`BannerWidth`)

Setze `Options.BannerWidth` auf einen positiven Integer, um jede Banner-Zeile auf diese Maximalbreite **hart abzuschneiden**. Ein Wert von `0` (Standard) bedeutet keine Beschränkung.

```go
opts := banner.Options{
    ServiceName: "my-very-long-service-name",
    BannerWidth: 60, // Jede Zeile wird nach 60 Zeichen abgeschnitten
}
```

Das ist nützlich, wenn der generierte Banner zu breit für dein Terminal oder dein Log-System ist.

---

## Environment-Erkennung (`GO_STARTYUPY_ENV`)

goStartyUpy unterstützt eine **automatische Environment-Erkennung** über die Umgebungsvariable `GO_STARTYUPY_ENV`. Das Verhalten ist wie folgt:

| Szenario | `Options.Environment` | `GO_STARTYUPY_ENV` | Ergebnis im Banner |
|----------|----------------------|--------------------|--------------------|
| Explizit gesetzt | `"production"` | egal | Kein Suffix angezeigt |
| Aus Env-Var erkannt | `""` (leer) | `"staging"` | Suffix `(staging)` wird angezeigt |
| Nichts gesetzt | `""` (leer) | nicht gesetzt / leer | Kein Suffix angezeigt |

**Regel:** Der Environment-Suffix (z.B. `(staging)`, `(dev)`) erscheint im Banner-Header **nur dann**, wenn der Wert tatsächlich aus der Umgebungsvariable `GO_STARTYUPY_ENV` stammt. Wenn `Options.Environment` explizit im Code gesetzt wird, wird **kein** Suffix angezeigt – der Wert wird intern verwendet, aber nicht im Banner sichtbar gemacht.

**Warum dieses Design?**

- Explizit gesetzte Werte im Code sind dem Entwickler bekannt → kein visueller Hinweis nötig.
- Werte aus Umgebungsvariablen könnten unerwartet sein (z.B. falsche Konfiguration in einer CI/CD-Pipeline) → visueller Hinweis im Banner hilft beim Debugging.

**Beispiel mit Umgebungsvariable:**

```bash
export GO_STARTYUPY_ENV=staging
go run ./cmd/myservice/
```

Der Banner zeigt dann:

```
 :: goStartyUpy :: (staging)
```

**Beispiel ohne Umgebungsvariable (explizit):**

```go
opts := banner.Options{
    ServiceName: "my-service",
    Environment: "production",  // Explizit → kein Suffix im Banner
}
```

---

## Options-Referenz (vollständig)

Die `banner.Options`-Struct steuert alle Aspekte des Banner-Renderings. Hier ist jedes Feld mit Typ, Default-Wert und Erklärung dokumentiert:

| Feld | Typ | Default | Beschreibung |
|------|-----|---------|-------------|
| `ServiceName` | `string` | `""` | Name des Services. Wird im Banner und in der Info-Sektion angezeigt. Wenn leer, wird `"SERVICE"` als Fallback verwendet. |
| `Environment` | `string` | `""` | Laufzeitumgebung (z.B. `"production"`, `"staging"`). Wenn leer, wird `GO_STARTYUPY_ENV` geprüft. Erscheint in der Info-Sektion. |
| `Banner` | `string` | `""` | Eigenes ASCII-Art. Wenn gesetzt, wird die automatische Banner-Generierung übersprungen und dieser Text direkt verwendet. |
| `BannerStyle` | `string` | `"spring"` | Steuert den automatisch generierten Banner-Style: `"spring"`, `"classic"`, `"box"`, `"mini"`, `"block"`. Wird ignoriert wenn `Banner` gesetzt ist. |
| `BannerWidth` | `int` | `0` | Maximale Breite pro Banner-Zeile. `0` = keine Beschränkung. Positive Werte schneiden jede Zeile hart ab. |
| `Separator` | `string` | `"═"` (Unicode) | Zeichen für die Trennlinie zwischen Banner und Info-Sektion. Im ASCII-Only-Modus wird `"="` verwendet. |
| `ASCIIOnly` | `bool` | `false` | Wenn `true`, werden alle Unicode-Zeichen (Box-Drawing, Separator) durch plain ASCII ersetzt. |
| `Color` | `bool` | `false` | Wenn `true`, wird die Ausgabe mit ANSI-Escape-Sequenzen eingefärbt. Standardmäßig reiner Text. |
| `Extra` | `map[string]string` | `nil` | Zusätzliche Key/Value-Paare, die in der Info-Sektion angezeigt werden (z.B. `"HTTP": ":8080"`). |
| `Tagline1` | `string` | `""` | Überschreibt die erste Tagline im Classic-Style. Wenn leer, wird der Default generiert. |
| `Tagline2` | `string` | `""` | Überschreibt die zweite Tagline im Classic-Style. Wenn leer, wird der Default generiert. |
| `ShowDetails` | `*bool` | `nil` | Steuert die Anzeige des Detail-Blocks im Classic-Style. `nil` = anzeigen, `&false` = ausblenden. |

**Interne Felder** (nicht direkt vom Benutzer gesetzt):

| Feld | Typ | Beschreibung |
|------|-----|-------------|
| `EnvironmentFromEnv` | `bool` | Wird intern auf `true` gesetzt, wenn das Environment aus `GO_STARTYUPY_ENV` stammt. Steuert die Suffix-Anzeige. |

---

## Check-System

Das `checks`-Paket bietet ein vollständiges Startup-Check-System, um Abhängigkeiten (Datenbanken, Caches, HTTP-Dienste) vor dem Akzeptieren von Traffic zu verifizieren.

### `Check`-Interface

Jeder Startup-Check implementiert das `Check`-Interface:

```go
type Check interface {
    Name() string                        // Menschenlesbarer Name des Checks
    Run(ctx context.Context) Result      // Führt den Check aus, gibt Result zurück
}
```

Das `Result`-Struct enthält das Ergebnis:

```go
type Result struct {
    Name     string        // Name des Checks
    OK       bool          // true = bestanden, false = fehlgeschlagen
    Duration time.Duration // Ausführungsdauer
    Error    string        // Fehlermeldung (leer bei Erfolg)
}
```

**Wichtig:** Checks paniken **niemals**. Alle Panics innerhalb von Check-Funktionen werden automatisch abgefangen und als `Result` mit `OK: false` und einer entsprechenden Fehlermeldung zurückgegeben.

### `Runner` — Check-Ausführung

Der `Runner` führt Checks mit konfigurierbarem Timeout aus. Er unterstützt sowohl **parallele** als auch **sequenzielle** Ausführung:

```go
runner := checks.Runner{
    TimeoutPerCheck: 2 * time.Second,  // Timeout pro einzelnem Check
    Parallel:        true,              // true = parallel, false = sequenziell
}
results := runner.Run(ctx, check1, check2, check3)
```

| Feld | Typ | Default | Beschreibung |
|------|-----|---------|-------------|
| `TimeoutPerCheck` | `time.Duration` | `0` | Timeout pro Check. `0` = kein zusätzliches Timeout (nur der übergebene Context). |
| `Parallel` | `bool` | `false` | `true` = alle Checks laufen in eigenen Goroutinen gleichzeitig. `false` = sequenzielle Ausführung in Eingabereihenfolge. |

**`DefaultRunner()`** gibt einen vorkonigurierten Runner mit 2 Sekunden Timeout und paralleler Ausführung zurück:

```go
runner := checks.DefaultRunner()
// Äquivalent zu:
// runner := checks.Runner{TimeoutPerCheck: 2 * time.Second, Parallel: true}
```

**Ergebnis-Reihenfolge:** Unabhängig vom Ausführungsmodus (parallel oder sequenziell) werden die Ergebnisse **immer in der gleichen Reihenfolge** wie die Eingabe-Checks zurückgegeben. Das macht die Ausgabe deterministisch und testbar.

### Built-in Checks (4 Typen)

#### `SQLPingCheck` — SQL-Datenbank

Pingt eine `*sql.DB`-Verbindung via `PingContext`. Nützlich für PostgreSQL, MySQL, SQLite und jeden anderen `database/sql`-kompatiblen Treiber.

```go
check := checks.SQLPingCheck{
    DB:        db,           // *sql.DB Handle (muss nicht nil sein)
    NameLabel: "postgres",   // Menschenlesbarer Name
}
```

| Feld | Typ | Beschreibung |
|------|-----|-------------|
| `DB` | `*sql.DB` | Das Datenbank-Handle. Wenn `nil`, schlägt der Check mit `"sql.DB is nil"` fehl. |
| `NameLabel` | `string` | Name des Checks in der Ausgabe. |

#### `TCPDialCheck` — TCP-Port

Prüft, ob ein TCP-Endpunkt erreichbar ist, indem eine Verbindung hergestellt und sofort wieder geschlossen wird. Ideal für Datenbanken, Caches oder andere TCP-basierte Dienste.

```go
check := checks.TCPDialCheck{
    Address: "localhost:5432",   // host:port Format
    Label:   "postgres-tcp",     // Menschenlesbarer Name
}
```

| Feld | Typ | Beschreibung |
|------|-----|-------------|
| `Address` | `string` | TCP-Adresse im `host:port`-Format (z.B. `"localhost:5432"`, `"redis:6379"`). |
| `Label` | `string` | Name des Checks in der Ausgabe. |

#### `HTTPGetCheck` — HTTP-Endpunkt

Führt einen HTTP-GET-Request aus und prüft, ob der Status-Code in einem erwarteten Bereich liegt. Nützlich für Health-Endpoints anderer Services.

```go
check := checks.HTTPGetCheck{
    URL:               "http://localhost:8080/healthz",
    Label:             "api-health",
    ExpectedStatusMin: 200,   // Optional, Default: 200
    ExpectedStatusMax: 299,   // Optional, Default: 399
}
```

| Feld | Typ | Default | Beschreibung |
|------|-----|---------|-------------|
| `URL` | `string` | — | Vollständige URL zum Proben (z.B. `"http://localhost:8080/healthz"`). |
| `Label` | `string` | — | Name des Checks in der Ausgabe. |
| `ExpectedStatusMin` | `int` | `200` | Untere Grenze (inklusiv) des akzeptablen Status-Code-Bereichs. |
| `ExpectedStatusMax` | `int` | `399` | Obere Grenze (inklusiv) des akzeptablen Status-Code-Bereichs. |

**Hinweis:** Der HTTP-Client verwendet kein eigenes Timeout – er verlässt sich auf den Context-Deadline des Runners, damit das Verhalten über alle Check-Typen konsistent ist.

#### `RedisPingCheck` — Redis via Raw TCP

Sendet einen RESP-kodierten `PING`-Befehl über eine rohe TCP-Verbindung und erwartet `+PONG` als Antwort. **Kein Redis-Client oder externe Dependency nötig** — funktioniert mit jedem Redis-kompatiblen Server.

```go
check := checks.RedisPingCheck{
    Address: "localhost:6379",  // host:port des Redis-Servers
    Label:   "redis-ping",     // Menschenlesbarer Name
}
```

| Feld | Typ | Beschreibung |
|------|-----|-------------|
| `Address` | `string` | TCP-Adresse des Redis-Servers im `host:port`-Format. |
| `Label` | `string` | Name des Checks in der Ausgabe. |

**Technisches Detail:** Der Check sendet das RESP-Array `*1\r\n$4\r\nPING\r\n` und erwartet `+PONG\r\n`. Bei unerwarteten Antworten schlägt der Check fehl mit der tatsächlichen Reply im Fehlertext.

### Custom Checks erstellen

#### Funktions-basierter Check (`checks.New`)

Der einfachste Weg, einen eigenen Check zu erstellen. Übergib einen Label-String und eine Funktion, die `error` zurückgibt (`nil` = bestanden):

```go
envCheck := checks.New("env-DATABASE_URL", func(ctx context.Context) error {
    if os.Getenv("DATABASE_URL") == "" {
        return fmt.Errorf("DATABASE_URL is not set")
    }
    return nil
})
```

Wenn `label` leer ist, wird `"custom"` als Fallback verwendet. Wenn `fn` `nil` ist, schlägt der Check immer mit `"nil check function"` fehl.

#### Boolean-Check (`checks.Bool`)

Für Checks, die ein Boolean-Ergebnis plus optionalen Fehler zurückgeben:

```go
featureFlag := checks.Bool("feature-flag", func(ctx context.Context) (bool, error) {
    return os.Getenv("ENABLE_NEW_UI") == "true", nil
})
```

Der Check besteht **nur**, wenn `ok == true` **und** `err == nil`. Wenn `ok == false` und `err == nil`, wird der Fehler `"check returned false"` erzeugt.

#### Gruppierter Check (`checks.NewGroup`)

Fasst mehrere Checks in einen einzelnen zusammen. Die Gruppe besteht **nur**, wenn **alle** Kinder bestehen:

```go
deps := checks.NewGroup("dependencies", checks.GroupOptions{
    Parallel:        true,                  // Kinder parallel ausführen
    TimeoutPerCheck: 3 * time.Second,       // Timeout pro Kind-Check
},
    checks.SQLPingCheck{DB: db, NameLabel: "postgres"},
    checks.TCPDialCheck{Address: "localhost:6379", Label: "redis-tcp"},
    checks.RedisPingCheck{Address: "localhost:6379", Label: "redis-ping"},
)
```

| `GroupOptions`-Feld | Typ | Default | Beschreibung |
|---------------------|-----|---------|-------------|
| `Parallel` | `bool` | `false` | `true` = Kind-Checks parallel ausführen. |
| `TimeoutPerCheck` | `time.Duration` | `0` | Timeout pro Kind-Check. `0` = kein zusätzliches Timeout. |

Bei Fehlern enthält der Error-String eine kompakte Zusammenfassung: `"2 failing: postgres: connection refused; redis-tcp: dial timeout"`.

#### Eigenes `Check`-Interface implementieren

Für komplexere Szenarien kannst du das `Check`-Interface direkt implementieren:

```go
type MyCustomCheck struct {
    // eigene Felder
}

func (c MyCustomCheck) Name() string { return "my-custom" }

func (c MyCustomCheck) Run(ctx context.Context) checks.Result {
    start := time.Now()
    // ... deine Logik ...
    return checks.Result{
        Name:     c.Name(),
        OK:       true,
        Duration: time.Since(start),
    }
}
```

### Parallele vs. sequenzielle Ausführung

| Modus | `Runner.Parallel` | Verhalten |
|-------|-------------------|-----------|
| Parallel | `true` | Jeder Check läuft in einer eigenen Goroutine. Der Runner wartet, bis alle fertig sind. Schnellster Gesamtdurchlauf. |
| Sequenziell | `false` | Checks laufen nacheinander in der Eingabereihenfolge. Ein langsamer Check blockiert die nachfolgenden. |

In **beiden Modi** werden die Ergebnisse in der **gleichen Reihenfolge** wie die Eingabe zurückgegeben.

---

## Modul-Version vs. Build-Version

goStartyUpy unterscheidet strikt zwischen **zwei verschiedenen Versionswerten**, die unabhängig voneinander sind:

| Wert | Paket | Zweck | Gesetzt durch |
|------|-------|-------|--------------|
| `version.ModuleVersion` | `version` | Release-Version der **Library** selbst (z.B. `"0.1.0"`) | Im Quellcode (`version/version.go`) |
| `banner.Version` | `banner` | Build-Version des **Service-Binaries** (z.B. `"v1.2.3"`) | `-ldflags` zur Compile-Zeit |

**Warum zwei Versionen?**

- `ModuleVersion` sagt dir, welche Version von goStartyUpy du als Dependency verwendest.
- `banner.Version` sagt dir, welche Version deines eigenen Services gerade läuft.

Beides sind unabhängige Werte. Dein Service kann `goStartyUpy@v0.1.0` verwenden und trotzdem als `v3.7.2` getaggt sein.

```go
import "github.com/keksclan/goStartyUpy/version"

fmt.Println("goStartyUpy Library:", version.ModuleVersion) // "0.1.0"
```

---

## `BuildInfo`-Struct

`CurrentBuildInfo()` erstellt ein `BuildInfo`-Struct mit allen Build- und Runtime-Metadaten. Dieses Struct wird an `Render()` und `RenderWithChecks()` übergeben:

```go
info := banner.CurrentBuildInfo()
```

| Feld | Typ | Quelle | Beschreibung |
|------|-----|--------|-------------|
| `Version` | `string` | `-ldflags` | Build-Version des Services (Default: `"dev"`) |
| `BuildTime` | `string` | `-ldflags` | UTC-Build-Zeitstempel (Default: `"unknown"`) |
| `Commit` | `string` | `-ldflags` | Kurzer Git-Commit-Hash (Default: `"unknown"`) |
| `Branch` | `string` | `-ldflags` | Git-Branch (Default: `"unknown"`) |
| `Dirty` | `string` | `-ldflags` | `"true"` / `"false"` für uncommitted Changes |
| `GoVersion` | `string` | `runtime.Version()` | Go-Version (z.B. `"go1.26"`) |
| `OS` | `string` | `runtime.GOOS` | Betriebssystem (z.B. `"linux"`, `"darwin"`) |
| `Arch` | `string` | `runtime.GOARCH` | CPU-Architektur (z.B. `"amd64"`, `"arm64"`) |
| `PID` | `int` | `os.Getpid()` | Prozess-ID des laufenden Binaries |

---

## Render-Funktionen

Das `banner`-Paket bietet zwei Haupt-Render-Funktionen:

### `Render(opts Options, info BuildInfo) string`

Erzeugt den vollständigen Startup-Banner **ohne Checks**. Gibt den gesamten Banner als String zurück (bereit für `fmt.Print`):

```go
output := banner.Render(opts, info)
fmt.Print(output)
```

### `RenderWithChecks(opts Options, info BuildInfo, results []checks.Result) string`

Erzeugt den vollständigen Startup-Banner **mit Check-Ergebnissen**. Die Check-Ergebnisse werden als Liste am Ende des Banners angehängt:

```go
output := banner.RenderWithChecks(opts, info, results)
fmt.Print(output)
```

**Ausgabeformat der Checks:**

```
Checks:
  [OK]   postgres (12ms)
  [OK]   redis-tcp (3ms)
  [FAIL] kafka: connection refused (1.2s)

Startup Complete
```

- `[OK]` = Check bestanden (grün bei `Color: true`)
- `[FAIL]` = Check fehlgeschlagen (rot bei `Color: true`), mit Fehlermeldung

---

## Öffentliche API-Stabilität

Die folgenden Elemente gelten als **öffentliche API** und unterliegen den Versionierungsgarantien:

- Alle exportierten Typen, Funktionen, Variablen und Konstanten in den Paketen `banner`, `checks` und `version`.
- Das `Check`-Interface und sein Vertrag.
- Die Felder der Structs `Options`, `BuildInfo`, `Result`, `Runner`, `GroupOptions`.
- Die Felder der Built-in-Check-Structs (`SQLPingCheck`, `TCPDialCheck`, `HTTPGetCheck`, `RedisPingCheck`).

**Nicht Teil der öffentlichen API** (können sich ohne Vorankündigung ändern):

- Alle unexportierten (kleingeschriebenen) Bezeichner.
- Das `example/`-Verzeichnis.
- Das `scripts/`-Verzeichnis.
- Interne Font-Daten und Render-Hilfsfunktionen.

---

## Versionierungsstrategie

Dieses Projekt folgt [Semantic Versioning 2.0.0](https://semver.org/):

| Version-Teil | Wann? | Beispiel |
|-------------|-------|---------|
| **MAJOR** (`X.0.0`) | Inkompatible API-Änderungen (Entfernen/Umbenennen exportierter Symbole, Änderung von Funktionssignaturen, Breaking Changes am `Check`-Interface) | `1.0.0` → `2.0.0` |
| **MINOR** (`0.X.0`) | Neue Features, abwärtskompatibel (neue Check-Typen, neue `Options`-Felder, neue Hilfsfunktionen) | `0.1.0` → `0.2.0` |
| **PATCH** (`0.0.X`) | Abwärtskompatible Bug-Fixes und Dokumentationskorrekturen | `0.1.0` → `0.1.1` |

**Hinweis:** Solange das Modul bei `0.x.y` ist, kann sich die API zwischen Minor-Versionen ändern. Ein `1.0.0`-Release signalisiert eine stabile API-Verpflichtung.

---

## Release-Prozess

1. **Version aktualisieren:** Setze `ModuleVersion` in `version/version.go` auf die neue Version.
2. **CHANGELOG aktualisieren:** Verschiebe Einträge von `[Unreleased]` in eine neue Versions-Sektion mit Datum.
3. **Committen:**
   ```bash
   git add -A
   git commit -m "release: v0.2.0"
   ```
4. **Taggen und pushen:**
   ```bash
   git tag v0.2.0
   git push origin master v0.2.0
   ```
5. **Konsumenten können die Version pinnen:**
   ```bash
   go get github.com/keksclan/goStartyUpy@v0.2.0
   ```

---

## Beispielprogramme

Das `example/`-Verzeichnis enthält lauffähige Programme für verschiedene Anwendungsfälle:

| Beispiel | Beschreibung | Starten mit |
|----------|-------------|------------|
| `example/` | Vollständige Demo: Custom Checks, Groups, Built-in Checks | `make run-example` |
| `example/simple/` | Minimaler Banner ohne Checks | `go run ./example/simple/` |
| `example/custom_banner/` | Eigenes ASCII-Art als Banner | `go run ./example/custom_banner/` |
| `example/ascii_only/` | ASCII-Only-Modus für Terminals ohne Unicode | `go run ./example/ascii_only/` |
| `example/checks_demo/` | Alle Built-in-Check-Typen (SQL, TCP, HTTP, Redis) | `go run ./example/checks_demo/` |
| `example/custom_checks/` | Funktions-basierte, Boolean- und gruppierte Checks | `go run ./example/custom_checks/` |
| `example/font_preview/` | Druckt den Big-Font ASCII-Wordmark für einen Service-Namen | `go run ./example/font_preview/` |

```bash
# Einfachster Start:
go run ./example/simple/

# Vollständige Demo mit Build-Metadaten:
make run-example
```

---

## Ausgabe-Beispiel (Box-Style mit Checks)

Das folgende Beispiel zeigt die vollständige Ausgabe im Box-Style mit Build-Metadaten, Extra-Feldern und Check-Ergebnissen:

```
┌──────────────────────────────┐
│        ORDER-SERVICE         │
└──────────────────────────────┘
════════════════════════════════════════════════════════════
  Service     : order-service
  Environment : staging
  Version     : v1.2.3
  BuildTime   : 2026-02-24T09:00:00Z
  Commit      : abcdef1
  Branch      : master
  Dirty       : false
  Go          : go1.26
  OS/Arch     : linux/amd64
  PID         : 12345
  HTTP        : :8080

Checks:
  [OK]   postgres (12ms)
  [OK]   redis-tcp (3ms)
  [OK]   self-http (8ms)
  [OK]   redis-ping (2ms)

Startup Complete
```

**Aufbau der Ausgabe:**

1. **Banner** — ASCII-Art oder Box (je nach Style)
2. **Separator** — Trennlinie (`═══...` oder `===...` im ASCII-Modus)
3. **Info-Sektion** — Key/Value-Paare (Service, Environment, Version, BuildTime, Commit, Branch, Dirty, Go, OS/Arch, PID, plus alle `Extra`-Einträge)
4. **Checks** (nur bei `RenderWithChecks`) — Ergebnisse mit `[OK]`/`[FAIL]` Status und Dauer
5. **Footer** — `"Startup Complete"` oder Check-Zusammenfassung

---

## Sicherheitshinweis

Der Banner druckt ausschließlich **sichere, nicht-geheime** Informationen (Version, Adressen, PID etc.).

⚠️ **Gib niemals Secrets** (Passwörter, Tokens, API-Keys) über `Options.Extra` oder andere Felder weiter. Der Aufrufer ist dafür verantwortlich, was gedruckt wird.

---

## Tests

```bash
# Alle Unit-Tests ausführen:
go test ./...

# Via Makefile (identisch):
make test

# Linting (go vet + gofmt):
make lint
```

Das Projekt enthält Tests für:
- Banner-Rendering aller 6 Styles
- Font-Rendering und Fallback-Glyphen
- Build-Metadaten-Snapshot
- Formatierung und Separator
- Environment-Erkennung (explizit, aus Env-Var, nicht gesetzt)
- Check-Runner (parallel und sequenziell)
- Alle Built-in-Check-Typen
- FuncCheck, Bool-Check, Group-Check
- Panic-Recovery
- Unicode- und Edge-Case-Sicherheit

---

## Lizenz

Dieses Projekt steht unter der **MIT License with Attribution Requirement** — siehe [LICENSE](LICENSE) für den vollständigen Text.

**Kurzfassung:**
- Open-Source-Nutzung: Frei, Attribution optional aber geschätzt.
- Kommerzielle/Corporate-Nutzung: MIT-kompatibel, Attribution-Erwähnung erforderlich.

---

## Used by

Dieses Projekt wird verwendet von:
- [Keksclan](https://github.com/Keksclan) — Creator von goStartyUpy

➕ **Dein Projekt/Unternehmen hier hinzufügen?** Öffne einen Pull Request und editiere [`USED_BY.md`](USED_BY.md). Regeln:
- Alphabetisch einordnen
- 1 Zeile pro Eintrag, kein Marketing
- Format: `- [Name](URL) - Kurze Beschreibung [tags]`

---

## Tags / Topics

Die folgenden Tags beschreiben dieses Projekt und können als GitHub-Topics verwendet werden:

`golang` · `go-library` · `banner` · `startup-banner` · `cli` · `microservices` · `health-check` · `startup-checks` · `zero-dependencies` · `devops` · `ascii-art` · `build-metadata` · `spring-boot-style` · `production-ready` · `deterministic`
