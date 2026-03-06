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

Ao fazer o build da imagem (produção), use **`make build`**: a versão vem do Git (mesma lógica do dev) e a imagem é taggeada como `ecommerce-users:$(VERSION)`.

```bash
make build   # usa git describe --tags --always, gera ecommerce-users:1.0.1 (ou dev)
```

O `make up` também usa essa versão: o Makefile exporta `VERSION` e o `docker-compose.yml` passa `args.VERSION` para o build do serviço users.

Para uma tag explícita (sem Git):

```bash
VERSION=1.0.0 make build
```

---

## Versão em dev (automática a partir do Git)

No ambiente de desenvolvimento (`make dev`), a versão **não é mais hardcoded**: o Makefile exporta `VERSION` usando `git describe --tags --always` (última tag + commits, ex.: `1.0.1` ou `1.0.1-2-gabc123`). O `docker-compose.dev.yml` usa `VERSION: "${VERSION:-dev}"`.

- Ao rodar **`make dev`**, a API exibe no `meta.version` o resultado do Git (ex.: `1.0.0`, `1.0.1`, `1.0.1-2-gabc123`).
- Se não houver tags ou o comando falhar, cai para `"dev"`.
- Não é necessário alterar nenhum arquivo ao criar uma nova tag; basta rodar `make dev` de novo.

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
