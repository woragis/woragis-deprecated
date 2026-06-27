#!/usr/bin/env python3
"""
Create Animated GIFs for Technical Diagrams

This script creates animated GIFs showing workflow progression or state changes.
You can either:
1. Generate multiple images with slight variations and animate them
2. Create frames programmatically and combine into GIF
"""

import os
import argparse
from pathlib import Path
from typing import List

try:
    from PIL import Image, ImageDraw, ImageFont
    from generate_post_images import PostImageGenerator, create_gif_from_images
except ImportError:
    print("Error: Required packages not installed. Run: pip install -r requirements-images.txt")
    exit(1)


def create_workflow_gif(post_path: str, steps: List[str], output_path: str, provider: str = "openai"):
    """
    Create animated GIF showing workflow progression.
    
    Args:
        post_path: Path to markdown post
        steps: List of step descriptions for each frame
        output_path: Output GIF path
        provider: AI provider (openai or replicate)
    """
    generator = PostImageGenerator(provider=provider)
    
    # Read post to get context
    with open(post_path, 'r', encoding='utf-8') as f:
        content = f.read()
    
    title = generator._extract_title(content)
    topic = generator._extract_topic(content, Path(post_path))
    
    # Generate images for each step
    image_paths = []
    temp_dir = Path(post_path).parent / "images" / "temp"
    temp_dir.mkdir(parents=True, exist_ok=True)
    
    for i, step in enumerate(steps):
        # Build prompt with step information
        prompt = f"Technical workflow diagram, step {i+1} of {len(steps)}: {step}, {topic}, arrows showing progression, numbered sequence, clean white background, modern flat design"
        
        # Update generator config
        generator.generator.config.width = 1920
        generator.generator.config.height = 1080
        
        # Generate image
        step_image_path = temp_dir / f"step-{i+1}.png"
        success = generator.generator.generate(prompt, str(step_image_path))
        
        if success:
            image_paths.append(str(step_image_path))
        else:
            print(f"Warning: Failed to generate image for step {i+1}")
    
    if len(image_paths) < 2:
        print("Error: Need at least 2 images to create GIF")
        return False
    
    # Create GIF
    success = create_gif_from_images(image_paths, output_path, duration=1.0)
    
    # Cleanup temp files (optional)
    # for path in image_paths:
    #     os.remove(path)
    
    return success


def create_simple_animated_diagram(output_path: str, frames: int = 5):
    """
    Create a simple animated diagram programmatically.
    Example: Animated microservices connecting.
    """
    if Image is None:
        print("Error: Pillow required")
        return False
    
    width, height = 1920, 1080
    images = []
    
    for frame in range(frames):
        # Create new image
        img = Image.new('RGB', (width, height), color='white')
        draw = ImageDraw.Draw(img)
        
        # Draw boxes (microservices)
        box_size = 150
        spacing = 300
        start_x = 200
        start_y = height // 2 - box_size // 2
        
        for i in range(4):
            x = start_x + i * spacing
            y = start_y
            
            # Animate: boxes appear one by one
            if i <= frame:
                # Draw box
                draw.rectangle([x, y, x + box_size, y + box_size], 
                             outline='blue', width=5, fill='lightblue')
                
                # Draw label
                draw.text((x + box_size//2 - 30, y + box_size//2 - 10), 
                         f"Service {i+1}", fill='black')
                
                # Draw arrow to next box (if not last)
                if i < 3 and i < frame:
                    arrow_x = x + box_size
                    arrow_y = y + box_size // 2
                    draw.line([arrow_x, arrow_y, arrow_x + spacing - box_size, arrow_y], 
                            fill='blue', width=3)
                    # Arrowhead
                    draw.polygon([
                        (arrow_x + spacing - box_size - 10, arrow_y - 5),
                        (arrow_x + spacing - box_size, arrow_y),
                        (arrow_x + spacing - box_size - 10, arrow_y + 5)
                    ], fill='blue')
        
        images.append(img)
    
    # Save as GIF
    images[0].save(
        output_path,
        save_all=True,
        append_images=images[1:],
        duration=800,  # 0.8 seconds per frame
        loop=0,
    )
    
    print(f"✓ Animated diagram saved to: {output_path}")
    return True


def main():
    parser = argparse.ArgumentParser(description="Create animated GIFs for technical diagrams")
    parser.add_argument("--post", type=str, help="Path to markdown post")
    parser.add_argument("--steps", type=str, nargs="+", help="List of workflow steps")
    parser.add_argument("--output", type=str, required=True, help="Output GIF path")
    parser.add_argument("--provider", type=str, choices=["openai", "replicate"], default="openai")
    parser.add_argument("--simple", action="store_true", help="Create simple programmatic animation")
    parser.add_argument("--frames", type=int, default=5, help="Number of frames for simple animation")
    
    args = parser.parse_args()
    
    if args.simple:
        # Create simple programmatic animation
        create_simple_animated_diagram(args.output, args.frames)
    elif args.post and args.steps:
        # Create workflow GIF from post
        create_workflow_gif(args.post, args.steps, args.output, args.provider)
    else:
        parser.print_help()
        print("\nExample usage:")
        print("  # Simple animation:")
        print("  python create_animated_diagram.py --simple --output images/microservices.gif")
        print("\n  # Workflow animation:")
        print("  python create_animated_diagram.py --post posts/cicd/cicd-strategy.md \\")
        print("    --steps 'Build' 'Test' 'Deploy' 'Monitor' --output images/cicd-workflow.gif")


if __name__ == "__main__":
    main()

