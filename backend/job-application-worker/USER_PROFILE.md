# User Profile Data for AI

## Overview

The job application worker fetches comprehensive user profile data from the database to provide context to the AI when generating cover letters. This ensures personalized, relevant cover letters that highlight the applicant's actual experience and skills.

## Data Sources

The AI receives information from these database tables:

### 1. **Projects** (`projects` table)
- Project name, description, status
- Health score, MRR, CAC, LTV, churn rate
- **Technologies** (from `project_technologies` table)
  - Technology name, version, category, purpose

### 2. **Posts** (`posts` table)
- Title, content, excerpt
- Published date, featured status
- **Skills** (from `post_skills` join table)
  - Skills associated with each post

### 3. **Technical Writings** (`technical_writings` table)
- Title, description, content
- Type (article, tutorial, guide, etc.)
- Platform (Medium, Dev.to, personal blog, etc.)
- URL, published date

### 4. **Case Studies** (`case_studies` table)
- Title, problem, context, solution
- Approach (JSON array)
- Technologies (JSON array)
- Lessons learned (JSON array)

### 5. **Project Case Studies** (`project_case_studies` table)
- Linked to specific projects
- Problem, solution, approach
- Technologies and lessons learned

### 6. **Certifications** (`certifications` table)
- Name, issuer, issue date, expiry date
- Description, category
- Status (active, expired, etc.)

### 7. **Skills** (`skills` table)
- Fetched from projects and posts (via joins)
- Name, category, description
- Grouped by category (backend, frontend, database, etc.)

### 8. **Interests** (`interests` table)
- Title, description
- Featured interests only

## How It Works

### Fetching Process

```javascript
// In database.js
async fetchUserProfile(userId) {
  // 1. Fetch projects with their technologies
  // 2. Fetch posts with associated skills
  // 3. Fetch technical writings
  // 4. Fetch case studies
  // 5. Fetch project case studies
  // 6. Fetch certifications
  // 7. Fetch skills (from projects and posts)
  // 8. Fetch interests
  
  return profile; // Comprehensive profile object
}
```

### Profile Structure

```javascript
{
  projects: [
    {
      name: "E-commerce Platform",
      description: "Full-stack e-commerce solution",
      status: "completed",
      techStack: [
        { name: "React", version: "18.0", category: "frontend" },
        { name: "Node.js", version: "20.0", category: "backend" }
      ],
      metrics: { mrr: 5000, cac: 50, ltv: 1000 }
    }
  ],
  posts: [
    {
      title: "Building Scalable APIs",
      content: "...",
      skills: ["REST API", "Node.js", "Microservices"]
    }
  ],
  technicalWritings: [
    {
      title: "Understanding GraphQL",
      platform: "medium",
      url: "https://...",
      content: "..."
    }
  ],
  caseStudies: [
    {
      title: "Reducing API Response Time",
      problem: "...",
      solution: "...",
      technologies: ["Redis", "Node.js"],
      lessonsLearned: ["Caching is crucial", "..."]
    }
  ],
  certifications: [
    {
      name: "AWS Certified Solutions Architect",
      issuer: "Amazon Web Services",
      issueDate: "2023-01-15"
    }
  ],
  skills: [
    { name: "React", category: "frontend" },
    { name: "Node.js", category: "backend" }
  ],
  interests: [
    { title: "Distributed Systems", description: "..." }
  ]
}
```

## AI Prompt Construction

The profile data is structured into a comprehensive prompt:

```
Job Information:
- Company: Google
- Position: Senior Software Engineer
- Location: Remote

Applicant Profile:

## Projects & Experience:
- **Project Name**: Description
  Technologies: React, Node.js, PostgreSQL

## Case Studies:
- **Case Study Title**: Problem and solution
  Technologies: Redis, Node.js
  Key Learnings: Caching strategies, performance optimization

## Technical Writings:
- **Article Title** (Medium): Description

## Skills:
- **backend**: Node.js, Python, Go
- **frontend**: React, TypeScript

## Certifications:
- **AWS Certified** from Amazon (2023)

## Interests:
- **Distributed Systems**: Focus on scalability
```

## Benefits

1. **Personalized**: Uses actual user data, not generic templates
2. **Relevant**: Highlights technologies and experience matching job requirements
3. **Comprehensive**: Includes projects, writings, certifications, case studies
4. **Contextual**: AI understands the full scope of user's experience
5. **Dynamic**: Always uses latest data from database

## Data Limits

To keep prompts manageable and within token limits:

- **Projects**: Top 5 most recent
- **Posts**: Top 5 published posts
- **Technical Writings**: Top 5 most recent
- **Case Studies**: Top 3 featured
- **Project Case Studies**: Top 3
- **Skills**: Up to 50 (grouped by category)
- **Content Truncation**: Long content is limited to 1000 characters

## Example AI Output

With this profile data, the AI can generate cover letters like:

> "Dear Hiring Manager,
> 
> I am excited to apply for the Senior Software Engineer position at Google. 
> 
> In my recent work on an e-commerce platform using React and Node.js, I reduced API response times by 60% through strategic Redis caching—a case study I documented on Medium. My experience with distributed systems, demonstrated through my AWS certification and technical writings on scalability, aligns perfectly with Google's infrastructure challenges.
> 
> I've consistently worked with technologies like React, Node.js, and PostgreSQL across multiple projects, and I'm particularly interested in distributed systems, which matches Google's focus areas.
> 
> I would welcome the opportunity to discuss how my experience with scalable architectures and performance optimization can contribute to your team.
> 
> Sincerely,
> [Name]"

## Future Enhancements

- **Job Description Analysis**: Extract keywords from job description and prioritize matching profile items
- **Relevance Scoring**: Rank profile items by relevance to job requirements
- **Dynamic Truncation**: Adjust content limits based on token budget
- **Caching**: Cache profile data in Redis to reduce database queries

