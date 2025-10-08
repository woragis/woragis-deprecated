import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import '../bloc/education_bloc.dart';
import '../../domain/entities/education_entity.dart';

class CreateEducationPage extends StatefulWidget {
  final String? educationId;

  const CreateEducationPage({
    super.key,
    this.educationId,
  });

  @override
  State<CreateEducationPage> createState() => _CreateEducationPageState();
}

class _CreateEducationPageState extends State<CreateEducationPage> {
  final _formKey = GlobalKey<FormState>();
  final _titleController = TextEditingController();
  final _institutionController = TextEditingController();
  final _descriptionController = TextEditingController();
  final _fieldOfStudyController = TextEditingController();
  final _gradeController = TextEditingController();
  final _creditsController = TextEditingController();
  final _certificateIdController = TextEditingController();
  final _issuerController = TextEditingController();
  final _validityPeriodController = TextEditingController();
  final _pdfDocumentController = TextEditingController();
  final _verificationUrlController = TextEditingController();
  final _skillsController = TextEditingController();
  
  // Auto-validate after first submit attempt
  AutovalidateMode _autovalidateMode = AutovalidateMode.disabled;

  EducationType _selectedType = EducationType.degree;
  DegreeLevel? _selectedDegreeLevel;
  DateTime? _startDate;
  DateTime? _endDate;
  DateTime? _completionDate;
  bool _isVisible = true;
  int _order = 0;
  List<String> _skills = [];

  bool get _isEditing => widget.educationId != null;

  @override
  void initState() {
    super.initState();
    if (_isEditing) {
      _loadEducationData();
    }
  }

  void _loadEducationData() {
    // Load education data from BLoC state if available
    final state = context.read<EducationBloc>().state;
    if (state is EducationDetailLoaded) {
      final education = state.education;
      _titleController.text = education.title;
      _institutionController.text = education.institution;
      _descriptionController.text = education.description ?? '';
      _fieldOfStudyController.text = education.fieldOfStudy ?? '';
      _gradeController.text = education.grade ?? '';
      _creditsController.text = education.credits?.toString() ?? '';
      _certificateIdController.text = education.certificateId ?? '';
      _issuerController.text = education.issuer ?? '';
      _validityPeriodController.text = education.validityPeriod ?? '';
      _pdfDocumentController.text = education.pdfDocument ?? '';
      _verificationUrlController.text = education.verificationUrl ?? '';
      _skillsController.text = education.skills?.join(', ') ?? '';
      _selectedType = education.type;
      _selectedDegreeLevel = education.degreeLevel;
      _startDate = education.startDate;
      _endDate = education.endDate;
      _completionDate = education.completionDate;
      _isVisible = education.visible;
      _order = education.order;
      _skills = education.skills ?? [];
    }
  }

  @override
  void dispose() {
    _titleController.dispose();
    _institutionController.dispose();
    _descriptionController.dispose();
    _fieldOfStudyController.dispose();
    _gradeController.dispose();
    _creditsController.dispose();
    _certificateIdController.dispose();
    _issuerController.dispose();
    _validityPeriodController.dispose();
    _pdfDocumentController.dispose();
    _verificationUrlController.dispose();
    _skillsController.dispose();
    super.dispose();
  }

  void _parseSkills() {
    final skillsText = _skillsController.text.trim();
    if (skillsText.isNotEmpty) {
      _skills = skillsText.split(',').map((skill) => skill.trim()).where((skill) => skill.isNotEmpty).toList();
    } else {
      _skills = [];
    }
    setState(() {}); // Trigger rebuild to show parsed skills
  }

  void _saveEducation() {
    // Enable auto-validate after first submit attempt
    setState(() {
      _autovalidateMode = AutovalidateMode.onUserInteraction;
    });
    
    if (_formKey.currentState!.validate()) {
      _parseSkills();
      
      if (_isEditing) {
        // Update existing education
        context.read<EducationBloc>().add(UpdateEducation(
          id: widget.educationId!,
          title: _titleController.text,
          institution: _institutionController.text,
          description: _descriptionController.text.isEmpty ? null : _descriptionController.text,
          type: _selectedType.name,
          degreeLevel: _selectedDegreeLevel?.name,
          fieldOfStudy: _fieldOfStudyController.text.isEmpty ? null : _fieldOfStudyController.text,
          startDate: _startDate,
          endDate: _endDate,
          completionDate: _completionDate,
          grade: _gradeController.text.isEmpty ? null : _gradeController.text,
          credits: _creditsController.text.isEmpty ? null : int.tryParse(_creditsController.text),
          certificateId: _certificateIdController.text.isEmpty ? null : _certificateIdController.text,
          issuer: _issuerController.text.isEmpty ? null : _issuerController.text,
          validityPeriod: _validityPeriodController.text.isEmpty ? null : _validityPeriodController.text,
          pdfDocument: _pdfDocumentController.text.isEmpty ? null : _pdfDocumentController.text,
          verificationUrl: _verificationUrlController.text.isEmpty ? null : _verificationUrlController.text,
          skills: _skills.isEmpty ? null : _skills,
          order: _order,
          visible: _isVisible,
        ));
      } else {
        // Create new education
        context.read<EducationBloc>().add(CreateEducation(
          title: _titleController.text,
          institution: _institutionController.text,
          description: _descriptionController.text.isEmpty ? null : _descriptionController.text,
          type: _selectedType.name,
          degreeLevel: _selectedDegreeLevel?.name,
          fieldOfStudy: _fieldOfStudyController.text.isEmpty ? null : _fieldOfStudyController.text,
          startDate: _startDate,
          endDate: _endDate,
          completionDate: _completionDate,
          grade: _gradeController.text.isEmpty ? null : _gradeController.text,
          credits: _creditsController.text.isEmpty ? null : int.tryParse(_creditsController.text),
          certificateId: _certificateIdController.text.isEmpty ? null : _certificateIdController.text,
          issuer: _issuerController.text.isEmpty ? null : _issuerController.text,
          validityPeriod: _validityPeriodController.text.isEmpty ? null : _validityPeriodController.text,
          pdfDocument: _pdfDocumentController.text.isEmpty ? null : _pdfDocumentController.text,
          verificationUrl: _verificationUrlController.text.isEmpty ? null : _verificationUrlController.text,
          skills: _skills.isEmpty ? null : _skills,
          order: _order,
          visible: _isVisible,
        ));
      }
    }
  }

  Future<void> _selectDate(BuildContext context, String field) async {
    final DateTime? picked = await showDatePicker(
      context: context,
      initialDate: DateTime.now(),
      firstDate: DateTime(1900),
      lastDate: DateTime.now().add(const Duration(days: 365 * 10)),
    );
    if (picked != null) {
      setState(() {
        switch (field) {
          case 'start':
            _startDate = picked;
            break;
          case 'end':
            _endDate = picked;
            break;
          case 'completion':
            _completionDate = picked;
            break;
        }
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(_isEditing ? 'Edit Education' : 'Add Education'),
        actions: [
          BlocBuilder<EducationBloc, EducationState>(
            builder: (context, state) {
              final isLoading = state is EducationLoading;
              return TextButton(
                onPressed: isLoading ? null : _saveEducation,
                child: isLoading 
                  ? const SizedBox(
                      width: 16,
                      height: 16,
                      child: CircularProgressIndicator(strokeWidth: 2),
                    )
                  : const Text('Save'),
              );
            },
          ),
        ],
      ),
      body: BlocConsumer<EducationBloc, EducationState>(
        listener: (context, state) {
          if (state is EducationCreated) {
            ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(
                content: Row(
                  children: [
                    Icon(Icons.check_circle, color: Colors.white),
                    SizedBox(width: 12),
                    Expanded(child: Text('Education created successfully!')),
                  ],
                ),
                backgroundColor: Colors.green,
                behavior: SnackBarBehavior.floating,
                duration: Duration(seconds: 3),
              ),
            );
            context.pop();
          } else if (state is EducationUpdated) {
            ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(
                content: Row(
                  children: [
                    Icon(Icons.check_circle, color: Colors.white),
                    SizedBox(width: 12),
                    Expanded(child: Text('Education updated successfully!')),
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
                        Text('Failed to save education', style: TextStyle(fontWeight: FontWeight.bold)),
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
                  onPressed: _saveEducation,
                ),
              ),
            );
          }
        },
        builder: (context, state) {
          final isLoading = state is EducationLoading;
          return Stack(
            children: [
              Form(
        key: _formKey,
        autovalidateMode: _autovalidateMode,
        child: SingleChildScrollView(
          padding: const EdgeInsets.all(16),
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

                      // Title
                      TextFormField(
                        controller: _titleController,
                        decoration: const InputDecoration(
                          labelText: 'Title *',
                          hintText: 'e.g., Bachelor of Computer Science',
                        ),
                        validator: (value) {
                          if (value == null || value.isEmpty) {
                            return 'Please enter a title';
                          }
                          return null;
                        },
                      ),
                      const SizedBox(height: 16),

                      // Institution
                      TextFormField(
                        controller: _institutionController,
                        decoration: const InputDecoration(
                          labelText: 'Institution *',
                          hintText: 'e.g., University of Technology',
                        ),
                        validator: (value) {
                          if (value == null || value.isEmpty) {
                            return 'Please enter an institution';
                          }
                          return null;
                        },
                      ),
                      const SizedBox(height: 16),

                      // Description
                      TextFormField(
                        controller: _descriptionController,
                        decoration: const InputDecoration(
                          labelText: 'Description',
                          hintText: 'Brief description of the education program',
                        ),
                        maxLines: 3,
                      ),
                    ],
                  ),
                ),
              ),
              const SizedBox(height: 16),

              // Education Type and Level
              Card(
                child: Padding(
                  padding: const EdgeInsets.all(16),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Text(
                        'Education Type',
                        style: TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      const SizedBox(height: 16),

                      // Type
                      DropdownButtonFormField<EducationType>(
                        value: _selectedType,
                        decoration: const InputDecoration(
                          labelText: 'Type *',
                        ),
                        items: EducationType.values.map((type) {
                          return DropdownMenuItem(
                            value: type,
                            child: Text(type.name.toUpperCase()),
                          );
                        }).toList(),
                        onChanged: (value) {
                          setState(() {
                            _selectedType = value!;
                            if (value != EducationType.degree) {
                              _selectedDegreeLevel = null;
                            }
                          });
                        },
                      ),
                      const SizedBox(height: 16),

                      // Degree Level (only for degrees)
                      if (_selectedType == EducationType.degree)
                        DropdownButtonFormField<DegreeLevel?>(
                          value: _selectedDegreeLevel,
                          decoration: const InputDecoration(
                            labelText: 'Degree Level *',
                            helperText: 'Required for degree type',
                          ),
                          validator: (value) {
                            if (_selectedType == EducationType.degree && value == null) {
                              return 'Please select a degree level';
                            }
                            return null;
                          },
                          items: [
                            const DropdownMenuItem(
                              value: null,
                              child: Text('Select degree level'),
                            ),
                            ...DegreeLevel.values.map((level) {
                              return DropdownMenuItem(
                                value: level,
                                child: Text(level.name.toUpperCase()),
                              );
                            }),
                          ],
                          onChanged: (value) {
                            setState(() {
                              _selectedDegreeLevel = value;
                            });
                          },
                        ),
                    ],
                  ),
                ),
              ),
              const SizedBox(height: 16),

              // Academic Details
              Card(
                child: Padding(
                  padding: const EdgeInsets.all(16),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Text(
                        'Academic Details',
                        style: TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      const SizedBox(height: 16),

                      // Field of Study
                      TextFormField(
                        controller: _fieldOfStudyController,
                        decoration: const InputDecoration(
                          labelText: 'Field of Study',
                          hintText: 'e.g., Computer Science, Business Administration',
                        ),
                      ),
                      const SizedBox(height: 16),

                      // Grade
                      TextFormField(
                        controller: _gradeController,
                        decoration: const InputDecoration(
                          labelText: 'Grade/GPA',
                          hintText: 'e.g., 3.8/4.0, A+, 85%',
                        ),
                      ),
                      const SizedBox(height: 16),

                      // Credits
                      TextFormField(
                        controller: _creditsController,
                        decoration: const InputDecoration(
                          labelText: 'Credits',
                          hintText: 'e.g., 120',
                        ),
                        keyboardType: TextInputType.number,
                        inputFormatters: [FilteringTextInputFormatter.digitsOnly],
                        validator: (value) {
                          if (value != null && value.isNotEmpty) {
                            final credits = int.tryParse(value);
                            if (credits == null || credits < 0) {
                              return 'Please enter a valid number of credits';
                            }
                          }
                          return null;
                        },
                      ),
                    ],
                  ),
                ),
              ),
              const SizedBox(height: 16),

              // Dates
              Card(
                child: Padding(
                  padding: const EdgeInsets.all(16),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Text(
                        'Important Dates',
                        style: TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      const SizedBox(height: 16),

                      // Start Date
                      ListTile(
                        title: const Text('Start Date'),
                        subtitle: Text(_startDate != null ? _formatDate(_startDate!) : 'Select start date'),
                        trailing: const Icon(Icons.calendar_today),
                        onTap: () => _selectDate(context, 'start'),
                      ),

                      // End Date
                      ListTile(
                        title: const Text('End Date'),
                        subtitle: Text(_endDate != null ? _formatDate(_endDate!) : 'Select end date'),
                        trailing: const Icon(Icons.calendar_today),
                        onTap: () => _selectDate(context, 'end'),
                      ),

                      // Completion Date
                      ListTile(
                        title: const Text('Completion Date'),
                        subtitle: Text(_completionDate != null ? _formatDate(_completionDate!) : 'Select completion date'),
                        trailing: const Icon(Icons.calendar_today),
                        onTap: () => _selectDate(context, 'completion'),
                      ),
                    ],
                  ),
                ),
              ),
              const SizedBox(height: 16),

              // Certificate Details
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

                      // Certificate ID
                      TextFormField(
                        controller: _certificateIdController,
                        decoration: const InputDecoration(
                          labelText: 'Certificate ID',
                          hintText: 'e.g., CS-BS-2024-001',
                        ),
                      ),
                      const SizedBox(height: 16),

                      // Issuer
                      TextFormField(
                        controller: _issuerController,
                        decoration: const InputDecoration(
                          labelText: 'Issuer',
                          hintText: 'e.g., University of Technology',
                        ),
                      ),
                      const SizedBox(height: 16),

                      // Validity Period
                      TextFormField(
                        controller: _validityPeriodController,
                        decoration: const InputDecoration(
                          labelText: 'Validity Period',
                          hintText: 'e.g., Lifetime, 2 years',
                        ),
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

                      // PDF Document
                      TextFormField(
                        controller: _pdfDocumentController,
                        decoration: const InputDecoration(
                          labelText: 'PDF Document URL',
                          hintText: 'e.g., /documents/degree.pdf',
                        ),
                      ),
                      const SizedBox(height: 16),

                      // Verification URL
                      TextFormField(
                        controller: _verificationUrlController,
                        decoration: const InputDecoration(
                          labelText: 'Verification URL',
                          hintText: 'e.g., https://verify.university.edu/certificate',
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
                      ),
                    ],
                  ),
                ),
              ),
              const SizedBox(height: 16),

              // Skills
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

                      TextFormField(
                        controller: _skillsController,
                        decoration: const InputDecoration(
                          labelText: 'Skills',
                          hintText: 'e.g., Programming, Algorithms, Data Structures',
                          helperText: 'Separate skills with commas',
                        ),
                        maxLines: 2,
                        onChanged: (value) => _parseSkills(),
                      ),
                      if (_skills.isNotEmpty) ...[
                        const SizedBox(height: 12),
                        Wrap(
                          spacing: 8,
                          runSpacing: 8,
                          children: _skills.map((skill) {
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
                    ],
                  ),
                ),
              ),
              const SizedBox(height: 16),

              // Settings
              Card(
                child: Padding(
                  padding: const EdgeInsets.all(16),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      const Text(
                        'Settings',
                        style: TextStyle(
                          fontSize: 18,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      const SizedBox(height: 16),

                      // Order
                      TextFormField(
                        initialValue: _order.toString(),
                        decoration: const InputDecoration(
                          labelText: 'Order *',
                          helperText: 'Display order (lower numbers appear first)',
                        ),
                        keyboardType: TextInputType.number,
                        inputFormatters: [FilteringTextInputFormatter.digitsOnly],
                        validator: (value) {
                          if (value == null || value.isEmpty) {
                            return 'Please enter an order number';
                          }
                          final order = int.tryParse(value);
                          if (order == null || order < 0) {
                            return 'Please enter a valid order number';
                          }
                          return null;
                        },
                        onChanged: (value) {
                          _order = int.tryParse(value) ?? 0;
                        },
                      ),
                      const SizedBox(height: 16),

                      // Visibility
                      SwitchListTile(
                        title: const Text('Visible'),
                        subtitle: const Text('Show this education in your portfolio'),
                        value: _isVisible,
                        onChanged: (value) {
                          setState(() {
                            _isVisible = value;
                          });
                        },
                      ),
                    ],
                  ),
                ),
              ),
              const SizedBox(height: 24),

              // Save Button
              BlocBuilder<EducationBloc, EducationState>(
                builder: (context, state) {
                  final isLoading = state is EducationLoading;
                  return SizedBox(
                    width: double.infinity,
                    child: ElevatedButton(
                      onPressed: isLoading ? null : _saveEducation,
                      style: ElevatedButton.styleFrom(
                        padding: const EdgeInsets.symmetric(vertical: 16),
                      ),
                      child: isLoading
                        ? const Row(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              SizedBox(
                                width: 20,
                                height: 20,
                                child: CircularProgressIndicator(strokeWidth: 2),
                              ),
                              SizedBox(width: 12),
                              Text('Saving...'),
                            ],
                          )
                        : Text(
                            _isEditing ? 'Update Education' : 'Add Education',
                            style: const TextStyle(fontSize: 16),
                          ),
                    ),
                  );
                },
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
                              'Saving education...',
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

  String _formatDate(DateTime date) {
    return '${date.day}/${date.month}/${date.year}';
  }
}
