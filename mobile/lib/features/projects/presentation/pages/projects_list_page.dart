import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import '../../domain/entities/project_entity.dart';
import '../../../frameworks/domain/entities/framework_entity.dart';
import '../bloc/projects_bloc.dart';

class ProjectsListPage extends StatefulWidget {
  const ProjectsListPage({super.key});

  @override
  State<ProjectsListPage> createState() => _ProjectsListPageState();
}

class _ProjectsListPageState extends State<ProjectsListPage> {
  final _searchController = TextEditingController();
  bool _showFeaturedOnly = false;
  bool _showPublicOnly = false;
  bool _isGridView = true;
  final List<String> _selectedTechnologies = [];
  List<FrameworkEntity> _availableFrameworks = [];
  bool _isLoadingFrameworks = false;

  @override
  void initState() {
    super.initState();
    // Load projects when page initializes
    _loadProjects();
    // Load available frameworks for filtering
    _loadFrameworks();
  }

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  void _loadProjects() {
    context.read<ProjectsBloc>().add(GetProjectsRequested(
      search: _searchController.text.trim().isEmpty ? null : _searchController.text.trim(),
      featured: _showFeaturedOnly,
      public: _showPublicOnly,
      technologies: _selectedTechnologies.isEmpty ? null : _selectedTechnologies,
    ));
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Projects'),
        backgroundColor: Colors.blue.shade600,
        foregroundColor: Colors.white,
        elevation: 0,
        actions: [
          IconButton(
            onPressed: () => _showFilterDialog(context),
            icon: const Icon(Icons.filter_list),
          ),
          IconButton(
            onPressed: () {
              setState(() {
                _isGridView = !_isGridView;
              });
            },
            icon: Icon(_isGridView ? Icons.list : Icons.grid_view),
          ),
        ],
      ),
      body: Column(
        children: [
          // Search bar
          Padding(
            padding: const EdgeInsets.all(16.0),
            child: TextField(
              controller: _searchController,
              decoration: InputDecoration(
                hintText: 'Search projects...',
                prefixIcon: const Icon(Icons.search),
                suffixIcon: _searchController.text.isNotEmpty
                    ? IconButton(
                        onPressed: () {
                          _searchController.clear();
                          _performSearch();
                        },
                        icon: const Icon(Icons.clear),
                      )
                    : null,
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(12),
                ),
              ),
              onChanged: (value) => _performSearch(),
            ),
          ),

          // Filter chips
          if (_showFeaturedOnly || _showPublicOnly || _selectedTechnologies.isNotEmpty)
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 16),
              height: 50,
              child: ListView(
                scrollDirection: Axis.horizontal,
                children: [
                  if (_showFeaturedOnly)
                    Chip(
                      label: const Text('Featured'),
                      onDeleted: () => _removeFilter('featured'),
                      deleteIcon: const Icon(Icons.close, size: 18),
                    ),
                  if (_showPublicOnly)
                    Chip(
                      label: const Text('Public'),
                      onDeleted: () => _removeFilter('public'),
                      deleteIcon: const Icon(Icons.close, size: 18),
                    ),
                  ..._selectedTechnologies.map((tech) => Chip(
                    label: Text(tech),
                    onDeleted: () => _removeTechnology(tech),
                    deleteIcon: const Icon(Icons.close, size: 18),
                  )),
                ],
              ),
            ),

          // Projects list
          Expanded(
            child: BlocBuilder<ProjectsBloc, ProjectsState>(
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
                          'Error loading projects',
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
                            context.read<ProjectsBloc>().add(GetProjectsRequested());
                          },
                          child: const Text('Retry'),
                        ),
                      ],
                    ),
                  );
                } else if (state is ProjectsLoaded) {
                  if (state.projects.isEmpty) {
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
                          const Text(
                            'No projects found',
                            style: TextStyle(
                              fontSize: 18,
                              fontWeight: FontWeight.w500,
                            ),
                          ),
                          const SizedBox(height: 8),
                          Text(
                            'Create your first project to showcase your work!',
                            style: TextStyle(
                              color: Colors.grey.shade600,
                            ),
                            textAlign: TextAlign.center,
                          ),
                        ],
                      ),
                    );
                  }

                  return RefreshIndicator(
                    onRefresh: () async {
                      context.read<ProjectsBloc>().add(GetProjectsRequested());
                    },
                    child: _isGridView
                        ? _buildGridView(state.projects)
                        : _buildListView(state.projects),
                  );
                }

                return const SizedBox.shrink();
              },
            ),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => context.push('/projects/create'),
        backgroundColor: Colors.blue.shade600,
        child: const Icon(Icons.add, color: Colors.white),
      ),
    );
  }

  Widget _buildGridView(List<ProjectEntity> projects) {
    return GridView.builder(
      padding: const EdgeInsets.all(16),
      gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: 2,
        childAspectRatio: 0.8,
        crossAxisSpacing: 16,
        mainAxisSpacing: 16,
      ),
      itemCount: projects.length,
      itemBuilder: (context, index) {
        final project = projects[index];
        return _buildProjectCard(context, project);
      },
    );
  }

  Widget _buildListView(List<ProjectEntity> projects) {
    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: projects.length,
      itemBuilder: (context, index) {
        final project = projects[index];
        return _buildProjectListItem(context, project);
      },
    );
  }

  Widget _buildProjectCard(BuildContext context, ProjectEntity project) {
    return Card(
      clipBehavior: Clip.antiAlias,
      child: InkWell(
        onTap: () async {
          await context.push('/projects/${project.id}');
          // Reload list when returning from detail page
          _loadProjects();
        },
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Project image
            Expanded(
              flex: 3,
              child: Stack(
                children: [
                  Container(
                    width: double.infinity,
                    decoration: BoxDecoration(
                      image: DecorationImage(
                        image: NetworkImage(project.image),
                        fit: BoxFit.cover,
                      ),
                    ),
                  ),
                  if (project.featured)
                    Positioned(
                      top: 8,
                      right: 8,
                      child: Container(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 6,
                          vertical: 2,
                        ),
                        decoration: BoxDecoration(
                          color: Colors.amber.shade100,
                          borderRadius: BorderRadius.circular(8),
                        ),
                        child: Text(
                          'Featured',
                          style: TextStyle(
                            fontSize: 10,
                            color: Colors.amber.shade800,
                            fontWeight: FontWeight.w500,
                          ),
                        ),
                      ),
                    ),
                ],
              ),
            ),
            // Project info
            Expanded(
              flex: 2,
              child: Padding(
                padding: const EdgeInsets.all(12),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      project.title,
                      style: const TextStyle(
                        fontSize: 14,
                        fontWeight: FontWeight.bold,
                      ),
                      maxLines: 2,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 4),
                    Text(
                      project.description,
                      style: TextStyle(
                        fontSize: 12,
                        color: Colors.grey.shade600,
                      ),
                      maxLines: 2,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const Spacer(),
                    Row(
                      children: [
                        Icon(
                          Icons.visibility,
                          size: 12,
                          color: project.visible ? Colors.green : Colors.grey,
                        ),
                        const SizedBox(width: 4),
                        Text(
                          project.visible ? 'Visible' : 'Hidden',
                          style: TextStyle(
                            fontSize: 10,
                            color: project.visible ? Colors.green : Colors.grey,
                          ),
                        ),
                        const Spacer(),
                        if (project.frameworks != null && project.frameworks!.isNotEmpty)
                          Icon(
                            Icons.code,
                            size: 12,
                            color: Colors.blue.shade600,
                          ),
                      ],
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildProjectListItem(BuildContext context, ProjectEntity project) {
    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: InkWell(
        onTap: () async {
          await context.push('/projects/${project.id}');
          // Reload list when returning from detail page
          _loadProjects();
        },
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Row(
            children: [
              // Project image
              ClipRRect(
                borderRadius: BorderRadius.circular(8),
                child: Image.network(
                  project.image,
                  width: 80,
                  height: 80,
                  fit: BoxFit.cover,
                ),
              ),
              const SizedBox(width: 16),
              // Project info
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      children: [
                        Expanded(
                          child: Text(
                            project.title,
                            style: const TextStyle(
                              fontSize: 16,
                              fontWeight: FontWeight.bold,
                            ),
                          ),
                        ),
                        if (project.featured)
                          Container(
                            padding: const EdgeInsets.symmetric(
                              horizontal: 6,
                              vertical: 2,
                            ),
                            decoration: BoxDecoration(
                              color: Colors.amber.shade100,
                              borderRadius: BorderRadius.circular(8),
                            ),
                            child: Text(
                              'Featured',
                              style: TextStyle(
                                fontSize: 10,
                                color: Colors.amber.shade800,
                                fontWeight: FontWeight.w500,
                              ),
                            ),
                          ),
                      ],
                    ),
                    const SizedBox(height: 4),
                    Text(
                      project.description,
                      style: TextStyle(
                        color: Colors.grey.shade600,
                        fontSize: 14,
                      ),
                      maxLines: 2,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 8),
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
                        if (project.frameworks != null && project.frameworks!.isNotEmpty)
                          Row(
                            children: [
                              Icon(
                                Icons.code,
                                size: 16,
                                color: Colors.blue.shade600,
                              ),
                              const SizedBox(width: 4),
                              Text(
                                '${project.frameworks!.length} tech',
                                style: TextStyle(
                                  fontSize: 12,
                                  color: Colors.blue.shade600,
                                ),
                              ),
                            ],
                          ),
                      ],
                    ),
                  ],
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  void _performSearch() {
    final searchQuery = _searchController.text.trim();
    context.read<ProjectsBloc>().add(GetProjectsRequested(
      search: searchQuery.isEmpty ? null : searchQuery,
      featured: _showFeaturedOnly,
      public: _showPublicOnly,
      technologies: _selectedTechnologies.isEmpty ? null : _selectedTechnologies,
    ));
  }

  void _removeFilter(String filterType) {
    setState(() {
      if (filterType == 'featured') {
        _showFeaturedOnly = false;
      } else if (filterType == 'public') {
        _showPublicOnly = false;
      }
    });
    _performSearch();
  }

  void _removeTechnology(String technology) {
    setState(() {
      _selectedTechnologies.remove(technology);
    });
    _performSearch();
  }

  Future<void> _loadFrameworks() async {
    setState(() {
      _isLoadingFrameworks = true;
    });

    try {
      // For now, we'll use a simple list of common technologies
      // In a real implementation, this would come from the frameworks repository
      _availableFrameworks = [
        FrameworkEntity(
          id: '1',
          userId: 'user',
          name: 'Flutter',
          slug: 'flutter',
          description: 'Flutter framework',
          type: FrameworkType.framework,
          visible: true,
          public: true,
          order: 1,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        ),
        FrameworkEntity(
          id: '2',
          userId: 'user',
          name: 'React',
          slug: 'react',
          description: 'React framework',
          type: FrameworkType.framework,
          visible: true,
          public: true,
          order: 2,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        ),
        FrameworkEntity(
          id: '3',
          userId: 'user',
          name: 'Vue.js',
          slug: 'vuejs',
          description: 'Vue.js framework',
          type: FrameworkType.framework,
          visible: true,
          public: true,
          order: 3,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        ),
        FrameworkEntity(
          id: '4',
          userId: 'user',
          name: 'Node.js',
          slug: 'nodejs',
          description: 'Node.js runtime',
          type: FrameworkType.framework,
          visible: true,
          public: true,
          order: 4,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        ),
        FrameworkEntity(
          id: '5',
          userId: 'user',
          name: 'Python',
          slug: 'python',
          description: 'Python programming language',
          type: FrameworkType.language,
          visible: true,
          public: true,
          order: 5,
          createdAt: DateTime.now(),
          updatedAt: DateTime.now(),
        ),
      ];
    } catch (e) {
      // Handle error silently for now
    } finally {
      setState(() {
        _isLoadingFrameworks = false;
      });
    }
  }

  void _showFilterDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Filter Projects'),
        content: StatefulBuilder(
          builder: (context, setDialogState) {
            return SingleChildScrollView(
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  CheckboxListTile(
                    title: const Text('Featured only'),
                    value: _showFeaturedOnly,
                    onChanged: (value) {
                      setDialogState(() {
                        _showFeaturedOnly = value ?? false;
                      });
                    },
                  ),
                  CheckboxListTile(
                    title: const Text('Public only'),
                    value: _showPublicOnly,
                    onChanged: (value) {
                      setDialogState(() {
                        _showPublicOnly = value ?? false;
                      });
                    },
                  ),
                  const Divider(),
                  const Text('Technologies:', style: TextStyle(fontWeight: FontWeight.bold)),
                  const SizedBox(height: 8),
                  if (_isLoadingFrameworks)
                    const CircularProgressIndicator()
                  else
                    Wrap(
                      spacing: 8,
                      runSpacing: 8,
                      children: _availableFrameworks.map((framework) {
                        final isSelected = _selectedTechnologies.contains(framework.name);
                        return FilterChip(
                          label: Text(framework.name),
                          selected: isSelected,
                          onSelected: (selected) {
                            setDialogState(() {
                              if (selected) {
                                _selectedTechnologies.add(framework.name);
                              } else {
                                _selectedTechnologies.remove(framework.name);
                              }
                            });
                          },
                        );
                      }).toList(),
                    ),
                ],
              ),
            );
          },
        ),
        actions: [
          TextButton(
            onPressed: () {
              setState(() {
                _showFeaturedOnly = false;
                _showPublicOnly = false;
                _selectedTechnologies.clear();
              });
              _performSearch();
              Navigator.pop(context);
            },
            child: const Text('Clear'),
          ),
          ElevatedButton(
            onPressed: () {
              setState(() {});
              _performSearch();
              Navigator.pop(context);
            },
            child: const Text('Apply'),
          ),
        ],
      ),
    );
  }
}
