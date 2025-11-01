import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import '../bloc/frameworks_bloc.dart';
import '../../domain/entities/framework_entity.dart';

class FrameworksListPage extends StatefulWidget {
  const FrameworksListPage({super.key});

  @override
  State<FrameworksListPage> createState() => _FrameworksListPageState();
}

class _FrameworksListPageState extends State<FrameworksListPage> {
  final TextEditingController _searchController = TextEditingController();
  String? _selectedType;
  bool? _selectedVisible;
  bool? _selectedPublic;
  String _sortBy = 'order';
  String _sortOrder = 'asc';
  bool _isGridView = true;

  @override
  void initState() {
    super.initState();
    _loadFrameworks();
  }

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  void _loadFrameworks() {
    context.read<FrameworksBloc>().add(LoadFrameworks(
      page: 1,
      limit: 20,
      search: _searchController.text.isEmpty ? null : _searchController.text,
      type: _selectedType,
      visible: _selectedVisible,
      public: _selectedPublic,
      sortBy: _sortBy,
      sortOrder: _sortOrder,
    ));
  }

  void _applyFilters() {
    context.read<FrameworksBloc>().add(LoadFrameworks(
      page: 1,
      limit: 20,
      search: _searchController.text.isEmpty ? null : _searchController.text,
      type: _selectedType,
      visible: _selectedVisible,
      public: _selectedPublic,
      sortBy: _sortBy,
      sortOrder: _sortOrder,
    ));
  }

  void _clearFilters() {
    setState(() {
      _searchController.clear();
      _selectedType = null;
      _selectedVisible = null;
      _selectedPublic = null;
      _sortBy = 'order';
      _sortOrder = 'asc';
    });
    context.read<FrameworksBloc>().add(const LoadFrameworks());
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Frameworks & Technologies'),
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
              context.push('/frameworks/create');
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
                    hintText: 'Search frameworks...',
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
                            _selectedType = selected ? 'language' : null;
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

                      // Public Filter
                      FilterChip(
                        label: Text(_selectedPublic == true ? 'Public' : 'All'),
                        selected: _selectedPublic == true,
                        onSelected: (selected) {
                          setState(() {
                            _selectedPublic = selected ? true : null;
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
                            value: 'name',
                            child: Text('Name'),
                          ),
                          const PopupMenuItem(
                            value: 'type',
                            child: Text('Type'),
                          ),
                          const PopupMenuItem(
                            value: 'createdAt',
                            child: Text('Created Date'),
                          ),
                        ],
                      ),
                      const SizedBox(width: 8),

                      // Clear Filters
                      if (_selectedType != null ||
                          _selectedVisible != null ||
                          _selectedPublic != null ||
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

          // Frameworks List
          Expanded(
            child: BlocBuilder<FrameworksBloc, FrameworksState>(
              builder: (context, state) {
                if (state is FrameworksLoading) {
                  return const Center(child: CircularProgressIndicator());
                }

                if (state is FrameworksError) {
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
                            context.read<FrameworksBloc>().add(const RefreshFrameworks());
                          },
                          child: const Text('Retry'),
                        ),
                      ],
                    ),
                  );
                }

                if (state is FrameworksLoaded) {
                  if (state.frameworks.isEmpty) {
                    return Center(
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Icon(
                            Icons.code_off,
                            size: 64,
                            color: Colors.grey.shade400,
                          ),
                          const SizedBox(height: 16),
                          const Text(
                            'No frameworks found',
                            style: TextStyle(fontSize: 18, fontWeight: FontWeight.w500),
                          ),
                          const SizedBox(height: 8),
                          const Text(
                            'Add your first framework or technology',
                            style: TextStyle(color: Colors.grey),
                          ),
                          const SizedBox(height: 24),
                          ElevatedButton.icon(
                            onPressed: () {
                              context.push('/frameworks/create');
                            },
                            icon: const Icon(Icons.add),
                            label: const Text('Add Framework'),
                          ),
                        ],
                      ),
                    );
                  }

                  return RefreshIndicator(
                    onRefresh: () async {
                      context.read<FrameworksBloc>().add(const RefreshFrameworks());
                    },
                    child: _isGridView
                        ? _buildGridView(state.frameworks)
                        : _buildListView(state.frameworks),
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

  Widget _buildGridView(List<FrameworkEntity> frameworks) {
    return GridView.builder(
      padding: const EdgeInsets.all(16),
      gridDelegate: const SliverGridDelegateWithFixedCrossAxisCount(
        crossAxisCount: 2,
        childAspectRatio: 1.2,
        crossAxisSpacing: 16,
        mainAxisSpacing: 16,
      ),
      itemCount: frameworks.length,
      itemBuilder: (context, index) {
        final framework = frameworks[index];
        return _buildFrameworkCard(framework);
      },
    );
  }

  Widget _buildListView(List<FrameworkEntity> frameworks) {
    return ListView.builder(
      padding: const EdgeInsets.all(16),
      itemCount: frameworks.length,
      itemBuilder: (context, index) {
        final framework = frameworks[index];
        return _buildFrameworkListItem(framework);
      },
    );
  }

  Widget _buildFrameworkCard(FrameworkEntity framework) {
    return Card(
      child: InkWell(
        onTap: () async {
          await context.push('/frameworks/${framework.id}');
          // Reload list when returning from detail page
          _loadFrameworks();
        },
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Icon and Type
              Row(
                children: [
                  Container(
                    width: 40,
                    height: 40,
                    decoration: BoxDecoration(
                      color: framework.color != null
                          ? Color(int.parse(framework.color!.replaceFirst('#', '0xFF')))
                          : Colors.blue.shade100,
                      borderRadius: BorderRadius.circular(8),
                    ),
                    child: framework.icon != null
                        ? Icon(
                            Icons.code,
                            color: Colors.white,
                            size: 20,
                          )
                        : Icon(
                            Icons.code,
                            color: Colors.blue.shade700,
                            size: 20,
                          ),
                  ),
                  const Spacer(),
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                    decoration: BoxDecoration(
                      color: _getTypeColor(framework.type).withValues(alpha: 0.1),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Text(
                      framework.type.name.toUpperCase(),
                      style: TextStyle(
                        fontSize: 10,
                        fontWeight: FontWeight.bold,
                        color: _getTypeColor(framework.type),
                      ),
                    ),
                  ),
                ],
              ),
              const SizedBox(height: 12),

              // Name
              Text(
                framework.name,
                style: const TextStyle(
                  fontSize: 16,
                  fontWeight: FontWeight.w600,
                ),
                maxLines: 2,
                overflow: TextOverflow.ellipsis,
              ),
              const SizedBox(height: 4),

              // Description
              if (framework.description != null)
                Text(
                  framework.description!,
                  style: const TextStyle(
                    fontSize: 12,
                    color: Colors.grey,
                  ),
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                ),

              const Spacer(),

              // Status indicators
              Row(
                children: [
                  if (framework.visible)
                    Icon(
                      Icons.visibility,
                      size: 16,
                      color: Colors.green.shade600,
                    ),
                  if (framework.public)
                    Icon(
                      Icons.public,
                      size: 16,
                      color: Colors.blue.shade600,
                    ),
                  const Spacer(),
                  if (framework.proficiencyLevel != null)
                    Container(
                      padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                      decoration: BoxDecoration(
                        color: _getProficiencyColor(framework.proficiencyLevel!).withValues(alpha: 0.1),
                        borderRadius: BorderRadius.circular(8),
                      ),
                      child: Text(
                        framework.proficiencyLevel!.name.toUpperCase(),
                        style: TextStyle(
                          fontSize: 8,
                          fontWeight: FontWeight.bold,
                          color: _getProficiencyColor(framework.proficiencyLevel!),
                        ),
                      ),
                    ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildFrameworkListItem(FrameworkEntity framework) {
    return Card(
      margin: const EdgeInsets.only(bottom: 8),
      child: ListTile(
        leading: Container(
          width: 48,
          height: 48,
          decoration: BoxDecoration(
            color: framework.color != null
                ? Color(int.parse(framework.color!.replaceFirst('#', '0xFF')))
                : Colors.blue.shade100,
            borderRadius: BorderRadius.circular(8),
          ),
          child: framework.icon != null
              ? Icon(
                  Icons.code,
                  color: Colors.white,
                  size: 24,
                )
              : Icon(
                  Icons.code,
                  color: Colors.blue.shade700,
                  size: 24,
                ),
        ),
        title: Text(
          framework.name,
          style: const TextStyle(fontWeight: FontWeight.w600),
        ),
        subtitle: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            if (framework.description != null)
              Text(
                framework.description!,
                maxLines: 1,
                overflow: TextOverflow.ellipsis,
              ),
            const SizedBox(height: 4),
            Row(
              children: [
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                  decoration: BoxDecoration(
                    color: _getTypeColor(framework.type).withValues(alpha: 0.1),
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Text(
                    framework.type.name.toUpperCase(),
                    style: TextStyle(
                      fontSize: 10,
                      fontWeight: FontWeight.bold,
                      color: _getTypeColor(framework.type),
                    ),
                  ),
                ),
                if (framework.proficiencyLevel != null) ...[
                  const SizedBox(width: 8),
                  Container(
                    padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                    decoration: BoxDecoration(
                      color: _getProficiencyColor(framework.proficiencyLevel!).withValues(alpha: 0.1),
                      borderRadius: BorderRadius.circular(12),
                    ),
                    child: Text(
                      framework.proficiencyLevel!.name.toUpperCase(),
                      style: TextStyle(
                        fontSize: 10,
                        fontWeight: FontWeight.bold,
                        color: _getProficiencyColor(framework.proficiencyLevel!),
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
            if (framework.visible)
              Icon(
                Icons.visibility,
                size: 16,
                color: Colors.green.shade600,
              ),
            if (framework.public)
              Icon(
                Icons.public,
                size: 16,
                color: Colors.blue.shade600,
              ),
            const Icon(Icons.arrow_forward_ios, size: 16),
          ],
        ),
        onTap: () async {
          await context.push('/frameworks/${framework.id}');
          // Reload list when returning from detail page
          _loadFrameworks();
        },
      ),
    );
  }

  Color _getTypeColor(FrameworkType type) {
    switch (type) {
      case FrameworkType.language:
        return Colors.blue;
      case FrameworkType.framework:
        return Colors.green;
      case FrameworkType.library:
        return Colors.orange;
      case FrameworkType.tool:
        return Colors.purple;
      case FrameworkType.database:
        return Colors.red;
      case FrameworkType.other:
        return Colors.grey;
    }
  }

  Color _getProficiencyColor(ProficiencyLevel level) {
    switch (level) {
      case ProficiencyLevel.beginner:
        return Colors.green;
      case ProficiencyLevel.intermediate:
        return Colors.orange;
      case ProficiencyLevel.advanced:
        return Colors.red;
      case ProficiencyLevel.expert:
        return Colors.purple;
    }
  }
}
