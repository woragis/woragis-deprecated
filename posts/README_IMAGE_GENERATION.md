# AI Image Generation for Technical Posts

This directory contains tools and guides for generating AI-powered images for technical blog posts.

## Quick Start

### 1. Install Dependencies

```bash
pip install -r requirements-images.txt
```

### 2. Set Up API Keys

**For OpenAI DALL-E 3:**
```bash
export OPENAI_API_KEY="your-api-key-here"
```

**For Replicate (Stable Diffusion):**
```bash
export REPLICATE_API_TOKEN="your-token-here"
```

### 3. Generate Images

**Single post thumbnail:**
```bash
python generate_post_images.py --post posts/cicd/cicd-strategy.md --type thumbnail
```

**Single post diagram:**
```bash
python generate_post_images.py --post posts/cicd/cicd-strategy.md --type diagram
```

**Batch generate thumbnails for all posts:**
```bash
python generate_post_images.py --batch --type thumbnail
```

**Use Replicate instead of OpenAI:**
```bash
python generate_post_images.py --post posts/cicd/cicd-strategy.md --type thumbnail --provider replicate
```

## Image Types

### Thumbnails
- **Size**: 1200x630px (LinkedIn/Twitter optimal)
- **Use case**: Social media previews, post headers
- **Style**: Eye-catching, modern, technical

### Diagrams
- **Size**: 1920x1080px (blog post width)
- **Use case**: Architecture diagrams, workflow visualizations
- **Style**: Technical, clean, professional

## Output Structure

Images are saved to:
```
posts/
  cicd/
    images/
      cicd-strategy-thumbnail.png
      cicd-strategy-thumbnail-prompt.txt
      cicd-strategy-diagram.png
      cicd-strategy-diagram-prompt.txt
```

## Cost Estimates

### OpenAI DALL-E 3
- $0.04 per thumbnail (1024x1024)
- $0.08 per diagram (1024x1792)
- **100 posts**: ~$4-8

### Replicate (Stable Diffusion)
- ~$0.002-0.01 per image
- **100 posts**: ~$0.20-1.00

## Creating Animated GIFs

To create animated GIFs showing workflow progression:

1. Generate multiple diagram images with slight variations
2. Use the `create_gif_from_images()` function in the script
3. Or use external tools like ImageMagick or online GIF makers

Example workflow:
```python
from generate_post_images import create_gif_from_images

image_paths = [
    "images/workflow-step1.png",
    "images/workflow-step2.png",
    "images/workflow-step3.png",
    "images/workflow-step4.png",
]

create_gif_from_images(image_paths, "images/workflow-animated.gif", duration=0.5)
```

## Customizing Prompts

Edit `PromptBuilder` class in `generate_post_images.py` to customize:
- Color schemes
- Style templates
- Technical concept mappings
- Diagram types

## Best Practices

1. **Generate multiple variations**: Run the script 2-3 times and pick the best
2. **Refine prompts**: Save prompts and adjust based on results
3. **Consistent style**: Use same prompt templates across posts
4. **Post-process**: Add text labels, adjust colors, optimize file size
5. **Version control**: Keep prompt files with images

## Troubleshooting

### "OPENAI_API_KEY not found"
- Set environment variable: `export OPENAI_API_KEY="your-key"`
- Or pass via command line (not recommended for security)

### "Package not installed"
- Install missing packages: `pip install -r requirements-images.txt`

### Images not generating
- Check API key is valid
- Check internet connection
- Check API rate limits
- Try different provider (replicate vs openai)

## Next Steps

1. Generate images for existing posts
2. Create custom prompt templates for your style
3. Set up automated generation in CI/CD
4. Create animated GIFs for workflow posts
5. Optimize images for web (compress, resize)

## Resources

- [OpenAI DALL-E 3 API Docs](https://platform.openai.com/docs/guides/images)
- [Replicate API Docs](https://replicate.com/docs)
- [Pillow Documentation](https://pillow.readthedocs.io/)
- See `AI_IMAGE_GENERATION_GUIDE.md` for detailed recommendations

