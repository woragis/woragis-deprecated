import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import '../bloc/experience_bloc.dart';

class CreateExperiencePage extends StatefulWidget {
  final String? experienceId; // null for create, not null for edit

  const CreateExperiencePage({super.key, this.experienceId});

  @override
  State<CreateExperiencePage> createState() => _CreateExperiencePageState();
}

class _CreateExperiencePageState extends State<CreateExperiencePage> {
  final _formKey = GlobalKey<FormState>();
  final _titleController = TextEditingController();
  final _companyController = TextEditingController();
  final _periodController = TextEditingController();
  final _locationController = TextEditingController();
  final _descriptionController = TextEditingController();
  final _iconUrlController = TextEditingController();
  final _achievementController = TextEditingController();
  final _technologyController = TextEditingController();
  
  List<String> _achievements = [];
  List<String> _technologies = [];
  bool _isVisible = true;
  bool _isLoading = false;
  bool _isEditMode = false;
  AutovalidateMode _autovalidateMode = AutovalidateMode.disabled;

  @override
  void initState() {
    super.initState();
    _isEditMode = widget.experienceId != null;
    
    if (_isEditMode) {
      // Load experience for editing
      context.read<ExperienceBloc>().add(GetExperienceByIdRequested(widget.experienceId!));
    }
  }

  @override
  void dispose() {
    _titleController.dispose();
    _companyController.dispose();
    _periodController.dispose();
    _locationController.dispose();
    _descriptionController.dispose();
    _iconUrlController.dispose();
    _achievementController.dispose();
    _technologyController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(_isEditMode ? 'Edit Experience' : 'Add Experience'),
        backgroundColor: Colors.indigo.shade600,
        foregroundColor: Colors.white,
        elevation: 0,
        automaticallyImplyLeading: true,
        actions: [
          BlocBuilder<ExperienceBloc, ExperienceState>(
            builder: (context, state) {
              final isLoading = state is ExperienceLoading || _isLoading;
              return TextButton(
                onPressed: isLoading ? null : _saveExperience,
                child: isLoading
                    ? const SizedBox(
                        width: 16,
                        height: 16,
                        child: CircularProgressIndicator(
                          strokeWidth: 2,
                          valueColor: AlwaysStoppedAnimation<Color>(Colors.white),
                        ),
                      )
                    : const Text(
                        'Save',
                        style: TextStyle(
                          color: Colors.white,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
              );
            },
          ),
        ],
      ),
      body: BlocConsumer<ExperienceBloc, ExperienceState>(
        listener: (context, state) {
          if (state is ExperienceError) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Column(
                  mainAxisSize: MainAxisSize.min,
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Row(
                      children: [
                        Icon(Icons.error_outline, color: Colors.white),
                        SizedBox(width: 8),
                        Text('Failed to save experience', 
                            style: TextStyle(fontWeight: FontWeight.bold)),
                      ],
                    ),
                    const SizedBox(height: 4),
                    Text(state.message, style: const TextStyle(fontSize: 12)),
                  ],
                ),
                backgroundColor: Colors.red,
                duration: const Duration(seconds: 5),
                behavior: SnackBarBehavior.floating,
                action: SnackBarAction(
                  label: 'Retry',
                  textColor: Colors.white,
                  onPressed: _saveExperience,
                ),
              ),
            );
            setState(() {
              _isLoading = false;
            });
          } else if (state is ExperienceLoaded && _isEditMode) {
            _populateForm(state.experience);
          } else if (state is ExperienceCreated) {
            ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(
                content: Row(
                  children: [
                    Icon(Icons.check_circle, color: Colors.white),
                    SizedBox(width: 12),
                    Expanded(child: Text('Experience created successfully!')),
                  ],
                ),
                backgroundColor: Colors.green,
                behavior: SnackBarBehavior.floating,
                duration: Duration(seconds: 3),
              ),
            );
            context.push('/experience');
          } else if (state is ExperienceUpdated) {
            ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(
                content: Row(
                  children: [
                    Icon(Icons.check_circle, color: Colors.white),
                    SizedBox(width: 12),
                    Expanded(child: Text('Experience updated successfully!')),
                  ],
                ),
                backgroundColor: Colors.green,
                behavior: SnackBarBehavior.floating,
                duration: Duration(seconds: 3),
              ),
            );
            context.push('/experience');
          }
        },
        builder: (context, state) {
          final isLoading = state is ExperienceLoading || _isLoading;
          
          if (state is ExperienceLoading && !_isLoading && _isEditMode) {
            return const Center(child: CircularProgressIndicator());
          }

          return Stack(
            children: [
              SingleChildScrollView(
                padding: const EdgeInsets.all(16),
                child: Form(
                  key: _formKey,
                  autovalidateMode: _autovalidateMode,
                  child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Basic Information
                  Card(
                    child: Padding(
                      padding: const EdgeInsets.all(16),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          const Text(
                            'Basic Information',
                            style: TextStyle(
                              fontSize: 18,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                          const SizedBox(height: 16),
                          TextFormField(
                            controller: _titleController,
                            decoration: const InputDecoration(
                              labelText: 'Job Title *',
                              hintText: 'e.g., Senior Software Engineer',
                              border: OutlineInputBorder(),
                            ),
                            validator: (value) {
                              if (value == null || value.trim().isEmpty) {
                                return 'Job title is required';
                              }
                              return null;
                            },
                          ),
                          const SizedBox(height: 16),
                          TextFormField(
                            controller: _companyController,
                            decoration: const InputDecoration(
                              labelText: 'Company *',
                              hintText: 'e.g., Google, Microsoft',
                              border: OutlineInputBorder(),
                            ),
                            validator: (value) {
                              if (value == null || value.trim().isEmpty) {
                                return 'Company name is required';
                              }
                              return null;
                            },
                          ),
                          const SizedBox(height: 16),
                          TextFormField(
                            controller: _periodController,
                            decoration: const InputDecoration(
                              labelText: 'Employment Period *',
                              hintText: 'e.g., Jan 2020 - Present',
                              border: OutlineInputBorder(),
                            ),
                            validator: (value) {
                              if (value == null || value.trim().isEmpty) {
                                return 'Employment period is required';
                              }
                              return null;
                            },
                          ),
                          const SizedBox(height: 16),
                          TextFormField(
                            controller: _locationController,
                            decoration: const InputDecoration(
                              labelText: 'Location',
                              hintText: 'e.g., San Francisco, CA',
                              border: OutlineInputBorder(),
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                  const SizedBox(height: 16),

                  // Company Icon
                  Card(
                    child: Padding(
                      padding: const EdgeInsets.all(16),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          const Text(
                            'Company Icon',
                            style: TextStyle(
                              fontSize: 18,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                          const SizedBox(height: 16),
                          TextFormField(
                            controller: _iconUrlController,
                            decoration: const InputDecoration(
                              labelText: 'Icon URL',
                              hintText: 'https://example.com/company-logo.png',
                              border: OutlineInputBorder(),
                            ),
                            validator: (value) {
                              if (value != null && value.isNotEmpty) {
                                final uri = Uri.tryParse(value);
                                if (uri == null || (!uri.hasScheme || !uri.hasAuthority)) {
                                  return 'Please enter a valid URL';
                                }
                              }
                              return null;
                            },
                            onChanged: (value) {
                              setState(() {}); // Rebuild to show preview
                            },
                          ),
                          if (_iconUrlController.text.isNotEmpty) ...[
                            const SizedBox(height: 16),
                            Center(
                              child: CircleAvatar(
                                radius: 40,
                                backgroundColor: Colors.indigo.shade100,
                                backgroundImage: NetworkImage(_iconUrlController.text),
                                onBackgroundImageError: (exception, stackTrace) {
                                  // Handle image error
                                },
                                child: _iconUrlController.text.isEmpty
                                    ? Icon(Icons.work, size: 40, color: Colors.indigo.shade600)
                                    : null,
                              ),
                            ),
                          ],
                        ],
                      ),
                    ),
                  ),
                  const SizedBox(height: 16),

                  // Description
                  Card(
                    child: Padding(
                      padding: const EdgeInsets.all(16),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          const Text(
                            'Job Description',
                            style: TextStyle(
                              fontSize: 18,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                          const SizedBox(height: 16),
                          TextFormField(
                            controller: _descriptionController,
                            decoration: const InputDecoration(
                              labelText: 'Description *',
                              hintText: 'Describe your role and responsibilities...',
                              border: OutlineInputBorder(),
                              alignLabelWithHint: true,
                            ),
                            maxLines: 5,
                            validator: (value) {
                              if (value == null || value.trim().isEmpty) {
                                return 'Job description is required';
                              }
                              return null;
                            },
                          ),
                        ],
                      ),
                    ),
                  ),
                  const SizedBox(height: 16),

                  // Achievements
                  Card(
                    child: Padding(
                      padding: const EdgeInsets.all(16),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          const Text(
                            'Key Achievements',
                            style: TextStyle(
                              fontSize: 18,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                          const SizedBox(height: 16),
                          Row(
                            children: [
                              Expanded(
                                child: TextFormField(
                                  controller: _achievementController,
                                  decoration: const InputDecoration(
                                    labelText: 'Add Achievement',
                                    hintText: 'e.g., Led a team of 5 developers',
                                    border: OutlineInputBorder(),
                                  ),
                                  onFieldSubmitted: (value) {
                                    if (value.trim().isNotEmpty) {
                                      _addAchievement(value.trim());
                                    }
                                  },
                                ),
                              ),
                              const SizedBox(width: 8),
                              IconButton(
                                onPressed: () {
                                  if (_achievementController.text.trim().isNotEmpty) {
                                    _addAchievement(_achievementController.text.trim());
                                  }
                                },
                                icon: const Icon(Icons.add),
                                style: IconButton.styleFrom(
                                  backgroundColor: Colors.indigo.shade600,
                                  foregroundColor: Colors.white,
                                ),
                              ),
                            ],
                          ),
                          if (_achievements.isNotEmpty) ...[
                            const SizedBox(height: 16),
                            Wrap(
                              spacing: 8,
                              runSpacing: 8,
                              children: _achievements.asMap().entries.map((entry) {
                                return Chip(
                                  label: Text(entry.value),
                                  deleteIcon: const Icon(Icons.close, size: 18),
                                  onDeleted: () => _removeAchievement(entry.key),
                                  backgroundColor: Colors.indigo.shade50,
                                );
                              }).toList(),
                            ),
                          ],
                        ],
                      ),
                    ),
                  ),
                  const SizedBox(height: 16),

                  // Technologies
                  Card(
                    child: Padding(
                      padding: const EdgeInsets.all(16),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          const Text(
                            'Technologies & Skills',
                            style: TextStyle(
                              fontSize: 18,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                          const SizedBox(height: 16),
                          Row(
                            children: [
                              Expanded(
                                child: TextFormField(
                                  controller: _technologyController,
                                  decoration: const InputDecoration(
                                    labelText: 'Add Technology',
                                    hintText: 'e.g., React, Node.js, Python',
                                    border: OutlineInputBorder(),
                                  ),
                                  onFieldSubmitted: (value) {
                                    if (value.trim().isNotEmpty) {
                                      _addTechnology(value.trim());
                                    }
                                  },
                                ),
                              ),
                              const SizedBox(width: 8),
                              IconButton(
                                onPressed: () {
                                  if (_technologyController.text.trim().isNotEmpty) {
                                    _addTechnology(_technologyController.text.trim());
                                  }
                                },
                                icon: const Icon(Icons.add),
                                style: IconButton.styleFrom(
                                  backgroundColor: Colors.indigo.shade600,
                                  foregroundColor: Colors.white,
                                ),
                              ),
                            ],
                          ),
                          if (_technologies.isNotEmpty) ...[
                            const SizedBox(height: 16),
                            Wrap(
                              spacing: 8,
                              runSpacing: 8,
                              children: _technologies.asMap().entries.map((entry) {
                                return Chip(
                                  label: Text(entry.value),
                                  deleteIcon: const Icon(Icons.close, size: 18),
                                  onDeleted: () => _removeTechnology(entry.key),
                                  backgroundColor: Colors.indigo.shade100,
                                );
                              }).toList(),
                            ),
                          ],
                        ],
                      ),
                    ),
                  ),
                  const SizedBox(height: 16),

                  // Visibility Settings
                  Card(
                    child: Padding(
                      padding: const EdgeInsets.all(16),
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          const Text(
                            'Visibility Settings',
                            style: TextStyle(
                              fontSize: 18,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                          const SizedBox(height: 16),
                          SwitchListTile(
                            title: const Text('Make this experience visible'),
                            subtitle: const Text('Visible experiences will be shown on your portfolio'),
                            value: _isVisible,
                            onChanged: (value) {
                              setState(() {
                                _isVisible = value;
                              });
                            },
                            activeColor: Colors.indigo.shade600,
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
                      onPressed: _isLoading ? null : _saveExperience,
                      style: ElevatedButton.styleFrom(
                        backgroundColor: Colors.indigo.shade600,
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
                          : Text(
                              _isEditMode ? 'Update Experience' : 'Create Experience',
                              style: const TextStyle(fontSize: 16),
                            ),
                    ),
                  ),
                ],
                  ),
                ),
              ),
              // Loading overlay
              if (isLoading)
                Container(
                  color: Colors.black.withValues(alpha: 0.3),
                  child: const Center(
                    child: Card(
                      child: Padding(
                        padding: EdgeInsets.all(24),
                        child: Column(
                          mainAxisSize: MainAxisSize.min,
                          children: [
                            CircularProgressIndicator(),
                            SizedBox(height: 16),
                            Text(
                              'Saving experience...',
                              style: TextStyle(fontSize: 16, fontWeight: FontWeight.w500),
                            ),
                          ],
                        ),
                      ),
                    ),
                  ),
                ),
            ],
          );
        },
      ),
    );
  }

  void _populateForm(dynamic experience) {
    _titleController.text = experience.title;
    _companyController.text = experience.company;
    _periodController.text = experience.period;
    _locationController.text = experience.location;
    _descriptionController.text = experience.description;
    _iconUrlController.text = experience.icon;
    _achievements = List<String>.from(experience.achievements);
    _technologies = List<String>.from(experience.technologies);
    _isVisible = experience.visible;
  }

  void _addAchievement(String achievement) {
    setState(() {
      _achievements.add(achievement);
      _achievementController.clear();
    });
  }

  void _removeAchievement(int index) {
    setState(() {
      _achievements.removeAt(index);
    });
  }

  void _addTechnology(String technology) {
    setState(() {
      _technologies.add(technology);
      _technologyController.clear();
    });
  }

  void _removeTechnology(int index) {
    setState(() {
      _technologies.removeAt(index);
    });
  }

  void _saveExperience() {
    // Enable auto-validate after first submit attempt
    setState(() {
      _autovalidateMode = AutovalidateMode.onUserInteraction;
    });
    
    if (_formKey.currentState!.validate()) {
      setState(() {
        _isLoading = true;
      });

      if (_isEditMode) {
        context.read<ExperienceBloc>().add(UpdateExperienceRequested(
          id: widget.experienceId!,
          title: _titleController.text.trim(),
          company: _companyController.text.trim(),
          period: _periodController.text.trim(),
          location: _locationController.text.trim(),
          description: _descriptionController.text.trim(),
          achievements: _achievements,
          technologies: _technologies,
          icon: _iconUrlController.text.trim(),
          visible: _isVisible,
        ));
      } else {
        context.read<ExperienceBloc>().add(CreateExperienceRequested(
          title: _titleController.text.trim(),
          company: _companyController.text.trim(),
          period: _periodController.text.trim(),
          location: _locationController.text.trim(),
          description: _descriptionController.text.trim(),
          achievements: _achievements,
          technologies: _technologies,
          icon: _iconUrlController.text.trim(),
          order: 0, // Default order, can be changed later
          visible: _isVisible,
        ));
      }
    }
  }
}
