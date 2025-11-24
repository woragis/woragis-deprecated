<script lang="ts">
	import { onMount } from 'svelte';
	import { ChevronLeft, ChevronRight, Star, Linkedin, Play, ExternalLink, FolderKanban, Code2 } from 'lucide-svelte';
	import type { Testimonial, TestimonialType } from '$lib/types/testimonial';

	interface Props {
		testimonials: Testimonial[];
		loading?: boolean;
	}

	let { testimonials = [], loading = false }: Props = $props();

	let currentIndex = $state(0);
	let autoplayInterval: ReturnType<typeof setInterval> | null = null;

	onMount(() => {
		if (testimonials.length > 1) {
			startAutoplay();
		}
		return () => {
			stopAutoplay();
		};
	});

	function nextTestimonial() {
		if (testimonials.length === 0) return;
		currentIndex = (currentIndex + 1) % testimonials.length;
	}

	function prevTestimonial() {
		if (testimonials.length === 0) return;
		currentIndex = currentIndex === 0 ? testimonials.length - 1 : currentIndex - 1;
	}

	function goToTestimonial(index: number) {
		if (index >= 0 && index < testimonials.length) {
			currentIndex = index;
		}
	}

	function startAutoplay() {
		if (testimonials.length <= 1) return;
		autoplayInterval = setInterval(() => {
			nextTestimonial();
		}, 5000); // Change every 5 seconds
	}

	function stopAutoplay() {
		if (autoplayInterval) {
			clearInterval(autoplayInterval);
			autoplayInterval = null;
		}
	}

	function renderStars(rating?: number) {
		if (!rating) return [];
		return Array.from({ length: 5 }, (_, i) => i < rating);
	}

	function getTypeLabel(type: TestimonialType): string {
		const labels: Record<TestimonialType, string> = {
			general: 'General',
			project_specific: 'Project Specific',
			skill_specific: 'Skill Specific'
		};
		return labels[type] || type;
	}

	function getTypeColor(type: TestimonialType): string {
		const colors: Record<TestimonialType, string> = {
			general: 'bg-gray-600/20 text-gray-300 border-gray-500/30',
			project_specific: 'bg-blue-600/20 text-blue-300 border-blue-500/30',
			skill_specific: 'bg-purple-600/20 text-purple-300 border-purple-500/30'
		};
		return colors[type] || colors.general;
	}
</script>

<div class="w-full">
	{#if loading}
		<div class="flex items-center justify-center py-20">
			<div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-blue-500"></div>
		</div>
	{:else if testimonials.length === 0}
		<div class="text-center py-20">
			<p class="text-gray-400 text-lg mb-2">No testimonials available</p>
			<p class="text-gray-500 text-sm">Check back later</p>
		</div>
	{:else}
		<div
			class="relative bg-gradient-to-br from-gray-800/50 to-gray-900/50 backdrop-blur-sm rounded-xl p-8 border border-gray-700"
			onmouseenter={stopAutoplay}
			onmouseleave={startAutoplay}
		>
			<!-- Testimonial Content -->
			<div class="relative min-h-[200px]">
				{#each testimonials as testimonial, index}
					<div
						class="transition-all duration-500 ease-in-out {index === currentIndex
							? 'opacity-100 translate-x-0'
							: 'opacity-0 absolute inset-0 translate-x-8 pointer-events-none'}"
					>
						{#if index === currentIndex}
							<div class="flex flex-col items-center text-center">
								<!-- Author Photo -->
								{#if testimonial.authorPhoto}
									<img
										src={testimonial.authorPhoto}
										alt={testimonial.authorName}
										class="w-20 h-20 rounded-full object-cover mb-4 border-2 border-blue-500/50 shadow-lg"
									/>
								{:else}
									<div
										class="w-20 h-20 rounded-full bg-gradient-to-br from-blue-600 to-indigo-600 flex items-center justify-center text-white text-2xl font-bold mb-4 border-2 border-blue-500/50 shadow-lg"
									>
										{testimonial.authorName.charAt(0).toUpperCase()}
									</div>
								{/if}

								<!-- Rating -->
								{#if testimonial.rating}
									<div class="flex items-center gap-1 mb-4">
										{#each renderStars(testimonial.rating) as isFilled}
											<Star
												class="w-5 h-5 {isFilled
													? 'text-yellow-400 fill-yellow-400'
													: 'text-gray-600'}"
											/>
										{/each}
									</div>
								{/if}

								<!-- Type Badge -->
								{#if testimonial.type && testimonial.type !== 'general'}
									<div class="mb-4">
										<span
											class="px-3 py-1 text-xs rounded-full border {getTypeColor(testimonial.type)}"
										>
											{getTypeLabel(testimonial.type)}
										</span>
									</div>
								{/if}

								<!-- Content -->
								<blockquote class="text-gray-200 text-lg mb-4 max-w-3xl leading-relaxed">
									"{testimonial.content}"
								</blockquote>

								<!-- Context -->
								{#if testimonial.context}
									<div class="mb-4 max-w-3xl">
										<p class="text-sm text-gray-400 italic">
											{testimonial.context}
										</p>
									</div>
								{/if}

								<!-- Entity Links -->
								{#if testimonial.entityLinks && testimonial.entityLinks.length > 0}
									<div class="flex flex-wrap items-center justify-center gap-2 mb-4">
										{#each testimonial.entityLinks as link}
											<a
												href={link.entityType === 'project' ? `/projects/${link.entityId}` : `/skills#${link.entityId}`}
												class="flex items-center gap-1 px-3 py-1 text-xs rounded-full bg-gray-700/50 text-gray-300 border border-gray-600 hover:border-blue-500/50 hover:text-blue-300 transition-colors"
											>
												{#if link.entityType === 'project'}
													<FolderKanban class="w-3 h-3" />
													<span>Project</span>
												{:else}
													<Code2 class="w-3 h-3" />
													<span>Skill</span>
												{/if}
												<ExternalLink class="w-3 h-3" />
											</a>
										{/each}
									</div>
								{/if}

								<!-- Author Info -->
								<div class="flex flex-col items-center gap-2">
									<div class="flex items-center gap-2">
										<h4 class="text-white font-semibold text-xl">
											{testimonial.authorName}
										</h4>
										{#if testimonial.linkedinUrl}
											<a
												href={testimonial.linkedinUrl}
												target="_blank"
												rel="noopener noreferrer"
												class="text-blue-400 hover:text-blue-300 transition-colors"
												aria-label="LinkedIn Profile"
											>
												<Linkedin class="w-5 h-5" />
											</a>
										{/if}
									</div>
									<div class="flex items-center gap-2 text-gray-400">
										{#if testimonial.authorRole}
											<span>{testimonial.authorRole}</span>
										{/if}
										{#if testimonial.authorRole && testimonial.authorCompany}
											<span>•</span>
										{/if}
										{#if testimonial.authorCompany}
											<span>{testimonial.authorCompany}</span>
										{/if}
									</div>

									<!-- Video Link -->
									{#if testimonial.videoUrl}
										<a
											href={testimonial.videoUrl}
											target="_blank"
											rel="noopener noreferrer"
											class="mt-2 flex items-center gap-2 px-4 py-2 text-sm rounded-lg bg-blue-600/20 text-blue-300 border border-blue-500/30 hover:bg-blue-600/30 transition-colors"
										>
											<Play class="w-4 h-4" />
											<span>Watch Video</span>
										</a>
									{/if}
								</div>
							</div>
						{/if}
					</div>
				{/each}
			</div>

			<!-- Navigation Buttons -->
			{#if testimonials.length > 1}
				<button
					onclick={prevTestimonial}
					class="absolute left-4 top-1/2 -translate-y-1/2 w-10 h-10 rounded-full bg-gray-700/80 hover:bg-gray-600/80 border border-gray-600 flex items-center justify-center text-white transition-all duration-200 hover:scale-110"
					aria-label="Previous testimonial"
				>
					<ChevronLeft class="w-6 h-6" />
				</button>
				<button
					onclick={nextTestimonial}
					class="absolute right-4 top-1/2 -translate-y-1/2 w-10 h-10 rounded-full bg-gray-700/80 hover:bg-gray-600/80 border border-gray-600 flex items-center justify-center text-white transition-all duration-200 hover:scale-110"
					aria-label="Next testimonial"
				>
					<ChevronRight class="w-6 h-6" />
				</button>
			{/if}

			<!-- Dots Indicator -->
			{#if testimonials.length > 1}
				<div class="flex justify-center gap-2 mt-6">
					{#each testimonials as _, index}
						<button
							onclick={() => goToTestimonial(index)}
							class="w-2 h-2 rounded-full transition-all duration-200 {index === currentIndex
								? 'bg-blue-500 w-8'
								: 'bg-gray-600 hover:bg-gray-500'}"
							aria-label="Go to testimonial {index + 1}"
						></button>
					{/each}
				</div>
			{/if}
		</div>
	{/if}
</div>

<style>
	@keyframes fadeIn {
		from {
			opacity: 0;
			transform: translateY(20px);
		}
		to {
			opacity: 1;
			transform: translateY(0);
		}
	}
</style>

