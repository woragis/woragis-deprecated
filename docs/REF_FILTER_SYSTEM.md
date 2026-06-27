# Sistema de Filtro por Referência (Ref Filter System)

## Visão Geral

Sistema de filtragem de conteúdo na landing page baseado em parâmetros de URL (`?ref=xxx`) que permite mostrar apenas conteúdo relevante para uma área específica de expertise. O filtro persiste por 24 horas e ajuda a evitar a percepção de "generalista sem profundidade" ao focar o conteúdo exibido.

## Objetivo

Quando um recrutador ou visitante acessa a landing page através de um link específico (ex: `www.woragis.me/?ref=do` para DevOps), ele verá apenas conteúdo relacionado àquela área de expertise. Isso cria uma experiência personalizada e focada, demonstrando profundidade técnica na área de interesse.

## Códigos de Referência

| Código | Área | Descrição |
|--------|------|-----------|
| `do` | DevOps | Conteúdo relacionado a DevOps, CI/CD, infraestrutura, containers, etc. |
| `ds` | Data Science | Projetos e posts sobre machine learning, análise preditiva, modelos, etc. |
| `da` | Data Analysis | Análises de dados, visualizações, dashboards, ETL, etc. |
| `go` | Golang | Projetos backend em Go, APIs, microserviços, etc. |
| `py` | Python | Projetos e conteúdo em Python |
| `js` | JavaScript/TypeScript | Projetos frontend/backend em JS/TS |
| `nx` | Next.js | Projetos e artigos sobre Next.js |
| `sk` | SvelteKit | Conteúdo relacionado a SvelteKit |
| `aw` | AWS | Projetos e conhecimentos sobre AWS |
| `sb` | Spring Boot | Projetos Java/Spring Boot |

## Arquitetura do Sistema

### 1. Fluxo de Funcionamento

```
Usuário acessa → www.woragis.me/?ref=do
    ↓
Sistema detecta parâmetro `ref` na URL
    ↓
Valida código de referência (do, ds, da, etc.)
    ↓
Salva no localStorage/cookie com timestamp
    ↓
Filtra conteúdo baseado no código
    ↓
Exibe apenas conteúdo com tag correspondente
    ↓
Filtro expira após 24h
```

### 2. Persistência do Filtro

**Opção 1: LocalStorage (Recomendado)**
- **Vantagens**: Mais simples, não envia dados no header HTTP
- **Estrutura**:
  ```json
  {
    "refFilter": "do",
    "timestamp": 1234567890,
    "expiresAt": 1234654290
  }
  ```

**Opção 2: Cookies**
- **Vantagens**: Funciona mesmo com JavaScript desabilitado (se implementado no backend)
- **Estrutura**: Cookie `ref_filter` com valor `do|timestamp|expiresAt`

**Recomendação**: Usar LocalStorage com fallback para cookies se necessário.

### 3. Estrutura de Tags no Backend

Cada item de conteúdo (projeto, post, tech writing, etc.) deve ter:

#### Campos no Banco de Dados

```typescript
interface ContentItem {
  id: string;
  title: string;
  // ... outros campos
  
  // Sistema de tags
  tags: string[];              // Tags principais (do, ds, da, go, etc.)
  isPublic: boolean;           // Visível para todos (sem filtro)
  priority: number;            // Prioridade quando isPublic = true (1-10)
  
  // Metadados
  createdAt: Date;
  updatedAt: Date;
}
```

#### Lógica de Tags

- **Tags múltiplas**: Um projeto pode ter múltiplas tags (ex: `["go", "aw"]` para um projeto Go na AWS)
- **Tag `public`**: Conteúdo visível para todos, independente do filtro
- **Sem tags**: Conteúdo não aparece em nenhum filtro (apenas se `isPublic = true`)

### 4. Conteúdo Público (Default)

Quando o usuário acessa **sem** parâmetro `ref` ou após expiração do filtro, deve ver conteúdo "default" prioritário:

#### Distribuição Sugerida

| Área | Quantidade | Prioridade |
|------|------------|------------|
| Data Science | 5 projetos/posts | Alta |
| Data Analysis | 4 projetos/posts | Alta |
| DevOps | 1-2 projetos/posts | Média |
| Golang | 2 projetos/posts | Média |

**Critérios para seleção:**
- `isPublic = true`
- Ordenado por `priority` (maior primeiro)
- Limite por área conforme tabela acima
- Mistura de tipos: projetos, posts, tech writings

## Implementação Técnica

### Frontend (SvelteKit)

#### 1. Hook/Store para Gerenciar Filtro

**Arquivo**: `src/lib/stores/refFilter.ts`

```typescript
// Estrutura conceitual
interface RefFilter {
  code: string | null;
  expiresAt: number | null;
  isValid: boolean;
}

// Funções principais:
- getRefFromURL(): string | null
- saveRefFilter(code: string): void
- getRefFilter(): RefFilter
- clearRefFilter(): void
- isFilterExpired(): boolean
```

#### 2. Detecção na Página Principal

**Arquivo**: `src/routes/+page.svelte` ou `+layout.svelte`

- No `onMount` ou `beforeMount`, verificar:
  1. Parâmetro `ref` na URL (`page.url.searchParams.get('ref')`)
  2. Se existe, validar código e salvar
  3. Se não existe, verificar localStorage por filtro ativo
  4. Se expirado, limpar e mostrar conteúdo default

#### 3. Componentes de Filtro

Cada seção que exibe conteúdo deve:
- Receber o filtro ativo como prop ou via store
- Filtrar dados antes de renderizar
- Mostrar indicador visual quando filtro está ativo (opcional)

**Componentes afetados:**
- `ProjectsShowcase.svelte`
- `BlogPostsSection.svelte`
- `TechnicalWritingsSection.svelte`
- `CaseStudiesSection.svelte`
- `SystemDesignsSection.svelte`
- `SocialMediaPostsSection.svelte`

#### 4. API Client Modificado

**Arquivo**: `src/lib/api/` (arquivos de API existentes)

- Adicionar parâmetro `refFilter` nas chamadas de API
- Backend filtra no servidor (mais eficiente)
- Ou filtrar no frontend após receber dados

### Backend (Go)

#### 1. Endpoints Modificados

Todos os endpoints que retornam conteúdo devem aceitar query parameter `ref`:

```
GET /api/projects?ref=do
GET /api/posts?ref=ds
GET /api/tech-writings?ref=go
```

#### 2. Lógica de Filtragem

**Arquivo**: `app/internal/handlers/` (handlers relevantes)

```go
// Pseudocódigo conceitual
func filterContent(items []ContentItem, refCode string) []ContentItem {
    if refCode == "" {
        // Retornar conteúdo público prioritário
        return getPublicPriorityContent(items)
    }
    
    // Filtrar por tag
    filtered := []ContentItem{}
    for _, item := range items {
        if contains(item.Tags, refCode) || item.IsPublic {
            filtered = append(filtered, item)
        }
    }
    
    return filtered
}
```

#### 3. Modelo de Dados

**Arquivo**: `app/internal/models/` (models relevantes)

Adicionar campos:
- `Tags []string` (JSON array no banco)
- `IsPublic bool`
- `Priority int`

#### 4. Migração de Banco de Dados

Criar migration para:
- Adicionar coluna `tags` (JSON/JSONB)
- Adicionar coluna `is_public` (BOOLEAN)
- Adicionar coluna `priority` (INTEGER)
- Criar índice em `tags` para busca eficiente (se suportado)

### Banco de Dados

#### Estrutura de Tabelas

```sql
-- Exemplo para tabela de projetos
ALTER TABLE projects 
ADD COLUMN tags JSONB DEFAULT '[]',
ADD COLUMN is_public BOOLEAN DEFAULT false,
ADD COLUMN priority INTEGER DEFAULT 0;

-- Índice para busca por tags (PostgreSQL)
CREATE INDEX idx_projects_tags ON projects USING GIN (tags);

-- Exemplo de dados
UPDATE projects 
SET tags = '["ds", "py"]', is_public = true, priority = 8
WHERE id = 'project-123';
```

## Experiência do Usuário

### 1. Indicadores Visuais

**Quando filtro está ativo:**
- Badge discreto no header: "Visualizando: DevOps"
- Opção para limpar filtro: "Ver tudo"
- Contador de itens filtrados (opcional)

**Quando sem filtro:**
- Exibir normalmente, sem indicadores

### 2. Transições

- Animações suaves ao aplicar/remover filtro
- Loading state durante filtragem
- Mensagem quando nenhum conteúdo encontrado

### 3. SEO e Compartilhamento

- Links com `?ref=xxx` são indexáveis
- Meta tags dinâmicas baseadas no filtro
- Open Graph tags ajustados por filtro

## Casos de Uso

### Caso 1: Recrutador DevOps
1. Acessa `www.woragis.me/?ref=do`
2. Vê apenas projetos/posts de DevOps
3. Filtro persiste por 24h em navegação interna
4. Após 24h, volta ao conteúdo default

### Caso 2: Visitante Sem Referência
1. Acessa `www.woragis.me`
2. Vê conteúdo público prioritário (mix de áreas)
3. Pode aplicar filtro manualmente (futuro)

### Caso 3: Múltiplas Tags
1. Projeto tem tags `["go", "aw"]`
2. Aparece tanto em `?ref=go` quanto `?ref=aw`
3. Não aparece em `?ref=ds`

### Caso 4: Conteúdo Público
1. Projeto tem `isPublic = true` e `tags = ["ds"]`
2. Aparece em `?ref=ds` (tag match)
3. Aparece sem filtro (isPublic)
4. Não aparece em `?ref=do`

## Validação e Segurança

### 1. Validação de Códigos

- Whitelist de códigos válidos
- Códigos inválidos são ignorados (fallback para default)
- Case-insensitive (DO = do)

### 2. Sanitização

- Validar formato do código (apenas letras minúsculas, 2 caracteres)
- Limpar parâmetros maliciosos
- Prevenir XSS em URLs

### 3. Performance

- Cache de conteúdo filtrado (se aplicável)
- Índices no banco para busca por tags
- Paginação mantida mesmo com filtro

## Roadmap de Implementação

### Fase 1: Fundação
- [ ] Criar store/hook para gerenciar filtro
- [ ] Implementar detecção de `ref` na URL
- [ ] Persistência em localStorage
- [ ] Expiração após 24h

### Fase 2: Backend
- [ ] Adicionar campos `tags`, `is_public`, `priority` nas tabelas
- [ ] Criar migration de banco de dados
- [ ] Implementar lógica de filtragem nos handlers
- [ ] Adicionar query parameter `ref` nos endpoints

### Fase 3: Frontend - Filtragem
- [ ] Modificar componentes para aceitar filtro
- [ ] Implementar filtragem de dados
- [ ] Adicionar indicadores visuais
- [ ] Testar todos os códigos de referência

### Fase 4: Conteúdo Default
- [ ] Marcar conteúdo prioritário no banco
- [ ] Implementar lógica de seleção default
- [ ] Testar distribuição de conteúdo

### Fase 5: Polimento
- [ ] Animações e transições
- [ ] Mensagens de estado vazio
- [ ] SEO e meta tags dinâmicas
- [ ] Testes E2E

### Fase 6: Admin/Content Management
- [ ] Interface para gerenciar tags
- [ ] Toggle `isPublic` e `priority`
- [ ] Preview de filtros

## Considerações Futuras

### 1. Filtro Manual
- Dropdown no header para aplicar filtro manualmente
- Útil para visitantes que chegam sem `ref`

### 2. Múltiplos Filtros
- `?ref=do,ds` para ver DevOps E Data Science
- Lógica de AND/OR

### 3. Analytics
- Rastrear quais filtros são mais usados
- Métricas de engajamento por área

### 4. Personalização Avançada
- Salvar preferências do usuário
- Recomendações baseadas em histórico

### 5. A/B Testing
- Testar diferentes distribuições de conteúdo default
- Otimizar prioridades

## Exemplos de Uso

### Exemplo 1: Link para Vaga DevOps
```
www.woragis.me/?ref=do
```

### Exemplo 2: Link para Vaga Data Science
```
www.woragis.me/?ref=ds
```

### Exemplo 3: Link para Vaga Full Stack (Next.js)
```
www.woragis.me/?ref=nx
```

### Exemplo 4: Sem Referência (Default)
```
www.woragis.me
```

## Notas de Implementação

### LocalStorage vs Cookies
- **Recomendação**: LocalStorage para simplicidade
- **Fallback**: Cookies se precisar de persistência entre domínios

### Expiração
- Calcular `expiresAt = Date.now() + (24 * 60 * 60 * 1000)`
- Verificar em cada carregamento de página
- Limpar automaticamente se expirado

### Compatibilidade
- Funciona mesmo se JavaScript estiver desabilitado (conteúdo default)
- Progressive enhancement: filtro é um "nice to have"

### Performance
- Filtragem no backend é mais eficiente (menos dados transferidos)
- Cache de resultados filtrados pode ser útil
- Considerar paginação mesmo com filtro ativo

## Conclusão

Este sistema permite criar uma experiência personalizada para cada visitante, demonstrando expertise focada em áreas específicas enquanto mantém flexibilidade para mostrar conteúdo diversificado quando apropriado. A implementação é modular e pode ser expandida conforme necessário.

