import 'dotenv/config';
import { Database } from './database.js';
import { CoverLetterService } from './coverLetter.js';
import { logger } from './utils/logger.js';

/**
 * Test cover letter generation with user profile data
 */
async function testCoverLetter() {
  console.log('📝 Testing Cover Letter Generation\n');

  const db = new Database();
  const coverLetterService = new CoverLetterService();

  try {
    await db.connect();
    console.log('✅ Database connected\n');

    // Get user ID from environment or use a test one
    // You'll need to find your actual user ID from the database
    const userId = process.env.TEST_USER_ID;
    
    if (!userId) {
      console.log('⚠️  TEST_USER_ID not set. Please set it in .env file');
      console.log('   You can find your user ID by querying: SELECT id FROM users WHERE email = \'your-email@example.com\';\n');
      process.exit(1);
    }

    console.log(`👤 Fetching profile for user: ${userId}\n`);

    // Fetch user profile
    const profile = await db.fetchUserProfile(userId);

    console.log('📋 Profile Data:');
    console.log(`  Projects: ${profile.projects.length}`);
    console.log(`  Posts: ${profile.posts.length}`);
    console.log(`  Technical Writings: ${profile.technicalWritings.length}`);
    console.log(`  Case Studies: ${profile.caseStudies.length}`);
    console.log(`  Certifications: ${profile.certifications.length}`);
    console.log(`  Skills: ${profile.skills.length}`);
    console.log(`  Interests: ${profile.interests.length}\n`);

    if (profile.projects.length === 0 && profile.posts.length === 0) {
      console.log('⚠️  No profile data found. The cover letter will be generic.\n');
    }

    // Test job information
    const jobInfo = {
      companyName: 'Google',
      jobTitle: 'Senior Software Engineer',
      jobDescription: `We are looking for a Senior Software Engineer to join our team. 
      
Requirements:
- 5+ years of experience in software development
- Strong proficiency in Go, Python, or Java
- Experience with distributed systems and microservices
- Knowledge of Kubernetes, Docker, and cloud platforms (AWS, GCP)
- Experience with database design and optimization
- Strong problem-solving and communication skills

Responsibilities:
- Design and develop scalable backend systems
- Work with cross-functional teams
- Mentor junior engineers
- Participate in code reviews and architecture decisions`,
      location: 'Remote, United States',
      requirements: ['Go', 'Python', 'Kubernetes', 'Docker', 'Microservices', 'Distributed Systems', 'AWS', 'GCP'],
    };

    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━');
    console.log('TEST 1: Generate Cover Letter in English');
    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n');

    console.log('Job Information:');
    console.log(`  Company: ${jobInfo.companyName}`);
    console.log(`  Position: ${jobInfo.jobTitle}`);
    console.log(`  Location: ${jobInfo.location}\n`);

    console.log('🤖 Generating cover letter with AI...\n');
    const coverLetter = await coverLetterService.generateCoverLetter(profile, jobInfo);

    console.log('📝 Generated Cover Letter:');
    console.log('═'.repeat(80));
    console.log(coverLetter);
    console.log('═'.repeat(80));
    console.log(`\n📊 Statistics:`);
    console.log(`  Length: ${coverLetter.length} characters`);
    console.log(`  Word count: ~${coverLetter.split(/\s+/).length} words`);
    console.log(`  Paragraphs: ${coverLetter.split(/\n\n/).length}\n`);

    // Check if cover letter mentions profile data
    console.log('🔍 Checking if cover letter uses profile data...\n');
    
    let mentionsProfile = false;
    if (profile.projects.length > 0) {
      const projectMentioned = profile.projects.some(p => 
        coverLetter.toLowerCase().includes(p.name.toLowerCase())
      );
      if (projectMentioned) {
        console.log('  ✅ Mentions projects from profile');
        mentionsProfile = true;
      }
    }

    if (profile.skills.length > 0) {
      const skillMentioned = profile.skills.some(s => 
        coverLetter.toLowerCase().includes(s.name.toLowerCase())
      );
      if (skillMentioned) {
        console.log('  ✅ Mentions skills from profile');
        mentionsProfile = true;
      }
    }

    if (profile.technicalWritings.length > 0) {
      const writingMentioned = profile.technicalWritings.some(w => 
        coverLetter.toLowerCase().includes(w.title.toLowerCase())
      );
      if (writingMentioned) {
        console.log('  ✅ Mentions technical writings');
        mentionsProfile = true;
      }
    }

    if (!mentionsProfile && (profile.projects.length > 0 || profile.skills.length > 0)) {
      console.log('  ⚠️  Cover letter may not be using profile data effectively');
    }

    console.log('\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━');
    console.log('TEST 2: Generate Cover Letter in Portuguese');
    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n');

    // Generate Portuguese version
    const portuguesePrompt = coverLetterService.buildPrompt(profile, jobInfo);
    const portuguesePromptWithLanguage = portuguesePrompt.replace(
      'Write the cover letter now:',
      'IMPORTANT: Write the cover letter in Portuguese (Brazilian Portuguese - pt-BR). Use professional Brazilian business language. Write the cover letter now:'
    );

    console.log('🤖 Generating Portuguese cover letter...\n');

    const axios = (await import('axios')).default;
    const response = await axios.post(
      `${process.env.AI_SERVICE_URL || 'http://ai-service:8000'}/api/chat/completions`,
      {
        provider: 'openai',
        model: 'gpt-4o-mini',
        temperature: 0.7,
        messages: [
          {
            role: 'user',
            content: portuguesePromptWithLanguage,
          },
        ],
        max_tokens: 1500,
      },
      { timeout: 30000 }
    );

    const portugueseCoverLetter = response.data.message?.content || response.data.choices?.[0]?.message?.content;

    console.log('📝 Cover Letter em Português:');
    console.log('═'.repeat(80));
    console.log(portugueseCoverLetter.trim());
    console.log('═'.repeat(80));
    console.log(`\n📊 Statistics:`);
    console.log(`  Length: ${portugueseCoverLetter.trim().length} characters`);
    console.log(`  Word count: ~${portugueseCoverLetter.trim().split(/\s+/).length} words\n`);

    // Detect language
    const hasPortugueseWords = /(experiência|oportunidade|empresa|desenvolvimento|tecnologia)/i.test(portugueseCoverLetter);
    if (hasPortugueseWords) {
      console.log('✅ Cover letter appears to be in Portuguese\n');
    } else {
      console.log('⚠️  Cover letter may not be in Portuguese\n');
    }

    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━');
    console.log('✅ Cover Letter Tests Completed!');
    console.log('━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n');

  } catch (error) {
    console.error('❌ Test failed:', error.message);
    console.error('Stack:', error.stack);
    process.exit(1);
  } finally {
    await db.disconnect();
  }
}

testCoverLetter().catch(console.error);

