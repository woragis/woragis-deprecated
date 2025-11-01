import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import '../bloc/settings_bloc.dart';

class SocialMediaSettingsPage extends StatefulWidget {
  const SocialMediaSettingsPage({super.key});

  @override
  State<SocialMediaSettingsPage> createState() => _SocialMediaSettingsPageState();
}

class _SocialMediaSettingsPageState extends State<SocialMediaSettingsPage> {
  final _formKey = GlobalKey<FormState>();
  final _githubController = TextEditingController();
  final _linkedinController = TextEditingController();
  final _twitterController = TextEditingController();
  final _instagramController = TextEditingController();
  final _facebookController = TextEditingController();
  final _youtubeController = TextEditingController();
  final _discordController = TextEditingController();
  final _twitchController = TextEditingController();
  final _tiktokController = TextEditingController();
  
  bool _isLoading = false;

  @override
  void initState() {
    super.initState();
    // Load social media settings
    context.read<SettingsBloc>().add(GetSocialMediaSettingsRequested());
  }

  @override
  void dispose() {
    _githubController.dispose();
    _linkedinController.dispose();
    _twitterController.dispose();
    _instagramController.dispose();
    _facebookController.dispose();
    _youtubeController.dispose();
    _discordController.dispose();
    _twitchController.dispose();
    _tiktokController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Social Media'),
        backgroundColor: Colors.purple.shade600,
        foregroundColor: Colors.white,
        elevation: 0,
        automaticallyImplyLeading: true,
        actions: [
          TextButton(
            onPressed: _isLoading ? null : _saveSettings,
            child: Text(
              'Save',
              style: TextStyle(
                color: Colors.white,
                fontWeight: FontWeight.bold,
              ),
            ),
          ),
        ],
      ),
      body: BlocConsumer<SettingsBloc, SettingsState>(
        listener: (context, state) {
          if (state is SettingsError) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Text(state.message),
                backgroundColor: Colors.red,
              ),
            );
          } else if (state is SocialMediaSettingsLoaded) {
            _populateForm(state.settings);
          } else if (state is SocialMediaSettingsUpdated) {
            ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(
                content: Text('Social media settings updated successfully'),
                backgroundColor: Colors.green,
              ),
            );
            setState(() {
              _isLoading = false;
            });
          }
        },
        builder: (context, state) {
          if (state is SettingsLoading && !_isLoading) {
            return const Center(child: CircularProgressIndicator());
          }

          return SingleChildScrollView(
            padding: const EdgeInsets.all(16),
            child: Form(
              key: _formKey,
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Header
                  Card(
                    child: Padding(
                      padding: const EdgeInsets.all(16),
                      child: Column(
                        children: [
                          Icon(
                            Icons.share,
                            size: 48,
                            color: Colors.purple.shade600,
                          ),
                          const SizedBox(height: 16),
                          const Text(
                            'Social Media Links',
                            style: TextStyle(
                              fontSize: 20,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                          const SizedBox(height: 8),
                          Text(
                            'Connect your social media accounts to showcase your online presence',
                            style: TextStyle(
                              color: Colors.grey.shade600,
                              fontSize: 14,
                            ),
                            textAlign: TextAlign.center,
                          ),
                        ],
                      ),
                    ),
                  ),
                  const SizedBox(height: 16),

                  // Professional Platforms
                  Card(
                    child: Padding(
                      padding: const EdgeInsets.all(16),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          const Text(
                            'Professional Platforms',
                            style: TextStyle(
                              fontSize: 18,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                          const SizedBox(height: 16),
                          _buildSocialMediaField(
                            controller: _githubController,
                            label: 'GitHub',
                            icon: Icons.code,
                            placeholder: 'https://github.com/username',
                            color: Colors.grey.shade800,
                          ),
                          const SizedBox(height: 16),
                          _buildSocialMediaField(
                            controller: _linkedinController,
                            label: 'LinkedIn',
                            icon: Icons.business,
                            placeholder: 'https://linkedin.com/in/username',
                            color: Colors.blue.shade700,
                          ),
                          const SizedBox(height: 16),
                          _buildSocialMediaField(
                            controller: _twitterController,
                            label: 'Twitter',
                            icon: Icons.alternate_email,
                            placeholder: 'https://twitter.com/username',
                            color: Colors.blue.shade400,
                          ),
                        ],
                      ),
                    ),
                  ),
                  const SizedBox(height: 16),

                  // Creative Platforms
                  Card(
                    child: Padding(
                      padding: const EdgeInsets.all(16),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          const Text(
                            'Creative Platforms',
                            style: TextStyle(
                              fontSize: 18,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                          const SizedBox(height: 16),
                          _buildSocialMediaField(
                            controller: _instagramController,
                            label: 'Instagram',
                            icon: Icons.camera_alt,
                            placeholder: 'https://instagram.com/username',
                            color: Colors.pink.shade600,
                          ),
                          const SizedBox(height: 16),
                          _buildSocialMediaField(
                            controller: _youtubeController,
                            label: 'YouTube',
                            icon: Icons.play_circle,
                            placeholder: 'https://youtube.com/@username',
                            color: Colors.red.shade600,
                          ),
                          const SizedBox(height: 16),
                          _buildSocialMediaField(
                            controller: _tiktokController,
                            label: 'TikTok',
                            icon: Icons.music_note,
                            placeholder: 'https://tiktok.com/@username',
                            color: Colors.black,
                          ),
                        ],
                      ),
                    ),
                  ),
                  const SizedBox(height: 16),

                  // Gaming & Community
                  Card(
                    child: Padding(
                      padding: const EdgeInsets.all(16),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          const Text(
                            'Gaming & Community',
                            style: TextStyle(
                              fontSize: 18,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                          const SizedBox(height: 16),
                          _buildSocialMediaField(
                            controller: _discordController,
                            label: 'Discord',
                            icon: Icons.chat,
                            placeholder: 'https://discord.gg/server',
                            color: Colors.indigo.shade600,
                          ),
                          const SizedBox(height: 16),
                          _buildSocialMediaField(
                            controller: _twitchController,
                            label: 'Twitch',
                            icon: Icons.live_tv,
                            placeholder: 'https://twitch.tv/username',
                            color: Colors.purple.shade600,
                          ),
                          const SizedBox(height: 16),
                          _buildSocialMediaField(
                            controller: _facebookController,
                            label: 'Facebook',
                            icon: Icons.facebook,
                            placeholder: 'https://facebook.com/username',
                            color: Colors.blue.shade800,
                          ),
                        ],
                      ),
                    ),
                  ),
                  const SizedBox(height: 24),

                  // Preview Section
                  Card(
                    child: Padding(
                      padding: const EdgeInsets.all(16),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          const Text(
                            'Preview',
                            style: TextStyle(
                              fontSize: 18,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                          const SizedBox(height: 16),
                          Wrap(
                            spacing: 12,
                            runSpacing: 12,
                            children: [
                              if (_githubController.text.isNotEmpty)
                                _buildPreviewChip('GitHub', Colors.grey.shade800, Icons.code),
                              if (_linkedinController.text.isNotEmpty)
                                _buildPreviewChip('LinkedIn', Colors.blue.shade700, Icons.business),
                              if (_twitterController.text.isNotEmpty)
                                _buildPreviewChip('Twitter', Colors.blue.shade400, Icons.alternate_email),
                              if (_instagramController.text.isNotEmpty)
                                _buildPreviewChip('Instagram', Colors.pink.shade600, Icons.camera_alt),
                              if (_youtubeController.text.isNotEmpty)
                                _buildPreviewChip('YouTube', Colors.red.shade600, Icons.play_circle),
                              if (_tiktokController.text.isNotEmpty)
                                _buildPreviewChip('TikTok', Colors.black, Icons.music_note),
                              if (_discordController.text.isNotEmpty)
                                _buildPreviewChip('Discord', Colors.indigo.shade600, Icons.chat),
                              if (_twitchController.text.isNotEmpty)
                                _buildPreviewChip('Twitch', Colors.purple.shade600, Icons.live_tv),
                              if (_facebookController.text.isNotEmpty)
                                _buildPreviewChip('Facebook', Colors.blue.shade800, Icons.facebook),
                            ],
                          ),
                          if (_githubController.text.isEmpty &&
                              _linkedinController.text.isEmpty &&
                              _twitterController.text.isEmpty &&
                              _instagramController.text.isEmpty &&
                              _youtubeController.text.isEmpty &&
                              _tiktokController.text.isEmpty &&
                              _discordController.text.isEmpty &&
                              _twitchController.text.isEmpty &&
                              _facebookController.text.isEmpty)
                            Container(
                              padding: const EdgeInsets.all(16),
                              decoration: BoxDecoration(
                                color: Colors.grey.shade50,
                                borderRadius: BorderRadius.circular(8),
                                border: Border.all(color: Colors.grey.shade300),
                              ),
                              child: Row(
                                children: [
                                  Icon(
                                    Icons.info_outline,
                                    color: Colors.grey.shade600,
                                  ),
                                  const SizedBox(width: 8),
                                  Text(
                                    'No social media links added yet',
                                    style: TextStyle(
                                      color: Colors.grey.shade600,
                                      fontStyle: FontStyle.italic,
                                    ),
                                  ),
                                ],
                              ),
                            ),
                        ],
                      ),
                    ),
                  ),
                  const SizedBox(height: 24),

                  // Save Button
                  SizedBox(
                    width: double.infinity,
                    child: ElevatedButton(
                      onPressed: _isLoading ? null : _saveSettings,
                      style: ElevatedButton.styleFrom(
                        backgroundColor: Colors.purple.shade600,
                        foregroundColor: Colors.white,
                        padding: const EdgeInsets.symmetric(vertical: 16),
                      ),
                      child: _isLoading
                          ? const SizedBox(
                              height: 20,
                              width: 20,
                              child: CircularProgressIndicator(
                                strokeWidth: 2,
                                valueColor: AlwaysStoppedAnimation<Color>(
                                  Colors.white,
                                ),
                              ),
                            )
                          : const Text(
                              'Save Settings',
                              style: TextStyle(fontSize: 16),
                            ),
                    ),
                  ),
                ],
              ),
            ),
          );
        },
      ),
    );
  }

  Widget _buildSocialMediaField({
    required TextEditingController controller,
    required String label,
    required IconData icon,
    required String placeholder,
    required Color color,
  }) {
    return TextFormField(
      controller: controller,
      decoration: InputDecoration(
        labelText: label,
        hintText: placeholder,
        border: const OutlineInputBorder(),
        prefixIcon: Icon(icon, color: color),
      ),
      validator: (value) {
        if (value != null && value.isNotEmpty) {
          if (Uri.tryParse(value)?.hasAbsolutePath != true) {
            return 'Please enter a valid URL';
          }
        }
        return null;
      },
    );
  }

  Widget _buildPreviewChip(String label, Color color, IconData icon) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
      decoration: BoxDecoration(
        color: color.withValues(alpha: 0.1),
        borderRadius: BorderRadius.circular(20),
        border: Border.all(color: color.withValues(alpha: 0.3)),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(icon, size: 16, color: color),
          const SizedBox(width: 6),
          Text(
            label,
            style: TextStyle(
              color: color,
              fontWeight: FontWeight.w500,
              fontSize: 12,
            ),
          ),
        ],
      ),
    );
  }

  void _populateForm(Map<String, String> settings) {
    _githubController.text = settings['github'] ?? '';
    _linkedinController.text = settings['linkedin'] ?? '';
    _twitterController.text = settings['twitter'] ?? '';
    _instagramController.text = settings['instagram'] ?? '';
    _facebookController.text = settings['facebook'] ?? '';
    _youtubeController.text = settings['youtube'] ?? '';
    _discordController.text = settings['discord'] ?? '';
    _twitchController.text = settings['twitch'] ?? '';
    _tiktokController.text = settings['tiktok'] ?? '';
  }

  void _saveSettings() {
    if (_formKey.currentState!.validate()) {
      setState(() {
        _isLoading = true;
      });

      final settings = {
        'github': _githubController.text.trim(),
        'linkedin': _linkedinController.text.trim(),
        'twitter': _twitterController.text.trim(),
        'instagram': _instagramController.text.trim(),
        'facebook': _facebookController.text.trim(),
        'youtube': _youtubeController.text.trim(),
        'discord': _discordController.text.trim(),
        'twitch': _twitchController.text.trim(),
        'tiktok': _tiktokController.text.trim(),
      };

      context.read<SettingsBloc>().add(UpdateSocialMediaSettingsRequested(settings));
    }
  }
}
