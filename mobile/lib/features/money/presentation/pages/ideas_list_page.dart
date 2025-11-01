import 'package:flutter/material.dart';
import 'package:flutter_bloc/flutter_bloc.dart';
import 'package:go_router/go_router.dart';
import '../../domain/entities/idea_entity.dart';
import '../bloc/money_bloc.dart';

class IdeasListPage extends StatefulWidget {
  const IdeasListPage({super.key});

  @override
  State<IdeasListPage> createState() => _IdeasListPageState();
}

class _IdeasListPageState extends State<IdeasListPage> {
  final _searchController = TextEditingController();
  bool _showFeaturedOnly = false;
  bool _showPublicOnly = false;

  @override
  void initState() {
    super.initState();
    // Load ideas when page initializes
    _loadIdeas();
  }

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  void _loadIdeas() {
    final searchQuery = _searchController.text.trim();
    context.read<MoneyBloc>().add(GetIdeasRequested(
      search: searchQuery.isEmpty ? null : searchQuery,
      featured: _showFeaturedOnly,
      public: _showPublicOnly,
    ));
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: const Text('Money Ideas'),
        backgroundColor: Colors.green.shade600,
        foregroundColor: Colors.white,
        elevation: 0,
        actions: [
          IconButton(
            onPressed: () => _showFilterDialog(context),
            icon: const Icon(Icons.filter_list),
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
                hintText: 'Search ideas...',
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
          if (_showFeaturedOnly || _showPublicOnly)
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
                ],
              ),
            ),

          // Ideas list
          Expanded(
            child: BlocBuilder<MoneyBloc, MoneyState>(
              builder: (context, state) {
                // Handle composite state
                if (state is MoneyDataState) {
                  final dataState = state;
                  
                  // Show loading indicator if loading and no cached data
                  if (dataState.isLoading && dataState.ideas == null) {
                    return const Center(child: CircularProgressIndicator());
                  }
                  
                  // Show error if there's an error and no cached data
                  if (dataState.error != null && dataState.ideas == null) {
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
                            'Error loading ideas',
                            style: TextStyle(
                              fontSize: 18,
                              color: Colors.red.shade700,
                            ),
                          ),
                          const SizedBox(height: 8),
                          Text(
                            dataState.error!,
                            style: TextStyle(color: Colors.red.shade600),
                            textAlign: TextAlign.center,
                          ),
                          const SizedBox(height: 16),
                          ElevatedButton(
                            onPressed: () {
                              context.read<MoneyBloc>().add(GetIdeasRequested());
                            },
                            child: const Text('Retry'),
                          ),
                        ],
                      ),
                    );
                  }
                  
                  // Show cached data if available
                  if (dataState.ideas != null) {
                    final ideas = dataState.ideas!;
                    
                    // Show loading overlay if still loading
                    if (dataState.isLoading) {
                      return Stack(
                        children: [
                          _buildIdeasList(ideas),
                          Container(
                            color: Colors.black.withValues(alpha: 0.3),
                            child: const Center(
                              child: CircularProgressIndicator(),
                            ),
                          ),
                        ],
                      );
                    }
                    
                    return _buildIdeasList(ideas);
                  }
                }
                
                // Legacy state handling
                if (state is MoneyLoading) {
                  return const Center(child: CircularProgressIndicator());
                } else if (state is MoneyError) {
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
                          'Error loading ideas',
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
                            context.read<MoneyBloc>().add(GetIdeasRequested());
                          },
                          child: const Text('Retry'),
                        ),
                      ],
                    ),
                  );
                } else if (state is IdeasLoaded) {
                  if (state.ideas.isEmpty) {
                    return Center(
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Icon(
                            Icons.lightbulb_outline,
                            size: 64,
                            color: Colors.grey.shade400,
                          ),
                          const SizedBox(height: 16),
                          const Text(
                            'No ideas found',
                            style: TextStyle(
                              fontSize: 18,
                              fontWeight: FontWeight.w500,
                            ),
                          ),
                          const SizedBox(height: 8),
                          Text(
                            'Create your first money-making idea!',
                            style: TextStyle(
                              color: Colors.grey.shade600,
                            ),
                          ),
                        ],
                      ),
                    );
                  }

                  return RefreshIndicator(
                    onRefresh: () async {
                      context.read<MoneyBloc>().add(GetIdeasRequested());
                    },
                    child: ListView.builder(
                      padding: const EdgeInsets.all(16),
                      itemCount: state.ideas.length,
                      itemBuilder: (context, index) {
                        final idea = state.ideas[index];
                        return _buildIdeaCard(context, idea);
                      },
                    ),
                  );
                }

                return const SizedBox.shrink();
              },
            ),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () => context.push('/money/ideas/create'),
        backgroundColor: Colors.green.shade600,
        child: const Icon(Icons.add, color: Colors.white),
      ),
    );
  }

  Widget _buildIdeasList(List<IdeaEntity> ideas) {
    if (ideas.isEmpty) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(
              Icons.lightbulb_outline,
              size: 64,
              color: Colors.grey.shade400,
            ),
            const SizedBox(height: 16),
            const Text(
              'No ideas found',
              style: TextStyle(
                fontSize: 18,
                fontWeight: FontWeight.w500,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              'Create your first money-making idea!',
              style: TextStyle(
                color: Colors.grey.shade600,
              ),
            ),
          ],
        ),
      );
    }

    return RefreshIndicator(
      onRefresh: () async {
        context.read<MoneyBloc>().add(GetIdeasRequested());
      },
      child: ListView.builder(
        padding: const EdgeInsets.all(16),
        itemCount: ideas.length,
        itemBuilder: (context, index) {
          final idea = ideas[index];
          return _buildIdeaCard(context, idea);
        },
      ),
    );
  }

  Widget _buildIdeaCard(BuildContext context, IdeaEntity idea) {
    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: InkWell(
        onTap: () async {
          await context.push('/money/ideas/${idea.id}');
          // Reload list when returning from detail page
          _loadIdeas();
        },
        borderRadius: BorderRadius.circular(12),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  Expanded(
                    child: Text(
                      idea.title,
                      style: const TextStyle(
                        fontSize: 18,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                  ),
                  if (idea.featured)
                    Container(
                      padding: const EdgeInsets.symmetric(
                        horizontal: 8,
                        vertical: 4,
                      ),
                      decoration: BoxDecoration(
                        color: Colors.amber.shade100,
                        borderRadius: BorderRadius.circular(12),
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
                ],
              ),
              if (idea.description != null) ...[
                const SizedBox(height: 8),
                Text(
                  idea.description!,
                  style: TextStyle(
                    color: Colors.grey.shade600,
                    fontSize: 14,
                  ),
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                ),
              ],
              const SizedBox(height: 12),
              Row(
                children: [
                  Icon(
                    Icons.visibility,
                    size: 16,
                    color: idea.visible ? Colors.green : Colors.grey,
                  ),
                  const SizedBox(width: 4),
                  Text(
                    idea.visible ? 'Visible' : 'Hidden',
                    style: TextStyle(
                      fontSize: 12,
                      color: idea.visible ? Colors.green : Colors.grey,
                    ),
                  ),
                  const SizedBox(width: 16),
                  Icon(
                    Icons.public,
                    size: 16,
                    color: idea.public ? Colors.blue : Colors.grey,
                  ),
                  const SizedBox(width: 4),
                  Text(
                    idea.public ? 'Public' : 'Private',
                    style: TextStyle(
                      fontSize: 12,
                      color: idea.public ? Colors.blue : Colors.grey,
                    ),
                  ),
                  const Spacer(),
                  Text(
                    _formatDate(idea.createdAt),
                    style: TextStyle(
                      fontSize: 12,
                      color: Colors.grey.shade500,
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

  void _performSearch() {
    final searchQuery = _searchController.text.trim();
    context.read<MoneyBloc>().add(GetIdeasRequested(
      search: searchQuery.isEmpty ? null : searchQuery,
      featured: _showFeaturedOnly,
      public: _showPublicOnly,
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

  void _showFilterDialog(BuildContext context) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Filter Ideas'),
        content: StatefulBuilder(
          builder: (context, setDialogState) {
            return Column(
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
              ],
            );
          },
        ),
        actions: [
          TextButton(
            onPressed: () {
              setState(() {
                _showFeaturedOnly = false;
                _showPublicOnly = false;
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
}
