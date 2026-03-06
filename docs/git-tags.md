# Git tags — guia rápido

Tags marcam um ponto específico no histórico (geralmente uma release). Este projeto usa tags para versionar (ex.: `1.0.0`) e em dev você pode exibir algo como `dev-1.0.0` na API.

---

## Criar uma tag

**Tag anotada** (recomendada: guarda autor, data e mensagem):

```bash
# Na branch/commit que você quer marcar como release
git tag -a 1.0.0 -m "Release 1.0.0: login + register"
```

**Tag leve** (só um nome no commit):

```bash
git tag 1.0.0
```

Sempre que possível use **anotada** (`-a`); ela é a que costuma ser usada em releases e por ferramentas (CI, changelog).

---

## Listar tags

```bash
git tag                    # todas as tags
git tag -l "1.*"           # tags que batem com o padrão (ex.: 1.0, 1.1)
git show 1.0.0              # detalhes da tag (commit, mensagem se anotada)
```

---

## Enviar tags para o remoto

Tags **não** vão no `git push` normal. É preciso enviá-las:

```bash
git push origin 1.0.0       # uma tag
git push origin --tags     # todas as tags que ainda não estão no remoto
```

---

## Usar a tag na versão do build

Ao fazer o build da imagem (produção), injete a versão com a tag atual:

```bash
# Exemplo: versão a partir da tag mais recente
VERSION=$(git describe --tags --always)
docker build --build-arg VERSION="$VERSION" -f docker/services/users/Dockerfile .
```

Ou com tag explícita:

```bash
docker build --build-arg VERSION=1.0.0 -f docker/services/users/Dockerfile .
```

---

## Versão em dev (dev-1.0.0)

No ambiente de desenvolvimento (Air, `make dev`), o binário não é buildado com `-ldflags`, então a versão vem da **variável de ambiente** `VERSION`.

- No **docker-compose.dev.yml** está algo como `VERSION: "dev-1.0.0"`.
- Quando você criar a tag `1.0.0`, pode deixar em dev como `dev-1.0.0` para indicar “ambiente dev, alinhado à release 1.0.0”.
- Ao lançar a próxima release (ex.: `1.1.0`), altere para `VERSION: "dev-1.1.0"` se quiser refletir a última release no dev.

---

## Boas práticas (versionamento semântico)

- **MAJOR.MINOR.PATCH** (ex.: `1.2.3`):
  - **MAJOR**: mudanças incompatíveis na API.
  - **MINOR**: nova funcionalidade compatível.
  - **PATCH**: correções compatíveis.
- Prefira tags com `v`: `v1.0.0` (muitos projetos usam assim). O `git describe --tags` entende e você pode tirar o `v` no script se precisar (`1.0.0`).
- Não mova tags já publicadas; crie uma nova (ex.: `1.0.1`) em vez de refazer `1.0.0`.

---

## Resumo de comandos

| Ação              | Comando |
|-------------------|--------|
| Criar tag anotada | `git tag -a 1.0.0 -m "Mensagem"` |
| Listar tags       | `git tag` |
| Ver detalhes      | `git show 1.0.0` |
| Enviar uma tag    | `git push origin 1.0.0` |
| Enviar todas      | `git push origin --tags` |
| Última tag/commit | `git describe --tags --always` |
