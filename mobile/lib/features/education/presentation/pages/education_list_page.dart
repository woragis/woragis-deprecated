import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import '../bloc/education_bloc.dart';
import '../../domain/entities/education_entity.dart';

class EducationListPage extends StatefulWidget {
  const EducationListPage({super.key});

  @override
  State<EducationListPage> createState() => _EducationListPageState();
}

class _EducationListPageState extends State<EducationListPage> {
  final TextEditingController _searchController = TextEditingController();
  String? _selectedType;
  String? _selectedInstitution;
  bool? _selectedVisible;
  String _sortBy = 'order';
  String _sortOrder = 'asc';
  bool _isGridView = true;

  @override
  void initState() {
    super.initState();
    context.read<EducationBloc>().add(const LoadEducationList());
  }

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  void _applyFilters() {
    context.read<EducationBloc>().add(LoadEducationList(
      page: 1,
      limit: 20,
      search: _searchController.text.isEmpty ? null : _searchController.text,
      type: _selectedType,
      institution: _selectedInstitution,
      visible: _selectedVisible,
      sortBy: _sortBy,
      sortOrder: _sortOrder,
    ));
  }

  void _clearFilters() {
    setState(() {
      _searchController.clear();
      _selectedType = null;
      _selectedInstitution = null;
      _selectedVisible = null;
      _sortBy = 'order';
      _sortOrder = 'asc';
    });
    context.read<EducationBloc>().add(const LoadEducationList());
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Education & Certifications'),
        actions: [
          IconButton(
            icon: Icon(_isGridView ? Icons.list : Icons.grid_view),
            onPressed: () {
              setState(() {
                _isGridView = !_isGridView;
              });
            },
          ),
          IconButton(
            icon: const Icon(Icons.add),
            onPressed: () {
              context.push('/education/create');
            },
          ),
        ],
      ),
      body: Column(
        children: [
          // Search and Filters
          Padding(
            padding: const EdgeInsets.all(16.0),
            child: Column(
              children: [
                // Search Bar
                TextField(
                  controller: _searchController,
                  decoration: InputDecoration(
                    hintText: 'Search education...',
                    prefixIcon: const Icon(Icons.search),
                    suffixIcon: _searchController.text.isNotEmpty
                        ? IconButton(
                            icon: const Icon(Icons.clear),
                            onPressed: () {
                              _searchController.clear();
                              _applyFilters();
                            },
                          )
                        : null,
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                  ),
                  onChanged: (value) {
                    setState(() {});
                    if (value.isEmpty) {
                      _applyFilters();
                    }
                  },
                  onSubmitted: (_) => _applyFilters(),
                ),
                const SizedBox(height: 12),

                // Filter Chips
                SingleChildScrollView(
                  scrollDirection: Axis.horizontal,
                  child: Row(
                    children: [
                      // Type Filter
                      FilterChip(
                        label: Text(_selectedType ?? 'All Types'),
                        selected: _selectedType != null,
                        onSelected: (selected) {
                          setState(() {
                            _selectedType = selected ? 'degree' : null;
                          });
                          _applyFilters();
                        },
                      ),
                      const SizedBox(width: 8),

                      // Institution Filter
                      FilterChip(
                        label: Text(_selectedInstitution ?? 'All Institutions'),
                        selected: _selectedInstitution != null,
                        onSelected: (selected) {
                          setState(() {
                            _selectedInstitution = selected ? 'University' : null;
                          });
                          _applyFilters();
                        },
                      ),
                      const SizedBox(width: 8),

                      // Visible Filter
                      FilterChip(
                        label: Text(_selectedVisible == true ? 'Visible' : 'All'),
                        selected: _selectedVisible == true,
                        onSelected: (selected) {
                          setState(() {
                            _selectedVisible = selected ? true : null;
                          });
                          _applyFilters();
                        },
                      ),
                      const SizedBox(width: 8),

                      // Sort Options
                      PopupMenuButton<String>(
                        icon: const Icon(Icons.sort),
                        onSelected: (value) {
                          setState(() {
                            _sortBy = value;
                          });
                          _applyFilters();
                        },
                        itemBuilder: (context) => [
                          const PopupMenuItem(
                            value: 'order',
                            child: Text('Order'),
                          ),
                          const PopupMenuItem(
                            value: 'title',
                            child: Text('Title'),
                          ),
                          const PopupMenuItem(
                            value: 'institution',
                            child: Text('Institution'),
                          ),
                          const PopupMenuItem(
                            value: 'startDate',
                            child: Text('Start Date'),
                          ),
                          const PopupMenuItem(
                            value: 'endDate',
                            child: Text('End Date'),
                          ),
                        ],
                      ),
                      const SizedBox(width: 8),

                      // Clear Filters
                      if (_selectedType != null ||
                          _selectedInstitution != null ||
                          _selectedVisible != null ||
                          _searchController.text.isNotEmpty)
                        ActionChip(
                          label: const Text('Clear'),
                          onPressed: _clearFilters,
                          backgroundColor: Colors.red.shade50,
                          labelStyle: TextStyle(color: Colors.red.shade700),
                        ),
                    ],
                  ),
                ),
              ],
            ),
          ),

          // Education List
          Expanded(
            child: BlocBuilder<EducationBloc, EducationState>(
              builder: (context, state) {
                if (state is EducationLoading) {
                  return const Center(child: CircularProgressIndicator());
                }

                if (state is EducationError) {
                  return SingleChildScrollView(
                    padding: const EdgeInsets.all(16),
                    child: Center(
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
                            maxLines: 10,
                            overflow: TextOverflow.ellipsis,
                          ),
                          const SizedBox(height: 16),
                          ElevatedButton(
                            onPressed: () {
                              context.read<EducationBloc>().add(const RefreshEducationList());
                            },
                            child: const Text('Retry'),
                          ),
                        ],
                      ),
                    ),
                  );
                }

                if (state is EducationLoaded) {
                  if (state.educationList.isEmpty) {
                    return SingleChildScrollView(
                      padding: const EdgeInsets.all(16),
                      child: Center(
                        child: Column(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            Icon(
                              Icons.school_outlined,
                              size: 64,
                              color: Colors.grey.shade400,
                            ),
                            const SizedBox(height: 16),
                            const Text(
                              'No education records found',
                              style: TextStyle(fontSize: 18, fontWeight: FontWeight.w500),
                              textAlign: TextAlign.center,
                            ),
                            const SizedBox(height: 8),
                            const Text(
                              'Add your first education or certification',
                              style: TextStyle(color: Colors.grey),
                              textAlign: TextAlign.center,
                            ),
                            const SizedBox(height: 24),
                            ElevatedButton.icon(
                              onPressed: () {
                                context.push('/education/create');
                              },
                              icon: const Icon(Icons.add),
                              label: const Text('Add Education'),
                            ),
                          ],
                        ),
                      ),
                    );
                  }

                  return RefreshIndicator(
                    onRefresh: () async {
                      context.read<EducationBloc>().add(const RefreshEducationList());
                    },
                    child: _isGridView
                        ? _buildGridView(state.educationList)
                        : _buildListView(state.educationList),
                  );
                }

                return const Center(child: CircularProgressIndicator());
              },
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildGridView(List<EducationEntity> educationList) {
    return GridView.builder(
      padding: const EdgeInsets.all(16),
      gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: 2,
        childAspectRatio: 1.1,
        crossAxisSpacing: 16,
        mainAxisSpacing: 16,
      ),
      itemCount: educationList.length,
      itemBuilder: (context, index) {
        final education = educationList[index];
        return _buildEducationCard(education);
      },
    );
  }

  Widget _buildListView(List<EducationEntity> educationList) {
    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: educationList.length,
      itemBuilder: (context, index) {
        final education = educationList[index];
        return _buildEducationListItem(education);
      },
    );
  }

  Widget _buildEducationCard(EducationEntity education) {
    return Card(
      child: InkWell(
        onTap: () {
          context.push('/education/${education.id}');
        },
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Type and Institution
              Row(
                children: [
                  Container(
                    width: 40,
                    height: 40,
                    decoration: BoxDecoration(
                      color: _getTypeColor(education.type).withValues(alpha: 0.1),
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: Icon(
                      _getTypeIcon(education.type),
                      color: _getTypeColor(education.type),
                      size: 20,
                    ),
                  ),
                  const Spacer(),
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                    decoration: BoxDecoration(
                      color: _getTypeColor(education.type).withValues(alpha: 0.1),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Text(
                      education.type.name.toUpperCase(),
                      style: TextStyle(
                        fontSize: 10,
                        fontWeight: FontWeight.bold,
                        color: _getTypeColor(education.type),
                      ),
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 12),

              // Title
              Text(
                education.title,
                style: const TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.w600,
                ),
                maxLines: 2,
                overflow: TextOverflow.ellipsis,
              ),
              const SizedBox(height: 4),

              // Institution
              Text(
                education.institution,
                style: const TextStyle(
                  fontSize: 12,
                  color: Colors.grey,
                  fontWeight: FontWeight.w500,
                ),
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
              ),

              const Spacer(),

              // Date and Status
              Row(
                children: [
                  if (education.startDate != null)
                    Icon(
                      Icons.calendar_today,
                      size: 12,
                      color: Colors.grey.shade600,
                    ),
                  if (education.startDate != null) ...[
                    const SizedBox(width: 4),
                    Text(
                      _formatDate(education.startDate!),
                      style: TextStyle(
                        fontSize: 10,
                        color: Colors.grey.shade600,
                      ),
                    ),
                  ],
                  const Spacer(),
                  if (education.visible)
                    Icon(
                      Icons.visibility,
                      size: 16,
                      color: Colors.green.shade600,
                    ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildEducationListItem(EducationEntity education) {
    return Card(
      margin: const EdgeInsets.only(bottom: 8),
      child: ListTile(
        leading: Container(
          width: 48,
          height: 48,
          decoration: BoxDecoration(
            color: _getTypeColor(education.type).withValues(alpha: 0.1),
            borderRadius: BorderRadius.circular(8),
          ),
          child: Icon(
            _getTypeIcon(education.type),
            color: _getTypeColor(education.type),
            size: 24,
          ),
        ),
        title: Text(
          education.title,
          style: const TextStyle(fontWeight: FontWeight.w600),
        ),
        subtitle: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              education.institution,
              style: const TextStyle(fontWeight: FontWeight.w500),
            ),
            const SizedBox(height: 4),
            Row(
              children: [
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                  decoration: BoxDecoration(
                    color: _getTypeColor(education.type).withValues(alpha: 0.1),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Text(
                    education.type.name.toUpperCase(),
                    style: TextStyle(
                      fontSize: 10,
                      fontWeight: FontWeight.bold,
                      color: _getTypeColor(education.type),
                    ),
                  ),
                ),
                if (education.degreeLevel != null) ...[
                  const SizedBox(width: 8),
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                    decoration: BoxDecoration(
                      color: Colors.blue.withValues(alpha: 0.1),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Text(
                      education.degreeLevel!.name.toUpperCase(),
                      style: const TextStyle(
                        fontSize: 10,
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
        trailing: Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            if (education.visible)
              Icon(
                Icons.visibility,
                size: 16,
                color: Colors.green.shade600,
              ),
            const Icon(Icons.arrow_forward_ios, size: 16),
          ],
        ),
        onTap: () {
          context.push('/education/${education.id}');
        },
      ),
    );
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

  String _formatDate(DateTime date) {
    return '${date.month}/${date.year}';
  }
}
