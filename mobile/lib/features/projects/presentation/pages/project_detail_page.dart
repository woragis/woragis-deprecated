import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import 'package:url_launcher/url_launcher.dart';
import 'package:image_picker/image_picker.dart';
import 'dart:io';
import '../bloc/projects_bloc.dart';

class ProjectDetailPage extends StatefulWidget {
  final String projectId;

  const ProjectDetailPage({
    super.key,
    required this.projectId,
  });

  @override
  State<ProjectDetailPage> createState() => _ProjectDetailPageState();
}

class _ProjectDetailPageState extends State<ProjectDetailPage> {
  bool _isEditing = false;
  bool _isSaving = false;
  late TextEditingController _titleController;
  late TextEditingController _descriptionController;
  late TextEditingController _longDescriptionController;
  late TextEditingController _contentController;
  late TextEditingController _githubUrlController;
  late TextEditingController _liveUrlController;
  late TextEditingController _videoUrlController;
  final ImagePicker _imagePicker = ImagePicker();
  String? _selectedImagePath;
  
  // Form validation errors
  String? _titleError;
  String? _descriptionError;
  String? _githubUrlError;
  String? _liveUrlError;
  String? _videoUrlError;

  @override
  void initState() {
    super.initState();
    _titleController = TextEditingController();
    _descriptionController = TextEditingController();
    _longDescriptionController = TextEditingController();
    _contentController = TextEditingController();
    _githubUrlController = TextEditingController();
    _liveUrlController = TextEditingController();
    _videoUrlController = TextEditingController();

    // Load the specific project
    context.read<ProjectsBloc>().add(GetProjectByIdRequested(widget.projectId));
  }

  @override
  void dispose() {
    _titleController.dispose();
    _descriptionController.dispose();
    _longDescriptionController.dispose();
    _contentController.dispose();
    _githubUrlController.dispose();
    _liveUrlController.dispose();
    _videoUrlController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Project Details'),
        backgroundColor: Colors.blue.shade600,
        foregroundColor: Colors.white,
        elevation: 0,
        actions: [
          if (!_isEditing)
            IconButton(
              onPressed: () => context.push('/projects/${widget.projectId}/edit'),
              icon: const Icon(Icons.edit),
            )
          else
            TextButton(
              onPressed: () {
                setState(() {
                  _isEditing = false;
                });
                // Reset form to original values
                context.read<ProjectsBloc>().add(GetProjectByIdRequested(widget.projectId));
              },
              child: const Text('Cancel'),
            ),
        ],
      ),
      body: BlocConsumer<ProjectsBloc, ProjectsState>(
        listener: (context, state) {
          if (state is ProjectsError) {
            setState(() {
              _isSaving = false;
            });
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Row(
                  children: [
                    const Icon(Icons.error_outline, color: Colors.white),
                    const SizedBox(width: 12),
                    Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          const Text(
                            'Error',
                            style: TextStyle(fontWeight: FontWeight.bold),
                          ),
                          Text(state.message),
                        ],
                      ),
                    ),
                  ],
                ),
                backgroundColor: Colors.red,
                duration: const Duration(seconds: 4),
              ),
            );
          } else if (state is ProjectUpdated) {
            setState(() {
              _isSaving = false;
              _isEditing = false;
              _selectedImagePath = null; // Reset selected image after save
            });
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Row(
                  children: [
                    const Icon(Icons.check_circle_outline, color: Colors.white),
                    const SizedBox(width: 12),
                    const Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          Text(
                            'Success',
                            style: TextStyle(fontWeight: FontWeight.bold),
                          ),
                          Text('Project updated successfully'),
                        ],
                      ),
                    ),
                  ],
                ),
                backgroundColor: Colors.green,
                duration: const Duration(seconds: 3),
              ),
            );
            // Reload the project to show updated data
            context.read<ProjectsBloc>().add(GetProjectByIdRequested(widget.projectId));
          } else if (state is ProjectDeleted) {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(
                content: Row(
                  children: [
                    const Icon(Icons.check_circle_outline, color: Colors.white),
                    const SizedBox(width: 12),
                    const Expanded(
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          Text(
                            'Success',
                            style: TextStyle(fontWeight: FontWeight.bold),
                          ),
                          Text('Project deleted successfully'),
                        ],
                      ),
                    ),
                  ],
                ),
                backgroundColor: Colors.green,
                duration: const Duration(seconds: 3),
              ),
            );
          } else if (state is ProjectsLoading) {
            setState(() {
              _isSaving = true;
            });
          }
        },
        builder: (context, state) {
          if (state is ProjectsLoading) {
            return const Center(child: CircularProgressIndicator());
          } else if (state is ProjectsError) {
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Icon(
                    Icons.error_outline,
                    size: 64,
                    color: Colors.red.shade300,
                  ),
                  const SizedBox(height: 16),
                  Text(
                    'Error loading project',
                    style: TextStyle(
                      fontSize: 18,
                      color: Colors.red.shade700,
                    ),
                  ),
                  const SizedBox(height: 8),
                  Text(
                    state.message,
                    style: TextStyle(color: Colors.red.shade600),
                    textAlign: TextAlign.center,
                  ),
                  const SizedBox(height: 16),
                  ElevatedButton(
                    onPressed: () {
                      context.read<ProjectsBloc>().add(GetProjectByIdRequested(widget.projectId));
                    },
                    child: const Text('Retry'),
                  ),
                ],
              ),
            );
          } else if (state is ProjectLoaded) {
            final project = state.project;
            
            // Populate controllers if not already done
            if (_titleController.text.isEmpty) {
              _titleController.text = project.title;
              _descriptionController.text = project.description;
              _longDescriptionController.text = project.longDescription ?? '';
              _contentController.text = project.content ?? '';
              _githubUrlController.text = project.githubUrl ?? '';
              _liveUrlController.text = project.liveUrl ?? '';
              _videoUrlController.text = project.videoUrl ?? '';
            }

            return Stack(
              children: [
                SingleChildScrollView(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                  // Project image
                  Container(
                    height: 250,
                    width: double.infinity,
                    decoration: BoxDecoration(
                      image: DecorationImage(
                        image: _selectedImagePath != null
                            ? FileImage(File(_selectedImagePath!))
                            : NetworkImage(project.image) as ImageProvider,
                        fit: BoxFit.cover,
                      ),
                    ),
                    child: Stack(
                      children: [
                        if (project.featured)
                          Positioned(
                            top: 16,
                            right: 16,
                            child: Container(
                              padding: const EdgeInsets.symmetric(
                                horizontal: 12,
                                vertical: 6,
                              ),
                              decoration: BoxDecoration(
                                color: Colors.amber.shade100,
                                borderRadius: BorderRadius.circular(16),
                              ),
                              child: Text(
                                'Featured',
                                style: TextStyle(
                                  fontSize: 12,
                                  color: Colors.amber.shade800,
                                  fontWeight: FontWeight.w500,
                                ),
                              ),
                            ),
                          ),
                        if (_isEditing)
                          Positioned(
                            bottom: 16,
                            right: 16,
                            child: FloatingActionButton.small(
                              onPressed: _pickImage,
                              backgroundColor: Colors.blue.shade600,
                              child: const Icon(Icons.camera_alt, color: Colors.white),
                            ),
                          ),
                      ],
                    ),
                  ),

                  Padding(
                    padding: const EdgeInsets.all(16),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        // Title Section
                        Row(
                          children: [
                            Expanded(
                              child: _isEditing
                                  ? TextField(
                                      controller: _titleController,
                                      decoration: InputDecoration(
                                        labelText: 'Title *',
                                        hintText: 'Enter project title',
                                        border: const OutlineInputBorder(),
                                        errorText: _titleError,
                                        errorMaxLines: 2,
                                        helperText: 'Required (3-200 characters)',
                                        helperStyle: TextStyle(
                                          fontSize: 12,
                                          color: Colors.grey.shade600,
                                        ),
                                      ),
                                      maxLength: 200,
                                      onChanged: (_) {
                                        // Clear error on change
                                        if (_titleError != null) {
                                          setState(() {
                                            _titleError = null;
                                          });
                                        }
                                      },
                                    )
                                  : Text(
                                      project.title,
                                      style: const TextStyle(
                                        fontSize: 28,
                                        fontWeight: FontWeight.bold,
                                      ),
                                    ),
                            ),
                          ],
                        ),
                        const SizedBox(height: 8),

                        // Status indicators
                        Row(
                          children: [
                            Icon(
                              Icons.visibility,
                              size: 16,
                              color: project.visible ? Colors.green : Colors.grey,
                            ),
                            const SizedBox(width: 4),
                            Text(
                              project.visible ? 'Visible' : 'Hidden',
                              style: TextStyle(
                                fontSize: 12,
                                color: project.visible ? Colors.green : Colors.grey,
                              ),
                            ),
                            const SizedBox(width: 16),
                            Icon(
                              Icons.public,
                              size: 16,
                              color: project.public ? Colors.blue : Colors.grey,
                            ),
                            const SizedBox(width: 4),
                            Text(
                              project.public ? 'Public' : 'Private',
                              style: TextStyle(
                                fontSize: 12,
                                color: project.public ? Colors.blue : Colors.grey,
                              ),
                            ),
                            const Spacer(),
                            Text(
                              _formatDate(project.createdAt),
                              style: TextStyle(
                                fontSize: 12,
                                color: Colors.grey.shade500,
                              ),
                            ),
                          ],
                        ),
                        const SizedBox(height: 16),

                        // Description Section
                        Card(
                          child: Padding(
                            padding: const EdgeInsets.all(16),
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                const Text(
                                  'Description',
                                  style: TextStyle(
                                    fontSize: 18,
                                    fontWeight: FontWeight.bold,
                                  ),
                                ),
                                const SizedBox(height: 8),
                                _isEditing
                                    ? TextField(
                                        controller: _descriptionController,
                                        decoration: InputDecoration(
                                          labelText: 'Description *',
                                          hintText: 'Enter project description',
                                          border: const OutlineInputBorder(),
                                          errorText: _descriptionError,
                                          errorMaxLines: 2,
                                          helperText: 'Required (10-500 characters)',
                                          helperStyle: TextStyle(
                                            fontSize: 12,
                                            color: Colors.grey.shade600,
                                          ),
                                        ),
                                        maxLines: 3,
                                        maxLength: 500,
                                        onChanged: (_) {
                                          // Clear error on change
                                          if (_descriptionError != null) {
                                            setState(() {
                                              _descriptionError = null;
                                            });
                                          }
                                        },
                                      )
                                    : Text(
                                        project.description,
                                        style: TextStyle(
                                          color: Colors.grey.shade700,
                                          fontSize: 16,
                                        ),
                                      ),
                              ],
                            ),
                          ),
                        ),
                        const SizedBox(height: 16),

                        // Long Description Section
                        if (project.longDescription != null || _isEditing) ...[
                          Card(
                            child: Padding(
                              padding: const EdgeInsets.all(16),
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  const Text(
                                    'Long Description',
                                    style: TextStyle(
                                      fontSize: 18,
                                      fontWeight: FontWeight.bold,
                                    ),
                                  ),
                                  const SizedBox(height: 8),
                                  _isEditing
                                      ? TextField(
                                          controller: _longDescriptionController,
                                          decoration: const InputDecoration(
                                            labelText: 'Long Description',
                                            border: OutlineInputBorder(),
                                          ),
                                          maxLines: 5,
                                        )
                                      : Text(
                                          project.longDescription ?? '',
                                          style: TextStyle(
                                            color: Colors.grey.shade700,
                                            fontSize: 16,
                                          ),
                                        ),
                                ],
                              ),
                            ),
                          ),
                          const SizedBox(height: 16),
                        ],

                        // Content Section
                        if (project.content != null || _isEditing) ...[
                          Card(
                            child: Padding(
                              padding: const EdgeInsets.all(16),
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  const Text(
                                    'Content',
                                    style: TextStyle(
                                      fontSize: 18,
                                      fontWeight: FontWeight.bold,
                                    ),
                                  ),
                                  const SizedBox(height: 8),
                                  _isEditing
                                      ? TextField(
                                          controller: _contentController,
                                          decoration: const InputDecoration(
                                            labelText: 'Content',
                                            border: OutlineInputBorder(),
                                          ),
                                          maxLines: 10,
                                        )
                                      : Container(
                                          width: double.infinity,
                                          padding: const EdgeInsets.all(12),
                                          decoration: BoxDecoration(
                                            color: Colors.grey.shade50,
                                            borderRadius: BorderRadius.circular(8),
                                            border: Border.all(color: Colors.grey.shade300),
                                          ),
                                          child: Text(
                                            project.content ?? '',
                                            style: TextStyle(
                                              color: Colors.grey.shade700,
                                              fontSize: 14,
                                            ),
                                          ),
                                        ),
                                ],
                              ),
                            ),
                          ),
                          const SizedBox(height: 16),
                        ],

                        // Frameworks Section
                        if (project.frameworks != null && project.frameworks!.isNotEmpty) ...[
                          Card(
                            child: Padding(
                              padding: const EdgeInsets.all(16),
                              child: Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  const Text(
                                    'Technologies',
                                    style: TextStyle(
                                      fontSize: 18,
                                      fontWeight: FontWeight.bold,
                                    ),
                                  ),
                                  const SizedBox(height: 12),
                                  Wrap(
                                    spacing: 8,
                                    runSpacing: 8,
                                    children: project.frameworks!.map((framework) {
                                      return Chip(
                                        label: Text(framework.name),
                                        backgroundColor: Colors.blue.shade50,
                                        labelStyle: TextStyle(
                                          color: Colors.blue.shade700,
                                          fontWeight: FontWeight.w500,
                                        ),
                                      );
                                    }).toList(),
                                  ),
                                ],
                              ),
                            ),
                          ),
                          const SizedBox(height: 16),
                        ],

                        // Links Section
                        Card(
                          child: Padding(
                            padding: const EdgeInsets.all(16),
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                const Text(
                                  'Links',
                                  style: TextStyle(
                                    fontSize: 18,
                                    fontWeight: FontWeight.bold,
                                  ),
                                ),
                                const SizedBox(height: 12),
                                if (project.githubUrl != null || _isEditing) ...[
                                  _buildLinkField(
                                    'GitHub',
                                    _githubUrlController,
                                    Icons.code,
                                    project.githubUrl,
                                    _githubUrlError,
                                  ),
                                  const SizedBox(height: 12),
                                ],
                                if (project.liveUrl != null || _isEditing) ...[
                                  _buildLinkField(
                                    'Live Demo',
                                    _liveUrlController,
                                    Icons.launch,
                                    project.liveUrl,
                                    _liveUrlError,
                                  ),
                                  const SizedBox(height: 12),
                                ],
                                if (project.videoUrl != null || _isEditing) ...[
                                  _buildLinkField(
                                    'Video',
                                    _videoUrlController,
                                    Icons.play_circle,
                                    project.videoUrl,
                                    _videoUrlError,
                                  ),
                                ],
                              ],
                            ),
                          ),
                        ),
                        const SizedBox(height: 16),

                        // Metadata Section
                        Card(
                          child: Padding(
                            padding: const EdgeInsets.all(16),
                            child: Column(
                              crossAxisAlignment: CrossAxisAlignment.start,
                              children: [
                                const Text(
                                  'Details',
                                  style: TextStyle(
                                    fontSize: 18,
                                    fontWeight: FontWeight.bold,
                                  ),
                                ),
                                const SizedBox(height: 12),
                                _buildMetadataRow('Slug', project.slug),
                                _buildMetadataRow('Order', project.order.toString()),
                                _buildMetadataRow('Created', _formatFullDate(project.createdAt)),
                                _buildMetadataRow('Updated', _formatFullDate(project.updatedAt)),
                              ],
                            ),
                          ),
                        ),
                        const SizedBox(height: 24),

                        // Action Buttons
                        if (_isEditing)
                          Row(
                            children: [
                              Expanded(
                                child: ElevatedButton.icon(
                                  onPressed: _isSaving ? null : () => _saveProject(),
                                  icon: _isSaving
                                      ? const SizedBox(
                                          width: 20,
                                          height: 20,
                                          child: CircularProgressIndicator(
                                            strokeWidth: 2,
                                            valueColor: AlwaysStoppedAnimation<Color>(Colors.white),
                                          ),
                                        )
                                      : const Icon(Icons.save),
                                  label: Text(_isSaving ? 'Saving...' : 'Save Changes'),
                                  style: ElevatedButton.styleFrom(
                                    backgroundColor: Colors.blue.shade600,
                                    foregroundColor: Colors.white,
                                    padding: const EdgeInsets.symmetric(vertical: 12),
                                  ),
                                ),
                              ),
                            ],
                          ),
                      ],
                    ),
                  ),
                    ],
                  ),
                ),
                // Loading Overlay
                if (_isSaving)
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
                                'Saving changes...',
                                style: TextStyle(
                                  fontSize: 16,
                                  fontWeight: FontWeight.w500,
                                ),
                              ),
                            ],
                          ),
                        ),
                      ),
                    ),
                  ),
              ],
            );
          }

          return const SizedBox.shrink();
        },
      ),
      floatingActionButton: !_isEditing
          ? FloatingActionButton.extended(
              onPressed: () => _showDeleteDialog(context),
              backgroundColor: Colors.red.shade600,
              icon: const Icon(Icons.delete, color: Colors.white),
              label: const Text('Delete', style: TextStyle(color: Colors.white)),
            )
          : null,
    );
  }

  Widget _buildLinkField(
    String label,
    TextEditingController controller,
    IconData icon,
    String? url,
    String? errorText,
  ) {
    return Row(
      children: [
        Icon(icon, size: 20, color: Colors.blue.shade600),
        const SizedBox(width: 8),
        Expanded(
          child: _isEditing
              ? TextField(
                  controller: controller,
                  decoration: InputDecoration(
                    labelText: label,
                    border: const OutlineInputBorder(),
                    hintText: 'https://...',
                    errorText: errorText,
                    errorMaxLines: 2,
                    helperText: 'Optional: Must be a valid URL',
                    helperStyle: TextStyle(
                      fontSize: 12,
                      color: Colors.grey.shade600,
                    ),
                  ),
                  onChanged: (_) {
                    // Clear error on change
                    if (errorText != null) {
                      setState(() {
                        if (label == 'GitHub') _githubUrlError = null;
                        if (label == 'Live Demo') _liveUrlError = null;
                        if (label == 'Video') _videoUrlError = null;
                      });
                    }
                  },
                )
              : url != null
                  ? InkWell(
                      onTap: () => _launchUrl(url),
                      child: Text(
                        url,
                        style: TextStyle(
                          color: Colors.blue.shade600,
                          decoration: TextDecoration.underline,
                        ),
                      ),
                    )
                  : Text(
                      'No $label provided',
                      style: TextStyle(color: Colors.grey.shade500),
                    ),
        ),
      ],
    );
  }

  Widget _buildMetadataRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 80,
            child: Text(
              '$label:',
              style: TextStyle(
                fontWeight: FontWeight.w500,
                color: Colors.grey.shade600,
              ),
            ),
          ),
          Expanded(
            child: Text(
              value,
              style: const TextStyle(
                fontWeight: FontWeight.w500,
              ),
            ),
          ),
        ],
      ),
    );
  }

  bool _validateForm() {
    bool isValid = true;
    setState(() {
      // Validate title
      if (_titleController.text.trim().isEmpty) {
        _titleError = 'Title is required';
        isValid = false;
      } else if (_titleController.text.trim().length < 3) {
        _titleError = 'Title must be at least 3 characters';
        isValid = false;
      } else if (_titleController.text.trim().length > 200) {
        _titleError = 'Title must be less than 200 characters';
        isValid = false;
      } else {
        _titleError = null;
      }

      // Validate description
      if (_descriptionController.text.trim().isEmpty) {
        _descriptionError = 'Description is required';
        isValid = false;
      } else if (_descriptionController.text.trim().length < 10) {
        _descriptionError = 'Description must be at least 10 characters';
        isValid = false;
      } else if (_descriptionController.text.trim().length > 500) {
        _descriptionError = 'Description must be less than 500 characters';
        isValid = false;
      } else {
        _descriptionError = null;
      }

      // Validate URLs if provided
      if (_githubUrlController.text.trim().isNotEmpty &&
          !_isValidUrl(_githubUrlController.text.trim())) {
        _githubUrlError = 'Invalid GitHub URL format';
        isValid = false;
      } else {
        _githubUrlError = null;
      }

      if (_liveUrlController.text.trim().isNotEmpty &&
          !_isValidUrl(_liveUrlController.text.trim())) {
        _liveUrlError = 'Invalid live URL format';
        isValid = false;
      } else {
        _liveUrlError = null;
      }

      if (_videoUrlController.text.trim().isNotEmpty &&
          !_isValidUrl(_videoUrlController.text.trim())) {
        _videoUrlError = 'Invalid video URL format';
        isValid = false;
      } else {
        _videoUrlError = null;
      }
    });

    if (!isValid) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Row(
            children: [
              Icon(Icons.warning_amber_outlined, color: Colors.white),
              SizedBox(width: 12),
              Expanded(
                child: Text('Please fix the validation errors before saving'),
              ),
            ],
          ),
          backgroundColor: Colors.orange,
          duration: Duration(seconds: 3),
        ),
      );
    }

    return isValid;
  }

  bool _isValidUrl(String url) {
    try {
      final uri = Uri.parse(url);
      return uri.hasScheme && (uri.scheme == 'http' || uri.scheme == 'https');
    } catch (e) {
      return false;
    }
  }

  void _saveProject() {
    if (!_validateForm()) {
      return;
    }

    setState(() {
      _isSaving = true;
    });

    context.read<ProjectsBloc>().add(UpdateProjectRequested(
      id: widget.projectId,
      title: _titleController.text.trim(),
      description: _descriptionController.text.trim(),
      longDescription: _longDescriptionController.text.trim().isEmpty
          ? null
          : _longDescriptionController.text.trim(),
      content: _contentController.text.trim().isEmpty
          ? null
          : _contentController.text.trim(),
      githubUrl: _githubUrlController.text.trim().isEmpty
          ? null
          : _githubUrlController.text.trim(),
      liveUrl: _liveUrlController.text.trim().isEmpty
          ? null
          : _liveUrlController.text.trim(),
      videoUrl: _videoUrlController.text.trim().isEmpty
          ? null
          : _videoUrlController.text.trim(),
    ));
  }

  void _showDeleteDialog(BuildContext context) {
    showDialog(
      context: context,
      barrierDismissible: false, // Prevent dismissing by tapping outside
      builder: (dialogContext) => AlertDialog(
        icon: Icon(
          Icons.warning_amber_rounded,
          color: Colors.red.shade600,
          size: 48,
        ),
        title: const Text(
          'Delete Project',
          style: TextStyle(
            fontWeight: FontWeight.bold,
          ),
        ),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const Text(
              'Are you sure you want to delete this project?',
              style: TextStyle(
                fontSize: 16,
                fontWeight: FontWeight.w500,
              ),
            ),
            const SizedBox(height: 12),
            Container(
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: Colors.red.shade50,
                borderRadius: BorderRadius.circular(8),
                border: Border.all(color: Colors.red.shade200),
              ),
              child: Row(
                children: [
                  Icon(Icons.info_outline, color: Colors.red.shade700, size: 20),
                  const SizedBox(width: 8),
                  const Expanded(
                    child: Text(
                      'This action cannot be undone. All project data will be permanently deleted.',
                      style: TextStyle(
                        fontSize: 14,
                        color: Colors.black87,
                      ),
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(dialogContext),
            child: const Text(
              'Cancel',
              style: TextStyle(fontSize: 16),
            ),
          ),
          ElevatedButton.icon(
            onPressed: () {
              Navigator.pop(dialogContext);
              // Show a brief confirmation message
              ScaffoldMessenger.of(context).showSnackBar(
                const SnackBar(
                  content: Row(
                    children: [
                      SizedBox(
                        width: 20,
                        height: 20,
                        child: CircularProgressIndicator(
                          strokeWidth: 2,
                          valueColor: AlwaysStoppedAnimation<Color>(Colors.white),
                        ),
                      ),
                      SizedBox(width: 12),
                      Text('Deleting project...'),
                    ],
                  ),
                  duration: Duration(seconds: 2),
                ),
              );
              context.read<ProjectsBloc>().add(DeleteProjectRequested(widget.projectId));
              context.pop(); // Go back to projects list
            },
            style: ElevatedButton.styleFrom(
              backgroundColor: Colors.red,
              foregroundColor: Colors.white,
              padding: const EdgeInsets.symmetric(horizontal: 20, vertical: 12),
            ),
            icon: const Icon(Icons.delete_forever),
            label: const Text(
              'Delete Permanently',
              style: TextStyle(
                fontSize: 16,
                fontWeight: FontWeight.w600,
              ),
            ),
          ),
        ],
      ),
    );
  }

  Future<void> _launchUrl(String url) async {
    try {
      final uri = Uri.parse(url);
      if (await canLaunchUrl(uri)) {
        await launchUrl(uri, mode: LaunchMode.externalApplication);
      } else {
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text('Could not launch URL: $url'),
              backgroundColor: Colors.red,
            ),
          );
        }
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Error launching URL: $e'),
            backgroundColor: Colors.red,
          ),
        );
      }
    }
  }

  Future<void> _pickImage() async {
    try {
      final XFile? image = await _imagePicker.pickImage(
        source: ImageSource.gallery,
        maxWidth: 1920,
        maxHeight: 1080,
        imageQuality: 85,
      );
      
      if (image != null) {
        setState(() {
          _selectedImagePath = image.path;
        });
        
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(
              content: Text('Image selected. Remember to save your changes.'),
              backgroundColor: Colors.green,
            ),
          );
        }
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Error picking image: $e'),
            backgroundColor: Colors.red,
          ),
        );
      }
    }
  }

  String _formatDate(DateTime date) {
    final now = DateTime.now();
    final difference = now.difference(date);

    if (difference.inDays == 0) {
      return 'Today';
    } else if (difference.inDays == 1) {
      return 'Yesterday';
    } else if (difference.inDays < 7) {
      return '${difference.inDays} days ago';
    } else {
      return '${date.day}/${date.month}/${date.year}';
    }
  }

  String _formatFullDate(DateTime date) {
    return '${date.day}/${date.month}/${date.year} at ${date.hour.toString().padLeft(2, '0')}:${date.minute.toString().padLeft(2, '0')}';
  }
}
