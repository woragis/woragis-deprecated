import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

class SettingsOverviewPage extends StatelessWidget {
  const SettingsOverviewPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Settings'),
        backgroundColor: Colors.grey.shade600,
        foregroundColor: Colors.white,
        elevation: 0,
      ),
      body: ListView(
        padding: const EdgeInsets.all(16),
        children: [
          // Core Profile Settings
          _buildSettingsCard(
            context,
            title: 'Core Profile',
            subtitle: 'Basic profile information and preferences',
            icon: Icons.person,
            color: Colors.blue,
            onTap: () => context.push('/settings/core-profile'),
          ),
          const SizedBox(height: 16),

          // Social Media Settings
          _buildSettingsCard(
            context,
            title: 'Social Media',
            subtitle: 'Connect your social media accounts',
            icon: Icons.share,
            color: Colors.purple,
            onTap: () => context.push('/settings/social-media'),
          ),
          const SizedBox(height: 16),

          // Contact Settings
          _buildSettingsCard(
            context,
            title: 'Contact Information',
            subtitle: 'Manage your contact details',
            icon: Icons.contact_mail,
            color: Colors.green,
            onTap: () => context.push('/settings/contact'),
          ),
          const SizedBox(height: 16),

          // Site Settings
          _buildSettingsCard(
            context,
            title: 'Site Settings',
            subtitle: 'Configure your portfolio site',
            icon: Icons.settings,
            color: Colors.orange,
            onTap: () => context.push('/settings/site'),
          ),
          const SizedBox(height: 24),

          // App Settings
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Icon(
                        Icons.phone_android,
                        color: Colors.grey.shade600,
                        size: 24,
                      ),
                      const SizedBox(width: 12),
                      const Text(
                        'App Settings',
                        style: TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                  _buildAppSettingTile(
                    context,
                    title: 'Theme',
                    subtitle: 'Choose your preferred theme',
                    icon: Icons.palette,
                    onTap: () {
                      // TODO: Implement theme settings
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(
                          content: Text('Theme settings coming soon!'),
                        ),
                      );
                    },
                  ),
                  _buildAppSettingTile(
                    context,
                    title: 'Language',
                    subtitle: 'Select your preferred language',
                    icon: Icons.language,
                    onTap: () {
                      // TODO: Implement language settings
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(
                          content: Text('Language settings coming soon!'),
                        ),
                      );
                    },
                  ),
                  _buildAppSettingTile(
                    context,
                    title: 'Notifications',
                    subtitle: 'Manage your notification preferences',
                    icon: Icons.notifications,
                    onTap: () {
                      // TODO: Implement notification settings
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(
                          content: Text('Notification settings coming soon!'),
                        ),
                      );
                    },
                  ),
                ],
              ),
            ),
          ),
          const SizedBox(height: 24),

          // Account Actions
          Card(
            child: Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Row(
                    children: [
                      Icon(
                        Icons.account_circle,
                        color: Colors.grey.shade600,
                        size: 24,
                      ),
                      const SizedBox(width: 12),
                      const Text(
                        'Account',
                        style: TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ],
                  ),
                  const SizedBox(height: 16),
                  _buildAppSettingTile(
                    context,
                    title: 'Change Password',
                    subtitle: 'Update your account password',
                    icon: Icons.lock,
                    onTap: () => context.push('/change-password'),
                  ),
                  _buildAppSettingTile(
                    context,
                    title: 'Export Data',
                    subtitle: 'Download your portfolio data',
                    icon: Icons.download,
                    onTap: () {
                      // TODO: Implement data export
                      ScaffoldMessenger.of(context).showSnackBar(
                        const SnackBar(
                          content: Text('Data export coming soon!'),
                        ),
                      );
                    },
                  ),
                  _buildAppSettingTile(
                    context,
                    title: 'Delete Account',
                    subtitle: 'Permanently delete your account',
                    icon: Icons.delete_forever,
                    onTap: () => _showDeleteAccountDialog(context),
                    isDestructive: true,
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildSettingsCard(
    BuildContext context, {
    required String title,
    required String subtitle,
    required IconData icon,
    required Color color,
    required VoidCallback onTap,
  }) {
    return Card(
      child: InkWell(
        onTap: onTap,
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Row(
            children: [
              Container(
                padding: const EdgeInsets.all(12),
                decoration: BoxDecoration(
                  color: color.withValues(alpha: 0.1),
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Icon(
                  icon,
                  color: color,
                  size: 24,
                ),
              ),
              const SizedBox(width: 16),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      title,
                      style: const TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      subtitle,
                      style: TextStyle(
                        fontSize: 14,
                        color: Colors.grey.shade600,
                      ),
                    ),
                  ],
                ),
              ),
              Icon(
                Icons.arrow_forward_ios,
                size: 16,
                color: Colors.grey.shade400,
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildAppSettingTile(
    BuildContext context, {
    required String title,
    required String subtitle,
    required IconData icon,
    required VoidCallback onTap,
    bool isDestructive = false,
  }) {
    return ListTile(
      leading: Icon(
        icon,
        color: isDestructive ? Colors.red : Colors.grey.shade600,
      ),
      title: Text(
        title,
        style: TextStyle(
          color: isDestructive ? Colors.red : null,
          fontWeight: FontWeight.w500,
        ),
      ),
      subtitle: Text(subtitle),
      trailing: const Icon(
        Icons.arrow_forward_ios,
        size: 16,
        color: Colors.grey,
      ),
      onTap: onTap,
    );
  }

  void _showDeleteAccountDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Delete Account'),
        content: const Text(
          'Are you sure you want to permanently delete your account? This action cannot be undone and will remove all your data.',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: const Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () {
              Navigator.pop(context);
              // TODO: Implement account deletion
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(
                  content: Text('Account deletion not implemented yet'),
                  backgroundColor: Colors.orange,
                ),
              );
            },
            style: ElevatedButton.styleFrom(
              backgroundColor: Colors.red,
              foregroundColor: Colors.white,
            ),
            child: const Text('Delete'),
          ),
        ],
      ),
    );
  }
}
