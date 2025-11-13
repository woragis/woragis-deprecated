# Translation System Implementation

This document describes the comprehensive translation system implemented for the portfolio project, supporting 8 languages with database-driven translations and fallback mechanisms.

## 🌍 Supported Languages

- **English (en)** - Default language
- **Spanish (es)** - Español
- **Portuguese (pt)** - Português
- **Italian (it)** - Italiano
- **French (fr)** - Français
- **Japanese (ja)** - 日本語
- **Chinese (zh)** - 中文
- **Korean (ko)** - 한국어

## 🏗️ Architecture Overview

The translation system uses a **hybrid approach** combining:
- **Static translations** for UI text (using next-intl)
- **Database translations** for dynamic content (blog posts, projects, etc.)
- **Fallback mechanisms** to ensure content is always available

## 📊 Database Schema

### Translation Tables

Each content type has its own translation table:

```sql
-- Blog Post Translations
blog_post_translations (
  blog_post_id UUID REFERENCES blog_posts(id),
  language_code TEXT CHECK (language_code IN ('en', 'es', 'pt', 'it', 'fr', 'ja', 'zh', 'ko')),
  title TEXT NOT NULL,
  excerpt TEXT NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  PRIMARY KEY (blog_post_id, language_code)
);

-- Project Translations
project_translations (
  project_id UUID REFERENCES projects(id),
  language_code TEXT CHECK (language_code IN ('en', 'es', 'pt', 'it', 'fr', 'ja', 'zh', 'ko')),
  title TEXT NOT NULL,
  description TEXT NOT NULL,
  long_description TEXT,
  content TEXT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  PRIMARY KEY (project_id, language_code)
);

-- Experience Translations
experience_translations (
  experience_id UUID REFERENCES experiences(id),
  language_code TEXT CHECK (language_code IN ('en', 'es', 'pt', 'it', 'fr', 'ja', 'zh', 'ko')),
  title TEXT NOT NULL,
  company TEXT NOT NULL,
  period TEXT NOT NULL,
  location TEXT NOT NULL,
  description TEXT NOT NULL,
  achievements TEXT NOT NULL, -- JSON string
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  PRIMARY KEY (experience_id, language_code)
);

-- Education Translations
education_translations (
  education_id UUID REFERENCES educations(id),
  language_code TEXT CHECK (language_code IN ('en', 'es', 'pt', 'it', 'fr', 'ja', 'zh', 'ko')),
  degree TEXT NOT NULL,
  school TEXT NOT NULL,
  year TEXT NOT NULL,
  description TEXT,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  PRIMARY KEY (education_id, language_code)
);

-- Testimonial Translations
testimonial_translations (
  testimonial_id UUID REFERENCES testimonials(id),
  language_code TEXT CHECK (language_code IN ('en', 'es', 'pt', 'it', 'fr', 'ja', 'zh', 'ko')),
  name TEXT NOT NULL,
  role TEXT NOT NULL,
  company TEXT NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  PRIMARY KEY (testimonial_id, language_code)
);

-- About Translations
about_translations (
  about_id UUID REFERENCES about_sections(id),
  language_code TEXT CHECK (language_code IN ('en', 'es', 'pt', 'it', 'fr', 'ja', 'zh', 'ko')),
  title TEXT NOT NULL,
  content TEXT NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  PRIMARY KEY (about_id, language_code)
);
```

## 🚀 Setup and Migration

### 1. Create Translation Tables

```bash
npm run db:migrate-translations
```

### 2. Migrate Existing Content

```bash
npm run db:migrate-content
```

This will create English translations for all existing content.

## 🔧 API Endpoints

### Translation Management

```typescript
// Get translation statistics
GET /api/translations

// Get translation status for specific content
GET /api/translations?contentType=blog&contentId=uuid

// Get content with translation
GET /api/translations/blog/[id]?lang=es&fallback=true

// Create translation
POST /api/translations/blog/[id]
{
  "languageCode": "es",
  "title": "Título en español",
  "excerpt": "Resumen en español",
  "content": "Contenido en español"
}

// Update translation
PUT /api/translations/blog/[id]
{
  "languageCode": "es",
  "title": "Título actualizado",
  "excerpt": "Resumen actualizado",
  "content": "Contenido actualizado"
}
```

### Content Endpoints with Language Support

All content endpoints now support a `lang` parameter:

```typescript
// Get blog posts in Spanish
GET /api/blog?lang=es

// Get projects in French
GET /api/projects?lang=fr

// Get experiences in Japanese
GET /api/experience?lang=ja
```

## 🎯 Frontend Integration

### Language Context

The enhanced `LanguageContext` provides:

```typescript
const { 
  language, 
  setLanguage, 
  t, 
  getTranslatedContent,
  isTranslationLoading 
} = useLanguage();

// Get translated content
const translatedBlog = await getTranslatedContent<BlogPost>(
  "blog", 
  blogId, 
  true // fallback to default
);
```

### Custom Hooks

```typescript
// Single translation
const { data, isLoading, error } = useTranslation<BlogPost>(
  "blog", 
  blogId, 
  { fallbackToDefault: true }
);

// Multiple translations
const { data, isLoading, error } = useMultipleTranslations<BlogPost>(
  "blog", 
  [blogId1, blogId2], 
  { fallbackToDefault: true }
);

// Translation statistics
const { stats, isLoading } = useTranslationStats();

// Translation status
const { status, isLoading } = useTranslationStatus("blog", blogId);
```

## 🎨 Admin Interface

### Translation Manager Component

The `TranslationManager` component provides a comprehensive interface for managing translations:

```typescript
<TranslationManager
  contentType="blog"
  contentId={blogId}
  contentTitle={blogTitle}
/>
```

Features:
- **Language selector** with completion status
- **Form-based editing** for each content type
- **Real-time validation**
- **Bulk operations**
- **Translation statistics**

## 📈 Translation Workflow

### 1. Content Creation
1. Create content in primary language (English)
2. Content is automatically available in English

### 2. Translation Process
1. Access admin interface
2. Select content to translate
3. Choose target language
4. Fill in translation form
5. Save translation

### 3. Content Display
1. User selects language preference
2. System fetches translated content
3. Falls back to English if translation missing
4. Displays content in user's preferred language

## 🔄 Fallback Strategy

The system implements a robust fallback strategy:

1. **Requested Language**: Try to fetch content in user's selected language
2. **Default Language**: If not available, fall back to English
3. **Original Content**: If no translations exist, show original content
4. **Error Handling**: Graceful degradation with error messages

## 📊 Translation Statistics

The system tracks:

- **Completion percentage** per language
- **Missing translations** by content type
- **Last updated** timestamps
- **Translation quality** metrics

## 🛠️ Development Tools

### Migration Scripts

```bash
# Create translation tables
npm run db:migrate-translations

# Migrate existing content
npm run db:migrate-content
```

### Type Safety

All translation operations are fully typed:

```typescript
type SupportedLanguage = "en" | "es" | "pt" | "it" | "fr" | "ja" | "zh" | "ko";

interface TranslationResponse<T> {
  content: T;
  language: SupportedLanguage;
  isTranslated: boolean;
  isFallback: boolean;
}
```

## 🚀 Future Enhancements

### Planned Features

1. **Auto-translation** using Google Translate API
2. **Translation memory** for common phrases
3. **Bulk import/export** functionality
4. **Translation workflow** (draft → review → published)
5. **SEO optimization** for different languages
6. **A/B testing** for translations
7. **Translation analytics** and insights

### Integration Points

- **Content Management Systems**
- **Translation Services** (Google Translate, DeepL)
- **SEO Tools** for multilingual optimization
- **Analytics** for translation performance

## 🔒 Security Considerations

- **Input validation** for all translation data
- **SQL injection protection** through parameterized queries
- **XSS prevention** with proper sanitization
- **Access control** for translation management
- **Rate limiting** for translation API endpoints

## 📝 Best Practices

### Content Creation
- Always create content in English first
- Use clear, translatable language
- Avoid cultural references that don't translate well
- Keep sentences simple and direct

### Translation Management
- Review translations before publishing
- Maintain consistency in terminology
- Use translation memory for common phrases
- Regular quality checks and updates

### Performance
- Cache translated content appropriately
- Use database indexes for language queries
- Implement lazy loading for non-critical translations
- Monitor translation API performance

## 🐛 Troubleshooting

### Common Issues

1. **Missing translations**: Check if content exists in requested language
2. **Fallback not working**: Verify English translation exists
3. **Performance issues**: Check database indexes and caching
4. **API errors**: Validate request parameters and authentication

### Debug Tools

```typescript
// Check translation status
const status = await translationService.getTranslationStatus("blog", blogId);

// Get translation statistics
const stats = await translationService.getTranslationStats();

// Validate translation data
const validation = translationService.validateTranslationData("blog", data);
```

## 📚 Resources

- [Next.js Internationalization](https://nextjs.org/docs/advanced-features/i18n)
- [next-intl Documentation](https://next-intl-docs.vercel.app/)
- [Drizzle ORM](https://orm.drizzle.team/)
- [Translation Best Practices](https://developers.google.com/style/translation)

---

This translation system provides a robust, scalable solution for multilingual content management while maintaining excellent performance and user experience.
