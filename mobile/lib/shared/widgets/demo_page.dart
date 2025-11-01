import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import '../styles/woragis_styles.dart';
import 'woragis_button.dart';
import 'woragis_card.dart';
import 'woragis_input.dart';

class WoragisDemoPage extends StatefulWidget {
  const WoragisDemoPage({super.key});

  @override
  State<WoragisDemoPage> createState() => _WoragisDemoPageState();
}

class _WoragisDemoPageState extends State<WoragisDemoPage>
    with TickerProviderStateMixin {
  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();
  final _searchController = TextEditingController();
  
  late AnimationController _animationController;
  late Animation<double> _fadeAnimation;
  late Animation<Offset> _slideAnimation;

  @override
  void initState() {
    super.initState();
    _animationController = AnimationController(
      duration: const Duration(milliseconds: 1000),
      vsync: this,
    );
    
    _fadeAnimation = Tween<double>(
      begin: 0.0,
      end: 1.0,
    ).animate(CurvedAnimation(
      parent: _animationController,
      curve: const Interval(0.0, 0.6, curve: Curves.easeOut),
    ));
    
    _slideAnimation = Tween<Offset>(
      begin: const Offset(0, 0.3),
      end: Offset.zero,
    ).animate(CurvedAnimation(
      parent: _animationController,
      curve: const Interval(0.0, 0.8, curve: Curves.easeOutCubic),
    ));
    
    _animationController.forward();
  }

  @override
  void dispose() {
    _animationController.dispose();
    _emailController.dispose();
    _passwordController.dispose();
    _searchController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;
    final isDark = theme.brightness == Brightness.dark;

    return Scaffold(
      backgroundColor: colorScheme.surface,
      body: SafeArea(
        child: FadeTransition(
          opacity: _fadeAnimation,
          child: SlideTransition(
            position: _slideAnimation,
            child: CustomScrollView(
              slivers: [
                // App Bar
                SliverAppBar(
                  expandedHeight: 120,
                  floating: false,
                  pinned: true,
                  backgroundColor: Colors.transparent,
                  elevation: 0,
                  flexibleSpace: FlexibleSpaceBar(
                    title: Text(
                      'Woragis Design System',
                      style: WoragisTextStyles.headlineMedium(colorScheme.onSurface),
                    ),
                    background: Container(
                      decoration: BoxDecoration(
                        gradient: LinearGradient(
                          colors: isDark
                              ? [
                                  colorScheme.primary.withValues(alpha: 0.1),
                                  colorScheme.secondary.withValues(alpha: 0.1),
                                ]
                              : [
                                  colorScheme.primary.withValues(alpha: 0.05),
                                  colorScheme.secondary.withValues(alpha: 0.05),
                                ],
                          begin: Alignment.topLeft,
                          end: Alignment.bottomRight,
                        ),
                      ),
                    ),
                  ),
                  actions: [
                    IconButton(
                      icon: Icon(
                        isDark ? Icons.light_mode : Icons.dark_mode,
                        color: colorScheme.onSurface,
                      ),
                      onPressed: () {
                        HapticFeedback.lightImpact();
                        // Toggle theme logic would go here
                      },
                    ),
                  ],
                ),

                // Content
                SliverPadding(
                  padding: const EdgeInsets.all(WoragisStyles.spacingMd),
                  sliver: SliverList(
                    delegate: SliverChildListDelegate([
                      // Hero Section
                      WoragisCard(
                        variant: WoragisCardVariant.gradient,
                        padding: const EdgeInsets.all(WoragisStyles.spacingXl),
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.start,
                          children: [
                            Text(
                              'Welcome to Woragis',
                              style: WoragisTextStyles.displaySmall(colorScheme.onSurface),
                            ),
                            const SizedBox(height: WoragisStyles.spacingSm),
                            Text(
                              'A modern Flutter design system inspired by your Next.js portfolio',
                              style: WoragisTextStyles.bodyLarge(
                                colorScheme.onSurface.withValues(alpha: 0.8),
                              ),
                            ),
                            const SizedBox(height: WoragisStyles.spacingLg),
                            Row(
                              children: [
                                Expanded(
                                  child: WoragisButton(
                                    text: 'Get Started',
                                    variant: WoragisButtonVariant.modern,
                                    size: WoragisButtonSize.large,
                                    icon: Icons.rocket_launch,
                                    onPressed: () {
                                      HapticFeedback.mediumImpact();
                                    },
                                  ),
                                ),
                                const SizedBox(width: WoragisStyles.spacingMd),
                                Expanded(
                                  child: WoragisButton(
                                    text: 'Learn More',
                                    variant: WoragisButtonVariant.outline,
                                    size: WoragisButtonSize.large,
                                    onPressed: () {
                                      HapticFeedback.lightImpact();
                                    },
                                  ),
                                ),
                              ],
                            ),
                          ],
                        ),
                      ),

                      const SizedBox(height: WoragisStyles.spacingXl),

                      // Buttons Section
                      _buildSection(
                        context,
                        'Buttons',
                        'Various button styles and variants',
                        [
                          _buildButtonGrid(context),
                        ],
                      ),

                      const SizedBox(height: WoragisStyles.spacingXl),

                      // Cards Section
                      _buildSection(
                        context,
                        'Cards',
                        'Different card styles and layouts',
                        [
                          _buildCardGrid(context),
                        ],
                      ),

                      const SizedBox(height: WoragisStyles.spacingXl),

                      // Inputs Section
                      _buildSection(
                        context,
                        'Inputs',
                        'Form inputs with different styles',
                        [
                          _buildInputExamples(context),
                        ],
                      ),

                      const SizedBox(height: WoragisStyles.spacingXl),

                      // Stats Section
                      _buildSection(
                        context,
                        'Statistics',
                        'Display data with style',
                        [
                          _buildStatsGrid(context),
                        ],
                      ),

                      const SizedBox(height: WoragisStyles.spacing3Xl),
                    ]),
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildSection(
    BuildContext context,
    String title,
    String subtitle,
    List<Widget> children,
  ) {
    final theme = Theme.of(context);
    final colorScheme = theme.colorScheme;

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          title,
          style: WoragisTextStyles.headlineMedium(colorScheme.onSurface),
        ),
        const SizedBox(height: WoragisStyles.spacingXs),
        Text(
          subtitle,
          style: WoragisTextStyles.bodyMedium(
            colorScheme.onSurface.withValues(alpha: 0.7),
          ),
        ),
        const SizedBox(height: WoragisStyles.spacingLg),
        ...children,
      ],
    );
  }

  Widget _buildButtonGrid(BuildContext context) {
    return Column(
      children: [
        // Primary Buttons
        Row(
          children: [
            Expanded(
              child: WoragisButton(
                text: 'Primary',
                variant: WoragisButtonVariant.primary,
                onPressed: () => HapticFeedback.lightImpact(),
              ),
            ),
            const SizedBox(width: WoragisStyles.spacingMd),
            Expanded(
              child: WoragisButton(
                text: 'Modern',
                variant: WoragisButtonVariant.modern,
                onPressed: () => HapticFeedback.lightImpact(),
              ),
            ),
          ],
        ),
        const SizedBox(height: WoragisStyles.spacingMd),
        
        // Secondary Buttons
        Row(
          children: [
            Expanded(
              child: WoragisButton(
                text: 'Outline',
                variant: WoragisButtonVariant.outline,
                onPressed: () => HapticFeedback.lightImpact(),
              ),
            ),
            const SizedBox(width: WoragisStyles.spacingMd),
            Expanded(
              child: WoragisButton(
                text: 'Glass',
                variant: WoragisButtonVariant.glass,
                onPressed: () => HapticFeedback.lightImpact(),
              ),
            ),
          ],
        ),
        const SizedBox(height: WoragisStyles.spacingMd),
        
        // Special Buttons
        Row(
          children: [
            Expanded(
              child: WoragisButton(
                text: 'Gradient',
                variant: WoragisButtonVariant.gradient,
                onPressed: () => HapticFeedback.lightImpact(),
              ),
            ),
            const SizedBox(width: WoragisStyles.spacingMd),
            Expanded(
              child: WoragisButton(
                text: 'Ghost',
                variant: WoragisButtonVariant.ghost,
                onPressed: () => HapticFeedback.lightImpact(),
              ),
            ),
          ],
        ),
      ],
    );
  }

  Widget _buildCardGrid(BuildContext context) {
    return Column(
      children: [
        // Feature Cards
        Row(
          children: [
            Expanded(
              child: WoragisFeatureCard(
                title: 'Modern Design',
                subtitle: 'Clean and contemporary',
                icon: Icon(
                  Icons.design_services,
                  color: Theme.of(context).colorScheme.primary,
                ),
                variant: WoragisCardVariant.modern,
                onTap: () => HapticFeedback.lightImpact(),
              ),
            ),
            const SizedBox(width: WoragisStyles.spacingMd),
            Expanded(
              child: WoragisFeatureCard(
                title: 'Glass Effect',
                subtitle: 'Beautiful transparency',
                icon: Icon(
                  Icons.blur_on,
                  color: Theme.of(context).colorScheme.secondary,
                ),
                variant: WoragisCardVariant.glass,
                onTap: () => HapticFeedback.lightImpact(),
              ),
            ),
          ],
        ),
        const SizedBox(height: WoragisStyles.spacingMd),
        
        // Single gradient card
        WoragisCard(
          variant: WoragisCardVariant.gradient,
          padding: const EdgeInsets.all(WoragisStyles.spacingLg),
          child: Row(
            children: [
              Container(
                padding: const EdgeInsets.all(WoragisStyles.spacingMd),
                decoration: BoxDecoration(
                  color: Theme.of(context).colorScheme.primary.withValues(alpha: 0.1),
                  borderRadius: BorderRadius.circular(WoragisStyles.radiusMd),
                ),
                child: Icon(
                  Icons.star,
                  color: Theme.of(context).colorScheme.primary,
                  size: 32,
                ),
              ),
              const SizedBox(width: WoragisStyles.spacingLg),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      'Premium Features',
                      style: WoragisTextStyles.titleLarge(
                        Theme.of(context).colorScheme.onSurface,
                      ),
                    ),
                    const SizedBox(height: WoragisStyles.spacingXs),
                    Text(
                      'Access to all advanced functionality',
                      style: WoragisTextStyles.bodyMedium(
                        Theme.of(context).colorScheme.onSurface.withValues(alpha: 0.7),
                      ),
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }

  Widget _buildInputExamples(BuildContext context) {
    return Column(
      children: [
        // Basic inputs
        WoragisInput(
          label: 'Email Address',
          hintText: 'Enter your email',
          controller: _emailController,
          variant: WoragisInputVariant.modern,
          prefixIcon: Icon(
            Icons.email_outlined,
            color: Theme.of(context).colorScheme.onSurface.withValues(alpha: 0.6),
          ),
        ),
        const SizedBox(height: WoragisStyles.spacingLg),
        
        WoragisPasswordInput(
          label: 'Password',
          hintText: 'Enter your password',
          controller: _passwordController,
          variant: WoragisInputVariant.modern,
        ),
        const SizedBox(height: WoragisStyles.spacingLg),
        
        // Search input
        WoragisSearchInput(
          hintText: 'Search components...',
          controller: _searchController,
          onChanged: (value) {
            // Search logic would go here
          },
        ),
      ],
    );
  }

  Widget _buildStatsGrid(BuildContext context) {
    return Row(
      children: [
        Expanded(
          child: WoragisStatsCard(
            label: 'Projects',
            value: '24',
            icon: Icon(
              Icons.work_outline,
              color: Theme.of(context).colorScheme.primary,
            ),
          ),
        ),
        const SizedBox(width: WoragisStyles.spacingMd),
        Expanded(
          child: WoragisStatsCard(
            label: 'Clients',
            value: '18',
            icon: Icon(
              Icons.people_outline,
              color: Theme.of(context).colorScheme.secondary,
            ),
          ),
        ),
        const SizedBox(width: WoragisStyles.spacingMd),
        Expanded(
          child: WoragisStatsCard(
            label: 'Experience',
            value: '3+',
            icon: Icon(
              Icons.trending_up,
              color: Theme.of(context).colorScheme.tertiary,
            ),
          ),
        ),
      ],
    );
  }
}
