import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import '../bloc/experience_bloc.dart';

class ExperienceListPage extends StatefulWidget {
  const ExperienceListPage({super.key});

  @override
  State<ExperienceListPage> createState() => _ExperienceListPageState();
}

class _ExperienceListPageState extends State<ExperienceListPage> {
  final _searchController = TextEditingController();
  final _companyFilterController = TextEditingController();
  
  bool _isGridView = false;
  bool _showOnlyVisible = true;
  String _selectedCompany = 'All Companies';
  List<String> _companies = ['All Companies'];

  @override
  void initState() {
    super.initState();
    // Load experience list
    _loadExperiences();
  }

  @override
  void dispose() {
    _searchController.dispose();
    _companyFilterController.dispose();
    super.dispose();
  }

  void _loadExperiences() {
    context.read<ExperienceBloc>().add(GetExperienceListRequested(
      search: _searchController.text.isEmpty ? null : _searchController.text,
      visible: _showOnlyVisible ? true : null,
      company: _selectedCompany == 'All Companies' ? null : _selectedCompany,
    ));
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Experience'),
        backgroundColor: Colors.indigo.shade600,
        foregroundColor: Colors.white,
        elevation: 0,
        actions: [
          IconButton(
            onPressed: () {
              setState(() {
                _isGridView = !_isGridView;
              });
            },
            icon: Icon(_isGridView ? Icons.list : Icons.grid_view),
          ),
          IconButton(
            onPressed: () => context.push('/experience/create'),
            icon: const Icon(Icons.add),
          ),
        ],
      ),
      body: Column(
        children: [
          // Search and Filter Section
          Container(
            padding: const EdgeInsets.all(16),
            color: Colors.grey.shade50,
            child: Column(
              children: [
                // Search Bar
                TextField(
                  controller: _searchController,
                  decoration: InputDecoration(
                    hintText: 'Search experiences...',
                    prefixIcon: const Icon(Icons.search),
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                    filled: true,
                    fillColor: Colors.white,
                  ),
                  onChanged: (value) {
                    context.read<ExperienceBloc>().add(GetExperienceListRequested(
                      search: value.isEmpty ? null : value,
                      visible: _showOnlyVisible ? true : null,
                      company: _selectedCompany == 'All Companies' ? null : _selectedCompany,
                    ));
                  },
                ),
                const SizedBox(height: 12),
                
                // Filters Row
                Row(
                  children: [
                    // Company Filter
                    Expanded(
                      child: DropdownButtonFormField<String>(
                        value: _selectedCompany,
                        decoration: InputDecoration(
                          labelText: 'Company',
                          border: OutlineInputBorder(
                            borderRadius: BorderRadius.circular(8),
                          ),
                          filled: true,
                          fillColor: Colors.white,
                          contentPadding: const EdgeInsets.symmetric(
                            horizontal: 12,
                            vertical: 8,
                          ),
                        ),
                        items: _companies.map((company) {
                          return DropdownMenuItem(
                            value: company,
                            child: Text(
                              company,
                              overflow: TextOverflow.ellipsis,
                            ),
                          );
                        }).toList(),
                        onChanged: (value) {
                          setState(() {
                            _selectedCompany = value ?? 'All Companies';
                          });
                          context.read<ExperienceBloc>().add(GetExperienceListRequested(
                            search: _searchController.text.isEmpty ? null : _searchController.text,
                            visible: _showOnlyVisible ? true : null,
                            company: _selectedCompany == 'All Companies' ? null : _selectedCompany,
                          ));
                        },
                      ),
                    ),
                    const SizedBox(width: 12),
                    
                    // Visibility Filter
                    FilterChip(
                      label: const Text('Visible Only'),
                      selected: _showOnlyVisible,
                      onSelected: (selected) {
                        setState(() {
                          _showOnlyVisible = selected;
                        });
                        context.read<ExperienceBloc>().add(GetExperienceListRequested(
                          search: _searchController.text.isEmpty ? null : _searchController.text,
                          visible: _showOnlyVisible ? true : null,
                          company: _selectedCompany == 'All Companies' ? null : _selectedCompany,
                        ));
                      },
                      selectedColor: Colors.indigo.shade100,
                      checkmarkColor: Colors.indigo.shade600,
                    ),
                  ],
                ),
              ],
            ),
          ),
          
          // Experience List
          Expanded(
            child: BlocConsumer<ExperienceBloc, ExperienceState>(
              listener: (context, state) {
                if (state is ExperienceError) {
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(
                      content: Text(state.message),
                      backgroundColor: Colors.red,
                    ),
                  );
                } else if (state is ExperienceListLoaded) {
                  // Update companies list for filter
                  final companies = state.experiences
                      .map((exp) => exp.company)
                      .toSet()
                      .toList()
                    ..sort();
                  _companies = ['All Companies', ...companies];
                  if (!_companies.contains(_selectedCompany)) {
                    _selectedCompany = 'All Companies';
                  }
                }
              },
              builder: (context, state) {
                if (state is ExperienceLoading) {
                  return const Center(child: CircularProgressIndicator());
                } else if (state is ExperienceListLoaded) {
                  if (state.experiences.isEmpty) {
                    return _buildEmptyState();
                  }
                  
                  return RefreshIndicator(
                    onRefresh: () async {
                      context.read<ExperienceBloc>().add(GetExperienceListRequested(
                        search: _searchController.text.isEmpty ? null : _searchController.text,
                        visible: _showOnlyVisible ? true : null,
                        company: _selectedCompany == 'All Companies' ? null : _selectedCompany,
                      ));
                    },
                    child: _isGridView
                        ? _buildGridView(state.experiences)
                        : _buildListView(state.experiences),
                  );
                } else if (state is ExperienceError) {
                  return _buildErrorState(state.message);
                }
                
                return const SizedBox.shrink();
              },
            ),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => context.go('/experience/create'),
        backgroundColor: Colors.indigo.shade600,
        child: const Icon(Icons.add, color: Colors.white),
      ),
    );
  }

  Widget _buildEmptyState() {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(
            Icons.work_outline,
            size: 64,
            color: Colors.grey.shade400,
          ),
          const SizedBox(height: 16),
          Text(
            'No experiences found',
            style: TextStyle(
              fontSize: 18,
              color: Colors.grey.shade600,
            ),
          ),
          const SizedBox(height: 8),
          Text(
            'Add your first work experience to get started',
            style: TextStyle(
              color: Colors.grey.shade500,
            ),
          ),
          const SizedBox(height: 24),
          ElevatedButton.icon(
            onPressed: () => context.push('/experience/create'),
            icon: const Icon(Icons.add),
            label: const Text('Add Experience'),
            style: ElevatedButton.styleFrom(
              backgroundColor: Colors.indigo.shade600,
              foregroundColor: Colors.white,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildErrorState(String message) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(
            Icons.error_outline,
            size: 64,
            color: Colors.red.shade400,
          ),
          const SizedBox(height: 16),
          Text(
            'Something went wrong',
            style: TextStyle(
              fontSize: 18,
              color: Colors.grey.shade600,
            ),
          ),
          const SizedBox(height: 8),
          Text(
            message,
            style: TextStyle(
              color: Colors.grey.shade500,
            ),
            textAlign: TextAlign.center,
          ),
          const SizedBox(height: 24),
          ElevatedButton.icon(
            onPressed: () {
              context.read<ExperienceBloc>().add(GetExperienceListRequested());
            },
            icon: const Icon(Icons.refresh),
            label: const Text('Try Again'),
            style: ElevatedButton.styleFrom(
              backgroundColor: Colors.indigo.shade600,
              foregroundColor: Colors.white,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildListView(List<dynamic> experiences) {
    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: experiences.length,
      itemBuilder: (context, index) {
        final experience = experiences[index];
        return Card(
          margin: const EdgeInsets.only(bottom: 12),
          child: ListTile(
            leading: CircleAvatar(
              backgroundColor: Colors.indigo.shade100,
              child: experience.icon.isNotEmpty
                  ? Image.network(
                      experience.icon,
                      errorBuilder: (context, error, stackTrace) =>
                          Icon(Icons.work, color: Colors.indigo.shade600),
                    )
                  : Icon(Icons.work, color: Colors.indigo.shade600),
            ),
            title: Text(
              experience.title,
              style: const TextStyle(
                fontWeight: FontWeight.bold,
              ),
            ),
            subtitle: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  '${experience.company} • ${experience.period}',
                  style: TextStyle(
                    color: Colors.grey.shade600,
                  ),
                ),
                if (experience.location.isNotEmpty)
                  Text(
                    experience.location,
                    style: TextStyle(
                      color: Colors.grey.shade500,
                      fontSize: 12,
                    ),
                  ),
                if (experience.technologies.isNotEmpty)
                  Wrap(
                    spacing: 4,
                    children: experience.technologies.take(3).map((tech) {
                      return Chip(
                        label: Text(
                          tech,
                          style: const TextStyle(fontSize: 10),
                        ),
                        backgroundColor: Colors.indigo.shade50,
                        materialTapTargetSize: MaterialTapTargetSize.shrinkWrap,
                      );
                    }).toList(),
                  ),
              ],
            ),
            trailing: Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                if (!experience.visible)
                  Icon(
                    Icons.visibility_off,
                    color: Colors.grey.shade400,
                    size: 16,
                  ),
                const SizedBox(width: 8),
                const Icon(Icons.arrow_forward_ios, size: 16),
              ],
            ),
            onTap: () async {
              await context.push('/experience/${experience.id}');
              // Reload list when returning from detail page
              _loadExperiences();
            },
          ),
        );
      },
    );
  }

  Widget _buildGridView(List<dynamic> experiences) {
    return GridView.builder(
      padding: const EdgeInsets.all(16),
      gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: 2,
        childAspectRatio: 0.8,
        crossAxisSpacing: 12,
        mainAxisSpacing: 12,
      ),
      itemCount: experiences.length,
      itemBuilder: (context, index) {
        final experience = experiences[index];
        return Card(
          child: InkWell(
            onTap: () async {
              await context.push('/experience/${experience.id}');
              // Reload list when returning from detail page
              _loadExperiences();
            },
            borderRadius: BorderRadius.circular(12),
            child: Padding(
              padding: const EdgeInsets.all(12),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  // Icon and visibility indicator
                  Row(
                    children: [
                      CircleAvatar(
                        radius: 20,
                        backgroundColor: Colors.indigo.shade100,
                        child: experience.icon.isNotEmpty
                            ? Image.network(
                                experience.icon,
                                errorBuilder: (context, error, stackTrace) =>
                                    Icon(Icons.work, color: Colors.indigo.shade600),
                              )
                            : Icon(Icons.work, color: Colors.indigo.shade600),
                      ),
                      const Spacer(),
                      if (!experience.visible)
                        Icon(
                          Icons.visibility_off,
                          color: Colors.grey.shade400,
                          size: 16,
                        ),
                    ],
                  ),
                  const SizedBox(height: 12),
                  
                  // Title
                  Text(
                    experience.title,
                    style: const TextStyle(
                      fontWeight: FontWeight.bold,
                      fontSize: 14,
                    ),
                    maxLines: 2,
                    overflow: TextOverflow.ellipsis,
                  ),
                  const SizedBox(height: 4),
                  
                  // Company
                  Text(
                    experience.company,
                    style: TextStyle(
                      color: Colors.grey.shade600,
                      fontSize: 12,
                    ),
                    maxLines: 1,
                    overflow: TextOverflow.ellipsis,
                  ),
                  const SizedBox(height: 4),
                  
                  // Period
                  Text(
                    experience.period,
                    style: TextStyle(
                      color: Colors.grey.shade500,
                      fontSize: 11,
                    ),
                  ),
                  const Spacer(),
                  
                  // Technologies
                  if (experience.technologies.isNotEmpty)
                    Wrap(
                      spacing: 2,
                      runSpacing: 2,
                      children: experience.technologies.take(2).map((tech) {
                        return Container(
                          padding: const EdgeInsets.symmetric(
                            horizontal: 6,
                            vertical: 2,
                          ),
                          decoration: BoxDecoration(
                            color: Colors.indigo.shade50,
                            borderRadius: BorderRadius.circular(8),
                          ),
                          child: Text(
                            tech,
                            style: TextStyle(
                              fontSize: 9,
                              color: Colors.indigo.shade700,
                            ),
                          ),
                        );
                      }).toList(),
                    ),
                ],
              ),
            ),
          ),
        );
      },
    );
  }
}
