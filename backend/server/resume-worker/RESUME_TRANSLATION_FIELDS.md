# Resume Translation Fields

This document outlines which fields in the resume generation system are translatable.

## Translatable Fields by Entity Type

### Projects (`project`)
- **name**: Project name
- **description**: Project description

### Certifications (`certification`)
- **name**: Certification name
- **description**: Certification description
- **category**: Category display name (if stored as text)

### Experiences (`experience`)
- **position**: Job title/position
- **description**: Job description/bullet points

### Posts/Publications (`post`)
- **title**: Publication title
- **excerpt**: Publication excerpt
- **content**: Publication content

## Hardcoded Sections (Translated via Helper)

### Education Section
- **Status values**: "In Progress" → "Em Andamento", "Completed" → "Concluído", etc.
- **Section headers**: "Education" → "Educação", "Relevant Coursework" → "Disciplinas Relevantes"

### Work Experience Section
- **Section headers**: "Work Experience" → "Experiência Profissional"

### Volunteer Work Section
- **Section headers**: "Volunteer Work" → "Trabalho Voluntário"

### Other Sections
- **Profile** → "Perfil"
- **Projects** → "Projetos"
- **Experience** → "Experiência"
- **Certifications** → "Certificações"
- **Publications** → "Publicações"
- **Languages** → "Idiomas"
- **Status** → "Status"

## How It Works

1. **Source Entity Translations**: When generating a resume, the system checks for translations of projects, certifications, experiences, and posts in the target language.

2. **Translation Helper**: The `TranslationHelper` class:
   - Fetches translations from the `translations` table
   - Maps language codes (e.g., "pt" → "pt-BR")
   - Applies translations to entity data before rendering

3. **Section Headers**: Section headers are translated using a mapping function that supports multiple languages.

4. **Status Values**: Education status values are translated using a mapping function.

## Usage

The resume generator automatically uses translations when:
- A `TranslationHelper` is provided
- Translations exist in the database with `status = 'completed'`
- The target language matches the translation language

## Adding New Languages

To add support for a new language:
1. Add the language constant to `translations/entity.go`
2. Update the `translate_section_header` method in `translation_helper.py`
3. Update the `translate_education_status` method in `translation_helper.py`

