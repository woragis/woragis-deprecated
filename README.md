# Woragis (deprecated)

**Arquivado.** Este repositório preserva o monólito em migração para microserviços.

| | |
|---|---|
| **Substituído por** | [woragis/woragis](https://github.com/woragis/woragis) — hub pessoal atual (`dev/woragis`) |
| **Microserviços backend** | Repos separados: `woragis-auth-backend`, `woragis-jobs-backend`, `woragis-posts-backend`, `woragis-management-backend`, etc. |
| **Status** | Sem desenvolvimento ativo — referência histórica |

## Estrutura (snapshot)

- `backend/` — orquestração + workers + submódulos de microserviços
- `frontend/` — frontends por domínio (jobs, posts, management)
- `socialmedia/` — dashboard de redes sociais
- `mobile/`, `landing/`, `docs/` — legado

Ver `SYSTEM_STATUS_REPORT.md` e `backend/README.md` para detalhes técnicos do último estado documentado.
