import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import '../bloc/education_bloc.dart';
import '../../domain/entities/education_entity.dart';

class EducationDetailPage extends StatefulWidget {
  final String educationId;

  const EducationDetailPage({
    super.key,
    required this.educationId,
  });

  @override
  State<EducationDetailPage> createState() => _EducationDetailPageState();
}

class _EducationDetailPageState extends State<EducationDetailPage> {
  EducationEntity? _lastLoadedEducation;

  @override
  void initState() {
    super.initState();
    // Load education data when the page initializes
    context.read<EducationBloc>().add(GetEducationById(widget.educationId));
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Education Details'),
        actions: [
          BlocBuilder<EducationBloc, EducationState>(
            builder: (context, state) {
              if (state is EducationDetailLoaded) {
                return IconButton(
                  icon: const Icon(Icons.edit),
                  onPressed: () {
                    context.push('/education/${widget.educationId}/edit');
                  },
                );
              }
              return const SizedBox.shrink();
            },
          ),
        ],
      ),
      body: BlocBuilder<EducationBloc, EducationState>(
        builder: (context, state) {
          if (state is EducationLoading) {
            return const Center(child: CircularProgressIndicator());
          }

          if (state is EducationError) {
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
                    state.message,
                    style: const TextStyle(fontSize: 16),
                    textAlign: TextAlign.center,
                  ),
                  const SizedBox(height: 16),
                  ElevatedButton(
                    onPressed: () {
                      context.read<EducationBloc>().add(GetEducationById(widget.educationId));
                    },
                    child: const Text('Retry'),
                  ),
                ],
              ),
            );
          }

          if (state is EducationDetailLoaded) {
            final education = state.education;
            _lastLoadedEducation = education;
            return _buildEducationDetail(education);
          }

          // Show loading overlay during delete operation
          if (state is EducationLoading) {
            return Stack(
              children: [
                _buildEducationDetail(_lastLoadedEducation ?? 
                  EducationEntity(
                    id: widget.educationId,
                    userId: '',
                    title: 'Loading...',
                    institution: '',
                    type: EducationType.degree,
                    order: 0,
                    visible: true,
                    createdAt: DateTime.now(),
                    updatedAt: DateTime.now(),
                  )),
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
                              'Processing...',
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
          }

          return const Center(child: Text('No education data available'));
        },
      ),
    );
  }

  Widget _buildEducationDetail(EducationEntity education) {
    return SingleChildScrollView(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Header Card
            Card(
              child: Padding(
                padding: const EdgeInsets.all(24),
                child: Column(
                  children: [
                    // Icon and Title
                    Row(
                      children: [
                        Container(
                          width: 80,
                          height: 80,
                          decoration: BoxDecoration(
                            color: _getTypeColor(education.type).withValues(alpha: 0.1),
                            borderRadius: BorderRadius.circular(16),
                          ),
                          child: Icon(
                            _getTypeIcon(education.type),
                            color: _getTypeColor(education.type),
                            size: 40,
                          ),
                        ),
                        const SizedBox(width: 16),
                        Expanded(
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              Text(
                                education.title,
                                style: const TextStyle(
                                  fontSize: 24,
                                  fontWeight: FontWeight.bold,
                                ),
                              ),
                              const SizedBox(height: 8),
                              Text(
                                education.institution,
                                style: const TextStyle(
                                  fontSize: 18,
                                  color: Colors.grey,
                                  fontWeight: FontWeight.w500,
                                ),
                              ),
                              const SizedBox(height: 8),
                              Row(
                                children: [
                                  Container(
                                    padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                                    decoration: BoxDecoration(
                                      color: _getTypeColor(education.type).withValues(alpha: 0.1),
                                      borderRadius: BorderRadius.circular(16),
                                    ),
                                    child: Text(
                                      education.type.name.toUpperCase(),
                                      style: TextStyle(
                                        fontSize: 12,
                                        fontWeight: FontWeight.bold,
                                        color: _getTypeColor(education.type),
                                      ),
                                    ),
                                  ),
                                  if (education.degreeLevel != null) ...[
                                    const SizedBox(width: 8),
                                    Container(
                                      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                                      decoration: BoxDecoration(
                                        color: Colors.blue.withValues(alpha: 0.1),
                                        borderRadius: BorderRadius.circular(16),
                                      ),
                                      child: Text(
                                        education.degreeLevel!.name.toUpperCase(),
                                        style: const TextStyle(
                                          fontSize: 12,
                                          fontWeight: FontWeight.bold,
                                          color: Colors.blue,
                                        ),
                                      ),
                                    ),
                                  ],
                                ],
                              ),
                            ],
                          ),
                        ),
                      ],
                    ),
                    const SizedBox(height: 16),

                    // Description
                    if (education.description != null)
                      Text(
                        education.description!,
                        style: const TextStyle(
                          fontSize: 16,
                          height: 1.5,
                        ),
                      ),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 16),

            // Education Details
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      'Education Details',
                      style: TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 16),
                    if (education.fieldOfStudy != null)
                      _buildDetailRow('Field of Study', education.fieldOfStudy!),
                    if (education.startDate != null)
                      _buildDetailRow('Start Date', _formatDate(education.startDate!)),
                    if (education.endDate != null)
                      _buildDetailRow('End Date', _formatDate(education.endDate!)),
                    if (education.completionDate != null)
                      _buildDetailRow('Completion Date', _formatDate(education.completionDate!)),
                    if (education.grade != null)
                      _buildDetailRow('Grade', education.grade!),
                    if (education.credits != null)
                      _buildDetailRow('Credits', education.credits.toString()),
                    _buildDetailRow('Order', education.order.toString()),
                    if (education.createdAt != null)
                      _buildDetailRow('Created', _formatDate(education.createdAt!)),
                    if (education.updatedAt != null)
                      _buildDetailRow('Updated', _formatDate(education.updatedAt!)),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 16),

            // Certificate Details
            if (education.certificateId != null || education.issuer != null || education.validityPeriod != null)
              Card(
                child: Padding(
                  padding: const EdgeInsets.all(16),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Text(
                        'Certificate Details',
                        style: TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      const SizedBox(height: 16),
                      if (education.certificateId != null)
                        _buildDetailRow('Certificate ID', education.certificateId!),
                      if (education.issuer != null)
                        _buildDetailRow('Issuer', education.issuer!),
                      if (education.validityPeriod != null)
                        _buildDetailRow('Validity Period', education.validityPeriod!),
                    ],
                  ),
                ),
              ),
            const SizedBox(height: 16),

            // Skills
            if (education.skills != null && education.skills!.isNotEmpty)
              Card(
                child: Padding(
                  padding: const EdgeInsets.all(16),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Text(
                        'Skills & Competencies',
                        style: TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      const SizedBox(height: 16),
                      Wrap(
                        spacing: 8,
                        runSpacing: 8,
                        children: education.skills!.map((skill) {
                          return Container(
                            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                            decoration: BoxDecoration(
                              color: Colors.blue.withValues(alpha: 0.1),
                              borderRadius: BorderRadius.circular(16),
                            ),
                            child: Text(
                              skill,
                              style: const TextStyle(
                                fontSize: 12,
                                fontWeight: FontWeight.w500,
                                color: Colors.blue,
                              ),
                            ),
                          );
                        }).toList(),
                      ),
                    ],
                  ),
                ),
              ),
            const SizedBox(height: 16),

            // Documents & Links
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      'Documents & Links',
                      style: TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 16),
                    if (education.pdfDocument != null)
                      ListTile(
                        leading: const Icon(Icons.picture_as_pdf, color: Colors.red),
                        title: const Text('Certificate PDF'),
                        subtitle: const Text('Download certificate'),
                        trailing: const Icon(Icons.download),
                        onTap: () {
                          // TODO: Implement PDF download
                          ScaffoldMessenger.of(context).showSnackBar(
                            const SnackBar(content: Text('PDF download coming soon!')),
                          );
                        },
                      ),
                    if (education.verificationUrl != null)
                      ListTile(
                        leading: const Icon(Icons.link, color: Colors.blue),
                        title: const Text('Verification Link'),
                        subtitle: const Text('Verify certificate online'),
                        trailing: const Icon(Icons.open_in_new),
                        onTap: () {
                          // TODO: Implement URL opening
                          ScaffoldMessenger.of(context).showSnackBar(
                            const SnackBar(content: Text('URL opening coming soon!')),
                          );
                        },
                      ),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 16),

            // Visibility Status
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      'Visibility',
                      style: TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 16),
                    Row(
                      children: [
                        Icon(
                          education.visible ? Icons.visibility : Icons.visibility_off,
                          color: education.visible ? Colors.green : Colors.grey,
                        ),
                        const SizedBox(width: 8),
                        Text(
                          education.visible ? 'Visible' : 'Hidden',
                          style: TextStyle(
                            color: education.visible ? Colors.green : Colors.grey,
                            fontWeight: FontWeight.w500,
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
            ),
            const SizedBox(height: 16),

            // Actions
            Card(
              child: Padding(
                padding: const EdgeInsets.all(16),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      'Actions',
                      style: TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 16),
                    Row(
                      children: [
                        Expanded(
                          child: OutlinedButton.icon(
                            onPressed: () {
                              context.push('/education/${education.id}/edit');
                            },
                            icon: const Icon(Icons.edit),
                            label: const Text('Edit'),
                          ),
                        ),
                        const SizedBox(width: 12),
                        Expanded(
                          child: BlocListener<EducationBloc, EducationState>(
                            listener: (context, state) {
                              if (state is EducationDeleted) {
                                ScaffoldMessenger.of(context).showSnackBar(
                                  const SnackBar(
                                    content: Row(
                                      children: [
                                        Icon(Icons.check_circle, color: Colors.white),
                                        SizedBox(width: 12),
                                        Expanded(child: Text('Education deleted successfully')),
                                      ],
                                    ),
                                    backgroundColor: Colors.green,
                                    behavior: SnackBarBehavior.floating,
                                    duration: Duration(seconds: 3),
                                  ),
                                );
                                context.pop();
                              } else if (state is EducationError) {
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
                                            Text('Failed to delete education', 
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
                                  ),
                                );
                              }
                            },
                            child: ElevatedButton.icon(
                              onPressed: () {
                                _showDeleteConfirmation(context, education);
                              },
                              icon: const Icon(Icons.delete),
                              label: const Text('Delete'),
                              style: ElevatedButton.styleFrom(
                                backgroundColor: Colors.red,
                                foregroundColor: Colors.white,
                              ),
                            ),
                          ),
                        ),
                      ],
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
    );
  }

  Widget _buildDetailRow(String label, String value) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          SizedBox(
            width: 120,
            child: Text(
              label,
              style: const TextStyle(
                fontWeight: FontWeight.w500,
                color: Colors.grey,
              ),
            ),
          ),
          Expanded(
            child: Text(
              value,
              style: const TextStyle(fontWeight: FontWeight.w500),
            ),
          ),
        ],
      ),
    );
  }

  String _formatDate(DateTime date) {
    return '${date.day}/${date.month}/${date.year}';
  }

  Color _getTypeColor(EducationType type) {
    switch (type) {
      case EducationType.degree:
        return Colors.blue;
      case EducationType.certificate:
        return Colors.green;
    }
  }

  IconData _getTypeIcon(EducationType type) {
    switch (type) {
      case EducationType.degree:
        return Icons.school;
      case EducationType.certificate:
        return Icons.verified;
    }
  }

  void _showDeleteConfirmation(BuildContext context, EducationEntity education) {
    showDialog(
      context: context,
      builder: (BuildContext context) {
        return AlertDialog(
          title: const Text('Delete Education'),
          content: Text('Are you sure you want to delete "${education.title}"? This action cannot be undone.'),
          actions: [
            TextButton(
              onPressed: () => Navigator.of(context).pop(),
              child: const Text('Cancel'),
            ),
            ElevatedButton(
              onPressed: () {
                Navigator.of(context).pop();
                context.read<EducationBloc>().add(DeleteEducation(education.id));
              },
              style: ElevatedButton.styleFrom(
                backgroundColor: Colors.red,
                foregroundColor: Colors.white,
              ),
              child: const Text('Delete'),
            ),
          ],
        );
      },
    );
  }
}
