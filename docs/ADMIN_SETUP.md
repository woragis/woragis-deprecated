# Admin Management System Setup

## 🚀 Features

- **Project Management**: Create, edit, delete, and reorder projects
- **Drag & Drop**: Reorder projects with intuitive drag and drop
- **Visibility Control**: Toggle project visibility on/off
- **Featured Projects**: Mark projects as featured for home page display
- **Settings Management**: Configure how many projects to show on home page
- **Dedicated Projects Page**: Full projects showcase with pagination
- **Database Integration**: SQLite with Drizzle ORM

## 📦 Database Setup

### 1. Install Dependencies

```bash
npm install
```

### 2. Generate Database Migration

```bash
npm run db:generate
```

### 3. Run Migration

```bash
npm run db:migrate
```

### 4. Seed Sample Data

```bash
npm run db:seed
```

### 5. Open Database Studio (Optional)

```bash
npm run db:studio
```

## 🗄️ Database Schema

### Projects Table

- `id`: Primary key (auto-increment)
- `title`: Project title
- `description`: Short description
- `longDescription`: Detailed description
- `technologies`: JSON array of technologies
- `image`: Image URL or emoji
- `githubUrl`: GitHub repository URL
- `liveUrl`: Live demo URL
- `featured`: Boolean (shown on home page)
- `order`: Display order (for drag & drop)
- `visible`: Boolean (public visibility)
- `createdAt`: Creation timestamp
- `updatedAt`: Last update timestamp

### Settings Table

- `id`: Primary key
- `key`: Setting key (e.g., "projects_per_page")
- `value`: Setting value
- `updatedAt`: Last update timestamp

## 🎯 Admin Features

### Project Management

- **Create**: Add new projects with full details
- **Edit**: Update existing project information
- **Delete**: Remove projects permanently
- **Reorder**: Drag and drop to change display order
- **Toggle Visibility**: Show/hide projects from public view
- **Toggle Featured**: Mark projects for home page display

### Settings

- **Projects Per Page**: Configure how many projects show on home page
- **Easy Configuration**: Simple dropdown interface

## 🔗 API Endpoints

### Admin API

- `GET /api/admin/projects` - Get all projects
- `POST /api/admin/projects` - Create new project
- `PUT /api/admin/projects` - Update project or reorder
- `DELETE /api/admin/projects?id={id}` - Delete project
- `GET /api/admin/settings` - Get settings
- `POST /api/admin/settings` - Update settings

### Public API

- `GET /api/projects` - Get visible projects
- `GET /api/projects?featured=true&limit=3` - Get featured projects

## 📱 Pages

### Admin Dashboard

- **URL**: `/admin`
- **Features**: Full project management interface
- **Access**: Currently open (add authentication as needed)

### Projects Page

- **URL**: `/projects`
- **Features**: Public projects showcase with pagination
- **Design**: Responsive grid layout

### Home Page

- **URL**: `/`
- **Features**: Shows featured projects (configurable limit)
- **Integration**: Links to full projects page

## 🎨 UI Components

### Admin Components

- `AdminDashboard`: Main admin interface
- `ProjectForm`: Create/edit project form
- `ProjectList`: Drag & drop project list

### Public Components

- `Projects`: Home page featured projects
- `ProjectsPage`: Full projects showcase

## 🔧 Configuration

### Environment Variables

Create a `.env.local` file:

```env
# For local development
DATABASE_URL="file:./local.db"
DATABASE_AUTH_TOKEN=""

# For production with Turso
# DATABASE_URL="libsql://your-database-name.turso.io"
# DATABASE_AUTH_TOKEN="your-auth-token"
```

### Database Options

#### Local SQLite (Development)

- Uses local file: `./local.db`
- No external dependencies
- Perfect for development

#### Turso (Production)

- Cloud-hosted SQLite
- Global edge locations
- Built-in backups
- Free tier available

## 🚀 Deployment

### 1. Set up Database

Choose your database option and configure environment variables.

### 2. Run Migrations

```bash
npm run db:migrate
```

### 3. Seed Data (Optional)

```bash
npm run db:seed
```

### 4. Deploy

Deploy to your preferred platform (Vercel, Netlify, etc.)

## 🔒 Security Considerations

### Current State

- Admin panel is currently open access
- No authentication implemented

### Recommended Security

1. Add authentication (NextAuth.js, Clerk, etc.)
2. Implement role-based access control
3. Add API rate limiting
4. Validate all inputs
5. Add CSRF protection

## 📊 Usage Examples

### Adding a New Project

1. Go to `/admin`
2. Click "Add Project"
3. Fill in project details
4. Set visibility and featured status
5. Save project

### Reordering Projects

1. Go to `/admin`
2. Drag projects by the grip handle
3. Drop in desired position
4. Order is automatically saved

### Managing Visibility

1. Go to `/admin`
2. Click the eye icon to toggle visibility
3. Hidden projects won't appear on public pages

### Configuring Home Page

1. Go to `/admin`
2. Adjust "Projects per page" setting
3. Only featured projects will show on home page

## 🎯 Next Steps

1. **Add Authentication**: Implement secure admin access
2. **Image Upload**: Add file upload for project images
3. **Categories**: Add project categorization
4. **Analytics**: Track project views and interactions
5. **SEO**: Add dynamic meta tags for projects
6. **Caching**: Implement Redis for better performance

## 🐛 Troubleshooting

### Database Connection Issues

- Check environment variables
- Ensure database file exists (for SQLite)
- Verify Turso credentials (for production)

### Migration Errors

- Delete existing database file
- Run `npm run db:generate` again
- Run `npm run db:migrate`

### Build Errors

- Check TypeScript types
- Verify all imports are correct
- Run `npm run lint` to check for issues
